//go:build cgo && c2go

package h3

import "testing"

func Test_h3index_makeDirectChild_ParityWithC(t *testing.T) {
	hs := []H3Index{0x8928308280fffff, 0x821c07fffffffff}
	for _, h := range hs {
		for d := int32(0); d <= 6; d++ {
			goH := makeDirectChild(h, d)
			cH := makeDirectChildC(h, d)
			if goH != cH {
				t.Fatalf("makeDirectChild mismatch h=%x d=%d: go=%x c=%x", uint64(h), d, uint64(goH), uint64(cH))
			}
		}
	}
}
