//go:build cgo && c2go

package c2go

import (
	"fmt"
	"testing"
)

func Test_cellToLocalIjk_parity(t *testing.T) {
	testCases := []struct {
		name   string
		origin H3Index
		h3     H3Index
	}{
		// Same cell - should return (0,0,0)
		{
			"Same cell - res 0",
			0x8007fffffffffff,
			0x8007fffffffffff,
		},
		{
			"Same cell - res 5",
			0x85283473fffffff,
			0x85283473fffffff,
		},
		// Adjacent cells at different resolutions
		{
			"Adjacent cells - res 1",
			0x8107fffffffffff,
			0x8117fffffffffff,
		},
		{
			"Adjacent cells - res 5",
			0x85283473fffffff,
			0x85283447fffffff,
		},
		// Cells on pentagon base cells
		{
			"Pentagon base cell 4 - res 3",
			0x83047ffffffffff,
			0x83047ffffffffff,
		},
		{
			"Pentagon to neighbor - res 2",
			0x82047ffffffffff,
			0x82027ffffffffff,
		},
		// Different base cells
		{
			"Different base cells - res 4",
			0x84194cfffffffff,
			0x8419a5fffffffff,
		},
		// Edge cases with various resolutions
		{
			"High resolution - res 10",
			0x8a194c0000000000,
			0x8a194c0000000001,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test Go implementation
			var goOut CoordIJK
			goErr := cellToLocalIjk(tc.origin, tc.h3, &goOut)

			// Test C implementation
			var cOut CoordIJK
			cErr := _cellToLocalIjkC(tc.origin, tc.h3, &cOut)

			// Compare errors
			if goErr != cErr {
				t.Errorf("Error mismatch: Go=%v, C=%v", goErr, cErr)
				return
			}

			// If both succeeded, compare results
			if goErr == E_SUCCESS {
				if goOut.I != cOut.I || goOut.J != cOut.J || goOut.K != cOut.K {
					t.Errorf("Result mismatch:\nGo:  {I:%d, J:%d, K:%d}\nC:   {I:%d, J:%d, K:%d}",
						goOut.I, goOut.J, goOut.K,
						cOut.I, cOut.J, cOut.K)
				}
			}
		})
	}
}

// Test error conditions specifically
func Test_cellToLocalIjk_errors_parity(t *testing.T) {
	testCases := []struct {
		name   string
		origin H3Index
		h3     H3Index
	}{
		{
			"Resolution mismatch",
			0x8007fffffffffff, // res 0
			0x8107fffffffffff, // res 1
		},
		{
			"Invalid origin - bad base cell",
			0x807ffffffffffff, // base cell 127, invalid
			0x8007fffffffffff,
		},
		{
			"Invalid h3 - bad base cell",
			0x8007fffffffffff,
			0x807ffffffffffff, // base cell 127, invalid
		},
		{
			"Non-neighbor base cells",
			0x8007fffffffffff, // base cell 0
			0x801ffffffffffff, // base cell 3, not adjacent to 0
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test Go implementation
			var goOut CoordIJK
			goErr := cellToLocalIjk(tc.origin, tc.h3, &goOut)

			// Test C implementation
			var cOut CoordIJK
			cErr := _cellToLocalIjkC(tc.origin, tc.h3, &cOut)

			// Compare errors - both should fail
			if goErr != cErr {
				t.Errorf("Error mismatch: Go=%v, C=%v", goErr, cErr)
			}

			// Both should return error (not E_SUCCESS)
			if goErr == E_SUCCESS || cErr == E_SUCCESS {
				t.Errorf("Expected error but got: Go=%v, C=%v", goErr, cErr)
			}
		})
	}
}

// Test pentagon distortion cases
func Test_cellToLocalIjk_pentagon_parity(t *testing.T) {
	// Pentagon base cells: 4, 14, 24, 38, 49, 58, 63, 72, 83, 97, 107, 117
	pentagonBaseCells := []int32{4, 14, 24, 38, 49, 58, 63, 72, 83, 97, 107, 117}

	for _, baseCell := range pentagonBaseCells {
		t.Run(fmt.Sprintf("Pentagon base cell %d", baseCell), func(t *testing.T) {
			// Create a cell on this pentagon base cell at res 3
			origin := H3Index(0x83000000000000 | (H3Index(baseCell) << H3_BC_OFFSET))
			h3 := origin

			// Test Go implementation
			var goOut CoordIJK
			goErr := cellToLocalIjk(origin, h3, &goOut)

			// Test C implementation
			var cOut CoordIJK
			cErr := _cellToLocalIjkC(origin, h3, &cOut)

			// Compare errors
			if goErr != cErr {
				t.Errorf("Error mismatch for pentagon base cell %d: Go=%v, C=%v", baseCell, goErr, cErr)
				return
			}

			// If both succeeded, compare results
			if goErr == E_SUCCESS {
				if goOut.I != cOut.I || goOut.J != cOut.J || goOut.K != cOut.K {
					t.Errorf("Result mismatch for pentagon base cell %d:\nGo:  {I:%d, J:%d, K:%d}\nC:   {I:%d, J:%d, K:%d}",
						baseCell, goOut.I, goOut.J, goOut.K, cOut.I, cOut.J, cOut.K)
				}
			}
		})
	}
}
