//go:build cgo && c2go

package h3

import (
	"testing"
)

func Test_gridDiskDistancesUnsafe_parity(t *testing.T) {
	// Test with various origin cells and k values
	testCases := []struct {
		name   string
		origin H3Index
		k      int32
	}{
		// Valid hexagon cells at different resolutions
		{"res0_k0", 0x8001fffffffffff, 0},
		{"res0_k1", 0x8001fffffffffff, 1},
		{"res0_k2", 0x8001fffffffffff, 2},
		{"res1_k1", 0x8101fffffffffff, 1},
		{"res2_k2", 0x8201fffffffffff, 2},
		{"res5_k3", 0x85283473fffffff, 3},
		{"res7_k5", 0x872830536ffffff, 5},
		{"res9_k10", 0x8928308280fffff, 10},
		{"res10_k15", 0x8a2834700007fff, 15},
		{"res12_k1", 0x8c283470003ffff, 1},
		{"res15_k0", 0x8f283470003dfff, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Calculate max size for allocation
			var maxSize int64
			err := maxGridDiskSize(tc.k, &maxSize)
			if err != E_SUCCESS {
				t.Fatalf("maxGridDiskSize failed: %v", err)
			}

			// Allocate arrays for Go implementation
			goOut := make([]H3Index, maxSize)
			goDistances := make([]int32, maxSize)

			// Allocate arrays for C implementation
			cOut := make([]H3Index, maxSize)
			cDistances := make([]int32, maxSize)

			// Call Go implementation
			goErr := gridDiskDistancesUnsafe(tc.origin, tc.k, goOut, goDistances)

			// Call C implementation
			cErr := gridDiskDistancesUnsafeC(tc.origin, tc.k, cOut, cDistances)

			// Compare errors
			if goErr != cErr {
				t.Errorf("Error mismatch: Go returned %v, C returned %v", goErr, cErr)
				return
			}

			// If successful, compare outputs
			if goErr == E_SUCCESS {
				// Compare cell outputs
				for i := int64(0); i < maxSize; i++ {
					if goOut[i] != cOut[i] {
						t.Errorf("Cell mismatch at index %d: Go=%x, C=%x", i, goOut[i], cOut[i])
					}
				}

				// Compare distance outputs
				for i := int64(0); i < maxSize; i++ {
					if goDistances[i] != cDistances[i] {
						t.Errorf("Distance mismatch at index %d: Go=%d, C=%d", i, goDistances[i], cDistances[i])
					}
				}
			}
		})
	}

	// Test with NULL distances (by passing empty slice)
	t.Run("null_distances", func(t *testing.T) {
		origin := H3Index(0x85283473fffffff)
		k := int32(2)

		var maxSize int64
		err := maxGridDiskSize(k, &maxSize)
		if err != E_SUCCESS {
			t.Fatalf("maxGridDiskSize failed: %v", err)
		}

		goOut := make([]H3Index, maxSize)
		cOut := make([]H3Index, maxSize)

		// Call with empty distances slice
		goErr := gridDiskDistancesUnsafe(origin, k, goOut, nil)
		cErr := gridDiskDistancesUnsafeC(origin, k, cOut, nil)

		if goErr != cErr {
			t.Errorf("Error mismatch with null distances: Go=%v, C=%v", goErr, cErr)
		}

		if goErr == E_SUCCESS {
			for i := int64(0); i < maxSize; i++ {
				if goOut[i] != cOut[i] {
					t.Errorf("Cell mismatch at index %d: Go=%x, C=%x", i, goOut[i], cOut[i])
				}
			}
		}
	})

	// Test negative k (should return E_DOMAIN)
	t.Run("negative_k", func(t *testing.T) {
		origin := H3Index(0x85283473fffffff)
		k := int32(-1)

		out := make([]H3Index, 1)
		distances := make([]int32, 1)

		goErr := gridDiskDistancesUnsafe(origin, k, out, distances)
		cErr := gridDiskDistancesUnsafeC(origin, k, out, distances)

		if goErr != E_DOMAIN || cErr != E_DOMAIN {
			t.Errorf("Expected E_DOMAIN for negative k, got Go=%v, C=%v", goErr, cErr)
		}
	})

	// Test with a pentagon (should return E_PENTAGON)
	t.Run("pentagon", func(t *testing.T) {
		// Pentagon at resolution 0
		pentagon := H3Index(0x8009fffffffffff)
		k := int32(1)

		var maxSize int64
		err := maxGridDiskSize(k, &maxSize)
		if err != E_SUCCESS {
			t.Fatalf("maxGridDiskSize failed: %v", err)
		}

		goOut := make([]H3Index, maxSize)
		goDistances := make([]int32, maxSize)
		cOut := make([]H3Index, maxSize)
		cDistances := make([]int32, maxSize)

		goErr := gridDiskDistancesUnsafe(pentagon, k, goOut, goDistances)
		cErr := gridDiskDistancesUnsafeC(pentagon, k, cOut, cDistances)

		// Both should return E_PENTAGON
		if goErr != E_PENTAGON || cErr != E_PENTAGON {
			t.Errorf("Expected E_PENTAGON for pentagon cell, got Go=%v, C=%v", goErr, cErr)
		}
	})
}
