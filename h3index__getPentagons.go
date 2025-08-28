package h3

// getPentagons appends pentagon indices at resolution res into dst and returns the slice.
// Ports H3_EXPORT(getPentagons) using the dst‑buffer pattern for performance.
// Ported from H3 C: h3Index.c::getPentagons
func getPentagons(dst []H3Index, res int32) ([]H3Index, H3Error) {
	if res < 0 || res > MAX_H3_RES {
		// Do not modify dst on error
		return dst, E_RES_DOMAIN
	}
	n := NUM_PENTAGONS
	var out []H3Index
	if cap(dst) >= n {
		out = dst[:n]
	} else {
		out = make([]H3Index, n)
	}
	i := 0
	for bc := int32(0); bc < NUM_BASE_CELLS; bc++ {
		if _isBaseCellPentagon(bc) {
			var p H3Index
			setH3Index(&p, res, bc, 0)
			out[i] = p
			i++
		}
	}
	return out, E_SUCCESS
}
