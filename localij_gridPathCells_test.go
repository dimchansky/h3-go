// Tests ported from H3 v4.4.0: src/apps/testapps/testGridPathCells.c.
// v4.5.0 delta incorporated: gridPathCells_pentagon was replaced
// upstream by gridPathCells_pentagonReverseInterpolation (the pair now
// succeeds via the end-anchored retry), the pinned
// knownFailureNotCoveredByReverseInterpolation pair was added, and the
// suite gained the shared assertPathValid helper.
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

// assertPathValid mirrors the 4.5.0 suite helper: positive size, exact
// endpoints, and areNeighborCells contiguity for every consecutive pair.
func assertPathValid(t *testing.T, start, end h3Index, path []h3Index) {
	t.Helper()
	if len(path) == 0 {
		t.Fatal("path size is positive")
	}
	if path[0] != start {
		t.Error("path starts with start index")
	}
	if path[len(path)-1] != end {
		t.Error("path ends with end index")
	}

	for i := 1; i < len(path); i++ {
		isNeighbor, err := areNeighborCells(path[i], path[i-1])
		if err != eSuccess {
			t.Fatalf("areNeighborCells: %v", err)
		}
		if !isNeighbor {
			t.Error("path is contiguous")
		}
	}
}

func Test_gridPathCells_pentagonReverseInterpolation(t *testing.T) {
	t.Parallel()

	start := h3Index(0x820807fffffffff)
	end := h3Index(0x8208e7fffffffff)

	var size int64
	if err := gridPathCellsSize(start, end, &size); err != eSuccess {
		t.Fatalf("gridPathCellsSize failed: %v", err)
	}

	path := make([]h3Index, size)
	if err := gridPathCells(path, start, end); err != eSuccess {
		t.Fatalf("gridPathCells failed: %v", err)
	}
	assertPathValid(t, start, end, path)
}

func Test_gridPathCells_knownFailureNotCoveredByReverseInterpolation(t *testing.T) {
	t.Parallel()

	// Known limitation case: there are still pairs where gridDistance
	// succeeds but interpolation fails in both origin-anchored local IJK
	// charts (anchored at start and anchored at end). Since gridPathCells
	// only performs these two interpolation attempts, it returns an
	// error. This pinned pair demonstrates the current approach is not
	// complete.
	start := h3Index(0x8411b61ffffffff)
	end := h3Index(0x84016d3ffffffff)

	var size int64
	if err := gridPathCellsSize(start, end, &size); err != eSuccess {
		t.Fatalf("gridPathCellsSize failed: %v", err)
	}

	path := make([]h3Index, size)
	if err := gridPathCells(path, start, end); err == eSuccess {
		t.Error("Expected gridPathCells to fail")
	}
}
