package google

import "time"

// HydrationStruct defines the hydration data type provided by `Google Fit`. Special thanks to https://github.com/bronnika/devto-google-fit/blob/main/models/models.go#L52C1-L56C2
type HydrationStruct struct {
	Amount    int       `json:"amount"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}
