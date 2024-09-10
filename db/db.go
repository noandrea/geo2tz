package db

import (
	"errors"
	"time"
)

type TzReply struct {
	TZ     string `json:"tz,omitempty"`
	Coords struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	} `json:"coords,omitempty"`
}

type ZoneReply struct {
	TzReply
	Local  time.Time `json:"local"`
	UTC    time.Time `json:"utc"`
	IsDST  bool      `json:"is_dst"`
	Offset int       `json:"offset"`
	Zone   string    `json:"zone"`
}

type TzDBIndex interface {
	Lookup(lat, lon float64) (TzReply, error)
	LookupZone(lat, lon float64) (ZoneReply, error)
	LookupTime(tzID string) (ZoneReply, error)
}

var (
	// ErrNotFound is returned when a timezone is not found
	ErrNotFound = errors.New("timezone not found")
	ErrInternal = errors.New("shape without tzID")
)
