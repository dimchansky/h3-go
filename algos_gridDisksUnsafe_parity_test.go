//go:build cgo && c2go

package h3

import (
	"testing"
)

func Test_gridDisksUnsafe_parity(t *testing.T) {
	// Test with various sets of input cells and k values
	testCases := []struct {
		name  string
		h3Set []H3Index
		k     int32
	}{
		// Single cell cases
		{"single_res0_k0", []H3Index{0x8001fffffffffff}, 0},
		{"single_res0_k1", []H3Index{0x8001fffffffffff}, 1},
		{"single_res0_k2", []H3Index{0x8001fffffffffff}, 2},
		{"single_res5_k3", []H3Index{0x85283473fffffff}, 3},
		{"single_res7_k5", []H3Index{0x872830536ffffff}, 5},

		// Multiple cell cases
		{"multi_res0_k1", []H3Index{0x8001fffffffffff, 0x8003fffffffffff}, 1},
		{"multi_res5_k2", []H3Index{0x85283473fffffff, 0x85283077fffffff}, 2},
		{"multi_various_k3", []H3Index{0x8001fffffffffff, 0x85283473fffffff, 0x872830536ffffff}, 3},

		// Large set with small k
		{"large_set_k1", []H3Index{
			0x8001fffffffffff, 0x8003fffffffffff, 0x8005fffffffffff,
			0x8007fffffffffff,
		}, 1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Calculate max size for allocation
			var segmentSize int64
			err := maxGridDiskSize(tc.k, &segmentSize)
			if err != E_SUCCESS {
				t.Fatalf("maxGridDiskSize failed: %v", err)
			}

			totalSize := segmentSize * int64(len(tc.h3Set))

			// Allocate arrays for Go implementation
			goOut := make([]H3Index, totalSize)

			// Allocate arrays for C implementation
			cOut := make([]H3Index, totalSize)

			// Call Go implementation
			goErr := gridDisksUnsafe(tc.h3Set, tc.k, goOut)

			// Call C implementation
			cErr := gridDisksUnsafeC(tc.h3Set, tc.k, cOut)

			// Compare errors
			if goErr != cErr {
				t.Errorf("Error mismatch: Go returned %v, C returned %v", goErr, cErr)
				return
			}

			// If successful, compare outputs
			if goErr == E_SUCCESS {
				// Compare cell outputs
				for i := int64(0); i < totalSize; i++ {
					if goOut[i] != cOut[i] {
						t.Errorf("Cell mismatch at index %d: Go=%x, C=%x", i, goOut[i], cOut[i])
					}
				}
			}
		})
	}

	// Test empty set (should return E_FAILED)
	t.Run("empty_set", func(t *testing.T) {
		h3Set := []H3Index{}
		k := int32(1)

		goErr := gridDisksUnsafe(h3Set, k, []H3Index{})
		cErr := gridDisksUnsafeC(h3Set, k, []H3Index{})

		if goErr != E_FAILED || cErr != E_FAILED {
			t.Errorf("Expected E_FAILED for empty set, got Go=%v, C=%v", goErr, cErr)
		}
	})

	// Test negative k (should return E_DOMAIN)
	t.Run("negative_k", func(t *testing.T) {
		h3Set := []H3Index{0x85283473fffffff}
		k := int32(-1)

		out := make([]H3Index, 1)

		goErr := gridDisksUnsafe(h3Set, k, out)
		cErr := gridDisksUnsafeC(h3Set, k, out)

		if goErr != E_DOMAIN || cErr != E_DOMAIN {
			t.Errorf("Expected E_DOMAIN for negative k, got Go=%v, C=%v", goErr, cErr)
		}
	})

	// Test with a pentagon (should return E_PENTAGON)
	t.Run("pentagon", func(t *testing.T) {
		// Pentagon at resolution 0
		pentagon := []H3Index{0x8009fffffffffff}
		k := int32(1)

		var segmentSize int64
		err := maxGridDiskSize(k, &segmentSize)
		if err != E_SUCCESS {
			t.Fatalf("maxGridDiskSize failed: %v", err)
		}

		goOut := make([]H3Index, segmentSize)
		cOut := make([]H3Index, segmentSize)

		goErr := gridDisksUnsafe(pentagon, k, goOut)
		cErr := gridDisksUnsafeC(pentagon, k, cOut)

		// Both should return E_PENTAGON
		if goErr != E_PENTAGON || cErr != E_PENTAGON {
			t.Errorf("Expected E_PENTAGON for pentagon cell, got Go=%v, C=%v", goErr, cErr)
		}
	})

	// Test with multiple cells where one is a pentagon
	t.Run("mixed_with_pentagon", func(t *testing.T) {
		// Mix pentagon with regular hexagon
		h3Set := []H3Index{0x8001fffffffffff, 0x8009fffffffffff} // hexagon, pentagon
		k := int32(1)

		var segmentSize int64
		err := maxGridDiskSize(k, &segmentSize)
		if err != E_SUCCESS {
			t.Fatalf("maxGridDiskSize failed: %v", err)
		}

		totalSize := segmentSize * int64(len(h3Set))
		goOut := make([]H3Index, totalSize)
		cOut := make([]H3Index, totalSize)

		goErr := gridDisksUnsafe(h3Set, k, goOut)
		cErr := gridDisksUnsafeC(h3Set, k, cOut)

		// Both should return E_PENTAGON when pentagon is encountered
		if goErr != E_PENTAGON || cErr != E_PENTAGON {
			t.Errorf("Expected E_PENTAGON when pentagon in set, got Go=%v, C=%v", goErr, cErr)
		}
	})

	// Test with k=0 (should return only origin cells)
	t.Run("k_zero", func(t *testing.T) {
		h3Set := []H3Index{0x85283473fffffff, 0x85283077fffffff}
		k := int32(0)

		// Each cell produces exactly 1 output (itself)
		goOut := make([]H3Index, len(h3Set))
		cOut := make([]H3Index, len(h3Set))

		goErr := gridDisksUnsafe(h3Set, k, goOut)
		cErr := gridDisksUnsafeC(h3Set, k, cOut)

		if goErr != cErr {
			t.Errorf("Error mismatch for k=0: Go=%v, C=%v", goErr, cErr)
		}

		if goErr == E_SUCCESS {
			for i := 0; i < len(h3Set); i++ {
				if goOut[i] != cOut[i] {
					t.Errorf("Cell mismatch at index %d for k=0: Go=%x, C=%x", i, goOut[i], cOut[i])
				}
			}
		}
	})
}
