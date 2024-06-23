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
	land rtree.RTreeG[string]
	sea  rtree.RTreeG[string]
	size int
}

// IsOcean checks if the timezone is for oceans
func IsOcean(label string) bool {
	return strings.HasPrefix(label, "Etc/GMT")
}

// Insert adds a new timezone bounding box to the index
func (g *Geo2TzRTreeIndex) Insert(min, max [2]float64, label string) {
	g.size++
	if IsOcean(label) {
		g.sea.Insert(min, max, label)
		return
	}
	g.land.Insert(min, max, label)
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
			minLat, minLng, maxLat, maxLng := p.Vertices[0].lat, p.Vertices[0].lng, p.Vertices[0].lat, p.Vertices[0].lng
			for _, v := range p.Vertices {
				if v.lng < minLng {
					minLng = v.lng
				}
				if v.lng > maxLng {
					maxLng = v.lng
				}
				if v.lat < minLat {
					minLat = v.lat
				}
				if v.lat > maxLat {
					maxLat = v.lat
				}
			}
			gri.Insert([2]float64{minLat, minLng}, [2]float64{maxLat, maxLng}, tz.Name)
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

	g.land.Search(
		[2]float64{lat, lng},
		[2]float64{lat, lng},
		func(min, max [2]float64, label string) bool {
			tzID = label
			return true
		},
	)

	if tzID == "" {
		g.sea.Search(
			[2]float64{lat, lng},
			[2]float64{lat, lng},
			func(min, max [2]float64, label string) bool {
				tzID = label
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

/*
GeoJSON processing
*/

// Polygon represents a polygon
// with a list of vertices [lat, lng]
type polygon struct {
	Vertices []vertex
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
