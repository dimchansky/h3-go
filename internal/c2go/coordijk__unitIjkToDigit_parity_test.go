//go:build cgo

package c2go

import "testing"

func Test_unitIjkToDigit_parity(t *testing.T) {
	tests := []struct {
		name  string
		coord CoordIJK
	}{
		{"center", CoordIJK{0, 0, 0}},
		{"k unit", CoordIJK{0, 0, 1}},
		{"j unit", CoordIJK{0, 1, 0}},
		{"jk unit", CoordIJK{0, 1, 1}},
		{"i unit", CoordIJK{1, 0, 0}},
		{"ik unit", CoordIJK{1, 0, 1}},
		{"ij unit", CoordIJK{1, 1, 0}},

		// Test various scaled coordinates
		{"scaled k", CoordIJK{0, 0, 3}},
		{"scaled j", CoordIJK{0, 2, 0}},
		{"scaled jk", CoordIJK{0, 2, 2}},
		{"scaled i", CoordIJK{5, 0, 0}},
		{"scaled ik", CoordIJK{3, 0, 3}},
		{"scaled ij", CoordIJK{4, 4, 0}},

		// Test coordinates that need normalization
		{"needs norm k", CoordIJK{1, 1, 2}},
		{"needs norm j", CoordIJK{1, 2, 1}},
		{"needs norm i", CoordIJK{2, 1, 1}},

		// Test various other coordinates
		{"random 1", CoordIJK{1, 2, 3}},
		{"random 2", CoordIJK{2, 1, 3}},
		{"random 3", CoordIJK{1, 0, 2}},
		{"negative", CoordIJK{-1, -2, -3}},
		{"mixed", CoordIJK{2, -1, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call C implementation
			gotC := Direction(_unitIjkToDigitC(&tt.coord))

			// Call Go implementation
			gotGo := _unitIjkToDigit(&tt.coord)

			// Compare results
			if gotGo != gotC {
				t.Errorf("_unitIjkToDigit() mismatch: Go=%d != C=%d for input{%d,%d,%d}",
					gotGo, gotC, tt.coord.I, tt.coord.J, tt.coord.K)
			}
		})
	}

	// Test all unit vectors directly
	t.Run("all_unit_vectors", func(t *testing.T) {
		for digit := CENTER_DIGIT; digit < NUM_DIGITS; digit++ {
			coord := UNIT_VECS[digit]

			gotC := Direction(_unitIjkToDigitC(&coord))
			gotGo := _unitIjkToDigit(&coord)

			if gotGo != gotC {
				t.Errorf("Unit vector %d: Go=%d != C=%d for coord{%d,%d,%d}",
					digit, gotGo, gotC, coord.I, coord.J, coord.K)
			}

			// Verify that unit vectors return their own digit
			if gotGo != digit {
				t.Errorf("Unit vector %d should return %d, got %d for coord{%d,%d,%d}",
					digit, digit, gotGo, coord.I, coord.J, coord.K)
			}
		}
	})
}
