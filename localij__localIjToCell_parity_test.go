//go:build cgo && c2go

package h3

import (
	"fmt"
	"testing"
)

func Test_localIjToCell_parity(t *testing.T) {
	testCases := []struct {
		name   string
		origin h3Index
		ij     CoordIJ
		mode   uint32
	}{
		// Origin cell - should return the origin itself
		{
			"Origin cell - res 0",
			0x8007fffffffffff,
			CoordIJ{I: 0, J: 0},
			0,
		},
		{
			"Origin cell - res 5",
			0x85283473fffffff,
			CoordIJ{I: 0, J: 0},
			0,
		},
		// Adjacent cells at different resolutions
		{
			"Adjacent I direction - res 1",
			0x8107fffffffffff,
			CoordIJ{I: 1, J: 0},
			0,
		},
		{
			"Adjacent J direction - res 1",
			0x8107fffffffffff,
			CoordIJ{I: 0, J: 1},
			0,
		},
		{
			"Adjacent diagonal - res 2",
			0x8207ffffffffff,
			CoordIJ{I: 1, J: 1},
			0,
		},
		// Negative coordinates
		{
			"Negative I - res 3",
			0x8307ffffffffff,
			CoordIJ{I: -1, J: 0},
			0,
		},
		{
			"Negative J - res 3",
			0x8307ffffffffff,
			CoordIJ{I: 0, J: -1},
			0,
		},
		{
			"Negative both - res 3",
			0x8307ffffffffff,
			CoordIJ{I: -1, J: -1},
			0,
		},
		// Pentagon base cells
		{
			"Pentagon base cell 4 - res 3",
			0x83047ffffffffff,
			CoordIJ{I: 0, J: 0},
			0,
		},
		{
			"Pentagon with offset - res 2",
			0x82047ffffffffff,
			CoordIJ{I: 1, J: 0},
			0,
		},
		// Larger offsets
		{
			"Large positive offset - res 4",
			0x84194a5bfffffff,
			CoordIJ{I: 3, J: 2},
			0,
		},
		{
			"Large negative offset - res 4",
			0x84194a5bfffffff,
			CoordIJ{I: -2, J: -3},
			0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test Go implementation
			var goOut h3Index
			goErr := localIjToCell(tc.origin, &tc.ij, tc.mode, &goOut)

			// Test C implementation
			var cOut h3Index
			cErr := _localIjToCellC(tc.origin, &tc.ij, tc.mode, &cOut)

			// Compare errors first
			if goErr != cErr {
				t.Errorf("Error mismatch: Go=%v, C=%v", goErr, cErr)
			}

			// If both succeeded, compare results
			if goErr == eSuccess {
				if goOut != cOut {
					t.Errorf("Result mismatch:\nGo:  %x\nC:   %x", goOut, cOut)
				}
			}
		})
	}
}

// Test error conditions specifically
func Test_localIjToCell_errors_parity(t *testing.T) {
	testCases := []struct {
		name   string
		origin h3Index
		ij     CoordIJ
		mode   uint32
	}{
		{
			"Invalid mode",
			0x8007fffffffffff,
			CoordIJ{I: 0, J: 0},
			1, // mode must be 0
		},
		{
			"Invalid origin - bad base cell",
			0x800fa7ffffffffff, // base cell 125, invalid
			CoordIJ{I: 0, J: 0},
			0,
		},
		{
			"Out of range coordinates",
			0x8007fffffffffff,
			CoordIJ{I: 1000, J: 1000}, // very large coordinates
			0,
		},
		{
			"Coordinates causing overflow",
			0x8007fffffffffff,
			CoordIJ{I: 0x7FFFFFFF, J: 0x7FFFFFFF}, // max int32 values
			0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test Go implementation
			var goOut h3Index
			goErr := localIjToCell(tc.origin, &tc.ij, tc.mode, &goOut)

			// Test C implementation
			var cOut h3Index
			cErr := _localIjToCellC(tc.origin, &tc.ij, tc.mode, &cOut)

			// Compare errors - both should fail
			if goErr != cErr {
				t.Errorf("Error mismatch: Go=%v, C=%v", goErr, cErr)
			}

			// Both should return error (not eSuccess)
			if goErr == eSuccess || cErr == eSuccess {
				t.Errorf("Expected error but got: Go=%v, C=%v", goErr, cErr)
			}
		})
	}
}

// Test pentagon cases specifically
func Test_localIjToCell_pentagon_parity(t *testing.T) {
	// Pentagon base cell indices
	pentagonBaseCells := []h3Index{
		0x8007fffffffffff, // 4
		0x800ffffffffffff, // 14
		0x8017fffffffffff, // 24
		0x8027fffffffffff, // 38
		0x802ffffffffffff, // 49
	}

	for _, origin := range pentagonBaseCells {
		t.Run(fmt.Sprintf("Pentagon base cell origin %x", origin), func(t *testing.T) {
			// Test with origin coordinates (should always work)
			ij := CoordIJ{I: 0, J: 0}
			var goOut h3Index
			goErr := localIjToCell(origin, &ij, 0, &goOut)

			var cOut h3Index
			cErr := _localIjToCellC(origin, &ij, 0, &cOut)

			// Compare errors first
			if goErr != cErr {
				t.Errorf("Error mismatch: Go=%v, C=%v", goErr, cErr)
			}

			// If both succeeded, compare results
			if goErr == eSuccess {
				if goOut != cOut {
					t.Errorf("Result mismatch:\nGo:  %x\nC:   %x", goOut, cOut)
				}

				// For origin coordinates, result should be the origin itself
				if goOut != origin {
					t.Errorf("Expected origin %x for (0,0) coordinates, got %x", origin, goOut)
				}
			}
		})
	}
}
