package h3

// latLngToVec3 converts latitude and longitude to a unit vec3d on the
// sphere.
// Ported from H3 C: vec3d.h::latLngToVec3.
func latLngToVec3(geo LatLng) vec3d {
	r := geo.Lat.Cos()
	out := vec3d{
		X: geo.Lng.Cos() * r,
		Y: geo.Lng.Sin() * r,
		Z: geo.Lat.Sin(),
	}
	return out
}
