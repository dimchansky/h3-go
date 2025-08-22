//go:build c2go

package c2go

import "testing"

func Test_h3index_getPentagons_ParityWithC(t *testing.T) {
    for _, res := range []int{0, 1, 5, 10, 15, -1, 16} {
        goOut, goErr := getPentagons(res)
        cOut, cErr := getPentagonsC(res)
        if goErr != cErr {
            t.Fatalf("getPentagons err mismatch res=%d: go=%d c=%d", res, goErr, cErr)
        }
        if goErr != _eSuccess {
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

