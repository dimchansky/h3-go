package h3

// gridRingUnsafe produces cells with distances of k from the origin cell.
// Output cells are placed in the provided array in no particular order.
// Elements of the output array may be left zero, as can happen when crossing a pentagon.
//
// A nonzero failure code may be returned in some cases, for example,
// if a pentagon is encountered. Failure cases may be fixed in future versions.
//
// Ported from H3 C: algos.c::gridRingUnsafe.
func gridRingUnsafe(origin H3Index, k int32, out []H3Index) H3Error {
	if k < 0 {
		return E_DOMAIN
	}
	// Short-circuit on 'identity' ring
	if k == 0 {
		out[0] = origin
		return E_SUCCESS
	}
	idx := int32(0)
	// Number of 60 degree ccw rotations to perform on the direction (based on
	// which faces have been crossed.)
	rotations := int32(0)
	// Scratch structure for checking for pentagons
	if isPentagon(origin) {
		// Pentagon was encountered; bail out as user doesn't want this.
		return E_PENTAGON
	}
	for ring := int32(0); ring < k; ring++ {
		neighborResult := h3NeighborRotations(
			origin, NEXT_RING_DIRECTION, &rotations, &origin)
		if neighborResult != E_SUCCESS {
			// Should not be possible because `origin` would have to be a
			// pentagon
			// TODO: Reachable via fuzzer
			return neighborResult
		}
		if isPentagon(origin) {
			return E_PENTAGON
		}
	}
	lastIndex := origin
	out[idx] = origin
	idx++
	for direction := int32(0); direction < 6; direction++ {
		for pos := int32(0); pos < k; pos++ {
			neighborResult := h3NeighborRotations(
				origin, DIRECTIONS[direction], &rotations, &origin)
			if neighborResult != E_SUCCESS {
				// Should not be possible because `origin` would have to be a
				// pentagon
				// TODO: Reachable via fuzzer
				return neighborResult
			}
			// Skip the very last index, it was already added. We do
			// however need to traverse to it because of the pentagonal
			// distortion check, below.
			if pos != k-1 || direction != 5 {
				out[idx] = origin
				idx++
				if isPentagon(origin) {
					return E_PENTAGON
				}
			}
		}
	}
	// Check that this matches the expected lastIndex, if it doesn't,
	// it indicates pentagonal distortion occurred and we should report
	// failure.
	if lastIndex != origin {
		return E_PENTAGON
	}
	return E_SUCCESS
}
