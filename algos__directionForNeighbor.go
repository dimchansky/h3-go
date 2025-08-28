package h3

// _directionForNeighbor returns the direction from origin to destination neighbor.
// Checks each neighbor, in order, to determine which direction the
// destination neighbor is located. Skips CENTER_DIGIT since that
// would be the origin; skips deleted K direction for pentagons.
// Ported from H3 C: algos.c::directionForNeighbor
func _directionForNeighbor(origin H3Index, destination H3Index) Direction {
	isPent := isPentagon(origin)
	// Checks each neighbor, in order, to determine which direction the
	// destination neighbor is located. Skips CENTER_DIGIT since that
	// would be the origin; skips deleted K direction for pentagons.
	var startDirection Direction
	if isPent {
		startDirection = J_AXES_DIGIT
	} else {
		startDirection = K_AXES_DIGIT
	}

	for direction := startDirection; direction < NUM_DIGITS; direction++ {
		var neighbor H3Index
		rotations := int32(0)
		neighborError := h3NeighborRotations(origin, direction, &rotations, &neighbor)
		if neighborError == E_SUCCESS && neighbor == destination {
			return direction
		}
	}
	return INVALID_DIGIT
}
