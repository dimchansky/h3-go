//go:build cgo

package h3

import "testing"

func Test_h3index_childPosToCell_ParityWithC(t *testing.T) {
	parents := []H3Index{0x821c07fffffffff}
	for _, parent := range parents {
		parentRes := getResolution(parent)
		for childRes := parentRes; childRes <= 6; childRes++ { // keep modest
			// Get children count to bound positions
			cnt, err := cellToChildrenSize(parent, childRes)
			if err != E_SUCCESS {
				t.Fatalf("childrenSize error: %d", err)
			}
			// test a few positions including edges
			tests := []int64{0}
			if cnt > 1 {
				tests = append(tests, 1, cnt-1)
			}
			for _, pos := range tests {
				goChild, goErr := childPosToCell(pos, parent, childRes)
				cChild, cErr := childPosToCellC(pos, parent, childRes)
				if goErr != H3Error(cErr) || goChild != cChild {
					t.Fatalf("childPosToCell mismatch parent=%x childRes=%d pos=%d: go(child=%x,err=%d) c(child=%x,err=%d)", uint64(parent), childRes, pos, uint64(goChild), goErr, uint64(cChild), cErr)
				}
			}
		}
	}
}
