//go:build cgo

package c2go

import (
	"testing"
)

func Test_h3index_isPentagon_ParityWithC(t *testing.T) {
	// Build a few base cell res=0 indexes, pentagon and non-pentagon
	pentBases := []int32{4, 14, 24}
	nonPentBases := []int32{0, 1, 2}
	// res=0
	for _, b := range pentBases {
		h := setH3IndexC(0, b, 7)
		goResult := isPentagon(h)
		cResult := isPentagonC(h)
		if goResult != cResult {
			t.Fatalf("isPentagon mismatch at res0 pent base=%d: Go=%t (%v), C=%v", b, goResult, goResult, cResult)
		}
	}
	for _, b := range nonPentBases {
		h := setH3IndexC(0, b, 7)
		goResult := isPentagon(h)
		cResult := isPentagonC(h)
		if goResult != cResult {
			t.Fatalf("isPentagon mismatch at res0 non-pent base=%d: Go=%t (%v), C=%v", b, goResult, goResult, cResult)
		}
	}
	// res>0, center path (leading digit 0) keeps pentagon if base is pentagon
	for _, b := range pentBases {
		h := setH3IndexC(3, b, 0) // all digits 0
		goResult := isPentagon(h)
		cResult := isPentagonC(h)
		if goResult != cResult {
			t.Fatalf("isPentagon mismatch at res3 center pent base=%d: Go=%t (%v), C=%v", b, goResult, goResult, cResult)
		}
	}
	// res>0, non-center digit should not be pentagon
	for _, b := range pentBases {
		h := setH3IndexC(3, b, 1) // digits 1
		goResult := isPentagon(h)
		cResult := isPentagonC(h)
		if goResult != cResult {
			t.Fatalf("isPentagon mismatch at res3 non-center pent base=%d: Go=%t (%v), C=%v", b, goResult, goResult, cResult)
		}
	}
}
