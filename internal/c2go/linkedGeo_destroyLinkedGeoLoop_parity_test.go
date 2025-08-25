//go:build cgo

package c2go

import (
	"testing"
)

func Test_destroyLinkedGeoLoop_parity(t *testing.T) {
	tests := []struct {
		name        string
		setupLoop   func() *LinkedGeoLoop
		description string
	}{
		{
			name: "nil loop",
			setupLoop: func() *LinkedGeoLoop {
				return nil
			},
			description: "destroyLinkedGeoLoop with nil loop should handle gracefully",
		},
		{
			name: "empty loop",
			setupLoop: func() *LinkedGeoLoop {
				return &LinkedGeoLoop{
					First: nil,
					Last:  nil,
					Next:  nil,
				}
			},
			description: "destroyLinkedGeoLoop with empty loop should handle gracefully",
		},
		{
			name: "single coordinate",
			setupLoop: func() *LinkedGeoLoop {
				coord := &LinkedLatLng{
					Vertex: LatLng{Lat: 37.775, Lng: -122.418},
					Next:   nil,
				}
				return &LinkedGeoLoop{
					First: coord,
					Last:  coord,
					Next:  nil,
				}
			},
			description: "destroyLinkedGeoLoop with single coordinate",
		},
		{
			name: "multiple coordinates",
			setupLoop: func() *LinkedGeoLoop {
				coord1 := &LinkedLatLng{
					Vertex: LatLng{Lat: 37.775, Lng: -122.418},
					Next:   nil,
				}
				coord2 := &LinkedLatLng{
					Vertex: LatLng{Lat: 37.776, Lng: -122.419},
					Next:   nil,
				}
				coord3 := &LinkedLatLng{
					Vertex: LatLng{Lat: 37.777, Lng: -122.420},
					Next:   nil,
				}

				// Link them together
				coord1.Next = coord2
				coord2.Next = coord3

				return &LinkedGeoLoop{
					First: coord1,
					Last:  coord3,
					Next:  nil,
				}
			},
			description: "destroyLinkedGeoLoop with multiple coordinates in a chain",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test loop
			goLoop := tt.setupLoop()

			// Create a copy for C function testing (deep copy)
			var cTestLoop *LinkedGeoLoop
			if goLoop != nil {
				cTestLoop = &LinkedGeoLoop{
					First: nil,
					Last:  nil,
					Next:  goLoop.Next,
				}

				// Deep copy coordinates
				if goLoop.First != nil {
					var firstCopy *LinkedLatLng
					var prevCopy *LinkedLatLng

					current := goLoop.First
					for current != nil {
						coordCopy := &LinkedLatLng{
							Vertex: current.Vertex,
							Next:   nil,
						}

						if firstCopy == nil {
							firstCopy = coordCopy
							cTestLoop.First = firstCopy
						} else {
							prevCopy.Next = coordCopy
						}

						prevCopy = coordCopy
						current = current.Next
					}
					cTestLoop.Last = prevCopy
				}
			}

			// Test Go implementation - check that references are cleared
			if goLoop != nil {
				// Count coordinates before
				beforeCount := 0
				current := goLoop.First
				for current != nil {
					beforeCount++
					current = current.Next
				}

				// Call Go function
				destroyLinkedGeoLoop(goLoop)

				// Verify loop references are cleared
				if goLoop.First != nil {
					t.Errorf("Go destroyLinkedGeoLoop: loop.First should be nil after destroy, got non-nil")
				}
				if goLoop.Last != nil {
					t.Errorf("Go destroyLinkedGeoLoop: loop.Last should be nil after destroy, got non-nil")
				}

				t.Logf("Go destroyLinkedGeoLoop: successfully processed %d coordinates", beforeCount)
			} else {
				// Test with nil - should not panic
				destroyLinkedGeoLoop(goLoop)
				t.Log("Go destroyLinkedGeoLoop: handled nil loop gracefully")
			}

			// Test C implementation - mainly check that it doesn't crash
			if cTestLoop != nil {
				// Count coordinates before for logging
				beforeCount := 0
				current := cTestLoop.First
				for current != nil {
					beforeCount++
					current = current.Next
				}

				// Call C function - this tests that our C wrapper works correctly
				destroyLinkedGeoLoopC(cTestLoop)

				t.Logf("C destroyLinkedGeoLoopC: successfully processed %d coordinates", beforeCount)
			} else {
				// Test with nil - should not crash
				destroyLinkedGeoLoopC(cTestLoop)
				t.Log("C destroyLinkedGeoLoopC: handled nil loop gracefully")
			}

			t.Logf("Test passed: %s", tt.description)
		})
	}
}

func Test_destroyLinkedGeoLoop_behavior(t *testing.T) {
	// Test that the function properly clears references to help GC
	coord1 := &LinkedLatLng{
		Vertex: LatLng{Lat: 37.775, Lng: -122.418},
		Next:   nil,
	}
	coord2 := &LinkedLatLng{
		Vertex: LatLng{Lat: 37.776, Lng: -122.419},
		Next:   nil,
	}

	coord1.Next = coord2

	loop := &LinkedGeoLoop{
		First: coord1,
		Last:  coord2,
		Next:  nil,
	}

	// Verify initial state
	if loop.First != coord1 {
		t.Error("Initial state: loop.First should point to coord1")
	}
	if loop.Last != coord2 {
		t.Error("Initial state: loop.Last should point to coord2")
	}
	if coord1.Next != coord2 {
		t.Error("Initial state: coord1.Next should point to coord2")
	}

	// Call the function
	destroyLinkedGeoLoop(loop)

	// Verify the function clears the loop's references
	if loop.First != nil {
		t.Error("After destroy: loop.First should be nil")
	}
	if loop.Last != nil {
		t.Error("After destroy: loop.Last should be nil")
	}

	// The function should also clear the coordinate chain references
	if coord1.Next != nil {
		t.Error("After destroy: coord1.Next should be nil")
	}
}
