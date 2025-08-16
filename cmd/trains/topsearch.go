package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"trains/internal/cache"
	"trains/internal/client"
	"trains/internal/parser"
	"trains/internal/types"
)

// runTopSearch handles the topsearch command
func runTopSearch(cmd *cobra.Command, args []string) error {
	// Get URL flag
	url, err := cmd.Flags().GetString("url")
	if err != nil {
		return fmt.Errorf("error getting url flag: %v", err)
	}
	
	// Get limit flag
	limit, err := cmd.Flags().GetInt("limit")
	if err != nil {
		return fmt.Errorf("error getting limit flag: %v", err)
	}
	
	// Get max-distance flag
	maxDistance, err := cmd.Flags().GetInt("max-distance")
	if err != nil {
		return fmt.Errorf("error getting max-distance flag: %v", err)
	}
	
	// Handle --no-cache flag (check parent flags since it's persistent)
	noCacheFlag, _ := cmd.Parent().Flags().GetBool("no-cache")
	if noCacheFlag {
		cacheEnabled = false
	}
	
	fmt.Printf("🚂 Starting top transit route analysis...\n")
	fmt.Printf("📍 URL: %s\n", url)
	fmt.Printf("💾 Cache: %t\n", cacheEnabled)
	fmt.Printf("📊 Limit: %d routes\n", limit)
	if maxDistance > 0 {
		fmt.Printf("📏 Max Distance: %d km\n", maxDistance)
	}
	fmt.Println()
	
	// Initialize cache directory if caching is enabled
	if cacheEnabled {
		if err := cache.InitCache(); err != nil {
			return fmt.Errorf("error initializing cache: %v", err)
		}
	}
	
	// Fetch the webpage with or without caching
	var htmlContent string
	
	if cacheEnabled {
		htmlContent, err = client.FetchWithCache(url)
	} else {
		fmt.Printf("🌐 Fetching from network (cache disabled): %s\n", url)
		htmlContent, err = client.FetchFromNetwork(url)
	}
	
	if err != nil {
		return fmt.Errorf("error fetching URL: %v", err)
	}
	
	// Check if we need to fetch multiple pages
	var allRoutes []types.TransitRoute
	
	if parser.ShouldFetchAllPages(url) {
		fmt.Printf("🔄 Detecting multi-page transit data, fetching all pages...\n")
		allRoutes, err = fetchAllTransitPages(url, cacheEnabled)
		if err != nil {
			return fmt.Errorf("error fetching all pages: %v", err)
		}
	} else {
		// Single page - parse normally
		allRoutes = parser.ParseTransitRoutes(htmlContent)
	}
	
	fmt.Printf("Found %d transit routes total\n", len(allRoutes))
	
	// Display results
	displayTransitRoutes(allRoutes, limit, maxDistance, url)
	
	return nil
}

// fetchAllTransitPages fetches all pages of transit data
func fetchAllTransitPages(baseURL string, cacheEnabled bool) ([]types.TransitRoute, error) {
	var allRoutes []types.TransitRoute
	pageNum := 1
	
	for {
		// Construct page URL
		separator := "?"
		if strings.Contains(baseURL, "?") {
			separator = "&"
		}
		pageURL := fmt.Sprintf("%s%spage=%d", baseURL, separator, pageNum)
		
		fmt.Printf("📄 Fetching page %d: %s\n", pageNum, pageURL)
		
		// Fetch page content
		var htmlContent string
		var err error
		
		if cacheEnabled {
			htmlContent, err = client.FetchWithCache(pageURL)
		} else {
			fmt.Printf("🌐 Fetching from network (cache disabled): %s\n", pageURL)
			htmlContent, err = client.FetchFromNetwork(pageURL)
		}
		
		if err != nil {
			return nil, fmt.Errorf("error fetching page %d: %v", pageNum, err)
		}
		
		// Parse routes from this page
		pageRoutes := parser.ParseTransitRoutes(htmlContent)
		
		// If no routes found on this page, we've reached the end
		if len(pageRoutes) == 0 {
			fmt.Printf("✅ Reached end of pages at page %d (no routes found)\n", pageNum)
			break
		}
		
		// Add routes from this page
		allRoutes = append(allRoutes, pageRoutes...)
		fmt.Printf("   Found %d routes on page %d (total: %d)\n", len(pageRoutes), pageNum, len(allRoutes))
		
		// Check if there's a "Next" link or pagination indicator
		// If this page has fewer routes than expected, it might be the last page
		if len(pageRoutes) < 10 { // Assuming typical page size is around 10+ routes
			fmt.Printf("✅ Likely reached last page (only %d routes found)\n", len(pageRoutes))
			break
		}
		
		pageNum++
		
		// Safety limit to prevent infinite loops
		if pageNum > 20 { // Reasonable upper limit
			fmt.Printf("⚠️  Reached safety limit of 20 pages\n")
			break
		}
	}
	
	fmt.Printf("🔄 Fetched %d pages with total %d routes\n", pageNum-1, len(allRoutes))
	return allRoutes, nil
}

