//go:build cgo

package c2go

import (
	"math"
	"testing"
)

func Test_posAngleRads_ParityWithC(t *testing.T) {
	twoPi := 2 * math.Pi
	cases := []float64{
		-0.0, 0.0, -1e-12, 1e-12,
		-math.Pi / 2, math.Pi / 2,
		-math.Pi + 1e-9, math.Pi - 1e-9,
		-twoPi + 1e-6, twoPi - 1e-6,
		twoPi + 1e-6, -twoPi - 1e-6,
	}
	for _, r := range cases {
		gotGo := _posAngleRads(r)
		gotC := _posAngleRadsC(r)
		if math.IsNaN(gotGo) != math.IsNaN(gotC) || (!math.IsNaN(gotGo) && gotGo != gotC) {
			t.Fatalf("posAngleRads parity mismatch: r=%.12g go=%.12g c=%.12g", r, gotGo, gotC)
		}
	}
}
