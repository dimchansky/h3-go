//go:build c2go

package c2go

import "testing"

func Test_ijkToHex2d_parity(t *testing.T) {
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
		{"normalized coords", CoordIJK{5, 3, 1}},
		{"zero i", CoordIJK{0, 5, 2}},
		{"zero j", CoordIJK{3, 0, 1}},
		{"zero k", CoordIJK{2, 4, 0}},
	}

	const tol = 1e-14

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call C implementation
			gotC := _ijkToHex2dC(&tt.coord)

			// Call Go implementation
			var gotGo Vec2d
			_ijkToHex2d(&tt.coord, &gotGo)

			// Compare results with tolerance
			if absFloat(gotGo.X-gotC.X) > tol || absFloat(gotGo.Y-gotC.Y) > tol {
				t.Errorf("_ijkToHex2d() mismatch: Go{%.15g,%.15g} != C{%.15g,%.15g} for input{%d,%d,%d}",
					gotGo.X, gotGo.Y, gotC.X, gotC.Y,
					tt.coord.I, tt.coord.J, tt.coord.K)
			}
		})
	}

	// Test that transformation is deterministic
	t.Run("deterministic", func(t *testing.T) {
		coord := CoordIJK{3, 1, 4}

		// Apply transformation twice
		var result1, result2 Vec2d
		_ijkToHex2d(&coord, &result1)
		_ijkToHex2d(&coord, &result2)

		if result1.X != result2.X || result1.Y != result2.Y {
			t.Errorf("_ijkToHex2d should be deterministic: first{%.15g,%.15g} != second{%.15g,%.15g}",
				result1.X, result1.Y, result2.X, result2.Y)
		}
	})
}

func absFloat(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
