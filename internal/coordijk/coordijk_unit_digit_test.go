package coordijk

import "testing"

// TestUnitIJKToDigitCValidation validates Go implementation against H3 C test cases
func TestUnitIJKToDigitCValidation(t *testing.T) {
    tests := []struct {
        name     string
        input    CoordIJK
        expected Direction
    }{
        {"CENTER_DIGIT", CoordIJK{0, 0, 0}, 0},
        {"K_AXES_DIGIT", CoordIJK{0, 0, 1}, 1},
        {"J_AXES_DIGIT", CoordIJK{0, 1, 0}, 2},
        {"JK_AXES_DIGIT", CoordIJK{0, 1, 1}, 3},
        {"I_AXES_DIGIT", CoordIJK{1, 0, 0}, 4},
        {"IK_AXES_DIGIT", CoordIJK{1, 0, 1}, 5},
        {"IJ_AXES_DIGIT", CoordIJK{1, 1, 0}, 6},
        {"INVALID_DIGIT (out of range)", CoordIJK{2, 0, 0}, 7},
        {"CENTER_DIGIT (unnormalized zero)", CoordIJK{2, 2, 2}, 0},
        {"CENTER_DIGIT (unnormalized zero)", CoordIJK{3, 3, 3}, 0},
        {"INVALID_DIGIT (no match)", CoordIJK{-1, 1, 0}, 7},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := UnitIJKToDigit(tt.input)
            if result != tt.expected {
                t.Errorf("UnitIJKToDigit(%v) = %v, want %v", tt.input, result, tt.expected)
            }
        })
    }
}

