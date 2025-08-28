//go:build cgo

package h3

import (
	"testing"
)

func Test_uncompactCellsSize_parity(t *testing.T) {
	tests := []struct {
		name         string
		compactedSet []H3Index
		numCompacted int64
		res          int32
		desc         string
	}{
		{
			name:         "empty set",
			compactedSet: []H3Index{},
			numCompacted: 0,
			res:          5,
			desc:         "Empty compacted set",
		},
		{
			name:         "single cell res 0 to res 1",
			compactedSet: []H3Index{0x8001fffffffffff},
			numCompacted: 1,
			res:          1,
			desc:         "Single resolution 0 hexagon expanded to resolution 1",
		},
		{
			name:         "single cell res 0 to res 2",
			compactedSet: []H3Index{0x8001fffffffffff},
			numCompacted: 1,
			res:          2,
			desc:         "Single resolution 0 hexagon expanded to resolution 2",
		},
		{
			name:         "single cell res 1 to res 2",
			compactedSet: []H3Index{0x8101fffffffffff},
			numCompacted: 1,
			res:          2,
			desc:         "Single resolution 1 hexagon expanded to resolution 2",
		},
		{
			name:         "pentagon cell res 0 to res 1",
			compactedSet: []H3Index{0x8083fffffffffff}, // Base cell 4 (pentagon)
			numCompacted: 1,
			res:          1,
			desc:         "Pentagon base cell expanded to resolution 1",
		},
		{
			name:         "pentagon cell res 0 to res 2",
			compactedSet: []H3Index{0x8083fffffffffff}, // Base cell 4 (pentagon)
			numCompacted: 1,
			res:          2,
			desc:         "Pentagon base cell expanded to resolution 2",
		},
		{
			name:         "multiple cells different types",
			compactedSet: []H3Index{0x8001fffffffffff, 0x8083fffffffffff},
			numCompacted: 2,
			res:          1,
			desc:         "Mixed hexagon and pentagon base cells",
		},
		{
			name:         "with H3_NULL entries",
			compactedSet: []H3Index{0x8001fffffffffff, H3_NULL, 0x8083fffffffffff, H3_NULL},
			numCompacted: 4,
			res:          1,
			desc:         "Compacted set with H3_NULL entries (should be skipped)",
		},
		{
			name:         "higher resolution expansion",
			compactedSet: []H3Index{0x8001fffffffffff},
			numCompacted: 1,
			res:          3,
			desc:         "Resolution 0 to resolution 3 (large expansion)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare input slice with sufficient capacity
			inputSet := make([]H3Index, len(tt.compactedSet))
			copy(inputSet, tt.compactedSet)

			// Call Go implementation
			goSize, goErr := uncompactCellsSize(inputSet, tt.numCompacted, tt.res)

			// Call C implementation
			cSize, cErrCode := uncompactCellsSizeC(inputSet, tt.numCompacted, tt.res)
			cErr := H3Error(cErrCode)

			// Compare errors
			if goErr != cErr {
				t.Errorf("Error mismatch: Go=%d, C=%d", goErr, cErr)
				return
			}

			if goErr != E_SUCCESS {
				t.Logf("Expected error for %s: %d", tt.desc, goErr)
				return
			}

			// Compare sizes
			if goSize != cSize {
				t.Errorf("Size mismatch: Go=%d, C=%d", goSize, cSize)
				return
			}

			t.Logf("Successfully calculated size %d for %s", goSize, tt.desc)
		})
	}
}

func Test_uncompactCellsSize_invalid_input_parity(t *testing.T) {
	invalidCases := []struct {
		name         string
		compactedSet []H3Index
		numCompacted int64
		res          int32
		desc         string
	}{
		{
			name:         "invalid resolution domain",
			compactedSet: []H3Index{0x8001fffffffffff},
			numCompacted: 1,
			res:          -1,
			desc:         "Negative resolution should cause domain error",
		},
		{
			name:         "resolution too high",
			compactedSet: []H3Index{0x8001fffffffffff},
			numCompacted: 1,
			res:          20,
			desc:         "Resolution beyond MAX_H3_RES should cause domain error",
		},
		{
			name:         "parent resolution higher than target",
			compactedSet: []H3Index{0x8301fffffffffff}, // Resolution 3
			numCompacted: 1,
			res:          2, // Lower than parent
			desc:         "Target resolution lower than parent should cause mismatch error",
		},
		{
			name:         "invalid H3 cell",
			compactedSet: []H3Index{0x1001fffffffffff}, // Invalid mode bits
			numCompacted: 1,
			res:          5,
			desc:         "Invalid H3 cell should cause validation error",
		},
	}

	for _, tt := range invalidCases {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare input slice
			inputSet := make([]H3Index, len(tt.compactedSet))
			copy(inputSet, tt.compactedSet)

			// Call Go implementation
			goSize, goErr := uncompactCellsSize(inputSet, tt.numCompacted, tt.res)

			// Call C implementation
			cSize, cErrCode := uncompactCellsSizeC(inputSet, tt.numCompacted, tt.res)
			cErr := H3Error(cErrCode)

			// Compare errors
			if goErr != cErr {
				t.Errorf("Error mismatch for %s: Go=%d, C=%d", tt.desc, goErr, cErr)
			}

			if goErr == E_SUCCESS {
				t.Logf("Unexpected success for invalid input %s, sizes: Go=%d, C=%d", tt.desc, goSize, cSize)
			} else {
				t.Logf("Expected error for %s: %d", tt.desc, goErr)
			}
		})
	}
}
