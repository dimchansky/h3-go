package h3

import "strconv"

// stringToH3 parses a hex string into an H3Index.
// Ported behavior: returns 0 index and non-zero error on failure.
// Ported from H3 C: h3Index.c::stringToH3.
func stringToH3(s string) (H3Index, H3Error) {
	u, err := strconv.ParseUint(s, 16, 64)
	if err != nil {
		return 0, 1
	}
	return H3Index(u), 0
}
