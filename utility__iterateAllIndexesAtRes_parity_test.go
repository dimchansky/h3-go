//go:build cgo && c2go

package h3

import (
	"testing"
)

func Test__iterateAllIndexesAtRes_parity(t *testing.T) {
	testCases := []struct {
		name string
		res  int32
	}{
		{"res0", 0},
		{"res1", 1},
		{"res2", 2},
		// Note: Higher resolutions produce millions of cells, so we keep tests small
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Calculate expected total cells at resolution
			totalCells, err := getNumCells(tc.res)
			if err != eSuccess {
				t.Fatalf("getNumCells failed: %v", err)
			}

			// For safety in tests, cap buffer size
			bufferSize := totalCells
			if bufferSize > 1000000 {
				t.Skip("Skipping test with more than 1M cells to avoid memory issues")
			}

			// Collect Go results
			goResults := make([]h3Index, 0, bufferSize)
			_iterateAllIndexesAtRes(tc.res, func(h h3Index) {
				goResults = append(goResults, h)
			})

			// Collect C results
			cBuffer := make([]h3Index, bufferSize)
			cCount := iterateAllIndexesAtResC(tc.res, cBuffer)
			cResults := cBuffer[:cCount]

			// Compare results
			if len(goResults) != len(cResults) {
				t.Fatalf("Count mismatch: Go=%d, C=%d", len(goResults), len(cResults))
			}

			// For large result sets, spot check rather than compare all
			if len(goResults) > 10000 {
				// Check first 100, last 100, and some samples in between
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
				// Check some samples in the middle
				step := len(goResults) / 100
				for i := step; i < len(goResults); i += step {
					if goResults[i] != cResults[i] {
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
