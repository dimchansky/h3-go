//go:build cgo && c2go

package h3

import (
	"math"
	"testing"
)

func Test_triangleEdgeLengthsToArea_ParityWithC(t *testing.T) {
	cases := [][3]float64{
		{0.1, 0.2, 0.3},
		{1.0, 1.1, 1.2},
		{0.5, 0.7, 0.9},
	}
	for _, tc := range cases {
		goVal := triangleEdgeLengthsToArea(tc[0], tc[1], tc[2])
		cVal := triangleEdgeLengthsToAreaC(tc[0], tc[1], tc[2])
		if math.Abs(goVal-cVal) > 1e-15 {
			t.Fatalf("triangleEdgeLengthsToArea mismatch: go=%.17g c=%.17g", goVal, cVal)
		}
	}
}
