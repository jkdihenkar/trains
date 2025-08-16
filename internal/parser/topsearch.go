package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"trains/internal/types"
)

// ParseTransitRoutes extracts transit routes from HTML
func ParseTransitRoutes(htmlContent string) []types.TransitRoute {
	var routes []types.TransitRoute
	
	// Look for table rows containing route data with multiline matching
	// The pattern matches table rows with station data
	rowPattern := regexp.MustCompile(`(?s)<tr[^>]*>.*?</tr>`)
	rows := rowPattern.FindAllString(htmlContent, -1)
	
	fmt.Printf("Found %d table rows to parse\n", len(rows))
	
	for i, row := range rows {
		// Skip header rows or rows without proper data
		if !strings.Contains(row, "Show") || !strings.Contains(row, "Kms") {
			continue
		}
		
		// Extract table cells
		cellPattern := regexp.MustCompile(`(?s)<td[^>]*>(.*?)</td>`)
		cellMatches := cellPattern.FindAllStringSubmatch(row, -1)
		
		if len(cellMatches) < 6 { // Need at least 6 cells: source, source_count, transit, transit_count, dest, distance
			fmt.Printf("Row %d: Found only %d cells, skipping\n", i, len(cellMatches))
			continue
		}
		
		// Structure: source | source_count | show_link+transit | transit_count | dest | distance
		// Extract source station info (cell 0)
		sourceStationPattern := regexp.MustCompile(`([A-Z\s]+)\s*<br>\s*\(([A-Z]+)\)`)
		sourceMatch := sourceStationPattern.FindStringSubmatch(cellMatches[0][1])
		if len(sourceMatch) < 3 {
			continue
		}
		
		// Extract source train count (cell 1)
		sourceCountStr := strings.TrimSpace(cellMatches[1][1])
		sourceCount, err := strconv.Atoi(sourceCountStr)
		if err != nil {
			continue
		}
		
		// Extract Show link and transit station from cell 2
		linkPattern := regexp.MustCompile(`href="([^"]*)"`)
		linkMatch := linkPattern.FindStringSubmatch(cellMatches[2][1])
		if len(linkMatch) < 2 {
			continue
		}
		
		// Extract transit station from the same cell (after the Show link)
		transitMatch := sourceStationPattern.FindStringSubmatch(cellMatches[2][1])
		if len(transitMatch) < 3 {
			continue
		}
		
		// Extract transit train count (cell 3)
		transitCountStr := strings.TrimSpace(cellMatches[3][1])
		transitCount, err := strconv.Atoi(transitCountStr)
		if err != nil {
			continue
		}
		
		// Extract destination station info (cell 4)
		destMatch := sourceStationPattern.FindStringSubmatch(cellMatches[4][1])
		if len(destMatch) < 3 {
			continue
		}
		
		// Extract distance (cell 5)
		distanceStr := strings.TrimSpace(cellMatches[5][1])
		
		// Build the route
		route := types.TransitRoute{
			SourceStation:      strings.TrimSpace(sourceMatch[1]),
			SourceStationCode:  sourceMatch[2],
			SourceTrainCount:   sourceCount,
			TransitStation:     strings.TrimSpace(transitMatch[1]),
			TransitStationCode: transitMatch[2],
			TransitTrainCount:  transitCount,
			DestStation:        strings.TrimSpace(destMatch[1]),
			DestStationCode:    destMatch[2],
			Distance:           distanceStr,
			ShowLink:           linkMatch[1],
		}
		
		routes = append(routes, route)
		
		// Debug first 3 routes
		if len(routes) <= 3 {
			fmt.Printf("Route %d: %s (%s) [%d] → %s (%s) [%d] → %s (%s) - %s\n",
				len(routes), route.SourceStation, route.SourceStationCode, route.SourceTrainCount,
				route.TransitStation, route.TransitStationCode, route.TransitTrainCount,
				route.DestStation, route.DestStationCode, route.Distance)
		}
	}
	
	return routes
}