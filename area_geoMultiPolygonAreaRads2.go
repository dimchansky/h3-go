package h3

// geoMultiPolygonAreaRads2 computes the area of a GeoMultiPolygon in
// radians^2: the sum of the areas of its polygons. No check is made that
// the polygons are disjoint.
// Ported from H3 C: area.c::geoMultiPolygonAreaRads2.
//
//nolint:unused // exercised by the cgo && c2go parity tests; consumed by the I-C multipolygon port (#34)
func geoMultiPolygonAreaRads2(mpoly geoMultiPolygon) (float64, h3Error) {
	var a adder

	for i := int32(0); i < mpoly.NumPolygons; i++ {
		term, err := geoPolygonAreaRads2(mpoly.Polygons[i])
		// NEVER(err) in C
		if err != eSuccess {
			return 0, err
		}
		kadd(&a, term)
	}

	return a.sum, eSuccess
}
