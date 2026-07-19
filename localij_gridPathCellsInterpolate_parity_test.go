//go:build cgo && c2go && h3v450

package h3

import "testing"

// Direct parity for gridPathCellsInterpolate (file-static in the 4.5.0
// localij.c, reached via the same-TU wrapper) and for gridPathCells over
// the exact 4.5.0 behavioral cases: the pentagon pair that succeeds only
// through reverse interpolation and the pinned pair where both attempts
// fail. All outputs are discrete (cells, error codes) and compare
// exactly.

func Test_gridPathCellsInterpolate_parity(t *testing.T) {
	// Deterministic same-res neighbor pairs around every listed base
	// cell's center child (base 58 is a pentagon, exercising distorted
	// neighborhoods). Short distances keep gridDistance computable.
	type pair struct{ start, end h3Index }
	var testPairs []pair
	for _, base := range []int32{2, 15, 37, 58, 79, 100, 121} {
		var cell h3Index
		setH3Index(&cell, 0, base, 0)
		child, err := cellToCenterChild(cell, 5)
		if err != eSuccess {
			t.Fatalf("cellToCenterChild(base %d): %v", base, err)
		}
		var ring [7]h3Index
		if err := gridRing(child, 1, ring[:]); err != eSuccess {
			t.Fatalf("gridRing(base %d): %v", base, err)
		}
		testPairs = append(testPairs,
			pair{child, ring[0]},
			pair{child, ring[3]},
			pair{ring[0], ring[3]},
		)
	}
	for _, p := range testPairs {
		var distance int64
		if err := gridDistance(p.start, p.end, &distance); err != eSuccess {
			t.Fatalf("gridDistance(%x,%x): %v", uint64(p.start), uint64(p.end), err)
		}
		if distance == 0 {
			continue
		}
		for _, dir := range []struct{ offset, step int64 }{
			{0, 1},         // forward fill
			{distance, -1}, // reverse fill
		} {
			goOut := make([]h3Index, distance+1)
			cOut := make([]h3Index, distance+1)
			goErr := gridPathCellsInterpolate(p.start, p.end, distance, goOut, dir.offset, dir.step)
			cErr := gridPathCellsInterpolateC(p.start, p.end, distance, cOut, dir.offset, dir.step)
			if goErr != cErr {
				t.Fatalf("interpolate(%x,%x,off=%d,step=%d): Go err=%v C err=%v",
					uint64(p.start), uint64(p.end), dir.offset, dir.step, goErr, cErr)
			}
			if goErr != eSuccess {
				continue
			}
			for i := range goOut {
				if goOut[i] != cOut[i] {
					t.Fatalf("interpolate(%x,%x,off=%d,step=%d)[%d]: Go=%x C=%x",
						uint64(p.start), uint64(p.end), dir.offset, dir.step, i,
						uint64(goOut[i]), uint64(cOut[i]))
				}
			}
		}
	}
}

func Test_gridPathCellsInterpolate_failingCases_parity(t *testing.T) {
	// Direct helper parity on the two 4.5.0 behavioral gridPathCells
	// pairs, calling the interpolation attempts exactly as
	// gridPathCells issues them: the start-anchored forward fill, then
	// the end-anchored reverse fill. For the pentagon pair the forward
	// attempt fails and the reverse attempt succeeds; for the pinned
	// pair both attempts fail. Error codes compare exactly on every
	// attempt; the successful reverse fill compares every written cell.
	cases := []struct {
		name            string
		start, end      h3Index
		reverseSucceeds bool
	}{
		{"pentagonReverseInterpolation", 0x820807fffffffff, 0x8208e7fffffffff, true},
		{"knownFailureNotCoveredByReverseInterpolation", 0x8411b61ffffffff, 0x84016d3ffffffff, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var distance int64
			if err := gridDistance(tc.start, tc.end, &distance); err != eSuccess {
				t.Fatalf("gridDistance: %v", err)
			}
			goOut := make([]h3Index, distance+1)
			cOut := make([]h3Index, distance+1)

			goFwd := gridPathCellsInterpolate(tc.start, tc.end, distance, goOut, 0, 1)
			cFwd := gridPathCellsInterpolateC(tc.start, tc.end, distance, cOut, 0, 1)
			if goFwd != cFwd {
				t.Fatalf("forward interpolate: Go err=%v C err=%v", goFwd, cFwd)
			}
			if goFwd == eSuccess {
				t.Fatalf("forward interpolate unexpectedly succeeded (case exists to exercise the failure)")
			}

			goRev := gridPathCellsInterpolate(tc.end, tc.start, distance, goOut, distance, -1)
			cRev := gridPathCellsInterpolateC(tc.end, tc.start, distance, cOut, distance, -1)
			if goRev != cRev {
				t.Fatalf("reverse interpolate: Go err=%v C err=%v", goRev, cRev)
			}
			if got := goRev == eSuccess; got != tc.reverseSucceeds {
				t.Fatalf("reverse interpolate: err=%v, want success=%v", goRev, tc.reverseSucceeds)
			}
			if goRev != eSuccess {
				return
			}
			for i := range goOut {
				if goOut[i] != cOut[i] {
					t.Fatalf("reverse interpolate[%d]: Go=%x C=%x", i, uint64(goOut[i]), uint64(cOut[i]))
				}
			}
		})
	}
}

func Test_gridPathCells_450Cases_parity(t *testing.T) {
	cases := []struct {
		name       string
		start, end h3Index
	}{
		// Succeeds only through the end-anchored reverse interpolation.
		{"pentagonReverseInterpolation", 0x820807fffffffff, 0x8208e7fffffffff},
		// Both interpolation attempts fail; error codes must match.
		{"knownFailureNotCoveredByReverseInterpolation", 0x8411b61ffffffff, 0x84016d3ffffffff},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var size int64
			if err := gridPathCellsSize(tc.start, tc.end, &size); err != eSuccess {
				t.Fatalf("gridPathCellsSize: %v", err)
			}
			goOut := make([]h3Index, size)
			cOut := make([]h3Index, size)
			goErr := gridPathCells(goOut, tc.start, tc.end)
			cErr := _gridPathCellsC(tc.start, tc.end, cOut)
			if goErr != cErr {
				t.Fatalf("gridPathCells(%x,%x): Go err=%v C err=%v",
					uint64(tc.start), uint64(tc.end), goErr, cErr)
			}
			if goErr != eSuccess {
				return
			}
			for i := range goOut {
				if goOut[i] != cOut[i] {
					t.Fatalf("gridPathCells(%x,%x)[%d]: Go=%x C=%x",
						uint64(tc.start), uint64(tc.end), i, uint64(goOut[i]), uint64(cOut[i]))
				}
			}
		})
	}
}
