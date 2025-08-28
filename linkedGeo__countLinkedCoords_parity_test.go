//go:build cgo

package h3

import (
	"testing"
)

func Test_countLinkedCoords_parity(t *testing.T) {
	testCases := []struct {
		name  string
		loop  *LinkedGeoLoop
		count int32
	}{
		{
			name: "empty loop",
			loop: &LinkedGeoLoop{
				First: nil,
				Last:  nil,
				Next:  nil,
			},
			count: 0,
		},
		{
			name: "single coordinate",
			loop: &LinkedGeoLoop{
				First: &LinkedLatLng{
					Vertex: LatLng{Lat: 37.775, Lng: -122.418},
					Next:   nil,
				},
				Last: nil, // Will be set to first
				Next: nil,
			},
			count: 1,
		},
		{
			name: "two coordinates",
			loop: createTestLoop([]LatLng{
				{Lat: 37.775, Lng: -122.418},
				{Lat: 37.776, Lng: -122.419},
			}),
			count: 2,
		},
		{
			name: "five coordinates",
			loop: createTestLoop([]LatLng{
				{Lat: 37.775, Lng: -122.418},
				{Lat: 37.776, Lng: -122.419},
				{Lat: 37.777, Lng: -122.420},
				{Lat: 37.778, Lng: -122.421},
				{Lat: 37.779, Lng: -122.422},
			}),
			count: 5,
		},
		{
			name:  "large loop",
			loop:  createTestLoop(make([]LatLng, 100)), // 100 default-initialized coordinates
			count: 100,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set up single coordinate case properly
			if tc.count == 1 && tc.loop.First != nil {
				tc.loop.Last = tc.loop.First
			}

			// Test Go implementation
			goResult := countLinkedCoords(tc.loop)

			// Test C implementation (skip empty loop as it would cause issues with C malloc/free)
			if tc.count > 0 {
				cResult := countLinkedCoordsC(tc.loop)

				if goResult != cResult {
					t.Errorf("Parity mismatch for %s: Go=%d, C=%d", tc.name, goResult, cResult)
				}
			}

			// Verify expected count
			if goResult != tc.count {
				t.Errorf("Go result mismatch for %s: expected %d, got %d", tc.name, tc.count, goResult)
			}
		})
	}
}

// Helper function to create a test loop from coordinates
func createTestLoop(coords []LatLng) *LinkedGeoLoop {
	if len(coords) == 0 {
		return &LinkedGeoLoop{First: nil, Last: nil, Next: nil}
	}

	first := &LinkedLatLng{
		Vertex: coords[0],
		Next:   nil,
	}

	current := first
	for i := 1; i < len(coords); i++ {
		next := &LinkedLatLng{
			Vertex: coords[i],
			Next:   nil,
		}
		current.Next = next
		current = next
	}

	return &LinkedGeoLoop{
		First: first,
		Last:  current,
		Next:  nil,
	}
}
