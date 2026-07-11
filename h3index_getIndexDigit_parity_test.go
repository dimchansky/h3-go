//go:build cgo && c2go

package h3

import "testing"

func Test_getIndexDigit_parity(t *testing.T) {
	cells := []h3Index{
		0x8928308280fffff, // res 9 hexagon
		0x8f2830828052d25, // res 15 hexagon
		0x8009fffffffffff, // res 0 pentagon (base cell 4)
		0,                 // null
		^h3Index(0),       // all ones
	}
	for _, h := range cells {
		for res := int32(-1); res <= maxH3Res+1; res++ {
			var goOut int32
			goErr := getIndexDigit(h, res, &goOut)
			cOut, cErr := getIndexDigitC(h, res)
			if uint32(goErr) != uint32(cErr) {
				t.Fatalf("getIndexDigit(%#x, %d): err Go %v != C %v", uint64(h), res, goErr, cErr)
			}
			if goErr == eSuccess && goOut != cOut {
				t.Fatalf("getIndexDigit(%#x, %d): Go %d != C %d", uint64(h), res, goOut, cOut)
			}
		}
	}
}
