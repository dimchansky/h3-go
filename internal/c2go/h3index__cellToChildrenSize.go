package c2go

// cellToChildrenSize returns the exact number of children for a cell at a given child resolution.
// Returns (count, err) with H3Error parity codes.
func cellToChildrenSize(h H3Index, childRes int) (int64, uint32) {
	if !_hasChildAtRes(h, childRes) {
		return 0, _eResDomain
	}
	n := childRes - getResolution(h)
    if isPentagonGo(h) {
		// 1 + 5 * (7^n - 1) / 6
		return 1 + 5*(_ipow(int64(7), int64(n))-1)/6, _eSuccess
	}
	return _ipow(int64(7), int64(n)), _eSuccess
}

func isPentagonGo(h H3Index) bool {
    return isBaseCellPentagon(getBaseCellNumber(h)) && _h3LeadingNonZeroDigit(h) == 0
}
