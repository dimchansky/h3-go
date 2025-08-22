//go:build c2go

package c2go

import "testing"

func Test_h3index_mode_highbit_ParityWithC(t *testing.T) {
	hs := []H3Index{0x0000000000000000, 0x8928308280fffff, 0x821c07fffffffff, 0xffffffffffffffff}
	// getMode / setMode parity
	for _, h := range hs {
		if getMode(h) != getModeC(h) {
			t.Fatalf("getMode mismatch for %x", uint64(h))
		}
		for v := 0; v <= 15; v++ {
			goH := setMode(h, v)
			cH := setModeC(h, v)
			if goH != cH {
				t.Fatalf("setMode mismatch h=%x v=%d: go=%x c=%x", uint64(h), v, uint64(goH), uint64(cH))
			}
		}
	}
	// getHighBit / setHighBit parity
	for _, h := range hs {
		if getHighBit(h) != getHighBitC(h) {
			t.Fatalf("getHighBit mismatch for %x", uint64(h))
		}
		for v := 0; v <= 1; v++ {
			goH := setHighBit(h, v)
			cH := setHighBitC(h, v)
			if goH != cH {
				t.Fatalf("setHighBit mismatch h=%x v=%d: go=%x c=%x", uint64(h), v, uint64(goH), uint64(cH))
			}
		}
	}
}
