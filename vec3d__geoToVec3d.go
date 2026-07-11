package h3

// _geoToVec3d converts geographic coordinates to a 3D unit vector.
// Ported from H3 C: vec3d.c::_geoToVec3d.
func _geoToVec3d(geo *LatLng, v *Vec3d) {
	sinLat, cosLat := geo.Lat.SinCos()
	sinLng, cosLng := geo.Lng.SinCos()

	v.Z = sinLat
	v.X = cosLng * cosLat
	v.Y = sinLng * cosLat
}
