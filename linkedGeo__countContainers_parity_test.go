//go:build cgo && c2go

package h3

import (
	"testing"
)

func Test_countContainers_parity(t *testing.T) {
	// Helper to create a simple square loop
	createSquareLoop := func(x, y float64) *LinkedGeoLoop {
		loop := &LinkedGeoLoop{}
		coords := []LatLng{
			{Lat: Rad(x), Lng: Rad(y)},
			{Lat: Rad(x), Lng: Rad(y + 0.1)},
			{Lat: Rad(x + 0.1), Lng: Rad(y + 0.1)},
			{Lat: Rad(x + 0.1), Lng: Rad(y)},
			{Lat: Rad(x), Lng: Rad(y)}, // Close the loop
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

	// Helper to create bbox for a loop
	createBboxForLoop := func(loop *LinkedGeoLoop) *BBox {
		bbox := &BBox{North: Deg(-90), South: Deg(90), East: Deg(-180), West: Deg(180)}

		// Simple bbox calculation (not exact but sufficient for testing)
		current := loop.First
		for current != nil {
			if current.Vertex.Lat.Rad() > bbox.North.Rad() {
				bbox.North = current.Vertex.Lat
			}
			if current.Vertex.Lat.Rad() < bbox.South.Rad() {
				bbox.South = current.Vertex.Lat
			}
			if current.Vertex.Lng.Rad() > bbox.East.Rad() {
				bbox.East = current.Vertex.Lng
			}
			if current.Vertex.Lng.Rad() < bbox.West.Rad() {
				bbox.West = current.Vertex.Lng
			}
			current = current.Next
		}

		return bbox
	}

	tests := []struct {
		name          string
		setupTest     func() (*LinkedGeoLoop, []*LinkedGeoPolygon, []*BBox)
		expectedCount int32
	}{
		{
			name: "empty polygons array",
			setupTest: func() (*LinkedGeoLoop, []*LinkedGeoPolygon, []*BBox) {
				loop := createSquareLoop(0, 0)
				return loop, []*LinkedGeoPolygon{}, []*BBox{}
			},
			expectedCount: 0,
		},
		{
			name: "loop not inside any polygon",
			setupTest: func() (*LinkedGeoLoop, []*LinkedGeoPolygon, []*BBox) {
				// Test loop
				testLoop := createSquareLoop(0, 0)

				// Polygon that doesn't contain the test loop
				polygon := &LinkedGeoPolygon{
					First: createSquareLoop(10, 10), // Far away
				}
				polygon.Last = polygon.First

				bbox := createBboxForLoop(polygon.First)

				return testLoop, []*LinkedGeoPolygon{polygon}, []*BBox{bbox}
			},
			expectedCount: 0,
		},
		{
			name: "loop is the same as polygon's first loop",
			setupTest: func() (*LinkedGeoLoop, []*LinkedGeoPolygon, []*BBox) {
				// Test loop
				testLoop := createSquareLoop(0, 0)

				// Polygon with the same loop as its first
				polygon := &LinkedGeoPolygon{
					First: testLoop, // Same loop - should not count as container
				}
				polygon.Last = polygon.First

				bbox := createBboxForLoop(polygon.First)

				return testLoop, []*LinkedGeoPolygon{polygon}, []*BBox{bbox}
			},
			expectedCount: 0,
		},
		{
			name: "loop inside one polygon",
			setupTest: func() (*LinkedGeoLoop, []*LinkedGeoPolygon, []*BBox) {
				// Small test loop inside
				testLoop := createSquareLoop(0.025, 0.025)

				// Larger polygon containing the test loop
				polygon := &LinkedGeoPolygon{
					First: createSquareLoop(0, 0), // Contains test loop
				}
				polygon.Last = polygon.First

				bbox := createBboxForLoop(polygon.First)

				return testLoop, []*LinkedGeoPolygon{polygon}, []*BBox{bbox}
			},
			expectedCount: 1,
		},
		{
			name: "loop inside multiple nested polygons",
			setupTest: func() (*LinkedGeoLoop, []*LinkedGeoPolygon, []*BBox) {
				// Small test loop at center
				testLoop := createSquareLoop(0.04, 0.04)

				// Create nested polygons (each contains the test loop)
				polygon1 := &LinkedGeoPolygon{
					First: createSquareLoop(0, 0), // Largest
				}
				polygon1.Last = polygon1.First

				polygon2 := &LinkedGeoPolygon{
					First: createSquareLoop(0.02, 0.02), // Medium
				}
				polygon2.Last = polygon2.First

				polygon3 := &LinkedGeoPolygon{
					First: createSquareLoop(0.03, 0.03), // Smallest but still contains test
				}
				polygon3.Last = polygon3.First

				bbox1 := createBboxForLoop(polygon1.First)
				bbox2 := createBboxForLoop(polygon2.First)
				bbox3 := createBboxForLoop(polygon3.First)

				return testLoop,
					[]*LinkedGeoPolygon{polygon1, polygon2, polygon3},
					[]*BBox{bbox1, bbox2, bbox3}
			},
			expectedCount: 3,
		},
		{
			name: "nil loop",
			setupTest: func() (*LinkedGeoLoop, []*LinkedGeoPolygon, []*BBox) {
				polygon := &LinkedGeoPolygon{
					First: createSquareLoop(0, 0),
				}
				polygon.Last = polygon.First

				bbox := createBboxForLoop(polygon.First)

				return nil, []*LinkedGeoPolygon{polygon}, []*BBox{bbox}
			},
			expectedCount: 0,
		},
		{
			name: "loop with nil first coordinate",
			setupTest: func() (*LinkedGeoLoop, []*LinkedGeoPolygon, []*BBox) {
				testLoop := &LinkedGeoLoop{
					First: nil,
					Last:  nil,
				}

				polygon := &LinkedGeoPolygon{
					First: createSquareLoop(0, 0),
				}
				polygon.Last = polygon.First

				bbox := createBboxForLoop(polygon.First)

				return testLoop, []*LinkedGeoPolygon{polygon}, []*BBox{bbox}
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
