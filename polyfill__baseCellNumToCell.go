package h3

// baseCellNumToCell gets a base cell by number, or H3_NULL if out of bounds.
// Creates an H3 index for the specified base cell at resolution 0.
// Ported from H3 C: polyfill.c::baseCellNumToCell
func baseCellNumToCell(baseCellNum int32) H3Index {
	if baseCellNum < 0 || baseCellNum >= NUM_BASE_CELLS {
		return H3_NULL
	}
	var baseCell H3Index
	setH3Index(&baseCell, 0, baseCellNum, 0)
	return baseCell
}
