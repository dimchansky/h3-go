package c2go

// isResolutionClassIII returns 1 if res is a Class III grid (odd), else 0.
// Ported from H3 C: h3Index.c::isResolutionClassIII
func isResolutionClassIII(res int) int { return res % 2 }
