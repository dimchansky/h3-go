//go:build cgo && c2go

package h3

import (
	"testing"
)

func Test_gridRingUnsafe_parity(t *testing.T) {
	// Test with valid cells at different resolutions
	testCases := []struct {
		name   string
		origin H3Index
		k      int32
	}{
		// Basic tests
		{"res5_k0", 0x85283473fffffff, 0},
		{"res5_k1", 0x85283473fffffff, 1},
		{"res5_k2", 0x85283473fffffff, 2},
		{"res5_k3", 0x85283473fffffff, 3},

		// Different resolutions
		{"res0_k1", 0x8009fffffffffff, 1},
		{"res1_k1", 0x81083ffffffffff, 1},
		{"res3_k2", 0x832834fffffffff, 2},
		{"res7_k1", 0x872834700ffffff, 1},
		{"res10_k1", 0x8a2834700007fff, 1},

		// Edge cases with larger k values
		{"res5_k5", 0x85283473fffffff, 5},
		{"res5_k10", 0x85283473fffffff, 10},
		{"res6_k4", 0x862834707ffffff, 4},

		// Different base cells
		{"base20_res5_k2", 0x8528b4b3fffffff, 2},
		{"base41_res5_k3", 0x8553d4b3fffffff, 3},
		{"base108_res5_k1", 0x85dc94b3fffffff, 1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Calculate the expected size
			var size int64
			if tc.k == 0 {
				size = 1
			} else {
				size = 6 * int64(tc.k)
			}

			// Create output buffers
			goOut := make([]H3Index, size)
			cOut := make([]H3Index, size)

			// Call Go implementation
			goErr := gridRingUnsafe(tc.origin, tc.k, goOut)

			// Call C implementation
			cErr := gridRingUnsafeC(tc.origin, tc.k, cOut)

			// Compare errors
			if goErr != cErr {
				t.Errorf("Error mismatch: Go=%v, C=%v", goErr, cErr)
				return
			}

			// If successful, compare outputs
			if goErr == E_SUCCESS {
				// Create maps to compare unordered sets
				goMap := make(map[H3Index]int)
				cMap := make(map[H3Index]int)

				for _, h := range goOut {
					if h != 0 {
						goMap[h]++
					}
				}

				for _, h := range cOut {
					if h != 0 {
						cMap[h]++
					}
				}

				// Compare maps
				if len(goMap) != len(cMap) {
					t.Errorf("Different number of non-zero cells: Go=%d, C=%d", len(goMap), len(cMap))
				}

				for h, goCount := range goMap {
					cCount, exists := cMap[h]
					if !exists {
						t.Errorf("Cell %x exists in Go output but not in C output", h)
					} else if goCount != cCount {
						t.Errorf("Cell %x count mismatch: Go=%d, C=%d", h, goCount, cCount)
					}
				}

				for h := range cMap {
					if _, exists := goMap[h]; !exists {
						t.Errorf("Cell %x exists in C output but not in Go output", h)
					}
				}
			}
		})
	}

	// Test error cases
	t.Run("negative_k", func(t *testing.T) {
		out := make([]H3Index, 1)
		goErr := gridRingUnsafe(0x85283473fffffff, -1, out)
		cErr := gridRingUnsafeC(0x85283473fffffff, -1, out)

		if goErr != cErr {
			t.Errorf("Error mismatch for negative k: Go=%v, C=%v", goErr, cErr)
		}
		if goErr != E_DOMAIN {
			t.Errorf("Expected E_DOMAIN for negative k, got %v", goErr)
		}
	})

	// Test with pentagon - should return E_PENTAGON
	t.Run("pentagon_origin", func(t *testing.T) {
		// Use a known pentagon index at resolution 5
		pentagonIndex := H3Index(0x85080003fffffff)
		out := make([]H3Index, 6)

		goErr := gridRingUnsafe(pentagonIndex, 1, out)
		cErr := gridRingUnsafeC(pentagonIndex, 1, out)

		if goErr != cErr {
			t.Errorf("Error mismatch for pentagon: Go=%v, C=%v", goErr, cErr)
		}
		if goErr != E_PENTAGON {
			t.Errorf("Expected E_PENTAGON for pentagon origin, got %v", goErr)
		}
	})

	// Test cells near pentagons - may encounter pentagon during traversal
	t.Run("near_pentagon", func(t *testing.T) {
		// Test a cell that's close to a pentagon with larger k values
		testCases := []struct {
			origin H3Index
			k      int32
		}{
			{0x85080107fffffff, 5}, // Near pentagon, k=5
			{0x8508024ffffffff, 3}, // Adjacent to pentagon, k=3
		}

		for _, tc := range testCases {
			var size int64
			if tc.k == 0 {
				size = 1
			} else {
				size = 6 * int64(tc.k)
			}

			goOut := make([]H3Index, size)
			cOut := make([]H3Index, size)

			goErr := gridRingUnsafe(tc.origin, tc.k, goOut)
			cErr := gridRingUnsafeC(tc.origin, tc.k, cOut)

			// Just compare errors - both should detect pentagon issues
			if goErr != cErr {
				t.Errorf("Error mismatch near pentagon (origin=%x, k=%d): Go=%v, C=%v",
					tc.origin, tc.k, goErr, cErr)
			}
		}
	})

	// Test invalid cell
	t.Run("invalid_cell", func(t *testing.T) {
		out := make([]H3Index, 6)
		goErr := gridRingUnsafe(0xffffffffffffffff, 1, out)
		cErr := gridRingUnsafeC(0xffffffffffffffff, 1, out)

		// Both should return an error for invalid cell
		if goErr == E_SUCCESS || cErr == E_SUCCESS {
			t.Errorf("Expected error for invalid cell: Go=%v, C=%v", goErr, cErr)
		}
	})
}
