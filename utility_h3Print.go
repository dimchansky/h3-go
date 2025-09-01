package h3

import "fmt"

// h3Print prints the H3Index as a hexadecimal integer.
// Ported from H3 C: utility.c::h3Print
func h3Print(h H3Index) {
	fmt.Printf("%x", uint64(h))
}
