//go:build cgo && c2go

package h3

import "testing"

func Test_baseCells_isBaseCellPentagon_ParityWithC(t *testing.T) {
	for base := int32(0); base < NUM_BASE_CELLS; base++ {
		goVal := _isBaseCellPentagon(base)
		cVal := isBaseCellPentagonC(base)
		if goVal != cVal {
			t.Fatalf("_isBaseCellPentagon mismatch for base=%d: go=%v c=%v", base, goVal, cVal)
		}
	}
	// Check a couple out-of-range values behave as false
	for _, base := range []int32{-1, NUM_BASE_CELLS, 200} {
		if _isBaseCellPentagon(base) {
			t.Fatalf("_isBaseCellPentagon out-of-range base=%d should be false", base)
		}
	}
}
