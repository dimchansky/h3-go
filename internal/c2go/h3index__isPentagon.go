package c2go

// isPentagon returns 1 if the index is a pentagon, else 0 (parity with C).
func isPentagon(h H3Index) int {
	basePent := _isBaseCellPentagon(getBaseCellNumber(h)) != 0
	leading := _h3LeadingNonZeroDigit(h)
	if basePent && leading == 0 {
		return 1
	}
	return 0
}
