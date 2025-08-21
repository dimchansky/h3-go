//go:build c2go

package c2go

import "testing"

func Test_degsToRads_ParityWithC(t *testing.T) {
    cases := []float64{-720, -360, -180, -90, -1, 0, 1, 90, 180, 360, 720, 123.456}
    for _, v := range cases {
        goVal := degsToRads(v)
        cVal := degsToRadsC(v)
        if goVal != cVal {
            t.Fatalf("degsToRads parity mismatch: in=%g go=%.15g c=%.15g", v, goVal, cVal)
        }
    }
}
