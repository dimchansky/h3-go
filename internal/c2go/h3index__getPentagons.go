package c2go

// getPentagons returns the pentagon indices at resolution res.
// Ports H3_EXPORT(getPentagons).
func getPentagons(res int) ([]H3Index, uint32) {
    if res < 0 || res > MAX_H3_RES {
        return nil, _eResDomain
    }
    out := make([]H3Index, 0, NUM_PENTAGONS)
    for bc := 0; bc < NUM_BASE_CELLS; bc++ {
        if _isBaseCellPentagon(bc) != 0 {
            var p H3Index
            setH3Index(&p, res, bc, 0)
            out = append(out, p)
        }
    }
    return out, _eSuccess
}

