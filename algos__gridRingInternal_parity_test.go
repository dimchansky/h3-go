//go:build cgo && c2go

package h3

import (
	"testing"
)

func Test__gridRingInternal_parity(t *testing.T) {
	testCases := []struct {
		name   string
		origin h3Index
		k      int32
	}{
		// Test basic cases
		{"k=0", 0x802bfffffffffff, 0},
		{"k=1", 0x802bfffffffffff, 1},
		{"k=2", 0x802bfffffffffff, 2},
		{"k=3", 0x802bfffffffffff, 3},
		{"k=5", 0x802bfffffffffff, 5},

		// Test different resolutions
		{"res0_k1", 0x8009fffffffffff, 1},
		{"res1_k1", 0x8109fffffffffff, 1},
		{"res2_k1", 0x8209fffffffffff, 1},
		{"res3_k1", 0x8309fffffffffff, 1},
		{"res4_k1", 0x8409fffffffffff, 1},
		{"res5_k1", 0x8509fffffffffff, 1},
		{"res6_k1", 0x8609fffffffffff, 1},
		{"res7_k1", 0x8709fffffffffff, 1},
		{"res8_k1", 0x8809fffffffffff, 1},
		{"res9_k1", 0x8909fffffffffff, 1},

		// Test edge cases
		{"res0_k0", 0x8009fffffffffff, 0},
		{"res0_k2", 0x8009fffffffffff, 2},
		{"res1_k0", 0x8109fffffffffff, 0},
		{"res1_k2", 0x8109fffffffffff, 2},

		// Test various cells at different locations
		{"cell1_k1", 0x89283080ddbffff, 1},
		{"cell1_k2", 0x89283080ddbffff, 2},
		{"cell1_k3", 0x89283080ddbffff, 3},
		{"cell2_k1", 0x8928308280fffff, 1},
		{"cell2_k2", 0x8928308280fffff, 2},
		{"cell3_k1", 0x89283480c27ffff, 1},
		{"cell3_k2", 0x89283480c27ffff, 2},

		// Test with cells that may be near pentagons (may fail with ePentagon)
		{"nearPentagon1", 0x821c07fffffffff, 1},
		{"nearPentagon2", 0x821c07fffffffff, 2},
		{"nearPentagon3", 0x8073fffffffffff, 1},
		{"nearPentagon4", 0x8073fffffffffff, 2},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Get the size for allocating arrays
			// For grid ring, the size is 6*k for k>0, and 1 for k=0
			var size int64
			if tc.k == 0 {
				size = 1
			} else {
				size = 6 * int64(tc.k)
			}

			// Allocate output arrays
			goOut := make([]h3Index, size)
			cOut := make([]h3Index, size)

			// Call Go implementation
			goErr := _gridRingInternal(tc.origin, tc.k, goOut)

			// Call C implementation
			cErr := _gridRingInternalC(tc.origin, tc.k, cOut)

			// Compare errors
			if goErr != cErr {
				t.Fatalf("Error mismatch: Go=%v, C=%v", goErr, cErr)
			}

			// If both succeeded, compare outputs
			if goErr == eSuccess {
				// Compare the cell arrays
				for i := int64(0); i < size; i++ {
					if goOut[i] != cOut[i] {
						t.Errorf("Output mismatch at index %d: Go=0x%x, C=0x%x", i, goOut[i], cOut[i])
					}
				}
			}
		})
	}

	// Test with invalid origin
	// Note: _gridRingInternal may not validate the origin index,
	// so we just check that both implementations behave the same
	t.Run("invalid_origin", func(t *testing.T) {
		origin := h3Index(0xFFFFFFFFFFFFFFFF) // Definitely invalid index
		k := int32(1)
		size := int64(6)

		goOut := make([]h3Index, size)
		cOut := make([]h3Index, size)

		goErr := _gridRingInternal(origin, k, goOut)
		cErr := _gridRingInternalC(origin, k, cOut)

		// The errors should match (whether success or failure)
		if goErr != cErr {
			t.Errorf("Error mismatch for invalid origin: Go=%v, C=%v", goErr, cErr)
		}

		// If both succeeded (unlikely but possible if no validation),
		// the outputs should still match
		if goErr == eSuccess {
			for i := int64(0); i < size; i++ {
				if goOut[i] != cOut[i] {
					t.Errorf("Output mismatch at index %d: Go=0x%x, C=0x%x", i, goOut[i], cOut[i])
				}
			}
		}
	})

	// Test larger rings
	t.Run("large_rings", func(t *testing.T) {
		origin := h3Index(0x89283080ddbffff)
		testKValues := []int32{10, 20, 30}

		for _, k := range testKValues {
			size := 6 * int64(k)
			goOut := make([]h3Index, size)
			cOut := make([]h3Index, size)

			goErr := _gridRingInternal(origin, k, goOut)
			cErr := _gridRingInternalC(origin, k, cOut)

			if goErr != cErr {
				t.Errorf("Error mismatch for k=%d: Go=%v, C=%v", k, goErr, cErr)
				continue
			}

			if goErr == eSuccess {
				// Verify outputs match
				for i := int64(0); i < size; i++ {
					if goOut[i] != cOut[i] {
						t.Errorf("Output mismatch at index %d for k=%d: Go=0x%x, C=0x%x", i, k, goOut[i], cOut[i])
						break // Only report first mismatch to avoid spam
					}
				}
			}
		}
	})
}
