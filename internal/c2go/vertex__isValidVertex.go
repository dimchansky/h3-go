package c2go

// _isValidVertex determines whether the input is a valid H3 vertex.
// Returns true if the vertex index is valid, false otherwise.
// Ported from H3 C: vertex.c::isValidVertex
func _isValidVertex(vertex H3Index) bool {
	if getMode(vertex) != H3_VERTEX_MODE {
		return false
	}

	vertexNum := getReservedBits(vertex)
	owner := vertex
	owner = setMode(owner, H3_CELL_MODE)
	owner = setReservedBits(owner, 0)

	if !isValidCell(owner) {
		return false
	}

	// The easiest way to ensure that the owner + vertex number is valid,
	// and that the vertex is canonical, is to recreate and compare.
	var canonical H3Index
	if _cellToVertex(owner, vertexNum, &canonical) != E_SUCCESS {
		return false
	}

	return vertex == canonical
}
