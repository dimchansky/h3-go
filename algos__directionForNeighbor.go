package h3

// _directionForNeighbor returns the direction from origin to destination neighbor.
// Checks each neighbor, in order, to determine which direction the
// destination neighbor is located. Skips centerDigit since that
// would be the origin; skips deleted K direction for pentagons.
// Ported from H3 C: algos.c::directionForNeighbor.
func _directionForNeighbor(origin h3Index, destination h3Index) direction {
	isPent := isPentagon(origin)
	// Checks each neighbor, in order, to determine which direction the
	// destination neighbor is located. Skips centerDigit since that
	// would be the origin; skips deleted K direction for pentagons.
	var startDirection direction
	if isPent {
		startDirection = jAxesDigit
	} else {
		startDirection = kAxesDigit
	}

	for direction := startDirection; direction < numDigits; direction++ {
		var neighbor h3Index
		rotations := int32(0)
		neighborError := h3NeighborRotations(origin, direction, &rotations, &neighbor)
		if neighborError == eSuccess && neighbor == destination {
			return direction
		}
	}
	return invalidDigit
}
