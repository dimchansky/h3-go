package c2go

// validateChildPos checks that childPos is in range for the parent's children at childRes.
// Mirrors static validateChildPos in h3Index.c
func validateChildPos(childPos int64, parent H3Index, childRes int) H3Error {
	maxChildCount, sizeErr := cellToChildrenSize(parent, childRes)
	if sizeErr != E_SUCCESS {
		return sizeErr
	}
	if childPos < 0 || childPos >= maxChildCount {
		return E_DOMAIN
	}
	return E_SUCCESS
}
