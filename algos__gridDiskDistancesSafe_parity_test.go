//go:build cgo

package h3

import (
	"testing"
)

func Test_gridDiskDistancesSafe_parity(t *testing.T) {
	testCases := []struct {
		name   string
		origin H3Index
		k      int32
	}{
		// Basic test cases
		{"k=0", 0x802bfffffffffff, 0},
		{"k=1", 0x802bfffffffffff, 1},
		{"k=2", 0x802bfffffffffff, 2},
		{"k=3", 0x802bfffffffffff, 3},

		// Test different resolutions
		{"res0_k1", 0x8009fffffffffff, 1},
		{"res1_k1", 0x8109fffffffffff, 1},
		{"res2_k1", 0x8209fffffffffff, 1},
		{"res3_k1", 0x8309fffffffffff, 1},

		// Test edge cases
		{"res0_k0", 0x8009fffffffffff, 0},
		{"res0_k2", 0x8009fffffffffff, 2},

		// Test various cells at different locations
		{"cell1_k1", 0x89283080ddbffff, 1},
		{"cell1_k2", 0x89283080ddbffff, 2},
		{"cell2_k1", 0x8928308280fffff, 1},
		{"cell2_k2", 0x8928308280fffff, 2},

		// Test larger k values
		{"k=5", 0x85283473fffffff, 5},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var maxIdx int64
			err := maxGridDiskSize(tc.k, &maxIdx)
			if err != E_SUCCESS {
				t.Fatalf("maxGridDiskSize failed: %v", err)
			}

			// Allocate buffers for both Go and C implementations
			outGo := make([]H3Index, maxIdx)
			distancesGo := make([]int32, maxIdx)
			outC := make([]H3Index, maxIdx)
			distancesC := make([]int32, maxIdx)

			// Call Go implementation
			errGo := gridDiskDistancesSafe(tc.origin, tc.k, outGo, distancesGo)

			// Call C implementation
			errC := gridDiskDistancesSafeC(tc.origin, tc.k, outC, distancesC)

			// Compare errors
			if errGo != errC {
				t.Errorf("Error mismatch: Go=%v, C=%v", errGo, errC)
			}

			if errGo != E_SUCCESS {
				// If both failed with same error, that's fine
				return
			}

			// Create maps to compare results since order is not guaranteed
			goResults := make(map[H3Index]int32)
			cResults := make(map[H3Index]int32)

			// Collect Go results
			for i := int64(0); i < maxIdx; i++ {
				if outGo[i] != 0 {
					goResults[outGo[i]] = distancesGo[i]
				}
			}

			// Collect C results
			for i := int64(0); i < maxIdx; i++ {
				if outC[i] != 0 {
					cResults[outC[i]] = distancesC[i]
				}
			}

			// Compare results
			if len(goResults) != len(cResults) {
				t.Errorf("Result count mismatch: Go=%d, C=%d", len(goResults), len(cResults))
			}

			for cell, goDist := range goResults {
				cDist, exists := cResults[cell]
				if !exists {
					t.Errorf("Cell %x present in Go results but not in C results", cell)
				} else if goDist != cDist {
					t.Errorf("Distance mismatch for cell %x: Go=%d, C=%d", cell, goDist, cDist)
				}
			}

			for cell := range cResults {
				if _, exists := goResults[cell]; !exists {
					t.Errorf("Cell %x present in C results but not in Go results", cell)
				}
			}
		})
	}
}
