//go:build c2go

package c2go

import "testing"

func Test_h3index_cellToChildrenSize_ParityWithC(t *testing.T) {
	hs := []H3Index{0x8928308280fffff, 0x821c07fffffffff}
	for _, h := range hs {
		baseRes := getResolution(h)
		cases := []int{baseRes - 1, baseRes, baseRes + 1, 10, 15, 16}
		for _, childRes := range cases {
			goCount, goErr := cellToChildrenSize(h, childRes)
			cCount, cErr := cellToChildrenSizeC(h, childRes)
			if goErr != H3Error(cErr) || goCount != cCount {
				t.Fatalf("cellToChildrenSize mismatch h=%x childRes=%d: go(count=%d,err=%d) c(count=%d,err=%d)", uint64(h), childRes, goCount, goErr, cCount, cErr)
			}
		}
	}
}
