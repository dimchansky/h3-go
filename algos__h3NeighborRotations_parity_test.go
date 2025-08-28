//go:build cgo

package h3

import (
	"testing"
)

func Test_h3NeighborRotations_parity(t *testing.T) {
	tests := []struct {
		name      string
		origin    H3Index
		dir       Direction
		rotations int32
	}{
		// Basic hexagon neighbors
		{"hex_k_axis", H3Index(0x8001fffffffffff), K_AXES_DIGIT, 0},
		{"hex_j_axis", H3Index(0x8001fffffffffff), J_AXES_DIGIT, 0},
		{"hex_jk_axis", H3Index(0x8001fffffffffff), JK_AXES_DIGIT, 0},
		{"hex_i_axis", H3Index(0x8001fffffffffff), I_AXES_DIGIT, 0},
		{"hex_ik_axis", H3Index(0x8001fffffffffff), IK_AXES_DIGIT, 0},
		{"hex_ij_axis", H3Index(0x8001fffffffffff), IJ_AXES_DIGIT, 0},

		// With some initial rotations
		{"hex_with_rotations_1", H3Index(0x8001fffffffffff), K_AXES_DIGIT, 1},
		{"hex_with_rotations_2", H3Index(0x8001fffffffffff), K_AXES_DIGIT, 2},
		{"hex_with_rotations_3", H3Index(0x8001fffffffffff), K_AXES_DIGIT, 3},

		// Different resolutions
		{"res1_hex", H3Index(0x8101fffffffffff), K_AXES_DIGIT, 0},
		{"res2_hex", H3Index(0x8201fffffffffff), K_AXES_DIGIT, 0},
		{"res3_hex", H3Index(0x8301fffffffffff), K_AXES_DIGIT, 0},

		// Pentagon base cells - should fail with E_PENTAGON for K_AXES_DIGIT
		{"pentagon_base_4_k", H3Index(0x8004fffffffffff), K_AXES_DIGIT, 0},  // Should fail
		{"pentagon_base_4_j", H3Index(0x8004fffffffffff), J_AXES_DIGIT, 0},  // Should work
		{"pentagon_base_14_k", H3Index(0x800efffffffffff), K_AXES_DIGIT, 0}, // Should fail
		{"pentagon_base_14_j", H3Index(0x800efffffffffff), J_AXES_DIGIT, 0}, // Should work

		// More pentagon base cells
		{"pentagon_base_24", H3Index(0x8018fffffffffff), J_AXES_DIGIT, 0},
		{"pentagon_base_38", H3Index(0x8026fffffffffff), J_AXES_DIGIT, 0},
		{"pentagon_base_49", H3Index(0x8031fffffffffff), J_AXES_DIGIT, 0},
		{"pentagon_base_58", H3Index(0x803afffffffffff), J_AXES_DIGIT, 0},

		// Edge cases with high rotations
		{"high_rotation_6", H3Index(0x8001fffffffffff), K_AXES_DIGIT, 6},
		{"high_rotation_7", H3Index(0x8001fffffffffff), K_AXES_DIGIT, 7},
		{"high_rotation_12", H3Index(0x8001fffffffffff), K_AXES_DIGIT, 12},

		// Different base cells
		{"base_cell_0", H3Index(0x8000fffffffffff), K_AXES_DIGIT, 0},
		{"base_cell_1", H3Index(0x8001fffffffffff), K_AXES_DIGIT, 0},
		{"base_cell_2", H3Index(0x8002fffffffffff), K_AXES_DIGIT, 0},
		{"base_cell_10", H3Index(0x800afffffffffff), K_AXES_DIGIT, 0},
		{"base_cell_20", H3Index(0x8014fffffffffff), K_AXES_DIGIT, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test C implementation
			cRotations := tt.rotations
			cOut, cResult := h3NeighborRotationsC(tt.origin, tt.dir, &cRotations)

			// Test Go implementation
			var goOut H3Index
			goRotations := tt.rotations
			goResult := h3NeighborRotations(tt.origin, tt.dir, &goRotations, &goOut)

			// Compare error results
			if cResult != goResult {
				t.Errorf("Error mismatch for origin=0x%x, dir=%d, rotations=%d: C=%v, Go=%v",
					tt.origin, tt.dir, tt.rotations, cResult, goResult)
				return
			}

			// If both succeeded, compare outputs
			if cResult == E_SUCCESS && goResult == E_SUCCESS {
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
	origins := []H3Index{
		H3Index(0x8001fffffffffff), // Basic hex
		H3Index(0x8101fffffffffff), // Res 1
		H3Index(0x8201fffffffffff), // Res 2
		H3Index(0x8004fffffffffff), // Pentagon base 4
		H3Index(0x800efffffffffff), // Pentagon base 14
		H3Index(0x8018fffffffffff), // Pentagon base 24
	}

	directions := []Direction{J_AXES_DIGIT, JK_AXES_DIGIT, K_AXES_DIGIT, IK_AXES_DIGIT, I_AXES_DIGIT, IJ_AXES_DIGIT}
	rotationsValues := []int32{0, 1, 2, 3, 4, 5}

	for _, origin := range origins {
		for _, dir := range directions {
			for _, rot := range rotationsValues {
				t.Run("comprehensive", func(t *testing.T) {
					// Test C implementation
					cRotations := rot
					cOut, cResult := h3NeighborRotationsC(origin, dir, &cRotations)

					// Test Go implementation
					var goOut H3Index
					goRotations := rot
					goResult := h3NeighborRotations(origin, dir, &goRotations, &goOut)

					// Compare error results
					if cResult != goResult {
						t.Errorf("Error mismatch for origin=0x%x, dir=%d, rotations=%d: C=%v, Go=%v",
							origin, dir, rot, cResult, goResult)
						return
					}

					// If both succeeded, compare outputs
					if cResult == E_SUCCESS && goResult == E_SUCCESS {
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
