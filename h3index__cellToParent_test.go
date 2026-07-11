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
			var child h3Index
			childErr := latLngToCell(&sf, res, &child)
			if childErr != eSuccess {
				t.Fatalf("latLngToCell failed for res=%d: %v", res, childErr)
			}

			parent, parentErr := cellToParent(child, res-step)
			if parentErr != eSuccess {
				t.Fatalf("cellToParent failed for child at res=%d, parent res=%d: %v", res, res-step, parentErr)
			}

			var comparisonParent h3Index
			comparisonErr := latLngToCell(&sf, res-step, &comparisonParent)
			if comparisonErr != eSuccess {
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

	var child h3Index
	childErr := latLngToCell(&sf, 5, &child)
	if childErr != eSuccess {
		t.Fatalf("latLngToCell failed: %v", childErr)
	}

	// Higher resolution fails
	_, err := cellToParent(child, 6)
	if err != eResMismatch {
		t.Errorf("Expected eResMismatch for higher resolution, got %v", err)
	}

	// Invalid resolution fails (negative)
	_, err = cellToParent(child, -1)
	if err != eResDomain {
		t.Errorf("Expected eResDomain for negative resolution, got %v", err)
	}

	// Invalid resolution fails (15 when child is at 5)
	_, err = cellToParent(child, 15)
	if err != eResMismatch {
		t.Errorf("Expected eResMismatch for resolution 15, got %v", err)
	}

	// Invalid resolution fails (16)
	_, err = cellToParent(child, 16)
	if err != eResDomain {
		t.Errorf("Expected eResDomain for resolution 16, got %v", err)
	}
}
