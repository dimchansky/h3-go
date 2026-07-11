// Tests ported from H3 v4.4.0: src/apps/testapps/testGridRingInternal.c.
package h3

import (
	"testing"
)

func Test_gridRingInternal_identityGridRing(t *testing.T) {
	t.Parallel()

	sf := LatLng{Lat: 0.659966917655, Lng: 2*3.14159 - 2.1364398519396}
	var sfHex h3Index
	err := latLngToCell(&sf, 9, &sfHex)
	if err != eSuccess {
		t.Fatalf("latLngToCell failed: %v", err)
	}

	k0 := []h3Index{0}
	result := _gridRingInternal(sfHex, 0, k0)
	if result != eSuccess {
		t.Fatalf("_gridRingInternal failed: %v", result)
	}
	if k0[0] != sfHex {
		t.Errorf("Expected identity k-ring to be the origin cell, got %x expected %x", k0[0], sfHex)
	}
}

func Test_gridRingInternal_negativeK(t *testing.T) {
	t.Parallel()

	sf := LatLng{Lat: 0.659966917655, Lng: 2*3.14159 - 2.1364398519396}
	var sfHex h3Index
	err := latLngToCell(&sf, 9, &sfHex)
	if err != eSuccess {
		t.Fatalf("latLngToCell failed: %v", err)
	}

	k0 := []h3Index{0}
	result := _gridRingInternal(sfHex, -1, k0)
	if result != eDomain {
		t.Errorf("Expected eDomain for negative k, got %v", result)
	}
}

func Test_gridRingInternal_gridDiskInvalid(t *testing.T) {
	t.Parallel()

	const k = 1000
	var kSz int64
	err := maxGridDiskSize(k, &kSz)
	if err != eSuccess {
		t.Fatalf("maxGridDiskSize failed: %v", err)
	}

	neighbors := make([]h3Index, kSz)
	result := _gridRingInternal(0x7fffffffffffffff, k, neighbors)
	if result != eCellInvalid {
		t.Errorf("Expected eCellInvalid for invalid input, got %v", result)
	}
}

func Test_gridRingInternal_gridDiskInvalidDigit(t *testing.T) {
	t.Parallel()

	const k = 2
	var kSz int64
	err := maxGridDiskSize(k, &kSz)
	if err != eSuccess {
		t.Fatalf("maxGridDiskSize failed: %v", err)
	}

	neighbors := make([]h3Index, kSz)
	result := _gridRingInternal(0x4d4b00fe5c5c3030, k, neighbors)
	if result != eCellInvalid {
		t.Errorf("Expected eCellInvalid for invalid input, got %v", result)
	}
}
