package h3

// gridPathCellsSize calculates the number of indexes in a line from the start
// index to the end index, to be used for allocating memory. Returns a negative
// number if the line cannot be computed.
// Ported from H3 C: localij.c::gridPathCellsSize.
func gridPathCellsSize(start H3Index, end H3Index, size *int64) H3Error {
	var distance int64
	distanceError := gridDistance(start, end, &distance)
	if distanceError != 0 {
		return distanceError
	}
	*size = distance + 1
	return 0 // E_SUCCESS
}
