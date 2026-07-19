package h3

import (
	"math"
	"testing"
)

// Regression test for issue #3: maxPolygonToCellsSizeExperimental on loops
// with NaN/Inf/huge-magnitude coordinates. Such inputs are in-domain for the
// upstream polygon fuzzers (raw doubles are used as radians), and before the
// lineCrossesLine NaN fix the Go port diverged from C on the Inf/NaN cases.
// Every expected size below is pinned to the value returned by H3 C 4.5.0
// (polyfill.c is byte-identical to 4.4.0, where these pins were first
// verified through the cgo parity harness); resolutions are kept low so the
// estimator's planet-wide scan on these degenerate bboxes stays cheap.
//
// The pathologically slow (~15 s) fuzz reproducer for the same code path is
// deliberately NOT replayed here — it is preserved under
// testdata/fuzz-findings/FuzzUpstreamPolygonOperations/ with replay
// instructions; both C and Go exhibit the same cost on it (C is slower), so
// the slowness is inherited upstream behavior, not a port defect.
func Test_maxPolygonToCellsSizeExperimental_nonFiniteCoordinates(t *testing.T) {
	t.Parallel()

	nan := Angle(math.NaN())
	posInf := Angle(math.Inf(1))

	cases := []struct {
		name string
		loop GeoLoop
		// wantSizes[res] is the C-verified size for that resolution.
		wantSizes []int64
	}{
		{
			name: "NaN lat",
			loop: GeoLoop{
				{Lat: nan, Lng: 0},
				{Lat: 0.1, Lng: 0.1},
				{Lat: 0.1, Lng: 0},
			},
			wantSizes: []int64{4, 6, 9, 16, 40},
		},
		{
			name: "+Inf lat",
			loop: GeoLoop{
				{Lat: posInf, Lng: 0},
				{Lat: 0.1, Lng: 0.1},
				{Lat: 0.1, Lng: 0},
			},
			wantSizes: []int64{5, 34, 237, 1658, 11605},
		},
		{
			name: "NaN lng",
			loop: GeoLoop{
				{Lat: 0.2, Lng: nan},
				{Lat: 0.1, Lng: 0.1},
				{Lat: 0.1, Lng: 0},
			},
			wantSizes: []int64{3, 6, 9, 16, 112},
		},
		{
			// Latitudes so large that cos(min(|north|, |south|)) in the
			// estimator's rough-area formula is negative, which defeats the
			// resolution-coarsening loop — the mechanism behind the slow
			// fuzz reproducer (issue #3), here at cheap resolutions.
			name: "huge-magnitude lat, negative rough area",
			loop: GeoLoop{
				{Lat: Angle(3.6864520893115203e+267), Lng: 0},
				{Lat: Angle(-4.782802678075986e+287), Lng: 0.1},
				{Lat: 0.1, Lng: 0},
			},
			wantSizes: []int64{21, 52, 135, 364, 953},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			polygon := GeoPolygon{GeoLoop: tc.loop}
			for res, want := range tc.wantSizes {
				size, err := maxPolygonToCellsSizeExperimental(&polygon, int32(res), 0)
				if err != eSuccess {
					t.Errorf("res %d: unexpected error %v", res, err)
					continue
				}
				if size != want {
					t.Errorf("res %d: size = %d, want %d (C-verified)", res, size, want)
				}
			}
		})
	}
}
