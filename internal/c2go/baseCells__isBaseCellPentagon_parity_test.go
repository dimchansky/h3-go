//go:build c2go

package c2go

import "testing"

func Test_baseCells_isBaseCellPentagon_ParityWithC(t *testing.T) {
	for base := 0; base < NUM_BASE_CELLS; base++ {
		goVal := _isBaseCellPentagon(base) != 0
		cVal := isBaseCellPentagonC(base) != 0
		if goVal != cVal {
			t.Fatalf("_isBaseCellPentagon mismatch for base=%d: go=%v c=%v", base, goVal, cVal)
		}
	}
	// Check a couple out-of-range values behave as false
	for _, base := range []int{-1, NUM_BASE_CELLS, 200} {
		if _isBaseCellPentagon(base) != 0 {
			t.Fatalf("_isBaseCellPentagon out-of-range base=%d should be false", base)
		}
	}
}
