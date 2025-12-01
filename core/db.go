package core

import (
	"errors"
)

type TzDBIndex interface {
	Lookup(lat, lon float64) (TimeZoneData, error)
}

var (
	// ErrNotFound is returned when a timezone is not found
	ErrNotFound = errors.New("timezone not found")
	ErrInternal = errors.New("shape without tzID")
)
