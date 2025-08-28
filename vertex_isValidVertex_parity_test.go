//go:build cgo

package h3

import (
	"testing"
)

func Test_isValidVertex_parity(t *testing.T) {
	tests := []struct {
		name   string
		vertex H3Index
	}{
		// Valid vertices from various resolutions and cell types
		{"valid hex vertex r0", 0x20283080bffffff},
		{"valid hex vertex r1", 0x21283080bffffff},
		{"valid hex vertex r2", 0x22283080bffffff},
		{"valid hex vertex r5", 0x25283080bffffff},
		{"valid pent vertex r0", 0x20080800fffffff},
		{"valid pent vertex r1", 0x21080800fffffff},
		{"valid pent vertex r2", 0x22080800fffffff},
		{"valid pent vertex r5", 0x25080800fffffff},

		// Invalid vertices - wrong mode
		{"cell mode instead of vertex", 0x10283080bffffff},
		{"edge mode instead of vertex", 0x30283080bffffff},

		// Invalid vertices - bad reserved bits (vertex numbers)
		{"invalid vertex number", 0x20283080bfffffff},

		// Null vertex
		{"null vertex", H3_NULL},

		// Invalid owner cells
		{"invalid owner cell", 0x2fffffffffffffff},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Get results from both implementations
			goResult := isValidVertex(tt.vertex)
			cResult := isValidVertexC(tt.vertex)

			// Compare results
			if goResult != cResult {
				t.Errorf("Result mismatch for vertex 0x%x: Go=%t, C=%t",
					tt.vertex, goResult, cResult)
			}
		})
	}
}

func Test_isValidVertex_comprehensive_parity(t *testing.T) {
	// Test a wider range of vertices by creating them from known cells
	// Note: Not all generated vertices will be canonical/valid, but we check for parity
	testCells := []H3Index{
		0x8001fffffffffff, // resolution 0 hex
		0x8008001ffffffff, // resolution 1 hex
		0x800c001ffffffff, // resolution 2 hex
		0x8009001ffffffff, // resolution 1 pentagon
		0x800d001ffffffff, // resolution 2 pentagon
	}

	for _, cell := range testCells {
		// Test all possible vertex numbers for this cell
		maxVertexNum := 6
		if isPentagon(cell) {
			maxVertexNum = 5
		}

		for vertexNum := 0; vertexNum < maxVertexNum; vertexNum++ {
			var vertex H3Index
			err := cellToVertex(cell, int32(vertexNum), &vertex)
			if err != E_SUCCESS {
				continue // Skip invalid combinations
			}

			goResult := isValidVertex(vertex)
			cResult := isValidVertexC(vertex)

			if goResult != cResult {
				t.Errorf("Result mismatch for vertex 0x%x (from cell 0x%x, vertexNum %d): Go=%t, C=%t",
					vertex, cell, vertexNum, goResult, cResult)
			}
		}
	}
}
