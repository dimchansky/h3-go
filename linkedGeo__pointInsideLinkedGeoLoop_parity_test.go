//go:build cgo

package h3

import (
	"math"
	"testing"
)

func Test_pointInsideLinkedGeoLoop_parity(t *testing.T) {
	testCases := []struct {
		name  string
		loop  *LinkedGeoLoop
		bbox  *BBox
		coord *LatLng
	}{
		{
			name: "point inside simple square",
			loop: createTestLoop([]LatLng{
				{Lat: 0.0, Lng: 0.0},
				{Lat: 0.0, Lng: 1.0},
				{Lat: 1.0, Lng: 1.0},
				{Lat: 1.0, Lng: 0.0},
			}),
			bbox: &BBox{
				North: 1.0, South: 0.0,
				East: 1.0, West: 0.0,
			},
			coord: &LatLng{Lat: 0.5, Lng: 0.5},
		},
		{
			name: "point outside simple square",
			loop: createTestLoop([]LatLng{
				{Lat: 0.0, Lng: 0.0},
				{Lat: 0.0, Lng: 1.0},
				{Lat: 1.0, Lng: 1.0},
				{Lat: 1.0, Lng: 0.0},
			}),
			bbox: &BBox{
				North: 1.0, South: 0.0,
				East: 1.0, West: 0.0,
			},
			coord: &LatLng{Lat: 1.5, Lng: 0.5},
		},
		{
			name: "point on edge",
			loop: createTestLoop([]LatLng{
				{Lat: 0.0, Lng: 0.0},
				{Lat: 0.0, Lng: 1.0},
				{Lat: 1.0, Lng: 1.0},
				{Lat: 1.0, Lng: 0.0},
			}),
			bbox: &BBox{
				North: 1.0, South: 0.0,
				East: 1.0, West: 0.0,
			},
			coord: &LatLng{Lat: 0.0, Lng: 0.5},
		},
		{
			name: "point inside triangle",
			loop: createTestLoop([]LatLng{
				{Lat: 0.0, Lng: 0.0},
				{Lat: 1.0, Lng: 0.5},
				{Lat: 0.0, Lng: 1.0},
			}),
			bbox: &BBox{
				North: 1.0, South: 0.0,
				East: 1.0, West: 0.0,
			},
			coord: &LatLng{Lat: 0.3, Lng: 0.5},
		},
		{
			name: "point outside triangle",
			loop: createTestLoop([]LatLng{
				{Lat: 0.0, Lng: 0.0},
				{Lat: 1.0, Lng: 0.5},
				{Lat: 0.0, Lng: 1.0},
			}),
			bbox: &BBox{
				North: 1.0, South: 0.0,
				East: 1.0, West: 0.0,
			},
			coord: &LatLng{Lat: 0.8, Lng: 0.5},
		},
		{
			name: "point outside bbox - fast fail",
			loop: createTestLoop([]LatLng{
				{Lat: 0.0, Lng: 0.0},
				{Lat: 0.0, Lng: 1.0},
				{Lat: 1.0, Lng: 1.0},
				{Lat: 1.0, Lng: 0.0},
			}),
			bbox: &BBox{
				North: 1.0, South: 0.0,
				East: 1.0, West: 0.0,
			},
			coord: &LatLng{Lat: 2.0, Lng: 2.0},
		},
		{
			name: "complex polygon - inside",
			loop: createTestLoop([]LatLng{
				{Lat: 0.0, Lng: 0.0},
				{Lat: 0.0, Lng: 2.0},
				{Lat: 1.0, Lng: 2.0},
				{Lat: 1.0, Lng: 1.0},
				{Lat: 2.0, Lng: 1.0},
				{Lat: 2.0, Lng: 0.0},
			}),
			bbox: &BBox{
				North: 2.0, South: 0.0,
				East: 2.0, West: 0.0,
			},
			coord: &LatLng{Lat: 0.5, Lng: 0.5},
		},
		{
			name: "complex polygon - outside in concave area",
			loop: createTestLoop([]LatLng{
				{Lat: 0.0, Lng: 0.0},
				{Lat: 0.0, Lng: 2.0},
				{Lat: 1.0, Lng: 2.0},
				{Lat: 1.0, Lng: 1.0},
				{Lat: 2.0, Lng: 1.0},
				{Lat: 2.0, Lng: 0.0},
			}),
			bbox: &BBox{
				North: 2.0, South: 0.0,
				East: 2.0, West: 0.0,
			},
			coord: &LatLng{Lat: 1.5, Lng: 1.5},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test Go implementation
			goResult := pointInsideLinkedGeoLoop(tc.loop, tc.bbox, tc.coord)

			// Test C implementation
			cResult := pointInsideLinkedGeoLoopC(tc.loop, tc.bbox, tc.coord)

			if goResult != cResult {
				t.Errorf("Parity mismatch for %s: Go=%v, C=%v", tc.name, goResult, cResult)
			}
		})
	}
}

func Test_pointInsideLinkedGeoLoop_transmeridian(t *testing.T) {
	// Test transmeridian polygon (crosses antimeridian)
	transmeridianLoop := createTestLoop([]LatLng{
		{Lat: 0.0, Lng: math.Pi - 0.1},
		{Lat: 0.0, Lng: -math.Pi + 0.1},
		{Lat: 1.0, Lng: -math.Pi + 0.1},
		{Lat: 1.0, Lng: math.Pi - 0.1},
	})

	// Create bbox for transmeridian case
	var bbox BBox
	bboxFromLinkedGeoLoop(transmeridianLoop, &bbox)

	testCases := []struct {
		name  string
		coord *LatLng
	}{
		{
			name:  "point inside transmeridian polygon",
			coord: &LatLng{Lat: 0.5, Lng: math.Pi - 0.05},
		},
		{
			name:  "point outside transmeridian polygon",
			coord: &LatLng{Lat: 0.5, Lng: 0.0},
		},
		{
			name:  "point inside transmeridian polygon (negative side)",
			coord: &LatLng{Lat: 0.5, Lng: -math.Pi + 0.05},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test Go implementation
			goResult := pointInsideLinkedGeoLoop(transmeridianLoop, &bbox, tc.coord)

			// Test C implementation
			cResult := pointInsideLinkedGeoLoopC(transmeridianLoop, &bbox, tc.coord)

			if goResult != cResult {
				t.Errorf("Parity mismatch for %s: Go=%v, C=%v", tc.name, goResult, cResult)
			}
		})
	}
}
