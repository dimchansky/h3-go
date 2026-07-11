//go:build cgo && c2go

package h3

import (
	"testing"
)

func Test_gridDisk_parity(t *testing.T) {
	testCases := []struct {
		name   string
		origin H3Index
		k      int32
	}{
		// Test basic cases
		{"k=0", 0x802bfffffffffff, 0},
		{"k=1", 0x802bfffffffffff, 1},
		{"k=2", 0x802bfffffffffff, 2},
		{"k=3", 0x802bfffffffffff, 3},

		// Test different resolutions
		{"res0_k1", 0x8009fffffffffff, 1},
		{"res1_k1", 0x8109fffffffffff, 1},
		{"res2_k1", 0x8209fffffffffff, 1},
		{"res3_k1", 0x8309fffffffffff, 1},
		{"res4_k1", 0x8409fffffffffff, 1},
		{"res5_k1", 0x8509fffffffffff, 1},

		// Test edge cases
		{"res0_k0", 0x8009fffffffffff, 0},
		{"res0_k2", 0x8009fffffffffff, 2},

		// Test various cells at different locations
		{"cell1_k1", 0x89283080ddbffff, 1},
		{"cell1_k2", 0x89283080ddbffff, 2},
		{"cell2_k1", 0x8928308280fffff, 1},
		{"cell2_k2", 0x8928308280fffff, 2},
		{"cell3_k1", 0x89283480c27ffff, 1},
		{"cell3_k2", 0x89283480c27ffff, 2},

		// Test with cells that may be near pentagons
		{"nearPentagon1", 0x821c07fffffffff, 1},
		{"nearPentagon2", 0x821c07fffffffff, 2},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Get the size for allocating arrays
			var maxSize int64
			err := maxGridDiskSize(tc.k, &maxSize)
			if err != E_SUCCESS {
				t.Fatalf("maxGridDiskSize failed: %v", err)
			}

			// Allocate output arrays
			goOut := make([]H3Index, maxSize)
			cOut := make([]H3Index, maxSize)

			// Call Go implementation
			goErr := gridDisk(tc.origin, tc.k, goOut)

			// Call C implementation
			cErr := gridDiskC(tc.origin, tc.k, cOut)

			// Compare errors
			if goErr != cErr {
				t.Fatalf("Error mismatch: Go=%v, C=%v", goErr, cErr)
			}

			// If both succeeded, compare outputs
			if goErr == E_SUCCESS {
				// Compare the cell arrays
				for i := int64(0); i < maxSize; i++ {
					if goOut[i] != cOut[i] {
						t.Errorf("Output mismatch at index %d: Go=0x%x, C=0x%x", i, goOut[i], cOut[i])
					}
				}
			}
		})
	}

	// Test error cases
	t.Run("error_cases", func(t *testing.T) {
		// Test with invalid k value
		t.Run("negative_k", func(t *testing.T) {
			out := make([]H3Index, 1)

			goErr := gridDisk(0x802bfffffffffff, -1, out)
			cErr := gridDiskC(0x802bfffffffffff, -1, out)

			if goErr != cErr {
				t.Fatalf("Error mismatch for negative k: Go=%v, C=%v", goErr, cErr)
			}
		})

		// Note: Empty output array test removed - this causes undefined behavior in C
		// and would cause a panic in Go, which is equivalent undefined behavior
	})
}
