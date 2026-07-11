//go:build cgo && c2go

package h3

import (
	"testing"
)

func Test_gridDiskDistancesInternal_parity(t *testing.T) {
	tests := []struct {
		name     string
		origin   H3Index
		k        int32
		maxIdx   int64
		curK     int32
		wantErr  H3Error
		skipTest bool
		reason   string
	}{
		// Test basic functionality with k=0 (origin only)
		{
			name:   "k=0 origin only",
			origin: 0x871fb3e05ffffff, // Valid H3 index
			k:      0,
			maxIdx: 1,
			curK:   0,
		},
		// Test with k=1 (origin + 1-ring neighbors)
		{
			name:   "k=1 with neighbors",
			origin: 0x871fb3e05ffffff,
			k:      1,
			maxIdx: 7, // 1 origin + up to 6 neighbors
			curK:   0,
		},
		// Test with k=2 (origin + 2-ring)
		{
			name:   "k=2 larger ring",
			origin: 0x871fb3e05ffffff,
			k:      2,
			maxIdx: 19, // Formula: 3*k*(k+1)+1 = 3*2*3+1 = 19
			curK:   0,
		},
		// Test with pentagon base cell (should handle E_PENTAGON gracefully)
		{
			name:   "pentagon base cell",
			origin: 0x80428b7ffffffff, // Pentagon base cell 4
			k:      1,
			maxIdx: 7,
			curK:   0,
		},
		// Test intermediate recursion step (curK > 0)
		{
			name:   "intermediate recursion step",
			origin: 0x871fb3e05ffffff,
			k:      2,
			maxIdx: 19,
			curK:   1, // Starting from ring 1
		},
		// Test edge case: k equals curK (should return immediately)
		{
			name:   "k equals curK",
			origin: 0x871fb3e05ffffff,
			k:      1,
			maxIdx: 7,
			curK:   1, // Already at max distance
		},
		// Test with resolution 0 cell
		{
			name:   "resolution 0 cell",
			origin: 0x8029fffffffffff, // Base cell 41, resolution 0
			k:      1,
			maxIdx: 7,
			curK:   0,
		},
		// Test with high resolution cell
		{
			name:   "high resolution cell",
			origin: 0x8f1fb3e05b4b6db, // Resolution 15 cell
			k:      1,
			maxIdx: 7,
			curK:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skipTest {
				t.Skip(tt.reason)
			}

			// Initialize output arrays for Go implementation
			outGo := make([]H3Index, tt.maxIdx)
			distancesGo := make([]int32, tt.maxIdx)

			// Initialize output arrays for C implementation
			outC := make([]H3Index, tt.maxIdx)
			distancesC := make([]int32, tt.maxIdx)

			// Call Go implementation
			errGo := _gridDiskDistancesInternal(tt.origin, tt.k, outGo, distancesGo, tt.maxIdx, tt.curK)

			// Call C implementation
			errC := _gridDiskDistancesInternalC(tt.origin, tt.k, outC, distancesC, tt.maxIdx, tt.curK)

			// Compare error codes
			if errGo != errC {
				t.Errorf("Error mismatch: Go=%v, C=%v", errGo, errC)
			}

			// If we expected a specific error, check it
			if tt.wantErr != 0 && errGo != tt.wantErr {
				t.Errorf("Expected error %v, got %v", tt.wantErr, errGo)
			}

			// If both succeeded, compare results
			if errGo == E_SUCCESS && errC == E_SUCCESS {
				// Count non-zero entries in both arrays
				countGo := 0
				countC := 0
				for i := int64(0); i < tt.maxIdx; i++ {
					if outGo[i] != 0 {
						countGo++
					}
					if outC[i] != 0 {
						countC++
					}
				}

				if countGo != countC {
					t.Errorf("Different number of results: Go=%d, C=%d", countGo, countC)
				}

				// Convert to sets for comparison (order doesn't matter in hash table)
				setGo := make(map[H3Index]int32)
				setC := make(map[H3Index]int32)

				for i := int64(0); i < tt.maxIdx; i++ {
					if outGo[i] != 0 {
						setGo[outGo[i]] = distancesGo[i]
					}
					if outC[i] != 0 {
						setC[outC[i]] = distancesC[i]
					}
				}

				// Compare sets
				if len(setGo) != len(setC) {
					t.Errorf("Different set sizes: Go=%d, C=%d", len(setGo), len(setC))
				}

				for cell, distGo := range setGo {
					if distC, exists := setC[cell]; !exists {
						t.Errorf("Cell %x found in Go but not in C", cell)
					} else if distGo != distC {
						t.Errorf("Distance mismatch for cell %x: Go=%d, C=%d", cell, distGo, distC)
					}
				}

				for cell := range setC {
					if _, exists := setGo[cell]; !exists {
						t.Errorf("Cell %x found in C but not in Go", cell)
					}
				}
			}
		})
	}
}

// Test specific scenarios that might cause issues
func Test_gridDiskDistancesInternal_edge_cases_parity(t *testing.T) {
	tests := []struct {
		name     string
		origin   H3Index
		k        int32
		maxIdx   int64
		curK     int32
		skipTest bool
		reason   string
	}{
		// Test with curK > k (should return immediately)
		{
			name:   "curK greater than k",
			origin: 0x871fb3e05ffffff,
			k:      1,
			maxIdx: 7,
			curK:   2, // curK > k
		},
		// Test with minimal hash table size
		{
			name:   "minimal hash table",
			origin: 0x871fb3e05ffffff,
			k:      0,
			maxIdx: 1, // Just enough space for origin
			curK:   0,
		},
		// Test that could cause hash collisions but with adequate space
		{
			name:   "hash collisions with adequate space",
			origin: 0x871fb3e05ffffff,
			k:      1,
			maxIdx: 10, // Small table but adequate for k=1 (max 7 results)
			curK:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skipTest {
				t.Skip(tt.reason)
			}

			// Initialize output arrays
			outGo := make([]H3Index, tt.maxIdx)
			distancesGo := make([]int32, tt.maxIdx)
			outC := make([]H3Index, tt.maxIdx)
			distancesC := make([]int32, tt.maxIdx)

			// Call both implementations
			errGo := _gridDiskDistancesInternal(tt.origin, tt.k, outGo, distancesGo, tt.maxIdx, tt.curK)
			errC := _gridDiskDistancesInternalC(tt.origin, tt.k, outC, distancesC, tt.maxIdx, tt.curK)

			// Compare results
			if errGo != errC {
				t.Errorf("Error mismatch: Go=%v, C=%v", errGo, errC)
			}

			// Both should handle edge cases the same way
			if errGo == E_SUCCESS && errC == E_SUCCESS {
				// Basic sanity check - both should have same number of non-zero entries
				countGo := 0
				countC := 0
				for i := int64(0); i < tt.maxIdx; i++ {
					if outGo[i] != 0 {
						countGo++
					}
					if outC[i] != 0 {
						countC++
					}
				}

				if countGo != countC {
					t.Errorf("Different result counts: Go=%d, C=%d", countGo, countC)
				}
			}
		})
	}
}
