# Trains CLI - Railway Route Analysis Tool

A powerful command-line tool for analyzing Indian Railways train routes, connections, and timetables with intelligent caching and connection optimization. Built with **Cobra CLI framework** for professional command-line experience.

## Features

- **Route Analysis**: Analyze train connections via intermediate stations
- **Top Route Discovery**: Find all possible transit routes between stations with smart sorting
- **Multi-page Fetching**: Automatically scrolls through all pages for comprehensive route discovery
- **Smart Caching**: 24-hour file-based caching system with MD5-hashed storage
- **Connection Optimization**: Find optimal train connections with realistic layover times (1-4 hours)
- **Running Days Validation**: Ensures connecting trains run on the same days
- **Day Filtering**: Filter results by specific day of the week (Monday, Tuesday, etc.)
- **Professional CLI**: Built with spf13/cobra for rich command-line experience
- **Auto-completion**: Bash, Zsh, Fish, and PowerShell completion support

## Installation

```bash
# Clone or download the project
cd trains/

# Install dependencies
go mod tidy

# Build the binary
go build -o trains .

# Run directly with Go (development)
go run main.go <command> [options]
```

## Usage

### Basic Commands

```bash
# Show help
./trains help

# Show help for specific commands
./trains help viasearch
./trains help topsearch

# Analyze train routes via intermediate stations
./trains viasearch --url="<etrain.info URL>"

# Find top transit routes between stations
./trains topsearch --url="<etrain.info transit URL>"

# Disable caching for fresh data
./trains viasearch --url="<URL>" --no-cache
./trains topsearch --url="<URL>" --no-cache
```

### Examples

```bash
# Analyze Valsad to Nanded via Kalyan Junction
./trains viasearch --url="https://etrain.info/trains/Valsad-BL-to-H-Sahib-Nanded-NED-via-Kalyan-Jn-KYN"

# Find top 8 routes from Valsad to Nanded
./trains topsearch --url="https://etrain.info/transit/BL-NED?page=1" --limit=8

# Find routes within 900km distance (single page)
./trains topsearch --url="https://etrain.info/transit/BL-NED?page=1" --max-distance=900

# Multi-page: fetch all pages and sort by shortest distance
./trains topsearch --url="https://etrain.info/transit/BL-NED"

# Multi-page with distance filter
./trains topsearch --url="https://etrain.info/transit/BL-NED" --max-distance=1000

# Using short flags
./trains viasearch -u="https://etrain.info/trains/Valsad-BL-to-H-Sahib-Nanded-NED-via-Kalyan-Jn-KYN"
./trains topsearch -u="https://etrain.info/transit/BL-NED?page=1" -l=5 -m=800

# Analyze with cache disabled (multiple ways)
./trains viasearch --url="<URL>" --no-cache
./trains topsearch --url="<URL>" --cache=false

# Filter by specific day of the week (viasearch only)
./trains viasearch --url="<URL>" --day=wednesday
./trains viasearch --url="<URL>" -d=fri
```

## Command Reference

### `viasearch`
Analyzes train routes via intermediate stations and finds optimal connections.

**Flags:**
- `-u, --url string`: URL to fetch train data from (required)  
- `-d, --day string`: Filter by day of week (sun, mon, tue, wed, thu, fri, sat)
- `-h, --help`: Help for viasearch command

**Features:**
- Finds connections under 19 hours total journey time
- Validates layover times between 1-4 hours for realistic transfers  
- Checks running days compatibility between connecting trains
- **Day filtering**: Filter connections by specific day of the week
- Provides detailed connection analysis with timings and days

### `topsearch`
Finds and ranks all possible transit routes between two stations.

**Flags:**
- `-u, --url string`: URL to fetch transit route data from (required)
- `-l, --limit int`: Limit number of routes to show (default: 10)
- `-m, --max-distance int`: Maximum distance in kilometers (0 = no limit)
- `-h, --help`: Help for topsearch command

**Features:**
- Discovers all available transit routes between stations
- **Smart sorting**: Distance-based (shortest first) for multi-page URLs, train availability for single pages
- **Multi-page fetching**: Automatically scrolls through all pages when no page parameter is specified
- **Distance filtering**: Filter routes by maximum distance in kilometers
- Shows train counts, distances, and detailed route links
- Supports both single page and comprehensive multi-page analysis
- Perfect for route discovery and comparison

**Global Flags:**
- `--cache`: Enable/disable caching (default: true)
- `--no-cache`: Disable caching (same as --cache=false)

### Shell Completion

Enable shell completion for better user experience:

```bash
# Bash
source <(./trains completion bash)

# Zsh  
source <(./trains completion zsh)

# Fish
./trains completion fish | source

# PowerShell
./trains completion powershell | Out-String | Invoke-Expression
```

## Output Format

### ViaSearch Output
The viasearch command provides detailed connection analysis:

