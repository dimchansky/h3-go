package h3

import "fmt"

// h3Print prints the h3Index as a hexadecimal integer.
// Ported from H3 C: utility.c::h3Print.
//
//nolint:unused // ported from H3 C for parity completeness; exercised by cgo && c2go parity tests
func h3Print(h h3Index) {
	fmt.Printf("%x", uint64(h))
}
