package types

import (
	"fmt"
	"time"
)

// TrainData represents train information from etrain.info
type TrainData struct {
	Type               string `json:"typ"`
	Number             string `json:"num"`
	Name               string `json:"name"`
	SourceStationCode  string `json:"s"`    // Source station code
	SourceTime         string `json:"st"`   // Source time
	DestStationCode    string `json:"d"`    // Destination station code  
	DestTime           string `json:"dt"`   // Destination time
	TravelTime         string `json:"tt"`   // Travel time
	RunningDays        string `json:"dy"`   // Running days (1=Sun, 2=Mon, etc.)
	BookingInfo        string `json:"book"`
	ArrivalPlatform    int    `json:"arp"`
}

// RouteConnection represents a connection between two trains via intermediate station
type RouteConnection struct {
	Train1     TrainData
	Train2     TrainData
	TotalTime  string
	Connection string
}

// TransitRoute represents a transit route between stations
type TransitRoute struct {
	SourceStation      string
	SourceStationCode  string
	SourceTrainCount   int
	TransitStation     string
	TransitStationCode string
	TransitTrainCount  int
	DestStation        string
	DestStationCode    string
	Distance           string
	ShowLink           string
}

// CacheEntry represents a cached HTTP response
type CacheEntry struct {
	URL       string    `json:"url"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// String returns a string representation of TrainData
func (t TrainData) String() string {
	return fmt.Sprintf("%s %s (%s→%s at %s-%s)", t.Number, t.Name, t.SourceStationCode, t.DestStationCode, t.SourceTime, t.DestTime)
}

// String returns a string representation of RouteConnection
func (r RouteConnection) String() string {
	return fmt.Sprintf("%s + %s | Total: %s | %s", r.Train1.String(), r.Train2.String(), r.TotalTime, r.Connection)
}

// String returns a string representation of TransitRoute
func (t TransitRoute) String() string {
	return fmt.Sprintf("%s (%s) → %s (%s) → %s (%s) | Distance: %s | Trains: %d+%d=%d", 
		t.SourceStation, t.SourceStationCode, t.TransitStation, t.TransitStationCode, 
		t.DestStation, t.DestStationCode, t.Distance, t.SourceTrainCount, 
		t.TransitTrainCount, t.SourceTrainCount+t.TransitTrainCount)
}

// String returns a string representation of CacheEntry
func (c CacheEntry) String() string {
	return fmt.Sprintf("Cache[%s] from %v", c.URL, c.Timestamp.Format("2006-01-02 15:04:05"))
}