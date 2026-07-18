//go:build cgo && c2go && h3v450

package h3

import "testing"

// Parity tests for the H3 4.5.0 header-only vec3d implementation
// (vec3d.h::latLngToVec3, vec3ToLatLng, vec3LinComb, vec3Cross, vec3Dot,
// vec3NormSq, vec3Norm, vec3Normalize, vec3DistSq) against the exact C
// definitions. Pure-arithmetic helpers compare bit-exactly (the harness
// pins C to strict IEEE with -ffp-contract=off and the Go bodies force
// per-product rounding via explicit conversions); the
// trig-dependent conversions
// (latLngToVec3, vec3ToLatLng) admit a last-ulp difference because Go's
// math library and the platform libm legitimately differ by 1 ulp on
// some sin/cos/asin/atan2 inputs.

// vec3UlpClose reports whether two doubles are within a couple of ulps —
// the tolerance for libm-dependent comparisons only.
func vec3UlpClose(a, b float64) bool {
	if a == b {
		return true
	}
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	scale := 1.0
	if aa := a; aa < 0 {
		aa = -aa
		if aa > scale {
			scale = aa
		}
	} else if a > scale {
		scale = a
	}
	return diff <= 1e-15*scale
}

func vec3UlpCloseVec(a, b vec3d) bool {
	return vec3UlpClose(a.X, b.X) && vec3UlpClose(a.Y, b.Y) && vec3UlpClose(a.Z, b.Z)
}

var vec3ParityVectors = []vec3d{
	{0, 0, 0},
	{1, 0, 0},
	{0, 1, 0},
	{0, 0, 1},
	{1, 1, 1},
	{3, -4, 12},
	{0.2995487373996265, -0.7694664284942747, 0.5640850727882765},
	{1e-163, 0, 0},
	{2.220446049250313e-16 / 2, 0, 0},
	{-0.5, 0.25, -0.125},
}

var vec3ParityGeos = []LatLng{
	{Lat: Rad(0), Lng: Rad(0)},
	{Lat: Rad(0.5), Lng: Rad(-1.3)},
	{Lat: Rad(0.6), Lng: Rad(-1.2)},
	{Lat: Rad(-0.4), Lng: Rad(0.8)},
	{Lat: Rad(1.5), Lng: Rad(3.1)},
}

func Test_vec3Ops_parity(t *testing.T) {
	for _, g := range vec3ParityGeos {
		goV := latLngToVec3(g)
		cV := latLngToVec3C(g)
		if !vec3UlpCloseVec(goV, cV) {
			t.Errorf("latLngToVec3(%v): Go=%v C=%v", g, goV, cV)
		}
		goLL := vec3ToLatLng(goV)
		cLL := vec3ToLatLngC(goV)
		if !vec3UlpClose(goLL.Lat.Rad(), cLL.Lat.Rad()) || !vec3UlpClose(goLL.Lng.Rad(), cLL.Lng.Rad()) {
			t.Errorf("vec3ToLatLng(%v): Go=%v C=%v", goV, goLL, cLL)
		}
	}

	for i, v1 := range vec3ParityVectors {
		for _, v2 := range vec3ParityVectors {
			if got, want := vec3Dot(v1, v2), vec3DotC(v1, v2); got != want {
				t.Errorf("vec3Dot(%v,%v): Go=%v C=%v", v1, v2, got, want)
			}
			if got, want := vec3Cross(v1, v2), vec3CrossC(v1, v2); got != want {
				t.Errorf("vec3Cross(%v,%v): Go=%v C=%v", v1, v2, got, want)
			}
			if got, want := vec3DistSq(v1, v2), vec3DistSqC(v1, v2); got != want {
				t.Errorf("vec3DistSq(%v,%v): Go=%v C=%v", v1, v2, got, want)
			}
			a := float64(i) - 2.5
			if got, want := vec3LinComb(a, v1, -a, v2), vec3LinCombC(a, v1, -a, v2); got != want {
				t.Errorf("vec3LinComb(%v,%v,%v): Go=%v C=%v", a, v1, v2, got, want)
			}
		}
		if got, want := vec3NormSq(v1), vec3NormSqC(v1); got != want {
			t.Errorf("vec3NormSq(%v): Go=%v C=%v", v1, got, want)
		}
		if got, want := vec3Norm(v1), vec3NormC(v1); got != want {
			t.Errorf("vec3Norm(%v): Go=%v C=%v", v1, got, want)
		}
		goN, cN := v1, v1
		vec3Normalize(&goN)
		vec3NormalizeC(&cN)
		if goN != cN {
			t.Errorf("vec3Normalize(%v): Go=%v C=%v", v1, goN, cN)
		}
	}
}
