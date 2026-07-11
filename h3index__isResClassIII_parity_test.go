//go:build cgo && c2go

package h3

import (
	"testing"
)

func Test_h3index_isResClassIII_ParityWithC(t *testing.T) {
	bases := []H3Index{0x8928308280fffff, 0x821c07fffffffff}
	for _, b := range bases {
		for res := int32(0); res <= 15; res++ {
			h := setResolution(b, res)
			goVal := isResClassIII(h)
			cVal := isResClassIIIC(h)
			// Convert bool to int for comparison
			if goVal != cVal {
				t.Fatalf("isResClassIII mismatch for res=%d base=%x: go=%t (%v) c=%v", res, uint64(b), goVal, goVal, cVal)
			}
		}
	}
}
