package h3

import "math"

// vec3ToLatLng converts a unit vec3d on the sphere back to latitude and
// longitude.
// Ported from H3 C: vec3d.h::vec3ToLatLng.
func vec3ToLatLng(v vec3d) LatLng {
	out := LatLng{
		Lat: Rad(math.Asin(v.Z)),
		Lng: Rad(math.Atan2(v.Y, v.X)),
	}
	return out
}
