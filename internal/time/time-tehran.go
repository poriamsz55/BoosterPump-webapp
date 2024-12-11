package tehrantime

import "time"

// When retrieving, convert to Tehran time
func FormattedDateTime(t1 time.Time) time.Time {
	tehranLoc, _ := time.LoadLocation("Asia/Tehran")
	tehranTime := t1.In(tehranLoc)

	return tehranTime
}
