//go:build cgo && c2go

package c2go

import (
	"testing"
)

func Test_gridDistance_parity(t *testing.T) {
	testCases := []struct {
		name   string
		origin H3Index
		index  H3Index
	}{
		// Same cell - should return distance 0
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
			var goResult int64
			var cResult int64

			// Test Go implementation
			goErr := gridDistance(tc.origin, tc.index, &goResult)

			// Test C implementation
			cErr := _gridDistanceC(tc.origin, tc.index, &cResult)

			// Compare errors
			if goErr != cErr {
				t.Errorf("Error mismatch: Go=%d, C=%d", goErr, cErr)
				return
			}

			// If both implementations succeeded, compare results
			if goErr == 0 && cErr == 0 {
				if goResult != cResult {
					t.Errorf("Distance mismatch: Go=%d, C=%d", goResult, cResult)
				}
			}

			// Log successful cases for verification
			if goErr == 0 {
				t.Logf("Distance from %016x to %016x: %d", tc.origin, tc.index, goResult)
			} else {
				t.Logf("Distance calculation failed with error: %d", goErr)
			}
		})
	}
}

func Test_gridDistance_error_cases_parity(t *testing.T) {
	errorCases := []struct {
		name   string
		origin H3Index
		index  H3Index
	}{
		{
			"Different resolutions",
			0x8107fffffffffff, // res 1
			0x85283473fffffff, // res 5
		},
		{
			"Invalid origin",
			0x0000000000000000, // Invalid H3Index
			0x85283473fffffff,
		},
		{
			"Invalid index",
			0x85283473fffffff,
			0x0000000000000000, // Invalid H3Index
		},
	}

	for _, tc := range errorCases {
		t.Run(tc.name, func(t *testing.T) {
			var goResult int64
			var cResult int64

			// Test Go implementation
			goErr := gridDistance(tc.origin, tc.index, &goResult)

			// Test C implementation
			cErr := _gridDistanceC(tc.origin, tc.index, &cResult)

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
