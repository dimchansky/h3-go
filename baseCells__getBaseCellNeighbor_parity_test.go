//go:build cgo

package h3

import (
	"testing"
)

func Test_getBaseCellNeighbor_parity(t *testing.T) {
	tests := []struct {
		name     string
		baseCell int32
		dir      Direction
	}{
		// Test regular base cells (non-pentagon)
		{"base_0_center", 0, CENTER_DIGIT},
		{"base_0_j", 0, J_AXES_DIGIT},
		{"base_0_jk", 0, JK_AXES_DIGIT},
		{"base_0_k", 0, K_AXES_DIGIT},
		{"base_0_ik", 0, IK_AXES_DIGIT},
		{"base_0_i", 0, I_AXES_DIGIT},
		{"base_0_ij", 0, IJ_AXES_DIGIT},

		{"base_1_center", 1, CENTER_DIGIT},
		{"base_1_j", 1, J_AXES_DIGIT},
		{"base_1_jk", 1, JK_AXES_DIGIT},
		{"base_1_k", 1, K_AXES_DIGIT},
		{"base_1_ik", 1, IK_AXES_DIGIT},
		{"base_1_i", 1, I_AXES_DIGIT},
		{"base_1_ij", 1, IJ_AXES_DIGIT},

		// Test pentagon base cells (have INVALID_BASE_CELL neighbors)
		{"pentagon_4_center", 4, CENTER_DIGIT},
		{"pentagon_4_j", 4, J_AXES_DIGIT}, // Should return INVALID_BASE_CELL
		{"pentagon_4_jk", 4, JK_AXES_DIGIT},
		{"pentagon_4_k", 4, K_AXES_DIGIT},
		{"pentagon_4_ik", 4, IK_AXES_DIGIT},
		{"pentagon_4_i", 4, I_AXES_DIGIT},
		{"pentagon_4_ij", 4, IJ_AXES_DIGIT},

		{"pentagon_14_center", 14, CENTER_DIGIT},
		{"pentagon_14_j", 14, J_AXES_DIGIT}, // Should return INVALID_BASE_CELL
		{"pentagon_14_jk", 14, JK_AXES_DIGIT},
		{"pentagon_14_k", 14, K_AXES_DIGIT},
		{"pentagon_14_ik", 14, IK_AXES_DIGIT},
		{"pentagon_14_i", 14, I_AXES_DIGIT},
		{"pentagon_14_ij", 14, IJ_AXES_DIGIT},

		// Test other pentagon base cells
		{"pentagon_24_j", 24, J_AXES_DIGIT},   // Should return INVALID_BASE_CELL
		{"pentagon_38_j", 38, J_AXES_DIGIT},   // Should return INVALID_BASE_CELL
		{"pentagon_49_j", 49, J_AXES_DIGIT},   // Should return INVALID_BASE_CELL
		{"pentagon_58_j", 58, J_AXES_DIGIT},   // Should return INVALID_BASE_CELL
		{"pentagon_63_j", 63, J_AXES_DIGIT},   // Should return INVALID_BASE_CELL
		{"pentagon_72_j", 72, J_AXES_DIGIT},   // Should return INVALID_BASE_CELL
		{"pentagon_83_j", 83, J_AXES_DIGIT},   // Should return INVALID_BASE_CELL
		{"pentagon_97_j", 97, J_AXES_DIGIT},   // Should return INVALID_BASE_CELL
		{"pentagon_107_j", 107, J_AXES_DIGIT}, // Should return INVALID_BASE_CELL
		{"pentagon_117_j", 117, J_AXES_DIGIT}, // Should return INVALID_BASE_CELL

		// Test some edge cases - highest base cell
		{"base_121_center", 121, CENTER_DIGIT},
		{"base_121_j", 121, J_AXES_DIGIT},
		{"base_121_k", 121, K_AXES_DIGIT},
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
	for baseCell := int32(0); baseCell < NUM_BASE_CELLS; baseCell++ {
		for dir := CENTER_DIGIT; dir < NUM_DIGITS; dir++ {
			goResult := _getBaseCellNeighbor(baseCell, dir)
			cResult := _getBaseCellNeighborC(baseCell, dir)

			if goResult != cResult {
				t.Errorf("_getBaseCellNeighbor(baseCell=%d, dir=%d): Go=%d, C=%d",
					baseCell, dir, goResult, cResult)
			}
		}
	}
}
