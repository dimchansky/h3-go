//go:build cgo

package c2go

/*
#include <stdint.h>
#include "h3api.h"
#include "constants.h"
#include "h3Index.h"
#include "iterators.h"

// Wrapper for static _incrementResDigit function
// Based on iterators.c::_incrementResDigit implementation
static void incrementResDigitWrapper(H3Index *h, int res) {
    H3Index val = 1;
    val <<= H3_PER_DIGIT_OFFSET * (MAX_H3_RES - res);
    *h += val;
}

// Wrapper for static _null_iter function
// Based on iterators.c::_null_iter implementation
static IterCellsChildren nullIterWrapper(void) {
    IterCellsChildren iter;
    iter.h = H3_NULL;
    iter._parentRes = -1;
    iter._skipDigit = -1;
    return iter;
}
*/
import "C"

// IterCellsChildren mirrors the C IterCellsChildren struct for iterator state.
type IterCellsChildren struct {
	H         H3Index // Current H3 index
	ParentRes int     // Parent resolution
	SkipDigit int     // Skip digit for pentagons
}

// incrementResDigitC calls the C wrapper for _incrementResDigit.
func incrementResDigitC(h *H3Index, res int) {
	ch := C.H3Index(*h)
	C.incrementResDigitWrapper(&ch, C.int(res))
	*h = H3Index(ch)
}

// nullIterC calls the C wrapper for _null_iter.
func nullIterC() IterCellsChildren {
	citer := C.nullIterWrapper()
	return IterCellsChildren{
		H:         H3Index(citer.h),
		ParentRes: int(citer._parentRes),
		SkipDigit: int(citer._skipDigit),
	}
}
