package h3

// vec3Normalize normalizes a 3D vector in place.
// Ported from H3 C: vec3d.h::vec3Normalize.
func vec3Normalize(v *vec3d) {
	norm := vec3Norm(*v)

	// Norm can be zero either from true zero vector, or from squaring
	// underflowing to zero.
	// If the norm is nonzero, we normalize v using it.
	// If the norm is zero, we set the vector to be exactly zero.
	s := 0.0
	if norm > 0.0 {
		s = 1.0 / norm
	}

	v.X *= s
	v.Y *= s
	v.Z *= s
}
