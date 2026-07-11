//go:build cgo && c2go

package h3

import "testing"

func Test_cubeToIjk_parity(t *testing.T) {
	tests := []struct {
		name  string
		coord coordIJK
	}{
		{"origin", coordIJK{0, 0, 0}},
		{"cube unit", coordIJK{1, -1, 0}},          // Valid cube coordinate (sums to 0)
		{"cube unit 2", coordIJK{0, 1, -1}},        // Valid cube coordinate
		{"cube unit 3", coordIJK{-1, 0, 1}},        // Valid cube coordinate
		{"positive cube", coordIJK{3, -1, -2}},     // Valid cube coordinate
		{"negative cube", coordIJK{-2, 1, 1}},      // Valid cube coordinate
		{"invalid cube", coordIJK{1, 2, 3}},        // Invalid cube (doesn't sum to 0)
		{"large cube", coordIJK{10, -5, -5}},       // Valid cube coordinate
		{"asymmetric cube", coordIJK{7, -3, -4}},   // Valid cube coordinate
		{"negative invalid", coordIJK{-1, -2, -3}}, // Invalid cube
		{"mixed invalid", coordIJK{2, -1, 3}},      // Invalid cube
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call C implementation
			gotC := cubeToIjkC(&tt.coord)

			// Call Go implementation
			gotGo := tt.coord
			cubeToIjk(&gotGo)

			// Compare results
			if gotGo.I != gotC.I || gotGo.J != gotC.J || gotGo.K != gotC.K {
				t.Errorf("cubeToIjk() mismatch: Go{%d,%d,%d} != C{%d,%d,%d} for input{%d,%d,%d}",
					gotGo.I, gotGo.J, gotGo.K, gotC.I, gotC.J, gotC.K,
					tt.coord.I, tt.coord.J, tt.coord.K)
			}
		})
	}

	// Test that transformation is deterministic
	t.Run("deterministic", func(t *testing.T) {
		coord := coordIJK{3, -1, -2}

		// Apply transformation twice
		result1 := coord
		cubeToIjk(&result1)

		result2 := coord
		cubeToIjk(&result2)

		if result1.I != result2.I || result1.J != result2.J || result1.K != result2.K {
			t.Errorf("cubeToIjk should be deterministic: first{%d,%d,%d} != second{%d,%d,%d}",
				result1.I, result1.J, result1.K, result2.I, result2.J, result2.K)
		}
	})
}
