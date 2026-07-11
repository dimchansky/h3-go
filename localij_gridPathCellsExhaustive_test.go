// Tests ported from testGridPathCellsExhaustive.c
package h3

import (
	"testing"
)

// Maximum distances for each resolution, mirroring C constant.
var maxDistances = []int32{1, 2, 5, 12, 19, 26}

// Helper function to iterate all indexes at a given resolution (reused from other tests).
func iterateAllIndexesAtResForGridPath(t *testing.T, res int32, testFunc func(t *testing.T, h3 h3Index)) {
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

// Helper function to iterate partial indexes at a given resolution (limit to first N base cells).
func iterateAllIndexesAtResPartial(t *testing.T, res int32, testFunc func(t *testing.T, h3 h3Index), maxBaseCells int32) {
	t.Helper()

	// Get all base cells
	baseCells := make([]h3Index, numBaseCells)
	if err := getRes0Cells(baseCells); err != eSuccess {
		t.Fatalf("Failed to get res 0 cells: %v", err)
	}

	// Limit to the specified number of base cells
	limit := maxBaseCells
	if int32(len(baseCells)) < limit {
		limit = int32(len(baseCells))
	}

	if res == 0 {
		// For resolution 0, just test the base cells
		for i := int32(0); i < limit; i++ {
			testFunc(t, baseCells[i])
		}
		return
	}

	// For higher resolutions, get children of each base cell (up to limit)
	for i := int32(0); i < limit; i++ {
		baseCell := baseCells[i]
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

// Property-based testing of gridPathCells output (ported from gridPathCells_assertions).
func gridPathCells_assertions(t *testing.T, start, end h3Index) {
	t.Helper()

	var sz int64
	err := gridPathCellsSize(start, end, &sz)
	if err != eSuccess {
		t.Fatalf("gridPathCellsSize failed: %v", err)
	}
	if sz <= 0 {
		t.Fatalf("Expected valid size > 0, got %d", sz)
	}

	line := make([]h3Index, sz)
	err = gridPathCells(line, start, end)
	if err != eSuccess {
		t.Fatalf("gridPathCells failed: %v", err)
	}

	if line[0] != start {
		t.Errorf("Line should start with start index, got %x expected %x", line[0], start)
	}
	if line[sz-1] != end {
		t.Errorf("Line should end with end index, got %x expected %x", line[sz-1], end)
	}

	for i := int64(1); i < sz; i++ {
		if !isValidCell(line[i]) {
			t.Errorf("Index %d (%x) is not valid", i, line[i])
		}

		isNeighbor, err := areNeighborCells(line[i], line[i-1])
		if err != eSuccess {
			t.Errorf("areNeighborCells failed: %v", err)
		}
		if !isNeighbor {
			t.Errorf("Index %d (%x) should be a neighbor of previous index (%x)", i, line[i], line[i-1])
		}

		if i > 1 {
			isNeighbor, err := areNeighborCells(line[i], line[i-2])
			if err != eSuccess {
				t.Errorf("areNeighborCells failed: %v", err)
			}
			if isNeighbor {
				t.Errorf("Index %d (%x) should not be a neighbor of index before previous (%x)", i, line[i], line[i-2])
			}
		}
	}
}

// Tests for invalid gridPathCells input (ported from gridPathCells_invalid_assertions).
func gridPathCells_invalid_assertions(t *testing.T, start, end h3Index) {
	t.Helper()

	var sz int64
	err := gridPathCellsSize(start, end, &sz)
	if err == eSuccess {
		t.Errorf("Line size should be marked as invalid, but got success")
	}

	line := make([]h3Index, 1) // Small buffer, doesn't matter since it should fail
	err = gridPathCells(line, start, end)
	if err == eSuccess {
		t.Errorf("Line should be marked as invalid, but got success")
	}
}

// Test for lines from an index to all neighbors within a gridDisk (ported from gridPathCells_gridDisk_assertions).
func gridPathCells_gridDisk_assertions(t *testing.T, h3 h3Index) {
	t.Helper()

	r := getResolution(h3)
	if r > 5 {
		t.Skipf("Resolution %d not supported by test function (gridDisk)", r)
	}
	maxK := maxDistances[r]

	var sz int64
	err := maxGridDiskSize(maxK, &sz)
	if err != eSuccess {
		t.Fatalf("maxGridDiskSize failed: %v", err)
	}

	if isPentagon(h3) {
		return // Skip pentagons as in C code
	}

	neighbors := make([]h3Index, sz)
	err = gridDisk(h3, maxK, neighbors)
	if err != eSuccess {
		t.Fatalf("gridDisk failed: %v", err)
	}

	for i := int64(0); i < sz; i++ {
		if neighbors[i] == 0 {
			continue
		}

		var distance int64
		distanceError := gridDistance(h3, neighbors[i], &distance)
		if distanceError == eSuccess {
			gridPathCells_assertions(t, h3, neighbors[i])
		} else {
			gridPathCells_invalid_assertions(t, h3, neighbors[i])
		}
	}
}

// Main exhaustive test function (ported from gridPathCells_gridDisk test).
func Test_gridPathCells_gridDisk(t *testing.T) {
	t.Parallel()

	// Test all indexes at resolutions 0, 1, 2
	iterateAllIndexesAtResForGridPath(t, 0, gridPathCells_gridDisk_assertions)
	iterateAllIndexesAtResForGridPath(t, 1, gridPathCells_gridDisk_assertions)
	iterateAllIndexesAtResForGridPath(t, 2, gridPathCells_gridDisk_assertions)

	// Don't iterate all of res 3, to save time (use partial)
	iterateAllIndexesAtResPartial(t, 3, gridPathCells_gridDisk_assertions, 6)
	// Further resolutions aren't tested to save time.
}
