package c2go

import "math"

// _geoToVec3d converts geographic coordinates to a 3D unit vector.
// Ported from H3 C: vec3d.c::_geoToVec3d
func _geoToVec3d(geo *LatLng, v *Vec3d) {
	r := math.Cos(geo.Lat)

	v.Z = math.Sin(geo.Lat)
	v.X = math.Cos(geo.Lng) * r
	v.Y = math.Sin(geo.Lng) * r
}
