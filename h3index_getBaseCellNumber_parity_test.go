//go:build cgo && c2go

package h3

import (
	"math/rand"
	"testing"
)

func TestGetBaseCellNumber_parity(t *testing.T) {
	// Test cases with H3 indexes constructed for different base cells
	testBaseCells := []int32{1, 4, 7, 14, 24, 38, 49, 58, 63, 72, 83, 97, 107, 117}
	testCases := make([]h3Index, len(testBaseCells))
	for i, baseCell := range testBaseCells {
		testCases[i] = setH3IndexC(0, baseCell, 0)
	}

	for _, h3Index := range testCases {
		// Test the Go implementation
		goResult := getBaseCellNumber(h3Index)
		// Test the C implementation
		cResult := getBaseCellNumberC(h3Index)

		if goResult != cResult {
			t.Errorf("getBaseCellNumber parity failure for h3Index %016x: Go=%d, C=%d",
				h3Index, goResult, cResult)
		}
	}

	// Test with random H3 indexes at different resolutions
	r := rand.New(rand.NewSource(12345))
	for res := int32(0); res <= 15; res++ {
		for baseCell := int32(0); baseCell < numBaseCells; baseCell++ {
			// Create a valid H3 index
			h := setH3IndexC(res, baseCell, 0)

			// Add some random digits for higher resolutions
			if res > 0 {
				for digit := int32(1); digit <= res; digit++ {
					randomDigit := int32(r.Intn(7))
					// Avoid invalid digit 7 and handle pentagon special case
					if randomDigit == 7 || (_isBaseCellPentagon(baseCell) && randomDigit == 1 && digit == 1) {
						randomDigit = 0
					}
					h = h3SetIndexDigitC(h, digit, randomDigit)
				}
			}

			// Test the parity
			goResult := getBaseCellNumber(h)
			cResult := getBaseCellNumberC(h)

			if goResult != cResult {
				t.Errorf("getBaseCellNumber parity failure for h3Index %016x (res=%d, baseCell=%d): Go=%d, C=%d",
					h, res, baseCell, goResult, cResult)
			}
		}
	}

	// Test edge cases
	edgeCases := []h3Index{
		0x0,                    // h3Null
		setH3IndexC(0, 0, 0),   // Base cell 0
		setH3IndexC(0, 121, 0), // Max base cell (base cell 121)
		setH3IndexC(15, 0, 0),  // Max resolution (res 15)
	}

	for _, h3Index := range edgeCases {
		goResult := getBaseCellNumber(h3Index)
		cResult := getBaseCellNumberC(h3Index)

		if goResult != cResult {
			t.Errorf("getBaseCellNumber parity failure for edge case h3Index %016x: Go=%d, C=%d",
				h3Index, goResult, cResult)
		}
	}
}

// Test that getBaseCellNumber returns expected values for properly constructed indexes
func TestGetBaseCellNumber_knownValues(t *testing.T) {
	// Test each base cell by constructing a proper H3 index for it
	testBaseCells := []int32{1, 4, 7, 14, 24, 38, 49, 58, 63, 72, 83, 97, 107, 117}

	for _, expectedBaseCell := range testBaseCells {
		// Create a valid H3 index for this base cell at resolution 0
		h3Index := setH3IndexC(0, expectedBaseCell, 0)

		goResult := getBaseCellNumber(h3Index)
		cResult := getBaseCellNumberC(h3Index)

		if goResult != expectedBaseCell {
			t.Errorf("getBaseCellNumber returned %d for h3Index %016x (base cell %d), expected %d",
				goResult, h3Index, expectedBaseCell, expectedBaseCell)
		}

		if cResult != expectedBaseCell {
			t.Errorf("getBaseCellNumberC returned %d for h3Index %016x (base cell %d), expected %d",
				cResult, h3Index, expectedBaseCell, expectedBaseCell)
		}

		// Also test at higher resolution
		h3IndexRes5 := setH3IndexC(5, expectedBaseCell, 0)
		goResult5 := getBaseCellNumber(h3IndexRes5)
		cResult5 := getBaseCellNumberC(h3IndexRes5)

		if goResult5 != expectedBaseCell {
			t.Errorf("getBaseCellNumber returned %d for h3Index %016x (base cell %d, res 5), expected %d",
				goResult5, h3IndexRes5, expectedBaseCell, expectedBaseCell)
		}

		if cResult5 != expectedBaseCell {
			t.Errorf("getBaseCellNumberC returned %d for h3Index %016x (base cell %d, res 5), expected %d",
				cResult5, h3IndexRes5, expectedBaseCell, expectedBaseCell)
		}
	}
}
