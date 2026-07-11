//go:build cgo && c2go

package h3

import (
	"fmt"
	"testing"
)

func TestH3ToFaceIjkParity(t *testing.T) {
	testCases := []h3Index{
		// Base resolution cells
		0x8001fffffffffff, // res 0, base cell 1
		0x8007fffffffffff, // res 0, base cell 7 (pentagon)
		0x800dfffffffffff, // res 0, base cell 13 (pentagon)

		// Higher resolution cells
		0x81283ffffffffff, // res 1, regular hexagon
		0x8228bffffffffff, // res 2, regular hexagon
		0x832830fffffffff, // res 3, regular hexagon

		// Pentagon cases
		0x81703ffffffffff, // res 1, pentagon base cell
		0x82734bfffffffff, // res 2, pentagon base cell

		// Edge cases around pentagon boundaries
		0x8427cffffffffff, // res 4, near pentagon
		0x85283473fffffff, // res 5, regular case

		// High resolution case
		0x8a283082803ffff, // res 10
		0x8b283082800bfff, // res 11
	}

	for i, h3 := range testCases {
		t.Run(fmt.Sprintf("case_%d_0x%x", i, uint64(h3)), func(t *testing.T) {
			// Test Go implementation
			var goFijk faceIJK
			goErr := _h3ToFaceIjk(h3, &goFijk)

			// Test C implementation
			var cFijk faceIJK
			cErr := _h3ToFaceIjkC(h3, &cFijk)

			// Compare errors
			if goErr != h3Error(cErr) {
				t.Errorf("Error mismatch: Go=%d, C=%d", goErr, cErr)
			}

			// If no error, compare results
			if goErr == eSuccess && cErr == 0 {
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
			}
		})
	}
}
