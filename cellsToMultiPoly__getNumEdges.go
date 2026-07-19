package h3

// getNumEdges returns the total number of directed edges contributed by
// the cell set: 6 per cell, minus 1 for each pentagon.
// Ported from H3 C: cellsToMultiPoly.c::getNumEdges.
func getNumEdges(cells []h3Index, numCells int64) int64 {
	numEdges := 6 * numCells

	for i := int64(0); i < numCells; i++ {
		if isPentagon(cells[i]) {
			numEdges--
		}
	}
	return numEdges
}
