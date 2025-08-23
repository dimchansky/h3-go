package c2go

// isPentagon returns true if the index is a pentagon, else false.
// Ported from H3 C: h3Index.c::isPentagon
func isPentagon(h H3Index) bool {
	basePent := _isBaseCellPentagon(getBaseCellNumber(h))
	leading := _h3LeadingNonZeroDigit(h)
	return basePent && leading == 0
}
