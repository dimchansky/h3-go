package h3

import "math"

// _vec3AzimuthRads calculates the azimuth from p1 to p2.
// Ported from H3 C: faceijk.c::_vec3AzimuthRads.
func _vec3AzimuthRads(p1, p2 vec3d) float64 {
	var northDir, eastDir vec3d
	_vec3TangentBasis(p1, &northDir, &eastDir)

	// project p2 onto tangent plane at p1
	p2Proj := vec3LinComb(1.0, p2, -vec3Dot(p2, p1), p1)
	vec3Normalize(&p2Proj)

	return math.Atan2(vec3Dot(p2Proj, eastDir), vec3Dot(p2Proj, northDir))
}
