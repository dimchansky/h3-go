//go:build cgo && c2go

package h3

import (
	"testing"
)

func Test_findPolygonForHole_parity(t *testing.T) {
	// Helper to create a simple square loop
	createSquareLoop := func(x, y, size float64) *linkedGeoLoop {
		loop := &linkedGeoLoop{}
		coords := []LatLng{
			{Lat: Rad(x), Lng: Rad(y)},
			{Lat: Rad(x), Lng: Rad(y + size)},
			{Lat: Rad(x + size), Lng: Rad(y + size)},
			{Lat: Rad(x + size), Lng: Rad(y)},
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

	// Helper to create a linked list of polygons
	createPolygonList := func(polygons []*linkedGeoPolygon) *linkedGeoPolygon {
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
		setupTest     func() (*linkedGeoLoop, []*linkedGeoPolygon, []*bbox, int32) // returns hole, polygons, bboxes, expected index
		expectedIndex int
	}{
		{
			name: "no polygons",
			setupTest: func() (*linkedGeoLoop, []*linkedGeoPolygon, []*bbox, int32) {
				hole := createSquareLoop(0.04, 0.04, 0.02)
				return hole, []*linkedGeoPolygon{}, []*bbox{}, -1
			},
			expectedIndex: -1,
		},
		{
			name: "nil hole",
			setupTest: func() (*linkedGeoLoop, []*linkedGeoPolygon, []*bbox, int32) {
				polygon := &linkedGeoPolygon{
					First: createSquareLoop(0, 0, 0.1),
				}
				polygon.Last = polygon.First

				bb := createBboxForLoop(polygon.First)

				return nil, []*linkedGeoPolygon{polygon}, []*bbox{bb}, -1
			},
			expectedIndex: -1,
		},
		{
			name: "hole with nil first coordinate",
			setupTest: func() (*linkedGeoLoop, []*linkedGeoPolygon, []*bbox, int32) {
				hole := &linkedGeoLoop{
					First: nil,
					Last:  nil,
				}

				polygon := &linkedGeoPolygon{
					First: createSquareLoop(0, 0, 0.1),
				}
				polygon.Last = polygon.First

				bb := createBboxForLoop(polygon.First)

				return hole, []*linkedGeoPolygon{polygon}, []*bbox{bb}, -1
			},
			expectedIndex: -1,
		},
		{
			name: "hole not inside any polygon",
			setupTest: func() (*linkedGeoLoop, []*linkedGeoPolygon, []*bbox, int32) {
				// Hole far from any polygon
				hole := createSquareLoop(10, 10, 0.02)

				polygon := &linkedGeoPolygon{
					First: createSquareLoop(0, 0, 0.1),
				}
				polygon.Last = polygon.First

				bb := createBboxForLoop(polygon.First)

				return hole, []*linkedGeoPolygon{polygon}, []*bbox{bb}, -1
			},
			expectedIndex: -1,
		},
		{
			name: "hole inside single polygon",
			setupTest: func() (*linkedGeoLoop, []*linkedGeoPolygon, []*bbox, int32) {
				// Small hole inside polygon
				hole := createSquareLoop(0.04, 0.04, 0.02)

				polygon := &linkedGeoPolygon{
					First: createSquareLoop(0, 0, 0.1), // Contains hole
				}
				polygon.Last = polygon.First

				bb := createBboxForLoop(polygon.First)

				return hole, []*linkedGeoPolygon{polygon}, []*bbox{bb}, 0
			},
			expectedIndex: 0,
		},
		{
			name: "hole inside multiple nested polygons - find most deeply nested",
			setupTest: func() (*linkedGeoLoop, []*linkedGeoPolygon, []*bbox, int32) {
				// Small hole at center
				hole := createSquareLoop(0.045, 0.045, 0.01)

				// Create nested polygons (each contains the hole)
				polygon1 := &linkedGeoPolygon{
					First: createSquareLoop(0, 0, 0.1), // Largest
				}
				polygon1.Last = polygon1.First

				polygon2 := &linkedGeoPolygon{
					First: createSquareLoop(0.02, 0.02, 0.06), // Medium
				}
				polygon2.Last = polygon2.First

				polygon3 := &linkedGeoPolygon{
					First: createSquareLoop(0.03, 0.03, 0.04), // Smallest but still contains hole
				}
				polygon3.Last = polygon3.First

				bbox1 := createBboxForLoop(polygon1.First)
				bbox2 := createBboxForLoop(polygon2.First)
				bbox3 := createBboxForLoop(polygon3.First)

				// The most deeply nested polygon (polygon3) should be selected
				return hole,
					[]*linkedGeoPolygon{polygon1, polygon2, polygon3},
					[]*bbox{bbox1, bbox2, bbox3}, 2
			},
			expectedIndex: 2,
		},
		{
			name: "hole inside some but not all polygons",
			setupTest: func() (*linkedGeoLoop, []*linkedGeoPolygon, []*bbox, int32) {
				// Hole positioned to be inside polygon1 and polygon2 but not polygon3
				hole := createSquareLoop(0.025, 0.025, 0.01)

				polygon1 := &linkedGeoPolygon{
					First: createSquareLoop(0, 0, 0.1), // Contains hole
				}
				polygon1.Last = polygon1.First

				polygon2 := &linkedGeoPolygon{
					First: createSquareLoop(0.01, 0.01, 0.08), // Contains hole (more nested)
				}
				polygon2.Last = polygon2.First

				polygon3 := &linkedGeoPolygon{
					First: createSquareLoop(0.05, 0.05, 0.04), // Does NOT contain hole
				}
				polygon3.Last = polygon3.First

				bbox1 := createBboxForLoop(polygon1.First)
				bbox2 := createBboxForLoop(polygon2.First)
				bbox3 := createBboxForLoop(polygon3.First)

				// polygon2 should be selected (more deeply nested than polygon1)
				return hole,
					[]*linkedGeoPolygon{polygon1, polygon2, polygon3},
					[]*bbox{bbox1, bbox2, bbox3}, 1
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
			var expected *linkedGeoPolygon
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
