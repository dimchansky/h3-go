package h3

// vec3DistSq returns the squared distance between two 3D vectors.
// Ported from H3 C: vec3d.h::vec3DistSq.
func vec3DistSq(v1, v2 vec3d) float64 {
	d := vec3LinComb(1.0, v1, -1.0, v2)
	return vec3NormSq(d)
}
