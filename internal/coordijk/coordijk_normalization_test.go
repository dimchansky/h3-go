package coordijk

import "testing"

// TestCoordIJKNormalizeCValidation validates Go implementation against H3 C test cases.
func TestCoordIJKNormalizeCValidation(t *testing.T) {
	tests := []struct {
		name     string
		input    CoordIJK
		expected CoordIJK
	}{
		{"normalize (0,0,0)", CoordIJK{0, 0, 0}, CoordIJK{0, 0, 0}},
		{"normalize (1,0,0)", CoordIJK{1, 0, 0}, CoordIJK{1, 0, 0}},
		{"normalize (0,1,0)", CoordIJK{0, 1, 0}, CoordIJK{0, 1, 0}},
		{"normalize (0,0,1)", CoordIJK{0, 0, 1}, CoordIJK{0, 0, 1}},
		{"normalize (1,1,0)", CoordIJK{1, 1, 0}, CoordIJK{1, 1, 0}},
		{"normalize (1,0,1)", CoordIJK{1, 0, 1}, CoordIJK{1, 0, 1}},
		{"normalize (0,1,1)", CoordIJK{0, 1, 1}, CoordIJK{0, 1, 1}},
		{"normalize (1,1,1)", CoordIJK{1, 1, 1}, CoordIJK{0, 0, 0}},
		{"normalize (2,2,2)", CoordIJK{2, 2, 2}, CoordIJK{0, 0, 0}},
		{"normalize (3,1,2)", CoordIJK{3, 1, 2}, CoordIJK{2, 0, 1}},
		{"normalize (-1,0,0)", CoordIJK{-1, 0, 0}, CoordIJK{0, 1, 1}},
		{"normalize (0,-1,0)", CoordIJK{0, -1, 0}, CoordIJK{1, 0, 1}},
		{"normalize (0,0,-1)", CoordIJK{0, 0, -1}, CoordIJK{1, 1, 0}},
		{"normalize (-1,-1,-1)", CoordIJK{-1, -1, -1}, CoordIJK{0, 0, 0}},
		{"normalize (5,-2,3)", CoordIJK{5, -2, 3}, CoordIJK{7, 0, 5}},
		{"normalize (-3,4,1)", CoordIJK{-3, 4, 1}, CoordIJK{0, 7, 4}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input
			result.Normalize()
			if result != tt.expected {
				t.Errorf("CoordIJK.Normalize() = %v, want %v", result, tt.expected)
			}
		})
	}
}
