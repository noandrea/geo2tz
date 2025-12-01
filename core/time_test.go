package core

import (
	"testing"

	"github.com/noandrea/geo2tz/v2/helpers"
	"github.com/stretchr/testify/assert"
)

// TestGeo2TzTreeIndex_LookupZone tests the LookupZone function
// since the timezone is not always the same as the expected one, we need to check the reference timezone
func TestGeo2TzTimeInfo_ComputeTimeInfo(t *testing.T) {
	var tests []struct {
		Input struct {
			Lat float64 `json:"lat"`
			Lon float64 `json:"lon"`
		} `json:"input"`
		Expected TimeZoneData `json:"expected"`
	}

	// load the database
	gsi, err := NewGeo2TzRTreeIndexFromGeoJSON("../tzdata/timezones.zip")
	assert.NoError(t, err)

	// load the coordinates
	err = helpers.LoadJSON("testdata/time_info.json", &tests)
	assert.NoError(t, err)
	assert.NotEmpty(t, tests)

	for _, tt := range tests {
		t.Run(tt.Expected.TZ, func(t *testing.T) {
			got, err := gsi.Lookup(tt.Input.Lat, tt.Input.Lon)
			assert.NoError(t, err)
			err = ComputeTimeData(&got, tt.Expected.TimeInfo.UTCTime)
			assert.NoError(t, err)
			assert.Equalf(t, tt.Expected, got, "expected %s got %s", tt.Expected, got)
		})
	}
}
