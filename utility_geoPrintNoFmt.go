package h3

import "fmt"

// geoPrintNoFmt prints a LatLng coordinate in degrees format without parentheses.
// Uses geoToStringDegsNoFmt internally to format the output.
// Ported from H3 C: utility.c::geoPrintNoFmt.
//
//nolint:unused // ported from H3 C for parity completeness; exercised by cgo && c2go parity tests
func geoPrintNoFmt(p *LatLng) {
	str := geoToStringDegsNoFmt(p)
	fmt.Print(str)
}
