//go:build cgo && c2go

package h3

import (
	"testing"
)

func Test_faceIjkToCellBoundary_parity(t *testing.T) {
	tests := []struct {
		name   string
		fijk   faceIJK
		res    int32
		start  int32
		length int32
	}{
		{
			name:   "res 0 center cell full boundary",
			fijk:   faceIJK{Face: 0, Coord: coordIJK{0, 0, 0}},
			res:    0,
			start:  0,
			length: numHexVerts,
		},
		{
			name:   "res 1 cell full boundary",
			fijk:   faceIJK{Face: 1, Coord: coordIJK{1, 0, 0}},
			res:    1,
			start:  0,
			length: numHexVerts,
		},
		{
			name:   "res 2 cell full boundary",
			fijk:   faceIJK{Face: 2, Coord: coordIJK{1, 1, 0}},
			res:    2,
			start:  0,
			length: numHexVerts,
		},
		{
			name:   "res 3 class III cell full boundary",
			fijk:   faceIJK{Face: 3, Coord: coordIJK{2, 1, 0}},
			res:    3,
			start:  0,
			length: numHexVerts,
		},
		{
			name:   "res 4 class II cell full boundary",
			fijk:   faceIJK{Face: 4, Coord: coordIJK{3, 2, 1}},
			res:    4,
			start:  0,
			length: numHexVerts,
		},
		{
			name:   "partial boundary starting at vertex 2, length 3",
			fijk:   faceIJK{Face: 0, Coord: coordIJK{1, 0, 0}},
			res:    1,
			start:  2,
			length: 3,
		},
		{
			name:   "single vertex",
			fijk:   faceIJK{Face: 0, Coord: coordIJK{0, 0, 0}},
			res:    0,
			start:  0,
			length: 1,
		},
		{
			name:   "edge case near face boundary",
			fijk:   faceIJK{Face: 5, Coord: coordIJK{6, 3, 0}},
			res:    5,
			start:  0,
			length: numHexVerts,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare Go boundary result
			var goBoundary CellBoundary

			// Prepare C boundary result
			var cBoundary CellBoundary

			// Call Go implementation
			_faceIjkToCellBoundary(&tt.fijk, tt.res, tt.start, tt.length, &goBoundary)

			// Call C implementation
			_faceIjkToCellBoundaryC(&tt.fijk, tt.res, tt.start, tt.length, &cBoundary)

			// Compare number of vertices
			if goBoundary.numVerts != cBoundary.numVerts {
				t.Errorf("Vertex count mismatch: Go=%d, C=%d", goBoundary.numVerts, cBoundary.numVerts)
			}

			// Compare each vertex with floating point tolerance
			minVerts := goBoundary.numVerts
			if cBoundary.numVerts < minVerts {
				minVerts = cBoundary.numVerts
			}

			const tolerance = 1e-12 // High precision tolerance for lat/lng comparisons
			for i := int32(0); i < minVerts; i++ {
				goVert := goBoundary.verts[i]
				cVert := cBoundary.verts[i]

				latDiff := goVert.Lat - cVert.Lat
				if latDiff < 0 {
					latDiff = -latDiff
				}
				if latDiff > tolerance {
					t.Errorf("Vertex %d latitude mismatch: Go=%.15f, C=%.15f, diff=%.15f",
						i, goVert.Lat, cVert.Lat, latDiff)
				}

				lngDiff := goVert.Lng - cVert.Lng
				if lngDiff < 0 {
					lngDiff = -lngDiff
				}
				if lngDiff > tolerance {
					t.Errorf("Vertex %d longitude mismatch: Go=%.15f, C=%.15f, diff=%.15f",
						i, goVert.Lng, cVert.Lng, lngDiff)
				}
			}

			t.Logf("Generated %d vertices for face=%d, ijk=(%d,%d,%d), res=%d",
				goBoundary.numVerts, tt.fijk.Face, tt.fijk.Coord.I, tt.fijk.Coord.J, tt.fijk.Coord.K, tt.res)
		})
	}
}

// Test edge crossing cases specifically for Class III resolutions
func Test_faceIjkToCellBoundary_edge_crossing_parity(t *testing.T) {
	tests := []struct {
		name string
		fijk faceIJK
		res  int32
	}{
		{
			name: "res 1 class III cell near face edge",
			fijk: faceIJK{Face: 0, Coord: coordIJK{2, 0, 0}},
			res:  1,
		},
		{
			name: "res 3 class III cell that might cross faces",
			fijk: faceIJK{Face: 1, Coord: coordIJK{5, 2, 0}},
			res:  3,
		},
		{
			name: "res 5 class III cell at higher resolution",
			fijk: faceIJK{Face: 2, Coord: coordIJK{12, 6, 3}},
			res:  5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Only test Class III resolutions
			if !isResolutionClassIII(tt.res) {
				t.Skipf("Skipping non-Class III resolution %d", tt.res)
			}

			// Prepare boundary results
			var goBoundary CellBoundary

			var cBoundary CellBoundary

			// Test full boundary
			_faceIjkToCellBoundary(&tt.fijk, tt.res, 0, numHexVerts, &goBoundary)
			_faceIjkToCellBoundaryC(&tt.fijk, tt.res, 0, numHexVerts, &cBoundary)

			// Compare results
			if goBoundary.numVerts != cBoundary.numVerts {
				t.Errorf("Vertex count mismatch: Go=%d, C=%d", goBoundary.numVerts, cBoundary.numVerts)
			}

			const tolerance = 1e-12
			for i := int32(0); i < goBoundary.numVerts && i < cBoundary.numVerts; i++ {
				goVert := goBoundary.verts[i]
				cVert := cBoundary.verts[i]

				latDiff := goVert.Lat - cVert.Lat
				if latDiff < 0 {
					latDiff = -latDiff
				}
				lngDiff := goVert.Lng - cVert.Lng
				if lngDiff < 0 {
					lngDiff = -lngDiff
				}

				if latDiff > tolerance || lngDiff > tolerance {
					t.Errorf("Vertex %d mismatch: Go=(%.15f,%.15f), C=(%.15f,%.15f)",
						i, goVert.Lat, goVert.Lng, cVert.Lat, cVert.Lng)
				}
			}

			t.Logf("Class III boundary generated %d vertices (potential edge crossings)",
				goBoundary.numVerts)
		})
	}
}
