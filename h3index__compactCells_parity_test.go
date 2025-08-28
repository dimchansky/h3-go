//go:build cgo

package h3

import (
	"testing"
)

func Test_compactCells_parity(t *testing.T) {
	tests := []struct {
		name   string
		h3Set  []H3Index
		numHex int64
		desc   string
	}{
		{
			name:   "empty set",
			h3Set:  []H3Index{},
			numHex: 0,
			desc:   "Empty input set",
		},
		{
			name:   "single hexagon res 0",
			h3Set:  []H3Index{0x8001fffffffffff},
			numHex: 1,
			desc:   "Single resolution 0 hexagon (no compaction possible)",
		},
		{
			name:   "single hexagon res 1",
			h3Set:  []H3Index{0x8101fffffffffff},
			numHex: 1,
			desc:   "Single resolution 1 hexagon",
		},
		{
			name:   "resolution 2 siblings",
			h3Set:  []H3Index{0x8220000ffffffff, 0x8220001ffffffff, 0x8220002ffffffff, 0x8220003ffffffff, 0x8220004ffffffff, 0x8220005ffffffff, 0x8220006ffffffff},
			numHex: 7,
			desc:   "Seven resolution 2 children that should compact to parent",
		},
		{
			name:   "mixed resolutions (invalid for compacting)",
			h3Set:  []H3Index{0x8001fffffffffff, 0x8101fffffffffff},
			numHex: 2,
			desc:   "Mixed resolution cells (should fail or handle gracefully)",
		},
		{
			name:   "pentagon children",
			h3Set:  []H3Index{0x81083ffffffffff, 0x81087ffffffffff, 0x8108bffffffffff, 0x8108fffffffffff, 0x81093ffffffffff, 0x81097ffffffffff},
			numHex: 6,
			desc:   "Pentagon children that should compact (only 6 children for pentagon)",
		},
		{
			name:   "partial set - no compaction",
			h3Set:  []H3Index{0x8220000ffffffff, 0x8220001ffffffff, 0x8220002ffffffff},
			numHex: 3,
			desc:   "Incomplete set of children, no compaction should occur",
		},
		{
			name:   "duplicate cells",
			h3Set:  []H3Index{0x8220000ffffffff, 0x8220000ffffffff},
			numHex: 2,
			desc:   "Duplicate input cells should return E_DUPLICATE_INPUT",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Skip if numHex is 0 but h3Set is not empty (for safety)
			if tt.numHex == 0 && len(tt.h3Set) > 0 {
				t.Skip("Skipping test with empty numHex but non-empty h3Set")
			}

			// Prepare Go compacted result
			goCompactedSet := make([]H3Index, max(tt.numHex, 1)) // Ensure at least size 1
			goErr := compactCells(tt.h3Set, goCompactedSet, tt.numHex)

			// Prepare C compacted result
			cCompactedSet := make([]H3Index, max(tt.numHex, 1)) // Ensure at least size 1
			cErr := H3Error(compactCellsC(tt.h3Set, cCompactedSet, tt.numHex))

			// Compare errors
			if goErr != cErr {
				t.Errorf("Error mismatch: Go=%d, C=%d", goErr, cErr)
				return
			}

			if goErr != E_SUCCESS {
				t.Logf("Expected error for %s: %d", tt.desc, goErr)
				return
			}

			// For successful compactions, compare the actual results
			// Both functions modify the array in place, so we need to count non-zero elements
			goCount := int64(0)
			for i := int64(0); i < tt.numHex; i++ {
				if goCompactedSet[i] != 0 {
					goCount++
				}
			}

			cCount := int64(0)
			for i := int64(0); i < tt.numHex; i++ {
				if cCompactedSet[i] != 0 {
					cCount++
				}
			}

			// Compare count of compacted cells
			if goCount != cCount {
				t.Errorf("Compacted count mismatch: Go=%d, C=%d", goCount, cCount)
				return
			}

			// Compare the actual compacted cells
			// Since the order might differ, we'll check that each Go result exists in C result
			for i := int64(0); i < tt.numHex; i++ {
				if goCompactedSet[i] == 0 {
					continue // Skip null entries
				}
				found := false
				for j := int64(0); j < tt.numHex; j++ {
					if goCompactedSet[i] == cCompactedSet[j] {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Go result cell 0x%x not found in C results", uint64(goCompactedSet[i]))
				}
			}

			// Also check that each C result exists in Go result
			for i := int64(0); i < tt.numHex; i++ {
				if cCompactedSet[i] == 0 {
					continue // Skip null entries
				}
				found := false
				for j := int64(0); j < goCount; j++ {
					if cCompactedSet[i] == goCompactedSet[j] {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("C result cell 0x%x not found in Go results", uint64(cCompactedSet[i]))
				}
			}

			t.Logf("Successfully compacted %d cells to %d cells for %s", tt.numHex, goCount, tt.desc)
		})
	}
}

func Test_compactCells_invalid_input_parity(t *testing.T) {
	invalidCases := []struct {
		name   string
		h3Set  []H3Index
		numHex int64
		desc   string
	}{
		{
			name:   "invalid cell in set",
			h3Set:  []H3Index{0x1001fffffffffff}, // Invalid mode bits
			numHex: 1,
			desc:   "Set contains invalid H3 cell",
		},
		{
			name:   "reserved bits set",
			h3Set:  []H3Index{0x8201fffffffffff | (1 << 56)}, // Reserved bits set
			numHex: 1,
			desc:   "Cell with reserved bits already set",
		},
	}

	for _, tt := range invalidCases {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare Go compacted result
			goCompactedSet := make([]H3Index, max(tt.numHex, 1))
			goErr := compactCells(tt.h3Set, goCompactedSet, tt.numHex)

			// Prepare C compacted result
			cCompactedSet := make([]H3Index, max(tt.numHex, 1))
			cErr := H3Error(compactCellsC(tt.h3Set, cCompactedSet, tt.numHex))

			// Compare errors
			if goErr != cErr {
				t.Errorf("Error mismatch for %s: Go=%d, C=%d", tt.desc, goErr, cErr)
			}

			if goErr == E_SUCCESS {
				t.Logf("Unexpected success for invalid input %s", tt.desc)
			} else {
				t.Logf("Expected error for %s: %d", tt.desc, goErr)
			}
		})
	}
}

// max returns the maximum of two int64 values
func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
