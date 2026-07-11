// Tests ported from H3 v4.4.0: src/apps/testapps/testH3NeighborRotations.c.
package h3

import (
	"fmt"
	"testing"
)

// TestOutput represents validation results from the h3NeighborRotations tests
// This corresponds to the TestOutput struct in the C code.
type TestOutput struct {
	ret0                   int
	ret0ValidationFailures int
	ret1                   int
	ret1ValidationFailures int
	ret2                   int
}

// doCell tests gridDiskUnsafe vs gridDiskDistancesSafe for a specific cell at various k values.
func doCell(t *testing.T, h h3Index, maxK int32, testOutput *TestOutput) {
	for k := int32(0); k < maxK; k++ {
		var maxSz int64
		err := maxGridDiskSize(k, &maxSz)
		if err != eSuccess {
			t.Fatalf("maxGridDiskSize failed: %v", err)
		}

		gridDiskInternalOutput := make([]h3Index, maxSz)
		gridDiskUnsafeOutput := make([]h3Index, maxSz)
		gridDiskInternalDistances := make([]int32, maxSz)

		// Call gridDiskDistancesSafe (equivalent to gridDiskDistancesInternal in C)
		err = gridDiskDistancesSafe(h, k, gridDiskInternalOutput, gridDiskInternalDistances)
		if err != eSuccess {
			t.Fatalf("gridDiskDistancesSafe failed: %v", err)
		}

		// Call gridDiskUnsafe
		gridDiskUnsafeFailed := gridDiskUnsafe(h, k, gridDiskUnsafeOutput)

		if gridDiskUnsafeFailed == eFailed {
			// TODO: Unreachable
			testOutput.ret2++
			continue
		} else if gridDiskUnsafeFailed == eSuccess {
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
		} else if gridDiskUnsafeFailed == ePentagon {
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

// recursiveH3IndexToGeo recursively generates all valid H3 indexes at a given resolution.
func recursiveH3IndexToGeo(t *testing.T, h h3Index, res int32, maxK int32, testOutput *TestOutput) {
	for d := int32(0); d < 7; d++ {
		current := h3SetIndexDigit(h, res, d)

		// skip the pentagonal deleted subsequence
		if _isBaseCellPentagon(int32(getBaseCell(current))) &&
			direction(_h3LeadingNonZeroDigit(current)) == kAxesDigit {
			continue
		}

		if res == getResolution(current) {
			doCell(t, current, maxK, testOutput)
		} else {
			recursiveH3IndexToGeo(t, current, res+1, maxK, testOutput)
		}
	}
}

// TestH3NeighborRotations tests gridDiskUnsafe vs. gridDiskDistancesSafe.
// Upstream CI (CMakeTests.cmake) runs testH3NeighborRotations three times
// with arguments 0, 1 and 2 (resolution) and the default maxK of 5; each run
// covers every index at that resolution across all 122 base cells.
func TestH3NeighborRotations(t *testing.T) {
	t.Parallel()

	for _, resolution := range []int32{0, 1, 2} {
		t.Run(fmt.Sprintf("resolution_%d", resolution), func(t *testing.T) {
			t.Parallel()

			const maxK = int32(5)
			testOutput := &TestOutput{0, 0, 0, 0, 0}

			for bc := int32(0); bc < numBaseCells; bc++ {
				rootCell := h3Index(h3CellMode << h3ModeOffset)
				rootCell = setBaseCell(rootCell, bc)

				if resolution == 0 {
					doCell(t, rootCell, maxK, testOutput)
				} else {
					rootRes := getResolution(rootCell)
					rootCell = setResolution(rootCell, resolution)
					recursiveH3IndexToGeo(t, rootCell, rootRes+1, maxK, testOutput)
				}
			}

			t.Logf("ret0: %d ret1: %d ret2: %d",
				testOutput.ret0, testOutput.ret1, testOutput.ret2)

			// ret2 should never occur, as it can only happen if we run over a pentagon
			if testOutput.ret2 > 0 || testOutput.ret0ValidationFailures > 0 ||
				testOutput.ret1ValidationFailures > 0 {
				t.Errorf("Resolution %d FAILED\nfailed0: %d\nfailed1: %d",
					resolution, testOutput.ret0ValidationFailures,
					testOutput.ret1ValidationFailures)
			}
		})
	}
}
