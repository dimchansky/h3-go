//go:build cgo

package h3

import (
	"testing"
)

func Test_gridDiskDistances_parity(t *testing.T) {
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

			// Test with distances array
			t.Run("with_distances", func(t *testing.T) {
				// Allocate output arrays
				goOut := make([]H3Index, maxSize)
				goDistances := make([]int32, maxSize)
				cOut := make([]H3Index, maxSize)
				cDistances := make([]int32, maxSize)

				// Call Go implementation
				goErr := gridDiskDistances(tc.origin, tc.k, goOut, goDistances)

				// Call C implementation
				cErr := gridDiskDistancesC(tc.origin, tc.k, cOut, cDistances)

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

					// Compare the distance arrays
					for i := int64(0); i < maxSize; i++ {
						if goDistances[i] != cDistances[i] {
							t.Errorf("Distance mismatch at index %d: Go=%d, C=%d", i, goDistances[i], cDistances[i])
						}
					}
				}
			})

			// Test without distances array (nil)
			t.Run("without_distances", func(t *testing.T) {
				// Allocate output arrays
				goOut := make([]H3Index, maxSize)
				cOut := make([]H3Index, maxSize)

				// Call Go implementation with nil distances
				goErr := gridDiskDistances(tc.origin, tc.k, goOut, nil)

				// Call C implementation with nil distances
				cErr := gridDiskDistancesC(tc.origin, tc.k, cOut, nil)

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
		})
	}

	// Test error cases
	t.Run("error_cases", func(t *testing.T) {
		// Test with invalid k value
		t.Run("negative_k", func(t *testing.T) {
			out := make([]H3Index, 1)
			distances := make([]int32, 1)

			goErr := gridDiskDistances(0x802bfffffffffff, -1, out, distances)
			cErr := gridDiskDistancesC(0x802bfffffffffff, -1, out, distances)

			if goErr != cErr {
				t.Fatalf("Error mismatch for negative k: Go=%v, C=%v", goErr, cErr)
			}
		})
	})
}
