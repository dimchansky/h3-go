package h3

import "fmt"

// coordIjkPrint prints the CoordIJK in [i, j, k] format.
// Ported from H3 C: utility.c::coordIjkPrint.
//
//nolint:unused // ported from H3 C for parity completeness; exercised by cgo && c2go parity tests
func coordIjkPrint(c *CoordIJK) {
	fmt.Printf("[%d, %d, %d]", c.I, c.J, c.K)
}
