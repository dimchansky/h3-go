//go:build cgo && c2go

package h3

import (
	"fmt"
	"testing"
)

func Test_h3Rotate60ccw_parity(t *testing.T) {
	tests := []struct {
		name string
		h    H3Index
	}{
		{"res_0", 0x8001fffffffffff},
		{"res_5", 0x850dab63fffffff},
		{"res_10", 0x8a1fb46622dffff},
		{"pentagon_res_1", 0x81083ffffffffff},
		{"pentagon_res_9", 0x89283082837ffff},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the Go implementation
			goResult := _h3Rotate60ccw(tt.h)
			// Test the C implementation
			cResult := h3Rotate60ccwC(tt.h)

			if goResult != cResult {
				t.Errorf("_h3Rotate60ccw() parity mismatch: got %016x, want %016x", goResult, cResult)
			}
		})
	}
}

func Test_h3Rotate60ccw_all_resolutions(t *testing.T) {
	// Test a representative index at each resolution
	testIndices := []H3Index{
		0x8001fffffffffff, // res 0
		0x81083ffffffffff, // res 1
		0x820807fffffffff, // res 2
		0x83080dfffffffff, // res 3
		0x840820fffffffff, // res 4
		0x850dab63fffffff, // res 5
		0x862834072ffffff, // res 6
		0x87283472bffffff, // res 7
		0x882834722ffffff, // res 8
		0x89283082837ffff, // res 9
		0x8a1fb46622dffff, // res 10
		0x8b1fb46622d7fff, // res 11
		0x8c1fb46622d6fff, // res 12
		0x8d1fb46622d69ff, // res 13
		0x8e1fb46622d691f, // res 14
		0x8f1fb46622d6919, // res 15
	}

	for i, h := range testIndices {
		t.Run(fmt.Sprintf("resolution_%d", i), func(t *testing.T) {
			goResult := _h3Rotate60ccw(h)
			cResult := h3Rotate60ccwC(h)

			if goResult != cResult {
				t.Errorf("_h3Rotate60ccw() parity mismatch at resolution %d: got %016x, want %016x", i, goResult, cResult)
			}
		})
	}
}
