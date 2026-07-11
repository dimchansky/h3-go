//go:build cgo && c2go

package h3

import (
	"testing"
)

func Test_uncompactCells_parity(t *testing.T) {
	tests := []struct {
		name         string
		compactedSet []h3Index
		numCompacted int64
		numOut       int64
		res          int32
		desc         string
	}{
		{
			name:         "empty set",
			compactedSet: []h3Index{},
			numCompacted: 0,
			numOut:       0,
			res:          5,
			desc:         "Empty compacted set",
		},
		{
			name:         "single cell res 0 to res 1",
			compactedSet: []h3Index{0x8001fffffffffff},
			numCompacted: 1,
			numOut:       7,
			res:          1,
			desc:         "Single resolution 0 hexagon expanded to resolution 1",
		},
		{
			name:         "single cell res 0 to res 2",
			compactedSet: []h3Index{0x8001fffffffffff},
			numCompacted: 1,
			numOut:       49,
			res:          2,
			desc:         "Single resolution 0 hexagon expanded to resolution 2",
		},
		{
			name:         "single cell res 1 to res 2",
			compactedSet: []h3Index{0x8101fffffffffff},
			numCompacted: 1,
			numOut:       7,
			res:          2,
			desc:         "Single resolution 1 hexagon expanded to resolution 2",
		},
		{
			name:         "pentagon cell res 0 to res 1",
			compactedSet: []h3Index{0x8083fffffffffff}, // Base cell 4 (pentagon)
			numCompacted: 1,
			numOut:       7, // Pentagon has 6 children + itself logic -> need to check exact count
			res:          1,
			desc:         "Pentagon base cell expanded to resolution 1",
		},
		{
			name:         "multiple cells different types",
			compactedSet: []h3Index{0x8001fffffffffff, 0x8083fffffffffff},
			numCompacted: 2,
			numOut:       14, // 7 + 7 from above
			res:          1,
			desc:         "Mixed hexagon and pentagon base cells",
		},
		{
			name:         "with h3Null entries",
			compactedSet: []h3Index{0x8001fffffffffff, h3Null, 0x8083fffffffffff},
			numCompacted: 3,
			numOut:       14, // Should skip h3Null, same as above
			res:          1,
			desc:         "Compacted set with h3Null entries (should be skipped by hasChildAtRes check)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare input slices
			inputCompacted := make([]h3Index, len(tt.compactedSet))
			copy(inputCompacted, tt.compactedSet)

			// Prepare Go output slice
			goOutSet := make([]h3Index, tt.numOut)
			goErr := uncompactCells(inputCompacted, tt.numCompacted, goOutSet, tt.numOut, tt.res)

			// Prepare C output slice
			cOutSet := make([]h3Index, tt.numOut)
			cErrCode := uncompactCellsC(inputCompacted, tt.numCompacted, cOutSet, tt.numOut, tt.res)
			cErr := h3Error(cErrCode)

			// Compare errors
			if goErr != cErr {
				t.Errorf("Error mismatch: Go=%d, C=%d", goErr, cErr)
				return
			}

			if goErr != eSuccess {
				t.Logf("Expected error for %s: %d", tt.desc, goErr)
				return
			}

			// Compare the output sets - they should be identical in content and order
			mismatchFound := false
			for i := int64(0); i < tt.numOut; i++ {
				if goOutSet[i] != cOutSet[i] {
					t.Errorf("Output mismatch at index %d: Go=0x%x, C=0x%x", i, uint64(goOutSet[i]), uint64(cOutSet[i]))
					mismatchFound = true
					if i > 10 { // Limit output for readability
						break
					}
				}
			}

			if !mismatchFound {
				// Count non-null results for logging
				count := int64(0)
				for i := int64(0); i < tt.numOut; i++ {
					if goOutSet[i] != h3Null {
						count++
					}
				}
				t.Logf("Successfully uncompacted %d cells for %s", count, tt.desc)
			}
		})
	}
}

func Test_uncompactCells_invalid_input_parity(t *testing.T) {
	invalidCases := []struct {
		name         string
		compactedSet []h3Index
		numCompacted int64
		numOut       int64
		res          int32
		desc         string
	}{
		{
			name:         "invalid resolution domain",
			compactedSet: []h3Index{0x8001fffffffffff},
			numCompacted: 1,
			numOut:       7,
			res:          -1,
			desc:         "Negative resolution should cause domain error",
		},
		{
			name:         "resolution too high",
			compactedSet: []h3Index{0x8001fffffffffff},
			numCompacted: 1,
			numOut:       100,
			res:          20,
			desc:         "Resolution beyond maxH3Res should cause domain error",
		},
		{
			name:         "parent resolution higher than target",
			compactedSet: []h3Index{0x8301fffffffffff}, // Resolution 3
			numCompacted: 1,
			numOut:       10,
			res:          2, // Lower than parent
			desc:         "Target resolution lower than parent should cause mismatch error",
		},
		{
			name:         "output buffer too small",
			compactedSet: []h3Index{0x8001fffffffffff}, // Should expand to 7 cells
			numCompacted: 1,
			numOut:       3, // Too small
			res:          1,
			desc:         "Output buffer smaller than required should cause memory bounds error",
		},
		{
			name:         "invalid H3 cell",
			compactedSet: []h3Index{0x1001fffffffffff}, // Invalid mode bits
			numCompacted: 1,
			numOut:       10,
			res:          5,
			desc:         "Invalid H3 cell should cause validation error",
		},
	}

	for _, tt := range invalidCases {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare input slices
			inputCompacted := make([]h3Index, len(tt.compactedSet))
			copy(inputCompacted, tt.compactedSet)

			// Prepare Go output slice
			goOutSet := make([]h3Index, tt.numOut)
			goErr := uncompactCells(inputCompacted, tt.numCompacted, goOutSet, tt.numOut, tt.res)

			// Prepare C output slice
			cOutSet := make([]h3Index, tt.numOut)
			cErrCode := uncompactCellsC(inputCompacted, tt.numCompacted, cOutSet, tt.numOut, tt.res)
			cErr := h3Error(cErrCode)

			// Compare errors
			if goErr != cErr {
				t.Errorf("Error mismatch for %s: Go=%d, C=%d", tt.desc, goErr, cErr)
			}

			if goErr == eSuccess {
				t.Logf("Unexpected success for invalid input %s", tt.desc)
			} else {
				t.Logf("Expected error for %s: %d", tt.desc, goErr)
			}
		})
	}
}
