//go:build cgo

package h3

import (
	"math"
	"testing"
)

func Test_constrainLat_ParityWithC(t *testing.T) {
	cases := []float64{
		-math.Pi, -math.Pi/2 - 1e-12, -3.0, -1.0, 0.0, 1.0, 3.0, math.Pi/2 + 1e-12, math.Pi,
	}
	for _, v := range cases {
		goVal := constrainLat(v)
		cVal := _constrainLatC(v)
		if goVal != cVal {
			t.Fatalf("constrainLat parity mismatch: in=%.12g go=%.12g c=%.12g", v, goVal, cVal)
		}
	}
}
