//go:build cgo

package h3

import "testing"

func Test_rotate60cw_parity(t *testing.T) {
	tests := []struct {
		name  string
		digit Direction
	}{
		// Valid direction rotations (clockwise)
		{"K -> JK", K_AXES_DIGIT},
		{"JK -> J", JK_AXES_DIGIT},
		{"J -> IJ", J_AXES_DIGIT},
		{"IJ -> I", IJ_AXES_DIGIT},
		{"I -> IK", I_AXES_DIGIT},
		{"IK -> K", IK_AXES_DIGIT},

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
			gotC := _rotate60cwC(tt.digit)

			// Call Go implementation
			gotGo := _rotate60cw(tt.digit)

			// Compare results
			if gotGo != gotC {
				t.Errorf("_rotate60cw() mismatch: Go=%d != C=%d for digit=%d",
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
				digit = _rotate60cw(digit)
			}

			if digit != original {
				t.Errorf("Six clockwise rotations should return to original: started=%d ended=%d",
					original, digit)
			}
		}
	})

	// Test rotation sequence correctness
	t.Run("rotation_sequence", func(t *testing.T) {
		// Starting from K_AXES_DIGIT, check the full clockwise sequence
		expected := []Direction{K_AXES_DIGIT, JK_AXES_DIGIT, J_AXES_DIGIT, IJ_AXES_DIGIT, I_AXES_DIGIT, IK_AXES_DIGIT}

		digit := K_AXES_DIGIT
		for i, expectedDigit := range expected {
			if digit != expectedDigit {
				t.Errorf("Rotation sequence step %d: expected=%d got=%d", i, expectedDigit, digit)
			}
			digit = _rotate60cw(digit)
		}

		// After 6 rotations, should be back to start
		if digit != K_AXES_DIGIT {
			t.Errorf("After 6 clockwise rotations: expected=%d got=%d", K_AXES_DIGIT, digit)
		}
	})
}
