package h3

// cellToVertexes gets all vertexes for the given cell.
// If the cell is a pentagon, will fill the final slot with h3Null.
// Ported from H3 C: vertex.c::cellToVertexes.
func cellToVertexes(cell h3Index, vertexes *[6]h3Index) h3Error {
	// Get all vertexes. If the cell is a pentagon, will fill the final slot
	// with h3Null.
	isPent := isPentagon(cell)
	for i := int32(0); i < numHexVerts; i++ {
		if i == 5 && isPent {
			vertexes[i] = h3Null
		} else {
			cellError := cellToVertex(cell, i, &vertexes[i])
			if cellError != eSuccess {
				return cellError
			}
		}
	}
	return eSuccess
}
