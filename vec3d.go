package h3

//lint:file-ignore U1000 Ignore all unused code
// 3D floating point vector functions.

import (
	"math"
)

// Vec3d holds 3D floating point structure
type Vec3d struct {
	X float64 // X component
	Y float64 // Y component
	Z float64 // Z component
}

// pointSquareDist calculates the square of the distance between two 3D coordinates.
// `v1` is the first 3D coordinate.
// `v2` is the second 3D coordinate.
// Returns the square of the distance between the given points.
func pointSquareDist(v1 *Vec3d, v2 *Vec3d) float64 {
	return square(v1.X-v2.X) + square(v1.Y-v2.Y) + square(v1.Z-v2.Z)
}

// geoToVec3d calculates the 3D coordinate on unit sphere from the latitude and longitude.
//
// `geo` is the latitude and longitude of the point.
// `v` is the 3D coordinate of the point.
func geoToVec3d(geo GeoCoord, v *Vec3d) {
	r := math.Cos(geo.Lat)

	v.Z = math.Sin(geo.Lat)
	v.X = math.Cos(geo.Lon) * r
	v.Y = math.Sin(geo.Lon) * r
}