```
🚂 Starting train route analysis...
📍 URL: https://etrain.info/trains/...
💾 Cache: true
📅 Day Filter: Wednesday (when using day filter)

💾 Cache hit for ... (cached 4m ago)
Found 81 trains

=== TRAIN CONNECTIONS FROM VALSAD TO NANDED VIA KALYAN (Available on Wednesday) ===

Found 3 connections under 19 hours available on Wednesday:

1. 11089 BGKT PUNE EXPRESS + 17617 TAPOVAN EXPRESS
   Valsad 01:08 → Kalyan 04:42 → Nanded 18:00
   Total Time: 16h 52m | Connection: Same day - 1h 45m layover (Wed)
   Days: Wed + Daily
```

### TopSearch Output
The topsearch command shows all available transit routes:

#### Single Page Mode
```
🚂 Starting top transit route analysis...
📍 URL: https://etrain.info/transit/BL-NED?page=1
💾 Cache: true
📊 Limit: 8 routes

Found 25 transit routes total
📊 Sorting by train availability (most trains first)...
Showing top 8 routes (sorted by total train availability):
```

#### Multi-Page Mode
```  
🚂 Starting top transit route analysis...
📍 URL: https://etrain.info/transit/BL-NED
💾 Cache: true
📊 Limit: 10 routes

🔄 Detecting multi-page transit data, fetching all pages...
📄 Fetching page 1: https://etrain.info/transit/BL-NED?page=1
📄 Fetching page 2: https://etrain.info/transit/BL-NED?page=2
📄 Fetching page 3: https://etrain.info/transit/BL-NED?page=3
...
🔄 Fetched 6 pages with total 138 routes
Found 138 transit routes total

=== TOP TRANSIT ROUTES ===

📊 Sorting by distance (shortest routes first)...
Showing top 10 routes (sorted by shortest distance):

1. VALSAD (BL) → KALYAN JN (KYN) → H SAHIB NANDED (NED)
   🚂 Trains: 15 + 4 = 19 total | 📏 Distance: 754 Kms
   🔗 Details: https://etrain.info/trains/...
```

## Connection Analysis Rules

1. **Total Journey Time**: Must be under 19 hours
2. **Layover Time**: Between 1-4 hours for realistic connections
3. **Running Days**: Both trains must run on at least one common day
4. **Same Day Connections**: Prioritized over next-day connections

## Day Filtering

The `--day` or `-d` flag allows you to filter connections that are available on specific days:

- **Supported formats**: Short (sun, mon, tue, wed, thu, fri, sat) or full names (sunday, monday, etc.)
- **Case insensitive**: `Wed`, `wed`, `WEDNESDAY` all work
- **Filtering logic**: Only shows connections where both trains run on the specified day
- **Smart matching**: Automatically matches day abbreviations in train schedules

### Examples:

```bash
# Wednesday connections only
./trains viasearch --url="<URL>" --day=wed
./trains viasearch --url="<URL>" -d=wednesday

# Sunday connections only  
./trains viasearch --url="<URL>" --day=sunday
./trains viasearch --url="<URL>" -d=sun
```

This feature is particularly useful for:
- **Weekend travel planning**: Filter by Saturday/Sunday
- **Specific day requirements**: Business travel on weekdays
- **Optimized results**: Reduce results to only relevant days

## Multi-Page Route Discovery

The `topsearch` command automatically detects when to fetch multiple pages:

### Single Page Mode
- **URL**: `https://etrain.info/transit/BL-NED?page=1` (page parameter specified)
- **Behavior**: Fetches only the specified page
- **Sorting**: By train availability (most trains first)
- **Use case**: Quick overview of top routes

### Multi-Page Mode  
- **URL**: `https://etrain.info/transit/BL-NED` (no page parameter)
- **Behavior**: Automatically fetches all available pages (typically 6-7 pages)
- **Sorting**: By distance (shortest routes first)
- **Use case**: Comprehensive analysis of all possible routes

### Multi-Page Benefits
- **Complete coverage**: Discovers all possible transit routes (100+ routes)
- **Distance optimization**: Shows shortest routes first for efficient travel
- **Smart caching**: Caches each page separately for faster subsequent runs
- **Progress tracking**: Shows real-time progress as pages are fetched

## Cache System

- **Storage**: `./cache/` directory with MD5-hashed filenames
- **Expiry**: 24 hours (configurable in code)
- **Benefits**: Faster subsequent runs, reduced server load
- **Control**: Use `-cache=false` to bypass cache

## Technical Details

- **Language**: Go 1.21+
- **CLI Framework**: spf13/cobra v1.9.1+
- **Dependencies**: Minimal external dependencies (cobra + pflag)
- **Data Source**: etrain.info HTML parsing
- **Cache Format**: JSON with URL, content, and timestamp

## Project Structure

```
trains/
├── main.go              # Main CLI application
├── cache/              # Cached responses (auto-created)
├── go.mod              # Go module file
├── trains              # Compiled binary
└── README.md           # This file
```

## Development

```bash
# Run tests
go test ./...

# Format code
go fmt ./...

# Build
go build -o trains .

# Clean cache
rm -rf cache/
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## License

This project is open source and available under the MIT License.