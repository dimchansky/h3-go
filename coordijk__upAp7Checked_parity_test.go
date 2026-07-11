//go:build cgo && c2go

package h3

import "testing"

func Test_upAp7Checked_parity(t *testing.T) {
	tests := []struct {
		name  string
		coord CoordIJK
	}{
		{"origin", CoordIJK{0, 0, 0}},
		{"unit i", CoordIJK{1, 0, 0}},
		{"unit j", CoordIJK{0, 1, 0}},
		{"unit k", CoordIJK{0, 0, 1}},
		{"positive coords", CoordIJK{7, 14, 21}}, // multiples of 7
		{"mixed coords", CoordIJK{3, -2, 5}},
		{"large coords", CoordIJK{35, 28, 42}}, // larger multiples of 7
		{"negative coords", CoordIJK{-7, -14, -21}},
		{"asymmetric", CoordIJK{8, 15, 22}},
		{"small values", CoordIJK{1, 2, 3}},
		{"aperture test", CoordIJK{3, 0, 1}}, // from down function
		{"aperture test 2", CoordIJK{1, 3, 0}},
		{"aperture test 3", CoordIJK{0, 1, 3}},
		{"normalized", CoordIJK{2, 1, 0}},
		{"needs normalization", CoordIJK{5, 3, 2}},
		// Overflow edge cases
		{"near max safe", CoordIJK{700000000, 0, 0}},
		{"near max safe j", CoordIJK{0, 700000000, 0}},
		{"balanced large", CoordIJK{300000000, 300000000, 0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call C implementation
			coordC := tt.coord
			coordCResult, errC := _upAp7CheckedC(&coordC)

			// Call Go implementation
			coordGo := tt.coord
			errGo := _upAp7Checked(&coordGo)

			// Compare errors first
			if errGo != errC {
				t.Errorf("_upAp7Checked() error mismatch: Go=%v != C=%v for input{%d,%d,%d}",
					errGo, errC, tt.coord.I, tt.coord.J, tt.coord.K)
				return
			}

			// If no error, compare coordinates
			if errGo == E_SUCCESS {
				if coordGo.I != coordCResult.I || coordGo.J != coordCResult.J || coordGo.K != coordCResult.K {
					t.Errorf("_upAp7Checked() mismatch: Go{%d,%d,%d} != C{%d,%d,%d} for input{%d,%d,%d}",
						coordGo.I, coordGo.J, coordGo.K, coordCResult.I, coordCResult.J, coordCResult.K,
						tt.coord.I, tt.coord.J, tt.coord.K)
				}
			}
		})
	}

	// Test that transformation is deterministic
	t.Run("deterministic", func(t *testing.T) {
		coord := CoordIJK{7, 14, 21}

		// Apply transformation twice
		result1 := coord
		err1 := _upAp7Checked(&result1)

		result2 := coord
		err2 := _upAp7Checked(&result2)

		if err1 != err2 {
			t.Errorf("_upAp7Checked should return same error: first=%v != second=%v", err1, err2)
		}

		if err1 == E_SUCCESS &&
			(result1.I != result2.I || result1.J != result2.J || result1.K != result2.K) {
			t.Errorf("_upAp7Checked should be deterministic: first{%d,%d,%d} != second{%d,%d,%d}",
				result1.I, result1.J, result1.K, result2.I, result2.J, result2.K)
		}
	})

	// Test overflow detection
	t.Run("overflow detection", func(t *testing.T) {
		overflowCases := []CoordIJK{
			{2000000000, 2000000000, 0}, // likely to overflow
			{-2000000000, -2000000000, 0},
			{2147483647, 0, 0}, // INT32_MAX
			{0, 2147483647, 0}, // INT32_MAX
		}

		for _, coord := range overflowCases {
			coordGo := coord
			errGo := _upAp7Checked(&coordGo)

			coordC := coord
			_, errC := _upAp7CheckedC(&coordC)

			if errGo != errC {
				t.Errorf("_upAp7Checked() overflow detection mismatch: Go=%v != C=%v for {%d,%d,%d}",
					errGo, errC, coord.I, coord.J, coord.K)
			}
		}
	})
}
