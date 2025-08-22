//go:build cgo && c2go

package c2go

/*
#include <stdint.h>
#include <stdbool.h>
#include "coordijk.h"

// Prototypes for the original C helpers in coordijk.c
void _ijkAdd(const CoordIJK* h1, const CoordIJK* h2, CoordIJK* sum);
void _ijkSub(const CoordIJK* h1, const CoordIJK* h2, CoordIJK* diff);
*/
import "C"

// _ijkAddC calls the original C implementation for adding IJK coordinates.
// Bridges to coordijk.c::_ijkAdd.
func _ijkAddC(h1, h2 *CoordIJK) CoordIJK {
	var ch1, ch2, sum C.CoordIJK
	ch1.i = C.int(h1.I)
	ch1.j = C.int(h1.J)
	ch1.k = C.int(h1.K)
	ch2.i = C.int(h2.I)
	ch2.j = C.int(h2.J)
	ch2.k = C.int(h2.K)
	C._ijkAdd(&ch1, &ch2, &sum)
	return CoordIJK{I: int(sum.i), J: int(sum.j), K: int(sum.k)}
}

// _ijkSubC calls the original C implementation for subtracting IJK coordinates.
// Bridges to coordijk.c::_ijkSub.
func _ijkSubC(h1, h2 *CoordIJK) CoordIJK {
	var ch1, ch2, diff C.CoordIJK
	ch1.i = C.int(h1.I)
	ch1.j = C.int(h1.J)
	ch1.k = C.int(h1.K)
	ch2.i = C.int(h2.I)
	ch2.j = C.int(h2.J)
	ch2.k = C.int(h2.K)
	C._ijkSub(&ch1, &ch2, &diff)
	return CoordIJK{I: int(diff.i), J: int(diff.j), K: int(diff.k)}
}
