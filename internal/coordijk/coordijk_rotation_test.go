package coordijk

import "testing"

// TestCoordIJKRotation tests Rotate60CCW and Rotate60CW functions.
func TestCoordIJKRotation(t *testing.T) {
	tests := []struct {
		name        string
		input       CoordIJK
		expectedCCW CoordIJK
		expectedCW  CoordIJK
	}{
		{"rotate (0,0,0)", CoordIJK{0, 0, 0}, CoordIJK{0, 0, 0}, CoordIJK{0, 0, 0}},
		{"rotate (1,0,0)", CoordIJK{1, 0, 0}, CoordIJK{1, 1, 0}, CoordIJK{1, 0, 1}},
		{"rotate (0,1,0)", CoordIJK{0, 1, 0}, CoordIJK{0, 1, 1}, CoordIJK{1, 1, 0}},
		{"rotate (0,0,1)", CoordIJK{0, 0, 1}, CoordIJK{1, 0, 1}, CoordIJK{0, 1, 1}},
		{"rotate (1,1,0)", CoordIJK{1, 1, 0}, CoordIJK{0, 1, 0}, CoordIJK{1, 0, 0}},
		{"rotate (1,0,1)", CoordIJK{1, 0, 1}, CoordIJK{1, 0, 0}, CoordIJK{0, 0, 1}},
		{"rotate (0,1,1)", CoordIJK{0, 1, 1}, CoordIJK{0, 0, 1}, CoordIJK{0, 1, 0}},
		{"rotate (2,-1,1)", CoordIJK{2, -1, 1}, CoordIJK{3, 1, 0}, CoordIJK{1, 0, 3}},
		{"rotate (-1,2,-1)", CoordIJK{-1, 2, -1}, CoordIJK{0, 3, 3}, CoordIJK{3, 3, 0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultCCW := tt.input
			resultCCW.Rotate60CCW()
			if resultCCW != tt.expectedCCW {
				t.Errorf("Rotate60CCW() = %v, want %v", resultCCW, tt.expectedCCW)
			}

			resultCW := tt.input
			resultCW.Rotate60CW()
			if resultCW != tt.expectedCW {
				t.Errorf("Rotate60CW() = %v, want %v", resultCW, tt.expectedCW)
			}
		})
	}
}
