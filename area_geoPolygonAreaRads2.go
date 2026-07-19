package h3

import "math"

// geoPolygonAreaRads2 computes the area of a GeoPolygon in radians^2.
//
// Outer GeoLoop vertices should be in counter-clockwise order, hole
// vertices in clockwise order; the result is the outer area minus the
// holes. No check is made that holes are disjoint or contained within
// the outer loop.
// Ported from H3 C: area.c::geoPolygonAreaRads2.
func geoPolygonAreaRads2(poly GeoPolygon) (float64, h3Error) {
	var a adder

	term, err := geoLoopAreaRads2(poly.GeoLoop)
	// NEVER(err) in C
	if err != eSuccess {
		return 0, err
	}
	kadd(&a, term)

	for i := 0; i < len(poly.Holes); i++ {
		term, err = geoLoopAreaRads2(poly.Holes[i])
		// NEVER(err) in C
		if err != eSuccess {
			return 0, err
		}

		// Due to clockwise order, holes will contribute area
		// of "everything except the hole", so adjust with -4*pi term.
		kadd(&a, term)
		kadd(&a, -4.0*math.Pi)
	}

	return a.sum, eSuccess
}
