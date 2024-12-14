package tehrantime

import (
	"time"
	_ "time/tzdata" // Embed timezone database
)

func FormattedDateTime(t1 time.Time) time.Time {
	tehranLoc, err := time.LoadLocation("Asia/Tehran")
	if err != nil {
		// Fallback to a fixed offset for Tehran (UTC+3:30)
		tehranLoc = time.FixedZone("Asia/Tehran", 3*60*60+30*60)
	}
	tehranTime := t1.In(tehranLoc)
	return tehranTime
}
