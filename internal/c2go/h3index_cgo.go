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
func getResolutionC(h H3Index) int32 { return int32(C.getResolution(C.H3Index(h))) }

// getBaseCellNumberC calls the original C implementation.
func getBaseCellNumberC(h H3Index) int32 { return int32(C.getBaseCellNumber(C.H3Index(h))) }

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
func getReservedBitsC(h H3Index) int32 { return int32(C.h3_get_reserved_bits_c(C.H3Index(h))) }

// setReservedBitsC exposes H3_SET_RESERVED_BITS.
func setReservedBitsC(h H3Index, v int32) H3Index {
	return H3Index(C.h3_set_reserved_bits_c(C.H3Index(h), C.int(v)))
}

// getIndexDigitC exposes H3_GET_INDEX_DIGIT.
func getIndexDigitC(h H3Index, res int32) int32 {
	return int32(C.h3_get_index_digit_c(C.H3Index(h), C.int(res)))
}

// setIndexDigitC exposes H3_SET_INDEX_DIGIT.
func setIndexDigitC(h H3Index, res int32, digit int32) H3Index {
	return H3Index(C.h3_set_index_digit_c(C.H3Index(h), C.int(res), C.int(digit)))
}

// h3LeadingNonZeroDigitC calls the original C implementation.
func h3LeadingNonZeroDigitC(h H3Index) int32 { return int32(C._h3LeadingNonZeroDigit(C.H3Index(h))) }

// zeroIndexDigitsC calls the original C implementation.
func zeroIndexDigitsC(h H3Index, start, end int32) H3Index {
	return H3Index(C.zero_index_digits_c(C.H3Index(h), C.int(start), C.int(end)))
}

// isResClassIIIC calls the original C implementation and returns int for parity.
func isResClassIIIC(h H3Index) bool { return int32(C.isResClassIII(C.H3Index(h))) != 0 }

// isPentagonC calls the original C implementation and returns int for parity.
func isPentagonC(h H3Index) bool { return int32(C.isPentagon(C.H3Index(h))) != 0 }

// pentagonCountC calls the original C implementation.
func pentagonCountC() int32 { return int32(C.pentagonCount()) }

// getModeC exposes H3_GET_MODE.
func getModeC(h H3Index) int32 { return int32(C.h3_get_mode_c(C.H3Index(h))) }

// setModeC exposes H3_SET_MODE.
func setModeC(h H3Index, v int32) H3Index { return H3Index(C.h3_set_mode_c(C.H3Index(h), C.int(v))) }

// getHighBitC exposes H3_GET_HIGH_BIT.
func getHighBitC(h H3Index) int32 { return int32(C.h3_get_high_bit_c(C.H3Index(h))) }

// setHighBitC exposes H3_SET_HIGH_BIT.
func setHighBitC(h H3Index, v int32) H3Index {
	return H3Index(C.h3_set_high_bit_c(C.H3Index(h), C.int(v)))
}

// isResolutionClassIIIC bridges to the C helper taking a resolution.
func isResolutionClassIIIC(res int32) int32 { return int32(C.isResolutionClassIII(C.int(res))) }

// setH3IndexC calls the original C implementation and returns the built index.
func setH3IndexC(res, baseCell, initDigit int32) H3Index {
	var out C.H3Index
	C.setH3Index(&out, C.int(res), C.int(baseCell), C.Direction(initDigit))
	return H3Index(out)
}

// makeDirectChildC calls the original C implementation.
func makeDirectChildC(h H3Index, cellNumber int32) H3Index {
	return H3Index(C.make_direct_child_c(C.H3Index(h), C.int(cellNumber)))
}

// cellToParentC calls the original C implementation.
func cellToParentC(h H3Index, parentRes int32) (H3Index, uint32) {
	var out C.H3Index
	err := C.cellToParent(C.H3Index(h), C.int(parentRes), &out)
	return H3Index(out), uint32(err)
}

// cellToCenterChildC calls the original C implementation.
func cellToCenterChildC(h H3Index, childRes int32) (H3Index, uint32) {
	var out C.H3Index
	err := C.cellToCenterChild(C.H3Index(h), C.int(childRes), &out)
	return H3Index(out), uint32(err)
}

// hasChildAtResC bridges to the static helper via our shim.
func hasChildAtResC(h H3Index, childRes int32) int32 {
	return int32(C.has_child_at_res_c(C.H3Index(h), C.int(childRes)))
}

// cellToChildrenSizeC calls the original C implementation.
func cellToChildrenSizeC(h H3Index, childRes int32) (int64, uint32) {
	var out C.int64_t
	err := C.cellToChildrenSize(C.H3Index(h), C.int(childRes), &out)
	return int64(out), uint32(err)
}

// cellToChildPosC calls the original C implementation.
func cellToChildPosC(child H3Index, parentRes int32) (int64, uint32) {
	var out C.int64_t
	err := C.cellToChildPos(C.H3Index(child), C.int(parentRes), &out)
	return int64(out), uint32(err)
}

