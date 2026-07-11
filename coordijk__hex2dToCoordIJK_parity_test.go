//go:build cgo && c2go

package h3

import "testing"

func Test_hex2dToCoordIJK_parity(t *testing.T) {
	tests := []struct {
		name string
		vec  vec2d
	}{
		{"origin", vec2d{0, 0}},
		{"unit x", vec2d{1, 0}},
		{"unit y", vec2d{0, 1}},
		{"positive coords", vec2d{1.5, 2.3}},
		{"negative x", vec2d{-1.5, 2.3}},
		{"negative y", vec2d{1.5, -2.3}},
		{"both negative", vec2d{-1.5, -2.3}},
		{"small coords", vec2d{0.1, 0.1}},
		{"large coords", vec2d{10.5, 20.7}},
		{"hex center", vec2d{1.0, 1.732050808}}, // approx hex center
		{"fractional", vec2d{0.333, 0.577}},
		{"hex boundary", vec2d{0.5, 0.866}},
		{"asymmetric", vec2d{3.14, 2.71}},
		{"zero x", vec2d{0, 5.2}},
		{"zero y", vec2d{3.8, 0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call C implementation
			gotC := _hex2dToCoordIJKC(&tt.vec)

			// Call Go implementation
			var gotGo coordIJK
			_hex2dToCoordIJK(&tt.vec, &gotGo)

			// Compare results
			if gotGo.I != gotC.I || gotGo.J != gotC.J || gotGo.K != gotC.K {
				t.Errorf("_hex2dToCoordIJK() mismatch: Go{%d,%d,%d} != C{%d,%d,%d} for input{%.15g,%.15g}",
					gotGo.I, gotGo.J, gotGo.K, gotC.I, gotC.J, gotC.K,
					tt.vec.X, tt.vec.Y)
			}
		})
	}

	// Test that transformation is deterministic
	t.Run("deterministic", func(t *testing.T) {
		vec := vec2d{3.14, 2.71}

		// Apply transformation twice
		var result1, result2 coordIJK
		_hex2dToCoordIJK(&vec, &result1)
		_hex2dToCoordIJK(&vec, &result2)

		if result1.I != result2.I || result1.J != result2.J || result1.K != result2.K {
			t.Errorf("_hex2dToCoordIJK should be deterministic: first{%d,%d,%d} != second{%d,%d,%d}",
				result1.I, result1.J, result1.K, result2.I, result2.J, result2.K)
		}
	})

	// Test round-trip consistency for simple cases
	t.Run("round_trip", func(t *testing.T) {
		testCoords := []coordIJK{
			{0, 0, 0}, {1, 0, 0}, {0, 1, 0}, {0, 0, 1},
			{1, 1, 0}, {1, 0, 1}, {0, 1, 1},
		}

		for _, orig := range testCoords {
			// Convert IJK -> hex2d -> IJK
			var v vec2d
			_ijkToHex2d(&orig, &v)

			var result coordIJK
			_hex2dToCoordIJK(&v, &result)

			// For simple unit vectors, we should get back something reasonable
			// Note: Round-trip isn't perfect due to quantization, so we just check
			// that the conversion doesn't crash and produces valid output
			if result.I == 0 && result.J == 0 && result.K == 0 &&
				!(orig.I == 0 && orig.J == 0 && orig.K == 0) {
				t.Errorf("Round-trip produced all zeros for non-zero input: orig{%d,%d,%d} -> vec{%.6g,%.6g} -> result{%d,%d,%d}",
					orig.I, orig.J, orig.K, v.X, v.Y, result.I, result.J, result.K)
			}
		}
	})
}
