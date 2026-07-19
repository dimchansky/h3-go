package h3

import "math"

// vec3Norm returns the norm of a 3D vector.
// Ported from H3 C: vec3d.h::vec3Norm.
func vec3Norm(v vec3d) float64 { return math.Sqrt(vec3NormSq(v)) }
