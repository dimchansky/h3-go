//go:build cgo

package h3

import "testing"

func Test_h3index_setH3Index_ParityWithC(t *testing.T) {
	cases := []struct {
		res  int32
		base int32
		init int32
	}{
		{0, 0, 0},
		{1, 1, 7},
		{2, 7, 3},
		{5, 57, 6},
		{10, 99, 2},
	}
	for _, c := range cases {
		var goH H3Index
		setH3Index(&goH, c.res, c.base, c.init)
		cH := setH3IndexC(c.res, c.base, c.init)
		if goH != cH {
			t.Fatalf("setH3Index mismatch res=%d base=%d init=%d: go=%x c=%x", c.res, c.base, c.init, uint64(goH), uint64(cH))
		}
	}
}
