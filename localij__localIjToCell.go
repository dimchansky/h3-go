package h3

// localIjToCell produces an H3 index from IJ coordinates anchored by an origin.
// The coordinate space used by this function may have deleted regions or warping
// due to pentagonal distortion. Failure may occur if the index is too far away
// from the origin or if the index is on the other side of a pentagon.
//
// This function's output is not guaranteed to be compatible across different
// versions of H3.
//
// Ported from H3 C: localij.c::localIjToCell
func localIjToCell(origin H3Index, ij *CoordIJ, mode uint32, out *H3Index) H3Error {
	if mode != 0 {
		return E_OPTION_INVALID
	}
	var ijk CoordIJK
	ijToIjkError := ijToIjk(ij, &ijk)
	if ijToIjkError != E_SUCCESS {
		return ijToIjkError
	}

	return localIjkToCell(origin, &ijk, out)
}
