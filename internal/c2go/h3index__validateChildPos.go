package c2go

// validateChildPos checks that childPos is in range for the parent's children at childRes.
// Mirrors static validateChildPos in h3Index.c
// Ported from H3 C: h3Index.c::validateChildPos
func validateChildPos(childPos int64, parent H3Index, childRes int32) H3Error {
	maxChildCount, sizeErr := cellToChildrenSize(parent, childRes)
	if sizeErr != E_SUCCESS {
		return sizeErr
	}
	if childPos < 0 || childPos >= maxChildCount {
		return E_DOMAIN
	}
	return E_SUCCESS
}
