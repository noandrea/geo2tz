package db

import (
	"errors"
	"time"
)

type TzDBIndex interface {
	Lookup(lat, lon float64) (string, error)
	LookupTime(tzID string) (local, utc time.Time, isDST bool, zone string, offset int, err error)
}

var (
	// ErrNotFound is returned when a timezone is not found
	ErrNotFound = errors.New("timezone not found")
	ErrInternal = errors.New("shape without tzID")
)
