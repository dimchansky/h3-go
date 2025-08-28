package h3

// _cellToVertexes gets all vertexes for the given cell.
// If the cell is a pentagon, will fill the final slot with H3_NULL.
// Ported from H3 C: vertex.c::cellToVertexes
func _cellToVertexes(cell H3Index, vertexes *[6]H3Index) H3Error {
	// Get all vertexes. If the cell is a pentagon, will fill the final slot
	// with H3_NULL.
	isPent := isPentagon(cell)
	for i := int32(0); i < NUM_HEX_VERTS; i++ {
		if i == 5 && isPent {
			vertexes[i] = H3_NULL
		} else {
			cellError := _cellToVertex(cell, i, &vertexes[i])
			if cellError != E_SUCCESS {
				return cellError
			}
		}
	}
	return E_SUCCESS
}
