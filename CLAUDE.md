# CLAUDE.md

This file provides guidance to Claude Code when working with the Trains CLI tool.

## Trains CLI Overview

A powerful Indian Railways route analysis tool with two main commands:
- `viasearch` - Analyze train connections via intermediate stations
- `topsearch` - Discover all possible transit routes between stations

## Command Usage

### ViaSearch (Connection Analysis)
```bash
# Analyze connections via specific intermediate station
./trains viasearch --url="https://etrain.info/trains/Valsad-BL-to-H-Sahib-Nanded-NED-via-Kalyan-Jn-KYN"

# Filter by specific day
./trains viasearch --url="<URL>" --day=wednesday
```

### TopSearch (Route Discovery)
```bash
# Single page: Quick analysis (sorted by train availability)
./trains topsearch --url="https://etrain.info/transit/BL-NED?page=1"

# Multi-page: Comprehensive analysis (sorted by shortest distance)
./trains topsearch --url="https://etrain.info/transit/BL-NED"

# With distance filter and limit
./trains topsearch --url="<URL>" --max-distance=800 --limit=10
```

## Key Features

### Smart Sorting Logic
- **Multi-page URLs** (no `?page=`): Sort by distance (shortest first)
- **Single-page URLs**: Sort by train availability (most trains first)  
- **When `--max-distance` used**: Always sort by distance

### Caching System
- 24-hour automatic caching in `./cache/` directory

## Common Workflows

### Two-Step Analysis for Best Train Selection

**Step 1**: Use `topsearch` to discover all possible routes
```bash
# Find all routes with distance/train availability analysis
./trains topsearch --url="https://etrain.info/transit/BL-NED" --limit=5
```

**Step 2**: Use `viasearch` with Details URLs for final train selection
```bash
# Copy the "Details" URL from topsearch output for deeper analysis
./trains viasearch --url="https://etrain.info/trains/Valsad-BL-to-H-Sahib-Nanded-NED-via-Kalyan-Jn-KYN" --day=wednesday
```

This two-step approach helps users:
1. **Discover** all route options with distance/train count comparison
2. **Analyze** specific routes with actual train timings, connections, and layovers
3. **Decide** on the best train combination based on schedule compatibility

### Finding Best Routes (Distance Priority)
```bash
# Get all routes sorted by shortest distance
./trains topsearch --url="https://etrain.info/transit/BL-NED"

# Within specific distance limit
./trains topsearch --url="https://etrain.info/transit/BL-NED" --max-distance=1000
```

### Finding Routes with Most Train Options
```bash
# Quick overview sorted by train availability
./trains topsearch --url="https://etrain.info/transit/BL-NED?page=1"
```

## Output Understanding

### TopSearch Results
- Lists routes: `SOURCE → TRANSIT → DESTINATION`
- Shows train counts: `15 + 4 = 19 total` (source trains + transit trains)
- Distance and details link provided for each route
- **Key**: Use the "Details" URLs (🔗) as input for `viasearch` to get specific train timings and connections

### ViaSearch Results  
- Shows specific train combinations with timings
- Validates realistic layover times (1-4 hours)
- Confirms running days compatibility
- Reports total journey time under 19 hours
