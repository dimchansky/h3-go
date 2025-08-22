package c2go

// isResClassIII returns 1 if the resolution of h is Class III (odd), else 0.
func isResClassIII(h H3Index) int {
	if getResolution(h)%2 != 0 {
		return 1
	}
	return 0
}
