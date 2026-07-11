//go:build cgo && c2go

package h3

import (
	"testing"
)

func Test_getBaseCellNeighbor_parity(t *testing.T) {
	tests := []struct {
		name     string
		baseCell int32
		dir      direction
	}{
		// Test regular base cells (non-pentagon)
		{"base_0_center", 0, centerDigit},
		{"base_0_j", 0, jAxesDigit},
		{"base_0_jk", 0, jkAxesDigit},
		{"base_0_k", 0, kAxesDigit},
		{"base_0_ik", 0, ikAxesDigit},
		{"base_0_i", 0, iAxesDigit},
		{"base_0_ij", 0, ijAxesDigit},

		{"base_1_center", 1, centerDigit},
		{"base_1_j", 1, jAxesDigit},
		{"base_1_jk", 1, jkAxesDigit},
		{"base_1_k", 1, kAxesDigit},
		{"base_1_ik", 1, ikAxesDigit},
		{"base_1_i", 1, iAxesDigit},
		{"base_1_ij", 1, ijAxesDigit},

		// Test pentagon base cells (have invalidBaseCell neighbors)
		{"pentagon_4_center", 4, centerDigit},
		{"pentagon_4_j", 4, jAxesDigit}, // Should return invalidBaseCell
		{"pentagon_4_jk", 4, jkAxesDigit},
		{"pentagon_4_k", 4, kAxesDigit},
		{"pentagon_4_ik", 4, ikAxesDigit},
		{"pentagon_4_i", 4, iAxesDigit},
		{"pentagon_4_ij", 4, ijAxesDigit},

		{"pentagon_14_center", 14, centerDigit},
		{"pentagon_14_j", 14, jAxesDigit}, // Should return invalidBaseCell
		{"pentagon_14_jk", 14, jkAxesDigit},
		{"pentagon_14_k", 14, kAxesDigit},
		{"pentagon_14_ik", 14, ikAxesDigit},
		{"pentagon_14_i", 14, iAxesDigit},
		{"pentagon_14_ij", 14, ijAxesDigit},

		// Test other pentagon base cells
		{"pentagon_24_j", 24, jAxesDigit},   // Should return invalidBaseCell
		{"pentagon_38_j", 38, jAxesDigit},   // Should return invalidBaseCell
		{"pentagon_49_j", 49, jAxesDigit},   // Should return invalidBaseCell
		{"pentagon_58_j", 58, jAxesDigit},   // Should return invalidBaseCell
		{"pentagon_63_j", 63, jAxesDigit},   // Should return invalidBaseCell
		{"pentagon_72_j", 72, jAxesDigit},   // Should return invalidBaseCell
		{"pentagon_83_j", 83, jAxesDigit},   // Should return invalidBaseCell
		{"pentagon_97_j", 97, jAxesDigit},   // Should return invalidBaseCell
		{"pentagon_107_j", 107, jAxesDigit}, // Should return invalidBaseCell
		{"pentagon_117_j", 117, jAxesDigit}, // Should return invalidBaseCell

		// Test some edge cases - highest base cell
		{"base_121_center", 121, centerDigit},
		{"base_121_j", 121, jAxesDigit},
		{"base_121_k", 121, kAxesDigit},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goResult := _getBaseCellNeighbor(tt.baseCell, tt.dir)
			cResult := _getBaseCellNeighborC(tt.baseCell, tt.dir)

			if goResult != cResult {
				t.Errorf("_getBaseCellNeighbor(%d, %d): Go=%d, C=%d",
					tt.baseCell, tt.dir, goResult, cResult)
			}
		})
	}
}

func Test_getBaseCellNeighbor_all_combinations(t *testing.T) {
	// Test all valid base cells with all directions
	for baseCell := int32(0); baseCell < numBaseCells; baseCell++ {
		for dir := centerDigit; dir < numDigits; dir++ {
			goResult := _getBaseCellNeighbor(baseCell, dir)
			cResult := _getBaseCellNeighborC(baseCell, dir)

			if goResult != cResult {
				t.Errorf("_getBaseCellNeighbor(baseCell=%d, dir=%d): Go=%d, C=%d",
					baseCell, dir, goResult, cResult)
			}
		}
	}
}
