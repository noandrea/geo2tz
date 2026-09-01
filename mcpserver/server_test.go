package mcpserver

import (
	"context"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/noandrea/geo2tz/v2/web"
	"github.com/stretchr/testify/assert"
)

func testSettings() web.ConfigSchema {
	return web.ConfigSchema{
		Tz: web.TzSchema{
			VersionFile:  "../tzdata/version.json",
			DatabaseName: "../tzdata/timezones.zip",
		},
		RuntimeVersion: "test",
	}
}

func TestNewServer(t *testing.T) {
	settings := web.ConfigSchema{
		Tz: web.TzSchema{
			VersionFile:  "../tzdata/version.json",
			DatabaseName: "file_not_found.zip",
		},
	}
	_, err := NewServer(settings)
	assert.ErrorIs(t, err, web.ErrDatabaseFileNotFound)

	settings = web.ConfigSchema{
		Tz: web.TzSchema{
			VersionFile:  "file_not_found.json",
			DatabaseName: "../tzdata/timezones.zip",
		},
	}
	_, err = NewServer(settings)
	assert.ErrorIs(t, err, web.ErrVersionFileNotFound)
}

func Test_getTimezone(t *testing.T) {
	server, err := NewServer(testSettings())
	assert.NoError(t, err)

	tests := []struct {
		name         string
		lat          float64
		lon          float64
		wantTimezone string
		wantErr      bool
	}{
		{
			"PASS: valid coordinates",
			51.477811,
			0,
			"Europe/London",
			false,
		},
		{
			"PASS: valid coordinates",
			41.9028,
			12.4964,
			"Europe/Rome",
			false,
		},
		{
			"FAIL: invalid latitude",
			100,
			11.831443,
			"",
			true,
		},
		{
			"FAIL: invalid longitude",
			43.42582,
			200,
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, out, lookupErr := server.getTimezone(context.Background(), nil, timezoneArgs{Lat: tt.lat, Lon: tt.lon})
			if tt.wantErr {
				assert.Error(t, lookupErr)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, lookupErr)
				assert.Equal(t, tt.wantTimezone, out.Timezone)
				assert.Equal(t, tt.lat, out.Lat)
				assert.Equal(t, tt.lon, out.Lon)
			}
		})
	}
}

func Test_getTzVersion(t *testing.T) {
	server, err := NewServer(testSettings())
	assert.NoError(t, err)

	_, out, err := server.getTzVersion(context.Background(), nil, struct{}{})
	assert.NoError(t, err)
	assert.NotEmpty(t, out.Version)
	assert.NotEmpty(t, out.URL)
	assert.NotEmpty(t, out.GeoDataURL)
}

func Test_MCPClientSession(t *testing.T) {
	server, err := NewServer(testSettings())
	assert.NoError(t, err)

	clientTransport, serverTransport := mcp.NewInMemoryTransports()
	done := make(chan error, 1)
	go func() {
		done <- server.Run(context.Background(), serverTransport)
	}()

	client := mcp.NewClient(&mcp.Implementation{Name: "geo2tz-test-client", Version: "v0.0.0"}, nil)
	session, err := client.Connect(context.Background(), clientTransport, nil)
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, session.Close())
	}()

	// list the tools
	tools, err := session.ListTools(context.Background(), nil)
	assert.NoError(t, err)
	assert.Len(t, tools.Tools, 2)

	// call the get_timezone tool
	tz, err := session.CallTool(context.Background(), &mcp.CallToolParams{
		Name:      "get_timezone",
		Arguments: map[string]any{"lat": 41.9028, "lon": 12.4964},
	})
	assert.NoError(t, err)
	assert.False(t, tz.IsError)
	assert.Equal(t, map[string]any{"timezone": "Europe/Rome", "lat": 41.9028, "lon": 12.4964}, tz.StructuredContent)

	// call the get_tz_version tool
	version, err := session.CallTool(context.Background(), &mcp.CallToolParams{
		Name:      "get_tz_version",
		Arguments: map[string]any{},
	})
	assert.NoError(t, err)
	assert.False(t, version.IsError)
	versionMap, ok := version.StructuredContent.(map[string]any)
	assert.True(t, ok)
	assert.NotEmpty(t, versionMap["version"])

	// call the get_timezone tool with invalid coordinates
	invalid, err := session.CallTool(context.Background(), &mcp.CallToolParams{
		Name:      "get_timezone",
		Arguments: map[string]any{"lat": 100, "lon": 0},
	})
	assert.NoError(t, err)
	assert.True(t, invalid.IsError)

	assert.NoError(t, session.Close())
	assert.NoError(t, <-done)
}
