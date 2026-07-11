//go:build cgo && c2go

package h3

import (
	"testing"
)

func Test_getDirectedEdgeOrigin_parity(t *testing.T) {
	tests := []struct {
		name string
		edge h3Index
	}{
		// Valid directed edge cases - constructed from valid cell indices
		// Mode=2 (directed edge), with direction in reserved bits (1-6)
		{"valid_edge_dir1", h3Index(0x2001fffffffffff)}, // mode=2, direction=1
		{"valid_edge_dir2", h3Index(0x2002fffffffffff)}, // mode=2, direction=2
		{"valid_edge_dir3", h3Index(0x2003fffffffffff)}, // mode=2, direction=3
		{"valid_edge_dir6", h3Index(0x2006fffffffffff)}, // mode=2, direction=6

		// Edge with different resolution patterns
		{"valid_edge_res5", h3Index(0x2501ffffffffff)},     // mode=2, direction=1
		{"valid_edge_simple", h3Index(0x2101000000000000)}, // mode=2, direction=1

		// Invalid cases
		{"invalid_mode_cell", h3Index(0x1001fffffffffff)},  // mode=1 (cell), should fail
		{"invalid_mode_0", h3Index(0x0001fffffffffff)},     // mode=0, should fail
		{"invalid_mode_3", h3Index(0x3001fffffffffff)},     // mode=3, should fail
		{"invalid_reserved_0", h3Index(0x2000fffffffffff)}, // mode=2 but direction=0, should fail
		{"invalid_reserved_7", h3Index(0x2007fffffffffff)}, // mode=2 but direction=7, should fail

		// Edge cases
		{"zero_index", h3Index(0x0)},
		{"max_uint64", h3Index(0xffffffffffffffff)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Get C implementation result
			cOrigin, cErr := getDirectedEdgeOriginC(tt.edge)

			// Get Go implementation result
			goOrigin, goErr := getDirectedEdgeOrigin(tt.edge)

			// Compare errors
			if cErr != goErr {
				t.Errorf("Error mismatch: C=%v, Go=%v", cErr, goErr)
				return
			}

			// If there was an error, we're done
			if cErr != eSuccess {
				return
			}

			// Compare origins (should be identical)
			if cOrigin != goOrigin {
				t.Errorf("Origin mismatch: C=0x%x, Go=0x%x", cOrigin, goOrigin)
			}
		})
	}
}

func Test_getDirectedEdgeOrigin_constructed_edges_parity(t *testing.T) {
	// Test with some constructed directed edges based on known valid cells
	validCells := []h3Index{
		h3Index(0x8001fffffffffff), // Simple res=0 cell
		h3Index(0x8101fffffffffff), // Simple res=1 cell
		h3Index(0x8201fffffffffff), // Simple res=2 cell
		h3Index(0x8301fffffffffff), // Simple res=3 cell
	}

	// Valid directions for directed edges (1-6)
	validDirections := []int32{1, 2, 3, 4, 5, 6}

	for _, cell := range validCells {
		for _, direction := range validDirections {
			// Construct a directed edge from this cell
			edge := setMode(cell, h3DirectededgeMode)
			edge = setReservedBits(edge, direction)

			t.Run("constructed_edge", func(t *testing.T) {
				// Get C implementation result
				cOrigin, cErr := getDirectedEdgeOriginC(edge)

				// Get Go implementation result
				goOrigin, goErr := getDirectedEdgeOrigin(edge)

				// Compare errors
				if cErr != goErr {
					t.Errorf("Error mismatch for edge 0x%x: C=%v, Go=%v", edge, cErr, goErr)
					return
				}

				// If there was an error, we're done
				if cErr != eSuccess {
					return
				}

				// Compare origins
				if cOrigin != goOrigin {
					t.Errorf("Origin mismatch for edge 0x%x: C=0x%x, Go=0x%x", edge, cOrigin, goOrigin)
				}

				// The origin should be the original cell (with cell mode and cleared reserved bits)
				expectedOrigin := setMode(cell, h3CellMode)
				expectedOrigin = setReservedBits(expectedOrigin, 0)

				if goOrigin != expectedOrigin {
					t.Errorf("Origin doesn't match expected cell: got=0x%x, expected=0x%x", goOrigin, expectedOrigin)
				}
			})
		}
	}
}
