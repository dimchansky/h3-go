package h3

// cellToLocalIj produces IJ coordinates for an index anchored by an origin.
// The coordinate space used by this function may have deleted regions or warping
// due to pentagonal distortion. Coordinates are only comparable if they come from
// the same origin index.
//
// Failure may occur if the index is too far away from the origin or if the index
// is on the other side of a pentagon.
//
// This function's output is not guaranteed to be compatible across different
// versions of H3.
//
// Ported from H3 C: localij.c::cellToLocalIj.
func cellToLocalIj(origin H3Index, index H3Index, mode uint32, out *CoordIJ) H3Error {
	if mode != 0 {
		return E_OPTION_INVALID
	}
	var ijk CoordIJK
	failed := cellToLocalIjk(origin, index, &ijk)
	if failed != E_SUCCESS {
		return failed
	}

	ijkToIj(&ijk, out)

	return E_SUCCESS
}
