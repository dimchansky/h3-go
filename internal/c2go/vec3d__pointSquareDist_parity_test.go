//go:build c2go

package c2go

import (
	"math"
	"math/rand"
	"testing"
)

func Test_vec3d__pointSquareDist_Parity(t *testing.T) {
	rng := rand.New(rand.NewSource(7))
	cases := []struct{ a, b Vec3d }{
		{Vec3d{0, 0, 0}, Vec3d{0, 0, 0}},
		{Vec3d{1, 0, 0}, Vec3d{0, 1, 0}},
		{Vec3d{1, 2, 3}, Vec3d{4, 5, 6}},
		{Vec3d{-1, -2, -3}, Vec3d{3, 2, 1}},
		{Vec3d{1e-9, -1e-9, 2e-9}, Vec3d{-2e-9, 3e-9, -4e-9}},
	}
	for i := 0; i < 100; i++ {
		a := Vec3d{rng.NormFloat64(), rng.NormFloat64(), rng.NormFloat64()}
		b := Vec3d{rng.NormFloat64(), rng.NormFloat64(), rng.NormFloat64()}
		cases = append(cases, struct{ a, b Vec3d }{a: a, b: b})
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
