package h3

// vertexNumForDirection gets the first vertex number for a given direction.
// The neighbor in this direction is located between this vertex number and
// the next number in sequence.
// Returns the number for the first topological vertex, or INVALID_VERTEX_NUM
// if the direction is not valid for this cell.
// Ported from H3 C: vertex.c::vertexNumForDirection
func vertexNumForDirection(origin H3Index, direction Direction) int32 {
	isPent := isPentagon(origin)

	// Check for invalid directions
	if direction == CENTER_DIGIT || direction >= INVALID_DIGIT ||
		(isPent && direction == K_AXES_DIGIT) {
		return INVALID_VERTEX_NUM
	}

	// Determine the vertex rotations for this cell
	var rotations int32
	err := vertexRotations(origin, &rotations)
	if err != E_SUCCESS {
		return INVALID_VERTEX_NUM
	}

	// Find the appropriate vertex, rotating CCW if necessary
	if isPent {
		return (directionToVertexNumPent[direction] + NUM_PENT_VERTS - rotations) % NUM_PENT_VERTS
	} else {
		return (directionToVertexNumHex[direction] + NUM_HEX_VERTS - rotations) % NUM_HEX_VERTS
	}
}
