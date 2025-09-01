package h3

import "fmt"

// cellBoundaryPrint prints a CellBoundary in compact format with braces.
// Each vertex is printed as "lat lng " with spaces between vertices.
// Ported from H3 C: utility.c::cellBoundaryPrint
func cellBoundaryPrint(b *CellBoundary) {
	fmt.Print("{")
	for v := 0; v < int(b.NumVerts); v++ {
		str := geoToStringDegsNoFmt(&b.Verts[v])
		fmt.Printf("%s ", str)
	}
	fmt.Print("}")
}
