//go:build cgo && c2go && !h3v450

package h3

import (
	"fmt"
	"testing"
)

func Test_h3SetToVertexGraph_parity(t *testing.T) {
	// Test case 1: Empty set - compare Go vs C behavior
	t.Run("empty_set", func(t *testing.T) {
		// Test Go implementation
		var goGraph vertexGraph
		goErr := h3SetToVertexGraph([]h3Index{}, 0, &goGraph)

		// Test C implementation
		cResult := h3SetToVertexGraphCForParity([]h3Index{})

		// Compare error codes
		if goErr != cResult.Err {
			t.Errorf("Error code mismatch: Go=%v, C=%v", goErr, cResult.Err)
		}

		// Both should succeed
		if goErr != eSuccess {
			t.Errorf("Go: Expected eSuccess for empty set, got %v", goErr)
		}

		// Compare graph properties
		if cResult.Size != goGraph.Size {
			t.Errorf("Size mismatch: Go=%d, C=%d", goGraph.Size, cResult.Size)
		}
		if cResult.NumBuckets != goGraph.NumBuckets {
			t.Errorf("NumBuckets mismatch: Go=%d, C=%d", goGraph.NumBuckets, cResult.NumBuckets)
		}

		// Clean up Go graph
		destroyVertexGraph(&goGraph)
	})

	// Test case 2: Single hexagon - compare Go vs C behavior
	t.Run("single_hexagon", func(t *testing.T) {
		// Create a valid H3 index by converting lat/lng to cell
		testPoint := LatLng{Lat: 37.775, Lng: -122.418} // San Francisco
		var idx h3Index
		err := latLngToCell(&testPoint, 7, &idx) // Resolution 7
		if err != eSuccess {
			t.Skipf("Could not create H3 index: %v", err)
		}

		h3Set := []h3Index{idx}

		// Test Go implementation
		var goGraph vertexGraph
		goErr := h3SetToVertexGraph(h3Set, 1, &goGraph)

		// Test C implementation
		cResult := h3SetToVertexGraphCForParity(h3Set)

		// Compare error codes
		if goErr != cResult.Err {
			t.Errorf("Error code mismatch: Go=%v, C=%v", goErr, cResult.Err)
		}

		// Both should succeed
		if goErr != eSuccess {
			t.Errorf("Go: Expected eSuccess for single hexagon, got %v", goErr)
		}

		// Compare graph properties - both should have same size and bucket count
		if cResult.Size != goGraph.Size {
			t.Errorf("Size mismatch: Go=%d, C=%d", goGraph.Size, cResult.Size)
		}
		if cResult.NumBuckets != goGraph.NumBuckets {
			t.Errorf("NumBuckets mismatch: Go=%d, C=%d", goGraph.NumBuckets, cResult.NumBuckets)
		}
		if cResult.Res != goGraph.Res {
			t.Errorf("Resolution mismatch: Go=%d, C=%d", goGraph.Res, cResult.Res)
		}

		// For a single hexagon, both should have 6 edges
		expectedEdges := int32(6)
		if goGraph.Size != expectedEdges {
			t.Errorf("Go: Expected %d edges for single hexagon, got %d", expectedEdges, goGraph.Size)
		}
		if cResult.Size != expectedEdges {
			t.Errorf("C: Expected %d edges for single hexagon, got %d", expectedEdges, cResult.Size)
		}

		// Clean up Go graph
		destroyVertexGraph(&goGraph)
	})

	// Test case 3: Two adjacent hexagons - compare Go vs C behavior
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

		// Test Go implementation
		var goGraph vertexGraph
		goErr := h3SetToVertexGraph(h3Set, 2, &goGraph)

		// Test C implementation
		cResult := h3SetToVertexGraphCForParity(h3Set)

		// Compare error codes
		if goErr != cResult.Err {
			t.Errorf("Error code mismatch: Go=%v, C=%v", goErr, cResult.Err)
		}

		// Both should succeed
		if goErr != eSuccess {
			t.Errorf("Go: Expected eSuccess for two adjacent hexagons, got %v", goErr)
		}

		// Compare graph properties
		if cResult.Size != goGraph.Size {
			t.Errorf("Size mismatch: Go=%d, C=%d", goGraph.Size, cResult.Size)
		}
		if cResult.NumBuckets != goGraph.NumBuckets {
			t.Errorf("NumBuckets mismatch: Go=%d, C=%d", goGraph.NumBuckets, cResult.NumBuckets)
		}
		if cResult.Res != goGraph.Res {
			t.Errorf("Resolution mismatch: Go=%d, C=%d", goGraph.Res, cResult.Res)
		}

		// For two adjacent hexagons, both should have 10 edges (6+6-2 shared)
		expectedEdges := int32(10)
		if goGraph.Size != expectedEdges {
			t.Errorf("Go: Expected %d edges for two adjacent hexagons, got %d", expectedEdges, goGraph.Size)
		}
		if cResult.Size != expectedEdges {
			t.Errorf("C: Expected %d edges for two adjacent hexagons, got %d", expectedEdges, cResult.Size)
		}

		// Clean up Go graph
		destroyVertexGraph(&goGraph)
	})

	// Test case 4: Invalid hexagon - compare Go vs C error handling
	t.Run("invalid_hexagon", func(t *testing.T) {
		// Use an invalid H3 index that should cause both implementations to fail
		invalidH3 := h3Index(0xfffffffffffffff)
		h3Set := []h3Index{invalidH3}

		// Test Go implementation
		var goGraph vertexGraph
		goErr := h3SetToVertexGraph(h3Set, 1, &goGraph)

		// Test C implementation
		cResult := h3SetToVertexGraphCForParity(h3Set)

		// Compare error codes - both should fail with the same error
		if goErr != cResult.Err {
			t.Errorf("Error code mismatch: Go=%v, C=%v", goErr, cResult.Err)
		}

		// Both should fail due to invalid H3 index
		if goErr == eSuccess {
			t.Errorf("Go: Expected error for invalid H3 index, got eSuccess")
		}
		if cResult.Err == eSuccess {
			t.Errorf("C: Expected error for invalid H3 index, got eSuccess")
		}

		// Clean up Go graph if it was partially initialized
		if goGraph.Buckets != nil {
			destroyVertexGraph(&goGraph)
		}
	})

	// Test case 5: Multiple test cases with different resolutions
	t.Run("different_resolutions", func(t *testing.T) {
		resolutions := []int32{5, 8, 10}

		for _, res := range resolutions {
			t.Run(fmt.Sprintf("resolution_%d", res), func(t *testing.T) {
				// Create a valid H3 index at this resolution
				testPoint := LatLng{Lat: 40.689, Lng: -74.045} // New York
				var idx h3Index
				err := latLngToCell(&testPoint, res, &idx)
				if err != eSuccess {
					t.Skipf("Could not create H3 index at resolution %d: %v", res, err)
				}

				h3Set := []h3Index{idx}

				// Test Go implementation
				var goGraph vertexGraph
				goErr := h3SetToVertexGraph(h3Set, 1, &goGraph)

				// Test C implementation
				cResult := h3SetToVertexGraphCForParity(h3Set)

				// Compare results
				if goErr != cResult.Err {
					t.Errorf("Resolution %d: Error code mismatch: Go=%v, C=%v", res, goErr, cResult.Err)
				}
				if goErr != eSuccess {
					t.Errorf("Resolution %d: Go implementation failed: %v", res, goErr)
				}

				// Verify both have same properties
				if cResult.Size != goGraph.Size {
					t.Errorf("Resolution %d: Size mismatch: Go=%d, C=%d", res, goGraph.Size, cResult.Size)
				}
				if cResult.Res != goGraph.Res {
					t.Errorf("Resolution %d: Resolution mismatch: Go=%d, C=%d", res, goGraph.Res, cResult.Res)
				}

				// Clean up
				destroyVertexGraph(&goGraph)
			})
		}
	})
}
