//go:build cgo

package c2go

import "testing"

func Test_ijkNormalize_parity(t *testing.T) {
	tests := []struct {
		name  string
		coord CoordIJK
	}{
		{"zeros", CoordIJK{0, 0, 0}},
		{"all positive", CoordIJK{1, 2, 3}},
		{"negative i", CoordIJK{-1, 2, 3}},
		{"negative j", CoordIJK{1, -2, 3}},
		{"negative k", CoordIJK{1, 2, -3}},
		{"multiple negatives", CoordIJK{-1, -2, 3}},
		{"all negative", CoordIJK{-1, -2, -3}},
		{"needs min reduction", CoordIJK{5, 3, 7}},
		{"mixed normalize case", CoordIJK{-2, 1, -1}},
		{"large values", CoordIJK{-100, 50, -25}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call C implementation
			gotC := _ijkNormalizeC(&tt.coord)

			// Call Go implementation
			gotGo := tt.coord
			_ijkNormalize(&gotGo)

			// Compare results
			if gotGo.I != gotC.I || gotGo.J != gotC.J || gotGo.K != gotC.K {
				t.Errorf("_ijkNormalize() mismatch: Go{%d,%d,%d} != C{%d,%d,%d} for input{%d,%d,%d}",
					gotGo.I, gotGo.J, gotGo.K, gotC.I, gotC.J, gotC.K,
					tt.coord.I, tt.coord.J, tt.coord.K)
			}
		})
	}
}
