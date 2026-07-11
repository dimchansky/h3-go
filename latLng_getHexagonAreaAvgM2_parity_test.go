//go:build cgo && c2go

package h3

import (
	"testing"
)

func Test_getHexagonAreaAvgM2_parity(t *testing.T) {
	tests := []struct {
		name string
		res  int32
	}{
		{"resolution 0", 0},
		{"resolution 1", 1},
		{"resolution 2", 2},
		{"resolution 3", 3},
		{"resolution 4", 4},
		{"resolution 5", 5},
		{"resolution 6", 6},
		{"resolution 7", 7},
		{"resolution 8", 8},
		{"resolution 9", 9},
		{"resolution 10", 10},
		{"resolution 11", 11},
		{"resolution 12", 12},
		{"resolution 13", 13},
		{"resolution 14", 14},
		{"resolution 15", 15},
		{"invalid resolution -1", -1},
		{"invalid resolution 16", 16},
		{"invalid resolution 100", 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var goOut float64
			goErr := getHexagonAreaAvgM2(tt.res, &goOut)

			cOut, cErr := getHexagonAreaAvgM2C(tt.res)

			if goErr != cErr {
				t.Errorf("getHexagonAreaAvgM2(%d) error mismatch: Go=%v, C=%v", tt.res, goErr, cErr)
				return
			}

			if goErr == E_SUCCESS {
				if goOut != cOut {
					t.Errorf("getHexagonAreaAvgM2(%d) output mismatch: Go=%v, C=%v", tt.res, goOut, cOut)
				}
			}
		})
	}
}
