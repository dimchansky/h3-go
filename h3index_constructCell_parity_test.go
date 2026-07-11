//go:build cgo && c2go

package h3

import "testing"

func Test_constructCell_parity(t *testing.T) {
	type tc struct {
		res      int32
		baseCell int32
		digits   []int32
	}
	cases := []tc{
		{0, 0, nil},
		{0, 20, nil},
		{0, 121, nil},
		{0, 122, nil}, // bad base cell
		{0, -1, nil},  // bad base cell
		{-1, 0, nil},  // bad res
		{16, 0, nil},  // bad res
		{1, 0, []int32{0}},
		{1, 0, []int32{6}},
		{1, 0, []int32{7}},  // bad digit
		{1, 0, []int32{-1}}, // bad digit
		{1, 4, []int32{0}},  // pentagon, center: ok
		{1, 4, []int32{1}},  // pentagon, k-axis: deleted digit
		{2, 4, []int32{0, 1}},
		{2, 4, []int32{2, 1}}, // no longer pentagon after digit 2
		{3, 14, []int32{0, 0, 1}},
	}
	// Every digit sequence of the SF cell at various resolutions.
	for _, res := range []int32{1, 5, 9, 15} {
		var c h3Index
		if err := latLngToCell(&LatLng{Lat: Deg(37.775938728915946), Lng: Deg(-122.41795063018799)}, res, &c); err != eSuccess {
			t.Fatal(err)
		}
		digits := make([]int32, res)
		for r := int32(1); r <= res; r++ {
			digits[r-1] = h3GetIndexDigit(c, r)
		}
		cases = append(cases, tc{res, getBaseCellNumber(c), digits})
	}
	for _, c := range cases {
		var goOut h3Index
		goErr := constructCell(c.res, c.baseCell, c.digits, &goOut)
		cOut, cErr := constructCellC(c.res, c.baseCell, c.digits)
		if uint32(goErr) != uint32(cErr) {
			t.Fatalf("constructCell(%d, %d, %v): err Go %v != C %v", c.res, c.baseCell, c.digits, goErr, cErr)
		}
		if goErr == eSuccess && goOut != cOut {
			t.Fatalf("constructCell(%d, %d, %v): Go %#x != C %#x", c.res, c.baseCell, c.digits, uint64(goOut), uint64(cOut))
		}
	}
}
