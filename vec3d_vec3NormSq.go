package h3

// vec3NormSq returns the squared norm of a 3D vector.
// Ported from H3 C: vec3d.h::vec3NormSq.
func vec3NormSq(v vec3d) float64 { return vec3Dot(v, v) }
