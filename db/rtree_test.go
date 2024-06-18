package db

import (
	"testing"

	"github.com/noandrea/geo2tz/v2/helpers"
	"github.com/stretchr/testify/assert"
)

func TestGeo2TzTreeIndex_Lookup(t *testing.T) {
	var tests []struct {
		Tz  string  `json:"tz"`
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	}

	// load the database
	gsi, err := NewGeo2TzRTreeIndexFromGeoJSON("../tzdata/timezones.zip")
	assert.NoError(t, err)
	assert.NotEmpty(t, gsi.Size())

	// load the timezone references
	var tz2etc map[string]string
	err = helpers.LoadJSON("testdata/tz2etc.json", &tz2etc)
	assert.NoError(t, err)
	assert.NotEmpty(t, tz2etc)

	// load the coordinates
	err = helpers.LoadJSON("testdata/coordinates.json", &tests)
	assert.NoError(t, err)
	assert.NotEmpty(t, tests)

	for _, tt := range tests {
		t.Run(tt.Tz, func(t *testing.T) {
			got, err := gsi.Lookup(tt.Lat, tt.Lon)
			assert.NoError(t, err)

			// check if the expected timezone is in the same etc reference
			etcExpected, ok := tz2etc[tt.Tz]
			assert.True(t, ok, "timezone %s not found in tz2etc.json", tt.Tz)

			// get the reference timezone for the expected timezone
			etcGot, ok := tz2etc[got]

			if !ok {
				assert.Equal(t, etcExpected, got, "expected %s (%s) to be %s (%s) for https://www.google.com/maps/@%v,%v,12z", tt.Tz, etcExpected, got, etcGot, tt.Lat, tt.Lon)
			} else {
				assert.Equal(t, etcExpected, etcGot, "expected %s (%s)  to be %s (%s) for https://www.google.com/maps/@%v,%v,12z", tt.Tz, etcExpected, got, etcGot, tt.Lat, tt.Lon)
			}
		})
	}
}
