//go:build cgo && c2go

package c2go

/*
// Interop wrapper only; the original C sources are compiled via separate
// build-tagged C shim files in this package (see h3lib_*.c with //go:build c2go).
#include <stdint.h>
#include <stdbool.h>
#include "h3api.h"
#include "latLng.h"
// Normalize C bool to int for cgo comparisons when needed (toolchain-safe)
static int h3_bool_to_int(_Bool b) { return b ? 1 : 0; }

// Forward declarations for wrappers
double _posAngleRads(double);
double constrainLng(double);
double constrainLat(double);
double degsToRads(double);
double radsToDegs(double);

static double _posAngleRads_c_wrapper(double rads) { return _posAngleRads(rads); }
static double constrainLng_c_wrapper(double lng) { return constrainLng(lng); }
static double constrainLat_c_wrapper(double lat) { return constrainLat(lat); }
static double degsToRads_c_wrapper(double degrees) { return degsToRads(degrees); }
static double radsToDegs_c_wrapper(double radians) { return radsToDegs(radians); }
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

// degsToRadsC invokes the original C implementation for parity tests.
func degsToRadsC(degrees float64) float64 {
    return float64(C.degsToRads_c_wrapper(C.double(degrees)))
}

// radsToDegsC invokes the original C implementation for parity tests.
func radsToDegsC(radians float64) float64 {
    return float64(C.radsToDegs_c_wrapper(C.double(radians)))
}

// geoAlmostEqualThresholdC calls the original C function using plain doubles.
func geoAlmostEqualThresholdC(a, b LatLng, threshold float64) bool {
    var ca, cb C.LatLng
    ca.lat = C.double(a.lat)
    ca.lng = C.double(a.lng)
    cb.lat = C.double(b.lat)
    cb.lng = C.double(b.lng)
    // Convert C.bool to int via C helper for toolchain compatibility
    return C.h3_bool_to_int(C.geoAlmostEqualThreshold(&ca, &cb, C.double(threshold))) != 0
}

// geoAlmostEqualC calls the original C function using plain doubles.
func geoAlmostEqualC(a, b LatLng) bool {
    var ca, cb C.LatLng
    ca.lat = C.double(a.lat)
    ca.lng = C.double(a.lng)
    cb.lat = C.double(b.lat)
    cb.lng = C.double(b.lng)
    return C.h3_bool_to_int(C.geoAlmostEqual(&ca, &cb)) != 0
}
