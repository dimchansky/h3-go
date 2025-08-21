//go:build c2go

package c2go

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
        goVal := constrainLng(v)
        cVal := _constrainLngC(v)
        if goVal != cVal {
            t.Fatalf("constrainLng parity mismatch: in=%.12g go=%.12g c=%.12g", v, goVal, cVal)
        }
    }
}

