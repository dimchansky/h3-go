package h3

// gridRing produces cells with distances of exactly k from the origin cell.
// This function attempts to use the faster gridRingUnsafe algorithm first,
// but falls back to the slower _gridRingInternal if pentagon distortion is encountered.
// Ported from H3 C: algos.c::gridRing.
func gridRing(origin h3Index, k int32, out []h3Index) h3Error {
	// Optimistically try the faster gridRingUnsafe algorithm first
	failed := gridRingUnsafe(origin, k, out)
	if failed == eSuccess {
		return eSuccess
	}
	// Fast algo failed, fall back to slower, correct algo
	// and also wipe out array because contents untrustworthy
	for i := range out {
		out[i] = 0
	}
	return _gridRingInternal(origin, k, out)
}
