package c2go

// _pointSquareDist returns the squared distance between two 3D vectors.
// Mirrors H3's vec3d.c::_pointSquareDist implementation.
func _pointSquareDist(v1, v2 Vec3d) float64 {
	dx := v1.X - v2.X
	dy := v1.Y - v2.Y
	dz := v1.Z - v2.Z
	return dx*dx + dy*dy + dz*dz
}
