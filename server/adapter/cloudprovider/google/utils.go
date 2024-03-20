package google

import "time"

const (
	layout        = "Jan 2, 2006 at 3:04pm" // for time.Format
	nanosPerMilli = 1e6
) // Special thanks: https://github.com/googleapis/google-api-go-client/blob/main/examples/fitness.go#L18C1-L21C2

// millisToTime converts Unix millis to time.Time.
func millisToTime(t int64) time.Time {
	// Special thanks: https://github.com/googleapis/google-api-go-client/blob/main/examples/fitness.go#L36
	return time.Unix(0, t*nanosPerMilli)
}

// NanosToTime converts Unix nanos to time.Time. Special thanks to: https://github.com/bronnika/devto-google-fit/blob/main/google-api/init.go#L49
func NanosToTime(t int64) time.Time {
	return time.Unix(0, t)
}

// TimeToNanos coverts time.Time to Unix nanos. Special thanks to: https://github.com/bronnika/devto-google-fit/blob/main/google-api/init.go#L54
func TimeToNanos(time2 time.Time) int64 {
	return time2.UnixNano()
}
