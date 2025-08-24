package c2go

// gridDiskDistancesUnsafe produces indexes within k distance of the origin index.
// Output behavior is undefined when one of the indexes returned by this function is a pentagon or is in the pentagon distortion area.
// k-ring 0 is defined as the origin index, k-ring 1 is defined as k-ring 0 and all neighboring indexes, and so on.
// Output is placed in the provided array in order of increasing distance from the origin.
// The distances in hexagons is placed in the distances array at the same offset.
// Ported from H3 C: algos.c::gridDiskDistancesUnsafe
func gridDiskDistancesUnsafe(origin H3Index, k int32, out []H3Index, distances []int32) H3Error {
	// Return codes:
	// 1 Pentagon was encountered
	// 2 Pentagon distortion (deleted k subsequence) was encountered
	// Pentagon being encountered is not itself a problem; really the deleted
	// k-subsequence is the problem, but for compatibility reasons we fail on
	// the pentagon.
	if k < 0 {
		return E_DOMAIN
	}

	// k must be >= 0, so origin is always needed
	idx := int32(0)
	out[idx] = origin
	if len(distances) > 0 {
		distances[idx] = 0
	}
	idx++

	if isPentagon(origin) {
		// Pentagon was encountered; bail out as user doesn't want this.
		return E_PENTAGON
	}

	// 0 < ring <= k, current ring
	ring := int32(1)
	// 0 <= direction < 6, current side of the ring
	direction := int32(0)
	// 0 <= i < ring, current position on the side of the ring
	i := int32(0)
	// Number of 60 degree ccw rotations to perform on the direction (based on
	// which faces have been crossed.)
	rotations := int32(0)

	for ring <= k {
		if direction == 0 && i == 0 {
			// Not putting in the output set as it will be done later, at
			// the end of this ring.
			neighborResult := h3NeighborRotations(origin, NEXT_RING_DIRECTION, &rotations, &origin)
			if neighborResult != E_SUCCESS {
				// Should not be possible because `origin` would have to be a
				// pentagon
				// TODO: Reachable via fuzzer
				return neighborResult
			}

			if isPentagon(origin) {
				// Pentagon was encountered; bail out as user doesn't want this.
				return E_PENTAGON
			}
		}

		neighborResult := h3NeighborRotations(origin, DIRECTIONS[direction], &rotations, &origin)
		if neighborResult != E_SUCCESS {
			return neighborResult
		}
		out[idx] = origin
		if len(distances) > 0 {
			distances[idx] = ring
		}
		idx++

		i++
		// Check if end of this side of the k-ring
		if i == ring {
			i = 0
			direction++
			// Check if end of this ring.
			if direction == 6 {
				direction = 0
				ring++
			}
		}

		if isPentagon(origin) {
			// Pentagon was encountered; bail out as user doesn't want this.
			return E_PENTAGON
		}
	}
	return E_SUCCESS
}
