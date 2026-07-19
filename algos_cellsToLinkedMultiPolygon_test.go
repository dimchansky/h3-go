// Tests ported from H3 v4.4.0: src/apps/testapps/testCellsToLinkedMultiPolygon.c.
// v4.5.0 delta incorporated: the nested-donut cases assert
// order-independently via findLinkedPolygonByOuterCount (polygon order
// is an implementation artifact of the new arc-based algorithm),
// specificLeak expects E_CELL_INVALID (validateCellSet, record §7.1),
// and globalEquatorCells is new.
package h3

import (
	"testing"
)

// findLinkedPolygonByOuterCount searches a LinkedGeoPolygon and returns
// the first polygon whose outer loop has exactly n coords.
//
// Used for order-independent testing of LinkedGeoPolygon output.
// Returns the first match, so results are ambiguous if multiple
// polygons share the size outer loop.
func findLinkedPolygonByOuterCount(root *linkedGeoPolygon, n int32) *linkedGeoPolygon {
	poly := root
	for poly != nil {
		if poly.First != nil && countLinkedCoords(poly.First) == n {
			return poly
		}
		poly = poly.Next
	}
	return nil
}

func TestCellsToLinkedMultiPolygon_empty(t *testing.T) {
	t.Parallel()

	var polygon linkedGeoPolygon

	err := cellsToLinkedMultiPolygon(nil, 0, &polygon)
	if err != eSuccess {
		t.Fatalf("cellsToLinkedMultiPolygon with empty set failed: %v", err)
	}

	if countLinkedLoops(&polygon) != 0 {
		t.Errorf("Expected 0 loops, got %d", countLinkedLoops(&polygon))
	}

	destroyLinkedMultiPolygon(&polygon)
}

func TestCellsToLinkedMultiPolygon_singleHex(t *testing.T) {
	t.Parallel()

	var polygon linkedGeoPolygon
	set := []h3Index{0x890dab6220bffff}

	err := cellsToLinkedMultiPolygon(set, int32(len(set)), &polygon)
	if err != eSuccess {
		t.Fatalf("cellsToLinkedMultiPolygon with single hex failed: %v", err)
	}

	if countLinkedLoops(&polygon) != 1 {
		t.Errorf("Expected 1 loop, got %d", countLinkedLoops(&polygon))
	}

	if countLinkedCoords(polygon.First) != 6 {
		t.Errorf("Expected 6 coords, got %d", countLinkedCoords(polygon.First))
	}

	destroyLinkedMultiPolygon(&polygon)
}

func TestCellsToLinkedMultiPolygon_invalid(t *testing.T) {
	t.Parallel()

	var polygon linkedGeoPolygon
	set := []h3Index{0xfffffffffffffff}

	err := cellsToLinkedMultiPolygon(set, int32(len(set)), &polygon)
	if err != eCellInvalid {
		t.Errorf("Expected eCellInvalid, got %v", err)
	}
}

func TestCellsToLinkedMultiPolygon_contiguous2(t *testing.T) {
	t.Parallel()

	var polygon linkedGeoPolygon
	set := []h3Index{0x8928308291bffff, 0x89283082957ffff}

	err := cellsToLinkedMultiPolygon(set, int32(len(set)), &polygon)
	if err != eSuccess {
		t.Fatalf("cellsToLinkedMultiPolygon with contiguous 2 failed: %v", err)
	}

	if countLinkedLoops(&polygon) != 1 {
		t.Errorf("Expected 1 loop, got %d", countLinkedLoops(&polygon))
	}

	if countLinkedCoords(polygon.First) != 10 {
		t.Errorf("Expected 10 coords (12 total minus 2 shared), got %d", countLinkedCoords(polygon.First))
	}

	destroyLinkedMultiPolygon(&polygon)
}

