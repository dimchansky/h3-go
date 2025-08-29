// Tests ported from testH3NeighborRotations.c
package h3

import (
	"fmt"
	"testing"
)

// TestOutput represents validation results from the h3NeighborRotations tests
// This corresponds to the TestOutput struct in the C code
type TestOutput struct {
	ret0                   int
	ret0ValidationFailures int
	ret1                   int
	ret1ValidationFailures int
	ret2                   int
}

// doCell tests gridDiskUnsafe vs gridDiskDistancesSafe for a specific cell at various k values
func doCell(t *testing.T, h H3Index, maxK int32, testOutput *TestOutput) {
	for k := int32(0); k < maxK; k++ {
		var maxSz int64
		err := maxGridDiskSize(k, &maxSz)
		if err != E_SUCCESS {
			t.Fatalf("maxGridDiskSize failed: %v", err)
		}

		gridDiskInternalOutput := make([]H3Index, maxSz)
		gridDiskUnsafeOutput := make([]H3Index, maxSz)
		gridDiskInternalDistances := make([]int32, maxSz)

		// Call gridDiskDistancesSafe (equivalent to gridDiskDistancesInternal in C)
		err = gridDiskDistancesSafe(h, k, gridDiskInternalOutput, gridDiskInternalDistances)
		if err != E_SUCCESS {
			t.Fatalf("gridDiskDistancesSafe failed: %v", err)
		}

		// Call gridDiskUnsafe
		gridDiskUnsafeFailed := gridDiskUnsafe(h, k, gridDiskUnsafeOutput)

		if gridDiskUnsafeFailed == E_FAILED {
			// TODO: Unreachable
			testOutput.ret2++
			continue
		} else if gridDiskUnsafeFailed == E_SUCCESS {
			testOutput.ret0++
			startIdx := 0
			// i is the current ring number
			for i := int32(0); i <= k; i++ {
				// Number of hexagons on this ring
				n := i * 6
				if i == 0 {
					n = 1
				}

				for ii := int32(0); ii < n; ii++ {
					h2 := gridDiskUnsafeOutput[ii+int32(startIdx)]
					found := false

					for iii := int64(0); iii < maxSz; iii++ {
						if gridDiskInternalOutput[iii] == h2 &&
							gridDiskInternalDistances[iii] == i {
							found = true
							break
						}
					}

					if !found {
						// Failed to find a hexagon in both outputs, or it had
						// different distances.
						testOutput.ret0ValidationFailures++
						t.Logf("Failed validation for cell %016x", h)
						return
					}
				}

				startIdx += int(n)
			}
		} else if gridDiskUnsafeFailed == E_PENTAGON {
			testOutput.ret1++
			foundPent := false
			for i := int64(0); i < maxSz; i++ {
				if isPentagon(gridDiskInternalOutput[i]) {
					foundPent = true
					break
				}
			}

			if !foundPent {
				// Failed to find the pentagon that caused gridDiskUnsafe
				// to fail.
				t.Logf("NO C k=%d h=%016x", k, h)
				testOutput.ret1ValidationFailures++
				return
			}
		}
	}
}

// recursiveH3IndexToGeo recursively generates all valid H3 indexes at a given resolution
func recursiveH3IndexToGeo(t *testing.T, h H3Index, res int32, maxK int32, testOutput *TestOutput) {
	for d := int32(0); d < 7; d++ {
		current := setIndexDigit(h, res, d)

		// skip the pentagonal deleted subsequence
		if _isBaseCellPentagon(int32(getBaseCell(current))) &&
			Direction(_h3LeadingNonZeroDigit(current)) == K_AXES_DIGIT {
			continue
		}

		if res == getResolution(current) {
			doCell(t, current, maxK, testOutput)
		} else {
			recursiveH3IndexToGeo(t, current, res+1, maxK, testOutput)
		}
	}
}

// TestH3NeighborRotations tests gridDiskUnsafe vs. gridDiskDistancesSafe
// This is a comprehensive test that validates the consistency between
// the two grid disk algorithms across all base cells and resolutions
func TestH3NeighborRotations(t *testing.T) {
	t.Parallel()

	// Test parameters - using smaller values for reasonable test time
	const resolution = int32(1)
	const maxK = int32(3)

	testOutput := &TestOutput{0, 0, 0, 0, 0}

	// Generate test cases for all base cells
	for bc := int32(0); bc < NUM_BASE_CELLS; bc++ {
		rootCell := H3Index(H3_CELL_MODE << H3_MODE_OFFSET)
		rootCell = setBaseCell(rootCell, bc)

		if resolution == 0 {
			doCell(t, rootCell, maxK, testOutput)
		} else {
			rootRes := getResolution(rootCell)
			rootCell = setResolution(rootCell, resolution)
			recursiveH3IndexToGeo(t, rootCell, rootRes+1, maxK, testOutput)
		}
	}

	t.Logf("ret0: %d", testOutput.ret0)
	t.Logf("ret1: %d", testOutput.ret1)
	t.Logf("ret2: %d", testOutput.ret2)

	// ret2 should never occur, as it can only happen if we run over a pentagon
	if testOutput.ret2 > 0 || testOutput.ret0ValidationFailures > 0 ||
		testOutput.ret1ValidationFailures > 0 {
		t.Errorf("FAILED\nfailed0: %d\nfailed1: %d",
			testOutput.ret0ValidationFailures,
			testOutput.ret1ValidationFailures)
	}
}

// TestH3NeighborRotationsMultipleResolutions tests multiple resolutions
func TestH3NeighborRotationsMultipleResolutions(t *testing.T) {
	t.Parallel()

	// Test multiple resolutions with smaller maxK for performance
	resolutions := []int32{0, 1, 2}
	maxK := int32(2)

	for _, resolution := range resolutions {
		t.Run(fmt.Sprintf("resolution_%d", resolution), func(t *testing.T) {
			t.Parallel()

			testOutput := &TestOutput{0, 0, 0, 0, 0}

			// Test a subset of base cells for higher resolutions
			maxBaseCells := NUM_BASE_CELLS
			if resolution > 1 {
				maxBaseCells = 10 // Limit to first 10 base cells for performance
			}

			for bc := int32(0); bc < int32(maxBaseCells); bc++ {
				rootCell := H3Index(H3_CELL_MODE << H3_MODE_OFFSET)
				rootCell = setBaseCell(rootCell, bc)

				if resolution == 0 {
					doCell(t, rootCell, maxK, testOutput)
				} else {
					rootRes := getResolution(rootCell)
					rootCell = setResolution(rootCell, resolution)
					recursiveH3IndexToGeo(t, rootCell, rootRes+1, maxK, testOutput)
				}
			}

			if testOutput.ret2 > 0 || testOutput.ret0ValidationFailures > 0 ||
				testOutput.ret1ValidationFailures > 0 {
				t.Errorf("Resolution %d FAILED\nfailed0: %d\nfailed1: %d",
					resolution, testOutput.ret0ValidationFailures,
					testOutput.ret1ValidationFailures)
			}
		})
	}
}
