package coordijk

import "testing"

// TestCoordIJKApertureTransforms tests aperture transformation functions with C-derived expected results.
func TestCoordIJKApertureTransforms(t *testing.T) {
	tests := []struct {
		name             string
		input            CoordIJK
		expectedUpAp7    CoordIJK
		expectedUpAp7r   CoordIJK
		expectedDownAp7  CoordIJK
		expectedDownAp7r CoordIJK
		expectedDownAp3  CoordIJK
		expectedDownAp3r CoordIJK
	}{
		{"aperture (0,0,0)", CoordIJK{0, 0, 0}, CoordIJK{0, 0, 0}, CoordIJK{0, 0, 0}, CoordIJK{0, 0, 0}, CoordIJK{0, 0, 0}, CoordIJK{0, 0, 0}, CoordIJK{0, 0, 0}},
		{"aperture (7,0,0)", CoordIJK{7, 0, 0}, CoordIJK{3, 1, 0}, CoordIJK{3, 0, 1}, CoordIJK{21, 0, 7}, CoordIJK{21, 7, 0}, CoordIJK{14, 0, 7}, CoordIJK{14, 7, 0}},
		{"aperture (0,7,0)", CoordIJK{0, 7, 0}, CoordIJK{0, 3, 1}, CoordIJK{1, 3, 0}, CoordIJK{7, 21, 0}, CoordIJK{0, 21, 7}, CoordIJK{7, 14, 0}, CoordIJK{0, 14, 7}},
		{"aperture (1,0,0)", CoordIJK{1, 0, 0}, CoordIJK{0, 0, 0}, CoordIJK{0, 0, 0}, CoordIJK{3, 0, 1}, CoordIJK{3, 1, 0}, CoordIJK{2, 0, 1}, CoordIJK{2, 1, 0}},
		{"aperture (0,1,0)", CoordIJK{0, 1, 0}, CoordIJK{0, 0, 0}, CoordIJK{0, 0, 0}, CoordIJK{1, 3, 0}, CoordIJK{0, 3, 1}, CoordIJK{1, 2, 0}, CoordIJK{0, 2, 1}},
		{"aperture (0,0,1)", CoordIJK{0, 0, 1}, CoordIJK{0, 0, 0}, CoordIJK{0, 0, 0}, CoordIJK{0, 1, 3}, CoordIJK{1, 0, 3}, CoordIJK{0, 1, 2}, CoordIJK{1, 0, 2}},
		{"aperture (2,1,0)", CoordIJK{2, 1, 0}, CoordIJK{1, 1, 0}, CoordIJK{1, 0, 0}, CoordIJK{5, 1, 0}, CoordIJK{5, 4, 0}, CoordIJK{3, 0, 0}, CoordIJK{3, 3, 0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			upAp7 := tt.input
			upAp7.UpAp7()
			if upAp7 != tt.expectedUpAp7 {
				t.Errorf("UpAp7() = %v, want %v", upAp7, tt.expectedUpAp7)
			}

			upAp7r := tt.input
			upAp7r.UpAp7r()
			if upAp7r != tt.expectedUpAp7r {
				t.Errorf("UpAp7r() = %v, want %v", upAp7r, tt.expectedUpAp7r)
			}

			downAp7 := tt.input
			downAp7.DownAp7()
			if downAp7 != tt.expectedDownAp7 {
				t.Errorf("DownAp7() = %v, want %v", downAp7, tt.expectedDownAp7)
			}

			downAp7r := tt.input
			downAp7r.DownAp7r()
			if downAp7r != tt.expectedDownAp7r {
				t.Errorf("DownAp7r() = %v, want %v", downAp7r, tt.expectedDownAp7r)
			}

			downAp3 := tt.input
			downAp3.DownAp3()
			if downAp3 != tt.expectedDownAp3 {
				t.Errorf("DownAp3() = %v, want %v", downAp3, tt.expectedDownAp3)
			}

			downAp3r := tt.input
			downAp3r.DownAp3r()
			if downAp3r != tt.expectedDownAp3r {
				t.Errorf("DownAp3r() = %v, want %v", downAp3r, tt.expectedDownAp3r)
			}
		})
	}
}
