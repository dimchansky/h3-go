//go:build cgo && c2go

package c2go

import (
	"testing"
)

func Test_normalizeMultiPolygon_parity(t *testing.T) {
	// Helper to create a square loop with counter-clockwise winding (outer loop)
	// Using standard geographic coordinates: Lat = North/South, Lng = East/West
	createSquareLoopCCW := func(lat, lng, size float64) *LinkedGeoLoop {
		loop := &LinkedGeoLoop{}
		coords := []LatLng{
			{Lat: lat, Lng: lng},                     // Start at SW corner
			{Lat: lat, Lng: lng + size},              // Move east to SE corner  
			{Lat: lat + size, Lng: lng + size},       // Move north to NE corner
			{Lat: lat + size, Lng: lng},              // Move west to NW corner
			{Lat: lat, Lng: lng},                     // Back to start (CCW)
		}

		var prev *LinkedLatLng
		for i, coord := range coords {
			node := &LinkedLatLng{
				Vertex: coord,
				Next:   nil,
			}
			if i == 0 {
				loop.First = node
			} else {
				prev.Next = node
			}
			prev = node
			loop.Last = node
		}
		return loop
	}

	// Helper to create a square loop with clockwise winding (hole/inner loop)
	createSquareLoopCW := func(lat, lng, size float64) *LinkedGeoLoop {
		loop := &LinkedGeoLoop{}
		coords := []LatLng{
			{Lat: lat, Lng: lng},                     // Start at SW corner
			{Lat: lat + size, Lng: lng},              // Move north to NW corner
			{Lat: lat + size, Lng: lng + size},       // Move east to NE corner
			{Lat: lat, Lng: lng + size},              // Move south to SE corner
			{Lat: lat, Lng: lng},                     // Back to start (CW)
		}

		var prev *LinkedLatLng
		for i, coord := range coords {
			node := &LinkedLatLng{
				Vertex: coord,
				Next:   nil,
			}
			if i == 0 {
				loop.First = node
			} else {
				prev.Next = node
			}
			prev = node
			loop.Last = node
		}
		return loop
	}

	tests := []struct {
		name          string
		setupPolygon  func() *LinkedGeoPolygon
		expectedError H3Error
	}{
		// Note: Skipping nil polygon test as C function expects valid pointer
		// and would segfault with nil input. Go version handles it properly.
		{
			name: "empty polygon",
			setupPolygon: func() *LinkedGeoPolygon {
				return &LinkedGeoPolygon{}
			},
			expectedError: E_SUCCESS,
		},
		{
			name: "polygon with Next (multiple polygons) - should fail",
			setupPolygon: func() *LinkedGeoPolygon {
				root := &LinkedGeoPolygon{
					First: createSquareLoopCCW(0, 0, 0.1),
				}
				root.Last = root.First

				// Add another polygon (this should cause E_FAILED)
				root.Next = &LinkedGeoPolygon{
					First: createSquareLoopCCW(0.2, 0.2, 0.1),
				}
				root.Next.Last = root.Next.First

				return root
			},
			expectedError: E_FAILED,
		},
		{
			name: "single outer loop - should succeed",
			setupPolygon: func() *LinkedGeoPolygon {
				root := &LinkedGeoPolygon{
					First: createSquareLoopCCW(0, 0, 0.1),
				}
				root.Last = root.First
				return root
			},
			expectedError: E_SUCCESS,
		},
		{
			name: "outer loop with hole - should succeed",
			setupPolygon: func() *LinkedGeoPolygon {
				// Create outer loop (counter-clockwise)
				outerLoop := createSquareLoopCCW(0, 0, 0.1)

				// Create inner loop (clockwise - hole)
				innerLoop := createSquareLoopCW(0.02, 0.02, 0.06)

				// Link them together
				outerLoop.Next = innerLoop

				root := &LinkedGeoPolygon{
					First: outerLoop,
					Last:  innerLoop,
				}
				return root
			},
			expectedError: E_SUCCESS,
		},
		{
			name: "multiple outer loops with holes - should succeed",
			setupPolygon: func() *LinkedGeoPolygon {
				// Create first outer loop (counter-clockwise)
				outerLoop1 := createSquareLoopCCW(0, 0, 0.1)

				// Create hole for first polygon (clockwise)
				holeLoop1 := createSquareLoopCW(0.02, 0.02, 0.06)

				// Create second outer loop (counter-clockwise)
				outerLoop2 := createSquareLoopCCW(0.2, 0.2, 0.1)

				// Create hole for second polygon (clockwise)
				holeLoop2 := createSquareLoopCW(0.22, 0.22, 0.06)

				// Link all loops together
				outerLoop1.Next = holeLoop1
				holeLoop1.Next = outerLoop2
				outerLoop2.Next = holeLoop2

				root := &LinkedGeoPolygon{
					First: outerLoop1,
					Last:  holeLoop2,
				}
				return root
			},
			expectedError: E_SUCCESS,
		},
		{
			name: "orphaned hole (no containing polygon) - should fail",
			setupPolygon: func() *LinkedGeoPolygon {
				// Create outer loop in one area
				outerLoop := createSquareLoopCCW(0, 0, 0.1)

				// Create hole in completely different area (no container)
				orphanedHole := createSquareLoopCW(0.5, 0.5, 0.1)

				// Link them together
				outerLoop.Next = orphanedHole

				root := &LinkedGeoPolygon{
					First: outerLoop,
					Last:  orphanedHole,
				}
				return root
			},
			expectedError: E_FAILED, // Should fail because hole has no container
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test polygon
			goPolygon := tt.setupPolygon()

			// Create deep copy for C test (since C function modifies in-place)
			var cPolygon *LinkedGeoPolygon
			if goPolygon != nil {
				cPolygon = deepCopyLinkedGeoPolygon(goPolygon)
			}

			// Test Go implementation
			goResult := normalizeMultiPolygon(goPolygon)

			// Test C implementation - only if we have a valid polygon
			// C implementation can't handle nil input safely
			var cResult H3Error
			if cPolygon != nil {
				cResult = normalizeMultiPolygonC(cPolygon)

				// Compare results only when both implementations can run
				if goResult != cResult {
					t.Errorf("Result mismatch: Go=%v, C=%v", goResult, cResult)
				}

				// Check C result matches expected
				if cResult != tt.expectedError {
					t.Errorf("C result doesn't match expected: got %v, want %v", cResult, tt.expectedError)
				}
			} else {
				// For nil input, only test Go implementation
				// C implementation would segfault with nil input
			}

			// Check Go result matches expected (always check this)
			if goResult != tt.expectedError {
				t.Errorf("Go result doesn't match expected: got %v, want %v", goResult, tt.expectedError)
			}
		})
	}
}

