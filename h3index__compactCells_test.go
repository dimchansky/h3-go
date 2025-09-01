// Tests ported from testCompactCells.c
package h3

import (
	"testing"
)

var sunnyvale H3Index = 0x89283470c27ffff

var uncompactable = []H3Index{0x89283470803ffff, 0x8928347081bffff, 0x8928347080bffff}
var uncompactableWithZero = []H3Index{0x89283470803ffff, 0x8928347081bffff, 0, 0x8928347080bffff}

func TestCompactCells_roundtrip(t *testing.T) {
	t.Parallel()
	k := int32(9)
	var hexCount int64
	err := maxGridDiskSize(k, &hexCount)
	if err != E_SUCCESS {
		t.Fatalf("maxGridDiskSize failed: %v", err)
	}
	expectedCompactCount := 73

	// Generate a set of hexagons to compact
	sunnyvaleExpanded := make([]H3Index, hexCount)
	err = gridDisk(sunnyvale, k, sunnyvaleExpanded)
	if err != E_SUCCESS {
		t.Fatalf("gridDisk failed: %v", err)
	}

	compressed := make([]H3Index, hexCount)
	err = compactCells(sunnyvaleExpanded, compressed, hexCount)
	if err != E_SUCCESS {
		t.Fatalf("compactCells failed: %v", err)
	}

	count := 0
	for i := int64(0); i < hexCount; i++ {
		if compressed[i] != 0 {
			count++
		}
	}
	if count != expectedCompactCount {
		t.Errorf("got compacted count %d, expected %d", count, expectedCompactCount)
	}

	countUncompact, err := uncompactCellsSize(compressed, int64(count), 9)
	if err != E_SUCCESS {
		t.Fatalf("uncompactCellsSize failed: %v", err)
	}
	decompressed := make([]H3Index, countUncompact)
	err = uncompactCells(compressed, int64(count), decompressed, hexCount, 9)
	if err != E_SUCCESS {
		t.Fatalf("uncompactCells failed: %v", err)
	}

	count2 := 0
	for i := int64(0); i < hexCount; i++ {
		if decompressed[i] != 0 {
			count2++
		}
	}
	if count2 != int(hexCount) {
		t.Errorf("got uncompacted count %d, expected %d", count2, hexCount)
	}
}

func TestCompactCells_res0children(t *testing.T) {
	t.Parallel()
	var parent H3Index
	setH3Index(&parent, 0, 0, 0)

	arrSize, err := cellToChildrenSize(parent, 1)
	if err != E_SUCCESS {
		t.Fatalf("cellToChildrenSize failed: %v", err)
	}

	children := make([]H3Index, arrSize)
	err = cellToChildren(parent, 1, children)
	if err != E_SUCCESS {
		t.Fatalf("cellToChildren failed: %v", err)
	}

	compressed := make([]H3Index, arrSize)
	err = compactCells(children, compressed, arrSize)
	if err != E_SUCCESS {
		t.Fatalf("compactCells failed: %v", err)
	}
	if compressed[0] != parent {
		t.Errorf("expected parent %x, got %x", parent, compressed[0])
	}
	for idx := 1; idx < int(arrSize); idx++ {
		if compressed[idx] != 0 {
			t.Errorf("expected only 1 cell, but index %d is %x", idx, compressed[idx])
		}
	}
}

