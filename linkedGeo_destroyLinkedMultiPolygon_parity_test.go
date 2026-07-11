//go:build cgo && c2go

package h3

import (
	"testing"
)

func TestDestroyLinkedMultiPolygonParity(t *testing.T) {
	tests := []struct {
		name string
		// Test case description
		desc string
		// Function to create the test polygon
		createPolygon func() *linkedGeoPolygon
		// Function to verify the post-destruction state
		verify func(*testing.T, *linkedGeoPolygon)
	}{
		{
			name: "nil_polygon",
			desc: "Test destroying nil polygon - should handle gracefully",
			createPolygon: func() *linkedGeoPolygon {
				return nil
			},
			verify: func(t *testing.T, polygon *linkedGeoPolygon) {
				// Nothing to verify for nil input
			},
		},
		{
			name: "empty_polygon",
			desc: "Test destroying empty polygon with no loops",
			createPolygon: func() *linkedGeoPolygon {
				return &linkedGeoPolygon{
					First: nil,
					Last:  nil,
					Next:  nil,
				}
			},
			verify: func(t *testing.T, polygon *linkedGeoPolygon) {
				if polygon == nil {
					t.Error("Input polygon should not be nil after destruction")
					return
				}
				if polygon.First != nil {
					t.Error("Expected polygon.First to be nil after destruction")
				}
				if polygon.Last != nil {
					t.Error("Expected polygon.Last to be nil after destruction")
				}
				if polygon.Next != nil {
					t.Error("Expected polygon.Next to be nil after destruction")
				}
			},
		},
		{
			name: "single_polygon_single_loop",
			desc: "Test destroying single polygon with one loop containing coordinates",
			createPolygon: func() *linkedGeoPolygon {
				// Create coordinates
				coord1 := &linkedLatLng{Vertex: LatLng{Lat: 0.1, Lng: 0.2}, Next: nil}
				coord2 := &linkedLatLng{Vertex: LatLng{Lat: 0.3, Lng: 0.4}, Next: nil}
				coord3 := &linkedLatLng{Vertex: LatLng{Lat: 0.5, Lng: 0.6}, Next: nil}

				// Link coordinates
				coord1.Next = coord2
				coord2.Next = coord3

				// Create loop
				loop := &linkedGeoLoop{
					First: coord1,
					Last:  coord3,
					Next:  nil,
				}

				// Create polygon
				return &linkedGeoPolygon{
					First: loop,
					Last:  loop,
					Next:  nil,
				}
			},
			verify: func(t *testing.T, polygon *linkedGeoPolygon) {
				if polygon == nil {
					t.Error("Input polygon should not be nil after destruction")
					return
				}
				if polygon.First != nil {
					t.Error("Expected polygon.First to be nil after destruction")
				}
				if polygon.Last != nil {
					t.Error("Expected polygon.Last to be nil after destruction")
				}
				if polygon.Next != nil {
					t.Error("Expected polygon.Next to be nil after destruction")
				}
			},
		},
		{
			name: "single_polygon_multiple_loops",
			desc: "Test destroying single polygon with multiple loops",
			createPolygon: func() *linkedGeoPolygon {
				// Create first loop coordinates
				coord1a := &linkedLatLng{Vertex: LatLng{Lat: 0.1, Lng: 0.2}, Next: nil}
				coord1b := &linkedLatLng{Vertex: LatLng{Lat: 0.3, Lng: 0.4}, Next: nil}
				coord1a.Next = coord1b

				loop1 := &linkedGeoLoop{
					First: coord1a,
					Last:  coord1b,
					Next:  nil,
				}

				// Create second loop coordinates
				coord2a := &linkedLatLng{Vertex: LatLng{Lat: 0.5, Lng: 0.6}, Next: nil}
				coord2b := &linkedLatLng{Vertex: LatLng{Lat: 0.7, Lng: 0.8}, Next: nil}
				coord2a.Next = coord2b

				loop2 := &linkedGeoLoop{
					First: coord2a,
					Last:  coord2b,
					Next:  nil,
				}

				// Link loops
				loop1.Next = loop2

				// Create polygon
				return &linkedGeoPolygon{
					First: loop1,
					Last:  loop2,
					Next:  nil,
				}
			},
			verify: func(t *testing.T, polygon *linkedGeoPolygon) {
				if polygon == nil {
					t.Error("Input polygon should not be nil after destruction")
					return
				}
				if polygon.First != nil {
					t.Error("Expected polygon.First to be nil after destruction")
				}
				if polygon.Last != nil {
					t.Error("Expected polygon.Last to be nil after destruction")
				}
				if polygon.Next != nil {
					t.Error("Expected polygon.Next to be nil after destruction")
				}
			},
		},
		{
			name: "multiple_polygons",
			desc: "Test destroying linked list of multiple polygons - only first should be preserved",
			createPolygon: func() *linkedGeoPolygon {
				// Create first polygon
				coord1 := &linkedLatLng{Vertex: LatLng{Lat: 0.1, Lng: 0.2}, Next: nil}
				loop1 := &linkedGeoLoop{First: coord1, Last: coord1, Next: nil}
				polygon1 := &linkedGeoPolygon{First: loop1, Last: loop1, Next: nil}

				// Create second polygon
				coord2 := &linkedLatLng{Vertex: LatLng{Lat: 0.3, Lng: 0.4}, Next: nil}
				loop2 := &linkedGeoLoop{First: coord2, Last: coord2, Next: nil}
				polygon2 := &linkedGeoPolygon{First: loop2, Last: loop2, Next: nil}

				// Create third polygon
				coord3 := &linkedLatLng{Vertex: LatLng{Lat: 0.5, Lng: 0.6}, Next: nil}
				loop3 := &linkedGeoLoop{First: coord3, Last: coord3, Next: nil}
				polygon3 := &linkedGeoPolygon{First: loop3, Last: loop3, Next: nil}

				// Link polygons
				polygon1.Next = polygon2
				polygon2.Next = polygon3

				return polygon1
			},
			verify: func(t *testing.T, polygon *linkedGeoPolygon) {
				if polygon == nil {
					t.Error("Input polygon should not be nil after destruction")
					return
				}
				// The input polygon should have its loops destroyed but structure preserved
				if polygon.First != nil {
					t.Error("Expected polygon.First to be nil after destruction")
				}
				if polygon.Last != nil {
					t.Error("Expected polygon.Last to be nil after destruction")
				}
				if polygon.Next != nil {
					t.Error("Expected polygon.Next to be nil after destruction (other polygons should be freed)")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test Go implementation
			goPolygon := tt.createPolygon()
			destroyLinkedMultiPolygon(goPolygon)
			tt.verify(t, goPolygon)

			// Test C implementation - we can't really verify the state after C destruction
			// because it would involve invalid memory access. The C function test is mainly
			// to ensure it doesn't crash.
			cPolygon := tt.createPolygon()
			if cPolygon != nil {
				// The C function will free memory, so we can only test that it doesn't crash
				destroyLinkedMultiPolygonC(cPolygon)
				// We can't verify the C result since memory has been freed
			}
		})
	}
}

func TestDestroyLinkedMultiPolygonBehavior(t *testing.T) {
	// Test that verifies the specific behavior described in the C function:
	// "The caller is responsible for freeing memory allocated to input polygon struct"
	// This means the function should NOT free the input polygon itself, only its contents
	// and any additional polygons in the chain.

	t.Run("preserves_input_polygon_structure", func(t *testing.T) {
		// Create a polygon with some content
		coord := &linkedLatLng{Vertex: LatLng{Lat: 0.1, Lng: 0.2}, Next: nil}
		loop := &linkedGeoLoop{First: coord, Last: coord, Next: nil}
		inputPolygon := &linkedGeoPolygon{First: loop, Last: loop, Next: nil}

		// Store original address to verify it's not changed
		originalAddr := inputPolygon

		// Call destroy function
		destroyLinkedMultiPolygon(inputPolygon)

		// Verify input polygon structure still exists (same address)
		if inputPolygon != originalAddr {
			t.Error("Input polygon address should not change")
		}

		// Verify contents are cleared
		if inputPolygon.First != nil {
			t.Error("Expected polygon.First to be nil after destruction")
		}
		if inputPolygon.Last != nil {
			t.Error("Expected polygon.Last to be nil after destruction")
		}
		if inputPolygon.Next != nil {
			t.Error("Expected polygon.Next to be nil after destruction")
		}
	})

	t.Run("clears_coordinate_references", func(t *testing.T) {
		// Test that coordinate references are properly cleared to help GC
		coord1 := &linkedLatLng{Vertex: LatLng{Lat: 0.1, Lng: 0.2}, Next: nil}
		coord2 := &linkedLatLng{Vertex: LatLng{Lat: 0.3, Lng: 0.4}, Next: nil}
		coord3 := &linkedLatLng{Vertex: LatLng{Lat: 0.5, Lng: 0.6}, Next: nil}

		// Create a chain
		coord1.Next = coord2
		coord2.Next = coord3

		loop := &linkedGeoLoop{First: coord1, Last: coord3, Next: nil}
		polygon := &linkedGeoPolygon{First: loop, Last: loop, Next: nil}

		// Call destroy
		destroyLinkedMultiPolygon(polygon)

		// Verify coordinate chain is broken (helps GC)
		if coord1.Next != nil {
			t.Error("Expected coord1.Next to be nil after destruction")
		}
		if coord2.Next != nil {
			t.Error("Expected coord2.Next to be nil after destruction")
		}

		// Verify loop references are cleared
		if loop.First != nil {
			t.Error("Expected loop.First to be nil after destruction")
		}
		if loop.Last != nil {
			t.Error("Expected loop.Last to be nil after destruction")
		}
	})
}
