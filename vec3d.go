package h3

//lint:file-ignore U1000 Ignore all unused code
// 3D floating point vector functions.

import (
	"math"
)

// Vec3d holds 3D floating point structure
type Vec3d struct {
	x float64 // x component
	y float64 // y component
	z float64 // z component
}

// _pointSquareDist calculates the square of the distance between two 3D coordinates.
// `v1` is the first 3D coordinate.
// `v2` is the second 3D coordinate.
// Returns the square of the distance between the given points.
func _pointSquareDist(v1 *Vec3d, v2 *Vec3d) float64 {
	return square(v1.x-v2.x) + square(v1.y-v2.y) + square(v1.z-v2.z)
}

// _geoToVec3d calculates the 3D coordinate on unit sphere from the latitude and longitude.
//
// `geo` is the latitude and longitude of the point.
// `v` is the 3D coordinate of the point.
func _geoToVec3d(geo GeoCoord, v *Vec3d) {
	r := math.Cos(geo.lat)

	v.z = math.Sin(geo.lat)
	v.x = math.Cos(geo.lon) * r
	v.y = math.Sin(geo.lon) * r
}
