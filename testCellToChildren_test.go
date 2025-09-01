// Tests ported from testCellToChildren.c
package h3

import (
	"testing"
)

// assertNoDuplicates checks that there are no duplicate cells in the slice
func assertNoDuplicates(t *testing.T, cells []H3Index, testName string) {
	t.Helper()
	for i := 0; i < len(cells); i++ {
		if cells[i] == H3_NULL {
			continue
		}
		if !isValidCell(cells[i]) {
			t.Errorf("%s: cell at index %d (0x%x) must be valid H3 cell", testName, i, cells[i])
		}
		for j := i + 1; j < len(cells); j++ {
			if cells[i] == cells[j] {
				t.Errorf("%s: can't have duplicate cells in set: cells[%d] == cells[%d] == 0x%x", testName, i, j, cells[i])
			}
		}
	}
}

// assertSubset checks that set1 is a subset of set2
func assertSubset(t *testing.T, set1, set2 []H3Index, testName string) {
	t.Helper()
	assertNoDuplicates(t, set1, testName+" (set1)")

	for i := 0; i < len(set1); i++ {
		if set1[i] == H3_NULL {
			continue
		}

		present := false
		for j := 0; j < len(set2); j++ {
			if set1[i] == set2[j] {
				present = true
				break
			}
		}
		if !present {
			t.Errorf("%s: children must match - 0x%x not found in set2", testName, set1[i])
		}
	}
}

// assertSetsEqual asserts that two arrays of h3 cells are equal sets:
//   - No duplicate cells allowed.
//   - Ignore zero elements (so array sizes may be different).
//   - Ignore array order.
func assertSetsEqual(t *testing.T, set1, set2 []H3Index, testName string) {
	t.Helper()
	assertSubset(t, set1, set2, testName)
	assertSubset(t, set2, set1, testName)
}

// checkChildren helper function that mirrors the C version
func checkChildren(t *testing.T, h H3Index, res int32, expectedError H3Error, expected []H3Index, testName string) {
	t.Helper()

	numChildren, numChildrenError := cellToChildrenSize(h, res)
	if numChildrenError != expectedError {
		t.Errorf("%s: Expected error code %v, got %v", testName, expectedError, numChildrenError)
		return
	}
	if expectedError != E_SUCCESS {
		return
	}

	children := make([]H3Index, numChildren)
	err := cellToChildren(h, res, children)
	if err != E_SUCCESS {
		t.Fatalf("%s: cellToChildren failed: %v", testName, err)
	}

	assertSetsEqual(t, children, expected, testName)
}

func Test_oneResStep(t *testing.T) {
	t.Parallel()

	h := H3Index(0x88283080ddfffff)
	res := int32(9)

	expected := []H3Index{
		0x89283080dc3ffff, 0x89283080dc7ffff,
		0x89283080dcbffff, 0x89283080dcfffff,
		0x89283080dd3ffff, 0x89283080dd7ffff,
		0x89283080ddbffff,
	}

	checkChildren(t, h, res, E_SUCCESS, expected, "oneResStep")
}

func Test_multipleResSteps(t *testing.T) {
	t.Parallel()

	h := H3Index(0x88283080ddfffff)
	res := int32(10)

	expected := []H3Index{
		0x8a283080dd27fff, 0x8a283080dd37fff, 0x8a283080dc47fff,
		0x8a283080dcdffff, 0x8a283080dc5ffff, 0x8a283080dc27fff,
		0x8a283080ddb7fff, 0x8a283080dc07fff, 0x8a283080dd8ffff,
		0x8a283080dd5ffff, 0x8a283080dc4ffff, 0x8a283080dd47fff,
		0x8a283080dce7fff, 0x8a283080dd1ffff, 0x8a283080dceffff,
		0x8a283080dc6ffff, 0x8a283080dc87fff, 0x8a283080dcaffff,
		0x8a283080dd2ffff, 0x8a283080dcd7fff, 0x8a283080dd9ffff,
		0x8a283080dd6ffff, 0x8a283080dcc7fff, 0x8a283080dca7fff,
		0x8a283080dccffff, 0x8a283080dd77fff, 0x8a283080dc97fff,
		0x8a283080dd4ffff, 0x8a283080dd97fff, 0x8a283080dc37fff,
		0x8a283080dc8ffff, 0x8a283080dcb7fff, 0x8a283080dcf7fff,
		0x8a283080dd87fff, 0x8a283080dda7fff, 0x8a283080dc9ffff,
		0x8a283080dc77fff, 0x8a283080dc67fff, 0x8a283080dc57fff,
		0x8a283080ddaffff, 0x8a283080dd17fff, 0x8a283080dc17fff,
		0x8a283080dd57fff, 0x8a283080dc0ffff, 0x8a283080dd07fff,
		0x8a283080dc1ffff, 0x8a283080dd0ffff, 0x8a283080dc2ffff,
		0x8a283080dd67fff,
	}

	checkChildren(t, h, res, E_SUCCESS, expected, "multipleResSteps")
}

func Test_sameRes(t *testing.T) {
	t.Parallel()

	h := H3Index(0x88283080ddfffff)
	res := int32(8)

	expected := []H3Index{h}

	checkChildren(t, h, res, E_SUCCESS, expected, "sameRes")
}

func Test_childResTooCoarse(t *testing.T) {
	t.Parallel()

	h := H3Index(0x88283080ddfffff)
	res := int32(7)

	expected := []H3Index{0} // empty set; zeros are ignored

	checkChildren(t, h, res, E_RES_DOMAIN, expected, "childResTooCoarse")
}

func Test_childResTooFine(t *testing.T) {
	t.Parallel()

	h := H3Index(0x8f283080dcb0ae2) // res 15 cell
	res := int32(MAX_H3_RES + 1)

	expected := []H3Index{0} // empty set; zeros are ignored

	checkChildren(t, h, res, E_RES_DOMAIN, expected, "childResTooFine")
}

func Test_pentagonChildren(t *testing.T) {
	t.Parallel()

	h := H3Index(0x81083ffffffffff) // res 1 pentagon
	res := int32(3)

	expected := []H3Index{
		0x830800fffffffff, 0x830802fffffffff, 0x830803fffffffff,
		0x830804fffffffff, 0x830805fffffffff, 0x830806fffffffff,
		0x830810fffffffff, 0x830811fffffffff, 0x830812fffffffff,
		0x830813fffffffff, 0x830814fffffffff, 0x830815fffffffff,
		0x830816fffffffff, 0x830818fffffffff, 0x830819fffffffff,
		0x83081afffffffff, 0x83081bfffffffff, 0x83081cfffffffff,
		0x83081dfffffffff, 0x83081efffffffff, 0x830820fffffffff,
		0x830821fffffffff, 0x830822fffffffff, 0x830823fffffffff,
		0x830824fffffffff, 0x830825fffffffff, 0x830826fffffffff,
		0x830828fffffffff, 0x830829fffffffff, 0x83082afffffffff,
		0x83082bfffffffff, 0x83082cfffffffff, 0x83082dfffffffff,
		0x83082efffffffff, 0x830830fffffffff, 0x830831fffffffff,
		0x830832fffffffff, 0x830833fffffffff, 0x830834fffffffff,
		0x830835fffffffff, 0x830836fffffffff,
	}

	checkChildren(t, h, res, E_SUCCESS, expected, "pentagonChildren")
}
