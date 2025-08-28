//go:build cgo

package h3

import (
	"testing"
)

func Test_cellToChildren_parity(t *testing.T) {
	tests := []struct {
		name     string
		h        H3Index
		childRes int32
	}{
		{
			name:     "hexagon res 0 to res 1",
			h:        0x8001fffffffffff,
			childRes: 1,
		},
		{
			name:     "hexagon res 3 to res 4",
			h:        0x83080dfffffff,
			childRes: 4,
		},
		{
			name:     "pentagon res 0 to res 1",
			h:        0x8007fffffffffff,
			childRes: 1,
		},
		{
			name:     "hexagon res 5 to res 7",
			h:        0x85080c3ffffffff,
			childRes: 7,
		},
		{
			name:     "pentagon res 2 to res 3",
			h:        0x82026ffffffffff,
			childRes: 3,
		},
		{
			name:     "hexagon res 8 to res 9",
			h:        0x8830062838bffff,
			childRes: 9,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Get the number of children
			childrenSize, err := cellToChildrenSize(tt.h, tt.childRes)
			if err != E_SUCCESS {
				t.Fatalf("cellToChildrenSize failed: %v", err)
			}

			// Allocate arrays for both implementations
			goChildren := make([]H3Index, childrenSize)
			cChildren := make([]H3Index, childrenSize)

			// Call Go implementation
			goErr := cellToChildren(tt.h, tt.childRes, goChildren)
			if goErr != E_SUCCESS {
				t.Fatalf("Go cellToChildren failed: %v", goErr)
			}

			// Call C implementation
			cErr := cellToChildrenC(tt.h, tt.childRes, cChildren)
			if cErr != 0 {
				t.Fatalf("C cellToChildren failed: %v", cErr)
			}

			// Compare results
			if len(goChildren) != len(cChildren) {
				t.Errorf("Children count mismatch: Go=%d, C=%d", len(goChildren), len(cChildren))
			}

			// Compare each child - they should be in the same order
			for i := 0; i < len(goChildren) && i < len(cChildren); i++ {
				if goChildren[i] != cChildren[i] {
					t.Errorf("Child %d mismatch: Go=%016x, C=%016x", i, goChildren[i], cChildren[i])
				}
			}

			// Verify all children are valid
			for i, child := range goChildren {
				if child == H3_NULL {
					t.Errorf("Go implementation returned H3_NULL at index %d", i)
				}
				if getResolution(child) != tt.childRes {
					t.Errorf("Go child %d has wrong resolution: expected %d, got %d",
						i, tt.childRes, getResolution(child))
				}
			}

			t.Logf("Generated %d children for parent %016x at res %d",
				childrenSize, tt.h, tt.childRes)
		})
	}
}

// Test edge cases
func Test_cellToChildren_edge_cases_parity(t *testing.T) {
	tests := []struct {
		name     string
		h        H3Index
		childRes int32
	}{
		{
			name:     "same resolution (should have 1 child - itself)",
			h:        0x8001fffffffffff,
			childRes: 0,
		},
		{
			name:     "res 14 to res 15 (max resolution)",
			h:        0x8e1fb46622d691f, // Valid res 14 index (from rotate tests)
			childRes: 15,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Get the number of children
			childrenSize, err := cellToChildrenSize(tt.h, tt.childRes)
			if err != E_SUCCESS {
				// If cellToChildrenSize fails, both implementations should fail
				t.Logf("Expected failure for invalid resolution: %v", err)
				return
			}

			// Allocate arrays for both implementations
			goChildren := make([]H3Index, childrenSize)
			cChildren := make([]H3Index, childrenSize)

			// Call both implementations
			goErr := cellToChildren(tt.h, tt.childRes, goChildren)
			cErr := cellToChildrenC(tt.h, tt.childRes, cChildren)

			// Both should succeed or both should fail
			if (goErr == E_SUCCESS) != (cErr == 0) {
				t.Errorf("Error mismatch: Go=%v, C=%v", goErr, cErr)
			}

			if goErr == E_SUCCESS {
				// Compare results
				for i := 0; i < len(goChildren) && i < len(cChildren); i++ {
					if goChildren[i] != cChildren[i] {
						t.Errorf("Child %d mismatch: Go=%016x, C=%016x",
							i, goChildren[i], cChildren[i])
					}
				}
			}
		})
	}
}
