// Tests ported from H3 v4.4.0: src/apps/testapps/testCoordIjkInternal.c.
package h3

import "testing"

func Test_neighbor(t *testing.T) {
	t.Parallel()

	ijk := coordIJK{0, 0, 0}

	zero := coordIJK{0, 0, 0}
	i := coordIJK{1, 0, 0}

	_neighbor(&ijk, centerDigit)
	if !_ijkMatches(&ijk, &zero) {
		t.Errorf("Center neighbor is self: expected %v, got %v", zero, ijk)
	}

	_neighbor(&ijk, iAxesDigit)
	if !_ijkMatches(&ijk, &i) {
		t.Errorf("I neighbor as expected: expected %v, got %v", i, ijk)
	}

	_neighbor(&ijk, invalidDigit)
	if !_ijkMatches(&ijk, &i) {
		t.Errorf("Invalid neighbor is self: expected %v, got %v", i, ijk)
	}
}
