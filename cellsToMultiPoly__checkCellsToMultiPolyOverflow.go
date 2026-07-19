package h3

// C sizes on the LP64 platforms the parity oracle runs on: sizeof(Arc)
// = 48 (8 id + 2 bools + 6 padding + 3×8 pointers + 8 rank), a pointer
// is 8 bytes, and SIZE_MAX = 2^64-1. The Go port pins these so the
// overflow threshold matches the C oracle bit-for-bit (verified by the
// parity suite); Go's own allocator limits are unrelated.
const (
	cSizeofArc    = 48
	cSizeofArcPtr = 8
	cSizeMax      = ^uint64(0)
)

// checkCellsToMultiPolyOverflow checks for potential integer overflow in
// cellsToMultiPolygon allocations.
//
// Validates that the two largest allocations won't overflow:
//  1. arcs array: numArcs * sizeof(Arc) where numArcs ~= 6 * numCells
//  2. buckets array: numBuckets * sizeof(Arc *)
//     where numBuckets = numArcs * HASH_TABLE_MULTIPLIER
//
// Returns eSuccess if allocations are safe, eMemoryBounds if overflow
// would occur.
// Ported from H3 C: cellsToMultiPoly.h::checkCellsToMultiPolyOverflow.
func checkCellsToMultiPolyOverflow(numCells int64, hashMultiplier int64) h3Error {
	// Compute the maximum bytes per cell across both allocations
	arcsPerCell := uint64(6 * cSizeofArc)
	bucketsPerCell := uint64(6 * hashMultiplier * cSizeofArcPtr)
	maxBytesPerCell := arcsPerCell
	if bucketsPerCell > maxBytesPerCell {
		maxBytesPerCell = bucketsPerCell
	}

	// Check if maxBytesPerCell * numCells would overflow size_t, which is
	// what is used for allocations. Use SIZE_MAX since size_t may be 32
	// bits (the oracle platforms are 64-bit; see cSizeMax above).
	if numCells > 0 && uint64(numCells) > cSizeMax/maxBytesPerCell {
		return eMemoryBounds
	}

	return eSuccess
}
