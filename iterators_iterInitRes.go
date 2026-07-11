package h3

// iterInitRes creates an iterator for all cells at given resolution.
// Starts iterating from base cell 0 and progresses through all base cells.
// Ported from H3 C: iterators.c::iterInitRes.
func iterInitRes(res int32) IterCellsResolution {
	itC := iterInitBaseCellNum(0, res)

	return IterCellsResolution{
		H:           itC.H,
		baseCellNum: 0,
		res:         res,
		itC:         itC,
	}
}
