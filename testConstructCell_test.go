package h3

import "testing"

// Ported from H3 C: testConstructCell.c (added in 4.4.0).

// constructCellTestCase mirrors the upstream TestCase struct: x is either a
// valid expected H3 cell or an expected error code.
type constructCellTestCase struct {
	x      uint64
	res    int32
	bc     int32
	digits [15]int32
}

func runConstructCellCase(t *testing.T, tc constructCellTestCase) {
	t.Helper()

	var h h3Index
	validTcx := isValidCell(h3Index(tc.x))
	err := constructCell(tc.res, tc.bc, tc.digits[:], &h)

	gotExpectedValidCell := validTcx && err == eSuccess && h3Index(tc.x) == h
	gotExpectedError := !validTcx && tc.x == uint64(err)

	if !gotExpectedValidCell && !gotExpectedError {
		t.Errorf("constructCell(res=%d, bc=%d, digits=%v): got (%#x, %v), want %#x",
			tc.res, tc.bc, tc.digits, uint64(h), err, tc.x)
	}
}

func TestConstructCell_tableOfTests(t *testing.T) {
	t.Parallel()

	tests := []constructCellTestCase{
		// a few valid cell constructions
		{x: 0x8001fffffffffff, res: 0, bc: 0},
		{x: 0x8003fffffffffff, res: 0, bc: 1},
		{x: 0x80f3fffffffffff, res: 0, bc: 121},
		{x: 0x839253fffffffff, res: 3, bc: 73, digits: [15]int32{1, 2, 3}},
		{x: 0x821f67fffffffff, res: 2, bc: 15, digits: [15]int32{5, 4}},
		{x: 0x8155bffffffffff, res: 1, bc: 42, digits: [15]int32{6}},
		{x: 0x8f754e64992d6d8, res: 15, bc: 58,
			digits: [15]int32{5, 1, 6, 3, 1, 1, 1, 4, 4, 5, 5, 3, 3, 3, 0}},

		// tests around resolution
		{res: 16, bc: 0, x: uint64(eResDomain)},
		{res: 18, bc: 0, x: uint64(eResDomain)},
		{res: -1, bc: 0, x: uint64(eResDomain)},
		{res: 0, bc: 0, x: 0x8001fffffffffff},

		// tests around base cell
		{res: 0, bc: 122, x: uint64(eBaseCellDomain)},
		{res: 0, bc: -1, x: uint64(eBaseCellDomain)},
		{res: 0, bc: 259, x: uint64(eBaseCellDomain)},
		{res: 2, bc: 122, digits: [15]int32{1, 0}, x: uint64(eBaseCellDomain)},

		// tests around digits
		{res: 1, bc: 40, digits: [15]int32{-1}, x: uint64(eDigitDomain)},
		{res: 1, bc: 40, digits: [15]int32{7}, x: uint64(eDigitDomain)},
		{res: 1, bc: 40, digits: [15]int32{8}, x: uint64(eDigitDomain)},
		{res: 1, bc: 40, digits: [15]int32{17}, x: uint64(eDigitDomain)},

		// deleted subsequence tests: bc = 4 is a pentagon base cell
		{bc: 4, digits: [15]int32{0, 0, 0}, res: 3, x: 0x830800fffffffff},
		{bc: 4, digits: [15]int32{0, 0, 1}, res: 3, x: uint64(eDeletedDigit)},
		{bc: 4, digits: [15]int32{0, 0, 2}, res: 3, x: 0x830802fffffffff},

		// more deleted subsequence tests: bc = 5 is *not* a pentagon base cell
		{bc: 5, digits: [15]int32{0, 0, 0}, res: 3, x: 0x830a00fffffffff},
		{bc: 5, digits: [15]int32{0, 0, 1}, res: 3, x: 0x830a01fffffffff},
		{bc: 5, digits: [15]int32{0, 0, 2}, res: 3, x: 0x830a02fffffffff},
	}

	for _, tc := range tests {
		runConstructCellCase(t, tc)
	}
}

// passesConstructRoundtrip mirrors upstream passesRoundtrip: index ->
// components -> index.
func passesConstructRoundtrip(h h3Index) bool {
	res := getResolution(h)
	bc := getBaseCellNumber(h)
	var digits [15]int32

	for r := int32(1); r <= res; r++ {
		if getIndexDigit(h, r, &digits[r-1]) != eSuccess {
			return false
		}
	}

	var out h3Index
	if constructCell(res, bc, digits[:], &out) != eSuccess {
		return false
	}
	return out == h
}

func TestConstructCell_roundtrip(t *testing.T) {
	t.Parallel()

	// Test roundtrip for all cells at a few resolutions.
	for res := int32(0); res <= 4; res++ {
		allPassed := true
		iter := iterInitRes(res)
		for iter.H != h3Null {
			if !passesConstructRoundtrip(iter.H) {
				allPassed = false
			}
			iterStepRes(&iter)
		}
		if !allPassed {
			t.Errorf("not all cells at res %d passed the roundtrip", res)
		}
	}
}
