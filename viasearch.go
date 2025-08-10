package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	
	"github.com/spf13/cobra"
)

// runViaSearch handles the viasearch command
func runViaSearch(cmd *cobra.Command, args []string) error {
	// Get URL flag
	url, err := cmd.Flags().GetString("url")
	if err != nil {
		return fmt.Errorf("error getting url flag: %v", err)
	}
	
	// Get day filter flag
	dayFilter, err := cmd.Flags().GetString("day")
	if err != nil {
		return fmt.Errorf("error getting day flag: %v", err)
	}
	
	// Validate and normalize day filter
	if dayFilter != "" {
		dayFilter, err = validateAndNormalizeDay(dayFilter)
		if err != nil {
			return err
		}
	}
	
	// Handle --no-cache flag (check parent flags since it's persistent)
	noCacheFlag, _ := cmd.Parent().Flags().GetBool("no-cache")
	if noCacheFlag {
		cacheEnabled = false
	}
	
	fmt.Printf("🚂 Starting train route analysis...\n")
	fmt.Printf("📍 URL: %s\n", url)
	fmt.Printf("💾 Cache: %t\n", cacheEnabled)
	if dayFilter != "" {
		fmt.Printf("📅 Day Filter: %s\n", dayFilter)
	}
	fmt.Println()
	
	// Initialize cache directory if caching is enabled
	if cacheEnabled {
		if err := initCache(); err != nil {
			return fmt.Errorf("error initializing cache: %v", err)
		}
	}
	
	// Fetch the webpage with or without caching
	var htmlContent string
	
	if cacheEnabled {
		htmlContent, err = fetchWithCache(url)
	} else {
		fmt.Printf("🌐 Fetching from network (cache disabled): %s\n", url)
		htmlContent, err = fetchFromNetwork(url)
	}
	
	if err != nil {
		return fmt.Errorf("error fetching URL: %v", err)
	}
	
	// Parse train data from JavaScript objects in HTML
	trains := parseTrainData(htmlContent)
	
	fmt.Printf("Found %d trains\n", len(trains))
	
	// Extract route information from URL
	sourceStation, destinationStation, transitStation := extractRouteInfo(url)
	
	// Group trains by source-destination pairs
	connections := analyzeConnections(trains, dayFilter, sourceStation, destinationStation, transitStation)
	
	// Generate results
	generateConnections(connections, dayFilter, sourceStation, destinationStation, transitStation)
	
	return nil
}

// parseTrainData extracts train data from HTML JavaScript objects
func parseTrainData(htmlContent string) []TrainData {
	var trains []TrainData
	
	// Look for train data in JavaScript objects
	trainPattern := regexp.MustCompile(`data-train='({[^}]+})'`)
	matches := trainPattern.FindAllStringSubmatch(htmlContent, -1)
	
	fmt.Printf("Found %d train data objects\n", len(matches))
	
	for i, match := range matches {
		if len(match) > 1 {
			var train TrainData
			err := json.Unmarshal([]byte(match[1]), &train)
			if err != nil {
				fmt.Printf("Error parsing train %d: %v\n", i, err)
				continue
			}
			trains = append(trains, train)
			
			if i < 5 { // Debug first 5 trains
				fmt.Printf("Train %d: %s %s from %s(%s) to %s(%s) at %s-%s\n", 
					i, train.Num, train.Name, train.S, train.St, train.D, train.Dt, train.St, train.Dt)
			}
		}
	}
	
	return trains
}

// extractRouteInfo extracts source, destination, and transit station codes from URL
func extractRouteInfo(url string) (string, string, string) {
	// URL format: https://etrain.info/trains/Valsad-BL-to-H-Sahib-Nanded-NED-via-Kalyan-Jn-KYN
	// Extract source, destination, and transit station codes
	
	// Default fallback values
	defaultSource := "BL"
	defaultDestination := "NED" 
	defaultTransit := "KYN"
	
	// Split by "-via-" to separate main route from transit
	parts := strings.Split(url, "-via-")
	if len(parts) < 2 {
		return defaultSource, defaultDestination, defaultTransit
	}
	
	mainRoute := parts[0]
	viaPart := parts[1]
	
	// Extract transit station code (last part after final hyphen)
	transitStation := defaultTransit
	if strings.Contains(viaPart, "-") {
		lastHyphen := strings.LastIndex(viaPart, "-")
		if lastHyphen > 0 && lastHyphen < len(viaPart)-1 {
			transitStation = viaPart[lastHyphen+1:]
		}
	}
	
	// Extract source and destination from main route
	// Format: ...Source-SRCCODE-to-...Destination-DESTCODE-via-...
	toIndex := strings.Index(mainRoute, "-to-")
	if toIndex == -1 {
		return defaultSource, defaultDestination, transitStation
	}
	
	sourcePart := mainRoute[:toIndex]
	destPart := mainRoute[toIndex+4:] // Skip "-to-"
	
	// Extract source station code (last part before "-to-")
	sourceStation := defaultSource
	if strings.Contains(sourcePart, "-") {
		lastHyphen := strings.LastIndex(sourcePart, "-")
		if lastHyphen > 0 && lastHyphen < len(sourcePart)-1 {
			sourceStation = sourcePart[lastHyphen+1:]
		}
	}
	
	// Extract destination station code (last part before "-via-")
	destinationStation := defaultDestination
	if strings.Contains(destPart, "-") {
		lastHyphen := strings.LastIndex(destPart, "-")
		if lastHyphen > 0 && lastHyphen < len(destPart)-1 {
			destinationStation = destPart[lastHyphen+1:]
		}
	}
	
	return sourceStation, destinationStation, transitStation
}

