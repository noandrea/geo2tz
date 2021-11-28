package server

import (
	"context"
	"crypto/subtle"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/evanoberholster/timezoneLookup"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/time/rate"
)

// constant values for lat / lon
const (
	Latitude      = "lat"
	Longitude     = "lon"
	compareEquals = 1
)

var (
	tz timezoneLookup.TimezoneInterface
	e  *echo.Echo
)

// hash calculate the hash of a string
func hash(data ...interface{}) []byte {
	hash := blake2b.Sum256([]byte(fmt.Sprint(data...)))
	return hash[:]
}

// isEq check if the hash of the second value is equals to the first value
func isEq(expectedTokenHash []byte, actualToken string) bool {
	return subtle.ConstantTimeCompare(expectedTokenHash, hash(actualToken)) == compareEquals
}

// Start starts the web server
func Start(config ConfigSchema) (err error) {
	encoding, err := timezoneLookup.EncodingFromString(config.Tz.Encoding)
	if err != nil {
		log.Errorln("invalid encoding:", err)
		return
	}
	// open the database
	tz, err = timezoneLookup.LoadTimezones(
		timezoneLookup.Config{
			DatabaseType: config.Tz.DatabaseType, // memory or boltdb
			DatabaseName: config.Tz.DatabaseName, // Name without suffix
			Snappy:       config.Tz.Snappy,
			Encoding:     encoding, // json or msgpack
		})
	if err != nil {
		log.Errorln("failed to load timezones:", err)
		return
	}

	// check token authorization
	hashedToken := hash(config.Web.AuthTokenValue)
	authEnabled := false
	if len(config.Web.AuthTokenValue) > 0 {
		log.Info("authorization enabled, using request parameter:", config.Web.AuthTokenParamName)
		authEnabled = true
	} else {
		log.Info("authorization disabled")
	}

	// echo start
	e = echo.New()
	e.HideBanner = true
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	if config.Web.RateLimit > 0 {
		log.Infoln("rate limit enabled:", config.Web.RateLimit)
		e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(config.Web.RateLimit))))
	} else {
		log.Infoln("rate limit disabled")
	}

	// logger
	e.GET("/tz/:lat/:lon", func(c echo.Context) (err error) {
		// token verification
		if authEnabled {
			requestToken := c.QueryParam(config.Web.AuthTokenParamName)
			if !isEq(hashedToken, requestToken) {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{"message": "unauthorized"})
			}
		}
		//parse latitude
		lat, err := parseCoordinate(c.Param(Latitude), Latitude)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": fmt.Sprint(err)})
		}
		//parse longitude
		lon, err := parseCoordinate(c.Param(Longitude), Longitude)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": fmt.Sprint(err)})
		}
		// build coordinates object
		coords := timezoneLookup.Coord{
			Lat: lat,
			Lon: lon,
		}
		// query the coordinates
		res, err := tz.Query(coords)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": fmt.Sprint(err)})
		}
		return c.JSON(http.StatusOK, map[string]interface{}{"tz": res, "coords": coords})
	})
	err = e.Start(config.Web.ListenAddress)
	return
}

// parseCoordinate parse a string into a coordinate
func parseCoordinate(val, side string) (float32, error) {
	if strings.TrimSpace(val) == "" {
		return 0, fmt.Errorf("empty coordinates value")
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
