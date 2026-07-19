//go:build cgo && c2go && h3v450

package h3

import "testing"

// Parity for reverseDirectedEdge against the exact 4.5.0 C
// implementation: success/double-reversal over every origin's edges at
// several resolutions, plus the upstream invalid-edge,
// partial-validation, and fuzz-regression inputs. All comparisons are
// bit-exact (pure index manipulation).

func Test_reverseDirectedEdge_parity(t *testing.T) {
	for base := int32(0); base < 122; base++ {
		var cell h3Index
		setH3Index(&cell, 0, base, 0)
		for _, res := range []int32{0, 2, 9} {
			child, err := cellToCenterChild(cell, res)
			if err != eSuccess {
				t.Fatalf("cellToCenterChild(base %d, res %d): %v", base, res, err)
			}
			var edges [6]h3Index
			if err := originToDirectedEdges(child, edges[:]); err != eSuccess {
				t.Fatalf("originToDirectedEdges(%x): %v", uint64(child), err)
			}
			for _, edge := range edges {
				var goOut h3Index
				goErr := reverseDirectedEdge(edge, &goOut)
				cOut, cErr := reverseDirectedEdgeC(edge)
				if goErr != cErr || goOut != cOut {
					t.Fatalf("reverseDirectedEdge(%x): Go=(%x,%v) C=(%x,%v)",
						uint64(edge), uint64(goOut), goErr, uint64(cOut), cErr)
				}
				if goErr != eSuccess {
					continue
				}
				// Double reversal.
				var goRev h3Index
				goErr = reverseDirectedEdge(goOut, &goRev)
				cRev, cErr := reverseDirectedEdgeC(cOut)
				if goErr != cErr || goRev != cRev {
					t.Fatalf("reverseDirectedEdge(rev %x): Go=(%x,%v) C=(%x,%v)",
						uint64(goOut), uint64(goRev), goErr, uint64(cRev), cErr)
				}
			}
		}
	}

	// Upstream error/partial-validation inputs.
	var sf h3Index
	if err := latLngToCell(&sfGeo, 9, &sf); err != eSuccess {
		t.Fatalf("latLngToCell: %v", err)
	}
	ring := make([]h3Index, 7)
	if err := gridRing(sf, 1, ring); err != eSuccess {
		t.Fatalf("gridRing: %v", err)
	}
	edge, err := cellsToDirectedEdge(sf, ring[0])
	if err != eSuccess {
		t.Fatalf("cellsToDirectedEdge: %v", err)
	}
	inputs := []h3Index{
		setReservedBits(edge, int32(invalidDigit)), // invalid reserved bits
		h3Null,                      // null index
		edge + 1,                    // invalid per isValidDirectedEdge, still succeeds
		h3Index(0x1001fff7ff2fbfff), // fuzz regression -> E_NOT_NEIGHBORS
	}
	for _, in := range inputs {
		var goOut h3Index
		goErr := reverseDirectedEdge(in, &goOut)
		cOut, cErr := reverseDirectedEdgeC(in)
		if goErr != cErr || (goErr == eSuccess && goOut != cOut) {
			t.Errorf("reverseDirectedEdge(%x): Go=(%x,%v) C=(%x,%v)",
				uint64(in), uint64(goOut), goErr, uint64(cOut), cErr)
		}
	}

	// Malformed edge whose origin is a pentagon and whose reserved
	// direction is the deleted K axis: destination recovery traverses
	// the deleted direction, so both sides must fail with E_PENTAGON
	// (the propagated branch behind the public ErrPentagon family).
	var pent h3Index
	setH3Index(&pent, 9, 4, 0) // base cell 4 is a pentagon; center children stay pentagons
	pentKEdge := setReservedBits(setMode(pent, h3DirectededgeMode), int32(kAxesDigit))
	var pentOut h3Index
	pentErr := reverseDirectedEdge(pentKEdge, &pentOut)
	cPentOut, cPentErr := reverseDirectedEdgeC(pentKEdge)
	if pentErr != cPentErr || pentErr != ePentagon {
		t.Errorf("reverseDirectedEdge(pentagon K-edge %x): Go=(%x,%v) C=(%x,%v), want ePentagon from both",
			uint64(pentKEdge), uint64(pentOut), pentErr, uint64(cPentOut), cPentErr)
	}
}
