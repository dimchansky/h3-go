//go:build cgo && c2go

package c2go

/*
#include <stdbool.h>
#include "baseCells.h"
*/
import "C"

// isBaseCellPentagonC calls the original C helper _isBaseCellPentagon.
// Returns 1 if pentagon, 0 otherwise.
func isBaseCellPentagonC(base int) int { return int(C._isBaseCellPentagon(C.int(base))) }

// _isBaseCellPolarPentagonC calls the original C helper _isBaseCellPolarPentagon.
// Returns 1 if polar pentagon, 0 otherwise.
func _isBaseCellPolarPentagonC(baseCell int) bool {
	return bool(C._isBaseCellPolarPentagon(C.int(baseCell)))
}

// res0CellCountC calls the original C implementation res0CellCount.
func res0CellCountC() int { return int(C.res0CellCount()) }

// _baseCellToFaceIjkC calls the original C helper _baseCellToFaceIjk.
func _baseCellToFaceIjkC(baseCell int) FaceIJK {
	var h C.FaceIJK
	C._baseCellToFaceIjk(C.int(baseCell), &h)
	return FaceIJK{
		Face: int(h.face),
		Coord: CoordIJK{
			I: int(h.coord.i),
			J: int(h.coord.j),
			K: int(h.coord.k),
		},
	}
}

// _baseCellToCCWrot60C calls the original C helper _baseCellToCCWrot60.
func _baseCellToCCWrot60C(baseCell int, face int) int {
	return int(C._baseCellToCCWrot60(C.int(baseCell), C.int(face)))
}

// _baseCellIsCwOffsetC calls the original C helper _baseCellIsCwOffset.
func _baseCellIsCwOffsetC(baseCell int, testFace int) bool {
	return bool(C._baseCellIsCwOffset(C.int(baseCell), C.int(testFace)))
}
