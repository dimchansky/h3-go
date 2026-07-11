package h3

import "fmt"

// bboxPrint prints a BBox in compact format showing north, south, east, west bounds in degrees.
// Ported from H3 C: utility.c::bboxPrint.
//
//nolint:unused // ported from H3 C for parity completeness; exercised by cgo && c2go parity tests
func bboxPrint(bbox *BBox) {
	fmt.Printf("bbox {%.9f, %.9f, %.9f, %.9f}",
		bbox.North.Deg(), bbox.South.Deg(), bbox.East.Deg(), bbox.West.Deg())
}
