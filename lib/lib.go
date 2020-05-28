package lib

import (
	"fmt"
	"math/rand"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func init() {
	rand.Seed(time.Now().UnixNano())
}

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

// RndString generates a random string of the specified length
func RndString(size int) string {
	b := make([]rune, size)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
