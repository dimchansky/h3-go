//go:build cgo

package c2go

import "testing"

func Test_firstOneIndex_parity(t *testing.T) {
	tests := []struct {
		name string
		h    H3Index
	}{
		{"zero", 0},
		{"single bit 0", 1},
		{"single bit 1", 2},
		{"single bit 2", 4},
		{"single bit 63", 1 << 63},
		{"multiple bits low", 0xFF},
		{"multiple bits high", 0xFF00000000000000},
		{"mixed bits", 0x123456789ABCDEF0},
		{"all bits", 0xFFFFFFFFFFFFFFFF},
		
		// H3 index examples  
		{"valid h3 index", 0x8a1fb46622dffff},
		{"another h3 index", 0x8928308280fffff},
		{"pentagon index", 0x821c07fffffffff},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call C implementation
			gotC := firstOneIndexC(tt.h)
			
			// Call Go implementation
			gotGo := _firstOneIndex(tt.h)
			
			// Compare results
			if gotGo != gotC {
				t.Errorf("_firstOneIndex() mismatch: Go=%d != C=%d for h=%x", 
					gotGo, gotC, tt.h)
			}
		})
	}
	
	// Test edge case where h=0 should return -1 for Go, but C behavior may vary
	t.Run("zero_edge_case", func(t *testing.T) {
		h := H3Index(0)
		gotC := firstOneIndexC(h)
		gotGo := _firstOneIndex(h)
		
		// Document the C behavior for h=0
		t.Logf("_firstOneIndex(0): Go=%d, C=%d", gotGo, gotC)
		
		// For h=0, Go returns -1, C may return different values depending on implementation
		// We'll just ensure they're consistent for our implementation
		if gotGo != gotC {
			// If C returns something different, document it but allow the difference
			t.Logf("Note: C and Go differ for h=0 - this may be expected")
		}
	})
}