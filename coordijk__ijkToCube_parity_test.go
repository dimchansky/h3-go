//go:build cgo

package h3

import "testing"

func Test_ijkToCube_parity(t *testing.T) {
	tests := []struct {
		name  string
		coord CoordIJK
	}{
		{"origin", CoordIJK{0, 0, 0}},
		{"unit i", CoordIJK{1, 0, 0}},
		{"unit j", CoordIJK{0, 1, 0}},
		{"unit k", CoordIJK{0, 0, 1}},
		{"positive coords", CoordIJK{1, 2, 3}},
		{"negative coords", CoordIJK{-1, -2, -3}},
		{"mixed coords", CoordIJK{2, -1, 3}},
		{"large coords", CoordIJK{10, 20, 30}},
		{"asymmetric", CoordIJK{1, 4, 2}},
		{"all equal", CoordIJK{5, 5, 5}},
		{"zeros with non-zero", CoordIJK{0, 0, 7}},
		{"negative zeros", CoordIJK{-3, 0, 0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call C implementation
			gotC := ijkToCubeC(&tt.coord)

			// Call Go implementation
			gotGo := tt.coord
			ijkToCube(&gotGo)

			// Compare results
			if gotGo.I != gotC.I || gotGo.J != gotC.J || gotGo.K != gotC.K {
				t.Errorf("ijkToCube() mismatch: Go{%d,%d,%d} != C{%d,%d,%d} for input{%d,%d,%d}",
					gotGo.I, gotGo.J, gotGo.K, gotC.I, gotC.J, gotC.K,
					tt.coord.I, tt.coord.J, tt.coord.K)
			}
		})
	}

	// Test that transformation is deterministic
	t.Run("deterministic", func(t *testing.T) {
		coord := CoordIJK{3, 1, 4}

		// Apply transformation twice
		result1 := coord
		ijkToCube(&result1)

		result2 := coord
		ijkToCube(&result2)

		if result1.I != result2.I || result1.J != result2.J || result1.K != result2.K {
			t.Errorf("ijkToCube should be deterministic: first{%d,%d,%d} != second{%d,%d,%d}",
				result1.I, result1.J, result1.K, result2.I, result2.J, result2.K)
		}
	})
}
