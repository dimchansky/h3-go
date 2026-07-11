//go:build cgo && c2go

package h3

import (
	"testing"
)

func Test_cellToBoundary_parity(t *testing.T) {
	tests := []struct {
		name string
		h3   h3Index
		desc string
	}{
		{
			name: "valid hexagon resolution 0",
			h3:   0x8001fffffffffff,
			desc: "Base cell 1 at resolution 0 (hexagon)",
		},
		{
			name: "valid hexagon resolution 1",
			h3:   0x8101fffffffffff,
			desc: "Resolution 1 hexagon cell",
		},
		{
			name: "valid hexagon resolution 3",
			h3:   0x8301fffffffffff,
			desc: "Resolution 3 hexagon cell",
		},
		{
			name: "pentagon resolution 0",
			h3:   0x8201fffffffffff,
			desc: "Base cell 4 at resolution 0 (pentagon)",
		},
		{
			name: "pentagon resolution 1",
			h3:   0x81083ffffffffff,
			desc: "Resolution 1 pentagon cell (base cell 4)",
		},
		{
			name: "pentagon resolution 2",
			h3:   0x821c07fffffffff,
			desc: "Resolution 2 pentagon cell (base cell 14)",
		},
		{
			name: "high resolution hexagon",
			h3:   0x870830828ffffff,
			desc: "Resolution 7 hexagon for precision test",
		},
		{
			name: "valid cell at north pole region",
			h3:   0x8000fffffffffff,
			desc: "Base cell 0 (near north pole)",
		},
		{
			name: "different base cell hexagon",
			h3:   0x8012fffffffffff,
			desc: "Base cell 18 hexagon",
		},
		{
			name: "different pentagon base cell",
			h3:   0x8601fffffffffff,
			desc: "Base cell 58 pentagon",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare Go boundary result
			var goBoundary CellBoundary
			goBoundary.Verts = make([]LatLng, 10) // Allocate space for vertices

			// Prepare C boundary result
			var cBoundary CellBoundary
			cBoundary.Verts = make([]LatLng, 10) // Allocate space for vertices

			// Call Go implementation
			goErr := cellToBoundary(tt.h3, &goBoundary)

			// Call C implementation
			cErr := h3Error(cellToBoundaryC(tt.h3, &cBoundary))

			// Compare errors
			if goErr != cErr {
				t.Errorf("Error mismatch: Go=%d, C=%d", goErr, cErr)
				return
			}

			if goErr != eSuccess {
				t.Logf("Expected error for %s: %d", tt.desc, goErr)
				return
			}

			// Compare number of vertices
			if goBoundary.NumVerts != cBoundary.NumVerts {
				t.Errorf("Vertex count mismatch: Go=%d, C=%d", goBoundary.NumVerts, cBoundary.NumVerts)
				return
			}

			// Verify reasonable vertex count (5-10 for pentagons with edge intersections, 6-12 for hexagons)
			if goBoundary.NumVerts < 5 || goBoundary.NumVerts > 12 {
				t.Errorf("Unexpected vertex count: got=%d, expected 5-12", goBoundary.NumVerts)
			}

			// Compare each vertex with floating point tolerance
			const tolerance = 1e-12 // Floating point precision tolerance for lat/lng comparisons
			for i := int32(0); i < goBoundary.NumVerts; i++ {
				goVert := goBoundary.Verts[i]
				cVert := cBoundary.Verts[i]

				latDiff := goVert.Lat - cVert.Lat
				if latDiff < 0 {
					latDiff = -latDiff
				}
				if latDiff > tolerance {
					t.Errorf("Vertex %d latitude mismatch: Go=%.17f, C=%.17f, diff=%.17f",
						i, goVert.Lat, cVert.Lat, latDiff)
				}

				lngDiff := goVert.Lng - cVert.Lng
				if lngDiff < 0 {
					lngDiff = -lngDiff
				}
				if lngDiff > tolerance {
					t.Errorf("Vertex %d longitude mismatch: Go=%.17f, C=%.17f, diff=%.17f",
						i, goVert.Lng, cVert.Lng, lngDiff)
				}
			}

			t.Logf("Generated %d vertices for %s (h3Index: 0x%x)",
				goBoundary.NumVerts, tt.desc, uint64(tt.h3))
		})
	}
}

func Test_cellToBoundary_invalid_cells_parity(t *testing.T) {
	invalidCells := []struct {
		name string
		h3   h3Index
		desc string
	}{
		{
			name: "h3Null",
			h3:   0,
			desc: "Null H3 index",
		},
		{
			name: "invalid mode bits",
			h3:   0x1001fffffffffff,
			desc: "Invalid mode bits",
		},
		{
			name: "invalid base cell",
			h3:   0x80ffffffffffffe,
			desc: "Base cell out of range",
		},
	}

	for _, tt := range invalidCells {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare Go boundary result
			var goBoundary CellBoundary
			goBoundary.Verts = make([]LatLng, 10)

			// Prepare C boundary result
			var cBoundary CellBoundary
			cBoundary.Verts = make([]LatLng, 10)

			// Call Go implementation
			goErr := cellToBoundary(tt.h3, &goBoundary)

			// Call C implementation
			cErr := h3Error(cellToBoundaryC(tt.h3, &cBoundary))

			// Compare errors
			if goErr != cErr {
				t.Errorf("Error mismatch for %s: Go=%d, C=%d", tt.desc, goErr, cErr)
			}

			if goErr == eSuccess {
				t.Logf("Unexpected success for invalid cell %s", tt.desc)
			}
		})
	}
}
