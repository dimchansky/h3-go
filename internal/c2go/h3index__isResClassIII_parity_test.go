//go:build c2go

package c2go

import "testing"

// setResolutionForTest mutates the resolution bits of an H3Index for testing.
func setResolutionForTest(h H3Index, res int) H3Index {
	const H3_RES_OFFSET = 52
	const H3_RES_MASK = uint64(15) << H3_RES_OFFSET
	x := uint64(h)
	x &^= H3_RES_MASK
	x |= (uint64(res) & 15) << H3_RES_OFFSET
	return H3Index(x)
}

func Test_h3index_isResClassIII_ParityWithC(t *testing.T) {
	bases := []H3Index{0x8928308280fffff, 0x821c07fffffffff}
	for _, b := range bases {
		for res := 0; res <= 15; res++ {
			h := setResolutionForTest(b, res)
			goVal := isResClassIII(h)
			cVal := isResClassIIIC(h)
			if goVal != cVal {
				t.Fatalf("isResClassIII mismatch for res=%d base=%x: go=%d c=%d", res, uint64(b), goVal, cVal)
			}
		}
	}
}
