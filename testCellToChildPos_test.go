// Tests ported from testCellToChildPos.c
package h3

import (
	"testing"
)

// Helper function to iterate all indexes at a given resolution.
func iterateAllIndexesAtResForChildPos(t *testing.T, res int32, testFunc func(t *testing.T, h3 h3Index)) {
	t.Helper()

	// Get all base cells
	baseCells := make([]h3Index, numBaseCells)
	if err := getRes0Cells(baseCells); err != eSuccess {
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
		if err != eSuccess {
			continue // Some cells might not have children at certain resolutions
		}

		children := make([]h3Index, childrenSize)
		if err := cellToChildren(baseCell, res, children); err != eSuccess {
			continue
		}

		for _, child := range children {
			if child != h3Null {
				testFunc(t, child)
			}
		}
	}
}

// childPos_assertions tests cellToChildPos and childPosToCell for a given H3 index.
func childPos_assertions(t *testing.T, h3 h3Index) {
	t.Helper()

	parentRes := getResolution(h3)

	for resolutionOffset := int32(0); resolutionOffset < 4; resolutionOffset++ {
		childRes := parentRes + resolutionOffset
		numChildren, err := cellToChildrenSize(h3, childRes)
		if err != eSuccess {
			continue
		}

		children := make([]h3Index, numChildren)
		err = cellToChildren(h3, childRes, children)
		if err != eSuccess {
			continue
		}

		for i, child := range children {
			// Test cellToChildPos
			childPos, err := cellToChildPos(child, parentRes)
			if err != eSuccess {
				t.Errorf("cellToChildPos failed for child %d: %v", i, err)
				continue
			}
			if childPos != int64(i) {
				t.Errorf("childPos matches iteration index: expected %d, got %d", i, childPos)
			}

			// Test childPosToCell
			cell, err := childPosToCell(int64(i), h3, childRes)
			if err != eSuccess {
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
	child := h3Index(0x88283080ddfffff)

	// Test invalid resolution -1
	_, err := cellToChildPos(child, -1)
	if err != eResDomain {
		t.Errorf("error matches expected for invalid res: expected eResDomain, got %v", err)
	}

	// Test invalid resolution 42
	_, err = cellToChildPos(child, 42)
	if err != eResDomain {
		t.Errorf("error matches expected for invalid res: expected eResDomain, got %v", err)
	}

	// Test parent res finer than child
	_, err = cellToChildPos(child, 9)
	if err != eResMismatch {
		t.Errorf("error matches expected for parent res finer than child: expected eResMismatch, got %v", err)
	}
}

func Test_childPosToCell_res_errors(t *testing.T) {
	t.Parallel()

	// random res 8 cell
	parent := h3Index(0x88283080ddfffff)
	childPos := int64(27)

	// Test invalid resolution 42
	_, err := childPosToCell(childPos, parent, 42)
	if err != eResDomain {
		t.Errorf("error matches expected for invalid res: expected eResDomain, got %v", err)
	}

	// Test invalid resolution -1
	_, err = childPosToCell(childPos, parent, -1)
	if err != eResDomain {
		t.Errorf("error matches expected for invalid res: expected eResDomain, got %v", err)
	}

	// Test child res coarser than parent
	_, err = childPosToCell(childPos, parent, 7)
	if err != eResMismatch {
		t.Errorf("error matches expected for child res coarser than parent: expected eResMismatch, got %v", err)
	}
}

func Test_childPosToCell_childPos_errors(t *testing.T) {
	t.Parallel()

	// random res 8 cell
	parent := h3Index(0x88283080ddfffff)
	res := int32(10)

	// Test negative childPos
	_, err := childPosToCell(-1, parent, res)
	if err != eDomain {
		t.Errorf("error matches expected for negative childPos: expected eDomain, got %v", err)
	}

	// res is two steps down, so max valid child pos is 48
	_, err = childPosToCell(48, parent, res)
	if err != eSuccess {
		t.Errorf("No error for max valid child pos: expected eSuccess, got %v", err)
	}

	_, err = childPosToCell(49, parent, res)
	if err != eDomain {
		t.Errorf("error matches expected for childPos greater than max: expected eDomain, got %v", err)
	}
}

func Test_cellToChildPos_invalid_digit(t *testing.T) {
	t.Parallel()

	// random res 8 cell
	child := h3Index(0x88283080ddfffff)
	child = setIndexDigit(child, 6, int32(invalidDigit))

	_, err := cellToChildPos(child, 0)
	if err != eCellInvalid {
		t.Errorf("error matches expected for invalid cell: expected eCellInvalid, got %v", err)
	}
}

func Test_cellToChildPos_invalid_pentagon_digit(t *testing.T) {
	t.Parallel()

	// Res 7 hexagon child of a pentagon
	child := h3Index(0x870800006ffffff)
	child = setIndexDigit(child, 7, int32(invalidDigit))

	_, err := cellToChildPos(child, 0)
	if err != eCellInvalid {
		t.Errorf("error matches expected for invalid cell: expected eCellInvalid, got %v", err)
	}
}

func Test_cellToChildPos_invalid_pentagon_kaxis(t *testing.T) {
	t.Parallel()

	// Create a res 8 index located in a deleted subsequence of a pentagon.
	var child h3Index
	setH3Index(&child, 8, 4, int32(kAxesDigit))

	_, err := cellToChildPos(child, 0)
	if err != eCellInvalid {
		t.Errorf("error matches expected for invalid cell: expected eCellInvalid, got %v", err)
	}
}
