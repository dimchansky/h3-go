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
		{"k=kAllCellsAtRes15-1", kAllCellsAtRes15 - 1},
		{"k=kAllCellsAtRes15", kAllCellsAtRes15},
		{"k=kAllCellsAtRes15+1", kAllCellsAtRes15 + 1},
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
			if goErr == eSuccess {
				if goOut != cOut {
					t.Errorf("Output mismatch: Go returned %d, C returned %d", goOut, cOut)
				}
			}
		})
	}

	// Test negative k values (should return eDomain)
	negativeTests := []int32{-1, -10, -100}
	for _, k := range negativeTests {
		t.Run("negative k", func(t *testing.T) {
			var goOut, cOut int64

			goErr := maxGridDiskSize(k, &goOut)
			cErr := maxGridDiskSizeC(k, &cOut)

			if goErr != eDomain || cErr != eDomain {
				t.Errorf("Expected eDomain for k=%d, got Go=%v, C=%v", k, goErr, cErr)
			}
		})
	}
}