func TestCellsToLinkedMultiPolygon_nonContiguous2(t *testing.T) {
	t.Parallel()

	var polygon linkedGeoPolygon
	set := []h3Index{0x8928308291bffff, 0x89283082943ffff}

	err := cellsToLinkedMultiPolygon(set, int32(len(set)), &polygon)
	if err != eSuccess {
		t.Fatalf("cellsToLinkedMultiPolygon with non-contiguous 2 failed: %v", err)
	}

	if countLinkedPolygons(&polygon) != 2 {
		t.Errorf("Expected 2 polygons, got %d", countLinkedPolygons(&polygon))
	}

	if countLinkedLoops(&polygon) != 1 {
		t.Errorf("Expected 1 loop on first polygon, got %d", countLinkedLoops(&polygon))
	}

	if countLinkedCoords(polygon.First) != 6 {
		t.Errorf("Expected 6 coords on first polygon, got %d", countLinkedCoords(polygon.First))
	}

	if countLinkedLoops(polygon.Next) != 1 {
		t.Errorf("Expected 1 loop on second polygon, got %d", countLinkedLoops(polygon.Next))
	}

	if countLinkedCoords(polygon.Next.First) != 6 {
		t.Errorf("Expected 6 coords on second polygon, got %d", countLinkedCoords(polygon.Next.First))
	}

	destroyLinkedMultiPolygon(&polygon)
}

func TestCellsToLinkedMultiPolygon_contiguous3(t *testing.T) {
	t.Parallel()

	var polygon linkedGeoPolygon
	set := []h3Index{0x8928308288bffff, 0x892830828d7ffff, 0x8928308289bffff}

	err := cellsToLinkedMultiPolygon(set, int32(len(set)), &polygon)
	if err != eSuccess {
		t.Fatalf("cellsToLinkedMultiPolygon with contiguous 3 failed: %v", err)
	}

	if countLinkedLoops(&polygon) != 1 {
		t.Errorf("Expected 1 loop, got %d", countLinkedLoops(&polygon))
	}

	if countLinkedCoords(polygon.First) != 12 {
		t.Errorf("Expected 12 coords (18 total minus 6 shared), got %d", countLinkedCoords(polygon.First))
	}

	destroyLinkedMultiPolygon(&polygon)
}

func TestCellsToLinkedMultiPolygon_hole(t *testing.T) {
	t.Parallel()

	var polygon linkedGeoPolygon
	set := []h3Index{0x892830828c7ffff, 0x892830828d7ffff, 0x8928308289bffff,
		0x89283082813ffff, 0x8928308288fffff, 0x89283082883ffff}

	err := cellsToLinkedMultiPolygon(set, int32(len(set)), &polygon)
	if err != eSuccess {
		t.Fatalf("cellsToLinkedMultiPolygon with hole failed: %v", err)
	}

	if countLinkedLoops(&polygon) != 2 {
		t.Errorf("Expected 2 loops, got %d", countLinkedLoops(&polygon))
	}

	if countLinkedCoords(polygon.First) != 18 { // 6 * 3
		t.Errorf("Expected 18 coords on outer loop, got %d", countLinkedCoords(polygon.First))
	}

	if countLinkedCoords(polygon.First.Next) != 6 {
		t.Errorf("Expected 6 coords on inner loop, got %d", countLinkedCoords(polygon.First.Next))
	}

	destroyLinkedMultiPolygon(&polygon)
}

func TestCellsToLinkedMultiPolygon_pentagon(t *testing.T) {
	t.Parallel()

	var polygon linkedGeoPolygon
	set := []h3Index{0x851c0003fffffff}

	err := cellsToLinkedMultiPolygon(set, int32(len(set)), &polygon)
	if err != eSuccess {
		t.Fatalf("cellsToLinkedMultiPolygon with pentagon failed: %v", err)
	}

	if countLinkedLoops(&polygon) != 1 {
		t.Errorf("Expected 1 loop, got %d", countLinkedLoops(&polygon))
	}

	if countLinkedCoords(polygon.First) != 10 {
		t.Errorf("Expected 10 coords (distorted pentagon), got %d", countLinkedCoords(polygon.First))
	}

	destroyLinkedMultiPolygon(&polygon)
}

