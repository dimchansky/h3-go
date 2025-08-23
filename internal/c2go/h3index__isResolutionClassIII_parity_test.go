//go:build cgo

package c2go

import "testing"

func Test_h3index_isResolutionClassIII_ParityWithC(t *testing.T) {
	for res := -1; res <= 20; res++ {
		goVal := isResolutionClassIII(res)
		cVal := isResolutionClassIIIC(res)
		if goVal != cVal {
			t.Fatalf("isResolutionClassIII mismatch res=%d: go=%d c=%d", res, goVal, cVal)
		}
	}
}
