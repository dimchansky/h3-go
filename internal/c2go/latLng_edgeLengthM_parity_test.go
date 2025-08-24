//go:build cgo

package c2go

import (
	"fmt"
	"math"
	"testing"
)

func Test_edgeLengthM_parity(t *testing.T) {
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

		// Edges at different resolutions for length validation
		{"hex_edge_r3", 0x1401e001ffffffff},
		{"hex_edge_r4", 0x1502f001ffffffff},
		{"hex_edge_r5", 0x1603f001ffffffff},
	}

	const tolerance = 5e-8 // Tolerance for meters (accounts for conversion from km)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Get C implementation result
			cLength, cErr := edgeLengthMC(tt.edge)

			// Get Go implementation result
			var goLength float64
			goErr := edgeLengthM(tt.edge, &goLength)

			// Compare errors
			if cErr != goErr {
				t.Errorf("Error mismatch: C=%v, Go=%v", cErr, goErr)
				return
			}

			// If there was an error, we're done
			if cErr != E_SUCCESS {
				return
			}

			// Compare lengths with tolerance appropriate for meter units
			if math.Abs(cLength-goLength) > tolerance {
				t.Errorf("Length mismatch for edge %016x: C=%.7e, Go=%.7e, diff=%.2e",
					tt.edge, cLength, goLength, math.Abs(cLength-goLength))
			}

			// Sanity check: edge length should be positive
			if goLength <= 0 {
				t.Errorf("Edge length should be positive, got %.7e", goLength)
			}

			// Sanity check: edge length should be reasonable (less than half Earth circumference in m)
			maxReasonableLength := math.Pi * EARTH_RADIUS_KM * 1000 // ~20,000,000 m
			if goLength >= maxReasonableLength {
				t.Errorf("Edge length %.7e seems unreasonably large (>= %.7e m)", goLength, maxReasonableLength)
			}
		})
	}
}

func Test_edgeLengthM_invalidEdges_parity(t *testing.T) {
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
			// Get C implementation result
			cLength, cErr := edgeLengthMC(tt.edge)

			// Get Go implementation result
			var goLength float64
			goErr := edgeLengthM(tt.edge, &goLength)

			// For invalid edges, we primarily care that error codes match
			if goErr != cErr {
				t.Errorf("Error code mismatch for invalid edge: Go=%d, C=%d", goErr, cErr)
			}

			// If both returned success (shouldn't happen for these invalid cases), check lengths match
			if goErr == E_SUCCESS && cErr == E_SUCCESS {
				if math.Abs(cLength-goLength) > 5e-8 {
					t.Errorf("Length mismatch for invalid edge: C=%.7e, Go=%.7e", cLength, goLength)
				}
			}
		})
	}
}

func Test_edgeLengthM_comprehensive_parity(t *testing.T) {
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

			// Test edge length for each valid edge
			for i, edge := range edges {
				if edge == H3_NULL {
					continue // Skip null edges (pentagon case)
				}

				// Get C implementation result
				cLength, cErr := edgeLengthMC(edge)

				// Get Go implementation result
				var goLength float64
				goErr := edgeLengthM(edge, &goLength)

				if goErr != cErr {
					t.Errorf("Edge %d error mismatch: Go=%d, C=%d", i, goErr, cErr)
					continue
				}

				if goErr != E_SUCCESS {
					continue
				}

				// Check length precision for meter units
				const tolerance = 5e-8
				if math.Abs(cLength-goLength) > tolerance {
					t.Errorf("Edge %d length mismatch: C=%.7e, Go=%.7e, diff=%.2e",
						i, cLength, goLength, math.Abs(cLength-goLength))
				}

				// Sanity checks
				if goLength <= 0 {
					t.Errorf("Edge %d length should be positive, got %.7e", i, goLength)
				}
			}
		})
	}
}
