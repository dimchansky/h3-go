//go:build cgo && c2go

package h3

import (
	"testing"
)

func Test_findPolygonForHole_parity(t *testing.T) {
	// Helper to create a simple square loop
	createSquareLoop := func(x, y, size float64) *LinkedGeoLoop {
		loop := &LinkedGeoLoop{}
		coords := []LatLng{
			{Lat: Rad(x), Lng: Rad(y)},
			{Lat: Rad(x), Lng: Rad(y + size)},
			{Lat: Rad(x + size), Lng: Rad(y + size)},
			{Lat: Rad(x + size), Lng: Rad(y)},
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

	// Helper to create a linked list of polygons
	createPolygonList := func(polygons []*LinkedGeoPolygon) *LinkedGeoPolygon {
		if len(polygons) == 0 {
			return nil
		}

		first := polygons[0]
		current := first
		for i := 1; i < len(polygons); i++ {
			current.Next = polygons[i]
			current = polygons[i]
		}

		// Ensure last polygon has no next
		current.Next = nil

		return first
	}

	tests := []struct {
		name          string
		setupTest     func() (*LinkedGeoLoop, []*LinkedGeoPolygon, []*BBox, int32) // returns hole, polygons, bboxes, expected index
		expectedIndex int
	}{
		{
			name: "no polygons",
			setupTest: func() (*LinkedGeoLoop, []*LinkedGeoPolygon, []*BBox, int32) {
				hole := createSquareLoop(0.04, 0.04, 0.02)
				return hole, []*LinkedGeoPolygon{}, []*BBox{}, -1
			},
			expectedIndex: -1,
		},
		{
			name: "nil hole",
			setupTest: func() (*LinkedGeoLoop, []*LinkedGeoPolygon, []*BBox, int32) {
				polygon := &LinkedGeoPolygon{
					First: createSquareLoop(0, 0, 0.1),
				}
				polygon.Last = polygon.First

				bbox := createBboxForLoop(polygon.First)

				return nil, []*LinkedGeoPolygon{polygon}, []*BBox{bbox}, -1
			},
			expectedIndex: -1,
		},
		{
			name: "hole with nil first coordinate",
			setupTest: func() (*LinkedGeoLoop, []*LinkedGeoPolygon, []*BBox, int32) {
				hole := &LinkedGeoLoop{
					First: nil,
					Last:  nil,
				}

				polygon := &LinkedGeoPolygon{
					First: createSquareLoop(0, 0, 0.1),
				}
				polygon.Last = polygon.First

				bbox := createBboxForLoop(polygon.First)

				return hole, []*LinkedGeoPolygon{polygon}, []*BBox{bbox}, -1
			},
			expectedIndex: -1,
		},
		{
			name: "hole not inside any polygon",
			setupTest: func() (*LinkedGeoLoop, []*LinkedGeoPolygon, []*BBox, int32) {
				// Hole far from any polygon
				hole := createSquareLoop(10, 10, 0.02)

				polygon := &LinkedGeoPolygon{
					First: createSquareLoop(0, 0, 0.1),
				}
				polygon.Last = polygon.First

				bbox := createBboxForLoop(polygon.First)

				return hole, []*LinkedGeoPolygon{polygon}, []*BBox{bbox}, -1
			},
			expectedIndex: -1,
		},
		{
			name: "hole inside single polygon",
			setupTest: func() (*LinkedGeoLoop, []*LinkedGeoPolygon, []*BBox, int32) {
				// Small hole inside polygon
				hole := createSquareLoop(0.04, 0.04, 0.02)

				polygon := &LinkedGeoPolygon{
					First: createSquareLoop(0, 0, 0.1), // Contains hole
				}
				polygon.Last = polygon.First

				bbox := createBboxForLoop(polygon.First)

				return hole, []*LinkedGeoPolygon{polygon}, []*BBox{bbox}, 0
			},
			expectedIndex: 0,
		},
		{
			name: "hole inside multiple nested polygons - find most deeply nested",
			setupTest: func() (*LinkedGeoLoop, []*LinkedGeoPolygon, []*BBox, int32) {
				// Small hole at center
				hole := createSquareLoop(0.045, 0.045, 0.01)

				// Create nested polygons (each contains the hole)
				polygon1 := &LinkedGeoPolygon{
					First: createSquareLoop(0, 0, 0.1), // Largest
				}
				polygon1.Last = polygon1.First

				polygon2 := &LinkedGeoPolygon{
					First: createSquareLoop(0.02, 0.02, 0.06), // Medium
				}
				polygon2.Last = polygon2.First

				polygon3 := &LinkedGeoPolygon{
					First: createSquareLoop(0.03, 0.03, 0.04), // Smallest but still contains hole
				}
				polygon3.Last = polygon3.First

				bbox1 := createBboxForLoop(polygon1.First)
				bbox2 := createBboxForLoop(polygon2.First)
				bbox3 := createBboxForLoop(polygon3.First)

				// The most deeply nested polygon (polygon3) should be selected
				return hole,
					[]*LinkedGeoPolygon{polygon1, polygon2, polygon3},
					[]*BBox{bbox1, bbox2, bbox3}, 2
			},
			expectedIndex: 2,
		},
		{
			name: "hole inside some but not all polygons",
			setupTest: func() (*LinkedGeoLoop, []*LinkedGeoPolygon, []*BBox, int32) {
				// Hole positioned to be inside polygon1 and polygon2 but not polygon3
				hole := createSquareLoop(0.025, 0.025, 0.01)

				polygon1 := &LinkedGeoPolygon{
					First: createSquareLoop(0, 0, 0.1), // Contains hole
				}
				polygon1.Last = polygon1.First

				polygon2 := &LinkedGeoPolygon{
					First: createSquareLoop(0.01, 0.01, 0.08), // Contains hole (more nested)
				}
				polygon2.Last = polygon2.First

				polygon3 := &LinkedGeoPolygon{
					First: createSquareLoop(0.05, 0.05, 0.04), // Does NOT contain hole
				}
				polygon3.Last = polygon3.First

				bbox1 := createBboxForLoop(polygon1.First)
				bbox2 := createBboxForLoop(polygon2.First)
				bbox3 := createBboxForLoop(polygon3.First)

				// polygon2 should be selected (more deeply nested than polygon1)
				return hole,
					[]*LinkedGeoPolygon{polygon1, polygon2, polygon3},
					[]*BBox{bbox1, bbox2, bbox3}, 1
			},
			expectedIndex: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test data
			hole, polygons, bboxes, expectedIndex := tt.setupTest()

			// Create linked list from polygons array for Go implementation
			polygonList := createPolygonList(polygons)

			// Test Go implementation
			goResult := findPolygonForHole(hole, polygonList, bboxes)

			// Test C implementation (it uses arrays directly)
			cResult := findPolygonForHoleC(hole, polygons, bboxes)

			// Compare results
			if goResult != cResult {
				t.Errorf("Result mismatch: Go=%v, C=%v", goResult, cResult)
			}

			// Check expected result
			var expected *LinkedGeoPolygon
			if expectedIndex >= 0 && int(expectedIndex) < len(polygons) {
				expected = polygons[expectedIndex]
			}

			if goResult != expected {
				t.Errorf("Go result doesn't match expected: got %v, want %v", goResult, expected)
			}

			if cResult != expected {
				t.Errorf("C result doesn't match expected: got %v, want %v", cResult, expected)
			}
		})
	}
}
