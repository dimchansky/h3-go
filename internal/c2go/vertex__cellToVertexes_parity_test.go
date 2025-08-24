//go:build c2go

package c2go

import (
	"testing"
)

func Test_cellToVertexes_parity(t *testing.T) {
	testCases := []struct {
		name string
		cell H3Index
	}{
		{
			name: "hex cell",
			cell: 0x8930062838bffff,
		},
		{
			name: "pentagon cell",
			cell: 0x893006283807fff,
		},
		{
			name: "resolution 0 cell",
			cell: 0x8001fffffffffff,
		},
		{
			name: "resolution 0 pentagon",
			cell: 0x8007fffffffffff,
		},
		{
			name: "high resolution hex cell",
			cell: 0x8a30062838b7776,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call Go implementation
			var goVertexes [6]H3Index
			goErr := _cellToVertexes(tc.cell, &goVertexes)

			// Call C implementation
			var cVertexes [6]H3Index
			cErr := cellToVertexesC(tc.cell, &cVertexes)

			// Compare error codes
			if goErr != cErr {
				t.Errorf("_cellToVertexes(%#016x) error = %v, want %v",
					tc.cell, goErr, cErr)
				return
			}

			// If no error, compare outputs
			if goErr == E_SUCCESS {
				for i := 0; i < 6; i++ {
					if goVertexes[i] != cVertexes[i] {
						t.Errorf("_cellToVertexes(%#016x)[%d] = %#016x, want %#016x",
							tc.cell, i, goVertexes[i], cVertexes[i])
					}
				}
			}
		})
	}
}

func Test_cellToVertexes_extensive_parity(t *testing.T) {
	// Test a variety of cells across different resolutions
	testCells := []H3Index{
		0x8001fffffffffff, // Resolution 0
		0x8107fffffffffff, // Resolution 1
		0x8230062ffffffff, // Resolution 2
		0x8330062838fffff, // Resolution 3
		0x8430062838bffff, // Resolution 4
		0x8530062838b7fff, // Resolution 5
		0x8630062838b77ff, // Resolution 6
		0x8730062838b777f, // Resolution 7
		0x8830062838b7777, // Resolution 8
		0x8930062838b7776, // Resolution 9
		0x8a30062838b7776, // Resolution 10
		// Pentagon cells
		0x8007fffffffffff,  // Pentagon at resolution 0
		0x8107ffffffffffff, // Pentagon at resolution 1
		0x8230007ffffffff,  // Pentagon at resolution 2
		0x8330007fffffff,   // Pentagon at resolution 3
		0x8430007ffffffff,  // Pentagon at resolution 4
	}

	for _, cell := range testCells {
		// Call Go implementation
		var goVertexes [6]H3Index
		goErr := _cellToVertexes(cell, &goVertexes)

		// Call C implementation
		var cVertexes [6]H3Index
		cErr := cellToVertexesC(cell, &cVertexes)

		// Compare error codes
		if goErr != cErr {
			t.Errorf("_cellToVertexes(%#016x) error = %v, want %v",
				cell, goErr, cErr)
			continue
		}

		// If no error, compare outputs
		if goErr == E_SUCCESS {
			for i := 0; i < 6; i++ {
				if goVertexes[i] != cVertexes[i] {
					t.Errorf("_cellToVertexes(%#016x)[%d] = %#016x, want %#016x",
						cell, i, goVertexes[i], cVertexes[i])
				}
			}
		}
	}
}
