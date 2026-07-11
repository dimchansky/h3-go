package h3

// describeH3Error returns a string describing the H3Error.
// The string is statically allocated and should not be freed.
// Ported from H3 C: h3Index.c::describeH3Error.
func describeH3Error(err H3Error) string {
	// C checks err >= 0 && err <= 15; H3Error is uint32 so >= 0 always holds.
	if err <= 15 {
		return h3ErrorDescriptions[err]
	} else {
		return "Invalid error code"
	}
}
