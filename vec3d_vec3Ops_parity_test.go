//go:build cgo && c2go

package h3

import (
	"math"
	"testing"
)

// Parity tests for the H3 4.5.0 header-only vec3d implementation
// (vec3d.h::latLngToVec3, vec3ToLatLng, vec3LinComb, vec3Cross, vec3Dot,
// vec3NormSq, vec3Norm, vec3Normalize, vec3DistSq) against the exact C
// definitions. Pure-arithmetic helpers compare bit-exactly (the harness
// compiles the reference C with contraction disabled via
// -ffp-contract=off — see the Makefile test-c2go comment and
// CONTRIBUTING.md — and the Go bodies force per-product rounding via
// explicit conversions); the trig-dependent conversions (latLngToVec3,
// vec3ToLatLng) compare with vec3UlpClose because Go's math library and
// the platform libm legitimately differ by 1 ulp on some
// sin/cos/asin/atan2 inputs (measured ≤1 ulp on these cases).

// ulpDistance returns the number of representable float64 values
// separating a and b: the absolute difference of their positions on the
// ordered-bit number line (IEEE 754 bits with the negative half
// reflected), so adjacent doubles are 1 apart regardless of magnitude.
// Special cases: +0 and -0 are 0 apart; two NaNs are 0 apart (both
// sides agreeing on NaN counts as agreement); NaN vs non-NaN is
// MaxUint64. Infinities follow the ordered line: same-signed
// infinities are 0 apart, and an infinity is 1 apart from the
// same-signed largest finite value — so a caller that requires an
// infinity to agree only with itself (vec3UlpClose does) must check
// infinities before applying a distance bound.
func ulpDistance(a, b float64) uint64 {
	if math.IsNaN(a) || math.IsNaN(b) {
		if math.IsNaN(a) && math.IsNaN(b) {
			return 0
		}
		return math.MaxUint64
	}
	// Bias the ordered int64 line into uint64 so the subtraction below
	// cannot overflow (the full span exceeds int64).
	ua := orderedBits(a)
	ub := orderedBits(b)
	if ua >= ub {
		return ua - ub
	}
	return ub - ua
}

// orderedBits maps a non-NaN float64 to a uint64 that preserves numeric
// order: negative values (sign bit set) are reflected, then the whole
// line is biased by 2^63. -0 and +0 both map to the same point.
func orderedBits(f float64) uint64 {
	b := int64(math.Float64bits(f))
	if b < 0 {
		b = math.MinInt64 - b
	}
	return uint64(b) ^ (1 << 63)
}

// Tolerance for libm-dependent comparisons (see vec3UlpClose). Both
// bounds are measured, not claimed: sweeping every libm-dependent
// parity comparison at H3 4.5.0 (darwin/arm64, clang -ffp-contract=off
// vs Go math; 2026-07 corrective pass for #29-#32) observed at most
// 7 ulps in the relative regime and at most ~1.7e-16 absolute in the
// cancellation regime.
const (
	// vec3MaxUlp bounds the relative regime: values whose disagreement
	// is genuinely a few representable doubles apart.
	vec3MaxUlp = 8
	// vec3CancellationFloor bounds the cancellation regime: pipeline
	// outputs of small magnitude (~1e-4..1e-2 here) where a 1-ulp libm
	// difference upstream of a subtraction leaves an absolute residual
	// around zero. The residual is bounded by the inputs' magnitude
	// (~1 ulp of O(1) quantities), not by an ulp of the tiny result,
	// so it can span thousands of the result's ulps (2559 measured)
	// while staying below 2e-16 absolutely.
	vec3CancellationFloor = 1e-15
)

// vec3UlpClose reports whether two libm-dependent doubles agree up to
// the noise a 1-ulp libm difference can produce through the ported
// pipelines: within vec3MaxUlp ulps of each other, or within
// vec3CancellationFloor absolutely. It is NOT a pure ULP comparison —
// near-zero cancellation residuals measure thousands of ulps (see the
// constants above) — but special values are strict: NaN agrees only
// with NaN (via ulpDistance), an infinity only with itself (checked
// explicitly below — on the ordered-bit line an infinity is just 1
// apart from the same-signed largest finite value, so the ulp bound
// alone would wrongly accept that pair), and +0/-0 agree. Use only for
// libm-dependent comparisons; pure-arithmetic helpers compare with ==.
func vec3UlpClose(a, b float64) bool {
	if math.IsInf(a, 0) || math.IsInf(b, 0) {
		return a == b
	}
	if ulpDistance(a, b) <= vec3MaxUlp {
		return true
	}
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff <= vec3CancellationFloor
}

