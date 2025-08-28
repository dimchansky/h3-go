//go:build cgo

package h3

import "testing"

func Test_isBaseCellPolarPentagon_ParityWithC(t *testing.T) {
	testCases := []int32{
		// All base cells 0-121
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
		10, 11, 12, 13, 14, 15, 16, 17, 18, 19,
		20, 21, 22, 23, 24, 25, 26, 27, 28, 29,
		30, 31, 32, 33, 34, 35, 36, 37, 38, 39,
		40, 41, 42, 43, 44, 45, 46, 47, 48, 49,
		50, 51, 52, 53, 54, 55, 56, 57, 58, 59,
		60, 61, 62, 63, 64, 65, 66, 67, 68, 69,
		70, 71, 72, 73, 74, 75, 76, 77, 78, 79,
		80, 81, 82, 83, 84, 85, 86, 87, 88, 89,
		90, 91, 92, 93, 94, 95, 96, 97, 98, 99,
		100, 101, 102, 103, 104, 105, 106, 107, 108, 109,
		110, 111, 112, 113, 114, 115, 116, 117, 118, 119,
		120, 121,

		// Additional edge cases
		-1, -10, 122, 200, 1000,
	}

	for _, baseCell := range testCases {
		goResult := _isBaseCellPolarPentagon(baseCell)
		cResult := _isBaseCellPolarPentagonC(baseCell)

		if goResult != cResult {
			t.Fatalf("_isBaseCellPolarPentagon mismatch: baseCell=%d go=%t c=%t",
				baseCell, goResult, cResult)
		}

		// Verify specific expected results for polar pentagons
		if baseCell == 4 || baseCell == 117 {
			if !goResult {
				t.Errorf("Expected baseCell %d to be polar pentagon but got false", baseCell)
			}
		} else if baseCell >= 0 && baseCell <= 121 {
			// Valid base cell that's not polar pentagon
			if goResult {
				t.Errorf("Expected baseCell %d to NOT be polar pentagon but got true", baseCell)
			}
		}
	}
}
