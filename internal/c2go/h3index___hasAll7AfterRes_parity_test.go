//go:build cgo

package c2go

import "testing"

func Test_hasAll7AfterRes_parity(t *testing.T) {
	tests := []struct {
		name string
		h    H3Index
		res  int32
	}{
		// Valid H3 indexes with proper 7s after resolution
		{"valid h3 res 9", 0x8a1fb46622dffff, 9},
		{"valid h3 res 5", 0x8a1fb46622dffff, 5},
		{"pentagon res 5", 0x821c07fffffffff, 5},
		{"res 0", 0x8001fffffffffff, 0},

		// Resolution 15 cases (all digits used)
		{"res 15", 0x8f01fffffffffff, 15},
		{"res 15 any index", 0x8a1fb46622dffff, 15},

		// Cases where digits after res are not all 7s
		{"res 5 not all 7s", 0x8a1fb46622d0000, 5}, // Has 0s after res 5
		{"res 10 mixed", 0x8a1fb46622d1234, 10},    // Mixed digits after res 10

		// Edge cases
		{"zero index", 0, 5},
		{"all bits set", 0xFFFFFFFFFFFFFFFF, 10},

		// Test different resolution levels
		{"res 1", 0x8001fffffffffff, 1},
		{"res 14", 0x8f01fffffffffff, 14},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call C implementation
			gotC := hasAll7AfterResC(tt.h, tt.res)

			// Call Go implementation
			gotGo := _hasAll7AfterRes(tt.h, tt.res)

			// Compare results
			if gotGo != gotC {
				t.Errorf("_hasAll7AfterRes(%x, %d) mismatch: Go=%v != C=%v",
					tt.h, tt.res, gotGo, gotC)
			}
		})
	}

	// Test invalid resolution values
	t.Run("invalid_resolutions", func(t *testing.T) {
		h := H3Index(0x8a1fb46622dffff)

		testCases := []struct {
			res  int32
			name string
		}{
			{-1, "negative"},
			{16, "too_large"},
			{100, "way_too_large"},
		}

		for _, tc := range testCases {
			gotC := hasAll7AfterResC(h, tc.res)
			gotGo := _hasAll7AfterRes(h, tc.res)

			if gotGo != gotC {
				t.Errorf("_hasAll7AfterRes(%x, %d) mismatch for %s: Go=%v != C=%v",
					h, tc.res, tc.name, gotGo, gotC)
			}
		}
	})

	// Test all resolutions 0-15
	t.Run("all_resolutions", func(t *testing.T) {
		h := H3Index(0x8a1fb46622dffff)

		for res := int32(0); res <= 15; res++ {
			gotC := hasAll7AfterResC(h, res)
			gotGo := _hasAll7AfterRes(h, res)

			if gotGo != gotC {
				t.Errorf("_hasAll7AfterRes(%x, %d) mismatch at res %d: Go=%v != C=%v",
					h, res, res, gotGo, gotC)
			}
		}
	})
}
