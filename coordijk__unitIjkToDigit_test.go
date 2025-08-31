// Tests ported from testCoordIjkInternal.c
package h3

import "testing"

func Test_unitIjkToDigit(t *testing.T) {
	t.Parallel()
	
	zero := CoordIJK{0, 0, 0}
	i := CoordIJK{1, 0, 0}
	outOfRange := CoordIJK{2, 0, 0}
	unnormalizedZero := CoordIJK{2, 2, 2}

	if got := _unitIjkToDigit(&zero); got != CENTER_DIGIT {
		t.Errorf("Unit IJK to zero: expected %v, got %v", CENTER_DIGIT, got)
	}
	
	if got := _unitIjkToDigit(&i); got != I_AXES_DIGIT {
		t.Errorf("Unit IJK to I axis: expected %v, got %v", I_AXES_DIGIT, got)
	}
	
	if got := _unitIjkToDigit(&outOfRange); got != INVALID_DIGIT {
		t.Errorf("Unit IJK out of range: expected %v, got %v", INVALID_DIGIT, got)
	}
	
	if got := _unitIjkToDigit(&unnormalizedZero); got != CENTER_DIGIT {
		t.Errorf("Unnormalized unit IJK to zero: expected %v, got %v", CENTER_DIGIT, got)
	}
}