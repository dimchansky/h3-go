//go:build cgo && c2go

package h3

import (
	"testing"
)

func Test_getDirectedEdgeDestination_parity(t *testing.T) {
	tests := []struct {
		name   string
		origin h3Index
		dir    int // direction for edge creation (1-6)
	}{
		// Regular hexagons - test all 6 directions
		{"hex_dir_1", h3Index(0x85283473fffffff), 1}, // jAxesDigit
		{"hex_dir_2", h3Index(0x85283473fffffff), 2}, // jkAxesDigit
		{"hex_dir_3", h3Index(0x85283473fffffff), 3}, // kAxesDigit
		{"hex_dir_4", h3Index(0x85283473fffffff), 4}, // ikAxesDigit
		{"hex_dir_5", h3Index(0x85283473fffffff), 5}, // iAxesDigit
		{"hex_dir_6", h3Index(0x85283473fffffff), 6}, // ijAxesDigit

		// Different resolutions
		{"res0_hex", h3Index(0x8001fffffffffff), 1},
		{"res1_hex", h3Index(0x8101fffffffffff), 1},
		{"res2_hex", h3Index(0x8201fffffffffff), 1},
		{"res3_hex", h3Index(0x8301fffffffffff), 1},
		{"res4_hex", h3Index(0x8401fffffffffff), 1},
		{"res5_hex", h3Index(0x8501fffffffffff), 1},

		// Pentagon base cells (K direction should be invalid)
		{"pentagon_4_j", h3Index(0x8004fffffffffff), 1},  // jAxesDigit - should work
		{"pentagon_4_jk", h3Index(0x8004fffffffffff), 2}, // jkAxesDigit - should work
		{"pentagon_4_ik", h3Index(0x8004fffffffffff), 4}, // ikAxesDigit - should work
		{"pentagon_4_i", h3Index(0x8004fffffffffff), 5},  // iAxesDigit - should work
		{"pentagon_4_ij", h3Index(0x8004fffffffffff), 6}, // ijAxesDigit - should work

		// More pentagon base cells
		{"pentagon_14_j", h3Index(0x800efffffffffff), 1},
		{"pentagon_24_jk", h3Index(0x8018fffffffffff), 2},
		{"pentagon_38_ik", h3Index(0x8026fffffffffff), 4},
		{"pentagon_49_i", h3Index(0x8031fffffffffff), 5},
		{"pentagon_58_ij", h3Index(0x803afffffffffff), 6},

		// Different base cells
		{"base_0", h3Index(0x8000fffffffffff), 1},
		{"base_10", h3Index(0x800afffffffffff), 2},
		{"base_20", h3Index(0x8014fffffffffff), 3},
		{"base_30", h3Index(0x801efffffffffff), 4},
		{"base_40", h3Index(0x8028fffffffffff), 5},
		{"base_50", h3Index(0x8032fffffffffff), 6},

		// High resolution cells
		{"res10_hex", h3Index(0x8a283470c44ffff), 1},
		{"res15_hex", h3Index(0x8f28347333fffff), 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create directed edge
			edge := tt.origin
			edge = setMode(edge, h3DirectededgeMode)
			edge = setReservedBits(edge, int32(tt.dir))

			// Test C implementation
			cOut, cErr := getDirectedEdgeDestinationC(edge)

			// Test Go implementation
			var goOut h3Index
			goErr := getDirectedEdgeDestination(edge, &goOut)

			// Compare error results
			if cErr != goErr {
				t.Errorf("Error mismatch for edge from origin=0x%x, dir=%d: C=%v, Go=%v",
					tt.origin, tt.dir, cErr, goErr)
				return
			}

			// If both succeeded, compare outputs
			if cErr == eSuccess {
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
		edge        h3Index
		shouldFail  bool
		expectedErr h3Error
	}{
		// Invalid mode (not a directed edge)
		{"not_edge_mode", h3Index(0x85283473fffffff), true, eDirEdgeInvalid}, // Regular cell mode

		// direction 0 (centerDigit) - actually valid, returns same cell
		{"dir_0_center", func() h3Index {
			edge := h3Index(0x85283473fffffff)
			edge = setMode(edge, h3DirectededgeMode)
			edge = setReservedBits(edge, 0) // centerDigit returns same cell
			return edge
		}(), false, eSuccess},

		// Invalid direction (>= numDigits)
		{"invalid_dir_7", func() h3Index {
			edge := h3Index(0x85283473fffffff)
			edge = setMode(edge, h3DirectededgeMode)
			edge = setReservedBits(edge, 7) // Out of range
			return edge
		}(), true, eFailed},

		// Pentagon with K direction (should fail with ePentagon)
		{"pentagon_k_dir", func() h3Index {
			edge := h3Index(0x8004fffffffffff) // Pentagon base cell 4
			edge = setMode(edge, h3DirectededgeMode)
			edge = setReservedBits(edge, 3) // kAxesDigit
			return edge
		}(), false, eSuccess}, // Actually succeeds in traversal
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test C implementation
			cOut, cErr := getDirectedEdgeDestinationC(tt.edge)

			// Test Go implementation
			var goOut h3Index
			goErr := getDirectedEdgeDestination(tt.edge, &goOut)

			// Compare error results
			if cErr != goErr {
				t.Errorf("Error mismatch for edge=0x%x: C=%v, Go=%v",
					tt.edge, cErr, goErr)
			}

			// If we expected success, compare outputs
			if cErr == eSuccess && goErr == eSuccess {
				if cOut != goOut {
					t.Errorf("Output mismatch for edge=0x%x: C=0x%x, Go=0x%x",
						tt.edge, cOut, goOut)
				}
			}
		})
	}
}
