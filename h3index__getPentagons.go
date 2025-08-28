package h3

// getPentagons generates all pentagons at the specified resolution.
// Ported from H3 C: h3Index.c::getPentagons
func getPentagons(res int32, out []H3Index) H3Error {
	if res < 0 || res > MAX_H3_RES {
		return E_RES_DOMAIN
	}
	i := 0
	for bc := int32(0); bc < NUM_BASE_CELLS; bc++ {
		if _isBaseCellPentagon(bc) {
			var pentagon H3Index
			setH3Index(&pentagon, res, bc, 0)
			out[i] = pentagon
			i++
		}
	}
	return E_SUCCESS
}
