package h3

// cellToCenterChild produces the center child index for a given H3 index at childRes.
// Returns (child, err) with h3Error parity codes.
// Ported from H3 C: h3Index.c::cellToCenterChild.
func cellToCenterChild(h h3Index, childRes int32) (h3Index, h3Error) {
	parentRes := getResolution(h)
	if childRes < parentRes || childRes > maxH3Res {
		return 0, eResDomain
	}
	// Zero digits from parentRes+1..childRes
	h = _zeroIndexDigits(h, parentRes+1, childRes)
	// Set resolution to childRes
	h = setResolution(h, childRes)
	return h, eSuccess
}
