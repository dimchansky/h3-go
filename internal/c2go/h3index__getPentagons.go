package c2go

// getPentagons appends pentagon indices at resolution res into dst and returns the slice.
// Ports H3_EXPORT(getPentagons) using the dst‑buffer pattern for performance.
func getPentagons(dst []H3Index, res int) ([]H3Index, uint32) {
    if res < 0 || res > MAX_H3_RES {
        // Do not modify dst on error
        return dst, _eResDomain
    }
    n := NUM_PENTAGONS
    var out []H3Index
    if cap(dst) >= n {
        out = dst[:n]
    } else {
        out = make([]H3Index, n)
    }
    i := 0
    for bc := 0; bc < NUM_BASE_CELLS; bc++ {
        if _isBaseCellPentagon(bc) != 0 {
            var p H3Index
            setH3Index(&p, res, bc, 0)
            out[i] = p
            i++
        }
    }
    return out, _eSuccess
}
