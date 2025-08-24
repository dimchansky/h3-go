//go:build cgo && c2go

package c2go

import (
	"fmt"
	"testing"
)

func Test_cellToLocalIj_parity(t *testing.T) {
	testCases := []struct {
		name   string
		origin H3Index
		index  H3Index
		mode   uint32
	}{
		// Same cell - should return (0,0)
		{
			"Same cell - res 0",
			0x8007fffffffffff,
			0x8007fffffffffff,
			0,
		},
		{
			"Same cell - res 5",
			0x85283473fffffff,
			0x85283473fffffff,
			0,
		},
		// Adjacent cells at different resolutions
		{
			"Adjacent cells - res 1",
			0x8107fffffffffff,
			0x8117fffffffffff,
			0,
		},
		{
			"Adjacent cells - res 5",
			0x85283473fffffff,
			0x85283447fffffff,
			0,
		},
		// Cells on pentagon base cells
		{
			"Pentagon base cell 4 - res 3",
			0x83047ffffffffff,
			0x83047ffffffffff,
			0,
		},
		{
			"Pentagon to neighbor - res 2",
			0x82047ffffffffff,
			0x82027ffffffffff,
			0,
		},
		// Different base cells
		{
			"Cross-base-cell - res 2",
			0x8207ffffffffff,
			0x8217ffffffffff,
			0,
		},
		{
			"Cross-base-cell - res 4",
			0x84194a5bfffffff,
			0x84194a67fffffff,
			0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test Go implementation
			var goOut CoordIJ
			goErr := cellToLocalIj(tc.origin, tc.index, tc.mode, &goOut)

			// Test C implementation
			var cOut CoordIJ
			cErr := _cellToLocalIjC(tc.origin, tc.index, tc.mode, &cOut)

			// Compare errors first
			if goErr != cErr {
				t.Errorf("Error mismatch: Go=%v, C=%v", goErr, cErr)
			}

			// If both succeeded, compare results
			if goErr == E_SUCCESS {
				if goOut.I != cOut.I || goOut.J != cOut.J {
					t.Errorf("Result mismatch:\nGo:  {I:%d, J:%d}\nC:   {I:%d, J:%d}",
						goOut.I, goOut.J,
						cOut.I, cOut.J)
				}
			}
		})
	}
}

// Test error conditions specifically
func Test_cellToLocalIj_errors_parity(t *testing.T) {
	testCases := []struct {
		name   string
		origin H3Index
		index  H3Index
		mode   uint32
	}{
		{
			"Invalid mode",
			0x8007fffffffffff,
			0x8007fffffffffff,
			1, // mode must be 0
		},
		{
			"Resolution mismatch",
			0x8007fffffffffff, // res 0
			0x8107fffffffffff, // res 1
			0,
		},
		{
			"Invalid origin - bad base cell",
			0x800fa7ffffffffff, // base cell 125, invalid
			0x8007fffffffffff,
			0,
		},
		{
			"Invalid index - bad base cell",
			0x8007fffffffffff,
			0x800fa7ffffffffff, // base cell 125, invalid
			0,
		},
		{
			"Non-neighbor base cells",
			0x8007fffffffffff, // base cell 0
			0x801ffffffffffff, // base cell 3, not adjacent to 0
			0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test Go implementation
			var goOut CoordIJ
			goErr := cellToLocalIj(tc.origin, tc.index, tc.mode, &goOut)

			// Test C implementation
			var cOut CoordIJ
			cErr := _cellToLocalIjC(tc.origin, tc.index, tc.mode, &cOut)

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

// Test pentagon cases specifically
func Test_cellToLocalIj_pentagon_parity(t *testing.T) {
	// Pentagon base cell indices (24 different base cells)
	pentagonBaseCells := []H3Index{
		0x8007fffffffffff, // 4
		0x800ffffffffffff, // 14
		0x8017fffffffffff, // 24
		0x8027fffffffffff, // 38
		0x802ffffffffffff, // 49
	}

	for _, origin := range pentagonBaseCells {
		t.Run(fmt.Sprintf("Pentagon base cell origin %x", origin), func(t *testing.T) {
			// Test with same cell (should always work)
			var goOut CoordIJ
			goErr := cellToLocalIj(origin, origin, 0, &goOut)

			var cOut CoordIJ
			cErr := _cellToLocalIjC(origin, origin, 0, &cOut)

			// Compare errors first
			if goErr != cErr {
				t.Errorf("Error mismatch: Go=%v, C=%v", goErr, cErr)
			}

			// If both succeeded, compare results
			if goErr == E_SUCCESS {
				if goOut.I != cOut.I || goOut.J != cOut.J {
					t.Errorf("Result mismatch:\nGo:  {I:%d, J:%d}\nC:   {I:%d, J:%d}",
						goOut.I, goOut.J,
						cOut.I, cOut.J)
				}

				// For same cell, coordinates should be (0,0)
				if goOut.I != 0 || goOut.J != 0 {
					t.Errorf("Expected (0,0) for same cell, got (%d,%d)", goOut.I, goOut.J)
				}
			}
		})
	}
}
