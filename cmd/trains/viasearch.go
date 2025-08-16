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
		dayFilter, err = parser.ValidateAndNormalizeDay(dayFilter)
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
	
	// Parse train data from JavaScript objects in HTML
	trains := parser.ParseTrainData(htmlContent)
	
	fmt.Printf("Found %d trains\n", len(trains))
	
	// Extract route information from URL
	sourceStation, destinationStation, transitStation := parser.ExtractRouteInfo(url)
	
	// Group trains by source-destination pairs
	connections := analyzeConnections(trains, dayFilter, sourceStation, destinationStation, transitStation)
	
	// Generate results
	generateConnections(connections, dayFilter, sourceStation, destinationStation, transitStation)
	
	return nil
}

// separateTrainsByRoute separates trains into source-to-transit and transit-to-destination segments
func separateTrainsByRoute(trains []types.TrainData, sourceStation, transitStation, destinationStation string) ([]types.TrainData, []types.TrainData) {
	sourceToTransit := make([]types.TrainData, 0, len(trains)/2)
	transitToDestination := make([]types.TrainData, 0, len(trains)/2)
	
	for _, train := range trains {
		if train.SourceStationCode == sourceStation && train.DestStationCode == transitStation {
			sourceToTransit = append(sourceToTransit, train)
		} else if train.SourceStationCode == transitStation && train.DestStationCode == destinationStation {
			transitToDestination = append(transitToDestination, train)
		}
	}
	
	return sourceToTransit, transitToDestination
}

// isValidConnection checks if a connection is valid (not a "No Connection" type)
func isValidConnection(connection types.RouteConnection) bool {
	return connection.Connection != "No Connection - No common running days" && 
		   connection.Connection != "No Connection - Insufficient layover time" &&
		   connection.Connection != "No Connection - Layover too long (>4h)"
}

// analyzeConnections finds valid train connections
func analyzeConnections(trains []types.TrainData, dayFilter string, sourceStation string, destinationStation string, transitStation string) []types.RouteConnection {
	connections := make([]types.RouteConnection, 0, 10) // Estimate initial capacity
	
	// Separate trains by route segments
	sourceToTransit, transitToDestination := separateTrainsByRoute(trains, sourceStation, transitStation, destinationStation)
	
	fmt.Printf("%s to %s trains: %d\n", sourceStation, transitStation, len(sourceToTransit))
	fmt.Printf("%s to %s trains: %d\n", transitStation, destinationStation, len(transitToDestination))
	
	// Find valid connections
	for _, train1 := range sourceToTransit {
		for _, train2 := range transitToDestination {
			connection := analyzeConnection(train1, train2)
			if !isValidConnection(connection) {
				continue
			}
			
			// Apply day filter if specified
			if dayFilter != "" && !connectionMatchesDay(connection, dayFilter) {
				continue
			}
			
			connections = append(connections, connection)
		}
	}
	
	return connections
}

// analyzeConnection analyzes connection between two trains
func analyzeConnection(train1, train2 types.TrainData) types.RouteConnection {
	connection := types.RouteConnection{
		Train1: train1,
		Train2: train2,
	}
	
	// First check if trains have overlapping running days
	commonDays := parser.GetCommonRunningDays(train1.RunningDays, train2.RunningDays)
	if commonDays == "" {
		connection.Connection = "No Connection - No common running days"
		return connection
	}
	
	// Parse times
	time1 := parser.ParseTime(train1.DestTime) // Arrival at Kalyan
	time2 := parser.ParseTime(train2.SourceTime) // Departure from Kalyan
	
	// Check if connection is possible (layover between 1-4 hours)
	layoverMinutes := time2 - time1
	if layoverMinutes >= parser.MinLayoverMinutes && layoverMinutes <= parser.MaxLayoverMinutes {
		connection.Connection = fmt.Sprintf("Same day - %dh %dm layover (%s)", layoverMinutes/parser.MinutesPerHour, layoverMinutes%parser.MinutesPerHour, commonDays)
		
		// Calculate total journey time
		startTime := parser.ParseTime(train1.SourceTime)
		endTime := parser.ParseTime(train2.DestTime)
		if endTime < startTime {
			endTime += parser.MinutesPerDay // Next day
		}
		totalMinutes := endTime - startTime
		connection.TotalTime = fmt.Sprintf("%dh %dm", totalMinutes/parser.MinutesPerHour, totalMinutes%parser.MinutesPerHour)
	} else {
		// Check next day connection (layover between 1-4 hours)
		nextDayTime2 := time2 + parser.MinutesPerDay
		layover := nextDayTime2 - time1
		if layover >= parser.MinLayoverMinutes && layover <= parser.MaxLayoverMinutes {
			connection.Connection = fmt.Sprintf("Next day - %dh %dm layover (%s)", layover/parser.MinutesPerHour, layover%parser.MinutesPerHour, commonDays)
			
			startTime := parser.ParseTime(train1.SourceTime)
			endTime := parser.ParseTime(train2.DestTime) + parser.MinutesPerDay // Next day
			totalMinutes := endTime - startTime
			connection.TotalTime = fmt.Sprintf("%dh %dm", totalMinutes/parser.MinutesPerHour, totalMinutes%parser.MinutesPerHour)
		} else {
			// No valid connection
			if layoverMinutes < parser.MinLayoverMinutes || layover < parser.MinLayoverMinutes {
				connection.Connection = "No Connection - Insufficient layover time"
			} else {
				connection.Connection = "No Connection - Layover too long (>4h)"
			}
		}
	}
	
	return connection
}

// connectionMatchesDay checks if connection runs on specified day
func connectionMatchesDay(connection types.RouteConnection, dayFilter string) bool {
	// Check if the connection runs on the specified day
	commonDays := parser.GetCommonRunningDays(connection.Train1.RunningDays, connection.Train2.RunningDays)
	if commonDays == "" {
		return false
	}
	
	// Check if the dayFilter is in the common running days
	dayAbbreviation := parser.GetDayAbbreviation(dayFilter)
	return strings.Contains(commonDays, dayAbbreviation)
}

// generateConnections displays the connection results
func generateConnections(connections []types.RouteConnection, dayFilter string, sourceStation string, destinationStation string, transitStation string) {
	if dayFilter != "" {
		fmt.Printf("\n=== TRAIN CONNECTIONS FROM %s TO %s VIA %s (Available on %s) ===\n\n", sourceStation, destinationStation, transitStation, dayFilter)
	} else {
		fmt.Printf("\n=== TRAIN CONNECTIONS FROM %s TO %s VIA %s ===\n\n", sourceStation, destinationStation, transitStation)
	}
	
	// Filter connections under 19 hours
	var validConnections []types.RouteConnection
	for _, conn := range connections {
		if parser.IsUnder19Hours(conn.TotalTime) {
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
			i+1, conn.Train1.Number, conn.Train1.Name, conn.Train2.Number, conn.Train2.Name)
		fmt.Printf("   %s %s → %s %s → %s %s\n", 
			sourceStation, conn.Train1.SourceTime, transitStation, conn.Train1.DestTime, destinationStation, conn.Train2.DestTime)
		fmt.Printf("   Total Time: %s | Connection: %s\n", 
			conn.TotalTime, conn.Connection)
		fmt.Printf("   Days: %s + %s\n\n", 
			parser.FormatRunningDays(conn.Train1.RunningDays), parser.FormatRunningDays(conn.Train2.RunningDays))
	}
}