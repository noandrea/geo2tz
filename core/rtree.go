package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"archive/zip"

	"github.com/tidwall/rtree"
)

type Geo2TzRTreeIndex struct {
	max_lookups int
	land        rtree.RTreeG[timezoneGeo]
	sea         rtree.RTreeG[timezoneGeo]
}

// IsOcean checks if the timezone is for oceans
func IsOcean(label string) bool {
	return strings.HasPrefix(label, "Etc/GMT")
}

// Insert adds a new timezone bounding box to the index
func (g *Geo2TzRTreeIndex) Insert(min, max [2]float64, element timezoneGeo) {
	if IsOcean(element.Name) {
		g.sea.Insert(min, max, element)
		return
	}
	g.land.Insert(min, max, element)
}

// NewGeo2TzRTreeIndexFromGeoJSON creates a new Geo2TzRTreeIndex from a GeoJSON file
func NewGeo2TzRTreeIndexFromGeoJSON(geoJSONPath string) (*Geo2TzRTreeIndex, error) {
	// open the zip file
	zipFile, err := zip.OpenReader(geoJSONPath)
	if err != nil {
		return nil, err
	}
	defer zipFile.Close()

	// create a new shape index
	gri := &Geo2TzRTreeIndex{
		max_lookups: 30,
	}

	// this function will add the timezone polygons to the shape index
	iter := func(tz *timezoneGeo) error {
		for _, p := range tz.Polygons {
			gri.Insert([2]float64{p.MinLat, p.MinLng}, [2]float64{p.MaxLat, p.MaxLng}, *tz)
		}
		return nil
	}
	// iterate over the zip file
	for _, v := range zipFile.File {
		if strings.EqualFold(".json", v.Name[len(v.Name)-5:]) {
			if err := decodeJSON(v, iter); err != nil {
				return nil, err
			}
		}
	}
	// build the shape index
	return gri, nil
}

// Lookup returns the timezone ID for a given latitude and longitude
// if the timezone is not found, it returns an error
// It first searches in the land index, if not found, it searches in the sea index
func (g *Geo2TzRTreeIndex) Lookup(lat, lon float64) (tzd TimeZoneData, err error) {
	tzd.Coords.Lat = lat
	tzd.Coords.Lon = lon

	lookup_num := 0
	// search the land index
	g.land.Search(
		[2]float64{lat, lon},
		[2]float64{lat, lon},
		func(min, max [2]float64, data timezoneGeo) bool {
			lookup_num++
			if lookup_num >= g.max_lookups {
				return false
			}
			for _, p := range data.Polygons {
				if isPointInPolygonPIP(vertex{lat, lon}, p) {
					tzd.TZ = data.Name
					return false
				}
			}
			return true
		},
	)

	if tzd.TZ == "" {
		// if not found, search the sea index
		lookup_num = 0
		g.sea.Search(
			[2]float64{lat, lon},
			[2]float64{lat, lon},
			func(min, max [2]float64, data timezoneGeo) bool {
				lookup_num++
				if lookup_num >= g.max_lookups {
					return false
				}
				for _, p := range data.Polygons {
					if isPointInPolygonPIP(vertex{lat, lon}, p) {
						tzd.TZ = data.Name
						return false
					}
				}
				return true
			},
		)
	}

	if tzd.TZ == "" {
		err = ErrNotFound
	}
	return
}

// isPointInPolygonPIP checks if a point is inside a polygon using the Point in Polygon algorithm
func isPointInPolygonPIP(point vertex, polygon polygon) bool {
	oddNodes := false
	n := len(polygon.Vertices)
	for i := 0; i < n; i++ {
		j := (i + 1) % n
		vi := polygon.Vertices[i]
		vj := polygon.Vertices[j]
		// Check if the point lies on an edge of the polygon (including horizontal)
		if (vi.lng == vj.lng && vi.lng == point.lng && point.lat >= min(vi.lat, vj.lat) && point.lat <= max(vi.lat, vj.lat)) ||
			((vi.lat < point.lat && point.lat <= vj.lat) || (vj.lat < point.lat && point.lat <= vi.lat)) &&
				(point.lng < (vj.lng-vi.lng)*(point.lat-vi.lat)/(vj.lat-vi.lat)+vi.lng) {
			oddNodes = !oddNodes
		}
	}
	return oddNodes
}

