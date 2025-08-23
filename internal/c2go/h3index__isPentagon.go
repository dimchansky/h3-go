package c2go

// isPentagon returns 1 if the index is a pentagon, else 0 (parity with C).
// Ported from H3 C: h3Index.c::isPentagon
func isPentagon(h H3Index) int {
	basePent := _isBaseCellPentagon(getBaseCellNumber(h))
	leading := _h3LeadingNonZeroDigit(h)
	if basePent && leading == 0 {
		return 1
	}
	return 0
}
