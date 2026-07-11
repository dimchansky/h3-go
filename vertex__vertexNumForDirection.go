package h3

// vertexNumForDirection gets the first vertex number for a given direction.
// The neighbor in this direction is located between this vertex number and
// the next number in sequence.
// Returns the number for the first topological vertex, or invalidVertexNum
// if the direction is not valid for this cell.
// Ported from H3 C: vertex.c::vertexNumForDirection.
func vertexNumForDirection(origin h3Index, direction direction) int32 {
	isPent := isPentagon(origin)

	// Check for invalid directions
	if direction == centerDigit || direction >= invalidDigit ||
		(isPent && direction == kAxesDigit) {
		return invalidVertexNum
	}

	// Determine the vertex rotations for this cell
	var rotations int32
	err := vertexRotations(origin, &rotations)
	if err != eSuccess {
		return invalidVertexNum
	}

	// Find the appropriate vertex, rotating CCW if necessary
	if isPent {
		return (directionToVertexNumPent[direction] + numPentVerts - rotations) % numPentVerts
	} else {
		return (directionToVertexNumHex[direction] + numHexVerts - rotations) % numHexVerts
	}
}
