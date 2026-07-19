package h3

// _vec3ToClosestFace encodes a coordinate on the sphere to the
// corresponding icosahedral face and containing the squared euclidean
// distance to that face center.
//
// Vec3d v is expected to be on the unit sphere.
// Ported from H3 C: faceijk.c::_vec3ToClosestFace.
func _vec3ToClosestFace(v *vec3d, face *int32, sqd *float64) {
	*face = 0
	// The distance between two farthest points is 2.0, therefore the square of
	// the distance between two points should always be less or equal than 4.0 .
	*sqd = 5.0
	for f := int32(0); f < numIcosaFaces; f++ {
		sqdT := vec3DistSq(faceCenterPoint[f], *v)
		if sqdT < *sqd {
			*face = f
			*sqd = sqdT
		}
	}
}
