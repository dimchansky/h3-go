//go:build cgo

package h3

import (
	"testing"
)

func Test__iterateAllIndexesAtResPartial_parity(t *testing.T) {
	testCases := []struct {
		name      string
		res       int32
		baseCells int32
	}{
		// Test with different numbers of base cells
		{"res0_1cell", 0, 1},
		{"res0_5cells", 0, 5},
		{"res0_10cells", 0, 10},
		{"res0_50cells", 0, 50},
		{"res0_allcells", 0, NUM_BASE_CELLS},

		{"res1_1cell", 1, 1},
		{"res1_5cells", 1, 5},
		{"res1_10cells", 1, 10},
		{"res1_20cells", 1, 20},

		{"res2_1cell", 2, 1},
		{"res2_3cells", 2, 3},
		{"res2_5cells", 2, 5},

		// Test boundary conditions
		{"res0_0cells", 0, 0},
		{"res1_0cells", 1, 0},

		// Test exactly NUM_BASE_CELLS (maximum valid value)
		{"res0_maxcells", 0, NUM_BASE_CELLS},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Calculate expected buffer size
			// Maximum is baseCells * 7^res, but cap for safety
			maxCellsPerBase := _ipow(7, int64(tc.res))
			bufferSize := int64(tc.baseCells) * maxCellsPerBase
			if bufferSize > 100000 {
				bufferSize = 100000
			}
			if bufferSize < 1 {
				bufferSize = 1
			}

			// Collect Go results
			goResults := make([]H3Index, 0, bufferSize)
			_iterateAllIndexesAtResPartial(tc.res, func(h H3Index) {
				goResults = append(goResults, h)
			}, tc.baseCells)

			// Collect C results
			cBuffer := make([]H3Index, bufferSize)
			cCount := iterateAllIndexesAtResPartialC(tc.res, tc.baseCells, cBuffer)
			cResults := cBuffer[:cCount]

			// Compare results
			if len(goResults) != len(cResults) {
				t.Fatalf("Count mismatch: Go=%d, C=%d", len(goResults), len(cResults))
			}

			// For large result sets, spot check
			if len(goResults) > 10000 {
				// Check first 100 and last 100
				for i := 0; i < 100 && i < len(goResults); i++ {
					if goResults[i] != cResults[i] {
						t.Errorf("Index mismatch at position %d: Go=0x%x, C=0x%x",
							i, goResults[i], cResults[i])
					}
				}
				for i := len(goResults) - 100; i < len(goResults); i++ {
					if i >= 0 && goResults[i] != cResults[i] {
						t.Errorf("Index mismatch at position %d: Go=0x%x, C=0x%x",
							i, goResults[i], cResults[i])
					}
				}
			} else {
				// For smaller sets, compare everything
				for i := 0; i < len(goResults); i++ {
					if goResults[i] != cResults[i] {
						t.Errorf("Index mismatch at position %d: Go=0x%x, C=0x%x",
							i, goResults[i], cResults[i])
					}
				}
			}
		})
	}
}
