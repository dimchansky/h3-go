package h3

// cmp_uint64 compares H3Index values, interpreting them as uint64s.
//
// Note that, usually, we only use this ordering when we know that the
// cells in the set are all the same resolution.
// Ported from H3 C: cellsToMultiPoly.h::cmp_uint64.
func cmp_uint64(a, b h3Index) int32 {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}
