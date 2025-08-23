//go:build cgo && c2go

package c2go

import (
	"testing"
)

func Test_countLinkedPolygons_parity(t *testing.T) {
	tests := []struct {
		name          string
		setupPolygons func() *LinkedGeoPolygon
		expectedCount int
	}{
		{
			name: "nil polygon",
			setupPolygons: func() *LinkedGeoPolygon {
				return nil
			},
			expectedCount: 0,
		},
		{
			name: "single polygon",
			setupPolygons: func() *LinkedGeoPolygon {
				return &LinkedGeoPolygon{
					First: nil,
					Last:  nil,
					Next:  nil,
				}
			},
			expectedCount: 1,
		},
		{
			name: "two polygons",
			setupPolygons: func() *LinkedGeoPolygon {
				first := &LinkedGeoPolygon{
					First: nil,
					Last:  nil,
					Next:  nil,
				}
				second := &LinkedGeoPolygon{
					First: nil,
					Last:  nil,
					Next:  nil,
				}
				first.Next = second
				return first
			},
			expectedCount: 2,
		},
		{
			name: "three polygons",
			setupPolygons: func() *LinkedGeoPolygon {
				first := &LinkedGeoPolygon{
					First: nil,
					Last:  nil,
					Next:  nil,
				}
				second := &LinkedGeoPolygon{
					First: nil,
					Last:  nil,
					Next:  nil,
				}
				third := &LinkedGeoPolygon{
					First: nil,
					Last:  nil,
					Next:  nil,
				}
				first.Next = second
				second.Next = third
				return first
			},
			expectedCount: 3,
		},
		{
			name: "five polygons",
			setupPolygons: func() *LinkedGeoPolygon {
				polygons := make([]*LinkedGeoPolygon, 5)
				for i := 0; i < 5; i++ {
					polygons[i] = &LinkedGeoPolygon{
						First: nil,
						Last:  nil,
						Next:  nil,
					}
					if i > 0 {
						polygons[i-1].Next = polygons[i]
					}
				}
				return polygons[0]
			},
			expectedCount: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test data
			polygon := tt.setupPolygons()

			// Test Go implementation
			goResult := countLinkedPolygons(polygon)

			// Test C implementation
			cResult := countLinkedPolygonsC(polygon)

			// Compare results
			if goResult != cResult {
				t.Errorf("Result mismatch: Go=%d, C=%d", goResult, cResult)
			}

			// Also verify against expected count
			if goResult != tt.expectedCount {
				t.Errorf("Go result %d does not match expected %d", goResult, tt.expectedCount)
			}
			if cResult != tt.expectedCount {
				t.Errorf("C result %d does not match expected %d", cResult, tt.expectedCount)
			}
		})
	}

	// Additional test: verify both handle larger chains correctly
	t.Run("large chain", func(t *testing.T) {
		// Create a chain of 100 polygons
		const chainSize = 100
		var first, prev *LinkedGeoPolygon

		for i := 0; i < chainSize; i++ {
			polygon := &LinkedGeoPolygon{
				First: nil,
				Last:  nil,
				Next:  nil,
			}

			if first == nil {
				first = polygon
			} else {
				prev.Next = polygon
			}
			prev = polygon
		}

		// Test Go implementation
		goResult := countLinkedPolygons(first)

		// Test C implementation
		cResult := countLinkedPolygonsC(first)

		// Compare results
		if goResult != cResult {
			t.Errorf("Result mismatch for large chain: Go=%d, C=%d", goResult, cResult)
		}

		if goResult != chainSize {
			t.Errorf("Go result %d does not match expected %d", goResult, chainSize)
		}
		if cResult != chainSize {
			t.Errorf("C result %d does not match expected %d", cResult, chainSize)
		}
	})
}
