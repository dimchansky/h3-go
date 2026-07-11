package h3

// cellToCenterChild produces the center child index for a given H3 index at childRes.
// Returns (child, err) with H3Error parity codes.
// Ported from H3 C: h3Index.c::cellToCenterChild.
func cellToCenterChild(h H3Index, childRes int32) (H3Index, H3Error) {
	parentRes := getResolution(h)
	if childRes < parentRes || childRes > MAX_H3_RES {
		return 0, E_RES_DOMAIN
	}
	// Zero digits from parentRes+1..childRes
	h = _zeroIndexDigits(h, parentRes+1, childRes)
	// Set resolution to childRes
	h = setResolution(h, childRes)
	return h, E_SUCCESS
}
