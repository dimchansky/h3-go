package h3

// Element sizes for the overflow pre-check. The struct sizes are the
// LP64 C oracle's (sizeof(Arc) = 48: 8 id + 2 bools + 6 padding + 3×8
// pointers + 8 rank; 8-byte pointers); on 32-bit Go targets the real
// element sizes are smaller, so using these values makes the check
// strictly more conservative there (it rejects earlier than a 32-bit C
// build would), never less safe. cSizeMax is the platform's SIZE_MAX
// equivalent (^uint(0)): 2^64-1 on 64-bit targets — where it reproduces
// the C oracle's threshold bit-for-bit (the parity suite runs on
// 64-bit platforms only) — and 2^32-1 on 32-bit targets, where the
// stricter threshold guards Go's own allocation-length limits.
const (
	cSizeofArc    = 48
	cSizeofArcPtr = 8
	cSizeMax      = uint64(^uint(0))
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
	// bits (cSizeMax mirrors the platform width; see above).
	if numCells > 0 && uint64(numCells) > cSizeMax/maxBytesPerCell {
		return eMemoryBounds
	}

	return eSuccess
}
