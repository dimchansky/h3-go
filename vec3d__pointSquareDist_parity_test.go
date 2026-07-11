//go:build cgo && c2go

package h3

import (
	"math"
	"math/rand"
	"testing"
)

func Test_vec3d__pointSquareDist_Parity(t *testing.T) {
	rng := rand.New(rand.NewSource(7))
	cases := []struct{ a, b vec3d }{
		{vec3d{0, 0, 0}, vec3d{0, 0, 0}},
		{vec3d{1, 0, 0}, vec3d{0, 1, 0}},
		{vec3d{1, 2, 3}, vec3d{4, 5, 6}},
		{vec3d{-1, -2, -3}, vec3d{3, 2, 1}},
		{vec3d{1e-9, -1e-9, 2e-9}, vec3d{-2e-9, 3e-9, -4e-9}},
	}
	for i := 0; i < 100; i++ {
		a := vec3d{rng.NormFloat64(), rng.NormFloat64(), rng.NormFloat64()}
		b := vec3d{rng.NormFloat64(), rng.NormFloat64(), rng.NormFloat64()}
		cases = append(cases, struct{ a, b vec3d }{a: a, b: b})
	}

	const tol = 1e-14
	for i, tc := range cases {
		gotGo := _pointSquareDist(&tc.a, &tc.b)
		gotC := pointSquareDistC(tc.a, tc.b)
		if math.IsNaN(gotGo) || math.IsNaN(gotC) {
			t.Fatalf("case %d produced NaN: go=%v c=%v", i, gotGo, gotC)
		}
		if math.Abs(gotGo-gotC) > tol {
			t.Fatalf("case %d mismatch: go=%0.17g c=%0.17g (a=%+v b=%+v)", i, gotGo, gotC, tc.a, tc.b)
		}
	}
}
