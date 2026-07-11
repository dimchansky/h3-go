//go:build cgo && c2go

package h3

import (
	"testing"
)

func Test_h3NeighborRotations_parity(t *testing.T) {
	tests := []struct {
		name      string
		origin    h3Index
		dir       direction
		rotations int32
	}{
		// Basic hexagon neighbors
		{"hex_k_axis", h3Index(0x8001fffffffffff), kAxesDigit, 0},
		{"hex_j_axis", h3Index(0x8001fffffffffff), jAxesDigit, 0},
		{"hex_jk_axis", h3Index(0x8001fffffffffff), jkAxesDigit, 0},
		{"hex_i_axis", h3Index(0x8001fffffffffff), iAxesDigit, 0},
		{"hex_ik_axis", h3Index(0x8001fffffffffff), ikAxesDigit, 0},
		{"hex_ij_axis", h3Index(0x8001fffffffffff), ijAxesDigit, 0},

		// With some initial rotations
		{"hex_with_rotations_1", h3Index(0x8001fffffffffff), kAxesDigit, 1},
		{"hex_with_rotations_2", h3Index(0x8001fffffffffff), kAxesDigit, 2},
		{"hex_with_rotations_3", h3Index(0x8001fffffffffff), kAxesDigit, 3},

		// Different resolutions
		{"res1_hex", h3Index(0x8101fffffffffff), kAxesDigit, 0},
		{"res2_hex", h3Index(0x8201fffffffffff), kAxesDigit, 0},
		{"res3_hex", h3Index(0x8301fffffffffff), kAxesDigit, 0},

		// Pentagon base cells - should fail with ePentagon for kAxesDigit
		{"pentagon_base_4_k", h3Index(0x8004fffffffffff), kAxesDigit, 0},  // Should fail
		{"pentagon_base_4_j", h3Index(0x8004fffffffffff), jAxesDigit, 0},  // Should work
		{"pentagon_base_14_k", h3Index(0x800efffffffffff), kAxesDigit, 0}, // Should fail
		{"pentagon_base_14_j", h3Index(0x800efffffffffff), jAxesDigit, 0}, // Should work

		// More pentagon base cells
		{"pentagon_base_24", h3Index(0x8018fffffffffff), jAxesDigit, 0},
		{"pentagon_base_38", h3Index(0x8026fffffffffff), jAxesDigit, 0},
		{"pentagon_base_49", h3Index(0x8031fffffffffff), jAxesDigit, 0},
		{"pentagon_base_58", h3Index(0x803afffffffffff), jAxesDigit, 0},

		// Edge cases with high rotations
		{"high_rotation_6", h3Index(0x8001fffffffffff), kAxesDigit, 6},
		{"high_rotation_7", h3Index(0x8001fffffffffff), kAxesDigit, 7},
		{"high_rotation_12", h3Index(0x8001fffffffffff), kAxesDigit, 12},

		// Different base cells
		{"base_cell_0", h3Index(0x8000fffffffffff), kAxesDigit, 0},
		{"base_cell_1", h3Index(0x8001fffffffffff), kAxesDigit, 0},
		{"base_cell_2", h3Index(0x8002fffffffffff), kAxesDigit, 0},
		{"base_cell_10", h3Index(0x800afffffffffff), kAxesDigit, 0},
		{"base_cell_20", h3Index(0x8014fffffffffff), kAxesDigit, 0},

		// Failing case from cellToVertex debug - should cause parity test to fail
		{"cellToVertex_failing_case", h3Index(0x08015fffffffffff), iAxesDigit, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test C implementation
			cRotations := tt.rotations
			cOut, cResult := h3NeighborRotationsC(tt.origin, tt.dir, &cRotations)

			// Test Go implementation
			var goOut h3Index
			goRotations := tt.rotations
			goResult := h3NeighborRotations(tt.origin, tt.dir, &goRotations, &goOut)

			// Compare error results
			if cResult != goResult {
				t.Errorf("Error mismatch for origin=0x%x, dir=%d, rotations=%d: C=%v, Go=%v",
					tt.origin, tt.dir, tt.rotations, cResult, goResult)
				return
			}

			// If both succeeded, compare outputs
			if cResult == eSuccess && goResult == eSuccess {
				if cOut != goOut {
					t.Errorf("Output mismatch for origin=0x%x, dir=%d: C=0x%x, Go=0x%x",
						tt.origin, tt.dir, cOut, goOut)
				}
				if cRotations != goRotations {
					t.Errorf("Rotations mismatch for origin=0x%x, dir=%d: C=%d, Go=%d",
						tt.origin, tt.dir, cRotations, goRotations)
				}
			}
		})
	}
}

func Test_h3NeighborRotations_comprehensive_parity(t *testing.T) {
	// Test with a variety of origins, directions, and rotations
	origins := []h3Index{
		h3Index(0x8001fffffffffff), // Basic hex
		h3Index(0x8101fffffffffff), // Res 1
		h3Index(0x8201fffffffffff), // Res 2
		h3Index(0x8004fffffffffff), // Pentagon base 4
		h3Index(0x800efffffffffff), // Pentagon base 14
		h3Index(0x8018fffffffffff), // Pentagon base 24
	}

	directions := []direction{jAxesDigit, jkAxesDigit, kAxesDigit, ikAxesDigit, iAxesDigit, ijAxesDigit}
	rotationsValues := []int32{0, 1, 2, 3, 4, 5}

	for _, origin := range origins {
		for _, dir := range directions {
			for _, rot := range rotationsValues {
				t.Run("comprehensive", func(t *testing.T) {
					// Test C implementation
					cRotations := rot
					cOut, cResult := h3NeighborRotationsC(origin, dir, &cRotations)

					// Test Go implementation
					var goOut h3Index
					goRotations := rot
					goResult := h3NeighborRotations(origin, dir, &goRotations, &goOut)

					// Compare error results
					if cResult != goResult {
						t.Errorf("Error mismatch for origin=0x%x, dir=%d, rotations=%d: C=%v, Go=%v",
							origin, dir, rot, cResult, goResult)
						return
					}

					// If both succeeded, compare outputs
					if cResult == eSuccess && goResult == eSuccess {
						if cOut != goOut {
							t.Errorf("Output mismatch for origin=0x%x, dir=%d: C=0x%x, Go=0x%x",
								origin, dir, cOut, goOut)
						}
						if cRotations != goRotations {
							t.Errorf("Rotations mismatch for origin=0x%x, dir=%d: C=%d, Go=%d",
								origin, dir, cRotations, goRotations)
						}
					}
				})
			}
		}
	}
}
