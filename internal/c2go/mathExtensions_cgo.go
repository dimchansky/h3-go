//go:build cgo && c2go

package c2go

/*
// Interop with H3 C source: mathExtensions.c
// NOTE: Include directories must be provided via CGO_CPPFLAGS (see Makefile test-c2go).
#include <stdint.h>
#include "mathExtensions.c"

// Wrapper with distinct name to avoid symbol collision.
static int64_t _ipow_c_wrapper(int64_t base, int64_t exp) {
    return _ipow(base, exp);
}
*/
import "C"

// _ipowC invokes the original C implementation for parity tests.
// Add more wrappers from mathExtensions.c here as needed.
func _ipowC(base, exp int64) int64 {
    return int64(C._ipow_c_wrapper(C.longlong(base), C.longlong(exp)))
}
