package core

import (
	"fmt"
	"time"
)

type TimeZoneData struct {
	TZ       string    `json:"tz"`
	Coords   Coords    `json:"coords"`
	TimeInfo *TimeInfo `json:"time_info,omitempty"`
}

type Coords struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type TimeInfo struct {
	UTCTime   time.Time `json:"utc_time"`
	LocalTime time.Time `json:"local_time"`
	Sunset    time.Time `json:"sunset"`
	Sunrise   time.Time `json:"sunrise"`
	// ISOWeek    int       `json:"iso_week"`
	IsDST      bool `json:"is_dst"`
	IsDaylight bool `json:"is_daylight"`
}

// implement stringer for TimeZoneData
func (tzd TimeZoneData) String() string {
	return fmt.Sprintf("TZ: %s, Coords: %s, TimeInfo: %s", tzd.TZ, tzd.Coords, tzd.TimeInfo)
}

// implement stringer for coords
func (c Coords) String() string {
	return fmt.Sprintf("Lat: %f, Lon: %f", c.Lat, c.Lon)
}

// implement stringer for TimeInfo
func (ti *TimeInfo) String() string {
	return fmt.Sprintf("LocalTime: %s, UTCTime: %s, Sunset: %s, Sunrise: %s, IsDST: %t, IsDaylight: %t",
		ti.LocalTime.Format(time.RFC3339), ti.UTCTime.Format(time.RFC3339), ti.Sunset.Format(time.RFC3339), ti.Sunrise.Format(time.RFC3339), ti.IsDST, ti.IsDaylight)
}
