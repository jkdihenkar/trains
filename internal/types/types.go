package types

import "time"

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