// displayTransitRoutes displays and sorts transit routes
func displayTransitRoutes(routes []types.TransitRoute, limit int, maxDistance int, originalURL string) {
	fmt.Println("\n=== TOP TRANSIT ROUTES ===\n")
	
	if len(routes) == 0 {
		fmt.Println("No transit routes found.")
		return
	}
	
	// Filter by max distance if specified
	if maxDistance > 0 {
		var filteredRoutes []types.TransitRoute
		for _, route := range routes {
			distance := parser.ParseDistanceKm(route.Distance)
			if distance > 0 && distance <= maxDistance {
				filteredRoutes = append(filteredRoutes, route)
			}
		}
		routes = filteredRoutes
		fmt.Printf("After distance filtering (≤%d km): %d routes\n", maxDistance, len(routes))
		
		if len(routes) == 0 {
			fmt.Printf("No routes found within %d km distance limit.\n", maxDistance)
			return
		}
	}
	
	// Sort routes by distance (lowest first) when using multi-page fetching OR when max-distance filter is specified
	// Otherwise sort by train count for single pages
	shouldSortByDistance := parser.ShouldFetchAllPages(originalURL) || maxDistance > 0
	
	if shouldSortByDistance {
		fmt.Printf("📊 Sorting by distance (shortest routes first)...\n")
		// Sort by distance (ascending)
		for i := 0; i < len(routes)-1; i++ {
			for j := i + 1; j < len(routes); j++ {
				distanceI := parser.ParseDistanceKm(routes[i].Distance)
				distanceJ := parser.ParseDistanceKm(routes[j].Distance)
				if distanceJ < distanceI {
					routes[i], routes[j] = routes[j], routes[i]
				}
			}
		}
	} else {
		fmt.Printf("📊 Sorting by train availability (most trains first)...\n")
		// Sort routes by total train count (source + transit) in descending order
		// This prioritizes routes with more train options
		for i := 0; i < len(routes)-1; i++ {
			for j := i + 1; j < len(routes); j++ {
				totalI := routes[i].SourceTrainCount + routes[i].TransitTrainCount
				totalJ := routes[j].SourceTrainCount + routes[j].TransitTrainCount
				if totalJ > totalI {
					routes[i], routes[j] = routes[j], routes[i]
				}
			}
		}
	}
	
	// Limit results
	displayCount := len(routes)
	if limit > 0 && limit < len(routes) {
		displayCount = limit
		routes = routes[:limit]
	}
	
	if shouldSortByDistance {
		fmt.Printf("Showing top %d routes (sorted by shortest distance):\n\n", displayCount)
	} else {
		fmt.Printf("Showing top %d routes (sorted by total train availability):\n\n", displayCount)
	}
	
	for i, route := range routes {
		totalTrains := route.SourceTrainCount + route.TransitTrainCount
		
		fmt.Printf("%d. %s (%s) → %s (%s) → %s (%s)\n",
			i+1,
			route.SourceStation, route.SourceStationCode,
			route.TransitStation, route.TransitStationCode,
			route.DestStation, route.DestStationCode)
		
		fmt.Printf("   🚂 Trains: %d + %d = %d total | 📏 Distance: %s\n",
			route.SourceTrainCount, route.TransitTrainCount, totalTrains, route.Distance)
		
		fmt.Printf("   🔗 Details: https://etrain.info%s\n\n", route.ShowLink)
	}
	
	if len(routes) > displayCount {
		fmt.Printf("... and %d more routes available.\n", len(routes)-displayCount)
	}
}