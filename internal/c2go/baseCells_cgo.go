//go:build cgo

package c2go

/*
#include <stdbool.h>
#include "baseCells.h"
*/
import "C"

// isBaseCellPentagonC calls the original C helper _isBaseCellPentagon.
// Returns 1 if pentagon, 0 otherwise.
func isBaseCellPentagonC(base int32) bool { return int32(C._isBaseCellPentagon(C.int(base))) != 0 }

// _isBaseCellPolarPentagonC calls the original C helper _isBaseCellPolarPentagon.
// Returns 1 if polar pentagon, 0 otherwise.
func _isBaseCellPolarPentagonC(baseCell int32) bool {
	return bool(C._isBaseCellPolarPentagon(C.int(baseCell)))
}

// res0CellCountC calls the original C implementation res0CellCount.
func res0CellCountC() int32 { return int32(C.res0CellCount()) }

// _baseCellToFaceIjkC calls the original C helper _baseCellToFaceIjk.
func _baseCellToFaceIjkC(baseCell int32) FaceIJK {
	var h C.FaceIJK
	C._baseCellToFaceIjk(C.int(baseCell), &h)
	return FaceIJK{
		Face: int32(h.face),
		Coord: CoordIJK{
			I: int32(h.coord.i),
			J: int32(h.coord.j),
			K: int32(h.coord.k),
		},
	}
}

// _baseCellToCCWrot60C calls the original C helper _baseCellToCCWrot60.
func _baseCellToCCWrot60C(baseCell int32, face int32) int32 {
	return int32(C._baseCellToCCWrot60(C.int(baseCell), C.int(face)))
}

// _baseCellIsCwOffsetC calls the original C helper _baseCellIsCwOffset.
func _baseCellIsCwOffsetC(baseCell int32, testFace int32) bool {
	return bool(C._baseCellIsCwOffset(C.int(baseCell), C.int(testFace)))
}

// _faceIjkToBaseCellC calls the original C helper _faceIjkToBaseCell.
func _faceIjkToBaseCellC(h *FaceIJK) int32 {
	var cFijk C.FaceIJK
	cFijk.face = C.int(h.Face)
	cFijk.coord.i = C.int(h.Coord.I)
	cFijk.coord.j = C.int(h.Coord.J)
	cFijk.coord.k = C.int(h.Coord.K)
	return int32(C._faceIjkToBaseCell(&cFijk))
}

// _faceIjkToBaseCellCCWrot60C calls the original C helper _faceIjkToBaseCellCCWrot60.
func _faceIjkToBaseCellCCWrot60C(h *FaceIJK) int32 {
	var cFijk C.FaceIJK
	cFijk.face = C.int(h.Face)
	cFijk.coord.i = C.int(h.Coord.I)
	cFijk.coord.j = C.int(h.Coord.J)
	cFijk.coord.k = C.int(h.Coord.K)
	return int32(C._faceIjkToBaseCellCCWrot60(&cFijk))
}

// _getBaseCellNeighborC calls the original C helper _getBaseCellNeighbor.
func _getBaseCellNeighborC(baseCell int32, dir Direction) int32 {
	return int32(C._getBaseCellNeighbor(C.int(baseCell), C.Direction(dir)))
}

// _getBaseCellDirectionC calls the original C helper _getBaseCellDirection.
func _getBaseCellDirectionC(originBaseCell int32, neighboringBaseCell int32) Direction {
	return Direction(C._getBaseCellDirection(C.int(originBaseCell), C.int(neighboringBaseCell)))
}

// getRes0CellsC calls the original C implementation getRes0Cells.
func getRes0CellsC(out []H3Index) H3Error {
	if len(out) != NUM_BASE_CELLS {
		return E_FAILED // Need exactly NUM_BASE_CELLS slots
	}
	cOut := make([]C.H3Index, NUM_BASE_CELLS)
	err := C.getRes0Cells(&cOut[0])
	for i := 0; i < NUM_BASE_CELLS; i++ {
		out[i] = H3Index(cOut[i])
	}
	return H3Error(err)
}
