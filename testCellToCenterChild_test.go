// Tests ported from H3 v4.4.0: src/apps/testapps/testCellToCenterChild.c.
package h3

import (
	"testing"
)

func Test_centerChild_propertyTests(t *testing.T) {
	t.Parallel()

	var baseHex h3Index
	var baseCentroid LatLng
	setH3Index(&baseHex, 8, 4, 2)
	err := cellToLatLng(baseHex, &baseCentroid)
	if err != eSuccess {
		t.Fatalf("Failed to get base centroid: %v", err)
	}

	for res := int32(0); res <= maxH3Res-1; res++ {
		for childRes := res + 1; childRes <= maxH3Res; childRes++ {
			var centroid LatLng
			var idx h3Index

			err := latLngToCell(&baseCentroid, res, &idx)
			if err != eSuccess {
				t.Errorf("latLngToCell failed at res %d: %v", res, err)
				continue
			}

			err = cellToLatLng(idx, &centroid)
			if err != eSuccess {
				t.Errorf("cellToLatLng failed at res %d: %v", res, err)
				continue
			}

			var geoChild h3Index
			err = latLngToCell(&centroid, childRes, &geoChild)
			if err != eSuccess {
				t.Errorf("latLngToCell failed for geoChild at childRes %d: %v", childRes, err)
				continue
			}

			centerChild, err := cellToCenterChild(idx, childRes)
			if err != eSuccess {
				t.Errorf("cellToCenterChild failed at res %d childRes %d: %v", res, childRes, err)
				continue
			}

			if centerChild != geoChild {
				t.Errorf("center child should be same as indexed centroid at child resolution: res=%d, childRes=%d, centerChild=0x%x, geoChild=0x%x", res, childRes, centerChild, geoChild)
			}

			if getResolution(centerChild) != childRes {
				t.Errorf("center child should have correct resolution: expected %d, got %d", childRes, getResolution(centerChild))
			}

			parent, err := cellToParent(centerChild, res)
			if err != eSuccess {
				t.Errorf("cellToParent failed: %v", err)
				continue
			}

			if parent != idx {
				t.Errorf("parent at original resolution should be initial index: expected 0x%x, got 0x%x", idx, parent)
			}
		}
	}
}

func Test_centerChild_sameRes(t *testing.T) {
	t.Parallel()

	var baseHex h3Index
	setH3Index(&baseHex, 8, 4, 2)
	res := getResolution(baseHex)

	child, err := cellToCenterChild(baseHex, res)
	if err != eSuccess {
		t.Fatalf("cellToCenterChild failed: %v", err)
	}

	if child != baseHex {
		t.Errorf("center child at same resolution should return self: expected 0x%x, got 0x%x", baseHex, child)
	}
}

func Test_centerChild_invalidInputs(t *testing.T) {
	t.Parallel()

	var baseHex h3Index
	setH3Index(&baseHex, 8, 4, 2)
	res := getResolution(baseHex)

	// Test coarser resolution
	_, err := cellToCenterChild(baseHex, res-1)
	if err != eResDomain {
		t.Errorf("should fail at coarser resolution: expected eResDomain, got %v", err)
	}

	// Test negative resolution
	_, err = cellToCenterChild(baseHex, -1)
	if err != eResDomain {
		t.Errorf("should fail for negative resolution: expected eResDomain, got %v", err)
	}

	// Test beyond finest resolution
	_, err = cellToCenterChild(baseHex, maxH3Res+1)
	if err != eResDomain {
		t.Errorf("should fail beyond finest resolution: expected eResDomain, got %v", err)
	}
}