// deepCopyLinkedGeoPolygon creates a deep copy of a LinkedGeoPolygon structure
func deepCopyLinkedGeoPolygon(orig *LinkedGeoPolygon) *LinkedGeoPolygon {
	if orig == nil {
		return nil
	}

	// Helper function to deep copy a LinkedGeoLoop
	copyLoop := func(origLoop *LinkedGeoLoop) *LinkedGeoLoop {
		if origLoop == nil {
			return nil
		}

		newLoop := &LinkedGeoLoop{}

		var prevNode *LinkedLatLng
		currentOrig := origLoop.First

		for currentOrig != nil {
			newNode := &LinkedLatLng{
				Vertex: LatLng{
					Lat: currentOrig.Vertex.Lat,
					Lng: currentOrig.Vertex.Lng,
				},
				Next: nil,
			}

			if newLoop.First == nil {
				newLoop.First = newNode
			} else {
				prevNode.Next = newNode
			}

			prevNode = newNode
			newLoop.Last = newNode
			currentOrig = currentOrig.Next
		}

		return newLoop
	}

	newPolygon := &LinkedGeoPolygon{}

	// Copy all loops
	var prevLoop *LinkedGeoLoop
	currentOrigLoop := orig.First

	for currentOrigLoop != nil {
		newLoop := copyLoop(currentOrigLoop)

		if newPolygon.First == nil {
			newPolygon.First = newLoop
		} else {
			prevLoop.Next = newLoop
		}

		prevLoop = newLoop
		newPolygon.Last = newLoop
		currentOrigLoop = currentOrigLoop.Next
	}

	// Copy Next pointer (though normalizeMultiPolygon should fail if this exists)
	newPolygon.Next = deepCopyLinkedGeoPolygon(orig.Next)

	return newPolygon
}
