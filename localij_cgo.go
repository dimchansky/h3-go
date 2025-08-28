//go:build cgo && c2go

package h3

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

static H3Error cellToLocalIj_c_wrapper(H3Index origin, H3Index index, uint32_t mode, CoordIJ* out) {
    return cellToLocalIj(origin, index, mode, out);
}

static H3Error localIjToCell_c_wrapper(H3Index origin, const CoordIJ* ij, uint32_t mode, H3Index* out) {
    return localIjToCell(origin, ij, mode, out);
}

static H3Error gridDistance_c_wrapper(H3Index origin, H3Index index, int64_t* out) {
    return gridDistance(origin, index, out);
}

static H3Error gridPathCellsSize_c_wrapper(H3Index start, H3Index end, int64_t* size) {
    return gridPathCellsSize(start, end, size);
}

static H3Error gridPathCells_c_wrapper(H3Index start, H3Index end, H3Index* out) {
    return gridPathCells(start, end, out);
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

// _cellToLocalIjC wraps the C cellToLocalIj function.
// Provides direct access to the C implementation for parity testing.
func _cellToLocalIjC(origin H3Index, index H3Index, mode uint32, out *CoordIJ) H3Error {
	var cOut C.CoordIJ
	err := C.cellToLocalIj_c_wrapper(C.H3Index(origin), C.H3Index(index), C.uint32_t(mode), &cOut)

	// Convert C result to Go
	out.I = int32(cOut.i)
	out.J = int32(cOut.j)

	return H3Error(err)
}

// _localIjToCellC wraps the C localIjToCell function.
// Provides direct access to the C implementation for parity testing.
func _localIjToCellC(origin H3Index, ij *CoordIJ, mode uint32, out *H3Index) H3Error {
	var cIJ C.CoordIJ
	var cOut C.H3Index

	// Convert Go to C
	cIJ.i = C.int(ij.I)
	cIJ.j = C.int(ij.J)

	err := C.localIjToCell_c_wrapper(C.H3Index(origin), &cIJ, C.uint32_t(mode), &cOut)

	// Convert C result to Go
	*out = H3Index(cOut)

	return H3Error(err)
}

// _gridDistanceC wraps the C gridDistance function.
// Provides direct access to the C implementation for parity testing.
func _gridDistanceC(origin H3Index, index H3Index, out *int64) H3Error {
	var cOut C.int64_t
	err := C.gridDistance_c_wrapper(C.H3Index(origin), C.H3Index(index), &cOut)

	// Convert C result to Go
	*out = int64(cOut)

	return H3Error(err)
}

// _gridPathCellsSizeC wraps the C gridPathCellsSize function.
// Provides direct access to the C implementation for parity testing.
func _gridPathCellsSizeC(start H3Index, end H3Index, size *int64) H3Error {
	var cSize C.int64_t
	err := C.gridPathCellsSize_c_wrapper(C.H3Index(start), C.H3Index(end), &cSize)

	// Convert C result to Go
	*size = int64(cSize)

	return H3Error(err)
}

// _gridPathCellsC wraps the C gridPathCells function.
// Provides direct access to the C implementation for parity testing.
func _gridPathCellsC(start H3Index, end H3Index, out []H3Index) H3Error {
	if len(out) == 0 {
		return E_SUCCESS
	}

	// Allocate C array
	cOut := (*C.H3Index)(&out[0])

	err := C.gridPathCells_c_wrapper(C.H3Index(start), C.H3Index(end), cOut)

	return H3Error(err)
}
