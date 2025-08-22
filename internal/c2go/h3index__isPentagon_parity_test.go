//go:build c2go

package c2go

import (
	"testing"
)

func Test_h3index_isPentagon_ParityWithC(t *testing.T) {
	// Build a few base cell res=0 indexes, pentagon and non-pentagon
	pentBases := []int{4, 14, 24}
	nonPentBases := []int{0, 1, 2}
	// res=0
	for _, b := range pentBases {
		h := setH3IndexC(0, b, 7)
		if isPentagon(h) != isPentagonC(h) {
			t.Fatalf("isPentagon mismatch at res0 pent base=%d", b)
		}
	}
	for _, b := range nonPentBases {
		h := setH3IndexC(0, b, 7)
		if isPentagon(h) != isPentagonC(h) {
			t.Fatalf("isPentagon mismatch at res0 non-pent base=%d", b)
		}
	}
	// res>0, center path (leading digit 0) keeps pentagon if base is pentagon
	for _, b := range pentBases {
		h := setH3IndexC(3, b, 0) // all digits 0
		if isPentagon(h) != isPentagonC(h) {
			t.Fatalf("isPentagon mismatch at res3 center pent base=%d", b)
		}
	}
	// res>0, non-center digit should not be pentagon
	for _, b := range pentBases {
		h := setH3IndexC(3, b, 1) // digits 1
		if isPentagon(h) != isPentagonC(h) {
			t.Fatalf("isPentagon mismatch at res3 non-center pent base=%d", b)
		}
	}
}
