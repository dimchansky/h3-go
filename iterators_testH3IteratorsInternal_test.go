// Tests ported from H3 v4.4.0: src/apps/testapps/testH3IteratorsInternal.c.
package h3

import (
	"testing"
)

// Helper function to assert iterator is null (ported from assert_is_null_iterator).
func assertIsNullIterator(t *testing.T, iter iterCellsChildren, msg string) {
	t.Helper()
	if iter.H != h3Null {
		t.Errorf("%s: expected h3Null, got %x", msg, uint64(iter.H))
	}
	if iter.ParentRes != -1 {
		t.Errorf("%s: expected ParentRes=-1, got %d", msg, iter.ParentRes)
	}
	if iter.SkipDigit != -1 {
		t.Errorf("%s: expected SkipDigit=-1, got %d", msg, iter.SkipDigit)
	}
}

// Test number of cells at each resolution (ported from test_number).
func testNumber(t *testing.T, res int32) {
	t.Helper()

	count := int64(0)
	iter := iterInitRes(res)
	for iter.H != h3Null {
		count++
		iterStepRes(&iter)
	}

	expected, err := getNumCells(res)
	if err != eSuccess {
		t.Fatalf("getNumCells failed for res %d: %v", res, err)
	}

	if count != expected {
		t.Errorf("res %d: expected %d cells from iterator, got %d", res, expected, count)
	}
}

// Test that all iterated cells are valid (ported from test_valid).
func testValid(t *testing.T, res int32) {
	t.Helper()

	iter := iterInitRes(res)
	count := 0
	for iter.H != h3Null {
		if !isValidCell(iter.H) {
			t.Errorf("res %d, cell %d: iterator cell %x is not valid", res, count, uint64(iter.H))
		}
		iterStepRes(&iter)
		count++

		// Safety limit to prevent infinite loops (res 3 has ~41k cells)
		if count > 100000 {
			t.Fatalf("res %d: iterator seems stuck, stopped after %d iterations", res, count)
			break
		}
	}
}

// Test that all iterated cells have correct resolution (ported from test_resolution).
func testResolution(t *testing.T, res int32) {
	t.Helper()

	iter := iterInitRes(res)
	count := 0
	for iter.H != h3Null {
		cellRes := getResolution(iter.H)
		if cellRes != res {
			t.Errorf("res %d, cell %d: expected resolution %d, got %d for cell %x",
				res, count, res, cellRes, uint64(iter.H))
		}
		iterStepRes(&iter)
		count++

		// Safety limit to prevent infinite loops (res 3 has ~41k cells)
		if count > 100000 {
			t.Fatalf("res %d: iterator seems stuck, stopped after %d iterations", res, count)
			break
		}
	}
}

// Test that cells are iterated in order without duplicates (ported from test_ordered).
func testOrdered(t *testing.T, res int32) {
	t.Helper()

	iter := iterInitRes(res)
	if iter.H == h3Null {
		t.Errorf("res %d: iterator is null at start", res)
		return
	}

	prev := iter.H
	iterStepRes(&iter)
	count := 1

	for iter.H != h3Null {
		if prev >= iter.H {
			t.Errorf("res %d, cell %d: cells not in order: prev=%x curr=%x",
				res, count, uint64(prev), uint64(iter.H))
		}
		prev = iter.H
		iterStepRes(&iter)
		count++

		// Safety limit to prevent infinite loops (res 3 has ~41k cells)
		if count > 100000 {
			t.Fatalf("res %d: iterator seems stuck, stopped after %d iterations", res, count)
			break
		}
	}
}

func TestIteratorSetupInvalid(t *testing.T) {
	t.Parallel()

	// Test iterInitBaseCellNum with invalid inputs
	assertIsNullIterator(t, iterInitBaseCellNum(-1, 0), "iterInitBaseCellNum(-1, 0)")
	assertIsNullIterator(t, iterInitBaseCellNum(1000, 0), "iterInitBaseCellNum(1000, 0)")
	assertIsNullIterator(t, iterInitBaseCellNum(0, -1), "iterInitBaseCellNum(0, -1)")
	assertIsNullIterator(t, iterInitBaseCellNum(0, 100), "iterInitBaseCellNum(0, 100)")

	// Test iterInitParent with invalid inputs
	var iter iterCellsChildren
	iterInitParent(0, 0, &iter)
	assertIsNullIterator(t, iter, "iterInitParent(0, 0)")

	testCell := h3Index(0x85283473fffffff)
	iterInitParent(testCell, 0, &iter)
	assertIsNullIterator(t, iter, "iterInitParent(testCell, 0)")

	iterInitParent(testCell, 100, &iter)
	assertIsNullIterator(t, iter, "iterInitParent(testCell, 100)")
}

func TestNullIteratorBaseCell(t *testing.T) {
	t.Parallel()

	iter := iterInitBaseCellNum(-1, 0)
	assertIsNullIterator(t, iter, "initial state")
	iterStepChild(&iter)
	if iter.H != h3Null {
		t.Errorf("expected h3Null after step, got %x", uint64(iter.H))
	}
}

func TestNullIteratorRes(t *testing.T) {
	t.Parallel()

	iter := iterInitRes(-1)
	assertIsNullIterator(t, iter.itC, "initial state")
	iterStepRes(&iter)
	if iter.H != h3Null {
		t.Errorf("expected h3Null after step, got %x", uint64(iter.H))
	}
}

func TestIteratorCellCount(t *testing.T) {
	t.Parallel()

	testNumber(t, 0)
	testNumber(t, 1)
	testNumber(t, 2)
	testNumber(t, 3)
}

func TestIteratorCellValid(t *testing.T) {
	t.Parallel()

	testValid(t, 0)
	testValid(t, 1)
	testValid(t, 2)
	testValid(t, 3)
}

func TestIteratorCellResolution(t *testing.T) {
	t.Parallel()

	testResolution(t, 0)
	testResolution(t, 1)
	testResolution(t, 2)
	testResolution(t, 3)
}

func TestIteratorCellOrdered(t *testing.T) {
	t.Parallel()

	testOrdered(t, 0)
	testOrdered(t, 1)
	testOrdered(t, 2)
	testOrdered(t, 3)
}
