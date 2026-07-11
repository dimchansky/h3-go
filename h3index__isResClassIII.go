package h3

// isResClassIII returns true if the resolution of h is Class III (odd), else false.
// Ported from H3 C: h3Index.c::isResClassIII.
func isResClassIII(h h3Index) bool {
	return getResolution(h)%2 != 0
}
