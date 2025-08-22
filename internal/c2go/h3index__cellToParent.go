package c2go

// cellToParent produces the parent index for a given H3 index at parentRes.
// Returns (out, err) where err mirrors H3Error codes.
func cellToParent(h H3Index, parentRes int) (H3Index, H3Error) {
	childRes := getResolution(h)
	if parentRes < 0 || parentRes > MAX_H3_RES {
		return 0, E_RES_DOMAIN
	} else if parentRes > childRes {
		return 0, E_RES_MISMATCH
	} else if parentRes == childRes {
		return h, E_SUCCESS
	}
	// Set resolution to parentRes
	const resOffset = 52
	const resMask = uint64(15) << resOffset
	x := uint64(h)
	x &^= resMask
	x |= (uint64(parentRes) & 15) << resOffset
	parentH := H3Index(x)
	// Set digits above parentRes to 7
	for i := parentRes + 1; i <= childRes; i++ {
		parentH = setIndexDigit(parentH, i, int(H3_DIGIT_MASK))
	}
	return parentH, E_SUCCESS
}
