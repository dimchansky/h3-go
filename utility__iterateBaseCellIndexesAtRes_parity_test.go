//go:build cgo

package h3

import (
	"testing"
)

func Test__iterateBaseCellIndexesAtRes_parity(t *testing.T) {
	testCases := []struct {
		name     string
		res      int32
		baseCell int32
	}{
		// Test various resolutions for different base cells
		{"res0_base0", 0, 0},
		{"res0_base4", 0, 4},   // Pentagon
		{"res0_base14", 0, 14}, // Pentagon
		{"res0_base50", 0, 50},
		{"res0_base121", 0, 121}, // Last base cell

		{"res1_base0", 1, 0},
		{"res1_base4", 1, 4},   // Pentagon
		{"res1_base14", 1, 14}, // Pentagon
		{"res1_base50", 1, 50},
		{"res1_base121", 1, 121}, // Last base cell

		{"res2_base0", 2, 0},
		{"res2_base4", 2, 4},   // Pentagon
		{"res2_base14", 2, 14}, // Pentagon
		{"res2_base50", 2, 50},

		{"res3_base0", 3, 0},
		{"res3_base4", 3, 4},   // Pentagon
		{"res3_base14", 3, 14}, // Pentagon

		// Invalid base cells should produce no output
		{"invalid_base122", 1, 122},
		{"invalid_base_neg", 1, -1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Calculate buffer size - max cells for a base cell at given resolution
			// Each base cell can have at most 7^res children
			bufferSize := _ipow(7, int64(tc.res))
			if bufferSize > 100000 {
				bufferSize = 100000 // Cap for memory safety in tests
			}

			// Collect Go results
			goResults := make([]H3Index, 0, bufferSize)
			_iterateBaseCellIndexesAtRes(tc.res, func(h H3Index) {
				goResults = append(goResults, h)
			}, tc.baseCell)

			// Collect C results
			cBuffer := make([]H3Index, bufferSize)
			cCount := iterateBaseCellIndexesAtResC(tc.res, tc.baseCell, cBuffer)
			cResults := cBuffer[:cCount]

			// Compare results
			if len(goResults) != len(cResults) {
				t.Fatalf("Count mismatch: Go=%d, C=%d", len(goResults), len(cResults))
			}

			// Compare each index
			for i := 0; i < len(goResults); i++ {
				if goResults[i] != cResults[i] {
					t.Errorf("Index mismatch at position %d: Go=0x%x, C=0x%x",
						i, goResults[i], cResults[i])
				}
			}
		})
	}
}
