//go:build cgo && c2go

package h3

import (
	"testing"
)

func Test_adjustPentVertOverage_parity(t *testing.T) {
	// Test with simple cases that are known to work
	// Note: This function is very sensitive to input and is designed
	// specifically for pentagon vertex coordinates from _faceIjkPentToVerts
	tests := []struct {
		name string
		fijk FaceIJK
		res  int32
	}{
		{
			name: "no overage case face 0",
			fijk: FaceIJK{Face: 0, Coord: CoordIJK{I: 0, J: 0, K: 0}},
			res:  0,
		},
		{
			name: "simple case face 2",
			fijk: FaceIJK{Face: 2, Coord: CoordIJK{I: 2, J: 1, K: 0}},
			res:  2,
		},
		{
			name: "normalized case face 5",
			fijk: FaceIJK{Face: 5, Coord: CoordIJK{I: 1, J: 2, K: 0}},
			res:  4,
		},
		{
			name: "critical difference case",
			fijk: FaceIJK{Face: 5, Coord: CoordIJK{I: 715827882, J: 0, K: 1431655770}},
			res:  3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test Go implementation
			goFijk := tt.fijk
			goResult := _adjustPentVertOverage(&goFijk, tt.res)

			// Test C implementation
			cFijk := tt.fijk
			cResult := _adjustPentVertOverageC(&cFijk, tt.res)

			// Compare overage results
			if goResult != cResult {
				t.Errorf("overage result mismatch: Go=%v, C=%v", goResult, cResult)
			}

			// Compare final FaceIJK coordinates
			if goFijk.Face != cFijk.Face {
				t.Errorf("face mismatch: Go=%d, C=%d", goFijk.Face, cFijk.Face)
			}
			if goFijk.Coord.I != cFijk.Coord.I {
				t.Errorf("I coordinate mismatch: Go=%d, C=%d", goFijk.Coord.I, cFijk.Coord.I)
			}
			if goFijk.Coord.J != cFijk.Coord.J {
				t.Errorf("J coordinate mismatch: Go=%d, C=%d", goFijk.Coord.J, cFijk.Coord.J)
			}
			if goFijk.Coord.K != cFijk.Coord.K {
				t.Errorf("K coordinate mismatch: Go=%d, C=%d", goFijk.Coord.K, cFijk.Coord.K)
			}
		})
	}
}
