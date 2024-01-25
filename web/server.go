package web

import (
	"context"
	"crypto/subtle"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/evanoberholster/timezoneLookup/v2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/noandrea/geo2tz/v2/helpers"

	"golang.org/x/crypto/blake2b"
)

// constant valuses for lat / lon
const (
	Latitude        = "lat"
	Longitude       = "lon"
	compareEquals   = 1
	teardownTimeout = 10 * time.Second
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

type Server struct {
	config          ConfigSchema
	tzDB            timezoneLookup.Timezonecache
	tzRelease       TzRelease
	echo            *echo.Echo
	authEnabled     bool
	authHashedToken []byte
}

func (server *Server) Start() error {
	return server.echo.Start(server.config.Web.ListenAddress)
}

func (server *Server) Teardown() (err error) {
	server.tzDB.Close()
	ctx, cancel := context.WithTimeout(context.Background(), teardownTimeout)
	defer cancel()
	if server.echo != nil {
		err = server.echo.Shutdown(ctx)
	}
	return
}

func NewServer(config ConfigSchema) (*Server, error) {
	var server Server
	server.config = config
	server.echo = echo.New()
	// open the database
	f, err := os.Open(config.Tz.DatabaseName)
	if err != nil {
		err = errors.Join(ErrorDatabaseFileNotFound, fmt.Errorf("error opening the timezone database: %w", err))
		return nil, err
	}
	defer f.Close()

	// load the database
	if err = server.tzDB.Load(f); err != nil {
		err = errors.Join(ErrorDatabaseFileInvalid, fmt.Errorf("error loading the timezone database: %w", err))
		return nil, err
	}

	// check token authorization
	server.authHashedToken = hash(config.Web.AuthTokenValue)
	if len(config.Web.AuthTokenValue) > 0 {
		server.echo.Logger.Info("Authorization enabled")
		server.authEnabled = true
	} else {
		server.echo.Logger.Info("Authorization disabled")
	}

	server.echo.HideBanner = true
	server.echo.Use(middleware.CORS())
	server.echo.Use(middleware.Logger())
	server.echo.Use(middleware.Recover())

	// load the release info
	if err = helpers.LoadJSON(config.Tz.VersionFile, &server.tzRelease); err != nil {
		err = errors.Join(ErrorVersionFileNotFound, err, fmt.Errorf("error loading the timezone release info: %w", err))
		return nil, err
	}

	// register routes
	server.echo.GET("/tz/:lat/:lon", server.handleTzRequest)
	server.echo.GET("/tz/version", server.handleTzVersion)

	return &server, nil
}

func (server *Server) handleTzRequest(c echo.Context) error {
	// token verification
	if server.authEnabled {
		requestToken := c.QueryParam(server.config.Web.AuthTokenParamName)
		if !isEq(server.authHashedToken, requestToken) {
			server.echo.Logger.Errorf("request unauthorized, invalid token: %v", requestToken)
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{"message": "unauthorized"})
		}
	}
	// parse latitude
	lat, err := parseCoordinate(c.Param(Latitude), Latitude)
	if err != nil {
		server.echo.Logger.Errorf("error parsing latitude: %v", err)
		return c.JSON(http.StatusBadRequest, newErrResponse(err))
	}
	// parse longitude
	lon, err := parseCoordinate(c.Param(Longitude), Longitude)
	if err != nil {
		server.echo.Logger.Errorf("error parsing longitude: %v", err)
		return c.JSON(http.StatusBadRequest, newErrResponse(err))
	}

	// query the coordinates
	res, err := server.tzDB.Search(lat, lon)
	if err != nil {
		server.echo.Logger.Errorf("error querying the timezone db: %v", err)
		return c.JSON(http.StatusInternalServerError, newErrResponse(err))
	}
	if res.Name == "" {
		notFoundErr := fmt.Errorf("timezone not found for coordinates %f,%f", lat, lon)
		server.echo.Logger.Errorf("error querying the timezone db: %v", notFoundErr)
		return c.JSON(http.StatusNotFound, newErrResponse(notFoundErr))
	}
 
	tzr := newTzResponse(res.Name, lat, lon)
	return c.JSON(http.StatusOK, tzr)
}

func newTzResponse(tzName string, lat, lon float64) map[string]any {
	return map[string]any{"tz": tzName, "coords": map[string]float64{Latitude: lat, Longitude: lon}}
}

func newErrResponse(err error) map[string]any {
	return map[string]any{"message": err.Error()}
}

func (server *Server) handleTzVersion(c echo.Context) error {
	return c.JSON(http.StatusOK, server.tzRelease)
}

// parseCoordinate parse a string into a coordinate
func parseCoordinate(val, side string) (float64, error) {
	if strings.TrimSpace(val) == "" {
		return 0, fmt.Errorf("empty coordinates value")
	}

	c, err := strconv.ParseFloat(val, 64)
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
	return c, nil
}
