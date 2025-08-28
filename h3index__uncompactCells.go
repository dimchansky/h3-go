package h3

// uncompactCells takes a compressed set of cells and expands back to the
// original set of cells. Skips elements that are H3_NULL.
// The outSet slice is modified in place with the expanded cells.
// Ported from H3 C: h3Index.c::uncompactCells
func uncompactCells(compactedSet []H3Index, numCompacted int64, outSet []H3Index, numOut int64, res int32) H3Error {
	i := int64(0)

	for j := int64(0); j < numCompacted; j++ {
		if !_hasChildAtRes(compactedSet[j], res) {
			return E_RES_MISMATCH
		}

		// Initialize iterator for this compacted cell
		var iter IterCellsChildren
		iterInitParent(compactedSet[j], res, &iter)

		// Iterate through all children
		for iter.H != H3_NULL {
			if i >= numOut {
				return E_MEMORY_BOUNDS // went too far; abort!
			}
			outSet[i] = iter.H
			i++
			iterStepChild(&iter)
		}
	}
	return E_SUCCESS
}
