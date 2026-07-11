// Tests ported from testCoordIjkInternal.c
package h3

import "testing"

func Test_unitIjkToDigit(t *testing.T) {
	t.Parallel()

	zero := coordIJK{0, 0, 0}
	i := coordIJK{1, 0, 0}
	outOfRange := coordIJK{2, 0, 0}
	unnormalizedZero := coordIJK{2, 2, 2}

	if got := _unitIjkToDigit(&zero); got != centerDigit {
		t.Errorf("Unit IJK to zero: expected %v, got %v", centerDigit, got)
	}

	if got := _unitIjkToDigit(&i); got != iAxesDigit {
		t.Errorf("Unit IJK to I axis: expected %v, got %v", iAxesDigit, got)
	}

	if got := _unitIjkToDigit(&outOfRange); got != invalidDigit {
		t.Errorf("Unit IJK out of range: expected %v, got %v", invalidDigit, got)
	}

	if got := _unitIjkToDigit(&unnormalizedZero); got != centerDigit {
		t.Errorf("Unnormalized unit IJK to zero: expected %v, got %v", centerDigit, got)
	}
}
