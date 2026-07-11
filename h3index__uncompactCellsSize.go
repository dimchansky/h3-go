package h3

// uncompactCellsSize takes a compacted set of hexagons and provides
// the exact size of the uncompacted set of hexagons.
// Returns (numOut, err) where err mirrors h3Error codes.
// Ported from H3 C: h3Index.c::uncompactCellsSize.
func uncompactCellsSize(compactedSet []h3Index, numCompacted int64, res int32) (int64, h3Error) {
	numOut := int64(0)
	for i := int64(0); i < numCompacted; i++ {
		if compactedSet[i] == h3Null {
			continue
		}

		childrenSize, childrenError := cellToChildrenSize(compactedSet[i], res)
		if childrenError != eSuccess {
			// The parent res does not contain `res`.
			return 0, eResMismatch
		}
		numOut += childrenSize
	}
	return numOut, eSuccess
}
