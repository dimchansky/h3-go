//go:build cgo && c2go

package h3

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

// Forward declaration for _iterInitParent function
extern void _iterInitParent(H3Index h, int childRes, IterCellsChildren *iter);

// Forward declarations for iterator functions
extern IterCellsChildren iterInitBaseCellNum(int baseCellNum, int childRes);
extern IterCellsResolution iterInitRes(int res);
extern void iterStepRes(IterCellsResolution *iter);
*/
import "C"

// incrementResDigitC calls the C wrapper for _incrementResDigit.
func incrementResDigitC(h *H3Index, res int32) {
	ch := C.H3Index(*h)
	C.incrementResDigitWrapper(&ch, C.int(res))
	*h = H3Index(ch)
}

// nullIterC calls the C wrapper for _null_iter.
func nullIterC() IterCellsChildren {
	citer := C.nullIterWrapper()
	return IterCellsChildren{
		H:         H3Index(citer.h),
		ParentRes: int32(citer._parentRes),
		SkipDigit: int32(citer._skipDigit),
	}
}

// iterInitParentC calls the original C _iterInitParent function.
func iterInitParentC(h H3Index, childRes int32, iter *IterCellsChildren) {
	var citer C.IterCellsChildren
	C._iterInitParent(C.H3Index(h), C.int(childRes), &citer)
	iter.H = H3Index(citer.h)
	iter.ParentRes = int32(citer._parentRes)
	iter.SkipDigit = int32(citer._skipDigit)
}

// iterStepChildC calls the original C iterStepChild function.
func iterStepChildC(iter *IterCellsChildren) {
	var citer C.IterCellsChildren
	citer.h = C.H3Index(iter.H)
	citer._parentRes = C.int(iter.ParentRes)
	citer._skipDigit = C.int(iter.SkipDigit)
	C.iterStepChild(&citer)
	iter.H = H3Index(citer.h)
	iter.ParentRes = int32(citer._parentRes)
	iter.SkipDigit = int32(citer._skipDigit)
}

// iterInitBaseCellNumC calls the original C iterInitBaseCellNum function.
func iterInitBaseCellNumC(baseCellNum int32, childRes int32) IterCellsChildren {
	citer := C.iterInitBaseCellNum(C.int(baseCellNum), C.int(childRes))
	return IterCellsChildren{
		H:         H3Index(citer.h),
		ParentRes: int32(citer._parentRes),
		SkipDigit: int32(citer._skipDigit),
	}
}

// iterInitResC calls the original C iterInitRes function.
func iterInitResC(res int32) IterCellsResolution {
	citer := C.iterInitRes(C.int(res))
	return IterCellsResolution{
		H:           H3Index(citer.h),
		baseCellNum: int32(citer._baseCellNum),
		res:         int32(citer._res),
		itC: IterCellsChildren{
			H:         H3Index(citer._itC.h),
			ParentRes: int32(citer._itC._parentRes),
			SkipDigit: int32(citer._itC._skipDigit),
		},
	}
}

// iterStepResC calls the original C iterStepRes function.
func iterStepResC(iter *IterCellsResolution) {
	var citer C.IterCellsResolution
	citer.h = C.H3Index(iter.H)
	citer._baseCellNum = C.int(iter.baseCellNum)
	citer._res = C.int(iter.res)
	citer._itC.h = C.H3Index(iter.itC.H)
	citer._itC._parentRes = C.int(iter.itC.ParentRes)
	citer._itC._skipDigit = C.int(iter.itC.SkipDigit)

	C.iterStepRes(&citer)

	iter.H = H3Index(citer.h)
	iter.baseCellNum = int32(citer._baseCellNum)
	iter.res = int32(citer._res)
	iter.itC.H = H3Index(citer._itC.h)
	iter.itC.ParentRes = int32(citer._itC._parentRes)
	iter.itC.SkipDigit = int32(citer._itC._skipDigit)
}
