package h3

// cellToParent produces the parent index for a given H3 index at parentRes.
// Returns (out, err) where err mirrors H3Error codes.
// Ported from H3 C: h3Index.c::cellToParent.
func cellToParent(h H3Index, parentRes int32) (H3Index, H3Error) {
	childRes := getResolution(h)
	if parentRes < 0 || parentRes > MAX_H3_RES {
		return 0, E_RES_DOMAIN
	} else if parentRes > childRes {
		return 0, E_RES_MISMATCH
	} else if parentRes == childRes {
		return h, E_SUCCESS
	}
	// Set resolution to parentRes
	parentH := setResolution(h, parentRes)
	// Set digits above parentRes to 7
	for i := parentRes + 1; i <= childRes; i++ {
		parentH = setIndexDigit(parentH, i, int32(H3_DIGIT_MASK))
	}
	return parentH, E_SUCCESS
}
