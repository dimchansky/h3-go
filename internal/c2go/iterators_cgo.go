//go:build cgo

package c2go

/*
#include <stdint.h>
#include "h3api.h"
#include "constants.h"
#include "h3Index.h"

// Wrapper for static _incrementResDigit function
// Based on iterators.c::_incrementResDigit implementation
static void incrementResDigitWrapper(H3Index *h, int res) {
    H3Index val = 1;
    val <<= H3_PER_DIGIT_OFFSET * (MAX_H3_RES - res);
    *h += val;
}
*/
import "C"

// incrementResDigitC calls the C wrapper for _incrementResDigit.
func incrementResDigitC(h *H3Index, res int) {
	ch := C.H3Index(*h)
	C.incrementResDigitWrapper(&ch, C.int(res))
	*h = H3Index(ch)
}
