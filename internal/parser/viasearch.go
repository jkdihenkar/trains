package parser

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"trains/internal/types"
)

// ParseTrainData extracts train data from HTML JavaScript objects
func ParseTrainData(htmlContent string) []types.TrainData {
	var trains []types.TrainData
	
	// Look for train data in JavaScript objects
	trainPattern := regexp.MustCompile(`data-train='({[^}]+})'`)
	matches := trainPattern.FindAllStringSubmatch(htmlContent, -1)
	
	fmt.Printf("Found %d train data objects\n", len(matches))
	
	for i, match := range matches {
		if len(match) > 1 {
			var train types.TrainData
			err := json.Unmarshal([]byte(match[1]), &train)
			if err != nil {
				fmt.Printf("Error parsing train %d: %v\n", i, err)
				continue
			}
			trains = append(trains, train)
			
			if i < 5 { // Debug first 5 trains
				fmt.Printf("Train %d: %s %s from %s(%s) to %s(%s) at %s-%s\n", 
					i, train.Number, train.Name, train.SourceStationCode, train.SourceTime, train.DestStationCode, train.DestTime, train.SourceTime, train.DestTime)
			}
		}
	}
	
	return trains
}

// ExtractRouteInfo extracts source, destination, and transit station codes from URL
func ExtractRouteInfo(url string) (string, string, string) {
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