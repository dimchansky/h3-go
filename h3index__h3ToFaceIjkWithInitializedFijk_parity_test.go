//go:build cgo && c2go

package h3

import (
	"testing"
)

func Test_h3ToFaceIjkWithInitializedFijk_Parity(t *testing.T) {
	// Test with a variety of H3 indexes and initialized faceIJK addresses
	testCases := []struct {
		name    string
		h       h3Index
		face    int32
		i, j, k int32
	}{
		{"base cell 0, res 0", 0x8001fffffffffff, 0, 0, 0, 0},
		{"base cell 1, res 1", 0x81083ffffffffff, 0, 0, 0, 0},
		{"base cell 4, res 2", 0x820843fffffffff, 1, 0, 0, 0}, // Pentagon base cell
		{"base cell 10, res 3", 0x8308a3fffffffff, 2, 0, 0, 0},
		{"base cell 20, res 4", 0x8409423ffffffff, 3, 0, 0, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Initialize faceIJK for Go test
			goFijk := faceIJK{
				Face:  tc.face,
				Coord: coordIJK{I: tc.i, J: tc.j, K: tc.k},
			}

			// Initialize faceIJK for C test
			cFijk := faceIJK{
				Face:  tc.face,
				Coord: coordIJK{I: tc.i, J: tc.j, K: tc.k},
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
	pentagonH := h3Index(0x820843fffffffff) // Base cell 4 (pentagon)
	fijk := faceIJK{Face: 1, Coord: coordIJK{I: 0, J: 0, K: 0}}

	goResult := _h3ToFaceIjkWithInitializedFijk(pentagonH, &fijk)

	// Pentagon should have possible overage
	if goResult != 1 {
		t.Errorf("Pentagon should have possible overage, got %d", goResult)
	}

	// Test resolution 0 hexagon (should have no overage)
	hexRes0 := h3Index(0x8001fffffffffff) // Base cell 0, res 0
	fijk2 := faceIJK{Face: 0, Coord: coordIJK{I: 0, J: 0, K: 0}}

	goResult2 := _h3ToFaceIjkWithInitializedFijk(hexRes0, &fijk2)

	// Hex at res 0 with (0,0,0) coords should have no overage
	if goResult2 != 0 {
		t.Errorf("Hex res 0 at (0,0,0) should have no overage, got %d", goResult2)
	}
}
