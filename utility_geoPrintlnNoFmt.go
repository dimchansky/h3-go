package h3

import "fmt"

// geoPrintlnNoFmt prints a LatLng coordinate in degrees format without parentheses with a newline.
// Uses geoPrintNoFmt internally and adds a newline.
// Ported from H3 C: utility.c::geoPrintlnNoFmt
func geoPrintlnNoFmt(p *LatLng) {
	geoPrintNoFmt(p)
	fmt.Print("\n")
}
