// Tests ported from testCellToChildPos.c
package h3

import (
	"testing"
)

// Helper function to iterate all indexes at a given resolution
func iterateAllIndexesAtResForChildPos(t *testing.T, res int32, testFunc func(t *testing.T, h3 H3Index)) {
	t.Helper()

	// Get all base cells
	baseCells := make([]H3Index, NUM_BASE_CELLS)
	if err := getRes0Cells(baseCells); err != E_SUCCESS {
		t.Fatalf("Failed to get res 0 cells: %v", err)
	}

	if res == 0 {
		// For resolution 0, just test the base cells
		for _, cell := range baseCells {
			testFunc(t, cell)
		}
		return
	}

	// For higher resolutions, get children of each base cell
	for _, baseCell := range baseCells {
		childrenSize, err := cellToChildrenSize(baseCell, res)
		if err != E_SUCCESS {
			continue // Some cells might not have children at certain resolutions
		}

		children := make([]H3Index, childrenSize)
		if err := cellToChildren(baseCell, res, children); err != E_SUCCESS {
			continue
		}

		for _, child := range children {
			if child != H3_NULL {
				testFunc(t, child)
			}
		}
	}
}

// childPos_assertions tests cellToChildPos and childPosToCell for a given H3 index
func childPos_assertions(t *testing.T, h3 H3Index) {
	t.Helper()

	parentRes := getResolution(h3)

	for resolutionOffset := int32(0); resolutionOffset < 4; resolutionOffset++ {
		childRes := parentRes + resolutionOffset
		numChildren, err := cellToChildrenSize(h3, childRes)
		if err != E_SUCCESS {
			continue
		}

		children := make([]H3Index, numChildren)
		err = cellToChildren(h3, childRes, children)
		if err != E_SUCCESS {
			continue
		}

		for i, child := range children {
			// Test cellToChildPos
			childPos, err := cellToChildPos(child, parentRes)
			if err != E_SUCCESS {
				t.Errorf("cellToChildPos failed for child %d: %v", i, err)
				continue
			}
			if childPos != int64(i) {
				t.Errorf("childPos matches iteration index: expected %d, got %d", i, childPos)
			}

			// Test childPosToCell
			cell, err := childPosToCell(int64(i), h3, childRes)
			if err != E_SUCCESS {
				t.Errorf("childPosToCell failed for position %d: %v", i, err)
				continue
			}
			if cell != children[i] {
				t.Errorf("cell matches expected: expected 0x%x, got 0x%x", children[i], cell)
			}
		}
	}
}

func Test_childPos_correctness(t *testing.T) {
	t.Parallel()

	iterateAllIndexesAtResForChildPos(t, 0, childPos_assertions)
	iterateAllIndexesAtResForChildPos(t, 1, childPos_assertions)
	iterateAllIndexesAtResForChildPos(t, 2, childPos_assertions)
}

func Test_cellToChildPos_res_errors(t *testing.T) {
	t.Parallel()

	// random res 8 cell
	child := H3Index(0x88283080ddfffff)

	// Test invalid resolution -1
	_, err := cellToChildPos(child, -1)
	if err != E_RES_DOMAIN {
		t.Errorf("error matches expected for invalid res: expected E_RES_DOMAIN, got %v", err)
	}

	// Test invalid resolution 42
	_, err = cellToChildPos(child, 42)
	if err != E_RES_DOMAIN {
		t.Errorf("error matches expected for invalid res: expected E_RES_DOMAIN, got %v", err)
	}

	// Test parent res finer than child
	_, err = cellToChildPos(child, 9)
	if err != E_RES_MISMATCH {
		t.Errorf("error matches expected for parent res finer than child: expected E_RES_MISMATCH, got %v", err)
	}
}

func Test_childPosToCell_res_errors(t *testing.T) {
	t.Parallel()

	// random res 8 cell
	parent := H3Index(0x88283080ddfffff)
	childPos := int64(27)

	// Test invalid resolution 42
	_, err := childPosToCell(childPos, parent, 42)
	if err != E_RES_DOMAIN {
		t.Errorf("error matches expected for invalid res: expected E_RES_DOMAIN, got %v", err)
	}

	// Test invalid resolution -1
	_, err = childPosToCell(childPos, parent, -1)
	if err != E_RES_DOMAIN {
		t.Errorf("error matches expected for invalid res: expected E_RES_DOMAIN, got %v", err)
	}

	// Test child res coarser than parent
	_, err = childPosToCell(childPos, parent, 7)
	if err != E_RES_MISMATCH {
		t.Errorf("error matches expected for child res coarser than parent: expected E_RES_MISMATCH, got %v", err)
	}
}

func Test_childPosToCell_childPos_errors(t *testing.T) {
	t.Parallel()

	// random res 8 cell
	parent := H3Index(0x88283080ddfffff)
	res := int32(10)

	// Test negative childPos
	_, err := childPosToCell(-1, parent, res)
	if err != E_DOMAIN {
		t.Errorf("error matches expected for negative childPos: expected E_DOMAIN, got %v", err)
	}

	// res is two steps down, so max valid child pos is 48
	_, err = childPosToCell(48, parent, res)
	if err != E_SUCCESS {
		t.Errorf("No error for max valid child pos: expected E_SUCCESS, got %v", err)
	}

	_, err = childPosToCell(49, parent, res)
	if err != E_DOMAIN {
		t.Errorf("error matches expected for childPos greater than max: expected E_DOMAIN, got %v", err)
	}
}

func Test_cellToChildPos_invalid_digit(t *testing.T) {
	t.Parallel()

	// random res 8 cell
	child := H3Index(0x88283080ddfffff)
	child = setIndexDigit(child, 6, int32(INVALID_DIGIT))

	_, err := cellToChildPos(child, 0)
	if err != E_CELL_INVALID {
		t.Errorf("error matches expected for invalid cell: expected E_CELL_INVALID, got %v", err)
	}
}

func Test_cellToChildPos_invalid_pentagon_digit(t *testing.T) {
	t.Parallel()

	// Res 7 hexagon child of a pentagon
	child := H3Index(0x870800006ffffff)
	child = setIndexDigit(child, 7, int32(INVALID_DIGIT))

	_, err := cellToChildPos(child, 0)
	if err != E_CELL_INVALID {
		t.Errorf("error matches expected for invalid cell: expected E_CELL_INVALID, got %v", err)
	}
}

func Test_cellToChildPos_invalid_pentagon_kaxis(t *testing.T) {
	t.Parallel()

	// Create a res 8 index located in a deleted subsequence of a pentagon.
	var child H3Index
	setH3Index(&child, 8, 4, int32(K_AXES_DIGIT))

	_, err := cellToChildPos(child, 0)
	if err != E_CELL_INVALID {
		t.Errorf("error matches expected for invalid cell: expected E_CELL_INVALID, got %v", err)
	}
}
