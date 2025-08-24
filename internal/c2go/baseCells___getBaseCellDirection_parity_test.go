//go:build cgo

package c2go

import (
	"testing"
)

func Test_getBaseCellDirection_parity(t *testing.T) {
	tests := []struct {
		name                string
		originBaseCell      int32
		neighboringBaseCell int32
	}{
		// Test neighboring base cells
		{"base_0_to_1", 0, 1},
		{"base_0_to_2", 0, 2},
		{"base_0_to_3", 0, 3},
		{"base_0_to_4", 0, 4},
		{"base_0_to_5", 0, 5},
		{"base_0_to_8", 0, 8},
		{"base_0_to_0", 0, 0}, // Self reference (CENTER_DIGIT)

		{"base_1_to_0", 1, 0},
		{"base_1_to_2", 1, 2},
		{"base_1_to_3", 1, 3},
		{"base_1_to_6", 1, 6},
		{"base_1_to_7", 1, 7},
		{"base_1_to_9", 1, 9},
		{"base_1_to_1", 1, 1}, // Self reference

		// Test non-neighboring base cells (should return INVALID_DIGIT)
		{"base_0_to_10", 0, 10}, // Not neighbors
		{"base_0_to_50", 0, 50}, // Not neighbors
		{"base_1_to_25", 1, 25}, // Not neighbors
		{"base_5_to_99", 5, 99}, // Not neighbors

		// Test pentagon base cells
		{"pentagon_4_to_0", 4, 0},
		{"pentagon_4_to_3", 4, 3},
		{"pentagon_4_to_8", 4, 8},
		{"pentagon_4_to_12", 4, 12},
		{"pentagon_4_to_15", 4, 15},
		{"pentagon_4_to_4", 4, 4}, // Self reference

		// Test some edge cases
		{"base_121_to_113", 121, 113},
		{"base_121_to_116", 121, 116},
		{"base_121_to_118", 121, 118},
		{"base_121_to_119", 121, 119},
		{"base_121_to_120", 121, 120},
		{"base_121_to_117", 121, 117},
		{"base_121_to_121", 121, 121}, // Self reference
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goResult := _getBaseCellDirection(tt.originBaseCell, tt.neighboringBaseCell)
			cResult := _getBaseCellDirectionC(tt.originBaseCell, tt.neighboringBaseCell)

			if goResult != cResult {
				t.Errorf("_getBaseCellDirection(%d, %d): Go=%d, C=%d",
					tt.originBaseCell, tt.neighboringBaseCell, goResult, cResult)
			}
		})
	}
}

func Test_getBaseCellDirection_comprehensive(t *testing.T) {
	// Test that direction lookup is consistent with neighbor lookup
	for baseCell := int32(0); baseCell < NUM_BASE_CELLS; baseCell++ {
		for dir := CENTER_DIGIT; dir < NUM_DIGITS; dir++ {
			neighbor := _getBaseCellNeighbor(baseCell, dir)
			if neighbor != INVALID_BASE_CELL {
				// If there's a valid neighbor, the direction lookup should work
				foundDir := _getBaseCellDirection(baseCell, neighbor)
				cFoundDir := _getBaseCellDirectionC(baseCell, neighbor)

				if foundDir != cFoundDir {
					t.Errorf("_getBaseCellDirection(baseCell=%d, neighbor=%d): Go=%d, C=%d",
						baseCell, neighbor, foundDir, cFoundDir)
				}

				// The found direction should match the original direction
				if foundDir != dir {
					t.Errorf("Direction mismatch: baseCell=%d, expected_dir=%d, found_dir=%d, neighbor=%d",
						baseCell, dir, foundDir, neighbor)
				}
			}
		}
	}
}

func Test_getBaseCellDirection_invalid_cases(t *testing.T) {
	// Test cases where base cells are not neighbors
	testCases := []struct {
		origin   int32
		neighbor int32
	}{
		{0, 50},   // Far apart
		{1, 100},  // Far apart
		{10, 90},  // Far apart
		{50, 121}, // Far apart
	}

	for _, tc := range testCases {
		goResult := _getBaseCellDirection(tc.origin, tc.neighbor)
		cResult := _getBaseCellDirectionC(tc.origin, tc.neighbor)

		if goResult != cResult {
			t.Errorf("_getBaseCellDirection(%d, %d): Go=%d, C=%d",
				tc.origin, tc.neighbor, goResult, cResult)
		}

		// Both should return INVALID_DIGIT
		if goResult != INVALID_DIGIT {
			t.Errorf("Expected INVALID_DIGIT for non-neighbors %d->%d, got %d",
				tc.origin, tc.neighbor, goResult)
		}
	}
}