// analyzeConnections finds valid train connections
func analyzeConnections(trains []TrainData, dayFilter string, sourceStation string, destinationStation string, transitStation string) []RouteConnection {
	var connections []RouteConnection
	
	// Separate trains by route segments using dynamic stations
	var sourceToTransit []TrainData
	var transitToDestination []TrainData
	
	for _, train := range trains {
		if train.S == sourceStation && train.D == transitStation {
			sourceToTransit = append(sourceToTransit, train)
		} else if train.S == transitStation && train.D == destinationStation {
			transitToDestination = append(transitToDestination, train)
		}
	}
	
	fmt.Printf("%s to %s trains: %d\n", sourceStation, transitStation, len(sourceToTransit))
	fmt.Printf("%s to %s trains: %d\n", transitStation, destinationStation, len(transitToDestination))
	
	// Find valid connections
	for _, train1 := range sourceToTransit {
		for _, train2 := range transitToDestination {
			connection := analyzeConnection(train1, train2)
			if connection.Connection != "No Connection - No common running days" && 
			   connection.Connection != "No Connection - Insufficient layover time" &&
			   connection.Connection != "No Connection - Layover too long (>4h)" {
				
				// Apply day filter if specified
				if dayFilter != "" && !connectionMatchesDay(connection, dayFilter) {
					continue
				}
				
				connections = append(connections, connection)
			}
		}
	}
	
	return connections
}

// analyzeConnection analyzes connection between two trains
func analyzeConnection(train1, train2 TrainData) RouteConnection {
	connection := RouteConnection{
		Train1: train1,
		Train2: train2,
	}
	
	// First check if trains have overlapping running days
	commonDays := getCommonRunningDays(train1.Dy, train2.Dy)
	if commonDays == "" {
		connection.Connection = "No Connection - No common running days"
		return connection
	}
	
	// Parse times
	time1 := parseTime(train1.Dt) // Arrival at Kalyan
	time2 := parseTime(train2.St) // Departure from Kalyan
	
	// Check if connection is possible (layover between 1-4 hours)
	layoverMinutes := time2 - time1
	if layoverMinutes >= 60 && layoverMinutes <= 240 {
		connection.Connection = fmt.Sprintf("Same day - %dh %dm layover (%s)", layoverMinutes/60, layoverMinutes%60, commonDays)
		
		// Calculate total journey time
		startTime := parseTime(train1.St)
		endTime := parseTime(train2.Dt)
		if endTime < startTime {
			endTime += 24 * 60 // Next day
		}
		totalMinutes := endTime - startTime
		connection.TotalTime = fmt.Sprintf("%dh %dm", totalMinutes/60, totalMinutes%60)
	} else {
		// Check next day connection (layover between 1-4 hours)
		nextDayTime2 := time2 + 24*60
		layover := nextDayTime2 - time1
		if layover >= 60 && layover <= 240 {
			connection.Connection = fmt.Sprintf("Next day - %dh %dm layover (%s)", layover/60, layover%60, commonDays)
			
			startTime := parseTime(train1.St)
			endTime := parseTime(train2.Dt) + 24*60 // Next day
			totalMinutes := endTime - startTime
			connection.TotalTime = fmt.Sprintf("%dh %dm", totalMinutes/60, totalMinutes%60)
		} else {
			// No valid connection
			if layoverMinutes < 60 || layover < 60 {
				connection.Connection = "No Connection - Insufficient layover time"
			} else {
				connection.Connection = "No Connection - Layover too long (>4h)"
			}
		}
	}
	
	return connection
}

// connectionMatchesDay checks if connection runs on specified day
func connectionMatchesDay(connection RouteConnection, dayFilter string) bool {
	// Check if the connection runs on the specified day
	commonDays := getCommonRunningDays(connection.Train1.Dy, connection.Train2.Dy)
	if commonDays == "" {
		return false
	}
	
	// Check if the dayFilter is in the common running days
	dayAbbreviation := getDayAbbreviation(dayFilter)
	return strings.Contains(commonDays, dayAbbreviation)
}

// generateConnections displays the connection results
func generateConnections(connections []RouteConnection, dayFilter string, sourceStation string, destinationStation string, transitStation string) {
	if dayFilter != "" {
		fmt.Printf("\n=== TRAIN CONNECTIONS FROM %s TO %s VIA %s (Available on %s) ===\n\n", sourceStation, destinationStation, transitStation, dayFilter)
	} else {
		fmt.Printf("\n=== TRAIN CONNECTIONS FROM %s TO %s VIA %s ===\n\n", sourceStation, destinationStation, transitStation)
	}
	
	// Filter connections under 19 hours
	var validConnections []RouteConnection
	for _, conn := range connections {
		if isUnder19Hours(conn.TotalTime) {
			validConnections = append(validConnections, conn)
		}
	}
	
	if dayFilter != "" {
		fmt.Printf("Found %d connections under 19 hours available on %s:\n\n", len(validConnections), dayFilter)
	} else {
		fmt.Printf("Found %d connections under 19 hours:\n\n", len(validConnections))
	}
	
	for i, conn := range validConnections {
		fmt.Printf("%d. %s %s + %s %s\n", 
			i+1, conn.Train1.Num, conn.Train1.Name, conn.Train2.Num, conn.Train2.Name)
		fmt.Printf("   %s %s → %s %s → %s %s\n", 
			sourceStation, conn.Train1.St, transitStation, conn.Train1.Dt, destinationStation, conn.Train2.Dt)
		fmt.Printf("   Total Time: %s | Connection: %s\n", 
			conn.TotalTime, conn.Connection)
		fmt.Printf("   Days: %s + %s\n\n", 
			formatRunningDays(conn.Train1.Dy), formatRunningDays(conn.Train2.Dy))
	}
}