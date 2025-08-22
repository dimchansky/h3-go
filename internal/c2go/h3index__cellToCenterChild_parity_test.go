//go:build c2go

package c2go

import "testing"

func Test_h3index_cellToCenterChild_ParityWithC(t *testing.T) {
	hs := []H3Index{0x8928308280fffff, 0x821c07fffffffff}
	for _, h := range hs {
		baseRes := getResolution(h)
		cases := []int{baseRes - 1, baseRes, baseRes + 1, 15, 16}
		for _, childRes := range cases {
			goOut, goErr := cellToCenterChild(h, childRes)
			cOut, cErr := cellToCenterChildC(h, childRes)
			if goErr != H3Error(cErr) || goOut != cOut {
				t.Fatalf("cellToCenterChild mismatch h=%x childRes=%d: go(out=%x,err=%d) c(out=%x,err=%d)", uint64(h), childRes, uint64(goOut), goErr, uint64(cOut), cErr)
			}
		}
	}
}
