//go:build cgo

package h3

import (
	"testing"
)

func Test_iterInitParent_parity(t *testing.T) {
	tests := []struct {
		name     string
		h        H3Index
		childRes int32
	}{
		{
			"valid_parent_to_child",
			H3Index(0x85283473fffffff), // res 5 cell
			7,                          // child res 7
		},
		{
			"parent_res_0_child_res_5",
			H3Index(0x8003fffffffffff), // res 0 cell
			5,                          // child res 5
		},
		{
			"parent_res_10_child_res_15",
			H3Index(0x8a2834700007fff), // res 10 cell
			15,                         // child res 15 (max)
		},
		{
			"same_resolution",
			H3Index(0x85283473fffffff), // res 5 cell
			5,                          // same res 5
		},
		{
			"invalid_child_res_too_small",
			H3Index(0x85283473fffffff), // res 5 cell
			3,                          // child res < parent res
		},
		{
			"invalid_child_res_too_large",
			H3Index(0x85283473fffffff), // res 5 cell
			16,                         // child res > MAX_H3_RES
		},
		{
			"null_index",
			H3Index(0), // H3_NULL
			5,
		},
		{
			"pentagon_parent",
			H3Index(0x81283fffffffffff), // res 1 pentagon
			3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test Go implementation
			var iterGo IterCellsChildren
			iterInitParent(tt.h, tt.childRes, &iterGo)

			// Test C implementation
			var iterC IterCellsChildren
			iterInitParentC(tt.h, tt.childRes, &iterC)

			// Compare results
			if iterGo.H != iterC.H {
				t.Errorf("iterInitParent() H mismatch: Go=0x%x, C=0x%x", iterGo.H, iterC.H)
			}
			if iterGo.ParentRes != iterC.ParentRes {
				t.Errorf("iterInitParent() ParentRes mismatch: Go=%d, C=%d", iterGo.ParentRes, iterC.ParentRes)
			}
			if iterGo.SkipDigit != iterC.SkipDigit {
				t.Errorf("iterInitParent() SkipDigit mismatch: Go=%d, C=%d", iterGo.SkipDigit, iterC.SkipDigit)
			}
		})
	}

	// Test deterministic behavior
	t.Run("deterministic", func(t *testing.T) {
		h := H3Index(0x85283473fffffff)
		childRes := int32(7)

		var iter1, iter2 IterCellsChildren
		iterInitParent(h, childRes, &iter1)
		iterInitParent(h, childRes, &iter2)

		if iter1.H != iter2.H || iter1.ParentRes != iter2.ParentRes || iter1.SkipDigit != iter2.SkipDigit {
			t.Errorf("iterInitParent should be deterministic: first=%+v != second=%+v", iter1, iter2)
		}
	})
}
