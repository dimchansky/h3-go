//go:build cgo && c2go

package c2go

/*
#include <stdint.h>
#include <stdbool.h>
#include <stdlib.h>
#include "h3api.h"
*/
import "C"

import "unsafe"

// getResolutionC calls the original C implementation.
func getResolutionC(h H3Index) int { return int(C.getResolution(C.H3Index(h))) }

// getBaseCellNumberC calls the original C implementation.
func getBaseCellNumberC(h H3Index) int { return int(C.getBaseCellNumber(C.H3Index(h))) }

// stringToH3C calls the original C implementation.
func stringToH3C(s string) (H3Index, uint32) {
    cs := C.CString(s)
    defer C.free(unsafe.Pointer(cs))
    var out C.H3Index
    err := C.stringToH3(cs, &out)
    return H3Index(out), uint32(err)
}

// h3ToStringC calls the original C implementation (17 bytes incl. NUL).
func h3ToStringC(h H3Index) (string, uint32) {
    const sz = 17
    buf := C.malloc(sz)
    defer C.free(buf)
    err := C.h3ToString(C.H3Index(h), (*C.char)(buf), C.size_t(sz))
    return C.GoString((*C.char)(buf)), uint32(err)
}

// describeH3ErrorC calls the original C implementation.
func describeH3ErrorC(code uint32) string { return C.GoString(C.describeH3Error(C.uint(code))) }

