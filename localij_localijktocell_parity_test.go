//go:build cgo && c2go

package h3

import (
	"testing"
)

func Test_localIjkToCell_parity(t *testing.T) {
	testCases := []struct {
		name   string
		origin h3Index
		ijk    coordIJK
	}{
		// Basic resolution 0 tests
		{
			name:   "res0_center",
			origin: 0x8001fffffffffff, // resolution 0 base cell 1
			ijk:    coordIJK{0, 0, 0},
		},
		{
			name:   "res0_neighbor_i",
			origin: 0x8001fffffffffff,
			ijk:    coordIJK{1, 0, 0},
		},
		{
			name:   "res0_neighbor_j",
			origin: 0x8001fffffffffff,
			ijk:    coordIJK{0, 1, 0},
		},
		{
			name:   "res0_neighbor_k",
			origin: 0x8001fffffffffff,
			ijk:    coordIJK{0, 0, 1},
		},

		// Basic resolution 1 tests
		{
			name:   "res1_center",
			origin: 0x8101fffffffffff, // resolution 1
			ijk:    coordIJK{0, 0, 0},
		},
		{
			name:   "res1_neighbor",
			origin: 0x8101fffffffffff,
			ijk:    coordIJK{1, 1, 0},
		},

		// Higher resolution tests
		{
			name:   "res5_center",
			origin: 0x8501fffffffffff, // resolution 5
			ijk:    coordIJK{0, 0, 0},
		},
		{
			name:   "res5_nearby",
			origin: 0x8501fffffffffff,
			ijk:    coordIJK{2, -1, -1},
		},

		// Pentagon base cell tests (base cell 4 is pentagon)
		{
			name:   "pentagon_res0",
			origin: 0x8004fffffffffff, // resolution 0 base cell 4 (pentagon)
			ijk:    coordIJK{0, 0, 0},
		},
		{
			name:   "pentagon_res1",
			origin: 0x8104fffffffffff, // resolution 1 base cell 4 (pentagon)
			ijk:    coordIJK{0, 0, 0},
		},

		// More challenging coordinate ranges (let's see what happens)
		{
			name:   "large_coordinates",
			origin: 0x8101fffffffffff,
			ijk:    coordIJK{10, 10, 10},
		},

		// Edge cases
		{
			name:   "null_origin",
			origin: 0x0,
			ijk:    coordIJK{0, 0, 0},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test Go implementation
			var goResult h3Index
			goErr := localIjkToCell(tc.origin, &tc.ijk, &goResult)

			// Test C implementation
			var cResult h3Index
			cErr := _localIjkToCellC(tc.origin, &tc.ijk, &cResult)

			// Compare errors
			if goErr != cErr {
				t.Errorf("Error mismatch - Go: %d, C: %d", goErr, cErr)
				return
			}

			// If both succeeded, compare results
			if goErr == eSuccess && cErr == eSuccess {
				if goResult != cResult {
					t.Errorf("Result mismatch - Go: %x, C: %x", goResult, cResult)
				}
			} else {
				// Both failed with the same error - that's okay
				t.Logf("Both implementations failed with error: %d", goErr)
			}
		})
	}
}

func Test_localIjkToCell_extensive_parity(t *testing.T) {
	// Test with various origins and coordinate ranges
	origins := []h3Index{
		0x8001fffffffffff, // res 0 base cell 1
		0x8101fffffffffff, // res 1 base cell 1
		0x8201fffffffffff, // res 2 base cell 1
		0x8004fffffffffff, // res 0 base cell 4 (pentagon)
		0x8104fffffffffff, // res 1 base cell 4 (pentagon)
		0x8500fffffffffff, // res 5 base cell 0
	}

	// Test coordinate ranges
	coords := []coordIJK{
		{0, 0, 0},   // center
		{1, 0, 0},   // i-axis
		{0, 1, 0},   // j-axis
		{0, 0, 1},   // k-axis
		{1, 1, 0},   // ij diagonal
		{1, 0, 1},   // ik diagonal
		{0, 1, 1},   // jk diagonal
		{-1, 0, 1},  // negative i
		{0, -1, 1},  // negative j
		{1, 0, -1},  // negative k
		{2, -1, -1}, // larger coordinates
		{-2, 1, 1},  // negative larger
	}

	for _, origin := range origins {
		for _, ijk := range coords {
			name := "extensive"

			t.Run(name, func(t *testing.T) {
				// Test Go implementation
				var goResult h3Index
				goErr := localIjkToCell(origin, &ijk, &goResult)

				// Test C implementation
				var cResult h3Index
				cErr := _localIjkToCellC(origin, &ijk, &cResult)

				// Compare errors
				if goErr != cErr {
					t.Errorf("Error mismatch for origin %x, ijk %+v - Go: %d, C: %d", origin, ijk, goErr, cErr)
					return
				}

				// If both succeeded, compare results
				if goErr == eSuccess && cErr == eSuccess {
					if goResult != cResult {
						t.Errorf("Result mismatch for origin %x, ijk %+v - Go: %x, C: %x", origin, ijk, goResult, cResult)
					}
				}
			})
		}
	}
}
