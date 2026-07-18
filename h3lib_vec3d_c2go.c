//go:build cgo && c2go && !h3v450

// vec3d.c exists only in the H3 4.4.0 tree (4.5.0 moved vec3d into a
// header-only implementation). The 4.4.0 C library itself still needs
// this translation unit — faceijk.c's _geoToClosestFace links against
// _geoToVec3d/_pointSquareDist — so the shim stays until the I-I cutover
// deletes the 4.4.0 configuration, even though the Go-side vec3d
// wrappers retired with I-A (#29).
#include "vec3d.c"
