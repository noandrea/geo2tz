package web

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"encoding/json"

	"github.com/stretchr/testify/assert"
)

func Test_parseCoordinate(t *testing.T) {
	type c struct {
		val  string
		side string
	}
	tests := []struct {
		ll      c
		want    float64
		wantErr bool
	}{
		{c{"22", Latitude}, 22, false},
		{c{"78.312", Longitude}, 78.312, false},
		{c{"0x429c9fbe", Longitude}, 0, true}, // 78.312 in hex
		{c{"", Longitude}, 0, true},
		{c{"   ", Longitude}, 0, true},
		{c{"2e4", Longitude}, 0, true},
		{c{"not a number", Longitude}, 0, true},
		{c{"-90.1", Latitude}, 0, true},
		{c{"90.001", Latitude}, 0, true},
		{c{"-180.1", Longitude}, 0, true},
		{c{"180.001", Longitude}, 0, true},
		{c{"43.42582", Latitude}, 43.42582, false},
		{c{"11.831443", Longitude}, 11.831443, false},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprint(tt.ll), func(t *testing.T) {
			got, err := parseCoordinate(tt.ll.val, tt.ll.side)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseCoordinate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseCoordinate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hash(t *testing.T) {
	tests := []struct {
		name string
		args []interface{}
		want []byte
	}{
		{
			"one element",
			[]interface{}{
				"test1",
			},
			[]byte{229, 104, 55, 204, 215, 163, 141, 103, 149, 211, 10, 194, 171, 99, 236, 204, 140, 43, 87, 18, 137, 166, 45, 196, 6, 187, 98, 118, 126, 136, 176, 108},
		},
		{
			"two elements",
			[]interface{}{
				"test1",
				"test2",
			},
			[]byte{84, 182, 224, 44, 5, 184, 19, 24, 41, 163, 6, 53, 242, 3, 167, 200, 192, 113, 61, 137, 208, 241, 141, 225, 134, 61, 78, 124, 88, 254, 117, 159},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hash(tt.args...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("hash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isEq(t *testing.T) {
	type args struct {
		expectedTokenHash []byte
		actualToken       string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"PASS: token matches",
			args{
				[]byte{229, 104, 55, 204, 215, 163, 141, 103, 149, 211, 10, 194, 171, 99, 236, 204, 140, 43, 87, 18, 137, 166, 45, 196, 6, 187, 98, 118, 126, 136, 176, 108},
				"test1",
			},
			true,
		},
		{
			"FAIL: token mismatch",
			args{
				[]byte{84, 182, 224, 44, 5, 184, 19, 24, 41, 163, 6, 53, 242, 3, 167, 200, 192, 113, 61, 137, 208, 241, 141, 225, 134, 61, 78, 124, 88, 254, 117, 159},
				"test1",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isEq(tt.args.expectedTokenHash, tt.args.actualToken); got != tt.want {
				t.Errorf("isEq() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewServer(t *testing.T) {
	settings := ConfigSchema{
		Tz: TzSchema{
			VersionFile:  "file_not_found.json",
			DatabaseName: "../tzdata/timezones.zip",
		},
	}
	_, err := NewServer(settings)
	assert.ErrorIs(t, err, ErrorVersionFileNotFound)

	settings = ConfigSchema{
		Tz: TzSchema{
			VersionFile:  "../tzdata/version.json",
			DatabaseName: "timezone_not_found.db",
		},
	}
	_, err = NewServer(settings)
	assert.ErrorIs(t, err, ErrorDatabaseFileNotFound)
}

func Test_TzVersion(t *testing.T) {
	settings := ConfigSchema{
		Tz: TzSchema{
			VersionFile:  "../tzdata/version.json",
			DatabaseName: "../tzdata/timezones.zip",
		},
	}
	server, err := NewServer(settings)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := server.echo.NewContext(req, rec)
	c.SetPath("/tz/version")

	// Assertions
	if assert.NoError(t, server.handleTzVersion(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var version TzRelease
		reply := rec.Body.String()
		err = json.Unmarshal([]byte(reply), &version)
		assert.NoError(t, err)
		assert.NotEmpty(t, version.Version)
		assert.NotEmpty(t, version.URL)
		assert.NotEmpty(t, version.GeoDataURL)
	}
}

func Test_TzRequest(t *testing.T) {
	settings := ConfigSchema{
		Tz: TzSchema{
			VersionFile:  "../tzdata/version.json",
			DatabaseName: "../tzdata/timezones.zip",
		},
	}
	server, err := NewServer(settings)
	assert.NoError(t, err)

	tests := []struct {
		name      string
		lat       string
		lon       string
		wantCode  int
		wantReply string
	}{
		{
			"PASS: valid coordinates",
			"51.477811",
			"0",
			http.StatusOK,
			`{"coords":{"lat":51.477811,"lon":0},"tz":"Europe/London"}`,
		},
		{
			"PASS: valid coordinates",
			"41.9028",
			"12.4964",
			http.StatusOK,
			`{"coords":{"lat":41.9028,"lon":12.4964},"tz":"Europe/Rome"}`,
		},
		{
			"FAIL: invalid latitude",
			"100",
			"11.831443",
			http.StatusBadRequest,
			`{"message":"lat value 100 out of range (-90/+90)"}`,
		},
		{
			"FAIL: invalid longitude",
			"43.42582",
			"200",
			http.StatusBadRequest,
			`{"message":"lon value 200 out of range (-180/+180)"}`,
		},
		{
			"FAIL: invalid latitude and longitude",
			"100",
			"200",
			http.StatusBadRequest,
			`{"message":"lat value 100 out of range (-90/+90)"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := server.echo.NewContext(req, rec)
			c.SetPath("/tz/:lat/:lon")
			c.SetParamNames("lat", "lon")
			c.SetParamValues(tt.lat, tt.lon)

			// Assertions
			if assert.NoError(t, server.handleTzRequest(c)) {
				assert.Equal(t, tt.wantCode, rec.Code)
				assert.Equal(t, tt.wantReply, strings.TrimSpace(rec.Body.String()))
			}
		})
	}
}
