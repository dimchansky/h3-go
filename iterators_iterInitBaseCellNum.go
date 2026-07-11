package h3

// iterInitBaseCellNum creates an iterator for children of a base cell at given resolution.
// Returns null iterator if baseCellNum is invalid (< 0 or >= numBaseCells) or
// childRes is invalid (< 0 or > maxH3Res).
// Ported from H3 C: iterators.c::iterInitBaseCellNum.
func iterInitBaseCellNum(baseCellNum int32, childRes int32) iterCellsChildren {
	if baseCellNum < 0 || baseCellNum >= numBaseCells || childRes < 0 || childRes > maxH3Res {
		return nullIter()
	}

	var baseCell h3Index
	setH3Index(&baseCell, 0, baseCellNum, 0)

	var iter iterCellsChildren
	iterInitParent(baseCell, childRes, &iter)
	return iter
}