// childPosToCellC calls the original C implementation.
func childPosToCellC(childPos int64, parent H3Index, childRes int32) (H3Index, uint32) {
	var out C.H3Index
	err := C.childPosToCell(C.int64_t(childPos), C.H3Index(parent), C.int(childRes), &out)
	return H3Index(out), uint32(err)
}

// getPentagonsC calls the original C implementation and returns slice + err.
func getPentagonsC(res int32) ([]H3Index, uint32) {
	n := pentagonCountC()
	if n <= 0 {
		return nil, uint32(E_FAILED)
	}
	buf := make([]H3Index, n)
	err := C.getPentagons(C.int(res), (*C.H3Index)(&buf[0]))
	return buf, uint32(err)
}

// firstOneIndexC calls the original C implementation.
func firstOneIndexC(h H3Index) int32 {
	return int32(C.first_one_index_c(C.H3Index(h)))
}

// hasGoodTopBitsC calls the original C implementation.
func hasGoodTopBitsC(h H3Index) bool {
	return C.has_good_top_bits_c(C.H3Index(h)) != 0
}

// hasAny7UptoResC calls the original C implementation.
func hasAny7UptoResC(h H3Index, res int32) bool {
	return C.has_any_7_upto_res_c(C.H3Index(h), C.int(res)) != 0
}

// hasAll7AfterResC calls the original C implementation.
func hasAll7AfterResC(h H3Index, res int32) bool {
	return C.has_all_7_after_res_c(C.H3Index(h), C.int(res)) != 0
}

