//go:build c2go

package c2go

import "testing"

func Test_ipow_ParityWithC(t *testing.T) {
	cases := []struct {
		base int64
		exp  int64
		name string
	}{
		{base: 0, exp: 0, name: "0^0"},
		{base: 0, exp: 5, name: "0^5"},
		{base: 1, exp: 0, name: "1^0"},
		{base: 1, exp: 5, name: "1^5"},
		{base: 2, exp: 10, name: "2^10"},
		{base: -2, exp: 3, name: "(-2)^3"},
		{base: -3, exp: 0, name: "(-3)^0"},
		{base: 5, exp: 1, name: "5^1"},
		{base: -5, exp: 2, name: "(-5)^2"},
	}

	for _, tc := range cases {
		gotGo := _ipow(tc.base, tc.exp)
		gotC := _ipowC(tc.base, tc.exp)
		if gotGo != gotC {
			t.Fatalf("%s: mismatch go=%d c=%d (base=%d exp=%d)", tc.name, gotGo, gotC, tc.base, tc.exp)
		}
	}
}
