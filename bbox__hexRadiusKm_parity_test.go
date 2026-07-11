//go:build cgo && c2go

package h3

import (
	"math"
	"testing"
)

func Test_hexRadiusKm_parity(t *testing.T) {
	tests := []struct {
		name    string
		h3Index h3Index
	}{
		// Valid indices at different resolutions
		{"res0_index", 0x8001fffffffffff},
		{"res1_index", 0x8101000000000ff},
		{"res2_index", 0x820100000000fff},
		{"res5_index", 0x85283447fffffff},
		{"res10_index", 0x8a283447c3fffff},
		{"res15_index", 0x8f283447c3c3c3f},
		// Pentagon test cases
		{"res5_pentagon", 0x851c0003fffffff}, // Base cell 4 (pentagon) at res 5
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !isValidCell(tt.h3Index) {
				t.Skipf("Invalid test index: %x", tt.h3Index)
			}

			goResult := _hexRadiusKm(tt.h3Index)
			cResult := _hexRadiusKmC(tt.h3Index)

			// Use tight tolerance for radius calculations
			tolerance := 1e-9
			if math.Abs(goResult-cResult) > tolerance {
				t.Errorf("_hexRadiusKm(%x): Go=%f, C=%f, diff=%f > tolerance %f",
					tt.h3Index, goResult, cResult, math.Abs(goResult-cResult), tolerance)
			}

			// Sanity check - radius should be positive and reasonable
			if goResult <= 0 {
				t.Errorf("_hexRadiusKm(%x): radius should be positive, got %f",
					tt.h3Index, goResult)
			}
		})
	}
}
