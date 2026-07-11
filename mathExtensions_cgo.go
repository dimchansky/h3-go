//go:build cgo && c2go

package h3

/*
// Interop with H3 C source: mathExtensions.c
// NOTE: Include directories must be provided via CGO_CPPFLAGS (see Makefile test-c2go).
#include <stdint.h>
#include <stdbool.h>
#include "mathExtensions.h"
#include "mathExtensions.c"

// Wrapper with distinct name to avoid symbol collision.
static int64_t _ipow_c_wrapper(int64_t base, int64_t exp) {
    return _ipow(base, exp);
}

// Wrappers for overflow checking macros
static bool _add_int32s_overflows_c_wrapper(int32_t a, int32_t b) {
    return ADD_INT32S_OVERFLOWS(a, b);
}

static bool _sub_int32s_overflows_c_wrapper(int32_t a, int32_t b) {
    return SUB_INT32S_OVERFLOWS(a, b);
}
*/
import "C"

// _ipowC invokes the original C implementation for parity tests.
// Add more wrappers from mathExtensions.c here as needed.
func _ipowC(base, exp int64) int64 {
	// C.int64_t (not C.longlong): int64_t is `long` on linux/amd64 glibc but
	// `long long` on darwin, so the stdint spelling is the portable one.
	return int64(C._ipow_c_wrapper(C.int64_t(base), C.int64_t(exp)))
}

// addInt32sOverflowsC invokes the original C ADD_INT32S_OVERFLOWS macro for parity tests.
func addInt32sOverflowsC(a, b int32) bool {
	return bool(C._add_int32s_overflows_c_wrapper(C.int32_t(a), C.int32_t(b)))
}

// subInt32sOverflowsC invokes the original C SUB_INT32S_OVERFLOWS macro for parity tests.
func subInt32sOverflowsC(a, b int32) bool {
	return bool(C._sub_int32s_overflows_c_wrapper(C.int32_t(a), C.int32_t(b)))
}
