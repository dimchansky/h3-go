//go:build cgo && c2go

package c2go

/*
// Interop wrapper only; the original C sources are compiled via separate
// build-tagged C shim files in this package (see h3lib_*.c with //go:build c2go).
#include <stdint.h>

// Forward declaration for wrapper
double _posAngleRads(double);

static double _posAngleRads_c_wrapper(double rads) { return _posAngleRads(rads); }
*/
import "C"

// _posAngleRadsC invokes the original C implementation for parity tests.
func _posAngleRadsC(rads float64) float64 {
    return float64(C._posAngleRads_c_wrapper(C.double(rads)))
}