func TestCellsToLinkedMultiPolygon_twoRing(t *testing.T) {
	t.Parallel()

	var polygon linkedGeoPolygon
	// 2-ring, in order returned by k-ring algo
	set := []h3Index{
		0x8930062838bffff, 0x8930062838fffff, 0x89300628383ffff,
		0x8930062839bffff, 0x893006283d7ffff, 0x893006283c7ffff,
		0x89300628313ffff, 0x89300628317ffff, 0x893006283bbffff,
		0x89300628387ffff, 0x89300628397ffff, 0x89300628393ffff,
		0x89300628067ffff, 0x8930062806fffff, 0x893006283d3ffff,
		0x893006283c3ffff, 0x893006283cfffff, 0x8930062831bffff,
		0x89300628303ffff,
	}

	err := cellsToLinkedMultiPolygon(set, int32(len(set)), &polygon)
	if err != eSuccess {
		t.Fatalf("cellsToLinkedMultiPolygon with two ring failed: %v", err)
	}

	if countLinkedLoops(&polygon) != 1 {
		t.Errorf("Expected 1 loop, got %d", countLinkedLoops(&polygon))
	}

	expectedCoords := 6 * (2*2 + 1) // 6 * 5 = 30
	if countLinkedCoords(polygon.First) != int32(expectedCoords) {
		t.Errorf("Expected %d coords, got %d", expectedCoords, countLinkedCoords(polygon.First))
	}

	destroyLinkedMultiPolygon(&polygon)
}

func TestCellsToLinkedMultiPolygon_twoRingUnordered(t *testing.T) {
	t.Parallel()

	var polygon linkedGeoPolygon
	// 2-ring in random order
	set := []h3Index{
		0x89300628393ffff, 0x89300628383ffff, 0x89300628397ffff,
		0x89300628067ffff, 0x89300628387ffff, 0x893006283bbffff,
		0x89300628313ffff, 0x893006283cfffff, 0x89300628303ffff,
		0x89300628317ffff, 0x8930062839bffff, 0x8930062838bffff,
		0x8930062806fffff, 0x8930062838fffff, 0x893006283d3ffff,
		0x893006283c3ffff, 0x8930062831bffff, 0x893006283d7ffff,
		0x893006283c7ffff,
	}

	err := cellsToLinkedMultiPolygon(set, int32(len(set)), &polygon)
	if err != eSuccess {
		t.Fatalf("cellsToLinkedMultiPolygon with two ring unordered failed: %v", err)
	}

	if countLinkedLoops(&polygon) != 1 {
		t.Errorf("Expected 1 loop, got %d", countLinkedLoops(&polygon))
	}

	expectedCoords := 6 * (2*2 + 1) // 6 * 5 = 30
	if countLinkedCoords(polygon.First) != int32(expectedCoords) {
		t.Errorf("Expected %d coords, got %d", expectedCoords, countLinkedCoords(polygon.First))
	}

	destroyLinkedMultiPolygon(&polygon)
}

func TestCellsToLinkedMultiPolygon_nestedDonut(t *testing.T) {
	t.Parallel()

	var polygon linkedGeoPolygon
	// hollow 1-ring + hollow 3-ring around the same hex
	set := []h3Index{
		0x89283082813ffff, 0x8928308281bffff, 0x8928308280bffff,
		0x8928308280fffff, 0x89283082807ffff, 0x89283082817ffff,
		0x8928308289bffff, 0x892830828d7ffff, 0x892830828c3ffff,
		0x892830828cbffff, 0x89283082853ffff, 0x89283082843ffff,
		0x8928308284fffff, 0x8928308287bffff, 0x89283082863ffff,
		0x89283082867ffff, 0x8928308282bffff, 0x89283082823ffff,
		0x89283082837ffff, 0x892830828afffff, 0x892830828a3ffff,
		0x892830828b3ffff, 0x89283082887ffff, 0x89283082883ffff,
	}

	err := cellsToLinkedMultiPolygon(set, int32(len(set)), &polygon)
	if err != eSuccess {
		t.Fatalf("cellsToLinkedMultiPolygon with nested donut failed: %v", err)
	}

	if countLinkedPolygons(&polygon) != 2 {
		t.Errorf("Expected 2 polygons, got %d", countLinkedPolygons(&polygon))
	}

	big := findLinkedPolygonByOuterCount(&polygon, 42)
	if big == nil {
		t.Fatal("Found no big polygon by outer count")
	}
	if countLinkedLoops(big) != 2 {
		t.Errorf("Expected 2 loops on big polygon, got %d", countLinkedLoops(big))
	}
	if countLinkedCoords(big.First.Next) != 30 {
		t.Errorf("Expected 30 coords on big inner loop, got %d", countLinkedCoords(big.First.Next))
	}

	small := findLinkedPolygonByOuterCount(&polygon, 18)
	if small == nil {
		t.Fatal("Found no small polygon by outer count")
	}
	if countLinkedLoops(small) != 2 {
		t.Errorf("Expected 2 loops on small polygon, got %d", countLinkedLoops(small))
	}
	if countLinkedCoords(small.First.Next) != 6 {
		t.Errorf("Expected 6 coords on small inner loop, got %d", countLinkedCoords(small.First.Next))
	}

	destroyLinkedMultiPolygon(&polygon)
}

