//go:build cgo && c2go

package h3

import "testing"

func Test_ijkToHex2d_parity(t *testing.T) {
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
		{"normalized coords", coordIJK{5, 3, 1}},
		{"zero i", coordIJK{0, 5, 2}},
		{"zero j", coordIJK{3, 0, 1}},
		{"zero k", coordIJK{2, 4, 0}},
	}

	const tol = 1e-14

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call C implementation
			gotC := _ijkToHex2dC(&tt.coord)

			// Call Go implementation
			var gotGo vec2d
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
		coord := coordIJK{3, 1, 4}

		// Apply transformation twice
		var result1, result2 vec2d
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
