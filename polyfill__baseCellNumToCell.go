package h3

// baseCellNumToCell gets a base cell by number, or h3Null if out of bounds.
// Creates an H3 index for the specified base cell at resolution 0.
// Ported from H3 C: polyfill.c::baseCellNumToCell.
func baseCellNumToCell(baseCellNum int32) h3Index {
	if baseCellNum < 0 || baseCellNum >= numBaseCells {
		return h3Null
	}
	var baseCell h3Index
	setH3Index(&baseCell, 0, baseCellNum, 0)
	return baseCell
}
