//go:build cgo && c2go

package h3

import "testing"

func Test_unitIjkToDigit_parity(t *testing.T) {
	tests := []struct {
		name  string
		coord coordIJK
	}{
		{"center", coordIJK{0, 0, 0}},
		{"k unit", coordIJK{0, 0, 1}},
		{"j unit", coordIJK{0, 1, 0}},
		{"jk unit", coordIJK{0, 1, 1}},
		{"i unit", coordIJK{1, 0, 0}},
		{"ik unit", coordIJK{1, 0, 1}},
		{"ij unit", coordIJK{1, 1, 0}},

		// Test various scaled coordinates
		{"scaled k", coordIJK{0, 0, 3}},
		{"scaled j", coordIJK{0, 2, 0}},
		{"scaled jk", coordIJK{0, 2, 2}},
		{"scaled i", coordIJK{5, 0, 0}},
		{"scaled ik", coordIJK{3, 0, 3}},
		{"scaled ij", coordIJK{4, 4, 0}},

		// Test coordinates that need normalization
		{"needs norm k", coordIJK{1, 1, 2}},
		{"needs norm j", coordIJK{1, 2, 1}},
		{"needs norm i", coordIJK{2, 1, 1}},

		// Test various other coordinates
		{"random 1", coordIJK{1, 2, 3}},
		{"random 2", coordIJK{2, 1, 3}},
		{"random 3", coordIJK{1, 0, 2}},
		{"negative", coordIJK{-1, -2, -3}},
		{"mixed", coordIJK{2, -1, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call C implementation
			gotC := direction(_unitIjkToDigitC(&tt.coord))

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
		for digit := centerDigit; digit < numDigits; digit++ {
			coord := unitVecs[digit]

			gotC := direction(_unitIjkToDigitC(&coord))
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
