//go:build cgo && c2go

package h3

import "testing"

func Test_ijkDistance_parity(t *testing.T) {
	tests := []struct {
		name string
		c1   coordIJK
		c2   coordIJK
	}{
		{"same point", coordIJK{0, 0, 0}, coordIJK{0, 0, 0}},
		{"unit distance i", coordIJK{0, 0, 0}, coordIJK{1, 0, 0}},
		{"unit distance j", coordIJK{0, 0, 0}, coordIJK{0, 1, 0}},
		{"unit distance k", coordIJK{0, 0, 0}, coordIJK{0, 0, 1}},
		{"simple distance", coordIJK{1, 2, 3}, coordIJK{4, 5, 6}},
		{"negative coords", coordIJK{-1, -2, -3}, coordIJK{1, 2, 3}},
		{"mixed coords", coordIJK{10, -5, 2}, coordIJK{-3, 8, -1}},
		{"large distance", coordIJK{100, 200, 300}, coordIJK{400, 500, 600}},
		{"asymmetric", coordIJK{0, 0, 0}, coordIJK{5, 3, 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call C implementation
			gotC := ijkDistanceC(&tt.c1, &tt.c2)

			// Call Go implementation
			gotGo := ijkDistance(&tt.c1, &tt.c2)

			// Compare results
			if gotGo != gotC {
				t.Errorf("ijkDistance() mismatch: Go=%d != C=%d for c1=%+v, c2=%+v",
					gotGo, gotC, tt.c1, tt.c2)
			}

			// Verify symmetry (distance should be the same both ways)
			gotGoReverse := ijkDistance(&tt.c2, &tt.c1)
			if gotGo != gotGoReverse {
				t.Errorf("ijkDistance() not symmetric: %d != %d for c1=%+v, c2=%+v",
					gotGo, gotGoReverse, tt.c1, tt.c2)
			}
		})
	}
}
