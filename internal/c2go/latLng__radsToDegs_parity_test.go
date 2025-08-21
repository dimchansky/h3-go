//go:build c2go

package c2go

import "testing"

func Test_radsToDegs_ParityWithC(t *testing.T) {
	cases := []float64{-10, -6.283185307179586, -3.141592653589793, -1, 0, 1, 3.141592653589793, 6.283185307179586, 10, 1.23456}
	for _, v := range cases {
		goVal := radsToDegs(v)
		cVal := radsToDegsC(v)
		if goVal != cVal {
			t.Fatalf("radsToDegs parity mismatch: in=%g go=%.15g c=%.15g", v, goVal, cVal)
		}
	}
}
