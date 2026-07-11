//go:build cgo && c2go

package h3

import "testing"

func Test_rotate60ccw_parity(t *testing.T) {
	tests := []struct {
		name  string
		digit Direction
	}{
		// Valid direction rotations (counter-clockwise)
		{"K -> IK", K_AXES_DIGIT},
		{"IK -> I", IK_AXES_DIGIT},
		{"I -> IJ", I_AXES_DIGIT},
		{"IJ -> J", IJ_AXES_DIGIT},
		{"J -> JK", J_AXES_DIGIT},
		{"JK -> K", JK_AXES_DIGIT},

		// Center digit and invalid digit should return unchanged
		{"center unchanged", CENTER_DIGIT},
		{"invalid unchanged", INVALID_DIGIT},

		// Test edge cases and invalid values (but avoid negative due to type conversion)
		{"large value unchanged", 100},
		{"zero unchanged", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call C implementation
			gotC := _rotate60ccwC(tt.digit)

			// Call Go implementation
			gotGo := _rotate60ccw(tt.digit)

			// Compare results
			if gotGo != gotC {
				t.Errorf("_rotate60ccw() mismatch: Go=%d != C=%d for digit=%d",
					gotGo, gotC, tt.digit)
			}
		})
	}

	// Test that 6 rotations return to original
	t.Run("six_rotations", func(t *testing.T) {
		validDigits := []Direction{K_AXES_DIGIT, IK_AXES_DIGIT, I_AXES_DIGIT, IJ_AXES_DIGIT, J_AXES_DIGIT, JK_AXES_DIGIT}

		for _, original := range validDigits {
			digit := original
			for i := 0; i < 6; i++ {
				digit = _rotate60ccw(digit)
			}

			if digit != original {
				t.Errorf("Six counter-clockwise rotations should return to original: started=%d ended=%d",
					original, digit)
			}
		}
	})

	// Test that it's the inverse of clockwise rotation
	t.Run("inverse_of_clockwise", func(t *testing.T) {
		validDigits := []Direction{K_AXES_DIGIT, IK_AXES_DIGIT, I_AXES_DIGIT, IJ_AXES_DIGIT, J_AXES_DIGIT, JK_AXES_DIGIT}

		for _, digit := range validDigits {
			// Apply counter-clockwise then clockwise
			ccw_then_cw := _rotate60cw(_rotate60ccw(digit))

			// Apply clockwise then counter-clockwise
			cw_then_ccw := _rotate60ccw(_rotate60cw(digit))

			if ccw_then_cw != digit {
				t.Errorf("CCW then CW should return original: digit=%d got=%d", digit, ccw_then_cw)
			}

			if cw_then_ccw != digit {
				t.Errorf("CW then CCW should return original: digit=%d got=%d", digit, cw_then_ccw)
			}
		}
	})
}
