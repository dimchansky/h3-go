//go:build cgo && c2go

package h3

import (
	"testing"
)

func Test_faceIjkPentToCellBoundary_parity(t *testing.T) {
	tests := []struct {
		name   string
		fijk   faceIJK
		res    int32
		start  int32
		length int32
	}{
		{
			name:   "pentagon face 0 res 0",
			fijk:   faceIJK{Face: 0, Coord: coordIJK{I: 0, J: 0, K: 0}},
			res:    0,
			start:  0,
			length: numPentVerts,
		},
		{
			name:   "pentagon face 2 res 2 class II",
			fijk:   faceIJK{Face: 2, Coord: coordIJK{I: 0, J: 1, K: 0}},
			res:    2,
			start:  0,
			length: numPentVerts,
		},
		{
			name:   "pentagon face 10 res 4 class II",
			fijk:   faceIJK{Face: 10, Coord: coordIJK{I: 1, J: 1, K: 1}},
			res:    4,
			start:  0,
			length: numPentVerts,
		},
		{
			name:   "partial pentagon boundary starting at vertex 1, length 3",
			fijk:   faceIJK{Face: 0, Coord: coordIJK{I: 0, J: 0, K: 0}},
			res:    0,
			start:  1,
			length: 3,
		},
		{
			name:   "single pentagon vertex",
			fijk:   faceIJK{Face: 0, Coord: coordIJK{I: 0, J: 0, K: 0}},
			res:    0,
			start:  0,
			length: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare Go boundary result
			var goBoundary CellBoundary
			goBoundary.Verts = make([]LatLng, 20) // Generous allocation for potential intersection vertices

			// Prepare C boundary result
			var cBoundary CellBoundary
			cBoundary.Verts = make([]LatLng, 20) // Generous allocation for potential intersection vertices

			// Call Go implementation
			_faceIjkPentToCellBoundary(&tt.fijk, tt.res, tt.start, tt.length, &goBoundary)

			// Call C implementation
			_faceIjkPentToCellBoundaryC(&tt.fijk, tt.res, tt.start, tt.length, &cBoundary)

			// Compare number of vertices
			if goBoundary.NumVerts != cBoundary.NumVerts {
				t.Errorf("Vertex count mismatch: Go=%d, C=%d", goBoundary.NumVerts, cBoundary.NumVerts)
			}

			// Compare each vertex with floating point tolerance
			minVerts := goBoundary.NumVerts
			if cBoundary.NumVerts < minVerts {
				minVerts = cBoundary.NumVerts
			}

			const tolerance = 1e-12 // High precision tolerance for lat/lng comparisons
			for i := int32(0); i < minVerts; i++ {
				goVert := goBoundary.Verts[i]
				cVert := cBoundary.Verts[i]

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

			t.Logf("Generated %d vertices for pentagon face=%d, ijk=(%d,%d,%d), res=%d",
				goBoundary.NumVerts, tt.fijk.Face, tt.fijk.Coord.I, tt.fijk.Coord.J, tt.fijk.Coord.K, tt.res)
		})
	}
}

// TODO: Add Class III pentagon edge-crossing tests
//
// Class III pentagon boundary generation involves complex edge-crossing logic when
// pentagon edges cross icosahedron face boundaries. The current implementation
// encounters assertion failures in the C reference during testing, specifically:
//   assert(adjacentFaceDir[tmpFijk.face][fijk.face] == quadKI)
//
// This suggests that the edge-crossing logic assumes specific adjacency relationships
// between faces that may not hold for all pentagon vertex combinations. Further
// investigation is needed to:
//   1. Identify valid Class III pentagon faceIJK coordinates that satisfy adjacency requirements
//   2. Understand the geometric constraints that cause the assertion failure
//   3. Determine if this is a limitation of the test data or the algorithm itself
//
// The Class II cases (even resolutions) work correctly as they don't trigger
// the edge-crossing logic, providing reliable pentagon boundaries for those cases.
