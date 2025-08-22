package c2go

import "math"

// _geoToVec3d converts geographic coordinates to a 3D unit vector.
// Mirrors H3's vec3d.c::_geoToVec3d implementation (input as pointer).
func _geoToVec3d(geo *LatLng) Vec3d {
	// x = cos(lat) * cos(lng)
	// y = cos(lat) * sin(lng)
	// z = sin(lat)
	clat := math.Cos(geo.Lat)
	return Vec3d{
		X: clat * math.Cos(geo.Lng),
		Y: clat * math.Sin(geo.Lng),
		Z: math.Sin(geo.Lat),
	}
}
