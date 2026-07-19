package h3

// destroyGeoPolygon frees all allocated memory for a GeoPolygon. The
// caller is responsible for freeing memory allocated to the input
// GeoPolygon struct. In Go the garbage collector frees the loops;
// zeroing mirrors C's pointer-nulling.
// Ported from H3 C: algos.c::destroyGeoPolygon.
func destroyGeoPolygon(poly *GeoPolygon) {
	destroyGeoLoop(&poly.GeoLoop)
	for i := 0; i < len(poly.Holes); i++ {
		destroyGeoLoop(&poly.Holes[i])
	}
	poly.Holes = nil
}
