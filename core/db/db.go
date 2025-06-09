package db

import (
	"errors"

	"github.com/noandrea/geo2tz/v2/core"
)

type TzDBIndex interface {
	Lookup(lat, lon float64) (core.TimeZoneData, error)
}

var (
	// ErrNotFound is returned when a timezone is not found
	ErrNotFound = errors.New("timezone not found")
	ErrInternal = errors.New("shape without tzID")
)
