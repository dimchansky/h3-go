// Tests ported from testGridRing.c
package h3

import (
	"testing"
)

func Test_gridRing_identityGridRing(t *testing.T) {
	t.Parallel()

	sf := LatLng{Lat: 0.659966917655, Lng: 2*3.14159 - 2.1364398519396}
	var sfHex H3Index
	err := latLngToCell(&sf, 9, &sfHex)
	if err != E_SUCCESS {
		t.Fatalf("latLngToCell failed: %v", err)
	}

	k0 := []H3Index{0}
	result := gridRing(sfHex, 0, k0)
	if result != E_SUCCESS {
		t.Fatalf("gridRing failed: %v", result)
	}
	if k0[0] != sfHex {
		t.Errorf("Expected identity k-ring to be the origin cell, got %x expected %x", k0[0], sfHex)
	}
}

func Test_gridRing_ring1(t *testing.T) {
	t.Parallel()

	sf := LatLng{Lat: 0.659966917655, Lng: 2*3.14159 - 2.1364398519396}
	var sfHex H3Index
	err := latLngToCell(&sf, 9, &sfHex)
	if err != E_SUCCESS {
		t.Fatalf("latLngToCell failed: %v", err)
	}

	k1 := []H3Index{0, 0, 0, 0, 0, 0}
	expectedK1 := []H3Index{0x89283080ddbffff, 0x89283080c37ffff,
		0x89283080c27ffff, 0x89283080d53ffff,
		0x89283080dcfffff, 0x89283080dc3ffff}

	result := gridRing(sfHex, 1, k1)
	if result != E_SUCCESS {
		t.Fatalf("gridRing failed: %v", result)
	}

	for i := 0; i < 6; i++ {
		if k1[i] == 0 {
			t.Errorf("index %d is not populated", i)
			continue
		}

		inList := 0
		for j := 0; j < 6; j++ {
			if k1[i] == expectedK1[j] {
				inList++
			}
		}
		if inList != 1 {
			t.Errorf("index %x not found exactly once in expected set", k1[i])
		}
	}
}

