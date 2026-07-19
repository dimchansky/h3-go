package h3

// destroyGeoLoop frees all allocated memory for a GeoLoop. The caller
// is responsible for freeing memory allocated to the input GeoLoop
// struct. In Go the garbage collector frees the vertex array; nilling
// the slice mirrors C's verts = NULL / numVerts = 0 zeroing.
// Ported from H3 C: algos.c::destroyGeoLoop.
func destroyGeoLoop(loop *GeoLoop) {
	*loop = nil
}
