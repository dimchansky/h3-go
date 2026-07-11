//go:build cgo && c2go

package h3

import "testing"

func Test_upAp7r_parity(t *testing.T) {
	tests := []struct {
		name  string
		coord coordIJK
	}{
		{"origin", coordIJK{0, 0, 0}},
		{"unit i", coordIJK{1, 0, 0}},
		{"unit j", coordIJK{0, 1, 0}},
		{"unit k", coordIJK{0, 0, 1}},
		{"positive coords", coordIJK{7, 14, 21}}, // multiples of 7
		{"mixed coords", coordIJK{3, -2, 5}},
		{"large coords", coordIJK{35, 28, 42}}, // larger multiples of 7
		{"negative coords", coordIJK{-7, -14, -21}},
		{"asymmetric", coordIJK{8, 15, 22}},
		{"small values", coordIJK{1, 2, 3}},
		{"aperture test", coordIJK{3, 1, 0}}, // from downr function
		{"aperture test 2", coordIJK{0, 3, 1}},
		{"aperture test 3", coordIJK{1, 0, 3}},
		{"normalized", coordIJK{2, 1, 0}},
		{"needs normalization", coordIJK{5, 3, 2}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call C implementation
			gotC := _upAp7rC(&tt.coord)

			// Call Go implementation
			gotGo := tt.coord
			_upAp7r(&gotGo)

			// Compare results
			if gotGo.I != gotC.I || gotGo.J != gotC.J || gotGo.K != gotC.K {
				t.Errorf("_upAp7r() mismatch: Go{%d,%d,%d} != C{%d,%d,%d} for input{%d,%d,%d}",
					gotGo.I, gotGo.J, gotGo.K, gotC.I, gotC.J, gotC.K,
					tt.coord.I, tt.coord.J, tt.coord.K)
			}
		})
	}

	// Test that transformation is deterministic
	t.Run("deterministic", func(t *testing.T) {
		coord := coordIJK{7, 14, 21}

		// Apply transformation twice
		result1 := coord
		_upAp7r(&result1)

		result2 := coord
		_upAp7r(&result2)

		if result1.I != result2.I || result1.J != result2.J || result1.K != result2.K {
			t.Errorf("_upAp7r should be deterministic: first{%d,%d,%d} != second{%d,%d,%d}",
				result1.I, result1.J, result1.K, result2.I, result2.J, result2.K)
		}
	})
}
