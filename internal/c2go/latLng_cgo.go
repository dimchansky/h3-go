//go:build cgo && c2go

package c2go

/*
// Interop wrapper only; the original C sources are compiled via separate
// build-tagged C shim files in this package (see h3lib_*.c with //go:build c2go).
#include <stdint.h>

// Forward declarations for wrappers
double _posAngleRads(double);
double constrainLng(double);
double constrainLat(double);

static double _posAngleRads_c_wrapper(double rads) { return _posAngleRads(rads); }
static double constrainLng_c_wrapper(double lng) { return constrainLng(lng); }
static double constrainLat_c_wrapper(double lat) { return constrainLat(lat); }
*/
import "C"

// _posAngleRadsC invokes the original C implementation for parity tests.
func _posAngleRadsC(rads float64) float64 {
    return float64(C._posAngleRads_c_wrapper(C.double(rads)))
}

// _constrainLngC invokes the original C implementation for parity tests.
func _constrainLngC(lng float64) float64 {
    return float64(C.constrainLng_c_wrapper(C.double(lng)))
}

// _constrainLatC invokes the original C implementation for parity tests.
func _constrainLatC(lat float64) float64 {
    return float64(C.constrainLat_c_wrapper(C.double(lat)))
}
