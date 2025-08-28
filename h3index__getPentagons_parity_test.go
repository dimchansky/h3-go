//go:build cgo

package h3

import "testing"

func Test_h3index_getPentagons_ParityWithC(t *testing.T) {
	for _, res := range []int32{0, 1, 5, 10, 15, -1, 16} {
		// Use dst-buffer with insufficient capacity first
		var dst []H3Index
		goOut, goErr := getPentagons(dst, res)
		cOut, cErr := getPentagonsC(res)
		if goErr != H3Error(cErr) {
			t.Fatalf("getPentagons err mismatch res=%d: go=%d c=%d", res, goErr, cErr)
		}
		if goErr != E_SUCCESS {
			continue
		}
		if len(goOut) != len(cOut) {
			t.Fatalf("getPentagons length mismatch res=%d: go=%d c=%d", res, len(goOut), len(cOut))
		}
		for i := range goOut {
			if goOut[i] != cOut[i] {
				t.Fatalf("getPentagons mismatch res=%d idx=%d: go=%x c=%x", res, i, uint64(goOut[i]), uint64(cOut[i]))
			}
		}
		// Now reuse a dst-buffer with sufficient capacity and ensure identical output
		dst2 := make([]H3Index, 0, NUM_PENTAGONS)
		goOut2, goErr2 := getPentagons(dst2, res)
		if goErr2 != goErr {
			t.Fatalf("getPentagons err mismatch (cap reuse) res=%d: go2=%d go=%d", res, goErr2, goErr)
		}
		if goErr2 == E_SUCCESS {
			if len(goOut2) != len(goOut) {
				t.Fatalf("getPentagons length mismatch (cap reuse) res=%d: %d vs %d", res, len(goOut2), len(goOut))
			}
			for i := range goOut2 {
				if goOut2[i] != goOut[i] {
					t.Fatalf("getPentagons mismatch (cap reuse) res=%d idx=%d: out2=%x out=%x", res, i, uint64(goOut2[i]), uint64(goOut[i]))
				}
			}
		}
	}
}
