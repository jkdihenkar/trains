package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// ParseTime converts time string to minutes since midnight
func ParseTime(timeStr string) int {
	parts := strings.Split(timeStr, ":")
	if len(parts) != 2 {
		return 0
	}
	
	hours, _ := strconv.Atoi(parts[0])
	minutes, _ := strconv.Atoi(parts[1])
	return hours*60 + minutes
}

// ParseDistanceKm extracts distance in kilometers from distance string
func ParseDistanceKm(distanceStr string) int {
	// Parse distance string like "754 Kms" or "1038 Kms"
	distancePattern := regexp.MustCompile(`(\d+)\s*Kms?`)
	matches := distancePattern.FindStringSubmatch(distanceStr)
	if len(matches) >= 2 {
		if distance, err := strconv.Atoi(matches[1]); err == nil {
			return distance
		}
	}
	return 0
}

// GetCommonRunningDays finds common running days between two trains
func GetCommonRunningDays(days1, days2 string) string {
	dayNames := []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
	var commonDays []string
	
	// Ensure both strings are at least 7 characters
	if len(days1) < 7 || len(days2) < 7 {
		return ""
	}
	
	for i := 0; i < 7; i++ {
		if days1[i] == '1' && days2[i] == '1' {
			commonDays = append(commonDays, dayNames[i])
		}
	}
	
	if len(commonDays) == 0 {
		return ""
	}
	
	return strings.Join(commonDays, ",")
}

// FormatRunningDays formats running days string for display
func FormatRunningDays(dayStr string) string {
	days := []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
	var runningDays []string
	
	for i, char := range dayStr {
		if i < len(days) && char == '1' {
			runningDays = append(runningDays, days[i])
		}
	}
	
	if len(runningDays) == 7 {
		return "Daily"
	}
	return strings.Join(runningDays, ",")
}

// IsUnder19Hours checks if total time is under 19 hours
func IsUnder19Hours(timeStr string) bool {
	// Parse time string like "16h 52m"
	re := regexp.MustCompile(`(\d+)h\s*(\d+)m`)
	matches := re.FindStringSubmatch(timeStr)
	if len(matches) >= 3 {
		hours, _ := strconv.Atoi(matches[1])
		return hours < 19
	}
	return false
}

// ValidateAndNormalizeDay validates and normalizes day input
func ValidateAndNormalizeDay(day string) (string, error) {
	day = strings.ToLower(strings.TrimSpace(day))
	
	dayMap := map[string]string{
		"sun":       "Sunday",
		"sunday":    "Sunday",
		"mon":       "Monday", 
		"monday":    "Monday",
		"tue":       "Tuesday",
		"tuesday":   "Tuesday", 
		"wed":       "Wednesday",
		"wednesday": "Wednesday",
		"thu":       "Thursday",
		"thursday":  "Thursday",
		"fri":       "Friday",
		"friday":    "Friday",
		"sat":       "Saturday", 
		"saturday":  "Saturday",
	}
	
	if normalized, ok := dayMap[day]; ok {
		return normalized, nil
	}
	
	return "", fmt.Errorf("invalid day '%s'. Valid options: sun, mon, tue, wed, thu, fri, sat (or full names)", day)
}

// GetDayAbbreviation converts full day name to abbreviation
func GetDayAbbreviation(fullDayName string) string {
	dayMap := map[string]string{
		"Sunday":    "Sun",
		"Monday":    "Mon", 
		"Tuesday":   "Tue",
		"Wednesday": "Wed",
		"Thursday":  "Thu",
		"Friday":    "Fri",
		"Saturday":  "Sat",
	}
	
	if abbrev, ok := dayMap[fullDayName]; ok {
		return abbrev
	}
	return fullDayName
}

// ShouldFetchAllPages determines if URL requires multi-page fetching
func ShouldFetchAllPages(url string) bool {
	// Check if URL is a transit URL without page parameter
	// Examples: 
	// - https://etrain.info/transit/BL-NED (should fetch all pages)
	// - https://etrain.info/transit/BL-NED?page=1 (single page specified)
	return strings.Contains(url, "/transit/") && !strings.Contains(url, "page=")
}