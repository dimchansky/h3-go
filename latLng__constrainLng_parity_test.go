//go:build cgo

package h3

import (
	"math"
	"testing"
)

func Test_constrainLng_ParityWithC(t *testing.T) {
	twoPi := 2 * math.Pi
	cases := []float64{
		-twoPi - 1e-6, -twoPi, -math.Pi - 1e-12, -3.0,
		-1.0, 0.0, 1.0, 3.0, math.Pi + 1e-12, twoPi, twoPi + 1e-6,
	}
	for _, v := range cases {
		inputAngle := Rad(v)
		goVal := constrainLng(inputAngle)
		cVal := _constrainLngC(inputAngle)
		if goVal.Rad() != cVal.Rad() {
			t.Fatalf("constrainLng parity mismatch: in=%.12g go=%.12g c=%.12g", v, goVal.Rad(), cVal.Rad())
		}
	}
}