func TestCompactCells_allRes1(t *testing.T) {
	t.Parallel()
	const numRes0 = 122
	const numRes1 = 842
	cells0 := make([]H3Index, numRes0)
	cells1 := make([]H3Index, numRes1)
	out := make([]H3Index, numRes1)

	err := getRes0Cells(cells0)
	if err != E_SUCCESS {
		t.Fatalf("getRes0Cells failed: %v", err)
	}
	if cells0[0] != 0x8001fffffffffff {
		t.Errorf("got unexpected first res0 cell: %x", cells0[0])
	}

	err = uncompactCells(cells0, numRes0, cells1, numRes1, 1)
	if err != E_SUCCESS {
		t.Fatalf("uncompactCells failed: %v", err)
	}

	numUncompacted := int64(numRes1)
	err = compactCells(cells1, out, numUncompacted)
	if err != E_SUCCESS {
		t.Fatalf("compactCells failed: %v", err)
	}

	// Assert that the output of this function matches exactly the set of
	// res 0 cells
	foundCount := 0
	for res1Idx := 0; res1Idx < numRes1; res1Idx++ {
		compactedCell := out[res1Idx]

		if compactedCell != 0 {
			for res1DupIdx := 0; res1DupIdx < res1Idx; res1DupIdx++ {
				if out[res1DupIdx] == compactedCell {
					t.Errorf("Duplicated output found at indices %d and %d", res1DupIdx, res1Idx)
				}
			}

			found := false
			for res0Idx := 0; res0Idx < numRes0; res0Idx++ {
				if cells0[res0Idx] == compactedCell {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Res 0 cell %x not found in original set", compactedCell)
			}
			foundCount++
		}
	}
	if foundCount != numRes0 {
		t.Errorf("found %d res 0 cells, expected %d", foundCount, numRes0)
	}
}

func TestCompactCells_allRes1_variousRanges(t *testing.T) {
	t.Parallel()
	const numRes0 = 122
	const numRes1 = 842
	cells0 := make([]H3Index, numRes0)
	cells1 := make([]H3Index, numRes1)
	out := make([]H3Index, numRes1)

	err := getRes0Cells(cells0)
	if err != E_SUCCESS {
		t.Fatalf("getRes0Cells failed: %v", err)
	}
	if cells0[0] != 0x8001fffffffffff {
		t.Errorf("got unexpected first res0 cell: %x", cells0[0])
	}

	err = uncompactCells(cells0, numRes0, cells1, numRes1, 1)
	if err != E_SUCCESS {
		t.Fatalf("uncompactCells failed: %v", err)
	}

	// Test various (but not all possible combinations) ranges of res 1
	// cells
	for offset := int64(0); offset < numRes1; offset++ {
		for numUncompacted := numRes1 - offset; numUncompacted >= 0; numUncompacted-- {
			// Clear output array
			for i := range out {
				out[i] = 0
			}

			err = compactCells(cells1[offset:], out, numUncompacted)
			if err != E_SUCCESS {
				t.Errorf("compactCells failed at offset %d, numUncompacted %d: %v", offset, numUncompacted, err)
			}
		}
	}
}

func TestCompactCells_res0(t *testing.T) {
	t.Parallel()
	hexCount := NUM_BASE_CELLS

	res0Hexes := make([]H3Index, hexCount)
	for i := 0; i < hexCount; i++ {
		setH3Index(&res0Hexes[i], 0, int32(i), 0)
	}
	compressed := make([]H3Index, hexCount)
	err := compactCells(res0Hexes, compressed, int64(hexCount))
	if err != E_SUCCESS {
		t.Fatalf("compactCells failed: %v", err)
	}

	for i := 0; i < hexCount; i++ {
		// At resolution 0, it will be an exact copy.
		// The test is further optimizing that it will be in order (which
		// isn't guaranteed.)
		if compressed[i] != res0Hexes[i] {
			t.Errorf("got compressed[%d] = %x, expected %x", i, compressed[i], res0Hexes[i])
		}
	}

	countUncompact, err := uncompactCellsSize(compressed, int64(hexCount), 0)
	if err != E_SUCCESS {
		t.Fatalf("uncompactCellsSize failed: %v", err)
	}
	decompressed := make([]H3Index, countUncompact)
	err = uncompactCells(compressed, int64(hexCount), decompressed, int64(hexCount), 0)
	if err != E_SUCCESS {
		t.Fatalf("uncompactCells failed: %v", err)
	}

	count2 := 0
	for i := 0; i < hexCount; i++ {
		if decompressed[i] != 0 {
			count2++
		}
	}
	if count2 != hexCount {
		t.Errorf("got uncompacted count %d, expected %d", count2, hexCount)
	}
}

func TestCompactCells_uncompactable(t *testing.T) {
	t.Parallel()
	hexCount := int64(3)
	expectedCompactCount := 3

	compressed := make([]H3Index, hexCount)
	err := compactCells(uncompactable, compressed, hexCount)
	if err != E_SUCCESS {
		t.Fatalf("compactCells failed: %v", err)
	}

	count := 0
	for i := int64(0); i < hexCount; i++ {
		if compressed[i] != 0 {
			count++
		}
	}
	if count != expectedCompactCount {
		t.Errorf("got compacted count %d, expected %d", count, expectedCompactCount)
	}

	countUncompact, err := uncompactCellsSize(compressed, int64(count), 9)
	if err != E_SUCCESS {
		t.Fatalf("uncompactCellsSize failed: %v", err)
	}
	decompressed := make([]H3Index, countUncompact)
	err = uncompactCells(compressed, int64(count), decompressed, hexCount, 9)
	if err != E_SUCCESS {
		t.Fatalf("uncompactCells failed: %v", err)
	}

	count2 := 0
	for i := int64(0); i < hexCount; i++ {
		if decompressed[i] != 0 {
			count2++
		}
	}
	if count2 != int(hexCount) {
		t.Errorf("got uncompacted count %d, expected %d", count2, hexCount)
	}
}

func TestCompactCells_duplicate(t *testing.T) {
	t.Parallel()
	numHex := int64(10)
	someHexagons := make([]H3Index, 10)
	for i := 0; i < 10; i++ {
		setH3Index(&someHexagons[i], 5, 0, 2)
	}
	compressed := make([]H3Index, 10)

	err := compactCells(someHexagons, compressed, numHex)
	if err != E_DUPLICATE_INPUT {
		t.Errorf("compactCells should fail on duplicate input, got %v", err)
	}
}

func TestCompactCells_duplicateMinimum(t *testing.T) {
	t.Parallel()
	// Test that the minimum number of duplicate hexagons causes failure
	var h3 H3Index
	res := int32(10)
	// Arbitrary index
	setH3Index(&h3, res, 0, 2)

	arrSize, err := cellToChildrenSize(h3, res+1)
	if err != E_SUCCESS {
		t.Fatalf("cellToChildrenSize failed: %v", err)
	}
	arrSize++
	children := make([]H3Index, arrSize)

	err = cellToChildren(h3, res+1, children)
	if err != E_SUCCESS {
		t.Fatalf("cellToChildren failed: %v", err)
	}
	// duplicate one index
	children[arrSize-1] = children[0]

	output := make([]H3Index, arrSize)

	compactCellsResult := compactCells(children, output, arrSize)
	if compactCellsResult != E_DUPLICATE_INPUT {
		t.Errorf("compactCells should fail on duplicate input (single duplicate), got %v", compactCellsResult)
	}
}

func TestCompactCells_duplicatePentagonLimit(t *testing.T) {
	t.Parallel()
	// Test that the minimum number of duplicate hexagons causes failure
	var h3 H3Index
	res := int32(10)
	// Arbitrary pentagon parent cell
	setH3Index(&h3, res, 4, 0)

	arrSize, err := cellToChildrenSize(h3, res+1)
	if err != E_SUCCESS {
		t.Fatalf("cellToChildrenSize failed: %v", err)
	}
	arrSize++
	children := make([]H3Index, arrSize)

	err = cellToChildren(h3, res+1, children)
	if err != E_SUCCESS {
		t.Fatalf("cellToChildren failed: %v", err)
	}
	// duplicate one index
	centerChild, err := cellToCenterChild(h3, res+1)
	if err != E_SUCCESS {
		t.Fatalf("cellToCenterChild failed: %v", err)
	}
	children[arrSize-1] = centerChild

	output := make([]H3Index, arrSize)

	compactCellsResult := compactCells(children, output, arrSize)
	if compactCellsResult != E_DUPLICATE_INPUT {
		t.Errorf("compactCells should fail on duplicate input (pentagon parent), got %v", compactCellsResult)
	}
}

func TestCompactCells_duplicateIgnored(t *testing.T) {
	t.Parallel()
	// Test that duplicated cells are not rejected by compactCells.
	// This is not necessarily desired - just asserting the
	// existing behavior.
	var h3 H3Index
	res := int32(10)
	// Arbitrary index
	setH3Index(&h3, res, 0, 2)

	arrSize, err := cellToChildrenSize(h3, res+1)
	if err != E_SUCCESS {
		t.Fatalf("cellToChildrenSize failed: %v", err)
	}
	children := make([]H3Index, arrSize)

	err = cellToChildren(h3, res+1, children)
	if err != E_SUCCESS {
		t.Fatalf("cellToChildren failed: %v", err)
	}
	// duplicate one index
	children[arrSize-1] = children[0]

	output := make([]H3Index, arrSize)

	err = compactCells(children, output, arrSize)
	if err != E_SUCCESS {
		t.Fatalf("compactCells failed: %v", err)
	}
}

func TestCompactCells_empty(t *testing.T) {
	t.Parallel()
	err := compactCells(nil, nil, 0)
	if err != E_SUCCESS {
		t.Errorf("compactCells should succeed on empty input, got %v", err)
	}
}

func TestCompactCells_disparate(t *testing.T) {
	t.Parallel()
	// Exercises a case where compaction needs to be tested but none is
	// possible
	const numHex = 7
	disparate := make([]H3Index, numHex)
	for i := 0; i < numHex; i++ {
		setH3Index(&disparate[i], 1, int32(i), int32(CENTER_DIGIT))
	}
	output := make([]H3Index, numHex)

	err := compactCells(disparate, output, numHex)
	if err != E_SUCCESS {
		t.Errorf("compactCells should succeed on disparate input, got %v", err)
	}

	// Assumes that `output` is an exact copy of `disparate`, including
	// the ordering (which may not necessarily be the case)
	for i := 0; i < numHex; i++ {
		if disparate[i] != output[i] {
			t.Errorf("output[%d] = %x, expected %x", i, output[i], disparate[i])
		}
	}
}

func TestCompactCells_reservedBitsSet(t *testing.T) {
	t.Parallel()
	const numHex = 7
	bad := []H3Index{
		0x0010000000010000, 0x0180c6c6c6c61616, 0x1616ffffffffffff,
		0xffff8affffffffff, 0xffffffffffffc6c6, 0xffffffffffffffc6,
		0xc6c6c6c6c66fffe0,
	}
	output := make([]H3Index, numHex)

	err := compactCells(bad, output, numHex)
	if err != E_CELL_INVALID {
		t.Errorf("compactCells should return E_CELL_INVALID on bad input, got %v", err)
	}
}

func TestCompactCells_parentError(t *testing.T) {
	t.Parallel()
	const numHex = 3
	bad := make([]H3Index, numHex)
	output := make([]H3Index, numHex)
	bad[0] = setResolution(bad[0], 10)
	bad[1] = setResolution(bad[1], 5)

	err := compactCells(bad, output, numHex)
	if err != E_RES_MISMATCH {
		t.Errorf("compactCells should return E_RES_MISMATCH on bad input (parent error), got %v", err)
	}
}

func TestCompactCells_parentError2(t *testing.T) {
	t.Parallel()
	// This test primarily ensures memory is not leaked upon invalid input,
	// and ensures coverage of a very particular error branch in
	// compactCells. The particular error code is not important.
	const numHex = 43
	bad := []H3Index{0x2010202020202020,
		0x2100000000,
		0x7,
		0x400000000,
		0x20000000,
		0x5000000000,
		0x100321,
		0x2100000000,
		0x7,
		0x400000000,
		0x20000000,
		0x2100000000,
		0x7,
		0x400000000,
		0x20000000,
		0x5000000000,
		0x100321,
		0x20000000,
		0x5000000000,
		0x100321,
		0x2100000000,
		0x7,
		0x400000000,
		0x5000000000,
		0x100321,
		0x2100000000,
		0x7,
		0x400000000,
		0x20000000,
		0x5000000000,
		0x100321,
		0x2100000000,
		0x7,
		0x400000000,
		0x20000000,
		0x5000000000,
		0x100321,
		0x20000000,
		0x5000000000,
		0x100321,
		0x2100000000,
		0x7,
		0x400000000}
	output := make([]H3Index, 43)
	err := compactCells(bad, output, numHex)
	if err != E_SUCCESS {
		t.Fatalf("compactCells failed: %v", err)
	}
}

func TestUncompactCells_wrongRes(t *testing.T) {
	t.Parallel()
	numHex := int64(3)
	someHexagons := make([]H3Index, numHex)
	for i := int64(0); i < numHex; i++ {
		setH3Index(&someHexagons[i], 5, int32(i), 0)
	}

	_, err := uncompactCellsSize(someHexagons, numHex, 4)
	if err != E_RES_MISMATCH {
		t.Errorf("uncompactCellsSize should fail when given illogical resolutions, got %v", err)
	}
	_, err = uncompactCellsSize(someHexagons, numHex, -1)
	if err != E_RES_MISMATCH {
		t.Errorf("uncompactCellsSize should fail when given illegal resolutions, got %v", err)
	}
	_, err = uncompactCellsSize(someHexagons, numHex, MAX_H3_RES+1)
	if err != E_RES_MISMATCH {
		t.Errorf("uncompactCellsSize should fail when given resolutions beyond max, got %v", err)
	}

	uncompressed := make([]H3Index, numHex)
	uncompactCellsResult := uncompactCells(someHexagons, numHex, uncompressed, numHex, 0)
	if uncompactCellsResult != E_RES_MISMATCH {
		t.Errorf("uncompactCells should fail when given illogical resolutions, got %v", uncompactCellsResult)
	}
	uncompactCellsResult = uncompactCells(someHexagons, numHex, uncompressed, numHex, 6)
	if uncompactCellsResult != E_MEMORY_BOUNDS {
		t.Errorf("uncompactCells should fail when given too little buffer, got %v", uncompactCellsResult)
	}
	uncompactCellsResult = uncompactCells(someHexagons, numHex, uncompressed, numHex-1, 5)
	if uncompactCellsResult != E_MEMORY_BOUNDS {
		t.Errorf("uncompactCells should fail when given too little buffer (same resolution), got %v", uncompactCellsResult)
	}

	for i := int64(0); i < numHex; i++ {
		setH3Index(&someHexagons[i], MAX_H3_RES, int32(i), 0)
	}
	uncompactCellsResult = uncompactCells(someHexagons, numHex, uncompressed, numHex*7, MAX_H3_RES+1)
	if uncompactCellsResult != E_RES_MISMATCH {
		t.Errorf("uncompactCells should fail when given resolutions beyond max, got %v", uncompactCellsResult)
	}
}

func TestCompactCells_someHexagon(t *testing.T) {
	t.Parallel()
	var origin H3Index
	setH3Index(&origin, 1, 5, 0)

	childrenSz, err := uncompactCellsSize([]H3Index{origin}, 1, 2)
	if err != E_SUCCESS {
		t.Fatalf("uncompactCellsSize failed: %v", err)
	}
	children := make([]H3Index, childrenSz)
	err = uncompactCells([]H3Index{origin}, 1, children, childrenSz, 2)
	if err != E_SUCCESS {
		t.Fatalf("uncompactCells failed: %v", err)
	}

	result := make([]H3Index, childrenSz)
	err = compactCells(children, result, childrenSz)
	if err != E_SUCCESS {
		t.Fatalf("compactCells failed: %v", err)
	}

	found := 0
	for i := int64(0); i < childrenSz; i++ {
		if result[i] != 0 {
			found++
			if result[i] != origin {
				t.Errorf("compacted to wrong origin: got %x, expected %x", result[i], origin)
			}
		}
	}
	if found != 1 {
		t.Errorf("compacted to %d hexagons, expected 1", found)
	}
}

func TestUncompactCells_empty(t *testing.T) {
	t.Parallel()
	uncompactSz, err := uncompactCellsSize(nil, 0, 0)
	if err != E_SUCCESS {
		t.Fatalf("uncompactCellsSize should accept empty input, got %v", err)
	}
	if uncompactSz != 0 {
		t.Errorf("uncompactCellsSize should return 0 for empty input, got %d", uncompactSz)
	}
	err = uncompactCells(nil, 0, nil, 0, 0)
	if err != E_SUCCESS {
		t.Errorf("uncompactCells should accept empty input, got %v", err)
	}
}

func TestUncompactCells_onlyZero(t *testing.T) {
	t.Parallel()
	// uncompactCellsSize and uncompactCells both permit 0 indexes
	// in the input array, and skip them. When only a zero is
	// given, it's a no-op.

	origin := H3Index(0)

	childrenSz, err := uncompactCellsSize([]H3Index{origin}, 1, 2)
	if err != E_SUCCESS {
		t.Fatalf("uncompactCellsSize failed: %v", err)
	}
	children := make([]H3Index, childrenSz)
	err = uncompactCells([]H3Index{origin}, 1, children, childrenSz, 2)
	if err != E_SUCCESS {
		t.Fatalf("uncompactCells failed: %v", err)
	}
}

func TestUncompactCells_withZero(t *testing.T) {
	t.Parallel()
	// uncompactCellsSize and uncompactSize both permit 0 indexes
	// in the input array, and skip them.

	childrenSz, err := uncompactCellsSize(uncompactableWithZero, 4, 10)
	if err != E_SUCCESS {
		t.Fatalf("uncompactCellsSize failed: %v", err)
	}
	children := make([]H3Index, childrenSz)
	err = uncompactCells(uncompactableWithZero, 4, children, childrenSz, 10)
	if err != E_SUCCESS {
		t.Fatalf("uncompactCells failed: %v", err)
	}

	found := 0
	for i := int64(0); i < childrenSz; i++ {
		if children[i] != 0 {
			found++
		}
	}
	if found != int(childrenSz) {
		t.Errorf("uncompacted with zero to %d hexagons, expected %d", found, childrenSz)
	}
}

func TestCompactCells_pentagon(t *testing.T) {
	t.Parallel()
	var pentagon H3Index
	setH3Index(&pentagon, 1, 4, 0)

	childrenSz, err := uncompactCellsSize([]H3Index{pentagon}, 1, 2)
	if err != E_SUCCESS {
		t.Fatalf("uncompactCellsSize failed: %v", err)
	}
	children := make([]H3Index, childrenSz)
	err = uncompactCells([]H3Index{pentagon}, 1, children, childrenSz, 2)
	if err != E_SUCCESS {
		t.Fatalf("uncompactCells failed: %v", err)
	}

	result := make([]H3Index, childrenSz)
	err = compactCells(children, result, childrenSz)
	if err != E_SUCCESS {
		t.Fatalf("compactCells failed: %v", err)
	}

	found := 0
	for i := int64(0); i < childrenSz; i++ {
		if result[i] != 0 {
			found++
			if result[i] != pentagon {
				t.Errorf("compacted to wrong pentagon: got %x, expected %x", result[i], pentagon)
			}
		}
	}
	if found != 1 {
		t.Errorf("compacted to %d pentagons, expected 1", found)
	}
}

func TestUncompactCells_large_uncompact_size_hexagon(t *testing.T) {
	t.Parallel()
	cells := []H3Index{0x806dfffffffffff} // res 0 *hexagon*
	res := int32(15)

	expected := int64(4747561509943) // 7^15
	out, err := uncompactCellsSize(cells, 1, res)
	if err != E_SUCCESS {
		t.Fatalf("uncompactCellsSize failed: %v", err)
	}

	if out != expected {
		t.Errorf("uncompactCells size needs 64 bit int: got %d, expected %d", out, expected)
	}
}

func TestUncompactCells_large_uncompact_size_pentagon(t *testing.T) {
	t.Parallel()
	cells := []H3Index{0x8009fffffffffff} // res 0 *pentagon*
	res := int32(15)

	expected := int64(3956301258286) // 1 + 5*(7^15 - 1)/6
	out, err := uncompactCellsSize(cells, 1, res)
	if err != E_SUCCESS {
		t.Fatalf("uncompactCellsSize failed: %v", err)
	}

	if out != expected {
		t.Errorf("uncompactCells size needs 64 bit int: got %d, expected %d", out, expected)
	}
}
