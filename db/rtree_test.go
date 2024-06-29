package db

import (
	"archive/zip"
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
		NotFound bool    `json:"not_found,omitempty"`
	}

	// load the database
	gsi, err := NewGeo2TzRTreeIndexFromGeoJSON("../tzdata/timezones.zip")
	assert.NoError(t, err)
	assert.NotEmpty(t, gsi.Size())

	// load the coordinates
	err = helpers.LoadJSON("testdata/coordinates.json", &tests)
	assert.NoError(t, err)
	assert.NotEmpty(t, tests)

	for _, tt := range tests {
		t.Run(tt.Tz, func(t *testing.T) {
			got, err := gsi.Lookup(tt.Lat, tt.Lon)
			if tt.NotFound {
				assert.ErrorIs(t, err, ErrNotFound, "expected %s to be not_found for https://www.google.com/maps/@%v,%v,12z", got, tt.Lat, tt.Lon)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, got, tt.Tz, "expected %s to be %s for https://www.google.com/maps/@%v,%v,12z", tt.Tz, got, tt.Lat, tt.Lon)
		})
	}
}

// benchmark the lookup function
func BenchmarkGeo2TzTreeIndex_LookupZone(b *testing.B) {
	// load the database
	gsi, err := NewGeo2TzRTreeIndexFromGeoJSON("../tzdata/timezones.zip")
	assert.NoError(b, err)
	assert.NotEmpty(b, gsi.Size())

	// load the coordinates
	var tests []struct {
		Tz       string  `json:"tz"`
		Lat      float64 `json:"lat"`
		Lon      float64 `json:"lon"`
		NotFound bool    `json:"not_found,omitempty"`
	}
	err = helpers.LoadJSON("testdata/coordinates.json", &tests)
	assert.NoError(b, err)
	assert.NotEmpty(b, tests)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, tt := range tests {
			_, err := gsi.Lookup(tt.Lat, tt.Lon)
			if tt.NotFound {
				assert.ErrorIs(b, err, ErrNotFound)
				return
			}
			assert.NoError(b, err)
		}
	}
}

func Test_decodeJSON(t *testing.T) {

	expected := map[string][]int{
		"Africa/Bamako":    {29290},
		"America/New_York": {459, 31606, 17},
		"Asia/Tokyo":       {133, 129, 129, 139, 55, 127, 22, 17, 148, 18, 17, 162, 129, 129, 198, 424, 129, 33, 26, 92, 634, 754, 1019, 518, 149, 2408},
		"Australia/Sydney": {43621},
		"Europe/Rome":      {114, 137, 217, 51, 567, 53273, 238, 948},
	}

	zipFile, err := zip.OpenReader("testdata/timezones.zip")
	assert.NoError(t, err)

	iter := func(tz timezoneGeo) error {
		assert.Contains(t, expected, tz.Name)
		assert.Len(t, tz.Polygons, len(expected[tz.Name]))
		for i, p := range tz.Polygons {
			assert.Equal(t, expected[tz.Name][i], len(p.Vertices), "expected %d vertices, got %d", expected[tz.Name][i], len(p.Vertices))
		}
		return nil
	}

	err = decodeJSON(zipFile.File[0], iter)
	assert.NoError(t, err)

}
