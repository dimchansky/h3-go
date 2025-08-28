//go:build c2go

package h3

import (
	"testing"
)

func Test_cellToVertex_parity(t *testing.T) {
	testCases := []struct {
		name      string
		cell      H3Index
		vertexNum int32
	}{
		{
			name:      "hex cell vertex 0",
			cell:      0x8930062838bffff,
			vertexNum: 0,
		},
		{
			name:      "hex cell vertex 1",
			cell:      0x8930062838bffff,
			vertexNum: 1,
		},
		{
			name:      "hex cell vertex 2",
			cell:      0x8930062838bffff,
			vertexNum: 2,
		},
		{
			name:      "hex cell vertex 3",
			cell:      0x8930062838bffff,
			vertexNum: 3,
		},
		{
			name:      "hex cell vertex 4",
			cell:      0x8930062838bffff,
			vertexNum: 4,
		},
		{
			name:      "hex cell vertex 5",
			cell:      0x8930062838bffff,
			vertexNum: 5,
		},
		{
			name:      "pentagon cell vertex 0",
			cell:      0x893006283807fff,
			vertexNum: 0,
		},
		{
			name:      "pentagon cell vertex 1",
			cell:      0x893006283807fff,
			vertexNum: 1,
		},
		{
			name:      "pentagon cell vertex 2",
			cell:      0x893006283807fff,
			vertexNum: 2,
		},
		{
			name:      "pentagon cell vertex 3",
			cell:      0x893006283807fff,
			vertexNum: 3,
		},
		{
			name:      "pentagon cell vertex 4",
			cell:      0x893006283807fff,
			vertexNum: 4,
		},
		{
			name:      "resolution 0 cell vertex 0",
			cell:      0x8001fffffffffff,
			vertexNum: 0,
		},
		{
			name:      "resolution 0 cell vertex 5",
			cell:      0x8001fffffffffff,
			vertexNum: 5,
		},
		{
			name:      "invalid vertex number -1",
			cell:      0x8930062838bffff,
			vertexNum: -1,
		},
		{
			name:      "invalid vertex number 6",
			cell:      0x8930062838bffff,
			vertexNum: 6,
		},
		{
			name:      "pentagon invalid vertex number 5",
			cell:      0x893006283807fff,
			vertexNum: 5,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call Go implementation
			var goOut H3Index
			goErr := _cellToVertex(tc.cell, tc.vertexNum, &goOut)

			// Call C implementation
			cOut, cErr := cellToVertexC(tc.cell, tc.vertexNum)

			// Compare error codes
			if goErr != cErr {
				t.Errorf("_cellToVertex(%#016x, %d) error = %v, want %v",
					tc.cell, tc.vertexNum, goErr, cErr)
				return
			}

			// If no error, compare outputs
			if goErr == E_SUCCESS {
				if goOut != cOut {
					t.Errorf("_cellToVertex(%#016x, %d) = %#016x, want %#016x",
						tc.cell, tc.vertexNum, goOut, cOut)
				}
			}
		})
	}
}

func Test_cellToVertex_extensive_parity(t *testing.T) {
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
		// Pentagon cells
		0x8007fffffffffff,  // Pentagon at resolution 0
		0x8107ffffffffffff, // Pentagon at resolution 1
		0x8230007ffffffff,  // Pentagon at resolution 2
	}

	for _, cell := range testCells {
		isPent := isPentagon(cell)
		maxVertex := int32(NUM_HEX_VERTS)
		if isPent {
			maxVertex = NUM_PENT_VERTS
		}

		for vertexNum := int32(0); vertexNum < maxVertex; vertexNum++ {
			// Call Go implementation
			var goOut H3Index
			goErr := _cellToVertex(cell, vertexNum, &goOut)

			// Call C implementation
			cOut, cErr := cellToVertexC(cell, vertexNum)

			// Compare error codes
			if goErr != cErr {
				t.Errorf("_cellToVertex(%#016x, %d) error = %v, want %v",
					cell, vertexNum, goErr, cErr)
				continue
			}

			// If no error, compare outputs
			if goErr == E_SUCCESS {
				if goOut != cOut {
					t.Errorf("_cellToVertex(%#016x, %d) = %#016x, want %#016x",
						cell, vertexNum, goOut, cOut)
				}
			}
		}
	}
}
