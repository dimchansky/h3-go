//go:build cgo && c2go

package h3

import (
	"testing"
)

func Test_vertexNumForDirection_parity(t *testing.T) {
	// Test with various cells including hexagons and pentagons
	testCells := []h3Index{
		0x85283473fffffff, // res 5 hexagon
		0x8928308280fffff, // res 9 hexagon
		0x8a2a1072b587fff, // res 10 hexagon
		0x8b2a1072b59ffff, // res 11 hexagon
		0x821c07fffffffff, // res 1 pentagon
		0x841c001ffffffff, // res 4 pentagon
		0x851c003ffffffff, // res 5 pentagon
		0x891c000003fffff, // res 9 pentagon

		// Additional test cases at various resolutions
		0x8001fffffffffff, // res 0 hexagon
		0x8101fffffffffff, // res 1 hexagon
		0x8201fffffffffff, // res 2 hexagon
		0x8301fffffffffff, // res 3 hexagon
		0x8401fffffffffff, // res 4 hexagon
		0x8501fffffffffff, // res 5 hexagon
		0x8601fffffffffff, // res 6 hexagon
		0x8701fffffffffff, // res 7 hexagon
		0x8801fffffffffff, // res 8 hexagon
		0x8901fffffffffff, // res 9 hexagon
		0x8a01fffffffffff, // res 10 hexagon
		0x8b01fffffffffff, // res 11 hexagon
		0x8c01fffffffffff, // res 12 hexagon
		0x8d01fffffffffff, // res 13 hexagon
		0x8e01fffffffffff, // res 14 hexagon
		0x8f01fffffffffff, // res 15 hexagon
	}

	// Test all valid directions
	testDirections := []direction{
		centerDigit, kAxesDigit, jAxesDigit, jkAxesDigit,
		iAxesDigit, ikAxesDigit, ijAxesDigit, invalidDigit,
	}

	for _, cell := range testCells {
		if !isValidCell(cell) {
			continue
		}

		for _, direction := range testDirections {
			goOut := vertexNumForDirection(cell, direction)
			cOut := vertexNumForDirectionC(cell, direction)

			if goOut != cOut {
				t.Errorf("vertexNumForDirection(0x%x, %d) output mismatch: Go=%d, C=%d", cell, direction, goOut, cOut)
			}
		}
	}
}

func Test_vertexNumForDirection_parity_pentagons(t *testing.T) {
	// Test all pentagons at different resolutions
	for res := int32(0); res <= 15; res++ {
		pentagons := make([]h3Index, numPentagons)
		err := getPentagons(res, pentagons)
		if err != eSuccess {
			t.Fatalf("getPentagons failed for res=%d: %v", res, err)
		}

		for _, pentagon := range pentagons {
			// Test all valid directions
			testDirections := []direction{
				centerDigit, kAxesDigit, jAxesDigit, jkAxesDigit,
				iAxesDigit, ikAxesDigit, ijAxesDigit, invalidDigit,
			}

			for _, direction := range testDirections {
				goOut := vertexNumForDirection(pentagon, direction)
				cOut := vertexNumForDirectionC(pentagon, direction)

				if goOut != cOut {
					t.Errorf("vertexNumForDirection(pentagon 0x%x res=%d, dir=%d) output mismatch: Go=%d, C=%d",
						pentagon, res, direction, goOut, cOut)
				}
			}
		}
	}
}

func Test_vertexNumForDirection_parity_invalid_cases(t *testing.T) {
	// Test invalid cases that should return invalidVertexNum
	testCases := []struct {
		cell      h3Index
		direction direction
		desc      string
	}{
		{0x85283473fffffff, centerDigit, "centerDigit on hexagon"},
		{0x85283473fffffff, invalidDigit, "invalidDigit on hexagon"},
		{0x821c07fffffffff, centerDigit, "centerDigit on pentagon"},
		{0x821c07fffffffff, kAxesDigit, "kAxesDigit on pentagon"},
		{0x821c07fffffffff, invalidDigit, "invalidDigit on pentagon"},
		{h3Null, jAxesDigit, "NULL cell"},
		{0xffffffffffffffff, jAxesDigit, "invalid cell"},
	}

	for _, tc := range testCases {
		goOut := vertexNumForDirection(tc.cell, tc.direction)
		cOut := vertexNumForDirectionC(tc.cell, tc.direction)

		if goOut != cOut {
			t.Errorf("vertexNumForDirection(%s) output mismatch: Go=%d, C=%d", tc.desc, goOut, cOut)
		}

		// Note: The C implementation doesn't necessarily return invalidVertexNum for all invalid cells
		// as it may process them through the calculation logic. We only test for parity between Go and C.
	}
}
