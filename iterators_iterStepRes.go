package h3

// iterStepRes steps a iterCellsResolution iterator to the next cell.
// When the iteration is exhausted, the iterator's H field becomes h3Null.
// This iterator cycles through all cells at the target resolution by iterating
// through all base cells and their children.
// Ported from H3 C: iterators.c::iterStepRes.
func iterStepRes(itR *iterCellsResolution) {
	// reached the end of over iterator; emits h3Null from now on
	if itR.H == h3Null {
		return
	}

	// step child iterator
	iterStepChild(&itR.itC)

	// If the child iterator is exhausted and there are still
	// base cells remaining, we initialize the next base cell child iterator
	if itR.itC.H == h3Null && itR.baseCellNum+1 < numBaseCells {
		itR.baseCellNum += 1
		itR.itC = iterInitBaseCellNum(itR.baseCellNum, itR.res)
	}

	// This overall iterator reflects the next cell in the child iterator.
	// Note: This sets itR.H = h3Null if the base cells were
	// exhausted in the check above.
	itR.H = itR.itC.H
}
