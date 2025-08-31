// Tests ported from testGridDistanceExhaustive.c
package h3

import (
	"testing"
)

// Maximum distances for each resolution, mirroring C constant
var maxDistancesForGridDistance = []int32{1, 2, 5, 12, 19, 26}

// Helper function to iterate all indexes at a given resolution (reused from other tests)
func iterateAllIndexesAtResForGridDistance(t *testing.T, res int32, testFunc func(t *testing.T, h3 H3Index)) {
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

// Helper function to iterate partial indexes at a given resolution (limit to first N base cells)
func iterateAllIndexesAtResPartialForGridDistance(t *testing.T, res int32, testFunc func(t *testing.T, h3 H3Index), maxBaseCells int32) {
	t.Helper()

	// Get all base cells
	baseCells := make([]H3Index, NUM_BASE_CELLS)
	if err := getRes0Cells(baseCells); err != E_SUCCESS {
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

// Identity distance assertions (ported from gridDistance_identity_assertions)
func gridDistance_identity_assertions(t *testing.T, h3 H3Index) {
	t.Helper()

	var distance int64
	err := gridDistance(h3, h3, &distance)
	if err != E_SUCCESS {
		t.Fatalf("gridDistance failed: %v", err)
	}
	if distance != 0 {
		t.Errorf("Distance to self should be 0, got %d", distance)
	}
}

// Grid disk distance assertions (ported from gridDistance_gridDisk_assertions)
func gridDistance_gridDisk_assertions(t *testing.T, h3 H3Index) {
	t.Helper()

	r := getResolution(h3)
	if r > 5 {
		t.Skipf("Resolution %d not supported by test function (gridDisk)", r)
	}
	maxK := maxDistancesForGridDistance[r]

	var sz int64
	err := maxGridDiskSize(maxK, &sz)
	if err != E_SUCCESS {
		t.Fatalf("maxGridDiskSize failed: %v", err)
	}

	neighbors := make([]H3Index, sz)
	distances := make([]int32, sz)

	err = gridDiskDistances(h3, maxK, neighbors, distances)
	if err != E_SUCCESS {
		t.Fatalf("gridDiskDistances failed: %v", err)
	}

	for i := int64(0); i < sz; i++ {
		if neighbors[i] == 0 {
			continue
		}

		var calculatedDistance int64
		calculatedError := gridDistance(h3, neighbors[i], &calculatedDistance)

		// Don't consider indexes where gridDistance reports failure to generate
		if calculatedError == E_SUCCESS {
			if calculatedDistance != int64(distances[i]) {
				t.Errorf("gridDiskDistances (%d) does not match gridDistance (%d) for neighbor %x of origin %x", distances[i], calculatedDistance, neighbors[i], h3)
			}
		}
	}
}

// Main identity test (ported from gridDistance_identity test)
func Test_gridDistance_identity(t *testing.T) {
	t.Parallel()

	iterateAllIndexesAtResForGridDistance(t, 0, gridDistance_identity_assertions)
	iterateAllIndexesAtResForGridDistance(t, 1, gridDistance_identity_assertions)
	iterateAllIndexesAtResForGridDistance(t, 2, gridDistance_identity_assertions)
}

// Main grid disk test (ported from gridDistance_gridDisk test)
func Test_gridDistance_gridDisk(t *testing.T) {
	t.Parallel()

	iterateAllIndexesAtResForGridDistance(t, 0, gridDistance_gridDisk_assertions)
	iterateAllIndexesAtResForGridDistance(t, 1, gridDistance_gridDisk_assertions)
	iterateAllIndexesAtResForGridDistance(t, 2, gridDistance_gridDisk_assertions)

	// Don't iterate all of res 3, to save time
	iterateAllIndexesAtResPartialForGridDistance(t, 3, gridDistance_gridDisk_assertions, 27)
	// Further resolutions aren't tested to save time.
}