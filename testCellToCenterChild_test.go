// Tests ported from testCellToCenterChild.c
package h3

import (
	"testing"
)

func Test_centerChild_propertyTests(t *testing.T) {
	t.Parallel()

	var baseHex H3Index
	var baseCentroid LatLng
	setH3Index(&baseHex, 8, 4, 2)
	err := cellToLatLng(baseHex, &baseCentroid)
	if err != E_SUCCESS {
		t.Fatalf("Failed to get base centroid: %v", err)
	}

	for res := int32(0); res <= MAX_H3_RES-1; res++ {
		for childRes := res + 1; childRes <= MAX_H3_RES; childRes++ {
			var centroid LatLng
			var h3Index H3Index

			err := latLngToCell(&baseCentroid, res, &h3Index)
			if err != E_SUCCESS {
				t.Errorf("latLngToCell failed at res %d: %v", res, err)
				continue
			}

			err = cellToLatLng(h3Index, &centroid)
			if err != E_SUCCESS {
				t.Errorf("cellToLatLng failed at res %d: %v", res, err)
				continue
			}

			var geoChild H3Index
			err = latLngToCell(&centroid, childRes, &geoChild)
			if err != E_SUCCESS {
				t.Errorf("latLngToCell failed for geoChild at childRes %d: %v", childRes, err)
				continue
			}

			centerChild, err := cellToCenterChild(h3Index, childRes)
			if err != E_SUCCESS {
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
			if err != E_SUCCESS {
				t.Errorf("cellToParent failed: %v", err)
				continue
			}

			if parent != h3Index {
				t.Errorf("parent at original resolution should be initial index: expected 0x%x, got 0x%x", h3Index, parent)
			}
		}
	}
}

func Test_centerChild_sameRes(t *testing.T) {
	t.Parallel()

	var baseHex H3Index
	setH3Index(&baseHex, 8, 4, 2)
	res := getResolution(baseHex)

	child, err := cellToCenterChild(baseHex, res)
	if err != E_SUCCESS {
		t.Fatalf("cellToCenterChild failed: %v", err)
	}

	if child != baseHex {
		t.Errorf("center child at same resolution should return self: expected 0x%x, got 0x%x", baseHex, child)
	}
}

func Test_centerChild_invalidInputs(t *testing.T) {
	t.Parallel()

	var baseHex H3Index
	setH3Index(&baseHex, 8, 4, 2)
	res := getResolution(baseHex)

	// Test coarser resolution
	_, err := cellToCenterChild(baseHex, res-1)
	if err != E_RES_DOMAIN {
		t.Errorf("should fail at coarser resolution: expected E_RES_DOMAIN, got %v", err)
	}

	// Test negative resolution
	_, err = cellToCenterChild(baseHex, -1)
	if err != E_RES_DOMAIN {
		t.Errorf("should fail for negative resolution: expected E_RES_DOMAIN, got %v", err)
	}

	// Test beyond finest resolution
	_, err = cellToCenterChild(baseHex, MAX_H3_RES+1)
	if err != E_RES_DOMAIN {
		t.Errorf("should fail beyond finest resolution: expected E_RES_DOMAIN, got %v", err)
	}
}
