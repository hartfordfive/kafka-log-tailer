package lib

import (
	"fmt"
	"time"
)

const (
	dateTimeFormat = "2006-01-02T15:04:51.555Z"
	dateFormat     = "2006-01-02"
	timeFormat     = "15:04"
)

// FromUtcToLocalTime takes a UTC time string and returns the time in the users current timezone
func FromUtcToLocalTime(fromDateTime, localTZ string) string {
	// locallyEnteredDate := "2017-07-16"
	// locallyEnteredTime := "14:00"

	// // Build a time object from received fields (time objects include zone info)
	// // We are assuming the code is running on a server that is in the same zone as the current user
	// zone, _ := time.Now().Zone() // get the local zone
	// dateTimeZ := locallyEnteredDate + " " + locallyEnteredTime + " " + zone
	// dte, err := time.Parse(dateTimeFormat, dateTimeZ)
	// if err != nil {
	// 	log.Fatal("Error parsing entered datetime", err)
	// }
	// fmt.Println("dte:", dte) // dte is a legit time object
	// // Perhaps we are saving this in a database.
	// // A good database driver should save the time object as UTC in a time with zone field,
	// // and return a time object with UTC as zone.

	// // For the sake of this example, let's assume an object identical to `dte` is returned
	// // dte := ReceiveFromDatabase()

	// // Convert received date to local.
	// // Note the use of the convenient "Local" location https://golang.org/pkg/time/#LoadLocation.
	// localLoc, err := time.LoadLocation("Local")
	// if err != nil {
	// 	log.Fatal(`Failed to load location "Local"`)
	// }
	// localDateTime := dte.In(localLoc)

	// fmt.Println("Date:", localDateTime.Format(dateFormat))
	// fmt.Println("Time:", localDateTime.Format(timeFormat))

	//now := time.Now()
	//loc := now.Location()

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
