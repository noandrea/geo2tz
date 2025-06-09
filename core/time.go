package core

import (
	"time"

	"github.com/nathan-osman/go-sunrise"
)

func ComputeTimeData(tzd *TimeZoneData, at time.Time) error {
	location, err := time.LoadLocation(tzd.TZ)
	if err != nil {
		return err
	}
	currentLocalTime := at.In(location)
	currentUTCTime := currentLocalTime.UTC()

	sunrise, sunset := sunrise.SunriseSunset(
		tzd.Coords.Lat, tzd.Coords.Lon,
		currentLocalTime.Year(), currentLocalTime.Month(), currentLocalTime.Day(),
	)

	tzd.TimeInfo = &TimeInfo{
		LocalTime:  currentLocalTime,
		Sunrise:    sunrise.In(location),
		Sunset:     sunset.In(location),
		IsDST:      currentLocalTime.IsDST(),
		UTCTime:    currentUTCTime,
		IsDaylight: currentUTCTime.After(sunrise) && currentUTCTime.Before(sunset),
	}
	return nil
}
