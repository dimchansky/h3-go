// Tests ported from testGridRingUnsafe.c
package h3

import (
	"testing"
)

func Test_gridRingUnsafe_negativeK(t *testing.T) {
	t.Parallel()

	sf := LatLng{Lat: 0.659966917655, Lng: 2*3.14159 - 2.1364398519396}
	var sfHex h3Index
	err := latLngToCell(&sf, 9, &sfHex)
	if err != eSuccess {
		t.Fatalf("latLngToCell failed: %v", err)
	}

	k0 := []h3Index{0}
	result := gridRingUnsafe(sfHex, -1, k0)
	if result != eDomain {
		t.Errorf("Expected eDomain for negative k, got %v", result)
	}
}

func Test_gridRingUnsafe_identityGridRing(t *testing.T) {
	t.Parallel()

	sf := LatLng{Lat: 0.659966917655, Lng: 2*3.14159 - 2.1364398519396}
	var sfHex h3Index
	err := latLngToCell(&sf, 9, &sfHex)
	if err != eSuccess {
		t.Fatalf("latLngToCell failed: %v", err)
	}

	k0 := []h3Index{0}
	result := gridRingUnsafe(sfHex, 0, k0)
	if result != eSuccess {
		t.Fatalf("gridRingUnsafe failed: %v", result)
	}
	if k0[0] != sfHex {
		t.Errorf("Expected identity k-ring to be the origin cell, got %x expected %x", k0[0], sfHex)
	}
}

func Test_gridRingUnsafe_ring1(t *testing.T) {
	t.Parallel()

	sf := LatLng{Lat: 0.659966917655, Lng: 2*3.14159 - 2.1364398519396}
	var sfHex h3Index
	err := latLngToCell(&sf, 9, &sfHex)
	if err != eSuccess {
		t.Fatalf("latLngToCell failed: %v", err)
	}

	k1 := []h3Index{0, 0, 0, 0, 0, 0}
	expectedK1 := []h3Index{0x89283080ddbffff, 0x89283080c37ffff,
		0x89283080c27ffff, 0x89283080d53ffff,
		0x89283080dcfffff, 0x89283080dc3ffff}

	result := gridRingUnsafe(sfHex, 1, k1)
	if result != eSuccess {
		t.Fatalf("gridRingUnsafe failed: %v", result)
	}

	for i := 0; i < 6; i++ {
		if k1[i] == 0 {
			t.Errorf("Index %d is not populated", i)
			continue
		}

		inList := 0
		for j := 0; j < 6; j++ {
			if k1[i] == expectedK1[j] {
				inList++
			}
		}
		if inList != 1 {
			t.Errorf("Index %x at position %d not found in expected set (found %d times)", k1[i], i, inList)
		}
	}
}

