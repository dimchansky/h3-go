package h3

import "fmt"

// geoPrintln prints a LatLng coordinate in degrees format with a newline.
// Uses geoPrint internally and adds a newline.
// Ported from H3 C: utility.c::geoPrintln.
//
//nolint:unused // ported from H3 C for parity completeness; exercised by cgo && c2go parity tests
func geoPrintln(p *LatLng) {
	geoPrint(p)
	fmt.Print("\n")
}
