package h3

import "fmt"

// geoPrint prints a LatLng coordinate in degrees format.
// Uses geoToStringDegs internally to format the output.
// Ported from H3 C: utility.c::geoPrint.
//
//nolint:unused // ported from H3 C for parity completeness; exercised by cgo && c2go parity tests
func geoPrint(p *LatLng) {
	str := geoToStringDegs(p)
	fmt.Print(str)
}
