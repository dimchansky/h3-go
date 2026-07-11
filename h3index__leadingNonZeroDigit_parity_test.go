//go:build cgo && c2go

package h3

import "testing"

func Test_h3LeadingNonZeroDigit_ParityWithC(t *testing.T) {
	hs := []h3Index{0x8928308280fffff, 0x821c07fffffffff, 0x8a195da49a2ffff}
	for _, h := range hs {
		if _h3LeadingNonZeroDigit(h) != h3LeadingNonZeroDigitC(h) {
			t.Fatalf("_h3LeadingNonZeroDigit mismatch for %x", uint64(h))
		}
	}
}
