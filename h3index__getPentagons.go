package h3

// getPentagons generates all pentagons at the specified resolution.
// Ported from H3 C: h3Index.c::getPentagons.
func getPentagons(res int32, out []h3Index) h3Error {
	if res < 0 || res > maxH3Res {
		return eResDomain
	}
	i := 0
	for bc := int32(0); bc < numBaseCells; bc++ {
		if _isBaseCellPentagon(bc) {
			var pentagon h3Index
			setH3Index(&pentagon, res, bc, 0)
			out[i] = pentagon
			i++
		}
	}
	return eSuccess
}
