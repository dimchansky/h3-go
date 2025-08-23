//go:build cgo

package c2go

import "testing"

func Test_h3index_isResolutionClassIII_ParityWithC(t *testing.T) {
	for res := int32(-1); res <= 20; res++ {
		goVal := isResolutionClassIII(res)
		cVal := isResolutionClassIIIC(res)
		// Convert both to boolean semantics for comparison
		// C function returns res % 2 (which can be -1 for negative odd numbers)
		// Go function returns proper boolean
		goBool := goVal
		cBool := cVal != 0
		if goBool != cBool {
			t.Fatalf("isResolutionClassIII boolean mismatch res=%d: go=%t c_raw=%d c_bool=%t", res, goBool, cVal, cBool)
		}
	}
}
