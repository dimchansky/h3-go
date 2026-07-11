package h3

// _gridDiskDistancesInternal implements the internal algorithm for the safe but slow version of gridDiskDistances.
//
// Adds the origin cell to the output set (treating it as a hash set)
// and recurses to its neighbors, if needed.
//
// The function uses a hash table approach with linear probing to store cells and their distances,
// recursively exploring neighbors up to distance k from the origin.
//
// Ported from H3 C: algos.c::_gridDiskDistancesInternal.
func _gridDiskDistancesInternal(origin h3Index, k int32, out []h3Index, distances []int32, maxIdx int64, curK int32) h3Error {
	// Put origin in the output array. out is used as a hash set.
	// Note: In C, this is int64_t off = origin % maxIdx; where origin is uint64_t
	// The C behavior with signed/unsigned conversion is preserved here
	off := int64(uint64(origin) % uint64(maxIdx))
	for out[off] != 0 && out[off] != origin {
		off = (off + 1) % maxIdx
	}

	// We either got a free slot in the hash set or hit a duplicate
	// We might need to process the duplicate anyways because we got
	// here on a longer path before.
	if out[off] == origin && distances[off] <= curK {
		return eSuccess
	}

	out[off] = origin
	distances[off] = curK

	// Base case: reached an index k away from the origin.
	if curK >= k {
		return eSuccess
	}

	// Recurse to all neighbors in no particular order.
	for i := 0; i < 6; i++ {
		rotations := int32(0)
		var nextNeighbor h3Index
		neighborResult := h3NeighborRotations(origin, algosDirections[i], &rotations, &nextNeighbor)
		if neighborResult != ePentagon {
			// ePentagon is an expected case when trying to traverse off of
			// pentagons.
			if neighborResult != eSuccess {
				return neighborResult
			}
			neighborResult = _gridDiskDistancesInternal(nextNeighbor, k, out, distances, maxIdx, curK+1)
			if neighborResult != eSuccess {
				return neighborResult
			}
		}
	}
	return eSuccess
}
