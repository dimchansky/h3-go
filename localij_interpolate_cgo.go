//go:build cgo && c2go && h3v450

package h3

// Bridge to the file-static gridPathCellsInterpolate added to localij.c
// in H3 4.5.0, via the same-TU wrapper in h3lib_localij_c2go.c.

/*
#include <stdint.h>
#include "h3api.h"

H3Error h3goTest_gridPathCellsInterpolate(H3Index start, H3Index end,
                                          int64_t distance, H3Index *out,
                                          int64_t outOffset, int64_t outStep);
*/
import "C"

// gridPathCellsInterpolateC calls the original C implementation.
func gridPathCellsInterpolateC(start, end h3Index, distance int64, out []h3Index, outOffset, outStep int64) h3Error {
	if len(out) == 0 {
		return eFailed
	}
	return h3Error(C.h3goTest_gridPathCellsInterpolate(
		C.H3Index(start), C.H3Index(end), C.int64_t(distance),
		(*C.H3Index)(&out[0]), C.int64_t(outOffset), C.int64_t(outStep)))
}
