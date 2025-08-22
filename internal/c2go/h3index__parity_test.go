//go:build c2go

package c2go

import "testing"

func Test_string_h3_ParityWithC(t *testing.T) {
	// h3ToString parity
	h := H3Index(0x8928308280fffff)
	goStr, goErr := h3ToString(h)
	cStr, cErr := h3ToStringC(h)
	if goErr != cErr || goStr != cStr {
		t.Fatalf("h3ToString mismatch: go=(%s,%d) c=(%s,%d)", goStr, goErr, cStr, cErr)
	}
	// stringToH3 parity
	goH, goE := stringToH3(cStr)
	cH, cE := stringToH3C(cStr)
	if goE != H3Error(cE) || goH != cH {
		t.Fatalf("stringToH3 mismatch: go=(%x,%d) c=(%x,%d)", uint64(goH), goE, uint64(cH), cE)
	}
	// invalid
	_, goE = stringToH3("nothex")
	_, cE = stringToH3C("nothex")
	if (goE == 0) == (cE == 0) == false { // both should be errors
		t.Fatalf("stringToH3 invalid input mismatch: go=%d c=%d", goE, cE)
	}
}

func Test_getters_ParityWithC(t *testing.T) {
	hs := []H3Index{0x8928308280fffff, 0x821c07fffffffff, 0x8a195da49a2ffff}
	for _, h := range hs {
		if getResolution(h) != getResolutionC(h) {
			t.Fatalf("getResolution mismatch for %x", uint64(h))
		}
		if getBaseCellNumber(h) != getBaseCellNumberC(h) {
			t.Fatalf("getBaseCellNumber mismatch for %x", uint64(h))
		}
	}
}
