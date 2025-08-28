//go:build cgo

package h3

import (
	"testing"
)

func Test_baseCellIsCwOffset_parity(t *testing.T) {
	tests := []struct {
		name     string
		baseCell int32
		testFace int32
	}{
		// Pentagon base cells with their clockwise offset faces
		{"pentagon_4_face_neg1", 4, -1}, // Expected from C data
		{"pentagon_4_face_0", 4, 0},
		{"pentagon_4_face_1", 4, 1},

		{"pentagon_14_face_2", 14, 2}, // Expected from C data
		{"pentagon_14_face_6", 14, 6}, // Expected from C data
		{"pentagon_14_face_0", 14, 0}, // Should be false
		{"pentagon_14_face_1", 14, 1}, // Should be false

		{"pentagon_24_face_1", 24, 1}, // Expected from C data
		{"pentagon_24_face_5", 24, 5}, // Expected from C data
		{"pentagon_24_face_0", 24, 0}, // Should be false

		{"pentagon_38_face_0", 38, 0}, // Expected from C data
		{"pentagon_38_face_4", 38, 4}, // Expected from C data
		{"pentagon_38_face_1", 38, 1}, // Should be false

		{"pentagon_49_face_0", 49, 0}, // Expected from C data
		{"pentagon_49_face_9", 49, 9}, // Expected from C data
		{"pentagon_49_face_1", 49, 1}, // Should be false

		{"pentagon_58_face_2", 58, 2}, // Expected from C data
		{"pentagon_58_face_8", 58, 8}, // Expected from C data
		{"pentagon_58_face_0", 58, 0}, // Should be false

		{"pentagon_63_face_3", 63, 3}, // Expected from C data
		{"pentagon_63_face_7", 63, 7}, // Expected from C data
		{"pentagon_63_face_0", 63, 0}, // Should be false

		{"pentagon_72_face_1", 72, 1},   // Expected from C data
		{"pentagon_72_face_11", 72, 11}, // Expected from C data
		{"pentagon_72_face_0", 72, 0},   // Should be false

		{"pentagon_82_face_4", 82, 4},   // Expected from C data
		{"pentagon_82_face_10", 82, 10}, // Expected from C data
		{"pentagon_82_face_0", 82, 0},   // Should be false

		{"pentagon_83_face_2", 83, 2},   // Expected from C data
		{"pentagon_83_face_12", 83, 12}, // Expected from C data
		{"pentagon_83_face_0", 83, 0},   // Should be false

		{"pentagon_97_face_5", 97, 5},   // Expected from C data
		{"pentagon_97_face_13", 97, 13}, // Expected from C data
		{"pentagon_97_face_0", 97, 0},   // Should be false

		{"pentagon_107_face_3", 107, 3},   // Expected from C data
		{"pentagon_107_face_15", 107, 15}, // Expected from C data
		{"pentagon_107_face_0", 107, 0},   // Should be false

		{"pentagon_117_face_0", 117, 0},   // Expected from C data
		{"pentagon_117_face_18", 117, 18}, // Expected from C data
		{"pentagon_117_face_1", 117, 1},   // Should be false

		// Non-pentagon base cells (should always return false)
		{"non_pentagon_0_face_0", 0, 0},
		{"non_pentagon_0_face_1", 0, 1},
		{"non_pentagon_1_face_0", 1, 0},
		{"non_pentagon_2_face_0", 2, 0},
		{"non_pentagon_50_face_5", 50, 5},
		{"non_pentagon_100_face_10", 100, 10},
		{"non_pentagon_121_face_19", 121, 19},

		// Edge cases
		{"pentagon_4_face_negative", 4, -10},
		{"pentagon_4_face_large", 4, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goResult := _baseCellIsCwOffset(tt.baseCell, tt.testFace)
			cResult := _baseCellIsCwOffsetC(tt.baseCell, tt.testFace)

			if goResult != cResult {
				t.Errorf("_baseCellIsCwOffset(%d, %d): Go=%t, C=%t",
					tt.baseCell, tt.testFace, goResult, cResult)
			}
		})
	}
}

func Test_baseCellIsCwOffset_all_base_cells(t *testing.T) {
	// Test all base cells with a range of face values
	testFaces := []int32{-1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 100}

	for baseCell := int32(0); baseCell < NUM_BASE_CELLS; baseCell++ {
		for _, testFace := range testFaces {
			t.Run("", func(t *testing.T) {
				goResult := _baseCellIsCwOffset(baseCell, testFace)
				cResult := _baseCellIsCwOffsetC(baseCell, testFace)

				if goResult != cResult {
					t.Errorf("_baseCellIsCwOffset(%d, %d): Go=%t, C=%t",
						baseCell, testFace, goResult, cResult)
				}
			})
		}
	}
}
