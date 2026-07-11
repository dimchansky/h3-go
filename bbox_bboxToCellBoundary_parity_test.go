//go:build cgo && c2go

package h3

import (
	"math"
	"testing"
)

func Test_bboxToCellBoundary_parity(t *testing.T) {
	tests := []struct {
		name string
		bbox bbox
	}{
		{
			name: "Simple bbox around origin",
			bbox: bbox{
				North: Angle(0.1),
				South: Angle(-0.1),
				East:  Angle(0.1),
				West:  Angle(-0.1),
			},
		},
		{
			name: "Bbox spanning positive longitude",
			bbox: bbox{
				North: Angle(0.5),
				South: Angle(0.3),
				East:  Angle(1.2),
				West:  Angle(0.8),
			},
		},
		{
			name: "Bbox spanning negative longitude",
			bbox: bbox{
				North: Angle(-0.3),
				South: Angle(-0.5),
				East:  Angle(-0.8),
				West:  Angle(-1.2),
			},
		},
		{
			name: "Large bbox covering significant area",
			bbox: bbox{
				North: Angle(1.5),
				South: Angle(-1.5),
				East:  Angle(3.0),
				West:  Angle(-3.0),
			},
		},
		{
			name: "Very small bbox",
			bbox: bbox{
				North: Angle(0.001),
				South: Angle(0.0),
				East:  Angle(0.001),
				West:  Angle(0.0),
			},
		},
		{
			name: "Bbox near north pole",
			bbox: bbox{
				North: PiOver2,
				South: Angle(1.4),
				East:  Pi,
				West:  Angle(-Pi),
			},
		},
		{
			name: "Bbox near south pole",
			bbox: bbox{
				North: Angle(-1.4),
				South: -PiOver2,
				East:  Pi,
				West:  Angle(-Pi),
			},
		},
		{
			name: "Bbox crossing antimeridian",
			bbox: bbox{
				North: Angle(0.5),
				South: Angle(-0.5),
				East:  Angle(-2.5),
				West:  Angle(2.5),
			},
		},
		{
			name: "Square bbox",
			bbox: bbox{
				North: Angle(1.0),
				South: Angle(0.0),
				East:  Angle(1.0),
				West:  Angle(0.0),
			},
		},
		{
			name: "Zero-width bbox (line)",
			bbox: bbox{
				North: Angle(1.0),
				South: Angle(0.0),
				East:  Angle(0.5),
				West:  Angle(0.5),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call Go implementation
			goBoundary := bboxToCellBoundary(&tt.bbox)

			// Call C implementation
			cBoundary := bboxToCellBoundaryC(&tt.bbox)

			// Compare number of vertices
			if goBoundary.numVerts != cBoundary.numVerts {
				t.Errorf("NumVerts mismatch: Go=%d, C=%d", goBoundary.numVerts, cBoundary.numVerts)
				return
			}

			// Compare each vertex
			tolerance := 1e-15 // Very tight tolerance for simple coordinate assignments
			for i := int32(0); i < goBoundary.numVerts; i++ {
				goVert := goBoundary.verts[i]
				cVert := cBoundary.verts[i]

				if math.Abs(float64(goVert.Lat-cVert.Lat)) > tolerance {
					t.Errorf("Vertex %d Lat mismatch: Go=%.15f, C=%.15f, diff=%.15f",
						i, float64(goVert.Lat), float64(cVert.Lat), float64(goVert.Lat-cVert.Lat))
				}
				if math.Abs(float64(goVert.Lng-cVert.Lng)) > tolerance {
					t.Errorf("Vertex %d Lng mismatch: Go=%.15f, C=%.15f, diff=%.15f",
						i, float64(goVert.Lng), float64(cVert.Lng), float64(goVert.Lng-cVert.Lng))
				}
			}
		})
	}
}

