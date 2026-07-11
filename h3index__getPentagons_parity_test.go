//go:build cgo && c2go

package h3

import "testing"

func Test_h3index_getPentagons_ParityWithC(t *testing.T) {
	for _, res := range []int32{0, 1, 5, 10, 15, -1, 16} {
		// Use dst-buffer with insufficient capacity first
		goOut := make([]H3Index, NUM_PENTAGONS)
		goErr := getPentagons(res, goOut)
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
	}
}
