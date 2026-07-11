//go:build cgo && c2go

package h3

import "testing"

func Test_maxGridDiskSize_parity(t *testing.T) {
	testCases := []struct {
		name string
		k    int32
	}{
		{"k=0", 0},
		{"k=1", 1},
		{"k=2", 2},
		{"k=5", 5},
		{"k=10", 10},
		{"k=100", 100},
		{"k=1000", 1000},
		{"k=10000", 10000},
		{"k=MAX_INT32", 2147483647}, // Test overflow handling
		{"k=K_ALL_CELLS_AT_RES_15-1", K_ALL_CELLS_AT_RES_15 - 1},
		{"k=K_ALL_CELLS_AT_RES_15", K_ALL_CELLS_AT_RES_15},
		{"k=K_ALL_CELLS_AT_RES_15+1", K_ALL_CELLS_AT_RES_15 + 1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var goOut, cOut int64

			// Call Go implementation
			goErr := maxGridDiskSize(tc.k, &goOut)

			// Call C implementation
			cErr := maxGridDiskSizeC(tc.k, &cOut)

			// Compare errors
			if goErr != cErr {
				t.Errorf("Error mismatch: Go returned %v, C returned %v", goErr, cErr)
			}

			// Compare outputs if no error
			if goErr == E_SUCCESS {
				if goOut != cOut {
					t.Errorf("Output mismatch: Go returned %d, C returned %d", goOut, cOut)
				}
			}
		})
	}

	// Test negative k values (should return E_DOMAIN)
	negativeTests := []int32{-1, -10, -100}
	for _, k := range negativeTests {
		t.Run("negative k", func(t *testing.T) {
			var goOut, cOut int64

			goErr := maxGridDiskSize(k, &goOut)
			cErr := maxGridDiskSizeC(k, &cOut)

			if goErr != E_DOMAIN || cErr != E_DOMAIN {
				t.Errorf("Expected E_DOMAIN for k=%d, got Go=%v, C=%v", k, goErr, cErr)
			}
		})
	}
}
