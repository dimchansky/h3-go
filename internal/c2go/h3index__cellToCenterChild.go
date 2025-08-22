package c2go

// cellToCenterChild produces the center child index for a given H3 index at childRes.
// Returns (child, err) with H3Error parity codes.
func cellToCenterChild(h H3Index, childRes int) (H3Index, uint32) {
	parentRes := getResolution(h)
	if childRes < parentRes || childRes > _maxH3Res {
		return 0, _eResDomain
	}
	// Zero digits from parentRes+1..childRes
	h = _zeroIndexDigits(h, parentRes+1, childRes)
	// Set resolution to childRes
	const resOffset = 52
	const resMask = uint64(15) << resOffset
	x := uint64(h)
	x &^= resMask
	x |= (uint64(childRes) & 15) << resOffset
	return H3Index(x), _eSuccess
}
