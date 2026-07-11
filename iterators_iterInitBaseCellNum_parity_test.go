//go:build cgo && c2go

package h3

import (
	"testing"
)

func Test_iterInitBaseCellNum_parity(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		baseCellNum int32
		childRes    int32
	}{
		// Valid cases
		{"Base cell 0, res 0", 0, 0},
		{"Base cell 0, res 1", 0, 1},
		{"Base cell 0, res 3", 0, 3},
		{"Base cell 121, res 0", 121, 0},
		{"Base cell 121, res 1", 121, 1},
		{"Base cell 121, res 3", 121, 3},
		{"Base cell 24 (pentagon), res 1", 24, 1},
		{"Base cell 24 (pentagon), res 2", 24, 2},

		// Invalid cases (should return null iterator)
		{"Negative base cell", -1, 0},
		{"Base cell too large", 122, 0},
		{"Negative resolution", 0, -1},
		{"Resolution too large", 0, 16},
		{"Both invalid", -1, -1},
		{"Both at boundary", 122, 16},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goResult := iterInitBaseCellNum(tt.baseCellNum, tt.childRes)
			cResult := iterInitBaseCellNumC(tt.baseCellNum, tt.childRes)

			if goResult.H != cResult.H {
				t.Errorf("H mismatch: go=%x, c=%x", uint64(goResult.H), uint64(cResult.H))
			}
			if goResult.ParentRes != cResult.ParentRes {
				t.Errorf("ParentRes mismatch: go=%d, c=%d", goResult.ParentRes, cResult.ParentRes)
			}
			if goResult.SkipDigit != cResult.SkipDigit {
				t.Errorf("SkipDigit mismatch: go=%d, c=%d", goResult.SkipDigit, cResult.SkipDigit)
			}
		})
	}
}

func Test_iterInitBaseCellNum_exhaustive_parity(t *testing.T) {
	t.Parallel()

	// Test all base cells at resolution 0 and 1
	for baseCellNum := int32(0); baseCellNum < numBaseCells; baseCellNum++ {
		for childRes := int32(0); childRes <= 1; childRes++ {
			goResult := iterInitBaseCellNum(baseCellNum, childRes)
			cResult := iterInitBaseCellNumC(baseCellNum, childRes)

			if goResult.H != cResult.H {
				t.Errorf("H mismatch for baseCell=%d, childRes=%d: go=%x, c=%x",
					baseCellNum, childRes, uint64(goResult.H), uint64(cResult.H))
			}
			if goResult.ParentRes != cResult.ParentRes {
				t.Errorf("ParentRes mismatch for baseCell=%d, childRes=%d: go=%d, c=%d",
					baseCellNum, childRes, goResult.ParentRes, cResult.ParentRes)
			}
			if goResult.SkipDigit != cResult.SkipDigit {
				t.Errorf("SkipDigit mismatch for baseCell=%d, childRes=%d: go=%d, c=%d",
					baseCellNum, childRes, goResult.SkipDigit, cResult.SkipDigit)
			}
		}
	}
}
