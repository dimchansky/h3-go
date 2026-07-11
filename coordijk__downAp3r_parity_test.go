//go:build cgo && c2go

package h3

import "testing"

func Test_downAp3r_parity(t *testing.T) {
	tests := []struct {
		name  string
		coord coordIJK
	}{
		{"origin", coordIJK{0, 0, 0}},
		{"unit i", coordIJK{1, 0, 0}},
		{"unit j", coordIJK{0, 1, 0}},
		{"unit k", coordIJK{0, 0, 1}},
		{"positive coords", coordIJK{1, 2, 3}},
		{"mixed coords", coordIJK{3, -2, 5}},
		{"large coords", coordIJK{5, 4, 6}},
		{"negative coords", coordIJK{-1, -2, -3}},
		{"asymmetric", coordIJK{2, 3, 1}},
		{"small values", coordIJK{1, 1, 1}},
		{"normalized", coordIJK{2, 1, 0}},
		{"needs normalization", coordIJK{5, 3, 2}},
		{"zero i", coordIJK{0, 2, 1}},
		{"zero j", coordIJK{2, 0, 1}},
		{"zero k", coordIJK{2, 1, 0}},
		{"aperture test", coordIJK{3, 1, 0}}, // from down function
		{"aperture test 2", coordIJK{0, 3, 1}},
		{"aperture test 3", coordIJK{1, 0, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call C implementation
			gotC := _downAp3rC(&tt.coord)

			// Call Go implementation
			gotGo := tt.coord
			_downAp3r(&gotGo)

			// Compare results
			if gotGo.I != gotC.I || gotGo.J != gotC.J || gotGo.K != gotC.K {
				t.Errorf("_downAp3r() mismatch: Go{%d,%d,%d} != C{%d,%d,%d} for input{%d,%d,%d}",
					gotGo.I, gotGo.J, gotGo.K, gotC.I, gotC.J, gotC.K,
					tt.coord.I, tt.coord.J, tt.coord.K)
			}
		})
	}

	// Test that transformation is deterministic
	t.Run("deterministic", func(t *testing.T) {
		coord := coordIJK{1, 2, 3}

		// Apply transformation twice
		result1 := coord
		_downAp3r(&result1)

		result2 := coord
		_downAp3r(&result2)

		if result1.I != result2.I || result1.J != result2.J || result1.K != result2.K {
			t.Errorf("_downAp3r should be deterministic: first{%d,%d,%d} != second{%d,%d,%d}",
				result1.I, result1.J, result1.K, result2.I, result2.J, result2.K)
		}
	})

	// Test relationship between clockwise and counter-clockwise versions
	t.Run("ccw_vs_cw_relationship", func(t *testing.T) {
		testCoords := []coordIJK{
			{1, 0, 0}, {0, 1, 0}, {0, 0, 1}, {1, 1, 0},
			{2, 1, 0}, {1, 2, 0}, {3, 2, 1},
		}

		for _, coord := range testCoords {
			// Apply both versions
			ccw := coord
			_downAp3(&ccw)

			cw := coord
			_downAp3r(&cw)

			// They should generally be different (except for origin)
			if coord.I == 0 && coord.J == 0 && coord.K == 0 {
				// Origin should map to itself in both
				if ccw.I != 0 || ccw.J != 0 || ccw.K != 0 ||
					cw.I != 0 || cw.J != 0 || cw.K != 0 {
					t.Errorf("Origin should map to origin: ccw{%d,%d,%d} cw{%d,%d,%d}",
						ccw.I, ccw.J, ccw.K, cw.I, cw.J, cw.K)
				}
			} else {
				// Non-origin values typically produce different results
				// Just verify neither crashes
				if ccw.I == 0 && ccw.J == 0 && ccw.K == 0 {
					t.Errorf("CCW produced zeros for non-zero input: orig{%d,%d,%d} -> ccw{%d,%d,%d}",
						coord.I, coord.J, coord.K, ccw.I, ccw.J, ccw.K)
				}
				if cw.I == 0 && cw.J == 0 && cw.K == 0 {
					t.Errorf("CW produced zeros for non-zero input: orig{%d,%d,%d} -> cw{%d,%d,%d}",
						coord.I, coord.J, coord.K, cw.I, cw.J, cw.K)
				}
			}
		}
	})
}
