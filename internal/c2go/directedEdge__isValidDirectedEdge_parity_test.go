//go:build cgo

package c2go

import (
	"testing"
)

func Test_isValidDirectedEdge_parity(t *testing.T) {
	tests := []struct {
		name string
		edge H3Index
	}{
		// Valid directed edge cases - constructed from valid cell indices
		// Mode=2 (directed edge), with direction in reserved bits (1-6)
		{"valid_edge_dir1", H3Index(0x2001fffffffffff)}, // mode=2, direction=1
		{"valid_edge_dir2", H3Index(0x2002fffffffffff)}, // mode=2, direction=2
		{"valid_edge_dir3", H3Index(0x2003fffffffffff)}, // mode=2, direction=3
		{"valid_edge_dir4", H3Index(0x2004fffffffffff)}, // mode=2, direction=4
		{"valid_edge_dir5", H3Index(0x2005fffffffffff)}, // mode=2, direction=5
		{"valid_edge_dir6", H3Index(0x2006fffffffffff)}, // mode=2, direction=6

		// Edge with different resolution patterns
		{"valid_edge_res5", H3Index(0x2501ffffffffff)},     // mode=2, direction=1
		{"valid_edge_simple", H3Index(0x2101000000000000)}, // mode=2, direction=1

		// Invalid cases - wrong mode
		{"invalid_mode_cell", H3Index(0x1001fffffffffff)}, // mode=1 (cell), should fail
		{"invalid_mode_0", H3Index(0x0001fffffffffff)},    // mode=0, should fail
		{"invalid_mode_3", H3Index(0x3001fffffffffff)},    // mode=3, should fail

		// Invalid cases - bad direction
		{"invalid_reserved_0", H3Index(0x2000fffffffffff)}, // mode=2 but direction=0 (CENTER_DIGIT), should fail
		{"invalid_reserved_7", H3Index(0x2007fffffffffff)}, // mode=2 but direction=7 (INVALID_DIGIT), should fail

		// Edge cases
		{"zero_index", H3Index(0x0)},
		{"max_uint64", H3Index(0xffffffffffffffff)},

		// Pentagon cases - K_AXES_DIGIT (1) should be invalid for pentagons
		{"pentagon_k_axis", H3Index(0x2001000000000000)}, // should be invalid if underlying is pentagon
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Get C implementation result
			cResult := isValidDirectedEdgeC(tt.edge)

			// Get Go implementation result
			goResult := isValidDirectedEdge(tt.edge)

			// Compare results
			if cResult != goResult {
				t.Errorf("Result mismatch for edge 0x%x: C=%v, Go=%v", tt.edge, cResult, goResult)
			}
		})
	}
}

func Test_isValidDirectedEdge_constructed_edges_parity(t *testing.T) {
	// Test with some constructed directed edges based on known valid cells
	validCells := []H3Index{
		H3Index(0x8001fffffffffff), // Simple res=0 cell
		H3Index(0x8101fffffffffff), // Simple res=1 cell
		H3Index(0x8201fffffffffff), // Simple res=2 cell
		H3Index(0x8301fffffffffff), // Simple res=3 cell
	}

	// Valid directions for directed edges (1-6)
	validDirections := []int32{1, 2, 3, 4, 5, 6}

	for _, cell := range validCells {
		for _, direction := range validDirections {
			// Construct a directed edge from this cell
			edge := setMode(cell, H3_DIRECTEDEDGE_MODE)
			edge = setReservedBits(edge, direction)

			t.Run("constructed_edge", func(t *testing.T) {
				// Get C implementation result
				cResult := isValidDirectedEdgeC(edge)

				// Get Go implementation result
				goResult := isValidDirectedEdge(edge)

				// Compare results
				if cResult != goResult {
					t.Errorf("Result mismatch for edge 0x%x: C=%v, Go=%v", edge, cResult, goResult)
				}
			})
		}
	}
}

func Test_isValidDirectedEdge_pentagon_edges_parity(t *testing.T) {
	// Test with pentagon base cells and different directions
	// Pentagon base cells are: 4, 14, 24, 38, 49, 58, 63, 72, 83, 97, 107, 117
	pentagonBaseCells := []int32{4, 14, 24, 38, 49, 58}

	for _, baseCell := range pentagonBaseCells {
		// Create a simple pentagon cell at resolution 1
		pentagon := setMode(H3Index(0), H3_CELL_MODE)
		pentagon = H3Index(setBaseCell(uint64(pentagon), baseCell))
		pentagon = setResolutionForTest(pentagon, 1)

		// Test all directions
		for direction := int32(1); direction <= 6; direction++ {
			edge := setMode(pentagon, H3_DIRECTEDEDGE_MODE)
			edge = setReservedBits(edge, direction)

			t.Run("pentagon_edge", func(t *testing.T) {
				// Get C implementation result
				cResult := isValidDirectedEdgeC(edge)

				// Get Go implementation result
				goResult := isValidDirectedEdge(edge)

				// Compare results
				if cResult != goResult {
					t.Errorf("Result mismatch for pentagon edge 0x%x (baseCell=%d, dir=%d): C=%v, Go=%v",
						edge, baseCell, direction, cResult, goResult)
				}

				// K_AXES_DIGIT (1) should be invalid for pentagons
				if direction == int32(K_AXES_DIGIT) && goResult {
					t.Errorf("Pentagon edge with K_AXES_DIGIT should be invalid: edge=0x%x, baseCell=%d",
						edge, baseCell)
				}
			})
		}
	}
}
