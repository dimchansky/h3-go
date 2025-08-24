//go:build cgo

package c2go

import (
	"testing"
)

func Test_iterStepChild_parity(t *testing.T) {
	tests := []struct {
		name     string
		parent   H3Index
		childRes int32
		numSteps int // Number of steps to test
	}{
		{
			name:     "hexagon parent res 0 to res 1",
			parent:   0x8001fffffffffff, // res 0 hexagon
			childRes: 1,
			numSteps: 10, // Test first 10 children
		},
		{
			name:     "hexagon parent res 3 to res 4",
			parent:   0x83080dfffffff, // res 3 hexagon
			childRes: 4,
			numSteps: 10,
		},
		{
			name:     "pentagon parent res 0 to res 1",
			parent:   0x8007fffffffffff, // res 0 pentagon (base cell 4)
			childRes: 1,
			numSteps: 8, // Pentagons have 6 children (5 + center), test a bit more
		},
		{
			name:     "hexagon parent res 5 to res 7",
			parent:   0x85080c3ffffffff, // res 5 hexagon
			childRes: 7,
			numSteps: 49, // Should have exactly 7^2 = 49 children
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize Go iterator
			var goIter IterCellsChildren
			iterInitParent(tt.parent, tt.childRes, &goIter)

			// Initialize C iterator
			var cIter IterCellsChildren
			iterInitParentC(tt.parent, tt.childRes, &cIter)

			// Compare initial state
			if goIter.H != cIter.H {
				t.Errorf("Initial H mismatch: Go=%x, C=%x", goIter.H, cIter.H)
			}
			if goIter.ParentRes != cIter.ParentRes {
				t.Errorf("Initial ParentRes mismatch: Go=%d, C=%d", goIter.ParentRes, cIter.ParentRes)
			}
			if goIter.SkipDigit != cIter.SkipDigit {
				t.Errorf("Initial SkipDigit mismatch: Go=%d, C=%d", goIter.SkipDigit, cIter.SkipDigit)
			}

			// Step through children and compare
			for i := 0; i < tt.numSteps; i++ {
				// Step both iterators
				iterStepChild(&goIter)
				iterStepChildC(&cIter)

				// Compare state after step
				if goIter.H != cIter.H {
					t.Errorf("Step %d: H mismatch: Go=%x, C=%x", i+1, goIter.H, cIter.H)
				}
				if goIter.ParentRes != cIter.ParentRes {
					t.Errorf("Step %d: ParentRes mismatch: Go=%d, C=%d", i+1, goIter.ParentRes, cIter.ParentRes)
				}
				if goIter.SkipDigit != cIter.SkipDigit {
					t.Errorf("Step %d: SkipDigit mismatch: Go=%d, C=%d", i+1, goIter.SkipDigit, cIter.SkipDigit)
				}

				// If we hit H3_NULL, we're done iterating
				if goIter.H == H3_NULL {
					break
				}
			}
		})
	}
}

// Test exhaustive iteration to ensure we generate the correct number of children
func Test_iterStepChild_exhaustive_parity(t *testing.T) {
	tests := []struct {
		name     string
		parent   H3Index
		childRes int32
	}{
		{
			name:     "hexagon res 0 to res 1",
			parent:   0x8001fffffffffff,
			childRes: 1,
		},
		{
			name:     "pentagon res 0 to res 1",
			parent:   0x8007fffffffffff,
			childRes: 1,
		},
		{
			name:     "hexagon res 2 to res 3",
			parent:   0x82080fffffffff,
			childRes: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Collect all children using Go implementation
			var goChildren []H3Index
			var goIter IterCellsChildren
			iterInitParent(tt.parent, tt.childRes, &goIter)
			for goIter.H != H3_NULL {
				goChildren = append(goChildren, goIter.H)
				iterStepChild(&goIter)
			}

			// Collect all children using C implementation
			var cChildren []H3Index
			var cIter IterCellsChildren
			iterInitParentC(tt.parent, tt.childRes, &cIter)
			for cIter.H != H3_NULL {
				cChildren = append(cChildren, cIter.H)
				iterStepChildC(&cIter)
			}

			// Compare counts
			if len(goChildren) != len(cChildren) {
				t.Errorf("Child count mismatch: Go=%d, C=%d", len(goChildren), len(cChildren))
			}

			// Compare each child (order should be identical)
			for i := 0; i < len(goChildren) && i < len(cChildren); i++ {
				if goChildren[i] != cChildren[i] {
					t.Errorf("Child %d mismatch: Go=%x, C=%x", i, goChildren[i], cChildren[i])
				}
			}

			t.Logf("Generated %d children (Go=%d, C=%d)", len(goChildren), len(goChildren), len(cChildren))
		})
	}
}
