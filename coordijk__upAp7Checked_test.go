// Tests ported from testCoordIjkInternal.c
package h3

import (
	"math"
	"testing"
)

func Test_upAp7Checked(t *testing.T) {
	t.Parallel()
	
	var ijk CoordIJK
	
	_setIJK(&ijk, 0, 0, 0)
	if err := _upAp7Checked(&ijk); err != E_SUCCESS {
		t.Errorf("upAp7Checked(0, 0, 0): expected E_SUCCESS, got %v", err)
	}
	
	_setIJK(&ijk, math.MaxInt32, 0, 0)
	if err := _upAp7Checked(&ijk); err != E_FAILED {
		t.Errorf("i + i overflows: expected E_FAILED, got %v", err)
	}
	
	_setIJK(&ijk, math.MaxInt32/2, 0, 0)
	if err := _upAp7Checked(&ijk); err != E_FAILED {
		t.Errorf("i * 3 overflows: expected E_FAILED, got %v", err)
	}
	
	_setIJK(&ijk, 0, math.MaxInt32, 0)
	if err := _upAp7Checked(&ijk); err != E_FAILED {
		t.Errorf("j + j overflows: expected E_FAILED, got %v", err)
	}
	
	// This input should be invalid because j < 0
	_setIJK(&ijk, math.MaxInt32/3, -2, 0)
	if err := _upAp7Checked(&ijk); err != E_FAILED {
		t.Errorf("(i * 3) - j overflows: expected E_FAILED, got %v", err)
	}
	
	_setIJK(&ijk, math.MaxInt32/3, math.MaxInt32/2, 0)
	if err := _upAp7Checked(&ijk); err != E_FAILED {
		t.Errorf("i + (j * 2) overflows: expected E_FAILED, got %v", err)
	}
	
	// This input should be invalid because j < 0
	_setIJK(&ijk, -1, 0, 0)
	if err := _upAp7Checked(&ijk); err != E_SUCCESS {
		t.Errorf("i < 0 succeeds: expected E_SUCCESS, got %v", err)
	}
}