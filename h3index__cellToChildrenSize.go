package h3

// cellToChildrenSize returns the exact number of children for a cell at a given child resolution.
// Returns (count, err) with H3Error parity codes.
// Ported from H3 C: h3Index.c::cellToChildrenSize
func cellToChildrenSize(h H3Index, childRes int32) (int64, H3Error) {
	if !_hasChildAtRes(h, childRes) {
		return 0, E_RES_DOMAIN
	}
	n := childRes - getResolution(h)
	if isPentagon(h) {
		// 1 + 5 * (7^n - 1) / 6
		return 1 + 5*(_ipow(int64(7), int64(n))-1)/6, E_SUCCESS
	}
	return _ipow(int64(7), int64(n)), E_SUCCESS
}
