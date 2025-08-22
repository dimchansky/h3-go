//go:build c2go

package c2go

import "testing"

func Test_ijkScale_parity(t *testing.T) {
	tests := []struct {
		name   string
		coord  CoordIJK
		factor int
	}{
		{"zeros", CoordIJK{0, 0, 0}, 5},
		{"positive scale positive", CoordIJK{1, 2, 3}, 2},
		{"positive scale negative", CoordIJK{1, 2, 3}, -2},
		{"negative scale positive", CoordIJK{-1, -2, -3}, 2},
		{"negative scale negative", CoordIJK{-1, -2, -3}, -2},
		{"scale by zero", CoordIJK{5, 10, 15}, 0},
		{"scale by one", CoordIJK{5, 10, 15}, 1},
		{"mixed coords", CoordIJK{10, -20, 30}, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call C implementation
			gotC := _ijkScaleC(&tt.coord, tt.factor)

			// Call Go implementation
			gotGo := tt.coord
			_ijkScale(&gotGo, tt.factor)

			// Compare results
			if gotGo.I != gotC.I || gotGo.J != gotC.J || gotGo.K != gotC.K {
				t.Errorf("_ijkScale() mismatch: Go{%d,%d,%d} != C{%d,%d,%d}",
					gotGo.I, gotGo.J, gotGo.K, gotC.I, gotC.J, gotC.K)
			}
		})
	}
}