func Test_gridRingUnsafe_ring2(t *testing.T) {
	t.Parallel()

	sf := LatLng{Lat: 0.659966917655, Lng: 2*3.14159 - 2.1364398519396}
	var sfHex h3Index
	err := latLngToCell(&sf, 9, &sfHex)
	if err != eSuccess {
		t.Fatalf("latLngToCell failed: %v", err)
	}

	k2 := []h3Index{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	expectedK2 := []h3Index{
		0x89283080ca7ffff, 0x89283080cafffff, 0x89283080c33ffff,
		0x89283080c23ffff, 0x89283080c2fffff, 0x89283080d5bffff,
		0x89283080d43ffff, 0x89283080d57ffff, 0x89283080d1bffff,
		0x89283080dc7ffff, 0x89283080dd7ffff, 0x89283080dd3ffff}

	result := gridRingUnsafe(sfHex, 2, k2)
	if result != eSuccess {
		t.Fatalf("gridRingUnsafe failed: %v", result)
	}

	for i := 0; i < 12; i++ {
		if k2[i] == 0 {
			t.Errorf("Index %d is not populated", i)
			continue
		}

		inList := 0
		for j := 0; j < 12; j++ {
			if k2[i] == expectedK2[j] {
				inList++
			}
		}
		if inList != 1 {
			t.Errorf("Index %x at position %d not found in expected set (found %d times)", k2[i], i, inList)
		}
	}
}

func Test_gridRingUnsafe_nearPentagonRing1(t *testing.T) {
	t.Parallel()

	nearPentagon := h3Index(0x837405fffffffff)
	kp1 := []h3Index{0, 0, 0, 0, 0, 0}

	result := gridRingUnsafe(nearPentagon, 1, kp1)
	if result != ePentagon {
		t.Errorf("Expected ePentagon when hitting a pentagon, got %v", result)
	}
}

func Test_gridRingUnsafe_nearPentagonRing2(t *testing.T) {
	t.Parallel()

	nearPentagon := h3Index(0x837405fffffffff)
	kp2 := []h3Index{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	result := gridRingUnsafe(nearPentagon, 2, kp2)
	if result != ePentagon {
		t.Errorf("Expected ePentagon when hitting a pentagon, got %v", result)
	}
}

func Test_gridRingUnsafe_onPentagon(t *testing.T) {
	t.Parallel()

	var nearPentagon h3Index
	setH3Index(&nearPentagon, 0, 4, 0)
	kp2 := []h3Index{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	result := gridRingUnsafe(nearPentagon, 2, kp2)
	if result != ePentagon {
		t.Errorf("Expected ePentagon when starting at a pentagon, got %v", result)
	}
}

func Test_gridRingUnsafe_matches_gridDiskDistancesSafe(t *testing.T) {
	t.Parallel()

	for res := int32(0); res < 2; res++ {
		for i := int32(0); i < numBaseCells; i++ {
			var bc h3Index
			setH3Index(&bc, 0, i, 0)

			childrenSz, err := uncompactCellsSize([]h3Index{bc}, 1, res)
			if err != eSuccess {
				t.Fatalf("uncompactCellsSize failed: %v", err)
			}

			children := make([]h3Index, childrenSz)
			err = uncompactCells([]h3Index{bc}, 1, children, childrenSz, res)
			if err != eSuccess {
				t.Fatalf("uncompactCells failed: %v", err)
			}

			for j := int64(0); j < childrenSz; j++ {
				if children[j] == 0 {
					continue
				}

				for k := int32(0); k < 3; k++ {
					var kSz int64
					err := maxGridDiskSize(k, &kSz)
					if err != eSuccess {
						t.Fatalf("maxGridDiskSize failed: %v", err)
					}

					var ringSize int64
					err = _maxGridRingSize(k, &ringSize)
					if err != eSuccess {
						t.Fatalf("_maxGridRingSize failed: %v", err)
					}

					ring := make([]h3Index, ringSize)
					failed := gridRingUnsafe(children[j], k, ring)

					if failed == eSuccess {
						internalNeighbors := make([]h3Index, kSz)
						internalDistances := make([]int32, kSz)

						err := gridDiskDistancesSafe(children[j], k, internalNeighbors, internalDistances)
						if err != eSuccess {
							t.Fatalf("gridDiskDistancesSafe failed: %v", err)
						}

						found := 0
						internalFound := 0

						for iRing := int64(0); iRing < ringSize; iRing++ {
							if ring[iRing] != 0 {
								found++

								for iInternal := int64(0); iInternal < kSz; iInternal++ {
									if internalNeighbors[iInternal] == ring[iRing] {
										internalFound++

										if internalDistances[iInternal] != k {
											t.Errorf("Ring and internal disagree on distance: ring k=%d, internal distance=%d", k, internalDistances[iInternal])
										}

										break
									}
								}

								if found != internalFound {
									t.Errorf("Ring and internal implementations produce different output: found=%d, internalFound=%d", found, internalFound)
								}
							}
						}
					}
				}
			}
		}
	}
}
