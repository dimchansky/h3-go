//go:build cgo && c2go

package h3

import (
	"math"
	"testing"
)

func Test_v2dMag_ParityWithC(t *testing.T) {
	cases := []Vec2d{{0, 0}, {1, 0}, {0, 1}, {3, 4}, {-2.5, 4.5}}
	for _, v := range cases {
		goVal := _v2dMag(&v)
		cVal := v2dMagC(v)
		if math.Abs(goVal-cVal) > 1e-15 {
			t.Fatalf("_v2dMag mismatch: go=%.17g c=%.17g", goVal, cVal)
		}
	}
}
