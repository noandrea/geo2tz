package db

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/tidwall/rtree"
)

const (
	maxLookupAttempts = 30
	coordinateLen     = 2
	minLat            = -90.0
	maxLat            = 90.0
	minLng            = -180.0
	maxLng            = 180.0
)

type Geo2TzRTreeIndex struct {
	maxLookups int
	land       rtree.RTreeG[timezoneGeo]
	sea        rtree.RTreeG[timezoneGeo]
}

// IsOcean checks if the timezone is for oceans
func IsOcean(label string) bool {
	return strings.HasPrefix(label, "Etc/GMT")
}

// Insert adds a new timezone bounding box to the index
func (g *Geo2TzRTreeIndex) Insert(bboxMin, bboxMax [2]float64, element timezoneGeo) {
	if IsOcean(element.Name) {
		g.sea.Insert(bboxMin, bboxMax, element)
		return
	}
	g.land.Insert(bboxMin, bboxMax, element)
}

// NewGeo2TzRTreeIndexFromGeoJSON creates a new Geo2TzRTreeIndex from a GeoJSON file
func NewGeo2TzRTreeIndexFromGeoJSON(geoJSONPath string) (*Geo2TzRTreeIndex, error) {
	// open the zip file
	zipFile, err := zip.OpenReader(geoJSONPath)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := zipFile.Close(); closeErr != nil {
			fmt.Println("Error closing zip file:", closeErr)
		}
	}()

	// create a new shape index
	gri := &Geo2TzRTreeIndex{
		maxLookups: maxLookupAttempts,
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
			if err = decodeJSON(v, iter); err != nil {
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

	lookupNum := 0
	// search the land index
	g.land.Search(
		[2]float64{lat, lng},
		[2]float64{lat, lng},
		func(_, _ [2]float64, data timezoneGeo) bool {
			lookupNum++
			if lookupNum >= g.maxLookups {
				return false
			}
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
		lookupNum = 0
		g.sea.Search(
			[2]float64{lat, lng},
			[2]float64{lat, lng},
			func(_, _ [2]float64, data timezoneGeo) bool {
				lookupNum++
				if lookupNum >= g.maxLookups {
					return false
				}
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

// isPointInPolygonPIP checks if a point is inside a polygon using the Point in Polygon algorithm
func isPointInPolygonPIP(point vertex, polygon polygon) bool {
	oddNodes := false
	n := len(polygon.Vertices)
	for i := 0; i < n; i++ {
		j := (i + 1) % n
		vi := polygon.Vertices[i]
		vj := polygon.Vertices[j]
		// Check if the point lies on an edge of the polygon (including horizontal)
		verticalEdge := vi.lng == vj.lng
		onEdgeLng := vi.lng == point.lng
		withinEdgeLat := point.lat >= min(vi.lat, vj.lat) && point.lat <= max(vi.lat, vj.lat)
		crossesEdge := (vi.lat < point.lat && point.lat <= vj.lat) || (vj.lat < point.lat && point.lat <= vi.lat)
		intersectLng := point.lng < (vj.lng-vi.lng)*(point.lat-vi.lat)/(vj.lat-vi.lat)+vi.lng
		if (verticalEdge && onEdgeLng && withinEdgeLat) || (crossesEdge && intersectLng) {
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
		MaxLat:   minLat,
		MinLat:   maxLat,
		MaxLng:   minLng,
		MinLng:   maxLng,
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
	defer func() {
		if closeErr := rc.Close(); closeErr != nil {
			fmt.Println("Error closing read closer:", closeErr)
		}
	}()

	dec := json.NewDecoder(rc)

	var token json.Token
	for dec.More() {
		if token, err = dec.Token(); err != nil {
			break
		}
		if t, isFeatures := token.(string); isFeatures && t == "features" {
			if token, err = dec.Token(); err == nil {
				if delim, isList := token.(json.Delim); isList && delim == '[' {
					return decodeFeatures(dec, iter) // decode features
				}
			}
		}
	}
	return errors.New("error no features found")
}

func decodeFeatures(dec *json.Decoder, fn func(tz *timezoneGeo) error) (err error) {
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
		var polygons []polygon
		if polygons, err = geometryToPolygons(f.Geometry.Item, f.Geometry.Coordinates); err != nil {
			return err
		}
		if err = fn(&timezoneGeo{Name: f.Properties.TzID, Polygons: polygons}); err != nil {
			return err
		}
	}
	return nil
}

// geometryToPolygons converts a geojson geometry into the list of polygons of a timezone,
// the holes in the geometries are ignored
func geometryToPolygons(geometryType string, coordinates []any) ([]polygon, error) {
	switch geometryType {
	case "Polygon":
		p, err := coordinatesToPolygon(coordinates[0])
		if err != nil {
			return nil, err
		}
		return []polygon{p}, nil
	case "MultiPolygon":
		polygons := make([]polygon, 0, len(coordinates))
		for _, multi := range coordinates {
			multiPolygon, ok := multi.([]any)
			if !ok {
				return nil, fmt.Errorf("invalid multipolygon data, expected []any, got %T", multi)
			}
			p, pErr := coordinatesToPolygon(multiPolygon[0])
			if pErr != nil {
				return nil, pErr
			}
			polygons = append(polygons, p)
		}
		return polygons, nil
	default:
		return nil, nil
	}
}

func coordinatesToPolygon(raw any) (polygon, error) {
	container, isList := raw.([]any)
	if !isList {
		return polygon{}, fmt.Errorf("invalid polygon data, expected[][]any, got %T", raw)
	}

	p := newPolygon()
	for _, item := range container {
		coordinates, ok := item.([]any)
		if !ok {
			return p, fmt.Errorf("invalid container data, expected []any, got %T", item)
		}
		if len(coordinates) != coordinateLen {
			return p, fmt.Errorf("invalid point data, expected %d, got %v", coordinateLen, len(coordinates))
		}
		lat, ok := coordinates[1].(float64)
		if !ok {
			return p, fmt.Errorf("invalid lat data, float64, got %T", coordinates)
		}
		lng, ok := coordinates[0].(float64)
		if !ok {
			return p, fmt.Errorf("invalid lng data, float64, got %T", coordinates)
		}
		p.AddVertex(lat, lng)
	}
	return p, nil
}
