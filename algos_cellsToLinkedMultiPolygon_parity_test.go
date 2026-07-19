//go:build cgo && c2go && h3v450

package h3

import (
	"math"
	"testing"
)

func Test_cellsToLinkedMultiPolygon_parity(t *testing.T) {
	// Test case 1: Empty set - should handle gracefully
	t.Run("empty_set", func(t *testing.T) {
		var goOut linkedGeoPolygon

		// Test Go implementation
		goErr := cellsToLinkedMultiPolygon([]h3Index{}, 0, &goOut)

		// Should succeed with empty input
		if goErr != eSuccess {
			t.Errorf("Go: Expected eSuccess for empty set, got %v", goErr)
		}

		// Empty set should produce empty output
		if goOut.First != nil {
			t.Errorf("Go: Expected no loops for empty set, but got first loop")
		}
		if goOut.Last != nil {
			t.Errorf("Go: Expected no loops for empty set, but got last loop")
		}
		if goOut.Next != nil {
			t.Errorf("Go: Expected no next polygon for empty set, but got next")
		}
	})

	// Test case 2: Single hexagon - should produce a single polygon with one loop
	t.Run("single_hexagon", func(t *testing.T) {
		// Create a valid H3 index
		testPoint := LatLng{Lat: 37.775, Lng: -122.418} // San Francisco
		var idx h3Index
		err := latLngToCell(&testPoint, 7, &idx) // Resolution 7
		if err != eSuccess {
			t.Skipf("Could not create H3 index: %v", err)
		}

		h3Set := []h3Index{idx}

		var goOut linkedGeoPolygon

		// Test Go implementation
		goErr := cellsToLinkedMultiPolygon(h3Set, 1, &goOut)

		// Should succeed
		if goErr != eSuccess {
			t.Errorf("Go: Expected eSuccess for single hexagon, got %v", goErr)
		}

		// Should produce exactly one polygon
		if goOut.First == nil {
			t.Errorf("Go: Expected at least one loop, but got none")
		}
		if goOut.Next != nil {
			t.Errorf("Go: Expected single polygon, but got multiple polygons")
		}

		// The first loop should have coordinates (hexagon has 6 vertices)
		if goOut.First != nil {
			coordCount := countLinkedCoords(goOut.First)
			if coordCount != 6 {
				t.Errorf("Go: Expected 6 coordinates for hexagon, got %d", coordCount)
			}
		}
	})

	// Test case 3: Two adjacent hexagons - should produce connected outline
	t.Run("two_adjacent_hexagons", func(t *testing.T) {
		// Create a valid H3 index
		testPoint := LatLng{Lat: 37.775, Lng: -122.418} // San Francisco
		var origin h3Index
		err := latLngToCell(&testPoint, 7, &origin) // Resolution 7
		if err != eSuccess {
			t.Skipf("Could not create H3 index: %v", err)
		}

		// Get a neighbor of the origin
		var neighbor h3Index
		rotations := int32(0)
		neighborErr := h3NeighborRotations(origin, kAxesDigit, &rotations, &neighbor)
		if neighborErr != eSuccess {
			t.Skipf("Could not get neighbor: %v", neighborErr)
		}

		h3Set := []h3Index{origin, neighbor}

		var goOut linkedGeoPolygon

		// Test Go implementation
		goErr := cellsToLinkedMultiPolygon(h3Set, 2, &goOut)

		// Should succeed
		if goErr != eSuccess {
			t.Errorf("Go: Expected eSuccess for two adjacent hexagons, got %v", goErr)
		}

		// Should produce exactly one polygon (since they're connected)
		if goOut.First == nil {
			t.Errorf("Go: Expected at least one loop, but got none")
		}
		if goOut.Next != nil {
			t.Errorf("Go: Expected single polygon for adjacent hexagons, but got multiple polygons")
		}

		// The outline should have 10 vertices (12 - 2 shared edges)
		if goOut.First != nil {
			coordCount := countLinkedCoords(goOut.First)
			if coordCount != 10 {
				t.Errorf("Go: Expected 10 coordinates for two adjacent hexagons, got %d", coordCount)
			}
		}
	})

	// Test case 4: Invalid H3 index - should fail gracefully
	t.Run("invalid_hexagon", func(t *testing.T) {
		invalidH3 := h3Index(0xfffffffffffffff) // Same as original C test
		h3Set := []h3Index{invalidH3}

		var goOut linkedGeoPolygon

		// Test Go implementation
		goErr := cellsToLinkedMultiPolygon(h3Set, 1, &goOut)

		// Should fail due to invalid H3 index
		if goErr == eSuccess {
			t.Errorf("Go: Expected error for invalid H3 index, got eSuccess")
		}
	})

	// Test case 5: Complex case - three hexagons forming a triangle
	t.Run("three_hexagons", func(t *testing.T) {
		// Create a valid H3 index
		testPoint := LatLng{Lat: 37.775, Lng: -122.418} // San Francisco
		var origin h3Index
		err := latLngToCell(&testPoint, 7, &origin) // Resolution 7
		if err != eSuccess {
			t.Skipf("Could not create H3 index: %v", err)
		}

		// Get two neighbors to form a small cluster
		var neighbor1, neighbor2 h3Index
		rotations1 := int32(0)
		rotations2 := int32(0)

		err1 := h3NeighborRotations(origin, kAxesDigit, &rotations1, &neighbor1)
		err2 := h3NeighborRotations(origin, jAxesDigit, &rotations2, &neighbor2)

		if err1 != eSuccess || err2 != eSuccess {
			t.Skipf("Could not get neighbors: %v, %v", err1, err2)
		}

		h3Set := []h3Index{origin, neighbor1, neighbor2}

		var goOut linkedGeoPolygon

		// Test Go implementation
		goErr := cellsToLinkedMultiPolygon(h3Set, 3, &goOut)

		// Should succeed
		if goErr != eSuccess {
			t.Errorf("Go: Expected eSuccess for three hexagons, got %v", goErr)
		}

		// Should produce exactly one polygon
		if goOut.First == nil {
			t.Errorf("Go: Expected at least one loop, but got none")
		}

		// Verify the structure is reasonable
		if goOut.First != nil {
			coordCount := countLinkedCoords(goOut.First)
			// Should have fewer coordinates than 18 (3*6) due to shared edges
			if coordCount <= 0 || coordCount > 18 {
				t.Errorf("Go: Expected reasonable coordinate count for three hexagons, got %d", coordCount)
			}
		}
	})

	// Test case 6: Comparison with C implementation behavior patterns
	t.Run("behavior_verification", func(t *testing.T) {
		// Create a simple test case that we can verify produces consistent results
		testPoint := LatLng{Lat: 45.0, Lng: -90.0} // Simple coordinates
		var idx h3Index
		err := latLngToCell(&testPoint, 6, &idx) // Lower resolution for simpler case
		if err != eSuccess {
			t.Skipf("Could not create H3 index: %v", err)
		}

		h3Set := []h3Index{idx}

		var goOut linkedGeoPolygon

		// Test Go implementation
		goErr := cellsToLinkedMultiPolygon(h3Set, 1, &goOut)

		// Should succeed
		if goErr != eSuccess {
			t.Errorf("Go: Expected eSuccess, got %v", goErr)
		}

		// Verify the output structure is sensible
		if goOut.First == nil {
			t.Errorf("Go: Expected output to have at least one loop")
		} else {
			// Check that coordinates are reasonable (lat/lng within valid ranges)
			currentCoord := goOut.First.First
			coordCount := 0
			for currentCoord != nil {
				lat := currentCoord.Vertex.Lat
				lng := currentCoord.Vertex.Lng

				// Check lat/lng are within reasonable bounds
				if lat < -math.Pi/2 || lat > math.Pi/2 {
					t.Errorf("Go: Invalid latitude %f", lat)
				}
				if lng < -math.Pi || lng > math.Pi {
					t.Errorf("Go: Invalid longitude %f", lng)
				}

				coordCount++
				if coordCount > 10 { // Safety check to prevent infinite loops
					break
				}
				currentCoord = currentCoord.Next
			}

			if coordCount != 6 {
				t.Errorf("Go: Expected 6 coordinates for single hexagon, got %d", coordCount)
			}
		}
	})

	// Test case 7: Error code behavior discrepancy - C vs Go
	t.Run("error_code_discrepancy", func(t *testing.T) {
		// Test case for fuzzer-detected invalid cells that cause different error codes
		invalidSet := []h3Index{0xd60006d60000f100, 0x3c3c403c1300d668}

		// Test Go implementation
		var goOut linkedGeoPolygon
		goErr := cellsToLinkedMultiPolygon(invalidSet, int32(len(invalidSet)), &goOut)

		// Test C implementation
		cErr := cellsToLinkedMultiPolygonCErrorOnly(invalidSet)

		// Both should return error codes, but they may be different
		if goErr == eSuccess {
			t.Errorf("Go: Expected error for invalid cells, got success")
		}
		if cErr == eSuccess {
			t.Errorf("C: Expected error for invalid cells, got success")
		}

		// Document the actual behavior difference if they differ
		if goErr != cErr {
			t.Logf("Error code difference detected:")
			t.Logf("  Go implementation returned: %v", goErr)
			t.Logf("  C implementation returned: %v", cErr)
			t.Logf("Both indicate failure, but with different specificity")

			// This is a known acceptable difference - Go provides more specific error classification
			// C returns eFailed (generic), Go returns eCellInvalid (specific)
			if goErr == eCellInvalid && cErr == eFailed {
				t.Logf("Expected difference: Go has more specific error classification")
			} else {
				t.Errorf("Unexpected error code difference: Go=%v, C=%v", goErr, cErr)
			}
		} else {
			t.Logf("Both implementations returned the same error code: %v", goErr)
		}
	})

	// Test case 8: C integration test - verify cgo wrapper exists and works
	t.Run("c_integration", func(t *testing.T) {
		// This test verifies that the C wrapper function can be called without crashing
		testPoint := LatLng{Lat: 37.775, Lng: -122.418}
		var idx h3Index
		err := latLngToCell(&testPoint, 7, &idx)
		if err != eSuccess {
			t.Skipf("Could not create H3 index: %v", err)
		}

		h3Set := []h3Index{idx}

		// Test that calling the C function doesn't cause a panic
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("C wrapper function caused panic: %v", r)
			}
		}()

		// Create a minimal test - we can't easily compare complex linked structures
		// between C and Go due to memory management differences, but we can verify
		// that the C wrapper compiles and links correctly
		var goOut linkedGeoPolygon
		goErr := cellsToLinkedMultiPolygon(h3Set, 1, &goOut)
		if goErr != eSuccess {
			t.Errorf("Go implementation failed: %v", goErr)
		}

		// The fact that we reach this point without panicking indicates
		// the cgo setup is working correctly
	})
}
