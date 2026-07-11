package h3

// validateChildPos checks that childPos is in range for the parent's children at childRes.
// Mirrors static validateChildPos in h3Index.c
// Ported from H3 C: h3Index.c::validateChildPos.
func validateChildPos(childPos int64, parent h3Index, childRes int32) h3Error {
	maxChildCount, sizeErr := cellToChildrenSize(parent, childRes)
	if sizeErr != eSuccess {
		return sizeErr
	}
	if childPos < 0 || childPos >= maxChildCount {
		return eDomain
	}
	return eSuccess
}