func TestCellsToLinkedMultiPolygon_nestedDonutTransmeridian(t *testing.T) {
	t.Parallel()

	var polygon linkedGeoPolygon
	// hollow 1-ring + hollow 3-ring around the hex at (0, -180)
	set := []h3Index{
		0x897eb5722c7ffff, 0x897eb5722cfffff, 0x897eb572257ffff,
		0x897eb57220bffff, 0x897eb572203ffff, 0x897eb572213ffff,
		0x897eb57266fffff, 0x897eb5722d3ffff, 0x897eb5722dbffff,
		0x897eb573537ffff, 0x897eb573527ffff, 0x897eb57225bffff,
		0x897eb57224bffff, 0x897eb57224fffff, 0x897eb57227bffff,
		0x897eb572263ffff, 0x897eb572277ffff, 0x897eb57223bffff,
		0x897eb572233ffff, 0x897eb5722abffff, 0x897eb5722bbffff,
		0x897eb572287ffff, 0x897eb572283ffff, 0x897eb57229bffff,
	}

	err := cellsToLinkedMultiPolygon(set, int32(len(set)), &polygon)
	if err != eSuccess {
		t.Fatalf("cellsToLinkedMultiPolygon with nested donut transmeridian failed: %v", err)
	}

	if countLinkedPolygons(&polygon) != 2 {
		t.Errorf("Expected 2 polygons, got %d", countLinkedPolygons(&polygon))
	}

	big := findLinkedPolygonByOuterCount(&polygon, 42)
	if big == nil {
		t.Fatal("Found no big polygon by outer count")
	}
	if countLinkedLoops(big) != 2 {
		t.Errorf("Expected 2 loops on big polygon, got %d", countLinkedLoops(big))
	}
	if countLinkedCoords(big.First.Next) != 30 {
		t.Errorf("Expected 30 coords on big inner loop, got %d", countLinkedCoords(big.First.Next))
	}

	small := findLinkedPolygonByOuterCount(&polygon, 18)
	if small == nil {
		t.Fatal("Found no small polygon by outer count")
	}
	if countLinkedLoops(small) != 2 {
		t.Errorf("Expected 2 loops on small polygon, got %d", countLinkedLoops(small))
	}
	if countLinkedCoords(small.First.Next) != 6 {
		t.Errorf("Expected 6 coords on small inner loop, got %d", countLinkedCoords(small.First.Next))
	}

	destroyLinkedMultiPolygon(&polygon)
}

func TestCellsToLinkedMultiPolygon_contiguous2distorted(t *testing.T) {
	t.Parallel()

	var polygon linkedGeoPolygon
	set := []h3Index{0x894cc5365afffff, 0x894cc536537ffff}

	err := cellsToLinkedMultiPolygon(set, int32(len(set)), &polygon)
	if err != eSuccess {
		t.Fatalf("cellsToLinkedMultiPolygon with contiguous 2 distorted failed: %v", err)
	}

	if countLinkedLoops(&polygon) != 1 {
		t.Errorf("Expected 1 loop, got %d", countLinkedLoops(&polygon))
	}

	if countLinkedCoords(polygon.First) != 12 {
		t.Errorf("Expected 12 coords (14 total minus 2 shared), got %d", countLinkedCoords(polygon.First))
	}

	destroyLinkedMultiPolygon(&polygon)
}

