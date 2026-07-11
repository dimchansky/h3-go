// Tests ported from H3 v4.4.0: src/apps/testapps/testCoordIjkInternal.c.
package h3

import (
	"math"
	"testing"
)

func Test_upAp7rChecked(t *testing.T) {
	t.Parallel()

	var ijk coordIJK

	_setIJK(&ijk, 0, 0, 0)
	if err := _upAp7rChecked(&ijk); err != eSuccess {
		t.Errorf("upAp7rChecked(0, 0, 0): expected eSuccess, got %v", err)
	}

	_setIJK(&ijk, math.MaxInt32, 0, 0)
	if err := _upAp7rChecked(&ijk); err != eFailed {
		t.Errorf("i + i overflows: expected eFailed, got %v", err)
	}

	_setIJK(&ijk, 0, math.MaxInt32, 0)
	if err := _upAp7rChecked(&ijk); err != eFailed {
		t.Errorf("j + j overflows: expected eFailed, got %v", err)
	}

	_setIJK(&ijk, 0, math.MaxInt32/2, 0)
	if err := _upAp7rChecked(&ijk); err != eFailed {
		t.Errorf("3 * j overflows: expected eFailed, got %v", err)
	}

	_setIJK(&ijk, math.MaxInt32/2, math.MaxInt32/3, 0)
	if err := _upAp7rChecked(&ijk); err != eFailed {
		t.Errorf("(i * 2) + j overflows: expected eFailed, got %v", err)
	}

	// This input should be invalid because i < 0
	_setIJK(&ijk, -2, math.MaxInt32/3, 0)
	if err := _upAp7rChecked(&ijk); err != eFailed {
		t.Errorf("(j * 3) - 1 overflows: expected eFailed, got %v", err)
	}

	// This input should be invalid because j < 0
	_setIJK(&ijk, -1, 0, 0)
	if err := _upAp7rChecked(&ijk); err != eSuccess {
		t.Errorf("i < 0 succeeds: expected eSuccess, got %v", err)
	}
}
