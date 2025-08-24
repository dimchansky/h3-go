//go:build cgo

package c2go

import "testing"

func Test_normalizeLng_ParityWithC(t *testing.T) {
	cases := []struct {
		lng  float64
		mode LongitudeNormalization
	}{
		{0, NORMALIZE_NONE},
		{-1.0, NORMALIZE_EAST},
		{1.0, NORMALIZE_EAST},
		{1.0, NORMALIZE_WEST},
		{-1.0, NORMALIZE_WEST},
	}
	for _, tc := range cases {
		goVal := normalizeLng(tc.lng, tc.mode)
		cVal := normalizeLngC(tc.lng, tc.mode)
		if goVal != cVal {
			t.Fatalf("normalizeLng mismatch: in=%g mode=%d go=%g c=%g", tc.lng, tc.mode, goVal, cVal)
		}
	}
}
