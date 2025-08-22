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