func Test_gridRing_ring2(t *testing.T) {
	t.Parallel()

	sf := LatLng{Lat: 0.659966917655, Lng: 2*3.14159 - 2.1364398519396}
	var sfHex H3Index
	err := latLngToCell(&sf, 9, &sfHex)
	if err != E_SUCCESS {
		t.Fatalf("latLngToCell failed: %v", err)
	}

	k2 := []H3Index{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	expectedK2 := []H3Index{
		0x89283080ca7ffff, 0x89283080cafffff, 0x89283080c33ffff,
		0x89283080c23ffff, 0x89283080c2fffff, 0x89283080d5bffff,
		0x89283080d43ffff, 0x89283080d57ffff, 0x89283080d1bffff,
		0x89283080dc7ffff, 0x89283080dd7ffff, 0x89283080dd3ffff}

	result := gridRing(sfHex, 2, k2)
	if result != E_SUCCESS {
		t.Fatalf("gridRing failed: %v", result)
	}

	for i := 0; i < 12; i++ {
		if k2[i] == 0 {
			t.Errorf("index %d is not populated", i)
			continue
		}

		inList := 0
		for j := 0; j < 12; j++ {
			if k2[i] == expectedK2[j] {
				inList++
			}
		}
		if inList != 1 {
			t.Errorf("index %x not found exactly once in expected set", k2[i])
		}
	}
}

func Test_gridRing_PolarPentagon_k1(t *testing.T) {
	t.Parallel()

	var polar H3Index
	setH3Index(&polar, 0, 4, 0)
	k2 := []H3Index{0, 0, 0, 0, 0, 0}
	expectedK2 := []H3Index{0x8007fffffffffff, 0x8001fffffffffff,
		0x8011fffffffffff, 0x801ffffffffffff,
		0x8019fffffffffff, 0}

	result := gridRing(polar, 1, k2)
	if result != E_SUCCESS {
		t.Fatalf("gridRing failed: %v", result)
	}

	k2present := 0
	for i := 0; i < 6; i++ {
		if k2[i] != 0 {
			k2present++
			inList := 0
			for j := 0; j < 6; j++ {
				if k2[i] == expectedK2[j] {
					inList++
				}
			}
			if inList != 1 {
				t.Errorf("index %x not found exactly once in expected set", k2[i])
			}
		}
	}
	if k2present != 5 {
		t.Errorf("pentagon should have 5 neighbors in k-ring 1, got %d", k2present)
	}
}

func Test_gridRing_PolarPentagon_res1_k1(t *testing.T) {
	t.Parallel()

	var polar H3Index
	setH3Index(&polar, 1, 4, 0)
	k2 := []H3Index{0, 0, 0, 0, 0, 0}
	expectedK2 := []H3Index{0x81093ffffffffff, 0x81097ffffffffff,
		0x8108fffffffffff, 0x8108bffffffffff,
		0x8109bffffffffff, 0}

	result := gridRing(polar, 1, k2)
	if result != E_SUCCESS {
		t.Fatalf("gridRing failed: %v", result)
	}

	k2present := 0
	for i := 0; i < 6; i++ {
		if k2[i] != 0 {
			k2present++
			inList := 0
			for j := 0; j < 6; j++ {
				if k2[i] == expectedK2[j] {
					inList++
				}
			}
			if inList != 1 {
				t.Errorf("index %x not found exactly once in expected set", k2[i])
			}
		}
	}
	if k2present != 5 {
		t.Errorf("pentagon should have 5 neighbors in k-ring 1, got %d", k2present)
	}
}

func Test_gridRing_PolarPentagon_res1_k3(t *testing.T) {
	t.Parallel()

	var polar H3Index
	setH3Index(&polar, 1, 4, 0)
	k2 := make([]H3Index, 18)
	expectedK2 := []H3Index{0x811fbffffffffff,
		0x81003ffffffffff,
		0x81183ffffffffff,
		0x8111bffffffffff,
		0x81067ffffffffff,
		0x811e7ffffffffff,
		0x8101bffffffffff,
		0x81107ffffffffff,
		0x81063ffffffffff,
		0x811e3ffffffffff,
		0x8119bffffffffff,
		0x81103ffffffffff,
		0x81007ffffffffff,
		0x81187ffffffffff,
		0x8107bffffffffff,
		0,
		0,
		0}

	result := gridRing(polar, 3, k2)
	if result != E_SUCCESS {
		t.Fatalf("gridRing failed: %v", result)
	}

	k2present := 0
	for i := 0; i < 18; i++ {
		if k2[i] != 0 {
			k2present++
			inList := 0
			for j := 0; j < 18; j++ {
				if k2[i] == expectedK2[j] {
					inList++
				}
			}
			if inList != 1 {
				t.Errorf("index %x not found exactly once in expected set", k2[i])
			}
		}
	}
	if k2present != 15 {
		t.Errorf("pentagon should have 15 neighbors in k-ring 3, got %d", k2present)
	}
}

func Test_gridRing_Pentagon_res1_k4(t *testing.T) {
	t.Parallel()

	var pent H3Index
	setH3Index(&pent, 1, 14, 0)
	k2 := make([]H3Index, 24)
	expectedK2 := []H3Index{
		0x81227ffffffffff,
		0x81293ffffffffff,
		0x8136bffffffffff,
		0x81167ffffffffff,
		0x81477ffffffffff,
		0x810dbffffffffff,
		0x81473ffffffffff,
		0x81237ffffffffff,
		0x81127ffffffffff,
		0x8126bffffffffff,
		0x81177ffffffffff,
		0x810d3ffffffffff,
		0x8150fffffffffff,
		0x8102fffffffffff,
		0x8129bffffffffff,
		0x8102bffffffffff,
		0x81507ffffffffff,
		0x8136fffffffffff,
		0x8127bffffffffff,
		0x81137ffffffffff,
		0,
		0,
		0,
		0,
	}

	result := gridRing(pent, 4, k2)
	if result != E_SUCCESS {
		t.Fatalf("gridRing failed: %v", result)
	}

	k2present := 0
	for i := 0; i < 24; i++ {
		if k2[i] != 0 {
			k2present++
			inList := 0
			for j := 0; j < 24; j++ {
				if k2[i] == expectedK2[j] {
					inList++
				}
			}
			if inList != 1 {
				t.Errorf("index %x not found exactly once in expected set", k2[i])
			}
		}
	}
	if k2present != 20 {
		t.Errorf("pentagon should have 20 neighbors in k-ring 4, got %d", k2present)
	}
}

func Test_maxGridRingSize_invalid(t *testing.T) {
	t.Parallel()

	var sz int64
	result := _maxGridRingSize(-1, &sz)
	if result != E_DOMAIN {
		t.Errorf("Expected E_DOMAIN for negative k, got %v", result)
	}
}

func Test_maxGridRingSize_identity(t *testing.T) {
	t.Parallel()

	var sz int64
	result := _maxGridRingSize(0, &sz)
	if result != E_SUCCESS {
		t.Fatalf("_maxGridRingSize failed: %v", result)
	}
	if sz != 1 {
		t.Errorf("k = 0 should return 1, got %d", sz)
	}
}

func Test_maxGridRingSize(t *testing.T) {
	t.Parallel()

	var sz int64
	result := _maxGridRingSize(2, &sz)
	if result != E_SUCCESS {
		t.Fatalf("_maxGridRingSize failed: %v", result)
	}
	if sz != 12 {
		t.Errorf("k = 2 should return 12, got %d", sz)
	}
}

func Test_gridRing_matches_gridDiskDistancesSafe(t *testing.T) {
	t.Parallel()

	for res := int32(0); res < 2; res++ {
		for i := int32(0); i < NUM_BASE_CELLS; i++ {
			var bc H3Index
			setH3Index(&bc, 0, i, 0)

			childrenSz, err := uncompactCellsSize([]H3Index{bc}, 1, res)
			if err != E_SUCCESS {
				t.Fatalf("uncompactCellsSize failed: %v", err)
			}

			children := make([]H3Index, childrenSz)
			result := uncompactCells([]H3Index{bc}, 1, children, childrenSz, res)
			if result != E_SUCCESS {
				t.Fatalf("uncompactCells failed: %v", result)
			}

			for j := int64(0); j < childrenSz; j++ {
				if children[j] == 0 {
					continue
				}

				for k := int32(0); k < 3; k++ {
					var kSz int64
					err := maxGridDiskSize(k, &kSz)
					if err != E_SUCCESS {
						t.Fatalf("maxGridDiskSize failed: %v", err)
					}

					var ringSize int64
					err = _maxGridRingSize(k, &ringSize)
					if err != E_SUCCESS {
						t.Fatalf("_maxGridRingSize failed: %v", err)
					}

					ring := make([]H3Index, ringSize)
					result := gridRing(children[j], k, ring)
					if result != E_SUCCESS {
						t.Fatalf("gridRing failed: %v", result)
					}

					internalNeighbors := make([]H3Index, kSz)
					internalDistances := make([]int32, kSz)
					result = gridDiskDistancesSafe(children[j], k, internalNeighbors, internalDistances)
					if result != E_SUCCESS {
						t.Fatalf("gridDiskDistancesSafe failed: %v", result)
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
										t.Errorf("Ring and internal disagree on distance: ring has cell %x at distance %d, but internal has it at distance %d", ring[iRing], k, internalDistances[iInternal])
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