/*
GeoJSON processing
*/

type timezoneGeo struct {
	Name     string
	Polygons []polygon
}
type polygon struct {
	Vertices []vertex
	MaxLat   float64
	MinLat   float64
	MaxLng   float64
	MinLng   float64
}

func newPolygon() polygon {
	return polygon{
		Vertices: make([]vertex, 0),
		MaxLat:   -90,
		MinLat:   90,
		MaxLng:   -180,
		MinLng:   180,
	}
}

type vertex struct {
	lat, lng float64
}

func (p *polygon) AddVertex(lat, lng float64) {

	if lat > p.MaxLat {
		p.MaxLat = lat
	}
	if lat < p.MinLat {
		p.MinLat = lat
	}
	if lng > p.MaxLng {
		p.MaxLng = lng
	}
	if lng < p.MinLng {
		p.MinLng = lng
	}

	p.Vertices = append(p.Vertices, vertex{lat, lng})
}

func decodeJSON(f *zip.File, iter func(tz *timezoneGeo) error) (err error) {
	var rc io.ReadCloser
	if rc, err = f.Open(); err != nil {
		return err
	}
	defer rc.Close()

	dec := json.NewDecoder(rc)

	var token json.Token
	for dec.More() {
		if token, err = dec.Token(); err != nil {
			break
		}
		if t, ok := token.(string); ok && t == "features" {
			if token, err = dec.Token(); err == nil && token.(json.Delim) == '[' {
				return decodeFeatures(dec, iter) // decode features
			}
		}
	}
	return errors.New("error no features found")
}

func decodeFeatures(dec *json.Decoder, fn func(tz *timezoneGeo) error) error {
	var err error
	toPolygon := func(raw any) (polygon, error) {
		container, ok := raw.([]any)
		if !ok {
			return polygon{}, fmt.Errorf("invalid polygon data, expected[][]any, got %T", raw)
		}

		p := newPolygon()
		for _, c := range container {
			c, ok := c.([]any)
			if !ok {
				return p, fmt.Errorf("invalid container data, expected []any, got %T", c)
			}
			if len(c) != 2 {
				return p, fmt.Errorf("invalid point data, expected 2, got %v", len(c))
			}
			lat, ok := c[1].(float64)
			if !ok {
				return p, fmt.Errorf("invalid lat data, float64, got %T", c)
			}
			lng, ok := c[0].(float64)
			if !ok {
				return p, fmt.Errorf("invalid lng data, float64, got %T", c)
			}
			p.AddVertex(lat, lng)
		}
		return p, nil
	}

	var f struct {
		Type       string `json:"type"`
		Properties struct {
			TzID string `json:"tzid"`
		} `json:"properties"`
		Geometry struct {
			Item        string `json:"type"`
			Coordinates []any  `json:"coordinates"`
		} `json:"geometry"`
	}

	for dec.More() {
		if err = dec.Decode(&f); err != nil {
			return err
		}
		tg := &timezoneGeo{Name: f.Properties.TzID}
		switch f.Geometry.Item {
		case "Polygon":
			// we ignore the holes, that is why we only take the first block of coordinates
			p, err := toPolygon(f.Geometry.Coordinates[0])
			if err != nil {
				return err
			}
			tg.Polygons = []polygon{p}
		case "MultiPolygon":
			for _, multi := range f.Geometry.Coordinates {
				// we ignore the holes, that is why we only take the first block of coordinates
				p, err := toPolygon(multi.([]any)[0])
				if err != nil {
					return err
				}
				tg.Polygons = append(tg.Polygons, p)
			}
		}
		if err = fn(tg); err != nil {
			return err
		}
	}
	return nil
}
