package h3

// iterInitParent initializes a parent iterator in-place for iterating through
// the children of a given parent cell at a given resolution. This function
// sets up the iterator state for stepping through all child cells.
// Ported from H3 C: iterators.c::_iterInitParent.
func iterInitParent(h H3Index, childRes int32, iter *IterCellsChildren) {
	iter.ParentRes = getResolution(h)

	if childRes < iter.ParentRes || childRes > MAX_H3_RES || h == 0 {
		*iter = nullIter()
		return
	}

	iter.H = _zeroIndexDigits(h, iter.ParentRes+1, childRes)
	// Set resolution
	iter.H = setResolution(iter.H, childRes)

	if isPentagon(iter.H) {
		// The skip digit skips `1` for pentagons.
		// The "_skipDigit" moves to the left as we count up from the
		// child resolution to the parent resolution.
		iter.SkipDigit = childRes
	} else {
		// if not a pentagon, we can ignore "skip digit" logic
		iter.SkipDigit = -1
	}
}
