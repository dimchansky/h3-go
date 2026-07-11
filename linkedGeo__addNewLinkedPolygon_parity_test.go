//go:build cgo && c2go

package h3

import (
	"testing"
)

func Test_addNewLinkedPolygon_parity(t *testing.T) {
	testCases := []struct {
		name    string
		polygon *linkedGeoPolygon
	}{
		{
			name: "empty polygon",
			polygon: &linkedGeoPolygon{
				First: nil,
				Last:  nil,
				Next:  nil,
			},
		},
		{
			name: "polygon with one loop",
			polygon: &linkedGeoPolygon{
				First: &linkedGeoLoop{
					First: &linkedLatLng{
						Vertex: LatLng{Lat: 0.0, Lng: 0.0},
						Next:   nil,
					},
					Last: &linkedLatLng{
						Vertex: LatLng{Lat: 0.0, Lng: 0.0},
						Next:   nil,
					},
					Next: nil,
				},
				Last: &linkedGeoLoop{
					First: &linkedLatLng{
						Vertex: LatLng{Lat: 0.0, Lng: 0.0},
						Next:   nil,
					},
					Last: &linkedLatLng{
						Vertex: LatLng{Lat: 0.0, Lng: 0.0},
						Next:   nil,
					},
					Next: nil,
				},
				Next: nil,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test Go implementation
			goResult := addNewLinkedPolygon(tc.polygon)

			// Verify the Go implementation behavior
			if goResult == nil {
				t.Errorf("Go implementation returned nil")
			}
			if tc.polygon.Next != goResult {
				t.Errorf("Go implementation: polygon.Next should point to new polygon")
			}
			if goResult.First != nil || goResult.Last != nil || goResult.Next != nil {
				t.Errorf("Go implementation: new polygon should be initialized with nil fields")
			}

			// Reset for C test by creating a fresh polygon
			freshPolygon := &linkedGeoPolygon{
				First: tc.polygon.First,
				Last:  tc.polygon.Last,
				Next:  nil, // Reset Next for C test
			}

			// Test C implementation
			cResult := addNewLinkedPolygonC(freshPolygon)

			// Verify C implementation returns a valid result
			if cResult == nil {
				t.Errorf("C implementation returned nil")
			}

			// Note: We can't directly compare pointer equality between Go and C results
			// since they're in different memory spaces, but we can verify behavior
			if cResult.First != nil || cResult.Last != nil || cResult.Next != nil {
				t.Errorf("C implementation: new polygon should be initialized with nil fields")
			}
		})
	}
}

func Test_addNewLinkedPolygon_chaining(t *testing.T) {
	// Test that we can chain multiple polygons
	root := &linkedGeoPolygon{
		First: nil,
		Last:  nil,
		Next:  nil,
	}

	// Add first polygon
	second := addNewLinkedPolygon(root)
	if root.Next != second {
		t.Errorf("First addition: root.Next should point to second")
	}

	// Add third polygon
	third := addNewLinkedPolygon(second)
	if second.Next != third {
		t.Errorf("Second addition: second.Next should point to third")
	}

	// Verify chain structure
	if root.Next != second {
		t.Errorf("Chain structure: root should point to second")
	}
	if second.Next != third {
		t.Errorf("Chain structure: second should point to third")
	}
	if third.Next != nil {
		t.Errorf("Chain structure: third.Next should be nil")
	}
}

func Test_addNewLinkedPolygon_panic_on_non_nil_next(t *testing.T) {
	// Test that function panics when Next is not nil
	polygon := &linkedGeoPolygon{
		First: nil,
		Last:  nil,
		Next:  &linkedGeoPolygon{}, // Non-nil Next should cause panic
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic when polygon.Next is not nil")
		}
	}()

	addNewLinkedPolygon(polygon)
}
