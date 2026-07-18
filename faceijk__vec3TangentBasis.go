package h3

// _vec3TangentBasis computes the local north and east directions on the
// tangent plane at a point on the unit sphere.
//
// Will not work if p is at a pole, but icosahedron face centers
// are never at the poles.
// Ported from H3 C: faceijk.c::_vec3TangentBasis.
func _vec3TangentBasis(p vec3d, north, east *vec3d) {
	northPole := vec3d{0.0, 0.0, 1.0}
	*north = vec3LinComb(1.0, northPole, -vec3Dot(northPole, p), p)
	vec3Normalize(north)
	*east = vec3Cross(*north, p)
}
