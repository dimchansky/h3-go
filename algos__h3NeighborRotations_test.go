// Tests ported from testGridDiskInternal.c
package h3

import (
	"testing"
)

func Test_h3NeighborRotations_identity(t *testing.T) {
	t.Parallel()

	origin := H3Index(0x811d7ffffffffff)
	rotations := int32(0)
	var out H3Index

	result := h3NeighborRotations(origin, CENTER_DIGIT, &rotations, &out)
	if result != E_SUCCESS {
		t.Fatalf("h3NeighborRotations failed: %v", result)
	}

	if out != origin {
		t.Errorf("Moving to self should go to self, got %x expected %x", out, origin)
	}

	if rotations != 0 {
		t.Errorf("Expected rotations to be 0, got %d", rotations)
	}
}

func Test_h3NeighborRotations_rotationsOverflow(t *testing.T) {
	t.Parallel()

	var origin H3Index
	setH3Index(&origin, 0, 0, int32(CENTER_DIGIT))

	// A multiple of 6, so effectively no rotation. Very close to INT32_MAX.
	rotations := int32(2147483646)
	var out H3Index

	result := h3NeighborRotations(origin, K_AXES_DIGIT, &rotations, &out)
	if result != E_SUCCESS {
		t.Fatalf("h3NeighborRotations failed: %v", result)
	}

	var expected H3Index
	// Determined by looking at the base cell table
	setH3Index(&expected, 0, 1, int32(CENTER_DIGIT))

	if out != expected {
		t.Errorf("Expected neighbor %x, got %x", expected, out)
	}

	if rotations != 5 {
		t.Errorf("Expected rotations value 5, got %d", rotations)
	}
}

func Test_h3NeighborRotations_rotationsOverflow2(t *testing.T) {
	t.Parallel()

	var origin H3Index
	setH3Index(&origin, 0, 4, int32(CENTER_DIGIT))

	// This modulo 6 is 1.
	rotations := int32(INT32_MAX)
	var out H3Index

	// This will try to move in the K direction off of origin,
	// which will be adjusted to the IK direction.
	result := h3NeighborRotations(origin, JK_AXES_DIGIT, &rotations, &out)
	if result != E_SUCCESS {
		t.Fatalf("h3NeighborRotations failed: %v", result)
	}

	var expected H3Index
	// Determined by looking at the base cell table
	setH3Index(&expected, 0, 0, int32(CENTER_DIGIT))

	if out != expected {
		t.Errorf("Expected neighbor %x, got %x", expected, out)
	}

	// 1 (original value) + 4 (newRotations for IK direction) + 1 (applied
	// when adjusting to the IK direction) = 6, 6 modulo 6 = 0
	if rotations != 0 {
		t.Errorf("Expected rotations value 0, got %d", rotations)
	}
}

func Test_h3NeighborRotations_invalid(t *testing.T) {
	t.Parallel()

	origin := H3Index(0x811d7ffffffffff)
	rotations := int32(0)
	var out H3Index

	result := h3NeighborRotations(origin, -1, &rotations, &out)
	if result != E_FAILED {
		t.Errorf("Expected E_FAILED for invalid direction -1, got %v", result)
	}

	result = h3NeighborRotations(origin, 7, &rotations, &out)
	if result != E_FAILED {
		t.Errorf("Expected E_FAILED for invalid direction 7, got %v", result)
	}

	result = h3NeighborRotations(origin, 100, &rotations, &out)
	if result != E_FAILED {
		t.Errorf("Expected E_FAILED for invalid direction 100, got %v", result)
	}
}

func Test_cwOffsetPent(t *testing.T) {
	t.Parallel()

	// Try to find a case where h3NeighborRotations would not pass the
	// cwOffsetPent check, and would hit a line marked as unreachable.

	// To do this, we need to find a case that would move from one
	// non-pentagon base cell into the deleted k-subsequence of a pentagon
	// base cell, and neither of the cwOffsetPent values are the original
	// base cell's face.

	for pentagon := int32(0); pentagon < NUM_BASE_CELLS; pentagon++ {
		if !_isBaseCellPentagon(pentagon) {
			continue
		}

		for neighbor := int32(0); neighbor < NUM_BASE_CELLS; neighbor++ {
			// Get the face for the neighbor base cell
			neighborFace := baseCellData[neighbor].HomeFijk.Face

			// Only direction 2 needs to be checked, because that is the
			// only direction where we can move from digit 2 to digit 1, and
			// into the deleted k subsequence.
			neighborPentagon := _getBaseCellNeighbor(neighbor, J_AXES_DIGIT)

			if neighborPentagon != pentagon || _baseCellIsCwOffset(pentagon, neighborFace) {
				// This is the expected condition - either we don't move to this pentagon
				// or the cwOffsetPent check would pass
				continue
			}

			// If we get here, it means we found a case that would fail the cwOffsetPent check
			t.Errorf("cwOffsetPent check would fail: neighbor %d (face %d) -> pentagon %d",
				neighbor, neighborFace, pentagon)
		}
	}
}
