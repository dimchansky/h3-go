//go:build cgo

package h3

import "testing"

func Test_hasGoodTopBits_parity(t *testing.T) {
	tests := []struct {
		name string
		h    H3Index
	}{
		// Valid H3 indexes (should have good top bits)
		{"valid h3 index", 0x8a1fb46622dffff},
		{"another valid h3", 0x8928308280fffff},
		{"pentagon index", 0x821c07fffffffff},
		{"res 0 index", 0x8001fffffffffff},
		{"res 15 index", 0x8f01fffffffffff},

		// Invalid patterns (should fail top bits test)
		{"zero", 0},
		{"no mode bits", 0x0a1fb46622dffff},
		{"wrong mode", 0x4a1fb46622dffff},        // mode 2 instead of 1
		{"high bit set", 0xca1fb46622dffff},      // high bit set
		{"reserved bits set", 0x8e1fb46622dffff}, // reserved bits set
		{"all bits set", 0xFFFFFFFFFFFFFFFF},

		// Edge cases - specific bit patterns
		{"only mode bit", 0x1000000000000000},   // 0_0010_000 pattern
		{"correct pattern", 0x0800000000000000}, // 0_0001_000 pattern
		{"mode + reserved", 0x1800000000000000}, // 0_0011_000 pattern
		{"high + mode", 0x8800000000000000},     // 1_0001_000 pattern
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call C implementation
			gotC := hasGoodTopBitsC(tt.h)

			// Call Go implementation
			gotGo := _hasGoodTopBits(tt.h)

			// Compare results
			if gotGo != gotC {
				t.Errorf("_hasGoodTopBits(%x) mismatch: Go=%v != C=%v",
					tt.h, gotGo, gotC)
			}
		})
	}

	// Test the exact bit pattern that should pass
	t.Run("expected_pattern", func(t *testing.T) {
		// Create an index with exactly the right top 8 bits: 0b00001000
		h := H3Index(0x0800000000000000) // Top 8 bits = 0b00001000

		gotC := hasGoodTopBitsC(h)
		gotGo := _hasGoodTopBits(h)

		if gotGo != gotC {
			t.Errorf("_hasGoodTopBits(%x) mismatch: Go=%v != C=%v", h, gotGo, gotC)
		}

		// Both should return true for the correct pattern
		if !gotGo {
			t.Errorf("Expected _hasGoodTopBits(%x) to return true (Go), got %v", h, gotGo)
		}
		if !gotC {
			t.Errorf("Expected _hasGoodTopBits(%x) to return true (C), got %v", h, gotC)
		}
	})
}
