//go:build cgo && c2go

package h3

import "testing"

func Test_h3index_zeroIndexDigits_ParityWithC(t *testing.T) {
	hs := []H3Index{0x8928308280fffff, 0x821c07fffffffff, 0x8a195da49a2ffff}
	// start > end: no-op
	for _, h := range hs {
		goVal := _zeroIndexDigits(h, 10, 5)
		cVal := zeroIndexDigitsC(h, 10, 5)
		if goVal != cVal {
			t.Fatalf("_zeroIndexDigits mismatch (start>end) for %x: go=%x c=%x", uint64(h), uint64(goVal), uint64(cVal))
		}
	}
	// Boundaries and typical ranges
	starts := []int32{0, 1, 5, 15}
	ends := []int32{0, 1, 5, 15, 16}
	for _, h := range hs {
		for _, s := range starts {
			for _, e := range ends {
				goVal := _zeroIndexDigits(h, s, e)
				cVal := zeroIndexDigitsC(h, s, e)
				if goVal != cVal {
					t.Fatalf("_zeroIndexDigits mismatch h=%x start=%d end=%d: go=%x c=%x", uint64(h), s, e, uint64(goVal), uint64(cVal))
				}
			}
		}
	}
}
