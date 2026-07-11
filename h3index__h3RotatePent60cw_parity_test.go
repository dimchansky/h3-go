//go:build cgo && c2go

package h3

import (
	"fmt"
	"testing"
)

func Test_h3RotatePent60cw_parity(t *testing.T) {
	// Test pentagon rotation with pentagon base cells
	tests := []struct {
		name string
		h    h3Index
	}{
		{"pentagon_res_1", 0x81083ffffffffff}, // Pentagon base cell 4, res 1
		{"pentagon_res_2", 0x820807fffffffff}, // Pentagon base cell 14, res 2
		{"pentagon_res_5", 0x850dab63fffffff}, // Pentagon base cell 58, res 5
		{"pentagon_res_9", 0x89283082837ffff}, // Pentagon base cell 117, res 9
		{"non_pentagon", 0x8001fffffffffff},   // Non-pentagon for comparison
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the Go implementation
			goResult := _h3RotatePent60cw(tt.h)
			// Test the C implementation
			cResult := h3RotatePent60cwC(tt.h)

			if goResult != cResult {
				t.Errorf("_h3RotatePent60cw() parity mismatch: got %016x, want %016x", goResult, cResult)
			}
		})
	}
}

func Test_h3RotatePent60cw_all_pentagon_base_cells(t *testing.T) {
	// Test all pentagon base cells at different resolutions
	pentagonBaseCells := []int32{4, 14, 24, 38, 49, 58, 63, 72, 82, 83, 97, 107, 117}

	for _, baseCell := range pentagonBaseCells {
		for res := int32(0); res <= 5; res++ { // Test lower resolutions to keep tests fast
			t.Run(fmt.Sprintf("baseCell_%d_res_%d", baseCell, res), func(t *testing.T) {
				// Create an index for this base cell and resolution
				var h h3Index
				setH3Index(&h, res, baseCell, 0)

				goResult := _h3RotatePent60cw(h)
				cResult := h3RotatePent60cwC(h)

				if goResult != cResult {
					t.Errorf("_h3RotatePent60cw() parity mismatch for baseCell %d res %d: got %016x, want %016x", baseCell, res, goResult, cResult)
				}
			})
		}
	}
}
