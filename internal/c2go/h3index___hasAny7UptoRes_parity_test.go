//go:build cgo

package c2go

import "testing"

func Test_hasAny7UptoRes_parity(t *testing.T) {
	tests := []struct {
		name string
		h    H3Index
		res  int
	}{
		// Valid H3 indexes (should not have 7s in valid digits)
		{"valid h3 res 9", 0x8a1fb46622dffff, 9},
		{"valid h3 res 5", 0x8a1fb46622dffff, 5},
		{"pentagon res 5", 0x821c07fffffffff, 5},
		{"res 0", 0x8001fffffffffff, 0},
		{"res 15", 0x8f01fffffffffff, 15},

		// Indexes with 7s in digits (should detect invalid digits)
		{"has 7 at res 1", 0x81c1fffffffffff, 1},  // digit 1 = 7
		{"has 7 at res 2", 0x8201c7ffffffffff, 2}, // digit 2 = 7
		{"has 7 at res 10", 0x8a1fb46622dffff | (7 << (3 * (15 - 10))), 10},

		// Edge cases
		{"zero index", 0, 5},
		{"all bits set", 0xFFFFFFFFFFFFFFFF, 10},
		{"res 0 check", 0x8001fffffffffff, 0}, // Should not check any digits beyond base cell

		// Test different resolution levels
		{"res 1", 0x8001fffffffffff, 1},
		{"res 14", 0x8f01fffffffffff, 14},
		{"res 15", 0x8f01fffffffffff, 15},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call C implementation
			gotC := hasAny7UptoResC(tt.h, tt.res)

			// Call Go implementation
			gotGo := _hasAny7UptoRes(tt.h, tt.res)

			// Compare results
			if gotGo != gotC {
				t.Errorf("_hasAny7UptoRes(%x, %d) mismatch: Go=%v != C=%v",
					tt.h, tt.res, gotGo, gotC)
			}
		})
	}

	// Test invalid resolution values
	t.Run("invalid_resolutions", func(t *testing.T) {
		h := H3Index(0x8a1fb46622dffff)

		testCases := []struct {
			res  int
			name string
		}{
			{-1, "negative"},
			{16, "too_large"},
			// Skip very large values like 100 due to C undefined behavior
		}

		for _, tc := range testCases {
			gotC := hasAny7UptoResC(h, tc.res)
			gotGo := _hasAny7UptoRes(h, tc.res)

			if gotGo != gotC {
				t.Errorf("_hasAny7UptoRes(%x, %d) mismatch for %s: Go=%v != C=%v",
					h, tc.res, tc.name, gotGo, gotC)
			}
		}
	})

	// Test all resolutions 0-15 with valid H3 index
	t.Run("all_resolutions", func(t *testing.T) {
		h := H3Index(0x8a1fb46622dffff)

		for res := 0; res <= 15; res++ {
			gotC := hasAny7UptoResC(h, res)
			gotGo := _hasAny7UptoRes(h, res)

			if gotGo != gotC {
				t.Errorf("_hasAny7UptoRes(%x, %d) mismatch at res %d: Go=%v != C=%v",
					h, res, res, gotGo, gotC)
			}
		}
	})
}
