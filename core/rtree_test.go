package core

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
			assert.Equal(t, got.TZ, tt.Tz, "expected %s to be %s for https://www.google.com/maps/@%v,%v,12z", tt.Tz, got, tt.Lat, tt.Lon)
		})
	}
}

// benchmark the lookup function
func BenchmarkGeo2TzTreeIndex_LookupZone(b *testing.B) {
	// load the database
	gsi, err := NewGeo2TzRTreeIndexFromGeoJSON("../tzdata/timezones.zip")
	assert.NoError(b, err)

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

	// lis of expected polygons with the number of vertices
	expected := map[string][]int{
		"Africa/Bamako":    {29290},
		"America/New_York": {459, 31606, 17},
		"Asia/Tokyo":       {133, 129, 129, 139, 55, 127, 22, 17, 148, 18, 17, 162, 129, 129, 198, 424, 129, 33, 26, 92, 634, 754, 1019, 518, 149, 2408},
		"Australia/Sydney": {43621},
		"Europe/Berlin":    {151051, 5, 20, 150, 60, 605, 341},
		"Europe/Rome":      {114, 137, 217, 51, 567, 53273},
	}

	got := map[string][]int{}

	zipFile, err := zip.OpenReader("testdata/timezones.zip")
	assert.NoError(t, err)

	iter := func(tz *timezoneGeo) error {
		got[tz.Name] = []int{}
		for _, p := range tz.Polygons {
			got[tz.Name] = append(got[tz.Name], len(p.Vertices))
		}
		return nil
	}

	err = decodeJSON(zipFile.File[0], iter)
	assert.NoError(t, err)

	for expectedTz, expectedPolySizes := range expected {
		gotPolySizes, ok := got[expectedTz]
		assert.True(t, ok, "expected %s to be in the got map", expectedTz)
		areEq := assert.Equal(t, len(expectedPolySizes), len(gotPolySizes), "expected %s to have %d polygons, got %d", expectedTz, len(expectedPolySizes), len(gotPolySizes))

		if !areEq {
			continue
		}
		for i := range expectedPolySizes {
			assert.Equal(t, expectedPolySizes[i], gotPolySizes[i], "expected %s to have polygon %d with %d vertices, got %d", expectedTz, i, expectedPolySizes[i], gotPolySizes[i])
		}
	}
}
