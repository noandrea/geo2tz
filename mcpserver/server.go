// Package mcpserver implements an MCP (Model Context Protocol) server that
// exposes the geo2tz timezone lookup as tools for MCP clients (eg. AI assistants).
package mcpserver

import (
	"context"
	"errors"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/noandrea/geo2tz/v2/db"
	"github.com/noandrea/geo2tz/v2/helpers"
	"github.com/noandrea/geo2tz/v2/web"
)

// valid ranges for coordinates
const (
	minLatitude  = -90.0
	maxLatitude  = 90.0
	minLongitude = -180.0
	maxLongitude = 180.0
)

// Server is a geo2tz MCP server.
type Server struct {
	tzDB      db.TzDBIndex
	tzRelease web.TzRelease
	mcp       *mcp.Server
}

// timezoneArgs is the input for the get_timezone tool.
type timezoneArgs struct {
	Lat float64 `json:"lat" jsonschema:"the latitude of the location, in the range -90 to +90"`
	Lon float64 `json:"lon" jsonschema:"the longitude of the location, in the range -180 to +180"`
}

// timezoneResult is the output of the get_timezone tool.
type timezoneResult struct {
	Timezone string  `json:"timezone" jsonschema:"the IANA timezone identifier for the coordinates, for example Europe/Rome"`
	Lat      float64 `json:"lat" jsonschema:"the latitude queried"`
	Lon      float64 `json:"lon" jsonschema:"the longitude queried"`
}

// tzVersionResult is the output of the get_tz_version tool.
type tzVersionResult struct {
	Version    string `json:"version" jsonschema:"the version of the timezone database, for example 2025b"`
	URL        string `json:"url" jsonschema:"the release page of the timezone database version"`
	GeoDataURL string `json:"geo_data_url" jsonschema:"the download URL of the timezone database geojson"`
}

// NewServer creates a new MCP server with the geo2tz tools registered.
func NewServer(config web.ConfigSchema) (*Server, error) {
	var server Server

	// load the database
	tzDB, err := db.NewGeo2TzRTreeIndexFromGeoJSON(config.Tz.DatabaseName)
	if err != nil {
		return nil, errors.Join(web.ErrDatabaseFileNotFound, err)
	}
	server.tzDB = tzDB

	// load the release info
	if err = helpers.LoadJSON(config.Tz.VersionFile, &server.tzRelease); err != nil {
		err = errors.Join(web.ErrVersionFileNotFound, err, fmt.Errorf("error loading the timezone release info: %w", err))
		return nil, err
	}

	server.mcp = mcp.NewServer(&mcp.Implementation{
		Name:    "geo2tz",
		Version: config.RuntimeVersion,
	}, nil)

	// register tools
	mcp.AddTool(server.mcp, &mcp.Tool{
		Name: "get_timezone",
		Description: "returns the IANA timezone identifier for a pair of geo coordinates (latitude/longitude), " +
			"it works offline on the bundled timezone boundaries database",
	}, server.getTimezone)
	mcp.AddTool(server.mcp, &mcp.Tool{
		Name:        "get_tz_version",
		Description: "returns the version of the timezone boundaries database in use",
	}, server.getTzVersion)

	return &server, nil
}

// Run runs the server over the given transport until the client disconnects
// or the context is cancelled.
func (server *Server) Run(ctx context.Context, t mcp.Transport) error {
	return server.mcp.Run(ctx, t)
}

// getTimezone handles the get_timezone tool call.
func (server *Server) getTimezone(_ context.Context, _ *mcp.CallToolRequest, args timezoneArgs) (*mcp.CallToolResult, timezoneResult, error) {
	if args.Lat < minLatitude || args.Lat > maxLatitude {
		return nil, timezoneResult{}, fmt.Errorf("lat value %v out of range (-90/+90)", args.Lat)
	}
	if args.Lon < minLongitude || args.Lon > maxLongitude {
		return nil, timezoneResult{}, fmt.Errorf("lon value %v out of range (-180/+180)", args.Lon)
	}

	tzID, err := server.tzDB.Lookup(args.Lat, args.Lon)
	if err != nil {
		return nil, timezoneResult{}, fmt.Errorf("error querying the timezone db: %w", err)
	}
	return nil, timezoneResult{Timezone: tzID, Lat: args.Lat, Lon: args.Lon}, nil
}

// getTzVersion handles the get_tz_version tool call.
func (server *Server) getTzVersion(_ context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, tzVersionResult, error) {
	return nil, tzVersionResult{
		Version:    server.tzRelease.Version,
		URL:        server.tzRelease.URL,
		GeoDataURL: server.tzRelease.GeoDataURL,
	}, nil
}
