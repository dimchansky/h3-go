//go:build cgo && c2go

package h3

import (
	"testing"
)

func Test_adjustOverageClassII_CriticalDifference(t *testing.T) {
	// Test the specific case where Go and C implementations differ
	// This is the termination point where C returns noOverage but Go returns newFace

	testCase := struct {
		name             string
		input            faceIJK
		expectedC        faceIJK
		expectedOverageC overage
	}{
		name:             "termination_point",
		input:            faceIJK{Face: 5, Coord: coordIJK{I: 715827882, J: 0, K: 1431655770}},
		expectedC:        faceIJK{Face: 5, Coord: coordIJK{I: 715827882, J: 0, K: 1431655770}}, // C doesn't change it
		expectedOverageC: noOverage,                                                            // C returns 0 (noOverage)
	}

	res := int32(3)
	pentLeading4Go := false
	substrateGo := true
	pentLeading4C := int32(0)
	substrateC := int32(1)

	t.Run(testCase.name, func(t *testing.T) {
		// Test Go implementation
		fijkGo := testCase.input
		overageGo := _adjustOverageClassII(&fijkGo, res, pentLeading4Go, substrateGo)

		// Test C implementation
		fijkC := testCase.input
		overageC := _adjustOverageClassIIC(&fijkC, res, pentLeading4C, substrateC)

		t.Logf("Input: Face=%d, Coord={%d,%d,%d}",
			testCase.input.Face, testCase.input.Coord.I, testCase.input.Coord.J, testCase.input.Coord.K)
		t.Logf("Go result:  overage=%d, Face=%d, Coord={%d,%d,%d}",
			int(overageGo), fijkGo.Face, fijkGo.Coord.I, fijkGo.Coord.J, fijkGo.Coord.K)
		t.Logf("C result:   overage=%d, Face=%d, Coord={%d,%d,%d}",
			int(overageC), fijkC.Face, fijkC.Coord.I, fijkC.Coord.J, fijkC.Coord.K)

		// Check what C actually returns
		if overageC != testCase.expectedOverageC {
			t.Errorf("C overage unexpected: expected %d, got %d", int32(testCase.expectedOverageC), int32(overageC))
		}
		if fijkC != testCase.expectedC {
			t.Errorf("C result unexpected: expected Face=%d Coord={%d,%d,%d}, got Face=%d Coord={%d,%d,%d}",
				testCase.expectedC.Face, testCase.expectedC.Coord.I, testCase.expectedC.Coord.J, testCase.expectedC.Coord.K,
				fijkC.Face, fijkC.Coord.I, fijkC.Coord.J, fijkC.Coord.K)
		}

		// Compare Go vs C - this is the critical difference
		if overageC != overageGo {
			t.Logf("*** CRITICAL DIFFERENCE CONFIRMED! ***")
			t.Logf("C overage: %d, Go overage: %d", int32(overageC), int32(overageGo))
		}
		if fijkC != fijkGo {
			t.Logf("*** COORDINATE DIFFERENCE CONFIRMED! ***")
			t.Logf("C result:  Face=%d, Coord={%d,%d,%d}", fijkC.Face, fijkC.Coord.I, fijkC.Coord.J, fijkC.Coord.K)
			t.Logf("Go result: Face=%d, Coord={%d,%d,%d}", fijkGo.Face, fijkGo.Coord.I, fijkGo.Coord.J, fijkGo.Coord.K)
		}
	})
}
