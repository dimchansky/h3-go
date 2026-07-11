//go:build cgo && c2go

package h3

import (
	"testing"
)

func Test_iterStepRes_parity(t *testing.T) {
	t.Parallel()

	// Test stepping through resolution 0 (should have exactly 122 cells)
	t.Run("Resolution 0", func(t *testing.T) {
		goIter := iterInitRes(0)
		cIter := iterInitResC(0)

		count := 0
		maxSteps := 200 // Safety limit

		for goIter.H != H3_NULL && cIter.H != H3_NULL && count < maxSteps {
			// Check that Go and C iterators are in sync
			if goIter.H != cIter.H {
				t.Errorf("Step %d: H mismatch: go=%x, c=%x", count, uint64(goIter.H), uint64(cIter.H))
			}
			if goIter.baseCellNum != cIter.baseCellNum {
				t.Errorf("Step %d: baseCellNum mismatch: go=%d, c=%d", count, goIter.baseCellNum, cIter.baseCellNum)
			}
			if goIter.res != cIter.res {
				t.Errorf("Step %d: res mismatch: go=%d, c=%d", count, goIter.res, cIter.res)
			}

			// Step both iterators
			iterStepRes(&goIter)
			iterStepResC(&cIter)
			count++
		}

		// Both should be exhausted at the same time
		if goIter.H != cIter.H {
			t.Errorf("Final state mismatch: go=%x, c=%x", uint64(goIter.H), uint64(cIter.H))
		}

		// Resolution 0 should have exactly NUM_BASE_CELLS cells
		if count != NUM_BASE_CELLS {
			t.Errorf("Expected %d cells for resolution 0, got %d", NUM_BASE_CELLS, count)
		}
	})

	// Test stepping through resolution 1
	t.Run("Resolution 1", func(t *testing.T) {
		goIter := iterInitRes(1)
		cIter := iterInitResC(1)

		count := 0
		maxSteps := 1000 // Safety limit for resolution 1

		for goIter.H != H3_NULL && cIter.H != H3_NULL && count < maxSteps {
			// Check that Go and C iterators are in sync
			if goIter.H != cIter.H {
				t.Errorf("Step %d: H mismatch: go=%x, c=%x", count, uint64(goIter.H), uint64(cIter.H))
				break
			}
			if goIter.baseCellNum != cIter.baseCellNum {
				t.Errorf("Step %d: baseCellNum mismatch: go=%d, c=%d", count, goIter.baseCellNum, cIter.baseCellNum)
				break
			}

			// Step both iterators
			iterStepRes(&goIter)
			iterStepResC(&cIter)
			count++
		}

		// Both should be exhausted at the same time
		if goIter.H != cIter.H {
			t.Errorf("Final state mismatch: go=%x, c=%x", uint64(goIter.H), uint64(cIter.H))
		}

		// Resolution 1 should have 842 cells (7*122 - 12 = 842, accounting for pentagon skipped digits)
		if count != 842 {
			t.Errorf("Expected 842 cells for resolution 1, got %d", count)
		}
	})
}

func Test_iterStepRes_nullIterator_parity(t *testing.T) {
	t.Parallel()

	// Test that null iterators behave correctly
	goIter := iterInitRes(-1) // Invalid resolution -> null iterator
	cIter := iterInitResC(-1)

	if goIter.H != H3_NULL {
		t.Errorf("Expected null Go iterator, got H=%x", uint64(goIter.H))
	}
	if cIter.H != H3_NULL {
		t.Errorf("Expected null C iterator, got H=%x", uint64(cIter.H))
	}

	// Stepping null iterators should keep them null
	iterStepRes(&goIter)
	iterStepResC(&cIter)

	if goIter.H != H3_NULL {
		t.Errorf("Expected null Go iterator after step, got H=%x", uint64(goIter.H))
	}
	if cIter.H != H3_NULL {
		t.Errorf("Expected null C iterator after step, got H=%x", uint64(cIter.H))
	}
}

func Test_iterStepRes_ordering_parity(t *testing.T) {
	t.Parallel()

	// Test that cells are returned in the same order for Go and C
	goIter := iterInitRes(1)
	cIter := iterInitResC(1)

	cells := make([][2]H3Index, 0, 50) // Store first 50 cells from both iterators

	for i := 0; i < 50 && goIter.H != H3_NULL && cIter.H != H3_NULL; i++ {
		cells = append(cells, [2]H3Index{goIter.H, cIter.H})
		iterStepRes(&goIter)
		iterStepResC(&cIter)
	}

	for i, pair := range cells {
		if pair[0] != pair[1] {
			t.Errorf("Cell %d mismatch: go=%x, c=%x", i, uint64(pair[0]), uint64(pair[1]))
		}
	}
}
