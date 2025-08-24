package c2go

// directionForVertexNum gets the direction for a given vertex number. This returns the direction for
// the neighbor between the given vertex number and the next number in sequence.
// Returns INVALID_DIGIT if the vertex number is invalid.
// Ported from H3 C: vertex.c::directionForVertexNum
func directionForVertexNum(origin H3Index, vertexNum int32) Direction {
	isPent := isPentagon(origin)
	// Check for invalid vertexes
	maxVerts := int32(NUM_HEX_VERTS)
	if isPent {
		maxVerts = int32(NUM_PENT_VERTS)
	}
	if vertexNum < 0 || vertexNum > maxVerts-1 {
		return INVALID_DIGIT
	}

	// Determine the vertex rotations for this cell
	var rotations int32
	err := vertexRotations(origin, &rotations)
	if err != E_SUCCESS {
		return INVALID_DIGIT
	}

	// Find the appropriate direction, rotating CW if necessary
	if isPent {
		return vertexNumToDirectionPent[(vertexNum+rotations)%int32(NUM_PENT_VERTS)]
	} else {
		return vertexNumToDirectionHex[(vertexNum+rotations)%int32(NUM_HEX_VERTS)]
	}
}
