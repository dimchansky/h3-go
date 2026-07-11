// Tests ported from testBaseCellsInternal.c
package h3

import (
	"testing"
)

func TestBaseCellToCCWrot60(t *testing.T) {
	t.Parallel()
	// a few random spot-checks
	if _baseCellToCCWrot60(16, 0) != 0 {
		t.Error("Expected rotation 0 for base cell 16, face 0")
	}
	if _baseCellToCCWrot60(32, 0) != 3 {
		t.Error("Expected rotation 3 for base cell 32, face 0")
	}
	if _baseCellToCCWrot60(7, 3) != 1 {
		t.Error("Expected rotation 1 for base cell 7, face 3")
	}
}

func TestBaseCellToCCWrot60_invalid(t *testing.T) {
	t.Parallel()
	if _baseCellToCCWrot60(16, 42) != invalidRotations {
		t.Error("Should return invalid rotation for invalid face")
	}
	if _baseCellToCCWrot60(16, -1) != invalidRotations {
		t.Error("Should return invalid rotation for invalid face (negative)")
	}
	if _baseCellToCCWrot60(1, 0) != invalidRotations {
		t.Error("Should return invalid rotation for base cell not appearing on face")
	}
}

func TestIsBaseCellPentagon_invalid(t *testing.T) {
	t.Parallel()
	if _isBaseCellPentagon(-1) != false {
		t.Error("isBaseCellPentagon should handle negative")
	}
}
