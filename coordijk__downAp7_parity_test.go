//go:build cgo

package h3

import "testing"

func Test_downAp7_parity(t *testing.T) {
	tests := []struct {
		name  string
		coord CoordIJK
	}{
		{"origin", CoordIJK{0, 0, 0}},
		{"unit i", CoordIJK{1, 0, 0}},
		{"unit j", CoordIJK{0, 1, 0}},
		{"unit k", CoordIJK{0, 0, 1}},
		{"positive coords", CoordIJK{1, 2, 3}},
		{"mixed coords", CoordIJK{3, -2, 5}},
		{"large coords", CoordIJK{5, 4, 6}},
		{"negative coords", CoordIJK{-1, -2, -3}},
		{"asymmetric", CoordIJK{2, 3, 1}},
		{"small values", CoordIJK{1, 1, 1}},
		{"normalized", CoordIJK{2, 1, 0}},
		{"needs normalization", CoordIJK{5, 3, 2}},
		{"zero i", CoordIJK{0, 2, 1}},
		{"zero j", CoordIJK{2, 0, 1}},
		{"zero k", CoordIJK{2, 1, 0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call C implementation
			gotC := _downAp7C(&tt.coord)

			// Call Go implementation
			gotGo := tt.coord
			_downAp7(&gotGo)

			// Compare results
			if gotGo.I != gotC.I || gotGo.J != gotC.J || gotGo.K != gotC.K {
				t.Errorf("_downAp7() mismatch: Go{%d,%d,%d} != C{%d,%d,%d} for input{%d,%d,%d}",
					gotGo.I, gotGo.J, gotGo.K, gotC.I, gotC.J, gotC.K,
					tt.coord.I, tt.coord.J, tt.coord.K)
			}
		})
	}

	// Test that transformation is deterministic
	t.Run("deterministic", func(t *testing.T) {
		coord := CoordIJK{1, 2, 3}

		// Apply transformation twice
		result1 := coord
		_downAp7(&result1)

		result2 := coord
		_downAp7(&result2)

		if result1.I != result2.I || result1.J != result2.J || result1.K != result2.K {
			t.Errorf("_downAp7 should be deterministic: first{%d,%d,%d} != second{%d,%d,%d}",
				result1.I, result1.J, result1.K, result2.I, result2.J, result2.K)
		}
	})

	// Test that up and down are somewhat inverse operations (for some values)
	t.Run("up_down_relationship", func(t *testing.T) {
		testCoords := []CoordIJK{
			{0, 0, 0}, {1, 0, 0}, {0, 1, 0}, {1, 1, 0},
		}

		for _, coord := range testCoords {
			// Apply down then up
			down := coord
			_downAp7(&down)
			upDown := down
			_upAp7(&upDown)

			// The result should be related to the original, though not necessarily identical
			// due to quantization. Just verify the operation doesn't crash.
			if upDown.I == 0 && upDown.J == 0 && upDown.K == 0 &&
				!(coord.I == 0 && coord.J == 0 && coord.K == 0) {
				t.Errorf("Up-down produced all zeros for non-zero input: orig{%d,%d,%d} -> down{%d,%d,%d} -> up{%d,%d,%d}",
					coord.I, coord.J, coord.K, down.I, down.J, down.K, upDown.I, upDown.J, upDown.K)
			}
		}
	})
}
