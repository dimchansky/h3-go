//go:build cgo && c2go

package h3

import "testing"

func Test_ijkAdd_parity(t *testing.T) {
	tests := []struct {
		name string
		h1   CoordIJK
		h2   CoordIJK
	}{
		{"zeros", CoordIJK{0, 0, 0}, CoordIJK{0, 0, 0}},
		{"positive", CoordIJK{1, 2, 3}, CoordIJK{4, 5, 6}},
		{"negative", CoordIJK{-1, -2, -3}, CoordIJK{-4, -5, -6}},
		{"mixed", CoordIJK{10, -20, 30}, CoordIJK{-5, 15, -25}},
		{"large", CoordIJK{1000000, 2000000, -3000000}, CoordIJK{-500000, 1500000, 2500000}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call C implementation
			gotC := _ijkAddC(&tt.h1, &tt.h2)

			// Call Go implementation
			var gotGo CoordIJK
			_ijkAdd(&tt.h1, &tt.h2, &gotGo)

			// Compare results
			if gotGo.I != gotC.I || gotGo.J != gotC.J || gotGo.K != gotC.K {
				t.Errorf("_ijkAdd() mismatch: Go{%d,%d,%d} != C{%d,%d,%d}",
					gotGo.I, gotGo.J, gotGo.K, gotC.I, gotC.J, gotC.K)
			}
		})
	}
}
