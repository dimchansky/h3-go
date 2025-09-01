package h3

import "fmt"

// h3Println prints the H3Index as a hexadecimal integer with a newline.
// Ported from H3 C: utility.c::h3Println
func h3Println(h H3Index) {
	fmt.Printf("%x\n", uint64(h))
}
