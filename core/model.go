package core

import "time"

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
	LocalTime time.Time `json:"local_time"`
	UTCTime   time.Time `json:"utc_time"`
	Sunset    time.Time `json:"sunset"`
	Sunrise   time.Time `json:"sunrise"`
	// ISOWeek    int       `json:"iso_week"`
	IsDST      bool `json:"is_dst"`
	IsDaylight bool `json:"is_daylight"`
}
