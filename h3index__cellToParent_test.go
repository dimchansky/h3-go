// Tests ported from testCellToParent.c
package h3

import (
	"math"
	"testing"
)

func TestCellToParent_ancestorsForEachRes(t *testing.T) {
	t.Parallel()
	sf := LatLng{
		Lat: Angle(0.659966917655),
		Lng: Angle(2*math.Pi - 2.1364398519396),
	}

	for res := int32(1); res < 15; res++ {
		for step := int32(0); step < res; step++ {
			var child H3Index
			childErr := latLngToCell(&sf, res, &child)
			if childErr != E_SUCCESS {
				t.Fatalf("latLngToCell failed for res=%d: %v", res, childErr)
			}

			parent, parentErr := cellToParent(child, res-step)
			if parentErr != E_SUCCESS {
				t.Fatalf("cellToParent failed for child at res=%d, parent res=%d: %v", res, res-step, parentErr)
			}

			var comparisonParent H3Index
			comparisonErr := latLngToCell(&sf, res-step, &comparisonParent)
			if comparisonErr != E_SUCCESS {
				t.Fatalf("latLngToCell failed for comparison parent at res=%d: %v", res-step, comparisonErr)
			}

			if parent != comparisonParent {
				t.Errorf("Got unexpected parent for res=%d->%d: got %x, want %x", res, res-step, parent, comparisonParent)
			}
		}
	}
}

func TestCellToParent_invalidInputs(t *testing.T) {
	t.Parallel()
	sf := LatLng{
		Lat: Angle(0.659966917655),
		Lng: Angle(2*math.Pi - 2.1364398519396),
	}

	var child H3Index
	childErr := latLngToCell(&sf, 5, &child)
	if childErr != E_SUCCESS {
		t.Fatalf("latLngToCell failed: %v", childErr)
	}

	// Higher resolution fails
	_, err := cellToParent(child, 6)
	if err != E_RES_MISMATCH {
		t.Errorf("Expected E_RES_MISMATCH for higher resolution, got %v", err)
	}

	// Invalid resolution fails (negative)
	_, err = cellToParent(child, -1)
	if err != E_RES_DOMAIN {
		t.Errorf("Expected E_RES_DOMAIN for negative resolution, got %v", err)
	}

	// Invalid resolution fails (15 when child is at 5)
	_, err = cellToParent(child, 15)
	if err != E_RES_MISMATCH {
		t.Errorf("Expected E_RES_MISMATCH for resolution 15, got %v", err)
	}

	// Invalid resolution fails (16)
	_, err = cellToParent(child, 16)
	if err != E_RES_DOMAIN {
		t.Errorf("Expected E_RES_DOMAIN for resolution 16, got %v", err)
	}
}