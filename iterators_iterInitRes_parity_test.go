//go:build cgo

package h3

import (
	"testing"
)

func Test_iterInitRes_parity(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		res  int32
	}{
		{"Resolution 0", 0},
		{"Resolution 1", 1},
		{"Resolution 2", 2},
		{"Resolution 3", 3},
		{"Resolution 4", 4},
		{"Resolution 15", 15},

		// Invalid cases (should return null iterator)
		{"Negative resolution", -1},
		{"Resolution too large", 16},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goResult := iterInitRes(tt.res)
			cResult := iterInitResC(tt.res)

			if goResult.H != cResult.H {
				t.Errorf("H mismatch: go=%x, c=%x", uint64(goResult.H), uint64(cResult.H))
			}
			if goResult.baseCellNum != cResult.baseCellNum {
				t.Errorf("baseCellNum mismatch: go=%d, c=%d", goResult.baseCellNum, cResult.baseCellNum)
			}
			if goResult.res != cResult.res {
				t.Errorf("res mismatch: go=%d, c=%d", goResult.res, cResult.res)
			}
			if goResult.itC.H != cResult.itC.H {
				t.Errorf("itC.H mismatch: go=%x, c=%x", uint64(goResult.itC.H), uint64(cResult.itC.H))
			}
			if goResult.itC.ParentRes != cResult.itC.ParentRes {
				t.Errorf("itC.ParentRes mismatch: go=%d, c=%d", goResult.itC.ParentRes, cResult.itC.ParentRes)
			}
			if goResult.itC.SkipDigit != cResult.itC.SkipDigit {
				t.Errorf("itC.SkipDigit mismatch: go=%d, c=%d", goResult.itC.SkipDigit, cResult.itC.SkipDigit)
			}
		})
	}
}

func Test_iterInitRes_initialValues_parity(t *testing.T) {
	t.Parallel()

	// Test that valid resolutions start with the expected initial state
	for res := int32(0); res <= 4; res++ {
		goResult := iterInitRes(res)
		cResult := iterInitResC(res)

		// For valid resolutions, the iterator should:
		// 1. Start with base cell 0
		// 2. Have the correct resolution
		// 3. Have the first cell from base cell 0 at that resolution
		if res >= 0 && res <= MAX_H3_RES {
			if goResult.baseCellNum != 0 {
				t.Errorf("Expected baseCellNum=0 for res=%d, got %d", res, goResult.baseCellNum)
			}
			if cResult.baseCellNum != 0 {
				t.Errorf("Expected C baseCellNum=0 for res=%d, got %d", res, cResult.baseCellNum)
			}

			if goResult.res != res {
				t.Errorf("Expected res=%d, got %d", res, goResult.res)
			}
			if cResult.res != res {
				t.Errorf("Expected C res=%d, got %d", res, cResult.res)
			}
		}
	}
}
