// Tests ported from H3 v4.4.0: src/apps/testapps/testGridDisksUnsafe.c.
package h3

import (
	"testing"
)

func Test_gridDisksUnsafe_identityGridDiskCells(t *testing.T) {
	t.Parallel()

	sf := LatLng{Lat: 0.659966917655, Lng: 2*3.14159 - 2.1364398519396}
	var sfHex h3Index
	err := latLngToCell(&sf, 9, &sfHex)
	if err != eSuccess {
		t.Fatalf("latLngToCell failed: %v", err)
	}

	sfHexPtr := []h3Index{sfHex}
	k0 := make([]h3Index, 1)

	result := gridDisksUnsafe(sfHexPtr, 0, k0)
	if result != eSuccess {
		t.Fatalf("gridDisksUnsafe failed: %v", result)
	}

	if k0[0] != sfHex {
		t.Errorf("Expected identity k-ring to be the origin cell, got %x expected %x", k0[0], sfHex)
	}
}

func Test_gridDisksUnsafe_ring1of1(t *testing.T) {
	t.Parallel()

	k1 := []h3Index{0x89283080ddbffff, 0x89283080c37ffff, 0x89283080c27ffff,
		0x89283080d53ffff, 0x89283080dcfffff, 0x89283080dc3ffff}

	allKrings := make([]h3Index, 42) // 6 origins * 7 cells each (1+6)

	result := gridDisksUnsafe(k1, 1, allKrings)
	if result != eSuccess {
		t.Fatalf("gridDisksUnsafe failed: %v", result)
	}

	for i := 0; i < 42; i++ {
		if allKrings[i] == 0 {
			t.Errorf("Index %d should be populated but is 0", i)
		}
		if i%7 == 0 {
			index := i / 7
			if k1[index] != allKrings[i] {
				t.Errorf("The beginning of segment %d should be the correct hexagon %x, got %x", index, k1[index], allKrings[i])
			}
		}
	}
}

func Test_gridDisksUnsafe_ring2of1(t *testing.T) {
	t.Parallel()

	k1 := []h3Index{0x89283080ddbffff, 0x89283080c37ffff, 0x89283080c27ffff,
		0x89283080d53ffff, 0x89283080dcfffff, 0x89283080dc3ffff}

	// 6 origins * 19 cells each (1 + 6 + 12)
	allKrings2 := make([]h3Index, 6*(1+6+12))

	result := gridDisksUnsafe(k1, 2, allKrings2)
	if result != eSuccess {
		t.Fatalf("gridDisksUnsafe failed: %v", result)
	}

	for i := 0; i < 6*(1+6+12); i++ {
		if allKrings2[i] == 0 {
			t.Errorf("Index %d should be populated but is 0", i)
		}

		if i%(1+6+12) == 0 {
			index := i / (1 + 6 + 12)
			if k1[index] != allKrings2[i] {
				t.Errorf("The beginning of segment %d should be the correct hexagon %x, got %x", index, k1[index], allKrings2[i])
			}
		}
	}
}

func Test_gridDisksUnsafe_failed(t *testing.T) {
	t.Parallel()

	withPentagon := []h3Index{0x8029fffffffffff, 0x801dfffffffffff}

	// 2 origins * 7 cells each (1 + 6)
	allKrings := make([]h3Index, 2*(1+6))

	result := gridDisksUnsafe(withPentagon, 1, allKrings)
	if result != ePentagon {
		t.Errorf("Expected ePentagon error on gridDisksUnsafe, got %v", result)
	}
}

func Test_gridDisksUnsafe_invalid_k(t *testing.T) {
	t.Parallel()

	k1 := []h3Index{0x89283080ddbffff, 0x89283080c37ffff, 0x89283080c27ffff,
		0x89283080d53ffff, 0x89283080dcfffff, 0x89283080dc3ffff}

	// We need to provide a valid output slice to test the k validation
	out := make([]h3Index, 1)
	result := gridDisksUnsafe(k1, -1, out)
	if result != eDomain {
		t.Errorf("Expected eDomain for invalid k, got %v", result)
	}
}
