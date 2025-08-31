//go:build cgo

package h3

/*
#cgo CFLAGS: -I${SRCDIR}/testref/h3-4.3.0/src/h3lib/include -I${SRCDIR}/testref/h3-4.3.0/src/apps/applib/include
#cgo CFLAGS: -std=c99

#include <stdint.h>
#include <stdlib.h>
#include "h3api.h"
#include "utility.h"

// Callback type for H3Index iteration
typedef void (*h3_callback)(H3Index);

// Global variables to store Go callback results
static H3Index* callback_buffer = NULL;
static int callback_index = 0;
static int callback_capacity = 0;

// C callback function that stores H3Index values
static void store_h3index(H3Index h) {
    if (callback_buffer != NULL && callback_index < callback_capacity) {
        callback_buffer[callback_index++] = h;
    }
}

// Wrapper for iterateAllIndexesAtResPartial
static int iterateAllIndexesAtResPartial_wrapper(int res, int baseCells, H3Index* out, int outSize) {
    callback_buffer = out;
    callback_index = 0;
    callback_capacity = outSize;
    iterateAllIndexesAtResPartial(res, store_h3index, baseCells);
    int result = callback_index;
    callback_buffer = NULL;
    callback_index = 0;
    callback_capacity = 0;
    return result;
}

// Wrapper for iterateAllIndexesAtRes
static int iterateAllIndexesAtRes_wrapper(int res, H3Index* out, int outSize) {
    callback_buffer = out;
    callback_index = 0;
    callback_capacity = outSize;
    iterateAllIndexesAtRes(res, store_h3index);
    int result = callback_index;
    callback_buffer = NULL;
    callback_index = 0;
    callback_capacity = 0;
    return result;
}

// Wrapper for iterateBaseCellIndexesAtRes
static int iterateBaseCellIndexesAtRes_wrapper(int res, int baseCell, H3Index* out, int outSize) {
    callback_buffer = out;
    callback_index = 0;
    callback_capacity = outSize;
    iterateBaseCellIndexesAtRes(res, store_h3index, baseCell);
    int result = callback_index;
    callback_buffer = NULL;
    callback_index = 0;
    callback_capacity = 0;
    return result;
}
*/
import "C"
import (
	"unsafe"
)

// iterateAllIndexesAtResPartialC calls the C implementation of iterateAllIndexesAtResPartial
func iterateAllIndexesAtResPartialC(res int32, baseCells int32, out []H3Index) int {
	if len(out) == 0 {
		return 0
	}
	count := C.iterateAllIndexesAtResPartial_wrapper(
		C.int(res),
		C.int(baseCells),
		(*C.H3Index)(unsafe.Pointer(&out[0])),
		C.int(len(out)),
	)
	return int(count)
}

// iterateAllIndexesAtResC calls the C implementation of iterateAllIndexesAtRes
func iterateAllIndexesAtResC(res int32, out []H3Index) int {
	if len(out) == 0 {
		return 0
	}
	count := C.iterateAllIndexesAtRes_wrapper(
		C.int(res),
		(*C.H3Index)(unsafe.Pointer(&out[0])),
		C.int(len(out)),
	)
	return int(count)
}

// iterateBaseCellIndexesAtResC calls the C implementation of iterateBaseCellIndexesAtRes
func iterateBaseCellIndexesAtResC(res int32, baseCell int32, out []H3Index) int {
	if len(out) == 0 {
		return 0
	}
	count := C.iterateBaseCellIndexesAtRes_wrapper(
		C.int(res),
		C.int(baseCell),
		(*C.H3Index)(unsafe.Pointer(&out[0])),
		C.int(len(out)),
	)
	return int(count)
}
