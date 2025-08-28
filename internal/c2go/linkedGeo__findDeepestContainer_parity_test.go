//go:build cgo && c2go

package c2go

import (
	"testing"

	"github.com/dimchansky/h3-go/angle"
)

func Test_findDeepestContainer_parity(t *testing.T) {
	// Helper to create a simple square loop
	createSquareLoop := func(x, y, size float64) *LinkedGeoLoop {
		loop := &LinkedGeoLoop{}
		coords := []LatLng{
			{Lat: angle.Rad(x), Lng: angle.Rad(y)},
			{Lat: angle.Rad(x), Lng: angle.Rad(y + size)},
			{Lat: angle.Rad(x + size), Lng: angle.Rad(y + size)},
			{Lat: angle.Rad(x + size), Lng: angle.Rad(y)},
			{Lat: angle.Rad(x), Lng: angle.Rad(y)}, // Close the loop
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
		bbox := &BBox{North: angle.Deg(-90), South: angle.Deg(90), East: angle.Deg(-180), West: angle.Deg(180)}

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
		setupTest     func() ([]*LinkedGeoPolygon, []*BBox, int32) // returns polygons, bboxes, expected index
		expectedIndex int
	}{
		{
			name: "empty polygons array",
			setupTest: func() ([]*LinkedGeoPolygon, []*BBox, int32) {
				return []*LinkedGeoPolygon{}, []*BBox{}, -1 // -1 indicates nil expected
			},
			expectedIndex: -1,
		},
		{
			name: "single polygon",
			setupTest: func() ([]*LinkedGeoPolygon, []*BBox, int32) {
				polygon := &LinkedGeoPolygon{
					First: createSquareLoop(0, 0, 0.1),
				}
				polygon.Last = polygon.First

				bbox := createBboxForLoop(polygon.First)

				return []*LinkedGeoPolygon{polygon}, []*BBox{bbox}, 0
			},
			expectedIndex: 0,
		},
		{
			name: "two nested polygons - outer and inner",
			setupTest: func() ([]*LinkedGeoPolygon, []*BBox, int32) {
				// Larger outer polygon
				outerPolygon := &LinkedGeoPolygon{
					First: createSquareLoop(0, 0, 0.1), // Large outer square
				}
				outerPolygon.Last = outerPolygon.First

				// Smaller inner polygon (more deeply nested)
				innerPolygon := &LinkedGeoPolygon{
					First: createSquareLoop(0.02, 0.02, 0.06), // Smaller inner square
				}
				innerPolygon.Last = innerPolygon.First

				outerBbox := createBboxForLoop(outerPolygon.First)
				innerBbox := createBboxForLoop(innerPolygon.First)

				// The inner polygon should be the deepest (it has 1 container - the outer)
				// The outer polygon has 0 containers
				return []*LinkedGeoPolygon{outerPolygon, innerPolygon}, []*BBox{outerBbox, innerBbox}, 1
			},
			expectedIndex: 1,
		},
		{
			name: "three nested polygons - find the most deeply nested",
			setupTest: func() ([]*LinkedGeoPolygon, []*BBox, int32) {
				// Create three nested polygons: large, medium, small
				largePolygon := &LinkedGeoPolygon{
					First: createSquareLoop(0, 0, 0.1), // Largest
				}
				largePolygon.Last = largePolygon.First

				mediumPolygon := &LinkedGeoPolygon{
					First: createSquareLoop(0.02, 0.02, 0.06), // Medium
				}
				mediumPolygon.Last = mediumPolygon.First

				smallPolygon := &LinkedGeoPolygon{
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
				return []*LinkedGeoPolygon{largePolygon, mediumPolygon, smallPolygon},
					[]*BBox{largeBbox, mediumBbox, smallBbox}, 2
			},
			expectedIndex: 2,
		},
		{
			name: "non-nested polygons - return first",
			setupTest: func() ([]*LinkedGeoPolygon, []*BBox, int32) {
				// Two separate polygons that don't contain each other
				polygon1 := &LinkedGeoPolygon{
					First: createSquareLoop(0, 0, 0.05), // Left side
				}
				polygon1.Last = polygon1.First

				polygon2 := &LinkedGeoPolygon{
					First: createSquareLoop(0.1, 0.1, 0.05), // Right side, far away
				}
				polygon2.Last = polygon2.First

				bbox1 := createBboxForLoop(polygon1.First)
				bbox2 := createBboxForLoop(polygon2.First)

				// Neither polygon contains the other, so both have 0 containers
				// Function should return the first polygon
				return []*LinkedGeoPolygon{polygon1, polygon2}, []*BBox{bbox1, bbox2}, 0
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
