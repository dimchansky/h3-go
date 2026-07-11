package h3

import "fmt"

// coordIjkPrint prints the coordIJK in [i, j, k] format.
// Ported from H3 C: utility.c::coordIjkPrint.
//
//nolint:unused // ported from H3 C for parity completeness; exercised by cgo && c2go parity tests
func coordIjkPrint(c *coordIJK) {
	fmt.Printf("[%d, %d, %d]", c.I, c.J, c.K)
}
