//go:build cgo

package c2go

import (
	"testing"
)

func TestAdjustOverageClassIIParity(t *testing.T) {
	testCases := []struct {
		name         string
		fijk         FaceIJK
		res          int
		pentLeading4 bool
		substrate    bool
	}{
		{
			name: "no overage case",
			fijk: FaceIJK{Face: 0, Coord: CoordIJK{I: 1, J: 1, K: 0}},
			res:  2, pentLeading4: false, substrate: false,
		},
		{
			name: "face edge case",
			fijk: FaceIJK{Face: 0, Coord: CoordIJK{I: 7, J: 7, K: 0}},
			res:  2, pentLeading4: false, substrate: true,
		},
		{
			name: "overage jk quadrant",
			fijk: FaceIJK{Face: 0, Coord: CoordIJK{I: 1, J: 8, K: 8}},
			res:  2, pentLeading4: false, substrate: false,
		},
		{
			name: "overage ik quadrant",
			fijk: FaceIJK{Face: 0, Coord: CoordIJK{I: 8, J: 0, K: 8}},
			res:  2, pentLeading4: false, substrate: false,
		},
		{
			name: "overage ij quadrant",
			fijk: FaceIJK{Face: 0, Coord: CoordIJK{I: 8, J: 8, K: 0}},
			res:  2, pentLeading4: false, substrate: false,
		},
		{
			name: "pentagon leading4 case",
			fijk: FaceIJK{Face: 4, Coord: CoordIJK{I: 10, J: 0, K: 8}},
			res:  2, pentLeading4: true, substrate: false,
		},
		{
			name: "substrate grid case",
			fijk: FaceIJK{Face: 1, Coord: CoordIJK{I: 15, J: 15, K: 15}},
			res:  2, pentLeading4: false, substrate: true,
		},
		{
			name: "higher resolution",
			fijk: FaceIJK{Face: 5, Coord: CoordIJK{I: 50, J: 30, K: 30}},
			res:  4, pentLeading4: false, substrate: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test Go implementation
			goFijk := tc.fijk
			pentLeading4Int := 0
			if tc.pentLeading4 {
				pentLeading4Int = 1
			}
			substrateInt := 0
			if tc.substrate {
				substrateInt = 1
			}

			goResult := _adjustOverageClassII(&goFijk, tc.res, tc.pentLeading4, tc.substrate)

			// Test C implementation
			cFijk := tc.fijk
			cResult := _adjustOverageClassIIC(&cFijk, tc.res, pentLeading4Int, substrateInt)

			// Compare results
			if goResult != cResult {
				t.Errorf("Overage mismatch: Go=%d, C=%d", goResult, cResult)
			}
			if goFijk.Face != cFijk.Face {
				t.Errorf("Face mismatch: Go=%d, C=%d", goFijk.Face, cFijk.Face)
			}
			if goFijk.Coord.I != cFijk.Coord.I {
				t.Errorf("Coord.I mismatch: Go=%d, C=%d", goFijk.Coord.I, cFijk.Coord.I)
			}
			if goFijk.Coord.J != cFijk.Coord.J {
				t.Errorf("Coord.J mismatch: Go=%d, C=%d", goFijk.Coord.J, cFijk.Coord.J)
			}
			if goFijk.Coord.K != cFijk.Coord.K {
				t.Errorf("Coord.K mismatch: Go=%d, C=%d", goFijk.Coord.K, cFijk.Coord.K)
			}
		})
	}
}