func TestCellsToLinkedMultiPolygon_negativeHashedCoordinates(t *testing.T) {
	t.Parallel()

	var polygon linkedGeoPolygon
	set := []h3Index{0x88ad36c547fffff, 0x88ad36c467fffff}

	err := cellsToLinkedMultiPolygon(set, int32(len(set)), &polygon)
	if err != eSuccess {
		t.Fatalf("cellsToLinkedMultiPolygon with negative hashed coordinates failed: %v", err)
	}

	if countLinkedPolygons(&polygon) != 2 {
		t.Errorf("Expected 2 polygons, got %d", countLinkedPolygons(&polygon))
	}

	if countLinkedLoops(&polygon) != 1 {
		t.Errorf("Expected 1 loop on first polygon, got %d", countLinkedLoops(&polygon))
	}

	if countLinkedCoords(polygon.First) != 6 {
		t.Errorf("Expected 6 coords on first polygon, got %d", countLinkedCoords(polygon.First))
	}

	if countLinkedLoops(polygon.Next) != 1 {
		t.Errorf("Expected 1 loop on second polygon, got %d", countLinkedLoops(polygon.Next))
	}

	if countLinkedCoords(polygon.Next.First) != 6 {
		t.Errorf("Expected 6 coords on second polygon, got %d", countLinkedCoords(polygon.Next.First))
	}

	destroyLinkedMultiPolygon(&polygon)
}

func TestCellsToLinkedMultiPolygon_specificLeak(t *testing.T) {
	t.Parallel()

	// Test for a case where a leak can occur, detected by fuzzer.
	// The leak detection part will be enforced here by the Go garbage collector.
	var polygon linkedGeoPolygon
	set := []h3Index{0xd60006d60000f100, 0x3c3c403c1300d668}

	err := cellsToLinkedMultiPolygon(set, int32(len(set)), &polygon)
	if err != eCellInvalid {
		t.Errorf("Expected eCellInvalid for invalid cells, got %v", err)
	}
}

func TestCellsToLinkedMultiPolygon_globalEquatorCells(t *testing.T) {
	t.Parallel()

	// Global polygon: a continuous band of res-1 cells around the equator.
	cells := []h3Index{
		0x81807ffffffffff, 0x817efffffffffff, 0x81723ffffffffff,
		0x817ebffffffffff, 0x817c3ffffffffff, 0x817e3ffffffffff,
		0x817a3ffffffffff, 0x8166fffffffffff, 0x8172bffffffffff,
		0x816afffffffffff, 0x81933ffffffffff, 0x8168fffffffffff,
		0x8188fffffffffff, 0x81853ffffffffff, 0x817f7ffffffffff,
		0x8180bffffffffff, 0x81783ffffffffff, 0x81743ffffffffff,
		0x8170bffffffffff, 0x8173bffffffffff, 0x8179bffffffffff,
		0x817cbffffffffff, 0x8188bffffffffff, 0x81857ffffffffff,
		0x816f7ffffffffff, 0x8177bffffffffff, 0x81617ffffffffff,
		0x816f3ffffffffff, 0x8174bffffffffff, 0x8180fffffffffff,
		0x817a7ffffffffff, 0x81767ffffffffff, 0x81757ffffffffff,
		0x81957ffffffffff, 0x81787ffffffffff, 0x81847ffffffffff,
		0x81653ffffffffff, 0x817bbffffffffff, 0x816cfffffffffff,
		0x816abffffffffff, 0x815f3ffffffffff, 0x817c7ffffffffff,
		0x8168bffffffffff, 0x818cbffffffffff, 0x818cfffffffffff,
		0x818afffffffffff, 0x8174fffffffffff, 0x8172fffffffffff,
		0x8170fffffffffff, 0x816fbffffffffff, 0x81657ffffffffff,
		0x816c7ffffffffff, 0x8186bffffffffff, 0x81763ffffffffff,
		0x818a7ffffffffff, 0x8186fffffffffff, 0x81707ffffffffff,
		0x8182bffffffffff, 0x818f3ffffffffff, 0x8182fffffffffff,
	}

	var polygon linkedGeoPolygon
	if err := cellsToLinkedMultiPolygon(cells, int32(len(cells)), &polygon); err != eSuccess {
		t.Fatalf("cellsToLinkedMultiPolygon with global equator cells failed: %v", err)
	}
	destroyLinkedMultiPolygon(&polygon)
}

