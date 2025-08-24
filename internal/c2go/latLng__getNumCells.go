package c2go

// getNumCells computes the number of cells at a given resolution.
//
// This function calculates the total number of H3 cells at the specified
// resolution using the formula: 2 + 120 * 7^res, where the icosahedron
// has 2 pentagons at the poles and 120 faces, each subdivided 7^res times
// per resolution level.
// Ported from H3 C: latLng.c::getNumCells
func getNumCells(res int32) (int64, H3Error) {
	if res < 0 || res > MAX_H3_RES {
		return 0, E_RES_DOMAIN
	}
	return 2 + 120*_ipow(7, int64(res)), E_SUCCESS
}
