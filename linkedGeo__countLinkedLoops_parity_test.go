//go:build cgo

package h3

import (
	"testing"
)

func Test_countLinkedLoops_parity(t *testing.T) {
	testCases := []struct {
		name    string
		polygon *LinkedGeoPolygon
		count   int32
	}{
		{
			name: "empty polygon",
			polygon: &LinkedGeoPolygon{
				First: nil,
				Last:  nil,
				Next:  nil,
			},
			count: 0,
		},
		{
			name: "single loop",
			polygon: &LinkedGeoPolygon{
				First: &LinkedGeoLoop{
					First: nil, // Don't need coords for loop counting
					Last:  nil,
					Next:  nil,
				},
				Last: nil, // Will be set to first
				Next: nil,
			},
			count: 1,
		},
		{
			name:    "two loops",
			polygon: createTestPolygon(2),
			count:   2,
		},
		{
			name:    "five loops",
			polygon: createTestPolygon(5),
			count:   5,
		},
		{
			name:    "large polygon",
			polygon: createTestPolygon(50),
			count:   50,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set up single loop case properly
			if tc.count == 1 && tc.polygon.First != nil {
				tc.polygon.Last = tc.polygon.First
			}

			// Test Go implementation
			goResult := countLinkedLoops(tc.polygon)

			// Test C implementation (skip empty polygon as it would cause issues with C malloc/free)
			if tc.count > 0 {
				cResult := countLinkedLoopsC(tc.polygon)

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

// Helper function to create a test polygon with specified number of loops
func createTestPolygon(numLoops int32) *LinkedGeoPolygon {
	if numLoops == 0 {
		return &LinkedGeoPolygon{First: nil, Last: nil, Next: nil}
	}

	first := &LinkedGeoLoop{
		First: nil, // Don't need coordinate data for loop counting
		Last:  nil,
		Next:  nil,
	}

	current := first
	for i := int32(1); i < numLoops; i++ {
		next := &LinkedGeoLoop{
			First: nil,
			Last:  nil,
			Next:  nil,
		}
		current.Next = next
		current = next
	}

	return &LinkedGeoPolygon{
		First: first,
		Last:  current,
		Next:  nil,
	}
}
