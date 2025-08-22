package c2go

// Local H3Error constants used in ports
const (
    _eFailed = uint32(1)
    _eDomain = uint32(2)
)

// validateChildPos checks that childPos is in range for the parent's children at childRes.
// Mirrors static validateChildPos in h3Index.c
func validateChildPos(childPos int64, parent H3Index, childRes int) uint32 {
    maxChildCount, sizeErr := cellToChildrenSize(parent, childRes)
    if sizeErr != _eSuccess {
        return sizeErr
    }
    if childPos < 0 || childPos >= maxChildCount {
        return _eDomain
    }
    return _eSuccess
}
