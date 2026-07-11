//go:build cgo && c2go

package h3

import (
	"fmt"
	"math"
	"testing"
)

func Test_edgeLengthKm_parity(t *testing.T) {
	tests := []struct {
		name string
		edge h3Index
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

	const tolerance = 5e-11 // Tolerance for kilometers (accounts for conversion from radians)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Get C implementation result
			cLength, cErr := edgeLengthKmC(tt.edge)

			// Get Go implementation result
			var goLength float64
			goErr := edgeLengthKm(tt.edge, &goLength)

			// Compare errors
			if cErr != goErr {
				t.Errorf("Error mismatch: C=%v, Go=%v", cErr, goErr)
				return
			}

			// If there was an error, we're done
			if cErr != eSuccess {
				return
			}

			// Compare lengths with tolerance appropriate for km units
			if math.Abs(cLength-goLength) > tolerance {
				t.Errorf("Length mismatch for edge %016x: C=%.10e, Go=%.10e, diff=%.2e",
					tt.edge, cLength, goLength, math.Abs(cLength-goLength))
			}

			// Sanity check: edge length should be positive
			if goLength <= 0 {
				t.Errorf("Edge length should be positive, got %.10e", goLength)
			}

			// Sanity check: edge length should be reasonable (less than half Earth circumference in km)
			maxReasonableLength := math.Pi * earthRadiusKm // ~20,000 km
			if goLength >= maxReasonableLength {
				t.Errorf("Edge length %.10e seems unreasonably large (>= %.10e km)", goLength, maxReasonableLength)
			}
		})
	}
}

func Test_edgeLengthKm_invalidEdges_parity(t *testing.T) {
	tests := []struct {
		name string
		edge h3Index
	}{
		{"invalid_mode", 0x8001fffffffffff},           // cell mode instead of edge
		{"zero_edge", 0x0},                            // zero
		{"invalid_edge_high_bits", 0xfffffffffffffff}, // all bits set
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Get C implementation result
			cLength, cErr := edgeLengthKmC(tt.edge)

			// Get Go implementation result
			var goLength float64
			goErr := edgeLengthKm(tt.edge, &goLength)

			// For invalid edges, we primarily care that error codes match
			if goErr != cErr {
				t.Errorf("Error code mismatch for invalid edge: Go=%d, C=%d", goErr, cErr)
			}

			// If both returned success (shouldn't happen for these invalid cases), check lengths match
			if goErr == eSuccess && cErr == eSuccess {
				if math.Abs(cLength-goLength) > 5e-11 {
					t.Errorf("Length mismatch for invalid edge: C=%.10e, Go=%.10e", cLength, goLength)
				}
			}
		})
	}
}

func Test_edgeLengthKm_comprehensive_parity(t *testing.T) {
	// Test various directed edges created from valid cells
	testCells := []h3Index{
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
			var edges [6]h3Index
			err := originToDirectedEdges(cell, edges[:])
			if err != eSuccess {
				t.Fatalf("Failed to get directed edges: %v", err)
			}

			// Test edge length for each valid edge
			for i, edge := range edges {
				if edge == h3Null {
					continue // Skip null edges (pentagon case)
				}

				// Get C implementation result
				cLength, cErr := edgeLengthKmC(edge)

				// Get Go implementation result
				var goLength float64
				goErr := edgeLengthKm(edge, &goLength)

				if goErr != cErr {
					t.Errorf("Edge %d error mismatch: Go=%d, C=%d", i, goErr, cErr)
					continue
				}

				if goErr != eSuccess {
					continue
				}

				// Check length precision for km units
				const tolerance = 5e-11
				if math.Abs(cLength-goLength) > tolerance {
					t.Errorf("Edge %d length mismatch: C=%.10e, Go=%.10e, diff=%.2e",
						i, cLength, goLength, math.Abs(cLength-goLength))
				}

				// Sanity checks
				if goLength <= 0 {
					t.Errorf("Edge %d length should be positive, got %.10e", i, goLength)
				}
			}
		})
	}
}
