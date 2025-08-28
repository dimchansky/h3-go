//go:build cgo

package h3

import (
	"testing"
)

func Test_cellAreaM2_parity(t *testing.T) {
	tests := []struct {
		name string
		cell H3Index
	}{
		// Valid hexagon cells at different resolutions
		{"res0_base", 0x8001fffffffffff},
		{"res1_hex", 0x8101ffffffffffff},
		{"res2_hex", 0x8201ffffffffffff},
		{"res3_hex", 0x8301ffffffffffff},
		{"res5_hex", 0x8501ffffffffffff},
		{"res7_hex", 0x8701ffffffffffff},
		{"res10_hex", 0x8a01ffffffffffff},
		{"res15_hex", 0x8f01ffffffffffff},

		// Pentagon cells
		{"res0_pentagon", 0x8004fffffffffff}, // Base cell 4 is pentagon
		{"res1_pentagon", 0x81047ffffffffff},
		{"res5_pentagon", 0x85047ffffffffff},

		// More varied test cases across different base cells
		{"different_base1", 0x8009fffffffffff},
		{"different_base2", 0x800dfffffffffff},
		{"different_base3", 0x8019fffffffffff},
	}

	const tolerance = 1e-1 // Reasonable tolerance for m^2 calculations (larger values than km^2)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Get C implementation result
			cArea, cErr := cellAreaM2C(tt.cell)

			// Get Go implementation result
			goArea, goErr := cellAreaM2(tt.cell)

			// Compare errors
			if cErr != goErr {
				t.Errorf("Error mismatch: C=%v, Go=%v", cErr, goErr)
				return
			}

			// If there was an error, we're done
			if cErr != E_SUCCESS {
				return
			}

			// Compare areas with appropriate tolerance
			if !floatAlmostEqual(cArea, goArea, tolerance) {
				t.Errorf("Area mismatch for cell %016x: C=%.9e, Go=%.9e, diff=%.2e",
					tt.cell, cArea, goArea, abs(cArea-goArea))
			}
		})
	}
}
