package db

import (
	"testing"

	"github.com/noandrea/geo2tz/v2/helpers"
	"github.com/stretchr/testify/assert"
)

// TestGeo2TzTreeIndex_LookupZone tests the LookupZone function
// since the timezone is not always the same as the expected one, we need to check the reference timezone
func TestGeo2TzTreeIndex_LookupZone(t *testing.T) {
	var tests []struct {
		Tz       string  `json:"tz"`
		Lat      float64 `json:"lat"`
		Lon      float64 `json:"lon"`
		HasError bool    `json:"err,omitempty"`
	}

	// load the database
	gsi, err := NewGeo2TzRTreeIndexFromGeoJSON("../tzdata/timezones.zip")
	assert.NoError(t, err)
	assert.NotEmpty(t, gsi.Size())

	// load the timezone references
	var tzZones map[string]struct {
		Zone      string  `json:"zone"`
		UtcOffset float32 `json:"utc_offset_h"`
		Dst       struct {
			Start     string  `json:"start"`
			End       string  `json:"end"`
			Zone      string  `json:"zone"`
			UtcOffset float32 `json:"utc_offset_h"`
		} `json:"dst,omitempty"`
	}
	err = helpers.LoadJSON("testdata/zones.json", &tzZones)
	assert.NoError(t, err)
	assert.NotEmpty(t, tzZones)

	// load the coordinates
	err = helpers.LoadJSON("testdata/coordinates.json", &tests)
	assert.NoError(t, err)
	assert.NotEmpty(t, tests)

	for _, tt := range tests {
		t.Run(tt.Tz, func(t *testing.T) {
			got, err := gsi.Lookup(tt.Lat, tt.Lon)
			assert.NoError(t, err)

			if tt.HasError {
				t.Skip("skipping test as it is expected to fail (know error)")
			}

			// for oceans do exact match
			if IsOcean(got) {
				assert.Equal(t, tt.Tz, got, "expected %s to be %s for https://www.google.com/maps/@%v,%v,12z", tt.Tz, got, tt.Lat, tt.Lon)
				return
			}

			// get the zone for the expected timezone
			zoneExpected, ok := tzZones[tt.Tz]
			assert.True(t, ok, "timezone %s not found in zones.json", tt.Tz)

			// get the reference timezone for the expected timezone
			zoneGot, ok := tzZones[got]
			assert.True(t, ok, "timezone %s not found in zones.json", got)

			if !ok {
				assert.Equal(t, zoneExpected.Zone, got, "expected %s (%s) to be %s (%s) for https://www.google.com/maps/@%v,%v,12z", tt.Tz, zoneExpected.Zone, got, zoneGot.Zone, tt.Lat, tt.Lon)
			} else {
				assert.Equal(t, zoneExpected.Zone, zoneGot.Zone, "expected %s (%s)  to be %s (%s) for https://www.google.com/maps/@%v,%v,12z", tt.Tz, zoneExpected.Zone, got, zoneGot.Zone, tt.Lat, tt.Lon)
			}
		})
	}
}
