//go:build c2go

package c2go

import "testing"

func Test_setIJK_parity(t *testing.T) {
	tests := []struct {
		name    string
		i, j, k int
	}{
		{"zeros", 0, 0, 0},
		{"positive", 1, 2, 3},
		{"negative", -1, -2, -3},
		{"mixed", 10, -20, 30},
		{"large", 1000000, -2000000, 3000000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call C implementation
			gotC := _setIJKC(tt.i, tt.j, tt.k)

			// Call Go implementation
			var gotGo CoordIJK
			_setIJK(&gotGo, tt.i, tt.j, tt.k)

			// Compare results
			if gotGo.I != gotC.I || gotGo.J != gotC.J || gotGo.K != gotC.K {
				t.Errorf("_setIJK() mismatch: Go{%d,%d,%d} != C{%d,%d,%d}",
					gotGo.I, gotGo.J, gotGo.K, gotC.I, gotC.J, gotC.K)
			}
		})
	}
}
