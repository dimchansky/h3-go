package h3

// _maxGridRingSize returns the number of cells that result from the gridRing algorithm with the given k.
// For k=0, returns 1 (just the origin cell). For k>0, returns 6*k (the hollow ring of hexagons).
// Ported from H3 C: algos.c::maxGridRingSize.
func _maxGridRingSize(k int32, out *int64) h3Error {
	if k < 0 {
		return eDomain
	}
	if k == 0 {
		*out = 1
		return eSuccess
	}
	*out = 6 * int64(k)
	return eSuccess
}
