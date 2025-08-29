package h3

// iterStepRes steps a IterCellsResolution iterator to the next cell.
// When the iteration is exhausted, the iterator's H field becomes H3_NULL.
// This iterator cycles through all cells at the target resolution by iterating
// through all base cells and their children.
// Ported from H3 C: iterators.c::iterStepRes
func iterStepRes(itR *IterCellsResolution) {
	// reached the end of over iterator; emits H3_NULL from now on
	if itR.H == H3_NULL {
		return
	}

	// step child iterator
	iterStepChild(&itR.itC)

	// If the child iterator is exhausted and there are still
	// base cells remaining, we initialize the next base cell child iterator
	if itR.itC.H == H3_NULL && itR.baseCellNum+1 < NUM_BASE_CELLS {
		itR.baseCellNum += 1
		itR.itC = iterInitBaseCellNum(itR.baseCellNum, itR.res)
	}

	// This overall iterator reflects the next cell in the child iterator.
	// Note: This sets itR.H = H3_NULL if the base cells were
	// exhausted in the check above.
	itR.H = itR.itC.H
}
