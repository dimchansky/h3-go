//go:build cgo

package h3

import "testing"

func Test_neighbor_parity(t *testing.T) {
	tests := []struct {
		name  string
		coord CoordIJK
		digit Direction
	}{
		// Test with origin and all valid directions
		{"origin center", CoordIJK{0, 0, 0}, CENTER_DIGIT},
		{"origin k", CoordIJK{0, 0, 0}, K_AXES_DIGIT},
		{"origin j", CoordIJK{0, 0, 0}, J_AXES_DIGIT},
		{"origin jk", CoordIJK{0, 0, 0}, JK_AXES_DIGIT},
		{"origin i", CoordIJK{0, 0, 0}, I_AXES_DIGIT},
		{"origin ik", CoordIJK{0, 0, 0}, IK_AXES_DIGIT},
		{"origin ij", CoordIJK{0, 0, 0}, IJ_AXES_DIGIT},

		// Test with various starting coordinates
		{"unit i + k", CoordIJK{1, 0, 0}, K_AXES_DIGIT},
		{"unit j + i", CoordIJK{0, 1, 0}, I_AXES_DIGIT},
		{"unit k + j", CoordIJK{0, 0, 1}, J_AXES_DIGIT},
		{"positive + direction", CoordIJK{2, 3, 1}, IJ_AXES_DIGIT},
		{"negative + direction", CoordIJK{-1, -2, -3}, K_AXES_DIGIT},
		{"mixed + direction", CoordIJK{2, -1, 3}, JK_AXES_DIGIT},

		// Test with large coordinates
		{"large + direction", CoordIJK{10, 20, 30}, I_AXES_DIGIT},

		// Test edge cases
		{"invalid direction below", CoordIJK{1, 2, 3}, -1},
		{"invalid direction above", CoordIJK{1, 2, 3}, NUM_DIGITS},
		{"invalid direction way above", CoordIJK{1, 2, 3}, 10},

		// Test boundary directions
		{"boundary direction 0", CoordIJK{1, 1, 1}, CENTER_DIGIT},
		{"boundary direction 7", CoordIJK{1, 1, 1}, INVALID_DIGIT},
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

	// Test that _neighbor with CENTER_DIGIT doesn't change coordinates for valid digits
	t.Run("center_digit_no_change", func(t *testing.T) {
		testCoords := []CoordIJK{
			{0, 0, 0}, {1, 2, 3}, {-1, -2, -3}, {5, 0, 2},
		}

		for _, coord := range testCoords {
			original := coord
			_neighbor(&coord, CENTER_DIGIT)

			// CENTER_DIGIT should not change coordinates (it adds {0,0,0})
			if coord.I != original.I || coord.J != original.J || coord.K != original.K {
				t.Errorf("_neighbor with CENTER_DIGIT should not change coords: original{%d,%d,%d} got{%d,%d,%d}",
					original.I, original.J, original.K, coord.I, coord.J, coord.K)
			}
		}
	})

	// Test applying unit vector directions to origin
	t.Run("unit_vectors_from_origin", func(t *testing.T) {
		for digit := K_AXES_DIGIT; digit < NUM_DIGITS; digit++ {
			coord := CoordIJK{0, 0, 0}
			_neighbor(&coord, digit)

			expected := UNIT_VECS[digit]

			if coord.I != expected.I || coord.J != expected.J || coord.K != expected.K {
				t.Errorf("_neighbor from origin with digit %d: got{%d,%d,%d} want{%d,%d,%d}",
					digit, coord.I, coord.J, coord.K, expected.I, expected.J, expected.K)
			}
		}
	})

	// Test invalid directions don't modify coordinates
	t.Run("invalid_directions_no_change", func(t *testing.T) {
		invalidDigits := []Direction{-1, NUM_DIGITS, NUM_DIGITS + 1, 10, -5}
		testCoord := CoordIJK{3, 1, 4}

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