func Test_bboxToCellBoundary_vertex_order_parity(t *testing.T) {
	// Test specific vertex ordering (counter-clockwise)
	bbox := bbox{
		North: Angle(1.0),
		South: Angle(0.0),
		East:  Angle(1.0),
		West:  Angle(0.0),
	}

	// Call both implementations
	goBoundary := bboxToCellBoundary(&bbox)
	cBoundary := bboxToCellBoundaryC(&bbox)

	// Both should have 4 vertices
	if goBoundary.numVerts != 4 || cBoundary.numVerts != 4 {
		t.Fatalf("Expected 4 vertices, got Go=%d, C=%d", goBoundary.numVerts, cBoundary.numVerts)
	}

	// Expected vertex order (CCW): NE, NW, SW, SE
	expectedOrder := []struct {
		name string
		lat  Angle
		lng  Angle
	}{
		{"NE", bbox.North, bbox.East},
		{"NW", bbox.North, bbox.West},
		{"SW", bbox.South, bbox.West},
		{"SE", bbox.South, bbox.East},
	}

	tolerance := 1e-15
	for i, expected := range expectedOrder {
		// Check Go implementation
		if math.Abs(float64(goBoundary.verts[i].Lat-expected.lat)) > tolerance ||
			math.Abs(float64(goBoundary.verts[i].Lng-expected.lng)) > tolerance {
			t.Errorf("Go vertex %d (%s): expected (%.15f, %.15f), got (%.15f, %.15f)",
				i, expected.name, float64(expected.lat), float64(expected.lng),
				float64(goBoundary.verts[i].Lat), float64(goBoundary.verts[i].Lng))
		}

		// Check C implementation
		if math.Abs(float64(cBoundary.verts[i].Lat-expected.lat)) > tolerance ||
			math.Abs(float64(cBoundary.verts[i].Lng-expected.lng)) > tolerance {
			t.Errorf("C vertex %d (%s): expected (%.15f, %.15f), got (%.15f, %.15f)",
				i, expected.name, float64(expected.lat), float64(expected.lng),
				float64(cBoundary.verts[i].Lat), float64(cBoundary.verts[i].Lng))
		}
	}
}

func Test_bboxToCellBoundary_comprehensive_parity(t *testing.T) {
	// Test with various coordinate ranges and edge cases
	testCases := []bbox{
		// Normal cases
		{North: Angle(0.5), South: Angle(0.0), East: Angle(0.5), West: Angle(0.0)},
		{North: Angle(-0.5), South: Angle(-1.0), East: Angle(-0.5), West: Angle(-1.0)},

		// Edge cases with pole boundaries
		{North: PiOver2, South: Angle(1.0), East: Pi, West: Angle(-Pi)},
		{North: Angle(-1.0), South: -PiOver2, East: Pi, West: Angle(-Pi)},

		// Cases with longitude wrap
		{North: Angle(0.1), South: Angle(-0.1), East: Angle(-3.0), West: Angle(3.0)},

		// Very small bboxes
		{North: Angle(1e-10), South: Angle(0.0), East: Angle(1e-10), West: Angle(0.0)},

		// Large bboxes
		{North: Angle(1.5), South: Angle(-1.5), East: Angle(3.0), West: Angle(-3.0)},
	}

	for i, bbox := range testCases {
		t.Run("Case_"+string(rune(i+'0')), func(t *testing.T) {
			// Call both implementations
			goBoundary := bboxToCellBoundary(&bbox)
			cBoundary := bboxToCellBoundaryC(&bbox)

			// Compare results
			if goBoundary.numVerts != cBoundary.numVerts {
				t.Errorf("NumVerts mismatch: Go=%d, C=%d", goBoundary.numVerts, cBoundary.numVerts)
				return
			}

			tolerance := 1e-15
			for j := int32(0); j < goBoundary.numVerts; j++ {
				if math.Abs(float64(goBoundary.verts[j].Lat-cBoundary.verts[j].Lat)) > tolerance ||
					math.Abs(float64(goBoundary.verts[j].Lng-cBoundary.verts[j].Lng)) > tolerance {
					t.Errorf("Vertex %d mismatch: Go=(%.15f, %.15f), C=(%.15f, %.15f)",
						j, float64(goBoundary.verts[j].Lat), float64(goBoundary.verts[j].Lng),
						float64(cBoundary.verts[j].Lat), float64(cBoundary.verts[j].Lng))
				}
			}
		})
	}
}
