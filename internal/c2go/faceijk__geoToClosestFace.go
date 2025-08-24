package c2go

// _geoToClosestFace determines the icosahedral face and squared euclidean distance
// from a lat/lng point to its closest icosahedral face center.
// Mirrors static _geoToClosestFace in faceijk.c
// Ported from H3 C: faceijk.c::_geoToClosestFace
func _geoToClosestFace(g *LatLng, face *int32, sqd *float64) {
	var v3d Vec3d
	_geoToVec3d(g, &v3d)

	// determine the icosahedron face
	*face = 0
	// The distance between two farthest points is 2.0, therefore the square of
	// the distance between two points should always be less or equal than 4.0 .
	*sqd = 5.0
	for f := int32(0); f < NUM_ICOSA_FACES; f++ {
		sqdT := _pointSquareDist(&faceCenterPoint[f], &v3d)
		if sqdT < *sqd {
			*face = f
			*sqd = sqdT
		}
	}
}
