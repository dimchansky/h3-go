package h3

import "fmt"

// geoPrintlnNoFmt prints a LatLng coordinate in degrees format without parentheses with a newline.
// Uses geoPrintNoFmt internally and adds a newline.
// Ported from H3 C: utility.c::geoPrintlnNoFmt.
//
//nolint:unused // ported from H3 C for parity completeness; exercised by cgo && c2go parity tests
func geoPrintlnNoFmt(p *LatLng) {
	geoPrintNoFmt(p)
	fmt.Print("\n")
}
