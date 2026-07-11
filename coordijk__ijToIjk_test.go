// Tests ported from H3 v4.4.0: src/apps/testapps/testCoordIjInternal.c.
package h3

import "testing"

func Test_ijToIjk_zero(t *testing.T) {
	t.Parallel()

	ij := CoordIJ{0, 0}
	ijk := coordIJK{0, 0, 0}

	if err := ijToIjk(&ij, &ijk); err != eSuccess {
		t.Errorf("ijToIjk failed: %v", err)
	}

	if ijk.I != 0 {
		t.Errorf("ijk.i zero: expected 0, got %v", ijk.I)
	}
	if ijk.J != 0 {
		t.Errorf("ijk.j zero: expected 0, got %v", ijk.J)
	}
	if ijk.K != 0 {
		t.Errorf("ijk.k zero: expected 0, got %v", ijk.K)
	}
}
