//go:build cgo && c2go

package h3

import "testing"

func Test_ijkRotate60cw_parity(t *testing.T) {
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call C implementation
			gotC := _ijkRotate60cwC(&tt.coord)

			// Call Go implementation
			gotGo := tt.coord
			_ijkRotate60cw(&gotGo)

			// Compare results
			if gotGo.I != gotC.I || gotGo.J != gotC.J || gotGo.K != gotC.K {
				t.Errorf("_ijkRotate60cw() mismatch: Go{%d,%d,%d} != C{%d,%d,%d} for input{%d,%d,%d}",
					gotGo.I, gotGo.J, gotGo.K, gotC.I, gotC.J, gotC.K,
					tt.coord.I, tt.coord.J, tt.coord.K)
			}
		})
	}

	// Test that 6 rotations return to original (after normalization)
	t.Run("six_rotations", func(t *testing.T) {
		original := CoordIJK{1, 2, 3}
		coord := original

		// Apply 6 rotations
		for i := 0; i < 6; i++ {
			_ijkRotate60cw(&coord)
		}

		// Should be back to original (after normalization)
		var normalizedOriginal = original
		_ijkNormalize(&normalizedOriginal)

		if coord.I != normalizedOriginal.I || coord.J != normalizedOriginal.J || coord.K != normalizedOriginal.K {
			t.Errorf("Six rotations should return to normalized original: got{%d,%d,%d} want{%d,%d,%d}",
				coord.I, coord.J, coord.K, normalizedOriginal.I, normalizedOriginal.J, normalizedOriginal.K)
		}
	})

	// Test that clockwise and counter-clockwise are inverses
	t.Run("inverse_property", func(t *testing.T) {
		original := CoordIJK{3, 1, 4}

		// Apply cw then ccw
		coord1 := original
		_ijkRotate60cw(&coord1)
		_ijkRotate60ccw(&coord1)

		// Apply ccw then cw
		coord2 := original
		_ijkRotate60ccw(&coord2)
		_ijkRotate60cw(&coord2)

		// Both should return to normalized original
		var normalizedOriginal = original
		_ijkNormalize(&normalizedOriginal)

		if coord1.I != normalizedOriginal.I || coord1.J != normalizedOriginal.J || coord1.K != normalizedOriginal.K {
			t.Errorf("CW then CCW should return to normalized original: got{%d,%d,%d} want{%d,%d,%d}",
				coord1.I, coord1.J, coord1.K, normalizedOriginal.I, normalizedOriginal.J, normalizedOriginal.K)
		}

		if coord2.I != normalizedOriginal.I || coord2.J != normalizedOriginal.J || coord2.K != normalizedOriginal.K {
			t.Errorf("CCW then CW should return to normalized original: got{%d,%d,%d} want{%d,%d,%d}",
				coord2.I, coord2.J, coord2.K, normalizedOriginal.I, normalizedOriginal.J, normalizedOriginal.K)
		}
	})
}
