//go:build cgo

package c2go

/*
#include <stdint.h>
#include <stdbool.h>
#include <stdlib.h>
#include "h3api.h"
#include "h3Index.h"

// Inline C helpers to expose macro behavior for parity tests
static inline int h3_get_reserved_bits_c(H3Index h) { return H3_GET_RESERVED_BITS(h); }
static inline H3Index h3_set_reserved_bits_c(H3Index h, int v) { H3Index x = h; H3_SET_RESERVED_BITS(x, v); return x; }
static inline int h3_get_index_digit_c(H3Index h, int res) { return (int)H3_GET_INDEX_DIGIT(h, res); }
static inline H3Index h3_set_index_digit_c(H3Index h, int res, int digit) { H3Index x = h; H3_SET_INDEX_DIGIT(x, res, digit); return x; }
static inline int h3_get_mode_c(H3Index h) { return H3_GET_MODE(h); }
static inline H3Index h3_set_mode_c(H3Index h, int v) { H3Index x = h; H3_SET_MODE(x, v); return x; }
static inline int h3_get_high_bit_c(H3Index h) { return H3_GET_HIGH_BIT(h); }
static inline H3Index h3_set_high_bit_c(H3Index h, int v) { H3Index x = h; H3_SET_HIGH_BIT(x, v); return x; }
// Direct bridges for small helpers under test
static inline H3Index zero_index_digits_c(H3Index h, int start, int end) { return _zeroIndexDigits(h, start, end); }
// Prototype for non-exported helper in h3Index.c
extern H3Index makeDirectChild(H3Index h, int cellNumber);
static inline H3Index make_direct_child_c(H3Index h, int cellNumber) { return makeDirectChild(h, cellNumber); }
// Prototypes for cgo-exposed helpers in h3lib_h3Index_c2go.c
extern int has_child_at_res_c(H3Index h, int childRes);
extern int first_one_index_c(H3Index h);
extern int has_good_top_bits_c(H3Index h);
extern int has_any_7_upto_res_c(H3Index h, int res);
extern int has_all_7_after_res_c(H3Index h, int res);
extern int has_deleted_subsequence_c(H3Index h, int baseCell);
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

// getReservedBitsC exposes H3_GET_RESERVED_BITS.
func getReservedBitsC(h H3Index) int { return int(C.h3_get_reserved_bits_c(C.H3Index(h))) }

// setReservedBitsC exposes H3_SET_RESERVED_BITS.
func setReservedBitsC(h H3Index, v int) H3Index {
	return H3Index(C.h3_set_reserved_bits_c(C.H3Index(h), C.int(v)))
}

// getIndexDigitC exposes H3_GET_INDEX_DIGIT.
func getIndexDigitC(h H3Index, res int) int {
	return int(C.h3_get_index_digit_c(C.H3Index(h), C.int(res)))
}

// setIndexDigitC exposes H3_SET_INDEX_DIGIT.
func setIndexDigitC(h H3Index, res int, digit int) H3Index {
	return H3Index(C.h3_set_index_digit_c(C.H3Index(h), C.int(res), C.int(digit)))
}

// h3LeadingNonZeroDigitC calls the original C implementation.
func h3LeadingNonZeroDigitC(h H3Index) int { return int(C._h3LeadingNonZeroDigit(C.H3Index(h))) }

// zeroIndexDigitsC calls the original C implementation.
func zeroIndexDigitsC(h H3Index, start, end int) H3Index {
	return H3Index(C.zero_index_digits_c(C.H3Index(h), C.int(start), C.int(end)))
}

// isResClassIIIC calls the original C implementation and returns int for parity.
func isResClassIIIC(h H3Index) int { return int(C.isResClassIII(C.H3Index(h))) }

// isPentagonC calls the original C implementation and returns int for parity.
func isPentagonC(h H3Index) int { return int(C.isPentagon(C.H3Index(h))) }

// pentagonCountC calls the original C implementation.
func pentagonCountC() int { return int(C.pentagonCount()) }

// getModeC exposes H3_GET_MODE.
func getModeC(h H3Index) int { return int(C.h3_get_mode_c(C.H3Index(h))) }

// setModeC exposes H3_SET_MODE.
func setModeC(h H3Index, v int) H3Index { return H3Index(C.h3_set_mode_c(C.H3Index(h), C.int(v))) }

// getHighBitC exposes H3_GET_HIGH_BIT.
func getHighBitC(h H3Index) int { return int(C.h3_get_high_bit_c(C.H3Index(h))) }

// setHighBitC exposes H3_SET_HIGH_BIT.
func setHighBitC(h H3Index, v int) H3Index {
	return H3Index(C.h3_set_high_bit_c(C.H3Index(h), C.int(v)))
}

// isResolutionClassIIIC bridges to the C helper taking a resolution.
func isResolutionClassIIIC(res int) int { return int(C.isResolutionClassIII(C.int(res))) }

// setH3IndexC calls the original C implementation and returns the built index.
func setH3IndexC(res, baseCell, initDigit int) H3Index {
	var out C.H3Index
	C.setH3Index(&out, C.int(res), C.int(baseCell), C.Direction(initDigit))
	return H3Index(out)
}

// makeDirectChildC calls the original C implementation.
func makeDirectChildC(h H3Index, cellNumber int) H3Index {
	return H3Index(C.make_direct_child_c(C.H3Index(h), C.int(cellNumber)))
}

// cellToParentC calls the original C implementation.
func cellToParentC(h H3Index, parentRes int) (H3Index, uint32) {
	var out C.H3Index
	err := C.cellToParent(C.H3Index(h), C.int(parentRes), &out)
	return H3Index(out), uint32(err)
}

// cellToCenterChildC calls the original C implementation.
func cellToCenterChildC(h H3Index, childRes int) (H3Index, uint32) {
	var out C.H3Index
	err := C.cellToCenterChild(C.H3Index(h), C.int(childRes), &out)
	return H3Index(out), uint32(err)
}

// hasChildAtResC bridges to the static helper via our shim.
func hasChildAtResC(h H3Index, childRes int) int {
	return int(C.has_child_at_res_c(C.H3Index(h), C.int(childRes)))
}

// cellToChildrenSizeC calls the original C implementation.
func cellToChildrenSizeC(h H3Index, childRes int) (int64, uint32) {
	var out C.int64_t
	err := C.cellToChildrenSize(C.H3Index(h), C.int(childRes), &out)
	return int64(out), uint32(err)
}

// cellToChildPosC calls the original C implementation.
func cellToChildPosC(child H3Index, parentRes int) (int64, uint32) {
	var out C.int64_t
	err := C.cellToChildPos(C.H3Index(child), C.int(parentRes), &out)
	return int64(out), uint32(err)
}

// childPosToCellC calls the original C implementation.
func childPosToCellC(childPos int64, parent H3Index, childRes int) (H3Index, uint32) {
	var out C.H3Index
	err := C.childPosToCell(C.int64_t(childPos), C.H3Index(parent), C.int(childRes), &out)
	return H3Index(out), uint32(err)
}

// getPentagonsC calls the original C implementation and returns slice + err.
func getPentagonsC(res int) ([]H3Index, uint32) {
	n := pentagonCountC()
	if n <= 0 {
		return nil, uint32(E_FAILED)
	}
	buf := make([]H3Index, n)
	err := C.getPentagons(C.int(res), (*C.H3Index)(&buf[0]))
	return buf, uint32(err)
}

// firstOneIndexC calls the original C implementation.
func firstOneIndexC(h H3Index) int {
	return int(C.first_one_index_c(C.H3Index(h)))
}

// hasGoodTopBitsC calls the original C implementation.
func hasGoodTopBitsC(h H3Index) bool {
	return C.has_good_top_bits_c(C.H3Index(h)) != 0
}

// hasAny7UptoResC calls the original C implementation.
func hasAny7UptoResC(h H3Index, res int) bool {
	return C.has_any_7_upto_res_c(C.H3Index(h), C.int(res)) != 0
}

// hasAll7AfterResC calls the original C implementation.
func hasAll7AfterResC(h H3Index, res int) bool {
	return C.has_all_7_after_res_c(C.H3Index(h), C.int(res)) != 0
}

// hasDeletedSubsequenceC calls the original C implementation.
func hasDeletedSubsequenceC(h H3Index, baseCell int) bool {
	return C.has_deleted_subsequence_c(C.H3Index(h), C.int(baseCell)) != 0
}

// h3Rotate60ccwC calls the original C implementation.
func h3Rotate60ccwC(h H3Index) H3Index {
	return H3Index(C._h3Rotate60ccw(C.H3Index(h)))
}

// h3Rotate60cwC calls the original C implementation.
func h3Rotate60cwC(h H3Index) H3Index {
	return H3Index(C._h3Rotate60cw(C.H3Index(h)))
}

// h3RotatePent60ccwC calls the original C implementation.
func h3RotatePent60ccwC(h H3Index) H3Index {
	return H3Index(C._h3RotatePent60ccw(C.H3Index(h)))
}

// h3RotatePent60cwC calls the original C implementation.
func h3RotatePent60cwC(h H3Index) H3Index {
	return H3Index(C._h3RotatePent60cw(C.H3Index(h)))
}

// _h3ToFaceIjkWithInitializedFijkC calls the original C implementation.
func _h3ToFaceIjkWithInitializedFijkC(h H3Index, fijk *FaceIJK) int {
	var cFijk C.FaceIJK
	cFijk.face = C.int(fijk.Face)
	cFijk.coord.i = C.int(fijk.Coord.I)
	cFijk.coord.j = C.int(fijk.Coord.J)
	cFijk.coord.k = C.int(fijk.Coord.K)

	result := int(C._h3ToFaceIjkWithInitializedFijk(C.H3Index(h), &cFijk))

	// Update the Go struct with results
	fijk.Face = int(cFijk.face)
	fijk.Coord.I = int(cFijk.coord.i)
	fijk.Coord.J = int(cFijk.coord.j)
	fijk.Coord.K = int(cFijk.coord.k)

	return result
}

// _h3ToFaceIjkC calls the original C implementation.
func _h3ToFaceIjkC(h H3Index, fijk *FaceIJK) uint32 {
	var cFijk C.FaceIJK
	err := C._h3ToFaceIjk(C.H3Index(h), &cFijk)

	// Update the Go struct with results
	fijk.Face = int(cFijk.face)
	fijk.Coord.I = int(cFijk.coord.i)
	fijk.Coord.J = int(cFijk.coord.j)
	fijk.Coord.K = int(cFijk.coord.k)

	return uint32(err)
}

// cellToLatLngC calls the original C implementation.
func cellToLatLngC(h H3Index, g *LatLng) uint32 {
	var cg C.LatLng
	err := C.cellToLatLng(C.H3Index(h), &cg)

	g.Lat = float64(cg.lat)
	g.Lng = float64(cg.lng)

	return uint32(err)
}

// isValidCellC calls the original C implementation.
func isValidCellC(h H3Index) int {
	return int(C.isValidCell(C.H3Index(h)))
}

// latLngToCellC calls the original C implementation.
func latLngToCellC(g *LatLng, res int, out *H3Index) uint32 {
	var cg C.LatLng
	cg.lat = C.double(g.Lat)
	cg.lng = C.double(g.Lng)

	var cOut C.H3Index
	err := C.latLngToCell(&cg, C.int(res), &cOut)

	*out = H3Index(cOut)
	return uint32(err)
}

// debugGeoToFaceIjkC calls the original C _geoToFaceIjk implementation.
func debugGeoToFaceIjkC(g *LatLng, res int, fijk *FaceIJK) {
	var cg C.LatLng
	cg.lat = C.double(g.Lat)
	cg.lng = C.double(g.Lng)

	var cfijk C.FaceIJK
	C._geoToFaceIjk(&cg, C.int(res), &cfijk)

	fijk.Face = int(cfijk.face)
	fijk.Coord.I = int(cfijk.coord.i)
	fijk.Coord.J = int(cfijk.coord.j)
	fijk.Coord.K = int(cfijk.coord.k)
}

// debugFaceIjkToH3C calls the original C _faceIjkToH3 implementation.
func debugFaceIjkToH3C(fijk *FaceIJK, res int) H3Index {
	var cfijk C.FaceIJK
	cfijk.face = C.int(fijk.Face)
	cfijk.coord.i = C.int(fijk.Coord.I)
	cfijk.coord.j = C.int(fijk.Coord.J)
	cfijk.coord.k = C.int(fijk.Coord.K)

	return H3Index(C._faceIjkToH3(&cfijk, C.int(res)))
}
