// Tests ported from testCoordIjkInternal.c
package h3

import "testing"

func Test_neighbor(t *testing.T) {
	t.Parallel()

	ijk := CoordIJK{0, 0, 0}

	zero := CoordIJK{0, 0, 0}
	i := CoordIJK{1, 0, 0}

	_neighbor(&ijk, CENTER_DIGIT)
	if !_ijkMatches(&ijk, &zero) {
		t.Errorf("Center neighbor is self: expected %v, got %v", zero, ijk)
	}

	_neighbor(&ijk, I_AXES_DIGIT)
	if !_ijkMatches(&ijk, &i) {
		t.Errorf("I neighbor as expected: expected %v, got %v", i, ijk)
	}

	_neighbor(&ijk, INVALID_DIGIT)
	if !_ijkMatches(&ijk, &i) {
		t.Errorf("Invalid neighbor is self: expected %v, got %v", i, ijk)
	}
}
