//go:build cgo

package h3

import "testing"

func Test_ijkMatches_parity(t *testing.T) {
	tests := []struct {
		name string
		c1   CoordIJK
		c2   CoordIJK
		want bool
	}{
		{"same zeros", CoordIJK{0, 0, 0}, CoordIJK{0, 0, 0}, true},
		{"same positive", CoordIJK{1, 2, 3}, CoordIJK{1, 2, 3}, true},
		{"same negative", CoordIJK{-1, -2, -3}, CoordIJK{-1, -2, -3}, true},
		{"different i", CoordIJK{1, 2, 3}, CoordIJK{2, 2, 3}, false},
		{"different j", CoordIJK{1, 2, 3}, CoordIJK{1, 3, 3}, false},
		{"different k", CoordIJK{1, 2, 3}, CoordIJK{1, 2, 4}, false},
		{"completely different", CoordIJK{1, 2, 3}, CoordIJK{4, 5, 6}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call C implementation
			gotC := _ijkMatchesC(&tt.c1, &tt.c2)

			// Call Go implementation
			gotGo := _ijkMatches(&tt.c1, &tt.c2)

			// Compare results
			if gotGo != gotC {
				t.Errorf("_ijkMatches() mismatch: Go=%v != C=%v for c1=%+v, c2=%+v",
					gotGo, gotC, tt.c1, tt.c2)
			}

			// Also verify against expected result
			if gotGo != tt.want {
				t.Errorf("_ijkMatches() Go result wrong: got=%v want=%v for c1=%+v, c2=%+v",
					gotGo, tt.want, tt.c1, tt.c2)
			}
		})
	}
}
