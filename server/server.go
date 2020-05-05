package server

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/evanoberholster/timezoneLookup"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// constant valuses for lat / lon
const (
	Latitude  = "lat"
	Longitude = "lon"
)

var (
	tz timezoneLookup.TimezoneInterface
	e  *echo.Echo
)

// Start starts the web server
func Start(config ConfigSchema) (err error) {
	// open the database
	tz, err := timezoneLookup.LoadTimezones(
		timezoneLookup.Config{
			DatabaseType: config.Tz.DatabaseType, // memory or boltdb
			DatabaseName: config.Tz.DatabaseName, // Name without suffix
			Snappy:       config.Tz.Snappy,
			Encoding:     config.Tz.Encoding, // json or msgpack
		})
	if err != nil {
		return
	}
	// echo start
	e = echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// logger
	e.GET("/tz/:lat/:lon", func(c echo.Context) (err error) {
		//parse latitude
		lat, err := parseCoordinate(c.Param(Latitude), Latitude)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": err})
		}
		//parse longitude
		lon, err := parseCoordinate(c.Param(Longitude), Longitude)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": err})
		}
		// build coordinates object
		coords := timezoneLookup.Coord{
			Lat: float32(lat),
			Lon: float32(lon),
		}
		// query the coordinates
		res, err := tz.Query(coords)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": err})
		}
		return c.JSON(http.StatusOK, map[string]interface{}{"tz": res, "coords": coords})
	})
	err = e.Start(config.Web.ListenAddress)
	return
}

func parseCoordinate(val, side string) (float32, error) {
	if len(strings.TrimSpace(val)) == 0 {
		return 0, fmt.Errorf("Empty coordinate value")
	}
	c, err := strconv.ParseFloat(val, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid type for %s, a number is required (eg. 45.3123)", side)
	}
	switch side {
	case Latitude:
		if c < -90 || c > 90 {
			return 0, fmt.Errorf("%s value %s out of range (-90/+90)", side, val)
		}
	case Longitude:
		if c < -180 || c > 180 {
			return 0, fmt.Errorf("%s value %s out of range (-180/+180)", side, val)
		}
	}
	return float32(c), nil
}

// Teardown gracefully release resources
func Teardown() (err error) {
	if tz != nil {
		tz.Close()
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if e != nil {
		err = e.Shutdown(ctx)
	}
	return
}
