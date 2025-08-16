package main

import "github.com/spf13/cobra"

var (
	// Global flags
	cacheEnabled bool
	
	// Root command
	rootCmd = &cobra.Command{
		Use:   "trains",
		Short: "Railway route analysis tool",
		Long: `Trains CLI - A powerful command-line tool for analyzing Indian Railways 
train routes, connections, and timetables with intelligent caching and connection optimization.`,
	}
	
	// Via search command
	viaSearchCmd = &cobra.Command{
		Use:   "viasearch",
		Short: "Analyze train routes via intermediate stations",
		Long: `Analyze train connections via intermediate stations and find optimal routes.

This command fetches train data from etrain.info URLs and analyzes possible connections
with realistic layover times (1-4 hours) and matching running days.`,
		Example: `  trains viasearch -url="https://etrain.info/trains/Valsad-BL-to-H-Sahib-Nanded-NED-via-Kalyan-Jn-KYN"
  trains viasearch -url="https://etrain.info/trains/..." --no-cache
  trains viasearch -url="https://etrain.info/trains/..." -d=wed
  trains viasearch -url="https://etrain.info/trains/..." --day=sunday`,
		RunE: runViaSearch,
	}
	
	// Top search command
	topSearchCmd = &cobra.Command{
		Use:   "topsearch",
		Short: "Find top transit routes between stations",
		Long: `Find and analyze all possible transit routes between two stations.

This command fetches route data from etrain.info transit pages and shows all available
routes with their transit stations, train counts, and distances.`,
		Example: `  trains topsearch --url="https://etrain.info/transit/BL-NED?page=1"
  trains topsearch -u="https://etrain.info/transit/BL-NED?page=1" --no-cache
  trains topsearch --url="https://etrain.info/transit/BL-NED?page=1" --limit=8
  trains topsearch --url="https://etrain.info/transit/BL-NED?page=1" --max-distance=900`,
		RunE: runTopSearch,
	}
)

// initCommands initializes all CLI commands and flags
func initCommands() {
	// Add persistent flags to root command  
	rootCmd.PersistentFlags().BoolVar(&cacheEnabled, "cache", true, "Enable/disable caching")
	rootCmd.PersistentFlags().Bool("no-cache", false, "Disable caching (same as --cache=false)")
	
	// Add viasearch command
	rootCmd.AddCommand(viaSearchCmd)
	
	// Add topsearch command
	rootCmd.AddCommand(topSearchCmd)
	
	// Add flags specific to viasearch command
	viaSearchCmd.Flags().StringP("url", "u", "", "URL to fetch train data from (required)")
	viaSearchCmd.Flags().StringP("day", "d", "", "Filter by day of week (sun, mon, tue, wed, thu, fri, sat)")
	viaSearchCmd.MarkFlagRequired("url")
	
	// Add flags specific to topsearch command
	topSearchCmd.Flags().StringP("url", "u", "", "URL to fetch transit route data from (required)")
	topSearchCmd.Flags().IntP("limit", "l", 10, "Limit number of routes to show (default: 10)")
	topSearchCmd.Flags().IntP("max-distance", "m", 0, "Maximum distance in kilometers (0 = no limit)")
	topSearchCmd.MarkFlagRequired("url")
}