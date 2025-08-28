package h3

// isResolutionClassIII returns true if res is a Class III grid (odd), else false.
// Ported from H3 C: h3Index.c::isResolutionClassIII
// Note: C version returns res % 2, but we return proper boolean for Go idiomaticity
func isResolutionClassIII(res int32) bool { return res%2 != 0 }
