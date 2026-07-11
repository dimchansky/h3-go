package h3

// getIndexDigit returns the index digit at res, which starts with 1 for
// resolution 1, up to and including resolution 15.
//
// 0 is not a valid value for res because resolution 0 is specified by the
// base cell number, not an indexing digit. res may exceed the actual
// resolution of the index, in which case the actual digit stored in the
// index is returned (7 for valid cell indexes).
// Ported from H3 C: h3Index.c::getIndexDigit.
func getIndexDigit(h h3Index, res int32, out *int32) h3Error {
	if res < 1 || res > maxH3Res {
		return eResDomain
	}
	*out = h3GetIndexDigit(h, res)
	return eSuccess
}
