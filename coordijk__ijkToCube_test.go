// Tests ported from testCoordIjInternal.c
package h3

import "testing"

func Test_ijkToCube_roundtrip(t *testing.T) {
	t.Parallel()

	for dir := centerDigit; dir < numDigits; dir++ {
		ijk := coordIJK{0, 0, 0}
		_neighbor(&ijk, dir)
		original := coordIJK{ijk.I, ijk.J, ijk.K}

		ijkToCube(&ijk)
		cubeToIjk(&ijk)

		if !_ijkMatches(&ijk, &original) {
			t.Errorf("got same ijk coordinates back: original %v, recovered %v, direction %v", original, ijk, dir)
		}
	}
}
