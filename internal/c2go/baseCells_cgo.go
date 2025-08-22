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