func vec3UlpCloseVec(a, b vec3d) bool {
	return vec3UlpClose(a.X, b.X) && vec3UlpClose(a.Y, b.Y) && vec3UlpClose(a.Z, b.Z)
}

// Meta-tests for the comparator itself: the documented special-value
// policy (an infinity agrees only with itself; NaN only with NaN;
// +0/-0 agree) and both tolerance regimes.
func Test_ulpDistance(t *testing.T) {
	next := func(f float64) float64 { return math.Nextafter(f, math.Inf(1)) }
	cases := []struct {
		name string
		a, b float64
		want uint64
	}{
		{"equal", 1.5, 1.5, 0},
		{"adjacent", 1.0, next(1.0), 1},
		{"adjacent negative", -1.0, math.Nextafter(-1.0, math.Inf(-1)), 1},
		{"signed zeros", 0.0, math.Copysign(0, -1), 0},
		{"across zero", 5e-324, -5e-324, 2},
		{"equal +Inf", math.Inf(1), math.Inf(1), 0},
		{"equal -Inf", math.Inf(-1), math.Inf(-1), 0},
		{"+Inf is 1 past MaxFloat64 on the ordered line", math.Inf(1), math.MaxFloat64, 1},
		{"-Inf is 1 past -MaxFloat64", math.Inf(-1), -math.MaxFloat64, 1},
		{"NaN vs NaN", math.NaN(), math.NaN(), 0},
		{"NaN vs finite", math.NaN(), 1.0, math.MaxUint64},
	}
	for _, tc := range cases {
		if got := ulpDistance(tc.a, tc.b); got != tc.want {
			t.Errorf("%s: ulpDistance(%v, %v) = %d, want %d", tc.name, tc.a, tc.b, got, tc.want)
		}
		if got := ulpDistance(tc.b, tc.a); got != tc.want {
			t.Errorf("%s (swapped): ulpDistance(%v, %v) = %d, want %d", tc.name, tc.b, tc.a, got, tc.want)
		}
	}
}

func Test_vec3UlpClose(t *testing.T) {
	// 9 ulps above 1.0 is ~2e-15 — beyond both the ulp bound and the
	// absolute floor; 9 ulps above 1e-16 is far below the floor.
	up := func(f float64, n int) float64 {
		for i := 0; i < n; i++ {
			f = math.Nextafter(f, math.Inf(1))
		}
		return f
	}
	cases := []struct {
		name string
		a, b float64
		want bool
	}{
		{"equal finite", 1.5, 1.5, true},
		{"adjacent finite", 1.0, up(1.0, 1), true},
		{"at the ulp bound", 1.0, up(1.0, int(vec3MaxUlp)), true},
		{"beyond the ulp bound and the floor", 1.0, up(1.0, int(vec3MaxUlp)+1), false},
		{"many ulps but under the cancellation floor", 1e-16, up(1e-16, 100), true},
		{"across zero under the floor", 1e-16, -1e-16, true},
		{"signed zeros", 0.0, math.Copysign(0, -1), true},
		{"across zero beyond the floor", 1e-14, -1e-14, false},
		{"equal +Inf", math.Inf(1), math.Inf(1), true},
		{"equal -Inf", math.Inf(-1), math.Inf(-1), true},
		{"+Inf vs MaxFloat64", math.Inf(1), math.MaxFloat64, false},
		{"-Inf vs -MaxFloat64", math.Inf(-1), -math.MaxFloat64, false},
		{"opposite infinities", math.Inf(1), math.Inf(-1), false},
		{"Inf vs finite", math.Inf(1), 1.0, false},
		{"NaN vs NaN", math.NaN(), math.NaN(), true},
		{"NaN vs finite", math.NaN(), 1.0, false},
		{"NaN vs Inf", math.NaN(), math.Inf(1), false},
	}
	for _, tc := range cases {
		if got := vec3UlpClose(tc.a, tc.b); got != tc.want {
			t.Errorf("%s: vec3UlpClose(%v, %v) = %v, want %v", tc.name, tc.a, tc.b, got, tc.want)
		}
		if got := vec3UlpClose(tc.b, tc.a); got != tc.want {
			t.Errorf("%s (swapped): vec3UlpClose(%v, %v) = %v, want %v", tc.name, tc.b, tc.a, got, tc.want)
		}
	}
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
