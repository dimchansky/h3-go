package h3

// iterInitBaseCellNum creates an iterator for children of a base cell at given resolution.
// Returns null iterator if baseCellNum is invalid (< 0 or >= NUM_BASE_CELLS) or
// childRes is invalid (< 0 or > MAX_H3_RES).
// Ported from H3 C: iterators.c::iterInitBaseCellNum
func iterInitBaseCellNum(baseCellNum int32, childRes int32) IterCellsChildren {
	if baseCellNum < 0 || baseCellNum >= NUM_BASE_CELLS || childRes < 0 || childRes > MAX_H3_RES {
		return nullIter()
	}

	var baseCell H3Index
	setH3Index(&baseCell, 0, baseCellNum, 0)

	var iter IterCellsChildren
	iterInitParent(baseCell, childRes, &iter)
	return iter
}
