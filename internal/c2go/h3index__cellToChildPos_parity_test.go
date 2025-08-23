//go:build cgo

package c2go

import "testing"

func Test_h3index_cellToChildPos_ParityWithC(t *testing.T) {
	hs := []H3Index{0x8928308280fffff, 0x821c07fffffffff}
	for _, child := range hs {
		for parentRes := 0; parentRes <= getResolution(child); parentRes++ {
			goPos, goErr := cellToChildPos(child, parentRes)
			cPos, cErr := cellToChildPosC(child, parentRes)
			if goErr != H3Error(cErr) || goPos != cPos {
				t.Fatalf("cellToChildPos mismatch child=%x parentRes=%d: go(pos=%d,err=%d) c(pos=%d,err=%d)", uint64(child), parentRes, goPos, goErr, cPos, cErr)
			}
		}
	}
}
