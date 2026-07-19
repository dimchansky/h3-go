//go:build cgo && c2go

package h3

// Bridges to the H3 4.5.0 faceijk Vec3 projection pipeline:
// _vec3ToFaceIjk/_faceIjkToVec3 are extern (faceijk.h); the file-static
// helpers are reached through the h3goTest_* wrappers defined in the
// same translation unit (h3lib_faceijk_c2go.c).

/*
#include "h3api.h"
#include "faceijk.h"
#include "vec3d.h"

void h3goTest_vec3ToHex2d(const Vec3d *p, int res, int *face, Vec2d *v);
void h3goTest_vec3ToClosestFace(const Vec3d *v, int *face, double *sqd);
void h3goTest_hex2dToVec3(const Vec2d *v, int face, int res, int substrate, Vec3d *v3);
double h3goTest_vec3AzimuthRads(Vec3d p1, Vec3d p2);
void h3goTest_vec3TangentBasis(Vec3d p, Vec3d *north, Vec3d *east);
*/
import "C"

func fk3ToC(v vec3d) C.Vec3d {
	var cv C.Vec3d
	cv.x = C.double(v.X)
	cv.y = C.double(v.Y)
	cv.z = C.double(v.Z)
	return cv
}

func fk3FromC(cv C.Vec3d) vec3d {
	return vec3d{X: float64(cv.x), Y: float64(cv.y), Z: float64(cv.z)}
}

// _vec3ToFaceIjkC calls the original C implementation.
func _vec3ToFaceIjkC(p vec3d, res int32) faceIJK {
	var ch C.FaceIJK
	C._vec3ToFaceIjk(fk3ToC(p), C.int(res), &ch)
	return faceIJK{Face: int32(ch.face), Coord: coordIJK{I: int32(ch.coord.i), J: int32(ch.coord.j), K: int32(ch.coord.k)}}
}

// _faceIjkToVec3C calls the original C implementation.
func _faceIjkToVec3C(h faceIJK, res int32) vec3d {
	var ch C.FaceIJK
	ch.face = C.int(h.Face)
	ch.coord.i = C.int(h.Coord.I)
	ch.coord.j = C.int(h.Coord.J)
	ch.coord.k = C.int(h.Coord.K)
	var cv C.Vec3d
	C._faceIjkToVec3(&ch, C.int(res), &cv)
	return fk3FromC(cv)
}

// _vec3ToHex2dC calls the original C implementation (same-TU wrapper).
func _vec3ToHex2dC(p vec3d, res int32) (int32, vec2d) {
	cp := fk3ToC(p)
	var cface C.int
	var cv C.Vec2d
	C.h3goTest_vec3ToHex2d(&cp, C.int(res), &cface, &cv)
	return int32(cface), vec2d{X: float64(cv.x), Y: float64(cv.y)}
}

// _vec3ToClosestFaceC calls the original C implementation (same-TU wrapper).
func _vec3ToClosestFaceC(v vec3d) (int32, float64) {
	cv := fk3ToC(v)
	var cface C.int
	var csqd C.double
	C.h3goTest_vec3ToClosestFace(&cv, &cface, &csqd)
	return int32(cface), float64(csqd)
}

// _hex2dToVec3C calls the original C implementation (same-TU wrapper).
func _hex2dToVec3C(v vec2d, face, res, substrate int32) vec3d {
	var cv C.Vec2d
	cv.x = C.double(v.X)
	cv.y = C.double(v.Y)
	var cv3 C.Vec3d
	C.h3goTest_hex2dToVec3(&cv, C.int(face), C.int(res), C.int(substrate), &cv3)
	return fk3FromC(cv3)
}

// _vec3AzimuthRadsC calls the original C implementation (same-TU wrapper).
func _vec3AzimuthRadsC(p1, p2 vec3d) float64 {
	return float64(C.h3goTest_vec3AzimuthRads(fk3ToC(p1), fk3ToC(p2)))
}

// _vec3TangentBasisC calls the original C implementation (same-TU wrapper).
func _vec3TangentBasisC(p vec3d) (north, east vec3d) {
	var cn, ce C.Vec3d
	C.h3goTest_vec3TangentBasis(fk3ToC(p), &cn, &ce)
	return fk3FromC(cn), fk3FromC(ce)
}
