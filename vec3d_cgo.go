//go:build cgo && c2go

package h3

/*
#include <stdint.h>
#include <stdbool.h>
#include "vec3d.h"
// Prototype for the original C helper in vec3d.c
double _pointSquareDist(const Vec3d* v1, const Vec3d* v2);
double _square(double x);
*/
import "C"

// pointSquareDistC calls the original C implementation for squared distance.
// This bridges to vec3d.c::_pointSquareDist.

func pointSquareDistC(a, b vec3d) float64 {
	var ca, cb C.Vec3d
	ca.x = C.double(a.X)
	ca.y = C.double(a.Y)
	ca.z = C.double(a.Z)
	cb.x = C.double(b.X)
	cb.y = C.double(b.Y)
	cb.z = C.double(b.Z)
	return float64(C._pointSquareDist(&ca, &cb))
}

// _geoToVec3dC calls the original C implementation to convert a LatLng to vec3d.
// Bridges to vec3d.c::_geoToVec3d.
func _geoToVec3dC(geo *LatLng, v *vec3d) {
	var cg C.LatLng
	cg.lat = C.double(geo.Lat)
	cg.lng = C.double(geo.Lng)
	var cv C.Vec3d
	C._geoToVec3d(&cg, &cv)
	v.X = float64(cv.x)
	v.Y = float64(cv.y)
	v.Z = float64(cv.z)
}

// _squareC calls the original C implementation.
func _squareC(x float64) float64 {
	return float64(C._square(C.double(x)))
}
