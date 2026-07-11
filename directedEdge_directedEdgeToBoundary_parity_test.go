//go:build cgo && c2go

package h3

import (
	"fmt"
	"math"
	"testing"
)

func Test_directedEdgeToBoundary_parity(t *testing.T) {
	tests := []struct {
		name string
		edge H3Index
	}{
		// Valid directed edges at various resolutions
		{"hex_edge_r0", 0x11001fffffffffff},
		{"hex_edge_r1", 0x12008001ffffffff},
		{"hex_edge_r2", 0x1300c001ffffffff},
		{"pent_edge_r0", 0x11009001ffffffff}, // pentagon edge
		{"pent_edge_r1", 0x1200d001ffffffff}, // pentagon edge

		// Edges that might cross icosahedral faces (higher res for more complexity)
		{"complex_edge_r3", 0x1401e001ffffffff},
		{"complex_edge_r4", 0x1502f001ffffffff},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var goBoundary CellBoundary
			var cBoundary CellBoundary

			goErr := directedEdgeToBoundary(tt.edge, &goBoundary)
			cErr := directedEdgeToBoundaryC(tt.edge, &cBoundary)

			// Check error codes match
			if goErr != cErr {
				t.Errorf("Error code mismatch: Go=%d, C=%d", goErr, cErr)
				return
			}

			// If there was an error, no need to check further
			if goErr != E_SUCCESS {
				return
			}

			// Check number of vertices matches
			if goBoundary.NumVerts != cBoundary.NumVerts {
				t.Errorf("NumVerts mismatch: Go=%d, C=%d", goBoundary.NumVerts, cBoundary.NumVerts)
				return
			}

			// Check each vertex coordinate with appropriate tolerance
			const tolerance = 1e-12
			for i := int32(0); i < goBoundary.NumVerts; i++ {
				goLat := goBoundary.Verts[i].Lat.Rad()
				goLng := goBoundary.Verts[i].Lng.Rad()
				cLat := cBoundary.Verts[i].Lat.Rad()
				cLng := cBoundary.Verts[i].Lng.Rad()

				if math.Abs(goLat-cLat) > tolerance {
					t.Errorf("Vertex %d lat mismatch: Go=%.15f, C=%.15f, diff=%.15e",
						i, goLat, cLat, math.Abs(goLat-cLat))
				}

				if math.Abs(goLng-cLng) > tolerance {
					t.Errorf("Vertex %d lng mismatch: Go=%.15f, C=%.15f, diff=%.15e",
						i, goLng, cLng, math.Abs(goLng-cLng))
				}
			}
		})
	}
}

func Test_directedEdgeToBoundary_invalidEdges_parity(t *testing.T) {
	tests := []struct {
		name string
		edge H3Index
	}{
		{"invalid_mode", 0x8001fffffffffff},           // cell mode instead of edge
		{"zero_edge", 0x0},                            // zero
		{"invalid_edge_high_bits", 0xfffffffffffffff}, // all bits set
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var goBoundary CellBoundary
			var cBoundary CellBoundary

			goErr := directedEdgeToBoundary(tt.edge, &goBoundary)
			cErr := directedEdgeToBoundaryC(tt.edge, &cBoundary)

			// For invalid edges, we primarily care that error codes match
			if goErr != cErr {
				t.Errorf("Error code mismatch for invalid edge: Go=%d, C=%d", goErr, cErr)
			}

			// Note: Some invalid edges may return E_SUCCESS with 0 vertices,
			// so we check parity rather than expecting specific error codes
		})
	}
}

func Test_directedEdgeToBoundary_comprehensive_parity(t *testing.T) {
	// Test various directed edges created from valid cells
	testCells := []H3Index{
		0x8001fffffffffff, // resolution 0 hex
		0x8008001ffffffff, // resolution 1 hex
		0x800c001ffffffff, // resolution 2 hex
		0x8009001ffffffff, // resolution 1 pentagon
		0x800d001ffffffff, // resolution 2 pentagon
	}

	for _, cell := range testCells {
		cellName := fmt.Sprintf("cell_0x%x", cell)
		t.Run(cellName, func(t *testing.T) {
			// Get all directed edges for this cell
			var edges [6]H3Index
			err := originToDirectedEdges(cell, edges[:])
			if err != E_SUCCESS {
				t.Fatalf("Failed to get directed edges: %v", err)
			}

			// Test boundary for each valid edge
			for i, edge := range edges {
				if edge == H3_NULL {
					continue // Skip null edges (pentagon case)
				}

				var goBoundary CellBoundary
				var cBoundary CellBoundary

				goErr := directedEdgeToBoundary(edge, &goBoundary)
				cErr := directedEdgeToBoundaryC(edge, &cBoundary)

				if goErr != cErr {
					t.Errorf("Edge %d error mismatch: Go=%d, C=%d", i, goErr, cErr)
					continue
				}

				if goErr != E_SUCCESS {
					continue
				}

				// Check boundary properties
				if goBoundary.NumVerts != cBoundary.NumVerts {
					t.Errorf("Edge %d NumVerts mismatch: Go=%d, C=%d", i, goBoundary.NumVerts, cBoundary.NumVerts)
					continue
				}

				// Directed edges should have 2 vertices (sometimes 3 with distortion)
				if goBoundary.NumVerts < 2 || goBoundary.NumVerts > 3 {
					t.Errorf("Edge %d has unexpected vertex count: %d", i, goBoundary.NumVerts)
				}

				// Check coordinate precision
				const tolerance = 1e-12
				for j := int32(0); j < goBoundary.NumVerts; j++ {
					if math.Abs(goBoundary.Verts[j].Lat.Rad()-cBoundary.Verts[j].Lat.Rad()) > tolerance {
						t.Errorf("Edge %d vertex %d lat mismatch", i, j)
					}
					if math.Abs(goBoundary.Verts[j].Lng.Rad()-cBoundary.Verts[j].Lng.Rad()) > tolerance {
						t.Errorf("Edge %d vertex %d lng mismatch", i, j)
					}
				}
			}
		})
	}
}
