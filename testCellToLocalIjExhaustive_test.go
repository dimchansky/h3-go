// Tests ported from testCellToLocalIjExhaustive.c
package h3

import (
	"testing"
)

// Maximum distances for each resolution, mirroring C constant MAX_DISTANCES.
var maxDistancesForCellToLocalIj = []int32{1, 2, 5, 12, 19, 26}

// Traversal constants from C implementation algosDirections.
var directionsForCellToLocalIj = []CoordIJ{
	{I: 0, J: 1}, {I: -1, J: 0}, {I: -1, J: -1},
	{I: 0, J: -1}, {I: 1, J: 0}, {I: 1, J: 1},
}

// Next ring direction from C implementation nextRingDirection.
var nextRingDirectionForCellToLocalIj = CoordIJ{I: 1, J: 0}

// Helper function to iterate all indexes at a given resolution.
func iterateAllIndexesAtResForCellToLocalIj(t *testing.T, res int32, testFunc func(t *testing.T, h3 h3Index)) {
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
func iterateAllIndexesAtResPartialForCellToLocalIj(t *testing.T, res int32, testFunc func(t *testing.T, h3 h3Index), maxBaseCells int32) {
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

// Test that the local coordinates for an index map to itself (ported from localIjToH3_identity_assertions).
func localIjToH3_identity_assertions(t *testing.T, h3 h3Index) {
	t.Helper()

	var ij CoordIJ
	err := cellToLocalIj(h3, h3, 0, &ij)
	if err != eSuccess {
		t.Errorf("able to setup localIjToH3 test for %#016x, got error: %v", h3, err)
		return
	}

	var retrieved h3Index
	err = localIjToCell(h3, &ij, 0, &retrieved)
	if err != eSuccess {
		t.Errorf("got an index back from localIjToH3 for %#016x, got error: %v", h3, err)
		return
	}

	if h3 != retrieved {
		t.Errorf("round trip through local quadIJ space failed for %#016x: got %#016x", h3, retrieved)
	}
}

// Test that coordinates for an index match some simple rules about index digits (ported from h3ToLocalIj_coordinates_assertions).
func h3ToLocalIj_coordinates_assertions(t *testing.T, h3 h3Index) {
	t.Helper()

	r := getResolution(h3)

	var ij CoordIJ
	err := cellToLocalIj(h3, h3, 0, &ij)
	if err != eSuccess {
		t.Errorf("get ij for origin %#016x, got error: %v", h3, err)
		return
	}

	var ijk coordIJK
	err = ijToIjk(&ij, &ijk)
	if err != eSuccess {
		t.Errorf("ijToIjk failed for %#016x, got error: %v", h3, err)
		return
	}

	switch r {
	case 0:
		if !_ijkMatches(&ijk, &unitVecs[0]) {
			t.Errorf("res 0 cell at 0,0,0 failed for %#016x: got ijk=%+v", h3, ijk)
		}
	case 1:
		expected := unitVecs[getIndexDigit(h3, 1)]
		if !_ijkMatches(&ijk, &expected) {
			t.Errorf("res 1 cell at expected coordinates failed for %#016x: got ijk=%+v, expected=%+v", h3, ijk, expected)
		}
	case 2:
		expected := unitVecs[getIndexDigit(h3, 1)]
		_downAp7r(&expected)
		_neighbor(&expected, direction(getIndexDigit(h3, 2)))
		if !_ijkMatches(&ijk, &expected) {
			t.Errorf("res 2 cell at expected coordinates failed for %#016x: got ijk=%+v, expected=%+v", h3, ijk, expected)
		}
	default:
		t.Errorf("resolution supported by test function (coordinates) for %#016x: res=%d", h3, r)
	}
}

// Test the immediate neighbors of an index are at the expected locations in local quadIJ coordinate space (ported from h3ToLocalIj_neighbors_assertions).
func h3ToLocalIj_neighbors_assertions(t *testing.T, h3 h3Index) {
	t.Helper()

	origin := CoordIJ{I: 0, J: 0}
	err := cellToLocalIj(h3, h3, 0, &origin)
	if err != eSuccess {
		t.Errorf("got ij for origin %#016x, got error: %v", h3, err)
		return
	}

	var originIjk coordIJK
	err = ijToIjk(&origin, &originIjk)
	if err != eSuccess {
		t.Errorf("ijToIjk failed for origin %#016x, got error: %v", h3, err)
		return
	}

	for d := kAxesDigit; d < invalidDigit; d++ {
		if d == kAxesDigit && isPentagon(h3) {
			continue
		}

		var rotations int32
		var offset h3Index
		err = h3NeighborRotations(h3, d, &rotations, &offset)
		if err != eSuccess {
			continue // Some neighbors may not be valid
		}

		ij := CoordIJ{I: 0, J: 0}
		err = cellToLocalIj(h3, offset, 0, &ij)
		if err != eSuccess {
			t.Errorf("got ij for destination %#016x from origin %#016x, got error: %v", offset, h3, err)
			return
		}

		var ijk coordIJK
		err = ijToIjk(&ij, &ijk)
		if err != eSuccess {
			t.Errorf("ijToIjk failed for destination %#016x from origin %#016x, got error: %v", offset, h3, err)
			return
		}

		invertedIjk := coordIJK{I: 0, J: 0, K: 0}
		_neighbor(&invertedIjk, d)
		for i := 0; i < 3; i++ {
			_ijkRotate60ccw(&invertedIjk)
		}
		_ijkAdd(&invertedIjk, &ijk, &ijk)
		_ijkNormalize(&ijk)

		if !_ijkMatches(&ijk, &originIjk) {
			t.Errorf("back to origin failed for %#016x direction %d: got ijk=%+v, expected origin=%+v", h3, d, ijk, originIjk)
		}
	}
}

// Test the immediate neighbors of an index with invalid digits return error (ported from h3ToLocalIj_invalid_assertions).
func h3ToLocalIj_invalid_assertions(t *testing.T, h3 h3Index) {
	t.Helper()

	r := getResolution(h3)
	if r == 0 {
		t.Errorf("resolution supported by test function (invalid digits) for %#016x: res=%d", h3, r)
		return
	}
	if r > 5 {
		t.Errorf("resolution supported by test function (invalid digits) for %#016x: res=%d", h3, r)
		return
	}

	maxK := maxDistancesForCellToLocalIj[r]

	var sz int64
	err := maxGridDiskSize(maxK, &sz)
	if err != eSuccess {
		t.Errorf("maxGridDiskSize failed for %#016x k=%d, got error: %v", h3, maxK, err)
		return
	}

	neighbors := make([]h3Index, sz)
	distances := make([]int32, sz)

	err = gridDiskDistances(h3, maxK, neighbors, distances)
	if err != eSuccess {
		t.Errorf("gridDiskDistances failed for %#016x k=%d, got error: %v", h3, maxK, err)
		return
	}

	for i := int64(0); i < sz; i++ {
		if neighbors[i] == h3Null {
			continue
		}

		var ij CoordIJ
		// Don't consider indexes which we can't unfold in the first place
		if cellToLocalIj(h3, neighbors[i], 0, &ij) == eSuccess {
			for j := 0; j < 2; j++ {
				var dir direction
				if j == 0 {
					dir = invalidDigit
				} else {
					dir = kAxesDigit
				}
				// Valgrind / ASAN / UBSAN are used to test these assertions
				h3Invalid := h3
				h3Invalid = setIndexDigit(h3Invalid, 0, int32(dir))
				var ij2 CoordIJ
				cellToLocalIj(h3Invalid, neighbors[i], 0, &ij2) // Should not crash
				neighborInvalid := neighbors[i]
				neighborInvalid = setIndexDigit(neighborInvalid, 0, int32(dir))
				cellToLocalIj(h3, neighborInvalid, 0, &ij2) // Should not crash
				var out h3Index
				localIjToCell(h3Invalid, &ij, 0, &out) // Should not crash
			}
		}
	}
}

// Test that neighbors (k-ring) can be converted back to indexes in local quadIJ coordinate space (ported from localIjToH3_gridDisk_assertions).
func localIjToH3_gridDisk_assertions(t *testing.T, h3 h3Index) {
	t.Helper()

	r := getResolution(h3)
	if r > 5 {
		t.Errorf("resolution supported by test function (gridDisk) for %#016x: res=%d", h3, r)
		return
	}

	maxK := maxDistancesForCellToLocalIj[r]

	var sz int64
	err := maxGridDiskSize(maxK, &sz)
	if err != eSuccess {
		t.Errorf("maxGridDiskSize failed for %#016x k=%d, got error: %v", h3, maxK, err)
		return
	}

	neighbors := make([]h3Index, sz)
	distances := make([]int32, sz)

	err = gridDiskDistances(h3, maxK, neighbors, distances)
	if err != eSuccess {
		t.Errorf("gridDiskDistances failed for %#016x k=%d, got error: %v", h3, maxK, err)
		return
	}

	for i := int64(0); i < sz; i++ {
		if neighbors[i] == h3Null {
			continue
		}

		var ij CoordIJ
		// Don't consider indexes which we can't unfold in the first place
		if cellToLocalIj(h3, neighbors[i], 0, &ij) == eSuccess {
			var retrieved h3Index
			err = localIjToCell(h3, &ij, 0, &retrieved)
			if err != eSuccess {
				t.Errorf("retrieved index for unfolded coordinates failed for %#016x neighbor %#016x, got error: %v", h3, neighbors[i], err)
				return
			}
			if retrieved != neighbors[i] {
				t.Errorf("round trip neighboring index mismatch for %#016x: expected %#016x, got %#016x", h3, neighbors[i], retrieved)
			}
		}
	}
}

// Test traversing the local quadIJ coordinate space (ported from localIjToH3_traverse_assertions).
func localIjToH3_traverse_assertions(t *testing.T, h3 h3Index) {
	t.Helper()

	r := getResolution(h3)
	if r > 5 {
		t.Errorf("resolution supported by test function (traverse) for %#016x: res=%d", h3, r)
		return
	}

	k := maxDistancesForCellToLocalIj[r]

	var ij CoordIJ
	err := cellToLocalIj(h3, h3, 0, &ij)
	if err != eSuccess {
		t.Errorf("Got origin coordinates for %#016x, got error: %v", h3, err)
		return
	}

	// This logic is from gridDiskDistancesUnsafe.
	// 0 < ring <= k, current ring
	ring := int32(1)
	// 0 <= direction < 6, current side of the ring
	direction := int32(0)
	// 0 <= i < ring, current position on the side of the ring
	i := int32(0)

	for ring <= k {
		if direction == 0 && i == 0 {
			ij.I += nextRingDirectionForCellToLocalIj.I
			ij.J += nextRingDirectionForCellToLocalIj.J
		}

		ij.I += directionsForCellToLocalIj[direction].I
		ij.J += directionsForCellToLocalIj[direction].J

		var testH3 h3Index
		failed := localIjToCell(h3, &ij, 0, &testH3)
		if failed == eSuccess {
			if !isValidCell(testH3) {
				t.Errorf("test coordinates result in valid index failed for %#016x at ij={%d,%d}: got invalid cell %#016x", h3, ij.I, ij.J, testH3)
			}

			var expectedIj CoordIJ
			reverseFailed := cellToLocalIj(h3, testH3, 0, &expectedIj)
			// If it doesn't give a coordinate for this origin,index pair that's OK.
			if reverseFailed == eSuccess {
				if expectedIj.I != ij.I || expectedIj.J != ij.J {
					// Multiple coordinates for the same index can happen due to
					// pentagon distortion. In that case, the other coordinates
					// should also belong to the same index.
					var testTestH3 h3Index
					err = localIjToCell(h3, &expectedIj, 0, &testTestH3)
					if err != eSuccess {
						t.Errorf("converted coordinates again failed for %#016x at ij={%d,%d}, got error: %v", h3, expectedIj.I, expectedIj.J, err)
						return
					}
					if testH3 != testTestH3 {
						t.Errorf("index has normalizable coordinates in local quadIJ failed for %#016x: testH3=%#016x != testTestH3=%#016x", h3, testH3, testTestH3)
					}
				}
			}
		}

		i++
		// Check if end of this side of the k-ring
		if i == ring {
			i = 0
			direction++
			// Check if end of this ring.
			if direction == 6 {
				direction = 0
				ring++
			}
		}
	}
}

// Test localIjToH3_identity (ported from C test).
func TestLocalIjToH3_Identity(t *testing.T) {
	t.Parallel()
	iterateAllIndexesAtResForCellToLocalIj(t, 0, localIjToH3_identity_assertions)
	iterateAllIndexesAtResForCellToLocalIj(t, 1, localIjToH3_identity_assertions)
	iterateAllIndexesAtResForCellToLocalIj(t, 2, localIjToH3_identity_assertions)
}

// Test h3ToLocalIj_coordinates (ported from C test).
func TestH3ToLocalIj_Coordinates(t *testing.T) {
	t.Parallel()
	iterateAllIndexesAtResForCellToLocalIj(t, 0, h3ToLocalIj_coordinates_assertions)
	iterateAllIndexesAtResForCellToLocalIj(t, 1, h3ToLocalIj_coordinates_assertions)
	iterateAllIndexesAtResForCellToLocalIj(t, 2, h3ToLocalIj_coordinates_assertions)
}

// Test h3ToLocalIj_neighbors (ported from C test).
func TestH3ToLocalIj_Neighbors(t *testing.T) {
	t.Parallel()
	iterateAllIndexesAtResForCellToLocalIj(t, 0, h3ToLocalIj_neighbors_assertions)
	iterateAllIndexesAtResForCellToLocalIj(t, 1, h3ToLocalIj_neighbors_assertions)
	iterateAllIndexesAtResForCellToLocalIj(t, 2, h3ToLocalIj_neighbors_assertions)
}

// Test h3ToLocalIj_invalid (ported from C test).
func TestH3ToLocalIj_Invalid(t *testing.T) {
	t.Parallel()
	iterateAllIndexesAtResForCellToLocalIj(t, 1, h3ToLocalIj_invalid_assertions)
	iterateAllIndexesAtResForCellToLocalIj(t, 2, h3ToLocalIj_invalid_assertions)
}

// Test localIjToH3_gridDisk (ported from C test).
func TestLocalIjToH3_GridDisk(t *testing.T) {
	t.Parallel()
	iterateAllIndexesAtResForCellToLocalIj(t, 0, localIjToH3_gridDisk_assertions)
	iterateAllIndexesAtResForCellToLocalIj(t, 1, localIjToH3_gridDisk_assertions)
	iterateAllIndexesAtResForCellToLocalIj(t, 2, localIjToH3_gridDisk_assertions)
	// Don't iterate all of res 3, to save time
	iterateAllIndexesAtResPartialForCellToLocalIj(t, 3, localIjToH3_gridDisk_assertions, 27)
	// Further resolutions aren't tested to save time.
}

// Test localIjToH3_traverse (ported from C test).
func TestLocalIjToH3_Traverse(t *testing.T) {
	t.Parallel()
	iterateAllIndexesAtResForCellToLocalIj(t, 0, localIjToH3_traverse_assertions)
	iterateAllIndexesAtResForCellToLocalIj(t, 1, localIjToH3_traverse_assertions)
	iterateAllIndexesAtResForCellToLocalIj(t, 2, localIjToH3_traverse_assertions)
	// Don't iterate all of res 3, to save time
	iterateAllIndexesAtResPartialForCellToLocalIj(t, 3, localIjToH3_traverse_assertions, 27)
	// Further resolutions aren't tested to save time.
}
