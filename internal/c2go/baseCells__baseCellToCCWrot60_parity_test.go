//go:build cgo

package c2go

import (
	"testing"
)

func Test_baseCellToCCWrot60_parity(t *testing.T) {
	tests := []struct {
		name     string
		baseCell int32
		face     int32
	}{
		// Test edge cases
		{"invalid_face_negative", 0, -1},
		{"invalid_face_too_large", 0, NUM_ICOSA_FACES + 1},
		{"invalid_face_way_too_large", 0, 100},

		// Test a selection of valid base cell/face combinations
		// Base cell 0 appears on faces 0, 1, 2
		{"base_cell_0_face_0", 0, 0},
		{"base_cell_0_face_1", 0, 1},
		{"base_cell_0_face_2", 0, 2},
		{"base_cell_0_face_3", 0, 3}, // Should return INVALID_ROTATIONS

		// Base cell 4 (pentagon) appears on faces 0, 1, 2, 3, 4
		{"base_cell_4_face_0", 4, 0},
		{"base_cell_4_face_1", 4, 1},
		{"base_cell_4_face_2", 4, 2},
		{"base_cell_4_face_3", 4, 3},
		{"base_cell_4_face_4", 4, 4},
		{"base_cell_4_face_5", 4, 5}, // Should return INVALID_ROTATIONS

		// Test some more base cells across different faces
		{"base_cell_16_face_0", 16, 0},
		{"base_cell_16_face_1", 16, 1},
		{"base_cell_16_face_4", 16, 4},
		{"base_cell_16_face_5", 16, 5}, // Should return INVALID_ROTATIONS

		{"base_cell_24_face_0", 24, 0},
		{"base_cell_24_face_1", 24, 1},
		{"base_cell_24_face_5", 24, 5},

		{"base_cell_58_face_3", 58, 3},
		{"base_cell_58_face_4", 58, 4},
		{"base_cell_58_face_8", 58, 8},
		{"base_cell_58_face_16", 58, 16},

		// Test some high numbered base cells
		{"base_cell_117_face_6", 117, 6},
		{"base_cell_117_face_17", 117, 17},
		{"base_cell_117_face_18", 117, 18},

		{"base_cell_121_face_6", 121, 6},
		{"base_cell_121_face_17", 121, 17},
		{"base_cell_121_face_18", 121, 18},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goResult := _baseCellToCCWrot60(tt.baseCell, tt.face)
			cResult := _baseCellToCCWrot60C(tt.baseCell, tt.face)

			if goResult != cResult {
				t.Errorf("_baseCellToCCWrot60(%d, %d): Go=%d, C=%d",
					tt.baseCell, tt.face, goResult, cResult)
			}
		})
	}
}

func Test_baseCellToCCWrot60_all_valid_combinations(t *testing.T) {
	// Test all base cells against all faces to ensure comprehensive coverage
	for baseCell := int32(0); baseCell < NUM_BASE_CELLS; baseCell++ {
		for face := int32(0); face < NUM_ICOSA_FACES; face++ {
			t.Run("", func(t *testing.T) {
				goResult := _baseCellToCCWrot60(baseCell, face)
				cResult := _baseCellToCCWrot60C(baseCell, face)

				if goResult != cResult {
					t.Errorf("_baseCellToCCWrot60(%d, %d): Go=%d, C=%d",
						baseCell, face, goResult, cResult)
				}
			})
		}
	}
}
