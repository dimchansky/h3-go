//go:build cgo

package h3

import (
	"math"
	"testing"
)

func Test_bboxFromLinkedGeoLoop_parity(t *testing.T) {
	testCases := []struct {
		name string
		loop *LinkedGeoLoop
	}{
		{
			name: "empty loop",
			loop: &LinkedGeoLoop{
				First: nil,
				Last:  nil,
				Next:  nil,
			},
		},
		{
			name: "single coordinate",
			loop: createTestLoop([]LatLng{
				{Lat: 0.5, Lng: 1.0},
			}),
		},
		{
			name: "simple rectangle",
			loop: createTestLoop([]LatLng{
				{Lat: 0.0, Lng: 0.0},
				{Lat: 0.0, Lng: 1.0},
				{Lat: 1.0, Lng: 1.0},
				{Lat: 1.0, Lng: 0.0},
			}),
		},
		{
			name: "triangle",
			loop: createTestLoop([]LatLng{
				{Lat: 0.0, Lng: 0.0},
				{Lat: 1.0, Lng: 0.5},
				{Lat: 0.0, Lng: 1.0},
			}),
		},
		{
			name: "coordinates with negative values",
			loop: createTestLoop([]LatLng{
				{Lat: -1.0, Lng: -1.0},
				{Lat: -1.0, Lng: 1.0},
				{Lat: 1.0, Lng: 1.0},
				{Lat: 1.0, Lng: -1.0},
			}),
		},
		{
			name: "transmeridian case",
			loop: createTestLoop([]LatLng{
				{Lat: 0.0, Lng: math.Pi - 0.1},
				{Lat: 0.0, Lng: -math.Pi + 0.1},
				{Lat: 1.0, Lng: -math.Pi + 0.1},
				{Lat: 1.0, Lng: math.Pi - 0.1},
			}),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test Go implementation
			var goBbox BBox
			bboxFromLinkedGeoLoop(tc.loop, &goBbox)

			// Test C implementation (skip empty loop as it would cause issues with C malloc/free)
			if tc.loop.First != nil {
				var cBbox BBox
				bboxFromLinkedGeoLoopC(tc.loop, &cBbox)

				const tolerance = 1e-15

				if math.Abs(goBbox.North.Rad()-cBbox.North.Rad()) > tolerance {
					t.Errorf("North mismatch for %s: Go=%v, C=%v", tc.name, goBbox.North.Rad(), cBbox.North.Rad())
				}
				if math.Abs(goBbox.South.Rad()-cBbox.South.Rad()) > tolerance {
					t.Errorf("South mismatch for %s: Go=%v, C=%v", tc.name, goBbox.South.Rad(), cBbox.South.Rad())
				}
				if math.Abs(goBbox.East.Rad()-cBbox.East.Rad()) > tolerance {
					t.Errorf("East mismatch for %s: Go=%v, C=%v", tc.name, goBbox.East.Rad(), cBbox.East.Rad())
				}
				if math.Abs(goBbox.West.Rad()-cBbox.West.Rad()) > tolerance {
					t.Errorf("West mismatch for %s: Go=%v, C=%v", tc.name, goBbox.West.Rad(), cBbox.West.Rad())
				}
			}

			// Verify expected properties for specific test cases
			switch tc.name {
			case "empty loop":
				if goBbox != (BBox{}) {
					t.Errorf("Empty loop should produce zero bbox, got %+v", goBbox)
				}
			case "simple rectangle":
				if goBbox.South != 0.0 || goBbox.North != 1.0 || goBbox.West != 0.0 || goBbox.East != 1.0 {
					t.Errorf("Rectangle bbox mismatch: expected [0,0,1,1], got [%v,%v,%v,%v]",
						goBbox.South, goBbox.West, goBbox.North, goBbox.East)
				}
			}
		})
	}
}
