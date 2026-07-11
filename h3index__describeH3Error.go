package h3

// describeH3Error returns a string describing the h3Error.
// The string is statically allocated and should not be freed.
// Ported from H3 C: h3Index.c::describeH3Error.
func describeH3Error(err h3Error) string {
	// err is always non-negative because it is an unsigned integer (C 4.4.0
	// checks err < H3_ERROR_END).
	if err < h3ErrorEnd {
		return h3ErrorDescriptions[err]
	}
	return "Invalid error code"
}
