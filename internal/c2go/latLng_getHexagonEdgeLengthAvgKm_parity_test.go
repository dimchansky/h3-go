//go:build cgo

package c2go

import (
	"testing"
)

func Test_getHexagonEdgeLengthAvgKm_parity(t *testing.T) {
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
			goErr := getHexagonEdgeLengthAvgKm(tt.res, &goOut)

			cOut, cErr := getHexagonEdgeLengthAvgKmC(tt.res)

			if goErr != cErr {
				t.Errorf("getHexagonEdgeLengthAvgKm(%d) error mismatch: Go=%v, C=%v", tt.res, goErr, cErr)
				return
			}

			if goErr == E_SUCCESS {
				if goOut != cOut {
					t.Errorf("getHexagonEdgeLengthAvgKm(%d) output mismatch: Go=%v, C=%v", tt.res, goOut, cOut)
				}
			}
		})
	}
}
