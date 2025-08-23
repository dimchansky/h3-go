//go:build cgo

package c2go

import "testing"

func Test_hasDeletedSubsequence_parity(t *testing.T) {
	tests := []struct {
		name     string
		h        H3Index
		baseCell int32
	}{
		// Pentagon base cells (4, 14, 24, 38, 49, 58, 63, 72, 83, 97, 107, 117)
		{"pentagon 4", 0x8201fffffffffff, 4},
		{"pentagon 14", 0x821c07fffffffff, 14},
		{"pentagon 24", 0x8301fffffffffff, 24},
		{"pentagon 38", 0x8401fffffffffff, 38},
		{"pentagon 49", 0x8501fffffffffff, 49},
		{"pentagon 58", 0x8601fffffffffff, 58},
		{"pentagon 63", 0x8701fffffffffff, 63},
		{"pentagon 72", 0x8801fffffffffff, 72},
		{"pentagon 83", 0x8901fffffffffff, 83},
		{"pentagon 97", 0x8a01fffffffffff, 97},
		{"pentagon 107", 0x8b01fffffffffff, 107},
		{"pentagon 117", 0x8c01fffffffffff, 117},

		// Non-pentagon base cells (should always return false)
		{"non-pentagon 0", 0x8001fffffffffff, 0},
		{"non-pentagon 1", 0x8011fffffffffff, 1},
		{"non-pentagon 50", 0x8501fffffffffff, 50},   // Not pentagon
		{"non-pentagon 100", 0x8a01fffffffffff, 100}, // Not pentagon

		// Pentagon with various digit patterns
		{"pentagon 4 with 1 digit", 0x8201200000000000, 4}, // Has K_AXES_DIGIT (1)
		{"pentagon 4 all zeros", 0x8200000000000000, 4},    // All zeros after base cell
		{"pentagon 14 mixed", 0x821c076543210000, 14},

		// Edge cases
		{"invalid base cell negative", 0x8001fffffffffff, -1},
		{"invalid base cell large", 0x8001fffffffffff, 128},
		{"zero index pentagon", 0, 4},
		{"zero index non-pentagon", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call C implementation
			gotC := hasDeletedSubsequenceC(tt.h, tt.baseCell)

			// Call Go implementation
			gotGo := _hasDeletedSubsequence(tt.h, tt.baseCell)

			// Compare results
			if gotGo != gotC {
				t.Errorf("_hasDeletedSubsequence(%x, %d) mismatch: Go=%v != C=%v",
					tt.h, tt.baseCell, gotGo, gotC)
			}
		})
	}

	// Test all pentagon base cells
	t.Run("all_pentagon_base_cells", func(t *testing.T) {
		pentagons := []int32{4, 14, 24, 38, 49, 58, 63, 72, 83, 97, 107, 117}
		h := H3Index(0x8201fffffffffff)

		for _, baseCell := range pentagons {
			gotC := hasDeletedSubsequenceC(h, baseCell)
			gotGo := _hasDeletedSubsequence(h, baseCell)

			if gotGo != gotC {
				t.Errorf("Pentagon base cell %d: Go=%v != C=%v", baseCell, gotGo, gotC)
			}
		}
	})

	// Test all non-pentagon base cells (should return false)
	t.Run("all_non_pentagon_base_cells", func(t *testing.T) {
		h := H3Index(0x8001fffffffffff)
		pentagons := map[int32]bool{4: true, 14: true, 24: true, 38: true, 49: true, 58: true, 63: true, 72: true, 83: true, 97: true, 107: true, 117: true}

		for baseCell := int32(0); baseCell < 122; baseCell++ {
			if !pentagons[baseCell] {
				gotC := hasDeletedSubsequenceC(h, baseCell)
				gotGo := _hasDeletedSubsequence(h, baseCell)

				if gotGo != gotC {
					t.Errorf("Non-pentagon base cell %d: Go=%v != C=%v", baseCell, gotGo, gotC)
				}

				// Non-pentagons should always return false
				if gotGo {
					t.Errorf("Non-pentagon base cell %d should return false, got true", baseCell)
				}
			}
		}
	})
}
