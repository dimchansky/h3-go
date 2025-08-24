//go:build cgo && c2go

package c2go

/*
#include <stdbool.h>
#include "h3Index.h"
#include "coordijk.h"
#include "localij.h"

// Wrapper functions to call the original C functions directly
static H3Error cellToLocalIjk_c_wrapper(H3Index origin, H3Index h3, CoordIJK* out) {
    return cellToLocalIjk(origin, h3, out);
}

static H3Error localIjkToCell_c_wrapper(H3Index origin, const CoordIJK* ijk, H3Index* out) {
    return localIjkToCell(origin, ijk, out);
}
*/
import "C"

// _cellToLocalIjkC wraps the C cellToLocalIjk function.
// Provides direct access to the C implementation for parity testing.
func _cellToLocalIjkC(origin H3Index, h3 H3Index, out *CoordIJK) H3Error {
	var cOut C.CoordIJK
	err := C.cellToLocalIjk_c_wrapper(C.H3Index(origin), C.H3Index(h3), &cOut)

	// Convert C result to Go
	out.I = int32(cOut.i)
	out.J = int32(cOut.j)
	out.K = int32(cOut.k)

	return H3Error(err)
}

// _localIjkToCellC wraps the C localIjkToCell function.
// Provides direct access to the C implementation for parity testing.
func _localIjkToCellC(origin H3Index, ijk *CoordIJK, out *H3Index) H3Error {
	var cIJK C.CoordIJK
	var cOut C.H3Index

	// Convert Go to C
	cIJK.i = C.int(ijk.I)
	cIJK.j = C.int(ijk.J)
	cIJK.k = C.int(ijk.K)

	err := C.localIjkToCell_c_wrapper(C.H3Index(origin), &cIJK, &cOut)

	// Convert C result to Go
	*out = H3Index(cOut)

	return H3Error(err)
}
