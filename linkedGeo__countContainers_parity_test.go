//go:build cgo && c2go

package h3

import (
	"testing"
)

func Test_countContainers_parity(t *testing.T) {
	// Helper to create a simple square loop
	createSquareLoop := func(x, y float64) *linkedGeoLoop {
		loop := &linkedGeoLoop{}
		coords := []LatLng{
			{Lat: Rad(x), Lng: Rad(y)},
			{Lat: Rad(x), Lng: Rad(y + 0.1)},
			{Lat: Rad(x + 0.1), Lng: Rad(y + 0.1)},
			{Lat: Rad(x + 0.1), Lng: Rad(y)},
			{Lat: Rad(x), Lng: Rad(y)}, // Close the loop
		}

		var prev *linkedLatLng
		for i, coord := range coords {
			node := &linkedLatLng{
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

	// Helper to create bb for a loop
	createBboxForLoop := func(loop *linkedGeoLoop) *bbox {
		bb := &bbox{North: Deg(-90), South: Deg(90), East: Deg(-180), West: Deg(180)}

		// Simple bb calculation (not exact but sufficient for testing)
		current := loop.First
		for current != nil {
			if current.Vertex.Lat.Rad() > bb.North.Rad() {
				bb.North = current.Vertex.Lat
			}
			if current.Vertex.Lat.Rad() < bb.South.Rad() {
				bb.South = current.Vertex.Lat
			}
			if current.Vertex.Lng.Rad() > bb.East.Rad() {
				bb.East = current.Vertex.Lng
			}
			if current.Vertex.Lng.Rad() < bb.West.Rad() {
				bb.West = current.Vertex.Lng
			}
			current = current.Next
		}

		return bb
	}

	tests := []struct {
		name          string
		setupTest     func() (*linkedGeoLoop, []*linkedGeoPolygon, []*bbox)
		expectedCount int32
	}{
		{
			name: "empty polygons array",
			setupTest: func() (*linkedGeoLoop, []*linkedGeoPolygon, []*bbox) {
				loop := createSquareLoop(0, 0)
				return loop, []*linkedGeoPolygon{}, []*bbox{}
			},
			expectedCount: 0,
		},
		{
			name: "loop not inside any polygon",
			setupTest: func() (*linkedGeoLoop, []*linkedGeoPolygon, []*bbox) {
				// Test loop
				testLoop := createSquareLoop(0, 0)

				// Polygon that doesn't contain the test loop
				polygon := &linkedGeoPolygon{
					First: createSquareLoop(10, 10), // Far away
				}
				polygon.Last = polygon.First

				bb := createBboxForLoop(polygon.First)

				return testLoop, []*linkedGeoPolygon{polygon}, []*bbox{bb}
			},
			expectedCount: 0,
		},
		{
			name: "loop is the same as polygon's first loop",
			setupTest: func() (*linkedGeoLoop, []*linkedGeoPolygon, []*bbox) {
				// Test loop
				testLoop := createSquareLoop(0, 0)

				// Polygon with the same loop as its first
				polygon := &linkedGeoPolygon{
					First: testLoop, // Same loop - should not count as container
				}
				polygon.Last = polygon.First

				bb := createBboxForLoop(polygon.First)

				return testLoop, []*linkedGeoPolygon{polygon}, []*bbox{bb}
			},
			expectedCount: 0,
		},
		{
			name: "loop inside one polygon",
			setupTest: func() (*linkedGeoLoop, []*linkedGeoPolygon, []*bbox) {
				// Small test loop inside
				testLoop := createSquareLoop(0.025, 0.025)

				// Larger polygon containing the test loop
				polygon := &linkedGeoPolygon{
					First: createSquareLoop(0, 0), // Contains test loop
				}
				polygon.Last = polygon.First

				bb := createBboxForLoop(polygon.First)

				return testLoop, []*linkedGeoPolygon{polygon}, []*bbox{bb}
			},
			expectedCount: 1,
		},
		{
			name: "loop inside multiple nested polygons",
			setupTest: func() (*linkedGeoLoop, []*linkedGeoPolygon, []*bbox) {
				// Small test loop at center
				testLoop := createSquareLoop(0.04, 0.04)

				// Create nested polygons (each contains the test loop)
				polygon1 := &linkedGeoPolygon{
					First: createSquareLoop(0, 0), // Largest
				}
				polygon1.Last = polygon1.First

				polygon2 := &linkedGeoPolygon{
					First: createSquareLoop(0.02, 0.02), // Medium
				}
				polygon2.Last = polygon2.First

				polygon3 := &linkedGeoPolygon{
					First: createSquareLoop(0.03, 0.03), // Smallest but still contains test
				}
				polygon3.Last = polygon3.First

				bbox1 := createBboxForLoop(polygon1.First)
				bbox2 := createBboxForLoop(polygon2.First)
				bbox3 := createBboxForLoop(polygon3.First)

				return testLoop,
					[]*linkedGeoPolygon{polygon1, polygon2, polygon3},
					[]*bbox{bbox1, bbox2, bbox3}
			},
			expectedCount: 3,
		},
		{
			name: "nil loop",
			setupTest: func() (*linkedGeoLoop, []*linkedGeoPolygon, []*bbox) {
				polygon := &linkedGeoPolygon{
					First: createSquareLoop(0, 0),
				}
				polygon.Last = polygon.First

				bb := createBboxForLoop(polygon.First)

				return nil, []*linkedGeoPolygon{polygon}, []*bbox{bb}
			},
			expectedCount: 0,
		},
		{
			name: "loop with nil first coordinate",
			setupTest: func() (*linkedGeoLoop, []*linkedGeoPolygon, []*bbox) {
				testLoop := &linkedGeoLoop{
					First: nil,
					Last:  nil,
				}

				polygon := &linkedGeoPolygon{
					First: createSquareLoop(0, 0),
				}
				polygon.Last = polygon.First

				bb := createBboxForLoop(polygon.First)

				return testLoop, []*linkedGeoPolygon{polygon}, []*bbox{bb}
			},
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test data
			loop, polygons, bboxes := tt.setupTest()

			// Test Go implementation
			goResult := countContainers(loop, polygons, bboxes)

			// Test C implementation
			cResult := countContainersC(loop, polygons, bboxes)

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
}
