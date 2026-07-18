//go:build cgo && c2go

#include "faceijk.c"

// H3 4.5.0 tree only (area.c is its marker file): test-only wrappers
// exposing the file-static Vec3 pipeline helpers to the parity harness.
// They live in the same translation unit as the statics they call
// (docs/sync/4.4.0-to-4.5.0.md parity-first coverage requirement).
#if __has_include("area.c")
void h3goTest_vec3ToHex2d(const Vec3d *p, int res, int *face, Vec2d *v) {
    _vec3ToHex2d(p, res, face, v);
}
void h3goTest_vec3ToClosestFace(const Vec3d *v, int *face, double *sqd) {
    _vec3ToClosestFace(v, face, sqd);
}
void h3goTest_hex2dToVec3(const Vec2d *v, int face, int res, int substrate,
                          Vec3d *v3) {
    _hex2dToVec3(v, face, res, substrate, v3);
}
double h3goTest_vec3AzimuthRads(Vec3d p1, Vec3d p2) {
    return _vec3AzimuthRads(p1, p2);
}
void h3goTest_vec3TangentBasis(Vec3d p, Vec3d *north, Vec3d *east) {
    _vec3TangentBasis(p, north, east);
}
#endif