// hasDeletedSubsequenceC calls the original C implementation.
func hasDeletedSubsequenceC(h H3Index, baseCell int32) bool {
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
func _h3ToFaceIjkWithInitializedFijkC(h H3Index, fijk *FaceIJK) int32 {
	var cFijk C.FaceIJK
	cFijk.face = C.int(fijk.Face)
	cFijk.coord.i = C.int(fijk.Coord.I)
	cFijk.coord.j = C.int(fijk.Coord.J)
	cFijk.coord.k = C.int(fijk.Coord.K)

	result := int32(C._h3ToFaceIjkWithInitializedFijk(C.H3Index(h), &cFijk))

	// Update the Go struct with results
	fijk.Face = int32(cFijk.face)
	fijk.Coord.I = int32(cFijk.coord.i)
	fijk.Coord.J = int32(cFijk.coord.j)
	fijk.Coord.K = int32(cFijk.coord.k)

	return result
}

// _h3ToFaceIjkC calls the original C implementation.
func _h3ToFaceIjkC(h H3Index, fijk *FaceIJK) uint32 {
	var cFijk C.FaceIJK
	err := C._h3ToFaceIjk(C.H3Index(h), &cFijk)

	// Update the Go struct with results
	fijk.Face = int32(cFijk.face)
	fijk.Coord.I = int32(cFijk.coord.i)
	fijk.Coord.J = int32(cFijk.coord.j)
	fijk.Coord.K = int32(cFijk.coord.k)

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
func isValidCellC(h H3Index) bool {
	return int32(C.isValidCell(C.H3Index(h))) != 0
}

// latLngToCellC calls the original C implementation.
func latLngToCellC(g *LatLng, res int32, out *H3Index) uint32 {
	var cg C.LatLng
	cg.lat = C.double(g.Lat)
	cg.lng = C.double(g.Lng)

	var cOut C.H3Index
	err := C.latLngToCell(&cg, C.int(res), &cOut)

	*out = H3Index(cOut)
	return uint32(err)
}

// debugGeoToFaceIjkC calls the original C _geoToFaceIjk implementation.
func debugGeoToFaceIjkC(g *LatLng, res int32, fijk *FaceIJK) {
	var cg C.LatLng
	cg.lat = C.double(g.Lat)
	cg.lng = C.double(g.Lng)

	var cfijk C.FaceIJK
	C._geoToFaceIjk(&cg, C.int(res), &cfijk)

	fijk.Face = int32(cfijk.face)
	fijk.Coord.I = int32(cfijk.coord.i)
	fijk.Coord.J = int32(cfijk.coord.j)
	fijk.Coord.K = int32(cfijk.coord.k)
}

// debugFaceIjkToH3C calls the original C _faceIjkToH3 implementation.
func debugFaceIjkToH3C(fijk *FaceIJK, res int32) H3Index {
	var cfijk C.FaceIJK
	cfijk.face = C.int(fijk.Face)
	cfijk.coord.i = C.int(fijk.Coord.I)
	cfijk.coord.j = C.int(fijk.Coord.J)
	cfijk.coord.k = C.int(fijk.Coord.K)

	return H3Index(C._faceIjkToH3(&cfijk, C.int(res)))
}

// maxFaceCountC calls the original C implementation.
func maxFaceCountC(h H3Index, out *int32) uint32 {
	var cOut C.int
	err := C.maxFaceCount(C.H3Index(h), &cOut)
	*out = int32(cOut)
	return uint32(err)
}

// getIcosahedronFacesC calls the original C implementation.
func getIcosahedronFacesC(h H3Index, out []int32) uint32 {
	if len(out) == 0 {
		return uint32(C.E_FAILED)
	}
	cOut := (*C.int)(C.malloc(C.size_t(len(out)) * C.size_t(C.sizeof_int)))
	defer C.free(unsafe.Pointer(cOut))
	err := C.getIcosahedronFaces(C.H3Index(h), cOut)
	if err == 0 {
		// Copy results back to Go slice
		slice := (*[1 << 30]C.int)(unsafe.Pointer(cOut))[:len(out):len(out)]
		for i := range out {
			out[i] = int32(slice[i])
		}
	}
	return uint32(err)
}

// cellToChildrenC calls the original C implementation.
func cellToChildrenC(h H3Index, childRes int32, children []H3Index) uint32 {
	if len(children) == 0 {
		return uint32(C.E_FAILED)
	}
	cChildren := (*C.H3Index)(C.malloc(C.size_t(len(children)) * C.size_t(C.sizeof_H3Index)))
	defer C.free(unsafe.Pointer(cChildren))
	err := C.cellToChildren(C.H3Index(h), C.int(childRes), cChildren)
	if err == 0 {
		// Copy results back to Go slice
		slice := (*[1 << 30]C.H3Index)(unsafe.Pointer(cChildren))[:len(children):len(children)]
		for i := range children {
			children[i] = H3Index(slice[i])
		}
	}
	return uint32(err)
}

// cellToBoundaryC calls the original C implementation.
func cellToBoundaryC(h H3Index, cb *CellBoundary) uint32 {
	var cCb C.CellBoundary
	cCb.numVerts = 0

	err := C.cellToBoundary(C.H3Index(h), &cCb)

	if err == 0 {
		// Copy results back to Go struct
		cb.NumVerts = int32(cCb.numVerts)

		// Ensure Go slice has enough capacity
		if len(cb.Verts) < int(cb.NumVerts) {
			cb.Verts = make([]LatLng, cb.NumVerts)
		}

		// Copy vertices from C to Go
		for i := int32(0); i < cb.NumVerts; i++ {
			cb.Verts[i].Lat = float64(cCb.verts[i].lat)
			cb.Verts[i].Lng = float64(cCb.verts[i].lng)
		}
	}

	return uint32(err)
}

// compactCellsC calls the original C implementation.
func compactCellsC(h3Set []H3Index, compactedSet []H3Index, numHexes int64) uint32 {
	if numHexes == 0 {
		return uint32(C.E_SUCCESS)
	}
	if len(h3Set) < int(numHexes) || len(compactedSet) < int(numHexes) {
		return uint32(C.E_FAILED)
	}

	// Allocate C arrays
	cH3Set := (*C.H3Index)(C.malloc(C.size_t(numHexes) * C.size_t(C.sizeof_H3Index)))
	defer C.free(unsafe.Pointer(cH3Set))
	cCompactedSet := (*C.H3Index)(C.malloc(C.size_t(numHexes) * C.size_t(C.sizeof_H3Index)))
	defer C.free(unsafe.Pointer(cCompactedSet))

	// Copy Go slice to C array
	h3Slice := (*[1 << 30]C.H3Index)(unsafe.Pointer(cH3Set))[:numHexes:numHexes]
	for i := int64(0); i < numHexes; i++ {
		h3Slice[i] = C.H3Index(h3Set[i])
	}

	// Call C function
	err := C.compactCells(cH3Set, cCompactedSet, C.int64_t(numHexes))

	if err == 0 {
		// Copy results back to Go slice
		compactedSlice := (*[1 << 30]C.H3Index)(unsafe.Pointer(cCompactedSet))[:numHexes:numHexes]
		for i := int64(0); i < numHexes; i++ {
			compactedSet[i] = H3Index(compactedSlice[i])
		}
	}

	return uint32(err)
}

// uncompactCellsSizeC calls the original C implementation.
func uncompactCellsSizeC(compactedSet []H3Index, numCompacted int64, res int32) (int64, uint32) {
	if numCompacted == 0 {
		return 0, uint32(C.E_SUCCESS)
	}
	if len(compactedSet) < int(numCompacted) {
		return 0, uint32(C.E_FAILED)
	}

	// Allocate C array
	cCompactedSet := (*C.H3Index)(C.malloc(C.size_t(numCompacted) * C.size_t(C.sizeof_H3Index)))
	defer C.free(unsafe.Pointer(cCompactedSet))

	// Copy Go slice to C array
	compactedSlice := (*[1 << 30]C.H3Index)(unsafe.Pointer(cCompactedSet))[:numCompacted:numCompacted]
	for i := int64(0); i < numCompacted; i++ {
		compactedSlice[i] = C.H3Index(compactedSet[i])
	}

	// Call C function
	var out C.int64_t
	err := C.uncompactCellsSize(cCompactedSet, C.int64_t(numCompacted), C.int(res), &out)

	return int64(out), uint32(err)
}
