package h3

// _iterateAllIndexesAtResPartial iterates through all H3 indexes at the given resolution
// for base cells from 0 up to (but not including) the specified baseCells limit.
// The callback function is called for each valid H3 index encountered.
// This function is useful for testing or partial iteration through the H3 grid.
// Ported from H3 C: utility.c::iterateAllIndexesAtResPartial.
func _iterateAllIndexesAtResPartial(res int32, callback func(H3Index), baseCells int32) {
	// Assert equivalent: ensure baseCells doesn't exceed maximum
	if baseCells > NUM_BASE_CELLS {
		baseCells = NUM_BASE_CELLS
	}

	for i := int32(0); i < baseCells; i++ {
		_iterateBaseCellIndexesAtRes(res, callback, i)
	}
}
