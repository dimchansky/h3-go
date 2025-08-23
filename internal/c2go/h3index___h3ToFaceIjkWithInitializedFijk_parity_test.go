//go:build cgo

package c2go

import (
	"testing"
)

func Test_h3ToFaceIjkWithInitializedFijk_Parity(t *testing.T) {
	// Test with a variety of H3 indexes and initialized FaceIJK addresses
	testCases := []struct {
		name    string
		h       H3Index
		face    int
		i, j, k int
	}{
		{"base cell 0, res 0", 0x8001fffffffffff, 0, 0, 0, 0},
		{"base cell 1, res 1", 0x81083ffffffffff, 0, 0, 0, 0},
		{"base cell 4, res 2", 0x820843fffffffff, 1, 0, 0, 0}, // Pentagon base cell
		{"base cell 10, res 3", 0x8308a3fffffffff, 2, 0, 0, 0},
		{"base cell 20, res 4", 0x8409423ffffffff, 3, 0, 0, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Initialize FaceIJK for Go test
			goFijk := FaceIJK{
				Face:  tc.face,
				Coord: CoordIJK{I: tc.i, J: tc.j, K: tc.k},
			}

			// Initialize FaceIJK for C test
			cFijk := FaceIJK{
				Face:  tc.face,
				Coord: CoordIJK{I: tc.i, J: tc.j, K: tc.k},
			}

			// Run both implementations
			goResult := _h3ToFaceIjkWithInitializedFijk(tc.h, &goFijk)
			cResult := _h3ToFaceIjkWithInitializedFijkC(tc.h, &cFijk)

			// Compare return values
			if goResult != cResult {
				t.Errorf("Return value mismatch: Go=%d, C=%d", goResult, cResult)
			}

			// Compare face
			if goFijk.Face != cFijk.Face {
				t.Errorf("Face mismatch: Go=%d, C=%d", goFijk.Face, cFijk.Face)
			}

			// Compare coordinates
			if goFijk.Coord.I != cFijk.Coord.I {
				t.Errorf("I mismatch: Go=%d, C=%d", goFijk.Coord.I, cFijk.Coord.I)
			}
			if goFijk.Coord.J != cFijk.Coord.J {
				t.Errorf("J mismatch: Go=%d, C=%d", goFijk.Coord.J, cFijk.Coord.J)
			}
			if goFijk.Coord.K != cFijk.Coord.K {
				t.Errorf("K mismatch: Go=%d, C=%d", goFijk.Coord.K, cFijk.Coord.K)
			}
		})
	}
}

func Test_h3ToFaceIjkWithInitializedFijk_EdgeCases(t *testing.T) {
	// Test pentagon handling
	pentagonH := H3Index(0x820843fffffffff) // Base cell 4 (pentagon)
	fijk := FaceIJK{Face: 1, Coord: CoordIJK{I: 0, J: 0, K: 0}}

	goResult := _h3ToFaceIjkWithInitializedFijk(pentagonH, &fijk)

	// Pentagon should have possible overage
	if goResult != 1 {
		t.Errorf("Pentagon should have possible overage, got %d", goResult)
	}

	// Test resolution 0 hexagon (should have no overage)
	hexRes0 := H3Index(0x8001fffffffffff) // Base cell 0, res 0
	fijk2 := FaceIJK{Face: 0, Coord: CoordIJK{I: 0, J: 0, K: 0}}

	goResult2 := _h3ToFaceIjkWithInitializedFijk(hexRes0, &fijk2)

	// Hex at res 0 with (0,0,0) coords should have no overage
	if goResult2 != 0 {
		t.Errorf("Hex res 0 at (0,0,0) should have no overage, got %d", goResult2)
	}
}
