// Tests ported from H3 v4.4.0: src/apps/testapps/testGridPathCells.c.
package h3

import (
	"testing"
)

func Test_gridPathCells_acrossMultipleFaces(t *testing.T) {
	t.Parallel()

	start := h3Index(0x85285aa7fffffff)
	end := h3Index(0x851d9b1bfffffff)

	var size int64
	lineError := gridPathCellsSize(start, end, &size)
	if lineError != eFailed {
		t.Errorf("Expected eFailed for line not computable across multiple icosa faces, got %v", lineError)
	}
}

func Test_gridPathCells_pentagon(t *testing.T) {
	t.Parallel()

	start := h3Index(0x820807fffffffff)
	end := h3Index(0x8208e7fffffffff)

	var size int64
	err := gridPathCellsSize(start, end, &size)
	if err != eSuccess {
		t.Fatalf("gridPathCellsSize failed: %v", err)
	}

	path := make([]h3Index, size)
	result := gridPathCells(path, start, end)
	if result != ePentagon {
		t.Errorf("Expected ePentagon for line not computable due to pentagon distortion, got %v", result)
	}
}
