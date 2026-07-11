package h3

// isValidVertex determines whether the input is a valid H3 vertex.
// Returns true if the vertex index is valid, false otherwise.
// Ported from H3 C: vertex.c::isValidVertex.
func isValidVertex(vertex h3Index) bool {
	if getMode(vertex) != h3VertexMode {
		return false
	}

	vertexNum := getReservedBits(vertex)
	owner := vertex
	owner = setMode(owner, h3CellMode)
	owner = setReservedBits(owner, 0)

	if !isValidCell(owner) {
		return false
	}

	// The easiest way to ensure that the owner + vertex number is valid,
	// and that the vertex is canonical, is to recreate and compare.
	var canonical h3Index
	if cellToVertex(owner, vertexNum, &canonical) != eSuccess {
		return false
	}

	return vertex == canonical
}
