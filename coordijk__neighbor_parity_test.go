//go:build cgo && c2go

package h3

import "testing"

func Test_neighbor_parity(t *testing.T) {
	tests := []struct {
		name  string
		coord coordIJK
		digit direction
	}{
		// Test with origin and all valid directions
		{"origin center", coordIJK{0, 0, 0}, centerDigit},
		{"origin k", coordIJK{0, 0, 0}, kAxesDigit},
		{"origin j", coordIJK{0, 0, 0}, jAxesDigit},
		{"origin jk", coordIJK{0, 0, 0}, jkAxesDigit},
		{"origin i", coordIJK{0, 0, 0}, iAxesDigit},
		{"origin ik", coordIJK{0, 0, 0}, ikAxesDigit},
		{"origin ij", coordIJK{0, 0, 0}, ijAxesDigit},

		// Test with various starting coordinates
		{"unit i + k", coordIJK{1, 0, 0}, kAxesDigit},
		{"unit j + i", coordIJK{0, 1, 0}, iAxesDigit},
		{"unit k + j", coordIJK{0, 0, 1}, jAxesDigit},
		{"positive + direction", coordIJK{2, 3, 1}, ijAxesDigit},
		{"negative + direction", coordIJK{-1, -2, -3}, kAxesDigit},
		{"mixed + direction", coordIJK{2, -1, 3}, jkAxesDigit},

		// Test with large coordinates
		{"large + direction", coordIJK{10, 20, 30}, iAxesDigit},

		// Test edge cases
		{"invalid direction below", coordIJK{1, 2, 3}, -1},
		{"invalid direction above", coordIJK{1, 2, 3}, numDigits},
		{"invalid direction way above", coordIJK{1, 2, 3}, 10},

		// Test boundary directions
		{"boundary direction 0", coordIJK{1, 1, 1}, centerDigit},
		{"boundary direction 7", coordIJK{1, 1, 1}, invalidDigit},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call C implementation
			gotC := _neighborC(&tt.coord, tt.digit)

			// Call Go implementation
			gotGo := tt.coord
			_neighbor(&gotGo, tt.digit)

			// Compare results
			if gotGo.I != gotC.I || gotGo.J != gotC.J || gotGo.K != gotC.K {
				t.Errorf("_neighbor() mismatch: Go{%d,%d,%d} != C{%d,%d,%d} for input{%d,%d,%d} digit=%d",
					gotGo.I, gotGo.J, gotGo.K, gotC.I, gotC.J, gotC.K,
					tt.coord.I, tt.coord.J, tt.coord.K, tt.digit)
			}
		})
	}

	// Test that _neighbor with centerDigit doesn't change coordinates for valid digits
	t.Run("center_digit_no_change", func(t *testing.T) {
		testCoords := []coordIJK{
			{0, 0, 0}, {1, 2, 3}, {-1, -2, -3}, {5, 0, 2},
		}

		for _, coord := range testCoords {
			original := coord
			_neighbor(&coord, centerDigit)

			// centerDigit should not change coordinates (it adds {0,0,0})
			if coord.I != original.I || coord.J != original.J || coord.K != original.K {
				t.Errorf("_neighbor with centerDigit should not change coords: original{%d,%d,%d} got{%d,%d,%d}",
					original.I, original.J, original.K, coord.I, coord.J, coord.K)
			}
		}
	})

	// Test applying unit vector directions to origin
	t.Run("unit_vectors_from_origin", func(t *testing.T) {
		for digit := kAxesDigit; digit < numDigits; digit++ {
			coord := coordIJK{0, 0, 0}
			_neighbor(&coord, digit)

			expected := unitVecs[digit]

			if coord.I != expected.I || coord.J != expected.J || coord.K != expected.K {
				t.Errorf("_neighbor from origin with digit %d: got{%d,%d,%d} want{%d,%d,%d}",
					digit, coord.I, coord.J, coord.K, expected.I, expected.J, expected.K)
			}
		}
	})

	// Test invalid directions don't modify coordinates
	t.Run("invalid_directions_no_change", func(t *testing.T) {
		invalidDigits := []direction{-1, numDigits, numDigits + 1, 10, -5}
		testCoord := coordIJK{3, 1, 4}

		for _, digit := range invalidDigits {
			coord := testCoord
			_neighbor(&coord, digit)

			if coord.I != testCoord.I || coord.J != testCoord.J || coord.K != testCoord.K {
				t.Errorf("_neighbor with invalid digit %d should not change coords: original{%d,%d,%d} got{%d,%d,%d}",
					digit, testCoord.I, testCoord.J, testCoord.K, coord.I, coord.J, coord.K)
			}
		}
	})
}