func TestCellsToLinkedMultiPolygon_gridDiskResolutions(t *testing.T) {
	t.Parallel()

	// This is a center-face base cell, no pentagon siblings
	baseCell := h3Index(0x8073fffffffffff)
	var origin h3Index

	indexes := make([]h3Index, 19)

	for res := int32(1); res < 15; res++ {
		// Take the 2-disk of the center child at res
		var err h3Error
		origin, err = cellToCenterChild(baseCell, res)
		if err != eSuccess {
			t.Fatalf("cellToCenterChild failed at res %d: %v", res, err)
		}

		err = gridDisk(origin, 2, indexes)
		if err != eSuccess {
			t.Fatalf("gridDisk failed at res %d: %v", res, err)
		}

		// Test the polygon output
		var polygon linkedGeoPolygon
		err = cellsToLinkedMultiPolygon(indexes, int32(len(indexes)), &polygon)
		if err != eSuccess {
			t.Fatalf("cellsToLinkedMultiPolygon failed at res %d: %v", res, err)
		}

		if countLinkedPolygons(&polygon) != 1 {
			t.Errorf("Expected 1 polygon at res %d, got %d", res, countLinkedPolygons(&polygon))
		}

		if countLinkedLoops(&polygon) != 1 {
			t.Errorf("Expected 1 loop at res %d, got %d", res, countLinkedLoops(&polygon))
		}

		if countLinkedCoords(polygon.First) != 30 {
			t.Errorf("Expected 30 coords at res %d, got %d", res, countLinkedCoords(polygon.First))
		}

		destroyLinkedMultiPolygon(&polygon)
	}
}

func TestCellsToLinkedMultiPolygon_gridDiskResolutionsPentagon(t *testing.T) {
	t.Parallel()

	// This is a pentagon base cell
	baseCell := h3Index(0x8031fffffffffff)
	var origin h3Index

	diskIndexes := make([]h3Index, 7)
	indexes := make([]h3Index, 6)

	for res := int32(1); res < 15; res++ {
		// Take the 1-disk of the center child at res. Note: We can't take
		// the 2-disk here, as increased distortion around the pentagon will
		// still fail at res 1. TODO: Use a 2-ring, start at res 0
		// when output is correct.
		var err h3Error
		origin, err = cellToCenterChild(baseCell, res)
		if err != eSuccess {
			t.Fatalf("cellToCenterChild failed at res %d: %v", res, err)
		}

		err = gridDisk(origin, 1, diskIndexes)
		if err != eSuccess {
			t.Fatalf("gridDisk failed at res %d: %v", res, err)
		}

		// Filter out null entries and collect into indexes
		j := 0
		for i := 0; i < len(diskIndexes); i++ {
			if diskIndexes[i] != h3Null {
				indexes[j] = diskIndexes[i]
				j++
			}
		}

		if j != 6 {
			t.Fatalf("Expected 6 non-null indexes at res %d, got %d", res, j)
		}

		// Test the polygon output
		var polygon linkedGeoPolygon
		err = cellsToLinkedMultiPolygon(indexes[:6], 6, &polygon)
		if err != eSuccess {
			t.Fatalf("cellsToLinkedMultiPolygon failed at res %d: %v", res, err)
		}

		if countLinkedPolygons(&polygon) != 1 {
			t.Errorf("Expected 1 polygon at res %d, got %d", res, countLinkedPolygons(&polygon))
		}

		if countLinkedLoops(&polygon) != 1 {
			t.Errorf("Expected 1 loop at res %d, got %d", res, countLinkedLoops(&polygon))
		}

		if countLinkedCoords(polygon.First) != 15 {
			t.Errorf("Expected 15 coords at res %d, got %d", res, countLinkedCoords(polygon.First))
		}

		destroyLinkedMultiPolygon(&polygon)
	}
}
