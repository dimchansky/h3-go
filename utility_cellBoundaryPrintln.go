package h3

import "fmt"

// cellBoundaryPrintln prints a CellBoundary in multi-line format with braces and newlines.
// Each vertex is printed on a separate line with indentation.
// Ported from H3 C: utility.c::cellBoundaryPrintln.
//
//nolint:unused // ported from H3 C for parity completeness; exercised by cgo && c2go parity tests
func cellBoundaryPrintln(b *CellBoundary) {
	fmt.Print("{\n")
	for v := 0; v < int(b.numVerts); v++ {
		str := geoToStringDegsNoFmt(&b.verts[v])
		fmt.Printf("   %s\n", str)
	}
	fmt.Print("}\n")
}
