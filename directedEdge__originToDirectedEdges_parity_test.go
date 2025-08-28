//go:build cgo

package h3

import (
	"testing"
)

func Test_originToDirectedEdges_parity(t *testing.T) {
	tests := []struct {
		name   string
		origin H3Index
	}{
		// Valid hexagon cells at different resolutions
		{"simple_res0", H3Index(0x8001fffffffffff)}, // Simple res=0 cell
		{"simple_res1", H3Index(0x8101fffffffffff)}, // Simple res=1 cell
		{"simple_res2", H3Index(0x8201fffffffffff)}, // Simple res=2 cell
		{"simple_res5", H3Index(0x8501ffffffffff)},  // Simple res=5 cell

		// Different base cell patterns
		{"base_cell_0", H3Index(0x8001fffffffffff)},  // Base cell 0
		{"base_cell_1", H3Index(0x8021fffffffffff)},  // Base cell 1
		{"base_cell_10", H3Index(0x8281fffffffffff)}, // Base cell 10

		// Edge cases
		{"zero_index", H3Index(0x0)},             // Should fail but test for consistency
		{"max_base", H3Index(0x8f81fffffffffff)}, // High base cell
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare output slices
			cEdges := make([]H3Index, 6)
			goEdges := make([]H3Index, 6)

			// Call C implementation
			cErr := originToDirectedEdgesC(tt.origin, cEdges)

			// Call Go implementation
			goErr := originToDirectedEdges(tt.origin, goEdges)

			// Compare error codes
			if cErr != goErr {
				t.Errorf("Error mismatch for origin 0x%x: C=%v, Go=%v", tt.origin, cErr, goErr)
				return
			}

			// If successful, compare all edge results
			if cErr == E_SUCCESS {
				for i := 0; i < 6; i++ {
					if cEdges[i] != goEdges[i] {
						t.Errorf("Edge[%d] mismatch for origin 0x%x: C=0x%x, Go=0x%x",
							i, tt.origin, cEdges[i], goEdges[i])
					}
				}
			}
		})
	}
}

func Test_originToDirectedEdges_pentagon_parity(t *testing.T) {
	// Test with pentagon base cells
	// Pentagon base cells are: 4, 14, 24, 38, 49, 58, 63, 72, 83, 97, 107, 117
	pentagonBaseCells := []int32{4, 14, 24, 38, 49, 58, 63, 72, 83, 97, 107, 117}

	for _, baseCell := range pentagonBaseCells {
		// Create pentagon cells at different resolutions
		for res := int32(0); res <= 3; res++ {
			pentagon := setMode(H3Index(0), H3_CELL_MODE)
			pentagon = setBaseCell(pentagon, baseCell)
			pentagon = setResolution(pentagon, res)

			t.Run("pentagon_test", func(t *testing.T) {
				// Prepare output slices
				cEdges := make([]H3Index, 6)
				goEdges := make([]H3Index, 6)

				// Call C implementation
				cErr := originToDirectedEdgesC(pentagon, cEdges)

				// Call Go implementation
				goErr := originToDirectedEdges(pentagon, goEdges)

				// Compare error codes
				if cErr != goErr {
					t.Errorf("Error mismatch for pentagon 0x%x (baseCell=%d, res=%d): C=%v, Go=%v",
						pentagon, baseCell, res, cErr, goErr)
					return
				}

				// If successful, compare all edge results
				if cErr == E_SUCCESS {
					for i := 0; i < 6; i++ {
						if cEdges[i] != goEdges[i] {
							t.Errorf("Edge[%d] mismatch for pentagon 0x%x (baseCell=%d, res=%d): C=0x%x, Go=0x%x",
								i, pentagon, baseCell, res, cEdges[i], goEdges[i])
						}
					}

					// For pentagons, edge[0] should be H3_NULL (K-axis direction)
					if goEdges[0] != H3_NULL {
						t.Errorf("Pentagon edge[0] should be H3_NULL for pentagon 0x%x (baseCell=%d, res=%d), got 0x%x",
							pentagon, baseCell, res, goEdges[0])
					}

					// Check that other edges are properly formed directed edges
					for i := 1; i < 6; i++ {
						if goEdges[i] == H3_NULL {
							t.Errorf("Pentagon edge[%d] should not be H3_NULL for pentagon 0x%x (baseCell=%d, res=%d)",
								i, pentagon, baseCell, res)
						}
						// Check that the mode is correct
						if getMode(goEdges[i]) != H3_DIRECTEDEDGE_MODE {
							t.Errorf("Pentagon edge[%d] should have DIRECTEDEDGE mode for pentagon 0x%x (baseCell=%d, res=%d), got mode %d",
								i, pentagon, baseCell, res, getMode(goEdges[i]))
						}
					}
				}
			})
		}
	}
}

func Test_originToDirectedEdges_exact_buffer_parity(t *testing.T) {
	origin := H3Index(0x8001fffffffffff)

	// Test with exact buffer size (6 elements as expected by C function)
	t.Run("exact_buffer", func(t *testing.T) {
		cEdges := make([]H3Index, 6)
		goEdges := make([]H3Index, 6)

		cErr := originToDirectedEdgesC(origin, cEdges)
		goErr := originToDirectedEdges(origin, goEdges)

		if cErr != goErr {
			t.Errorf("Error mismatch for exact buffer: C=%v, Go=%v", cErr, goErr)
		}

		if cErr == E_SUCCESS {
			for i := 0; i < 6; i++ {
				if cEdges[i] != goEdges[i] {
					t.Errorf("Edge[%d] mismatch: C=0x%x, Go=0x%x", i, cEdges[i], goEdges[i])
				}
			}
		}
	})
}
