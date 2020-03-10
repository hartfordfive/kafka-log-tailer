package lib

import (
	"fmt"
	"time"
)

// FromUtcToLocalTime takes a UTC time string and returns the time in the users current timezone
func FromUtcToLocalTime(fromDateTime, localTZ string) string {
	t, err := time.Parse(time.RFC3339Nano, fromDateTime)
	if err != nil {
		return fromDateTime
	}

	loc, err := time.LoadLocation(localTZ)
	if err != nil {
		return fromDateTime
	}

	return fmt.Sprintf("%s", t.In(loc))
}
