//go:build cgo && c2go

package h3

import "testing"

func Test_ijkRotate60ccw_parity(t *testing.T) {
	tests := []struct {
		name  string
		coord coordIJK
	}{
		{"origin", coordIJK{0, 0, 0}},
		{"unit i", coordIJK{1, 0, 0}},
		{"unit j", coordIJK{0, 1, 0}},
		{"unit k", coordIJK{0, 0, 1}},
		{"positive coords", coordIJK{1, 2, 3}},
		{"negative coords", coordIJK{-1, -2, -3}},
		{"mixed coords", coordIJK{2, -1, 3}},
		{"large coords", coordIJK{10, 20, 30}},
		{"asymmetric", coordIJK{1, 4, 2}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call C implementation
			gotC := _ijkRotate60ccwC(&tt.coord)

			// Call Go implementation
			gotGo := tt.coord
			_ijkRotate60ccw(&gotGo)

			// Compare results
			if gotGo.I != gotC.I || gotGo.J != gotC.J || gotGo.K != gotC.K {
				t.Errorf("_ijkRotate60ccw() mismatch: Go{%d,%d,%d} != C{%d,%d,%d} for input{%d,%d,%d}",
					gotGo.I, gotGo.J, gotGo.K, gotC.I, gotC.J, gotC.K,
					tt.coord.I, tt.coord.J, tt.coord.K)
			}
		})
	}

	// Test that 6 rotations return to original (after normalization)
	t.Run("six_rotations", func(t *testing.T) {
		original := coordIJK{1, 2, 3}
		coord := original

		// Apply 6 rotations
		for i := 0; i < 6; i++ {
			_ijkRotate60ccw(&coord)
		}

		// Should be back to original (after normalization)
		var normalizedOriginal = original
		_ijkNormalize(&normalizedOriginal)

		if coord.I != normalizedOriginal.I || coord.J != normalizedOriginal.J || coord.K != normalizedOriginal.K {
			t.Errorf("Six rotations should return to normalized original: got{%d,%d,%d} want{%d,%d,%d}",
				coord.I, coord.J, coord.K, normalizedOriginal.I, normalizedOriginal.J, normalizedOriginal.K)
		}
	})
}
