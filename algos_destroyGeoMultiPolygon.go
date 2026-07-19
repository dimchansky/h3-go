package h3

// destroyGeoMultiPolygon frees all allocated memory for a
// GeoMultiPolygon. The caller is responsible for freeing memory
// allocated to the input GeoMultiPolygon struct. In Go the garbage
// collector frees the polygons; zeroing mirrors C's pointer-nulling.
// The C function is exported for the bindings; it stays internal here
// like the rest of the multipolygon machinery (record §13.2).
// Ported from H3 C: algos.c::destroyGeoMultiPolygon.
func destroyGeoMultiPolygon(mpoly *geoMultiPolygon) {
	for i := int32(0); i < mpoly.NumPolygons; i++ {
		destroyGeoPolygon(&mpoly.Polygons[i])
	}
	mpoly.Polygons = nil
	mpoly.NumPolygons = 0
}
