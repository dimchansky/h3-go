// Tests ported from testGridPathCells.c
package h3

import (
	"testing"
)

func Test_gridPathCells_acrossMultipleFaces(t *testing.T) {
	t.Parallel()
	
	start := H3Index(0x85285aa7fffffff)
	end := H3Index(0x851d9b1bfffffff)

	var size int64
	lineError := gridPathCellsSize(start, end, &size)
	if lineError != E_FAILED {
		t.Errorf("Expected E_FAILED for line not computable across multiple icosa faces, got %v", lineError)
	}
}

func Test_gridPathCells_pentagon(t *testing.T) {
	t.Parallel()
	
	start := H3Index(0x820807fffffffff)
	end := H3Index(0x8208e7fffffffff)

	var size int64
	err := gridPathCellsSize(start, end, &size)
	if err != E_SUCCESS {
		t.Fatalf("gridPathCellsSize failed: %v", err)
	}
	
	path := make([]H3Index, size)
	result := gridPathCells(path, start, end)
	if result != E_PENTAGON {
		t.Errorf("Expected E_PENTAGON for line not computable due to pentagon distortion, got %v", result)
	}
}