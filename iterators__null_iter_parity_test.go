//go:build cgo && c2go

package h3

import (
	"testing"
)

func Test_nullIter_parity(t *testing.T) {
	// Test Go implementation
	iterGo := nullIter()

	// Test C implementation
	iterC := nullIterC()

	// Compare results
	if iterGo.H != iterC.H {
		t.Errorf("nullIter() H mismatch: Go=0x%x, C=0x%x", iterGo.H, iterC.H)
	}
	if iterGo.ParentRes != iterC.ParentRes {
		t.Errorf("nullIter() ParentRes mismatch: Go=%d, C=%d", iterGo.ParentRes, iterC.ParentRes)
	}
	if iterGo.SkipDigit != iterC.SkipDigit {
		t.Errorf("nullIter() SkipDigit mismatch: Go=%d, C=%d", iterGo.SkipDigit, iterC.SkipDigit)
	}

	// Verify expected values
	if iterGo.H != 0 {
		t.Errorf("nullIter() H should be 0 (h3Null), got 0x%x", iterGo.H)
	}
	if iterGo.ParentRes != -1 {
		t.Errorf("nullIter() ParentRes should be -1, got %d", iterGo.ParentRes)
	}
	if iterGo.SkipDigit != -1 {
		t.Errorf("nullIter() SkipDigit should be -1, got %d", iterGo.SkipDigit)
	}

	// Test deterministic behavior
	t.Run("deterministic", func(t *testing.T) {
		iter1 := nullIter()
		iter2 := nullIter()

		if iter1.H != iter2.H || iter1.ParentRes != iter2.ParentRes || iter1.SkipDigit != iter2.SkipDigit {
			t.Errorf("nullIter should be deterministic: first=%+v != second=%+v", iter1, iter2)
		}
	})
}
