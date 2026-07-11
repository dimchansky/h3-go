// Tests ported from H3 v4.4.0: src/apps/testapps/testCoordIjInternal.c.
package h3

import "testing"

func Test_ijkToIj_zero(t *testing.T) {
	t.Parallel()

	ijk := coordIJK{0, 0, 0}
	ij := CoordIJ{0, 0}

	ijkToIj(&ijk, &ij)
	if ij.I != 0 {
		t.Errorf("ij.i zero: expected 0, got %v", ij.I)
	}
	if ij.J != 0 {
		t.Errorf("ij.j zero: expected 0, got %v", ij.J)
	}
}

func Test_ijkToIj_roundtrip(t *testing.T) {
	t.Parallel()

	for dir := centerDigit; dir < numDigits; dir++ {
		ijk := coordIJK{0, 0, 0}
		_neighbor(&ijk, dir)

		ij := CoordIJ{0, 0}
		ijkToIj(&ijk, &ij)

		recovered := coordIJK{0, 0, 0}
		if err := ijToIjk(&ij, &recovered); err != eSuccess {
			t.Errorf("ijToIjk failed for direction %v: %v", dir, err)
		}

		if !_ijkMatches(&ijk, &recovered) {
			t.Errorf("got same ijk coordinates back: original %v, recovered %v", ijk, recovered)
		}
	}
}
