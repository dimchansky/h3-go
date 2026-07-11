//go:build cgo && c2go

package h3

import (
	"math"
	"testing"
)

func Test_isClockwiseLinkedGeoLoop_parity(t *testing.T) {
	testCases := []struct {
		name string
		loop *LinkedGeoLoop
	}{
		{
			name: "clockwise square",
			loop: createTestLoop([]LatLng{
				{Lat: 0.0, Lng: 0.0},
				{Lat: 1.0, Lng: 0.0},
				{Lat: 1.0, Lng: 1.0},
				{Lat: 0.0, Lng: 1.0},
			}),
		},
		{
			name: "counter-clockwise square",
			loop: createTestLoop([]LatLng{
				{Lat: 0.0, Lng: 0.0},
				{Lat: 0.0, Lng: 1.0},
				{Lat: 1.0, Lng: 1.0},
				{Lat: 1.0, Lng: 0.0},
			}),
		},
		{
			name: "clockwise triangle",
			loop: createTestLoop([]LatLng{
				{Lat: 0.0, Lng: 0.0},
				{Lat: 1.0, Lng: 0.0},
				{Lat: 0.5, Lng: 1.0},
			}),
		},
		{
			name: "counter-clockwise triangle",
			loop: createTestLoop([]LatLng{
				{Lat: 0.0, Lng: 0.0},
				{Lat: 0.5, Lng: 1.0},
				{Lat: 1.0, Lng: 0.0},
			}),
		},
		{
			name: "single point",
			loop: createTestLoop([]LatLng{
				{Lat: 0.0, Lng: 0.0},
			}), // Single point has no winding
		},
		{
			name: "two points",
			loop: createTestLoop([]LatLng{
				{Lat: 0.0, Lng: 0.0},
				{Lat: 1.0, Lng: 1.0},
			}), // Two points form no area
		},
		{
			name: "clockwise rectangle",
			loop: createTestLoop([]LatLng{
				{Lat: 0.0, Lng: 0.0},
				{Lat: 2.0, Lng: 0.0},
				{Lat: 2.0, Lng: 1.0},
				{Lat: 0.0, Lng: 1.0},
			}),
		},
		{
			name: "counter-clockwise rectangle",
			loop: createTestLoop([]LatLng{
				{Lat: 0.0, Lng: 0.0},
				{Lat: 0.0, Lng: 1.0},
				{Lat: 2.0, Lng: 1.0},
				{Lat: 2.0, Lng: 0.0},
			}),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test Go implementation
			goResult := isClockwiseLinkedGeoLoop(tc.loop)

			// Test C implementation (skip very simple cases that may cause C issues)
			if len(getLoopCoords(tc.loop)) >= 3 {
				cResult := isClockwiseLinkedGeoLoopC(tc.loop)

				if goResult != cResult {
					t.Errorf("Parity mismatch for %s: Go=%v, C=%v", tc.name, goResult, cResult)
				}
			}

		})
	}
}

func Test_isClockwiseLinkedGeoLoop_transmeridian(t *testing.T) {
	// Test transmeridian polygon (crosses antimeridian)
	transmeridianClockwise := createTestLoop([]LatLng{
		{Lat: 0.0, Lng: math.Pi - 0.1},
		{Lat: 1.0, Lng: math.Pi - 0.1},
		{Lat: 1.0, Lng: -math.Pi + 0.1},
		{Lat: 0.0, Lng: -math.Pi + 0.1},
	})

	transmeridianCounterClockwise := createTestLoop([]LatLng{
		{Lat: 0.0, Lng: math.Pi - 0.1},
		{Lat: 0.0, Lng: -math.Pi + 0.1},
		{Lat: 1.0, Lng: -math.Pi + 0.1},
		{Lat: 1.0, Lng: math.Pi - 0.1},
	})

	testCases := []struct {
		name string
		loop *LinkedGeoLoop
	}{
		{
			name: "transmeridian clockwise",
			loop: transmeridianClockwise,
		},
		{
			name: "transmeridian counter-clockwise",
			loop: transmeridianCounterClockwise,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test Go implementation
			goResult := isClockwiseLinkedGeoLoop(tc.loop)

			// Test C implementation
			cResult := isClockwiseLinkedGeoLoopC(tc.loop)

			if goResult != cResult {
				t.Errorf("Parity mismatch for %s: Go=%v, C=%v", tc.name, goResult, cResult)
			}

			// Log results for inspection (transmeridian winding can be complex)
			t.Logf("%s: Go=%v, C=%v", tc.name, goResult, cResult)
		})
	}
}

// Helper function to get coordinates from a loop for length checking
func getLoopCoords(loop *LinkedGeoLoop) []LatLng {
	var coords []LatLng
	current := loop.First
	for current != nil {
		coords = append(coords, current.Vertex)
		current = current.Next
		if current == loop.First {
			break
		}
	}
	return coords
}
