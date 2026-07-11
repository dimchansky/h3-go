//go:build cgo && c2go

package h3

import (
	"testing"
)

func Test_findDeepestContainer_parity(t *testing.T) {
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

	tests := []struct {
		name          string
		setupTest     func() ([]*linkedGeoPolygon, []*bbox, int32) // returns polygons, bboxes, expected index
		expectedIndex int
	}{
		{
			name: "empty polygons array",
			setupTest: func() ([]*linkedGeoPolygon, []*bbox, int32) {
				return []*linkedGeoPolygon{}, []*bbox{}, -1 // -1 indicates nil expected
			},
			expectedIndex: -1,
		},
		{
			name: "single polygon",
			setupTest: func() ([]*linkedGeoPolygon, []*bbox, int32) {
				polygon := &linkedGeoPolygon{
					First: createSquareLoop(0, 0, 0.1),
				}
				polygon.Last = polygon.First

				bb := createBboxForLoop(polygon.First)

				return []*linkedGeoPolygon{polygon}, []*bbox{bb}, 0
			},
			expectedIndex: 0,
		},
		{
			name: "two nested polygons - outer and inner",
			setupTest: func() ([]*linkedGeoPolygon, []*bbox, int32) {
				// Larger outer polygon
				outerPolygon := &linkedGeoPolygon{
					First: createSquareLoop(0, 0, 0.1), // Large outer square
				}
				outerPolygon.Last = outerPolygon.First

				// Smaller inner polygon (more deeply nested)
				innerPolygon := &linkedGeoPolygon{
					First: createSquareLoop(0.02, 0.02, 0.06), // Smaller inner square
				}
				innerPolygon.Last = innerPolygon.First

				outerBbox := createBboxForLoop(outerPolygon.First)
				innerBbox := createBboxForLoop(innerPolygon.First)

				// The inner polygon should be the deepest (it has 1 container - the outer)
				// The outer polygon has 0 containers
				return []*linkedGeoPolygon{outerPolygon, innerPolygon}, []*bbox{outerBbox, innerBbox}, 1
			},
			expectedIndex: 1,
		},
		{
			name: "three nested polygons - find the most deeply nested",
			setupTest: func() ([]*linkedGeoPolygon, []*bbox, int32) {
				// Create three nested polygons: large, medium, small
				largePolygon := &linkedGeoPolygon{
					First: createSquareLoop(0, 0, 0.1), // Largest
				}
				largePolygon.Last = largePolygon.First

				mediumPolygon := &linkedGeoPolygon{
					First: createSquareLoop(0.02, 0.02, 0.06), // Medium
				}
				mediumPolygon.Last = mediumPolygon.First

				smallPolygon := &linkedGeoPolygon{
					First: createSquareLoop(0.03, 0.03, 0.04), // Smallest (most deeply nested)
				}
				smallPolygon.Last = smallPolygon.First

				largeBbox := createBboxForLoop(largePolygon.First)
				mediumBbox := createBboxForLoop(mediumPolygon.First)
				smallBbox := createBboxForLoop(smallPolygon.First)

				// Small polygon has 2 containers (large + medium)
				// Medium polygon has 1 container (large)
				// Large polygon has 0 containers
				// So small polygon should be selected
				return []*linkedGeoPolygon{largePolygon, mediumPolygon, smallPolygon},
					[]*bbox{largeBbox, mediumBbox, smallBbox}, 2
			},
			expectedIndex: 2,
		},
		{
			name: "non-nested polygons - return first",
			setupTest: func() ([]*linkedGeoPolygon, []*bbox, int32) {
				// Two separate polygons that don't contain each other
				polygon1 := &linkedGeoPolygon{
					First: createSquareLoop(0, 0, 0.05), // Left side
				}
				polygon1.Last = polygon1.First

				polygon2 := &linkedGeoPolygon{
					First: createSquareLoop(0.1, 0.1, 0.05), // Right side, far away
				}
				polygon2.Last = polygon2.First

				bbox1 := createBboxForLoop(polygon1.First)
				bbox2 := createBboxForLoop(polygon2.First)

				// Neither polygon contains the other, so both have 0 containers
				// Function should return the first polygon
				return []*linkedGeoPolygon{polygon1, polygon2}, []*bbox{bbox1, bbox2}, 0
			},
			expectedIndex: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test data
			polygons, bboxes, expectedIndex := tt.setupTest()

			// Test Go implementation
			goResult := findDeepestContainer(polygons, bboxes)

			// Test C implementation
			cResult := findDeepestContainerC(polygons, bboxes)

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
