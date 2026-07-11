package h3

// iterStepChild steps a iterCellsChildren to the next child cell.
// When the iteration is over, iterCellsChildren.H will be h3Null.
// Handles iterating through hexagon and pentagon cells.
// Ported from H3 C: iterators.c::iterStepChild.
func iterStepChild(it *iterCellsChildren) {
	// once h == h3Null, the iterator returns an infinite sequence of h3Null
	if it.H == h3Null {
		return
	}
	childRes := getResolution(it.H)
	// increment the digit at the child resolution
	incrementResDigit(&it.H, childRes)

	for i := childRes; i >= it.ParentRes; i-- {
		if i == it.ParentRes {
			// if we're modifying the parent resolution digit, then we're done
			*it = nullIter()
			return
		}
		// pentagonSkippedDigit == 1
		if i == it.SkipDigit && getIndexDigit(it.H, i) == int32(pentagonSkippedDigit) {
			// Then we are iterating through the children of a pentagon cell.
			// All children of a pentagon have the property that the first
			// nonzero digit between the parent and child resolutions is
			// not 1.
			// I.e., we never see a sequence like 00001.
			// Thus, we skip the `1` in this digit.
			incrementResDigit(&it.H, i)
			it.SkipDigit -= 1
			return
		}
		// invalidDigit == 7
		if getIndexDigit(it.H, i) == int32(invalidDigit) {
			// We have exhausted all the children for the current digit at
			// resolution i. Call incrementResDigit which will wrap the digit
			// from 7 to 0 and carry to the next resolution (i-1).
			incrementResDigit(&it.H, i)
		} else {
			return
		}
	}
	// This should never happen in valid iteration
	*it = nullIter()
}
