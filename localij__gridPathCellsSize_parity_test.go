//go:build cgo && c2go

package h3

import (
	"testing"
)

func Test_gridPathCellsSize_parity(t *testing.T) {
	testCases := []struct {
		name  string
		start h3Index
		end   h3Index
	}{
		// Same cell - should return size 1 (start == end)
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
			0x85283477fffffff,
		},
		// Further apart cells
		{
			"Distant cells - res 3",
			0x83754e64fffffff,
			0x83754e65fffffff,
		},
		{
			"Distant cells - res 7",
			0x87283470fffffff,
			0x87283471fffffff,
		},
		// Test with pentagon cells (known to be more complex)
		{
			"Pentagon origin - res 5",
			0x851c0003fffffff, // Pentagon base cell
			0x851c0007fffffff,
		},
		// Cross different base cells
		{
			"Different base cells - res 4",
			0x8428309ffffffff,
			0x8428301ffffffff,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var goSize int64
			var cSize int64

			// Test Go implementation
			goErr := gridPathCellsSize(tc.start, tc.end, &goSize)

			// Test C implementation
			cErr := _gridPathCellsSizeC(tc.start, tc.end, &cSize)

			// Compare errors
			if goErr != cErr {
				t.Errorf("Error mismatch: Go=%d, C=%d", goErr, cErr)
				return
			}

			// If both implementations succeeded, compare results
			if goErr == 0 && cErr == 0 {
				if goSize != cSize {
					t.Errorf("Size mismatch: Go=%d, C=%d", goSize, cSize)
				}

				// Verify the expected relationship: size = distance + 1
				var distance int64
				distErr := gridDistance(tc.start, tc.end, &distance)
				if distErr == 0 && goSize != distance+1 {
					t.Errorf("Size not equal to distance+1: size=%d, distance=%d", goSize, distance)
				}
			}

			// Log successful cases for verification
			if goErr == 0 {
				t.Logf("Path size from %016x to %016x: %d", tc.start, tc.end, goSize)
			} else {
				t.Logf("Path size calculation failed with error: %d", goErr)
			}
		})
	}
}

func Test_gridPathCellsSize_error_cases_parity(t *testing.T) {
	errorCases := []struct {
		name  string
		start h3Index
		end   h3Index
	}{
		{
			"Different resolutions",
			0x8107fffffffffff, // res 1
			0x85283473fffffff, // res 5
		},
		{
			"Invalid start",
			0x0000000000000000, // Invalid h3Index
			0x85283473fffffff,
		},
		{
			"Invalid end",
			0x85283473fffffff,
			0x0000000000000000, // Invalid h3Index
		},
	}

	for _, tc := range errorCases {
		t.Run(tc.name, func(t *testing.T) {
			var goSize int64
			var cSize int64

			// Test Go implementation
			goErr := gridPathCellsSize(tc.start, tc.end, &goSize)

			// Test C implementation
			cErr := _gridPathCellsSizeC(tc.start, tc.end, &cSize)

			// Both should fail with the same error
			if goErr != cErr {
				t.Errorf("Error code mismatch: Go=%d, C=%d", goErr, cErr)
			}

			// Both should return an error (non-zero)
			if goErr == 0 || cErr == 0 {
				t.Errorf("Expected error but got success: Go err=%d, C err=%d", goErr, cErr)
			}

			t.Logf("Expected error case: Go err=%d, C err=%d", goErr, cErr)
		})
	}
}
