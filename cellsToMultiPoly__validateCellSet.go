package h3

import "sort"

// validateCellSet validates the input cell set for cellsToMultiPolygon:
// numCells < 0 fails with eDomain, an invalid cell with eCellInvalid
// (checked before resolution), a resolution mismatch against cells[0]
// with eResMismatch, and a duplicate cell with eDuplicateInput.
// Ported from H3 C: cellsToMultiPoly.c::validateCellSet.
func validateCellSet(cells []h3Index, numCells int64) h3Error {
	if numCells < 0 {
		return eDomain
	}
	if numCells == 0 {
		return eSuccess
	}

	// Check that all cells are valid and have the same resolution
	res := getResolution(cells[0])
	for i := int64(0); i < numCells; i++ {
		if !isValidCell(cells[i]) {
			return eCellInvalid
		}
		if getResolution(cells[i]) != res {
			return eResMismatch
		}
	}

	// Check for duplicate cells by sorting a copy and looking for adjacent
	// duplicates
	if numCells >= 2 {
		cellsCopy := make([]h3Index, numCells)
		copy(cellsCopy, cells[:numCells])
		sort.Slice(cellsCopy, func(a, b int) bool {
			return cmp_uint64(cellsCopy[a], cellsCopy[b]) < 0
		})
		for i := int64(1); i < numCells; i++ {
			if cellsCopy[i] == cellsCopy[i-1] {
				return eDuplicateInput
			}
		}
	}

	return eSuccess
}
