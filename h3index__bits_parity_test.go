//go:build cgo && c2go

package h3

import "testing"

func Test_h3index_bits_ParityWithC(t *testing.T) {
	hs := []H3Index{0x8928308280fffff, 0x821c07fffffffff, 0x8a195da49a2ffff}
	// Reserved bits get
	for _, h := range hs {
		if getReservedBits(h) != getReservedBitsC(h) {
			t.Fatalf("getReservedBits mismatch for %x", uint64(h))
		}
	}
	// Reserved bits set
	for _, h := range hs {
		for v := int32(0); v <= 7; v++ {
			goH := setReservedBits(h, v)
			cH := setReservedBitsC(h, v)
			if goH != cH {
				t.Fatalf("setReservedBits mismatch for %x v=%d: go=%x c=%x", uint64(h), v, uint64(goH), uint64(cH))
			}
		}
	}
	// Index digit get/set at several resolutions
	for _, h := range hs {
		for res := int32(0); res <= 15; res++ {
			if getIndexDigit(h, res) != getIndexDigitC(h, res) {
				t.Fatalf("getIndexDigit mismatch h=%x res=%d", uint64(h), res)
			}
		}
		// Set to each possible digit and verify parity
		for res := int32(1); res <= 15; res++ {
			for d := int32(0); d <= 7; d++ {
				goH := setIndexDigit(h, res, d)
				cH := setIndexDigitC(h, res, d)
				if goH != cH {
					t.Fatalf("setIndexDigit mismatch h=%x res=%d d=%d: go=%x c=%x", uint64(h), res, d, uint64(goH), uint64(cH))
				}
			}
		}
	}
}
