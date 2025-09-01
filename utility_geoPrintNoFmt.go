package h3

import "fmt"

// geoPrintNoFmt prints a LatLng coordinate in degrees format without parentheses.
// Uses geoToStringDegsNoFmt internally to format the output.
// Ported from H3 C: utility.c::geoPrintNoFmt
func geoPrintNoFmt(p *LatLng) {
	str := geoToStringDegsNoFmt(p)
	fmt.Print(str)
}
