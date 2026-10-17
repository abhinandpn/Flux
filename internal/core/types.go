package core

import "time"

// CounterSample contains cumulative network counters
// captured at a specific time.
type CounterSample struct {
	Interface     string
	BytesReceived uint64
	BytesSent     uint64
	Timestamp     time.Time
}

// TrafficRate contains calculated network traffic rates.
// DownloadBPS and UploadBPS are measured in bytes per second.
type TrafficRate struct {
	Interface   string
	DownloadBPS float64
	UploadBPS   float64
}
