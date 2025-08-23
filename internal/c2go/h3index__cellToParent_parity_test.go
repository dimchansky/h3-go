//go:build cgo

package c2go

import "testing"

func Test_h3index_cellToParent_ParityWithC(t *testing.T) {
	hs := []H3Index{0x8928308280fffff, 0x821c07fffffffff}
	cases := []int{-1, 0, 1, 5, 15, 16}
	for _, h := range hs {
		for _, parentRes := range cases {
			goOut, goErr := cellToParent(h, parentRes)
			cOut, cErr := cellToParentC(h, parentRes)
			if goErr != H3Error(cErr) || goOut != cOut {
				t.Fatalf("cellToParent mismatch h=%x parentRes=%d: go(out=%x,err=%d) c(out=%x,err=%d)", uint64(h), parentRes, uint64(goOut), goErr, uint64(cOut), cErr)
			}
		}
	}
}
