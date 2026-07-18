//go:build cgo && c2go && h3v450

package h3

// Bridges to the H3 4.5.0 header-only vec3d implementation: vec3d.h
// defines these as static inline, so the cgo-generated translation unit
// including the header gets its own definitions — no C shim needed
// (docs/sync/4.4.0-to-4.5.0.md §15.1).

/*
#include "h3api.h"
#include "vec3d.h"
*/
import "C"

func cvtVec3ToC(v vec3d) C.Vec3d {
	var cv C.Vec3d
	cv.x = C.double(v.X)
	cv.y = C.double(v.Y)
	cv.z = C.double(v.Z)
	return cv
}

func cvtVec3FromC(cv C.Vec3d) vec3d {
	return vec3d{X: float64(cv.x), Y: float64(cv.y), Z: float64(cv.z)}
}

// latLngToVec3C calls the original C implementation.
func latLngToVec3C(geo LatLng) vec3d {
	var cg C.LatLng
	cg.lat = C.double(geo.Lat.Rad())
	cg.lng = C.double(geo.Lng.Rad())
	return cvtVec3FromC(C.latLngToVec3(cg))
}

// vec3ToLatLngC calls the original C implementation.
func vec3ToLatLngC(v vec3d) LatLng {
	cg := C.vec3ToLatLng(cvtVec3ToC(v))
	return LatLng{Lat: Rad(float64(cg.lat)), Lng: Rad(float64(cg.lng))}
}

// vec3LinCombC calls the original C implementation.
func vec3LinCombC(a float64, v1 vec3d, b float64, v2 vec3d) vec3d {
	return cvtVec3FromC(C.vec3LinComb(C.double(a), cvtVec3ToC(v1), C.double(b), cvtVec3ToC(v2)))
}

// vec3CrossC calls the original C implementation.
func vec3CrossC(v1, v2 vec3d) vec3d {
	return cvtVec3FromC(C.vec3Cross(cvtVec3ToC(v1), cvtVec3ToC(v2)))
}

// vec3DotC calls the original C implementation.
func vec3DotC(v1, v2 vec3d) float64 {
	return float64(C.vec3Dot(cvtVec3ToC(v1), cvtVec3ToC(v2)))
}

// vec3NormSqC calls the original C implementation.
func vec3NormSqC(v vec3d) float64 { return float64(C.vec3NormSq(cvtVec3ToC(v))) }

// vec3NormC calls the original C implementation.
func vec3NormC(v vec3d) float64 { return float64(C.vec3Norm(cvtVec3ToC(v))) }

// vec3NormalizeC calls the original C implementation.
func vec3NormalizeC(v *vec3d) {
	cv := cvtVec3ToC(*v)
	C.vec3Normalize(&cv)
	*v = cvtVec3FromC(cv)
}

// vec3DistSqC calls the original C implementation.
func vec3DistSqC(v1, v2 vec3d) float64 {
	return float64(C.vec3DistSq(cvtVec3ToC(v1), cvtVec3ToC(v2)))
}
