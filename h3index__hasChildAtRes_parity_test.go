//go:build cgo && c2go

package h3

import "testing"

func Test_h3index_hasChildAtRes_ParityWithC(t *testing.T) {
	hs := []H3Index{0x8928308280fffff, 0x821c07fffffffff}
	for _, h := range hs {
		for childRes := int32(-1); childRes <= 16; childRes++ {
			goVal := _hasChildAtRes(h, childRes)
			cVal := hasChildAtResC(h, childRes) != 0
			if goVal != cVal {
				t.Fatalf("_hasChildAtRes mismatch h=%x childRes=%d: go=%v c=%v", uint64(h), childRes, goVal, cVal)
			}
		}
	}
}
