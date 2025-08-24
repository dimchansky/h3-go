//go:build cgo && c2go

package c2go

/*
#cgo CFLAGS: -DH3_HAVE_VLA
#include <stdint.h>
#include "h3api.h"
#include "vertex.h"

// The vertexRotations function is static in vertex.c, so we need to include the implementation
#include "vertex.c"

// Wrapper to access the static function
static inline H3Error vertexRotations_wrapper(H3Index cell, int *out) {
    return vertexRotations(cell, out);
}
*/
import "C"

// vertexRotationsC wraps the C vertexRotations function.
func vertexRotationsC(cell H3Index, out *int32) H3Error {
	var cOut C.int
	err := H3Error(C.vertexRotations_wrapper(C.H3Index(cell), &cOut))
	*out = int32(cOut)
	return err
}
