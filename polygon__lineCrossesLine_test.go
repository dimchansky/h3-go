package h3

import (
	"math"
	"testing"
)

// Regression test for issue #3's NaN/Inf finding: C's lineCrossesLine ends
// with `return (test >= 0 && test <= 1)`, which is false when `test` is NaN.
// An earlier Go port used the negated form `if test < 0 || test > 1 { return
// false }; return true`, which returned true for NaN and broke parity with C
// on polygons containing NaN/Inf/huge-magnitude coordinates (in-domain for
// the upstream polygon fuzzers). All expectations below are verified against
// H3 C 4.4.0 via the parity harness (Test_lineCrossesLine_ParityWithC).
func Test_lineCrossesLine_nonFiniteCoordinates(t *testing.T) {
	t.Parallel()

	nan := Angle(math.NaN())
	posInf := Angle(math.Inf(1))
	negInf := Angle(math.Inf(-1))

	cases := []struct {
		name           string
		a1, a2, b1, b2 LatLng
		want           bool
	}{
		// NaN/Inf coordinates make the intersection parameter NaN; C
		// reports "no crossing" for every such segment pair.
		{"NaN lat on a1", LatLng{Lat: nan, Lng: 0}, LatLng{Lat: 0.1, Lng: 0.1}, LatLng{Lat: 0, Lng: 0.05}, LatLng{Lat: 0.2, Lng: 0.05}, false},
		{"+Inf lat on a1", LatLng{Lat: posInf, Lng: 0}, LatLng{Lat: 0.1, Lng: 0.1}, LatLng{Lat: 0, Lng: 0.05}, LatLng{Lat: 0.2, Lng: 0.05}, false},
		{"NaN lng on a1", LatLng{Lat: 0.2, Lng: nan}, LatLng{Lat: 0.1, Lng: 0.1}, LatLng{Lat: 0, Lng: 0.05}, LatLng{Lat: 0.2, Lng: 0.05}, false},
		{"-Inf lng on b1", LatLng{Lat: 0, Lng: 0}, LatLng{Lat: 0.2, Lng: 0.2}, LatLng{Lat: 0.1, Lng: negInf}, LatLng{Lat: 0.1, Lng: 0.1}, false},
		// Finite sanity checks: behavior unchanged for ordinary segments.
		{"finite crossing", LatLng{Lat: 0, Lng: 0}, LatLng{Lat: 1, Lng: 1}, LatLng{Lat: 0, Lng: 1}, LatLng{Lat: 1, Lng: 0}, true},
		{"finite non-crossing", LatLng{Lat: 0, Lng: 0}, LatLng{Lat: 1, Lng: 1}, LatLng{Lat: 0, Lng: 2}, LatLng{Lat: 1, Lng: 2}, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := lineCrossesLine(&tc.a1, &tc.a2, &tc.b1, &tc.b2); got != tc.want {
				t.Errorf("lineCrossesLine(%+v, %+v, %+v, %+v) = %v, want %v (C-verified)",
					tc.a1, tc.a2, tc.b1, tc.b2, got, tc.want)
			}
		})
	}
}
