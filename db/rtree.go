package db

import (
	"encoding/json"
	"errors"
	"io"
	"strings"

	"github.com/klauspost/compress/zip"
	"github.com/tidwall/rtree"
)

type Geo2TzRTreeIndex struct {
	land rtree.RTreeG[timezoneGeo]
	sea  rtree.RTreeG[timezoneGeo]
	size int
}

// IsOcean checks if the timezone is for oceans
func IsOcean(label string) bool {
	return strings.HasPrefix(label, "Etc/GMT")
}

// Insert adds a new timezone bounding box to the index
func (g *Geo2TzRTreeIndex) Insert(min, max [2]float64, element timezoneGeo) {
	g.size++
	if IsOcean(element.Name) {
		g.sea.Insert(min, max, element)
		return
	}
	g.land.Insert(min, max, element)
}

func NewGeo2TzRTreeIndexFromGeoJSON(geoJSONPath string) (*Geo2TzRTreeIndex, error) {
	// open the zip file
	zipFile, err := zip.OpenReader(geoJSONPath)
	if err != nil {
		return nil, err
	}
	defer zipFile.Close()

	// create a new shape index
	gri := &Geo2TzRTreeIndex{}

	// this function will add the timezone polygons to the shape index
	iter := func(tz timezoneGeo) error {
		for _, p := range tz.Polygons {
			gri.Insert([2]float64{p.MinLat, p.MinLng}, [2]float64{p.MaxLat, p.MaxLng}, tz)
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
func (g *Geo2TzRTreeIndex) Lookup(lat, lng float64) (tzID string, err error) {
	// search the land index
	g.land.Search(
		[2]float64{lat, lng},
		[2]float64{lat, lng},
		func(min, max [2]float64, data timezoneGeo) bool {
			for _, p := range data.Polygons {
				if isPointInPolygonPIP(vertex{lat, lng}, p) {
					tzID = data.Name
					return false
				}
			}
			return true
		},
	)

	if tzID == "" {
		// if not found, search the sea index
		g.sea.Search(
			[2]float64{lat, lng},
			[2]float64{lat, lng},
			func(min, max [2]float64, data timezoneGeo) bool {
				for _, p := range data.Polygons {
					if isPointInPolygonPIP(vertex{lat, lng}, p) {
						tzID = data.Name
						return false
					}
				}
				return true
			},
		)
	}

	if tzID == "" {
		err = ErrNotFound
	}
	return
}

func (g *Geo2TzRTreeIndex) Size() int {
	return g.size
}

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

// Polygon represents a polygon
// with a list of vertices [lat, lng]
type polygon struct {
	Vertices []vertex
	MaxLat   float64
	MinLat   float64
	MaxLng   float64
	MinLng   float64
}

type vertex struct {
	lat, lng float64
}

type GeoJSONFeature struct {
	Type       string `json:"type"`
	Properties struct {
		TzID string `json:"tzid"`
	} `json:"properties"`
	Geometry struct {
		Item        string        `json:"type"`
		Coordinates []interface{} `json:"coordinates"`
	} `json:"geometry"`
}

func (p *polygon) AddVertex(lat, lng float64) {
	if len(p.Vertices) == 0 {
		p.MaxLat = lat
		p.MinLat = lat
		p.MaxLng = lng
		p.MinLng = lng
	} else {
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
	}
	p.Vertices = append(p.Vertices, vertex{lat, lng})
}

type timezoneGeo struct {
	Name     string
	Polygons []polygon
}

func decodeJSON(f *zip.File, iter func(tz timezoneGeo) error) (err error) {
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

func decodeFeatures(dec *json.Decoder, fn func(tz timezoneGeo) error) error {
	var f GeoJSONFeature
	var err error

	for dec.More() {
		if err = dec.Decode(&f); err != nil {
			return err
		}
		var pp []polygon
		switch f.Geometry.Item {
		case "Polygon":
			pp = decodePolygons(f.Geometry.Coordinates)
		case "MultiPolygon":
			pp = decodeMultiPolygons(f.Geometry.Coordinates)
		}
		if err = fn(timezoneGeo{Name: f.Properties.TzID, Polygons: pp}); err != nil {
			return err
		}
	}

	return nil
}

// decodePolygons
// GeoJSON Spec https://geojson.org/geojson-spec.html
// Coordinates: [Longitude, Latitude]
func decodePolygons(polygons []interface{}) []polygon {
	var pp []polygon
	for _, points := range polygons {
		p := polygon{}
		for _, i := range points.([]interface{}) {
			if latlng, ok := i.([]interface{}); ok {
				p.AddVertex(latlng[1].(float64), latlng[0].(float64))
			}
		}
		pp = append(pp, p)
	}
	return pp
}

// decodeMultiPolygons
// GeoJSON Spec https://geojson.org/geojson-spec.html
// Coordinates: [Longitude, Latitude]
func decodeMultiPolygons(polygons []interface{}) []polygon {
	var pp []polygon
	for _, v := range polygons {
		p := polygon{}
		for _, points := range v.([]interface{}) { // 2
			for _, i := range points.([]interface{}) {
				if latlng, ok := i.([]interface{}); ok {
					p.AddVertex(latlng[1].(float64), latlng[0].(float64))
				}
			}
		}
		pp = append(pp, p)
	}
	return pp
}
