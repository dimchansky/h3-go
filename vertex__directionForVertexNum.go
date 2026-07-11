package h3

// directionForVertexNum gets the direction for a given vertex number. This returns the direction for
// the neighbor between the given vertex number and the next number in sequence.
// Returns invalidDigit if the vertex number is invalid.
// Ported from H3 C: vertex.c::directionForVertexNum.
func directionForVertexNum(origin h3Index, vertexNum int32) direction {
	isPent := isPentagon(origin)
	// Check for invalid vertexes
	maxVerts := int32(numHexVerts)
	if isPent {
		maxVerts = int32(numPentVerts)
	}
	if vertexNum < 0 || vertexNum > maxVerts-1 {
		return invalidDigit
	}

	// Determine the vertex rotations for this cell
	var rotations int32
	err := vertexRotations(origin, &rotations)
	if err != eSuccess {
		return invalidDigit
	}

	// Find the appropriate direction, rotating CW if necessary
	if isPent {
		return vertexNumToDirectionPent[(vertexNum+rotations)%int32(numPentVerts)]
	} else {
		return vertexNumToDirectionHex[(vertexNum+rotations)%int32(numHexVerts)]
	}
}
