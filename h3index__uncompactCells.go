package h3

// uncompactCells takes a compressed set of cells and expands back to the
// original set of cells. Skips elements that are h3Null.
// The outSet slice is modified in place with the expanded cells.
// Ported from H3 C: h3Index.c::uncompactCells.
func uncompactCells(compactedSet []h3Index, numCompacted int64, outSet []h3Index, numOut int64, res int32) h3Error {
	i := int64(0)

	for j := int64(0); j < numCompacted; j++ {
		if !_hasChildAtRes(compactedSet[j], res) {
			return eResMismatch
		}

		// Initialize iterator for this compacted cell
		var iter iterCellsChildren
		iterInitParent(compactedSet[j], res, &iter)

		// Iterate through all children
		for iter.H != h3Null {
			if i >= numOut {
				return eMemoryBounds // went too far; abort!
			}
			outSet[i] = iter.H
			i++
			iterStepChild(&iter)
		}
	}
	return eSuccess
}
