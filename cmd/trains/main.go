package main

import (
	"fmt"
	"os"
)

func main() {
	// Initialize CLI commands and flags
	initCommands()
	
	// Execute root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}