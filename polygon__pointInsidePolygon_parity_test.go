//go:build cgo

package h3

import "testing"

func Test_pointInsidePolygon_ParityWithC(t *testing.T) {
	outer := GeoLoop{{Lat: 0, Lng: 0}, {Lat: 0, Lng: 2}, {Lat: 2, Lng: 2}, {Lat: 2, Lng: 0}}
	hole := GeoLoop{{Lat: 0.5, Lng: 0.5}, {Lat: 0.5, Lng: 1.5}, {Lat: 1.5, Lng: 1.5}, {Lat: 1.5, Lng: 0.5}}
	poly := GeoPolygon{GeoLoop: outer, Holes: []GeoLoop{hole}}
	bboxes := make([]BBox, 1+len(poly.Holes))
	bboxFromGeoLoop(poly.GeoLoop, &bboxes[0])
	for i := range poly.Holes {
		bboxFromGeoLoop(poly.Holes[i], &bboxes[i+1])
	}
	pts := []LatLng{{Lat: 0.25, Lng: 0.25}, {Lat: 1.0, Lng: 1.0}, {Lat: 3, Lng: 3}}
	for _, p := range pts {
		goVal := pointInsidePolygon(poly, bboxes, &p)
		cVal := pointInsidePolygonC(poly, bboxes, p)
		if goVal != cVal {
			t.Fatalf("pointInsidePolygon mismatch for p=%+v: go=%v c=%v", p, goVal, cVal)
		}
	}
}
