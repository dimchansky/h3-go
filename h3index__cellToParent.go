package h3

// cellToParent produces the parent index for a given H3 index at parentRes.
// Returns (out, err) where err mirrors h3Error codes.
// Ported from H3 C: h3Index.c::cellToParent.
func cellToParent(h h3Index, parentRes int32) (h3Index, h3Error) {
	childRes := getResolution(h)
	if parentRes < 0 || parentRes > maxH3Res {
		return 0, eResDomain
	} else if parentRes > childRes {
		return 0, eResMismatch
	} else if parentRes == childRes {
		return h, eSuccess
	}
	// Set resolution to parentRes
	parentH := setResolution(h, parentRes)
	// Set digits above parentRes to 7
	for i := parentRes + 1; i <= childRes; i++ {
		parentH = setIndexDigit(parentH, i, int32(h3DigitMask))
	}
	return parentH, eSuccess
}
