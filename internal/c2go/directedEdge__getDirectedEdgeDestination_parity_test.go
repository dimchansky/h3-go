//go:build cgo

package c2go

import (
	"testing"
)

func Test_getDirectedEdgeDestination_parity(t *testing.T) {
	tests := []struct {
		name   string
		origin H3Index
		dir    int // Direction for edge creation (1-6)
	}{
		// Regular hexagons - test all 6 directions
		{"hex_dir_1", H3Index(0x85283473fffffff), 1}, // J_AXES_DIGIT
		{"hex_dir_2", H3Index(0x85283473fffffff), 2}, // JK_AXES_DIGIT
		{"hex_dir_3", H3Index(0x85283473fffffff), 3}, // K_AXES_DIGIT
		{"hex_dir_4", H3Index(0x85283473fffffff), 4}, // IK_AXES_DIGIT
		{"hex_dir_5", H3Index(0x85283473fffffff), 5}, // I_AXES_DIGIT
		{"hex_dir_6", H3Index(0x85283473fffffff), 6}, // IJ_AXES_DIGIT

		// Different resolutions
		{"res0_hex", H3Index(0x8001fffffffffff), 1},
		{"res1_hex", H3Index(0x8101fffffffffff), 1},
		{"res2_hex", H3Index(0x8201fffffffffff), 1},
		{"res3_hex", H3Index(0x8301fffffffffff), 1},
		{"res4_hex", H3Index(0x8401fffffffffff), 1},
		{"res5_hex", H3Index(0x8501fffffffffff), 1},

		// Pentagon base cells (K direction should be invalid)
		{"pentagon_4_j", H3Index(0x8004fffffffffff), 1},  // J_AXES_DIGIT - should work
		{"pentagon_4_jk", H3Index(0x8004fffffffffff), 2}, // JK_AXES_DIGIT - should work
		{"pentagon_4_ik", H3Index(0x8004fffffffffff), 4}, // IK_AXES_DIGIT - should work
		{"pentagon_4_i", H3Index(0x8004fffffffffff), 5},  // I_AXES_DIGIT - should work
		{"pentagon_4_ij", H3Index(0x8004fffffffffff), 6}, // IJ_AXES_DIGIT - should work

		// More pentagon base cells
		{"pentagon_14_j", H3Index(0x800efffffffffff), 1},
		{"pentagon_24_jk", H3Index(0x8018fffffffffff), 2},
		{"pentagon_38_ik", H3Index(0x8026fffffffffff), 4},
		{"pentagon_49_i", H3Index(0x8031fffffffffff), 5},
		{"pentagon_58_ij", H3Index(0x803afffffffffff), 6},

		// Different base cells
		{"base_0", H3Index(0x8000fffffffffff), 1},
		{"base_10", H3Index(0x800afffffffffff), 2},
		{"base_20", H3Index(0x8014fffffffffff), 3},
		{"base_30", H3Index(0x801efffffffffff), 4},
		{"base_40", H3Index(0x8028fffffffffff), 5},
		{"base_50", H3Index(0x8032fffffffffff), 6},

		// High resolution cells
		{"res10_hex", H3Index(0x8a283470c44ffff), 1},
		{"res15_hex", H3Index(0x8f28347333fffff), 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create directed edge
			edge := tt.origin
			edge = setMode(edge, H3_DIRECTEDEDGE_MODE)
			edge = setReservedBits(edge, int32(tt.dir))

			// Test C implementation
			cOut, cErr := getDirectedEdgeDestinationC(edge)

			// Test Go implementation
			var goOut H3Index
			goErr := getDirectedEdgeDestination(edge, &goOut)

			// Compare error results
			if cErr != goErr {
				t.Errorf("Error mismatch for edge from origin=0x%x, dir=%d: C=%v, Go=%v",
					tt.origin, tt.dir, cErr, goErr)
				return
			}

			// If both succeeded, compare outputs
			if cErr == E_SUCCESS {
				if cOut != goOut {
					t.Errorf("Output mismatch for edge from origin=0x%x, dir=%d: C=0x%x, Go=0x%x",
						tt.origin, tt.dir, cOut, goOut)
				}
			}
		})
	}
}

func Test_getDirectedEdgeDestination_invalid_parity(t *testing.T) {
	// Test invalid edges
	tests := []struct {
		name        string
		edge        H3Index
		shouldFail  bool
		expectedErr H3Error
	}{
		// Invalid mode (not a directed edge)
		{"not_edge_mode", H3Index(0x85283473fffffff), true, E_DIR_EDGE_INVALID}, // Regular cell mode

		// Direction 0 (CENTER_DIGIT) - actually valid, returns same cell
		{"dir_0_center", func() H3Index {
			edge := H3Index(0x85283473fffffff)
			edge = setMode(edge, H3_DIRECTEDEDGE_MODE)
			edge = setReservedBits(edge, 0) // CENTER_DIGIT returns same cell
			return edge
		}(), false, E_SUCCESS},

		// Invalid direction (>= NUM_DIGITS)
		{"invalid_dir_7", func() H3Index {
			edge := H3Index(0x85283473fffffff)
			edge = setMode(edge, H3_DIRECTEDEDGE_MODE)
			edge = setReservedBits(edge, 7) // Out of range
			return edge
		}(), true, E_FAILED},

		// Pentagon with K direction (should fail with E_PENTAGON)
		{"pentagon_k_dir", func() H3Index {
			edge := H3Index(0x8004fffffffffff) // Pentagon base cell 4
			edge = setMode(edge, H3_DIRECTEDEDGE_MODE)
			edge = setReservedBits(edge, 3) // K_AXES_DIGIT
			return edge
		}(), false, E_SUCCESS}, // Actually succeeds in traversal
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test C implementation
			cOut, cErr := getDirectedEdgeDestinationC(tt.edge)

			// Test Go implementation
			var goOut H3Index
			goErr := getDirectedEdgeDestination(tt.edge, &goOut)

			// Compare error results
			if cErr != goErr {
				t.Errorf("Error mismatch for edge=0x%x: C=%v, Go=%v",
					tt.edge, cErr, goErr)
			}

			// If we expected success, compare outputs
			if cErr == E_SUCCESS && goErr == E_SUCCESS {
				if cOut != goOut {
					t.Errorf("Output mismatch for edge=0x%x: C=0x%x, Go=0x%x",
						tt.edge, cOut, goOut)
				}
			}
		})
	}
}
