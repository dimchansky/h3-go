// Tests ported from testGridDisk.c
package h3

import (
	"testing"
)

// Test helper that checks gridDisk output matches gridDiskDistancesSafe output.
func gridDisk_equals_gridDiskDistancesSafe_assertions(t *testing.T, h3 H3Index) {
	for k := int32(0); k < 3; k++ {
		var kSz int64
		err := maxGridDiskSize(k, &kSz)
		if err != E_SUCCESS {
			t.Fatalf("maxGridDiskSize failed: %v", err)
		}

		neighbors := make([]H3Index, kSz)
		distances := make([]int32, kSz)
		err = gridDiskDistances(h3, k, neighbors, distances)
		if err != E_SUCCESS {
			t.Fatalf("gridDiskDistances failed: %v", err)
		}

		internalNeighbors := make([]H3Index, kSz)
		internalDistances := make([]int32, kSz)
		err = gridDiskDistancesSafe(h3, k, internalNeighbors, internalDistances)
		if err != E_SUCCESS {
			t.Fatalf("gridDiskDistancesSafe failed: %v", err)
		}

		found := 0
		internalFound := 0
		for iNeighbor := int64(0); iNeighbor < kSz; iNeighbor++ {
			if neighbors[iNeighbor] != 0 {
				found++

				for iInternal := int64(0); iInternal < kSz; iInternal++ {
					if internalNeighbors[iInternal] == neighbors[iNeighbor] {
						internalFound++

						if distances[iNeighbor] != internalDistances[iInternal] {
							t.Errorf("External and internal disagree on distance at neighbor %d: external=%d, internal=%d",
								iNeighbor, distances[iNeighbor], internalDistances[iInternal])
						}
						break
					}
				}
			}
		}

		if found != internalFound {
			t.Errorf("External and internal implementations produce different output: found=%d, internalFound=%d",
				found, internalFound)
		}
	}
}

func TestGridDisk0(t *testing.T) {
	t.Parallel()

	// Convert from radians to degrees
	sf := LatLng{Lat: Angle(0.659966917655), Lng: Angle(2*3.14159 - 2.1364398519396)}
	var sfHex0 H3Index
	err := latLngToCell(&sf, 0, &sfHex0)
	if err != E_SUCCESS {
		t.Fatalf("latLngToCell failed: %v", err)
	}

	k1 := make([]H3Index, 7)
	k1Dist := make([]int32, 7)
	expectedK1 := []H3Index{0x8029fffffffffff, 0x801dfffffffffff,
		0x8013fffffffffff, 0x8027fffffffffff,
		0x8049fffffffffff, 0x8051fffffffffff,
		0x8037fffffffffff}

	err = gridDiskDistances(sfHex0, 1, k1, k1Dist)
	if err != E_SUCCESS {
		t.Fatalf("gridDiskDistances failed: %v", err)
	}

	for i := 0; i < 7; i++ {
		if k1[i] == 0 {
			t.Errorf("index at position %d is not populated", i)
		}
		inList := 0
		for j := 0; j < 7; j++ {
			if k1[i] == expectedK1[j] {
				expectedDistance := int32(0)
				if k1[i] != sfHex0 {
					expectedDistance = 1
				}
				if k1Dist[i] != expectedDistance {
					t.Errorf("distance is not as expected: got %d, expected %d", k1Dist[i], expectedDistance)
				}
				inList++
			}
		}
		if inList != 1 {
			t.Errorf("index %x found %d times in expected set, expected 1", k1[i], inList)
		}
	}
}

func TestGridDisk0_PolarPentagon(t *testing.T) {
	t.Parallel()

	var polar H3Index
	setH3Index(&polar, 0, 4, 0)
	k2 := make([]H3Index, 7)
	k2Dist := make([]int32, 7)
	expectedK2 := []H3Index{0x8009fffffffffff,
		0x8007fffffffffff,
		0x8001fffffffffff,
		0x8011fffffffffff,
		0x801ffffffffffff,
		0x8019fffffffffff,
		0}

	err := gridDiskDistances(polar, 1, k2, k2Dist)
	if err != E_SUCCESS {
		t.Fatalf("gridDiskDistances failed: %v", err)
	}

	k2present := 0
	for i := 0; i < 7; i++ {
		if k2[i] != 0 {
			k2present++
			inList := 0
			for j := 0; j < 7; j++ {
				if k2[i] == expectedK2[j] {
					expectedDistance := int32(0)
					if k2[i] != polar {
						expectedDistance = 1
					}
					if k2Dist[i] != expectedDistance {
						t.Errorf("distance is not as expected: got %d, expected %d", k2Dist[i], expectedDistance)
					}
					inList++
				}
			}
			if inList != 1 {
				t.Errorf("index %x found %d times in expected set, expected 1", k2[i], inList)
			}
		}
	}
	if k2present != 6 {
		t.Errorf("pentagon has %d neighbors, expected 6 (5 neighbors + itself)", k2present)
	}
}

func TestGridDisk1_PolarPentagon(t *testing.T) {
	t.Parallel()

	var polar H3Index
	setH3Index(&polar, 1, 4, 0)
	k2 := make([]H3Index, 7)
	k2Dist := make([]int32, 7)
	expectedK2 := []H3Index{0x81083ffffffffff,
		0x81093ffffffffff,
		0x81097ffffffffff,
		0x8108fffffffffff,
		0x8108bffffffffff,
		0x8109bffffffffff,
		0}

	err := gridDiskDistances(polar, 1, k2, k2Dist)
	if err != E_SUCCESS {
		t.Fatalf("gridDiskDistances failed: %v", err)
	}

	k2present := 0
	for i := 0; i < 7; i++ {
		if k2[i] != 0 {
			k2present++
			inList := 0
			for j := 0; j < 7; j++ {
				if k2[i] == expectedK2[j] {
					expectedDistance := int32(0)
					if k2[i] != polar {
						expectedDistance = 1
					}
					if k2Dist[i] != expectedDistance {
						t.Errorf("distance is not as expected: got %d, expected %d", k2Dist[i], expectedDistance)
					}
					inList++
				}
			}
			if inList != 1 {
				t.Errorf("index %x found %d times in expected set, expected 1", k2[i], inList)
			}
		}
	}
	if k2present != 6 {
		t.Errorf("pentagon has %d neighbors, expected 6 (5 neighbors + itself)", k2present)
	}
}

func TestGridDisk1_PolarPentagon_k3(t *testing.T) {
	t.Parallel()

	var polar H3Index
	setH3Index(&polar, 1, 4, 0)
	k2 := make([]H3Index, 37)
	k2Dist := make([]int32, 37)
	expectedK2 := []H3Index{0x81013ffffffffff,
		0x811fbffffffffff,
		0x81193ffffffffff,
		0x81097ffffffffff,
		0x81003ffffffffff,
		0x81183ffffffffff,
		0x8111bffffffffff,
		0x81077ffffffffff,
		0x811f7ffffffffff,
		0x81067ffffffffff,
		0x81093ffffffffff,
		0x811e7ffffffffff,
		0x81083ffffffffff,
		0x81117ffffffffff,
		0x8101bffffffffff,
		0x81107ffffffffff,
		0x81073ffffffffff,
		0x811f3ffffffffff,
		0x81063ffffffffff,
		0x8108fffffffffff,
		0x811e3ffffffffff,
		0x8119bffffffffff,
		0x81113ffffffffff,
		0x81017ffffffffff,
		0x81103ffffffffff,
		0x8109bffffffffff,
		0x81197ffffffffff,
		0x81007ffffffffff,
		0x8108bffffffffff,
		0x81187ffffffffff,
		0x8107bffffffffff,
		0,
		0,
		0,
		0,
		0,
		0}
	expectedK2Dist := []int32{2, 3, 2, 1, 3, 3, 3, 2, 2, 3, 1, 3, 0,
		2, 3, 3, 2, 2, 3, 1, 3, 3, 2, 2, 3, 1,
		2, 3, 1, 3, 3, 0, 0, 0, 0, 0, 0}

	err := gridDiskDistances(polar, 3, k2, k2Dist)
	if err != E_SUCCESS {
		t.Fatalf("gridDiskDistances failed: %v", err)
	}

	k2present := 0
	for i := 0; i < 37; i++ {
		if k2[i] != 0 {
			k2present++
			inList := 0
			for j := 0; j < 37; j++ {
				if k2[i] == expectedK2[j] {
					if k2Dist[i] != expectedK2Dist[j] {
						t.Errorf("distance is not as expected at index %d: got %d, expected %d", i, k2Dist[i], expectedK2Dist[j])
					}
					inList++
				}
			}
			if inList != 1 {
				t.Errorf("index %x found %d times in expected set, expected 1", k2[i], inList)
			}
		}
	}
	if k2present != 31 {
		t.Errorf("pentagon has %d neighbors, expected 31 (30 neighbors + itself)", k2present)
	}
}

func TestGridDisk1_Pentagon_k4(t *testing.T) {
	t.Parallel()

	var pent H3Index
	setH3Index(&pent, 1, 14, 0)

	var maxSize int64
	err := maxGridDiskSize(4, &maxSize)
	if err != E_SUCCESS {
		t.Fatalf("maxGridDiskSize failed: %v", err)
	}

	k2 := make([]H3Index, maxSize)
	k2Dist := make([]int32, maxSize)
	expectedK2 := []H3Index{0x811d7ffffffffff,
		0x810c7ffffffffff,
		0x81227ffffffffff,
		0x81293ffffffffff,
		0x81133ffffffffff,
		0x8136bffffffffff,
		0x81167ffffffffff,
		0x811d3ffffffffff,
		0x810c3ffffffffff,
		0x81223ffffffffff,
		0x81477ffffffffff,
		0x8128fffffffffff,
		0x81367ffffffffff,
		0x8112fffffffffff,
		0x811cfffffffffff,
		0x8123bffffffffff,
		0x810dbffffffffff,
		0x8112bffffffffff,
		0x81473ffffffffff,
		0x8128bffffffffff,
		0x81363ffffffffff,
		0x811cbffffffffff,
		0x81237ffffffffff,
		0x810d7ffffffffff,
		0x81127ffffffffff,
		0x8137bffffffffff,
		0x81287ffffffffff,
		0x8126bffffffffff,
		0x81177ffffffffff,
		0x810d3ffffffffff,
		0x81233ffffffffff,
		0x8150fffffffffff,
		0x81123ffffffffff,
		0x81377ffffffffff,
		0x81283ffffffffff,
		0x8102fffffffffff,
		0x811c3ffffffffff,
		0x810cfffffffffff,
		0x8122fffffffffff,
		0x8113bffffffffff,
		0x81373ffffffffff,
		0x8129bffffffffff,
		0x8102bffffffffff,
		0x811dbffffffffff,
		0x810cbffffffffff,
		0x8122bffffffffff,
		0x81297ffffffffff,
		0x81507ffffffffff,
		0x8136fffffffffff,
		0x8127bffffffffff,
		0x81137ffffffffff,
		0,
		0}

	err = gridDiskDistances(pent, 4, k2, k2Dist)
	if err != E_SUCCESS {
		t.Fatalf("gridDiskDistances failed: %v", err)
	}

	k2present := 0
	for i := 0; i < len(k2); i++ {
		if k2[i] != 0 {
			k2present++
			inList := 0
			for j := 0; j < len(expectedK2); j++ {
				if k2[i] == expectedK2[j] {
					inList++
				}
			}
			if inList != 1 {
				t.Errorf("index %x found %d times in expected set, expected 1", k2[i], inList)
			}
		}
	}
	if k2present != 51 {
		t.Errorf("pentagon has %d neighbors, expected 51 (50 neighbors + itself)", k2present)
	}
}

func TestGridDisk_equals_gridDiskDistancesSafe(t *testing.T) {
	// Check that gridDiskDistances output matches gridDiskDistancesSafe,
	// since gridDiskDistances will sometimes use a different
	// implementation.

	for res := int32(0); res < 2; res++ {
		t.Run("", func(t *testing.T) {
			_iterateAllIndexesAtRes(res, func(h3 H3Index) {
				gridDisk_equals_gridDiskDistancesSafe_assertions(t, h3)
			})
		})
	}
}

func TestGridDiskInvalid(t *testing.T) {
	t.Parallel()

	k := int32(1000)
	var kSz int64
	err := maxGridDiskSize(k, &kSz)
	if err != E_SUCCESS {
		t.Fatalf("maxGridDiskSize failed: %v", err)
	}
	neighbors := make([]H3Index, kSz)
	err = gridDisk(0x7fffffffffffffff, k, neighbors)
	if err != E_CELL_INVALID {
		t.Errorf("gridDisk should return E_CELL_INVALID for invalid input, got %v", err)
	}
}

func TestGridDiskInvalidDigit(t *testing.T) {
	t.Parallel()

	k := int32(2)
	var kSz int64
	err := maxGridDiskSize(k, &kSz)
	if err != E_SUCCESS {
		t.Fatalf("maxGridDiskSize failed: %v", err)
	}
	neighbors := make([]H3Index, kSz)
	err = gridDisk(0x4d4b00fe5c5c3030, k, neighbors)
	if err != E_CELL_INVALID {
		t.Errorf("gridDisk should return E_CELL_INVALID for invalid input, got %v", err)
	}
}

func TestGridDiskDistances_invalidK(t *testing.T) {
	t.Parallel()

	index := H3Index(0x811d7ffffffffff)
	err := gridDiskDistances(index, -1, nil, nil)
	if err != E_DOMAIN {
		t.Errorf("gridDiskDistances should return E_DOMAIN for invalid k, got %v", err)
	}

	err = gridDiskDistancesUnsafe(index, -1, nil, nil)
	if err != E_DOMAIN {
		t.Errorf("gridDiskDistancesUnsafe should return E_DOMAIN for invalid k, got %v", err)
	}

	err = gridDiskDistancesSafe(index, -1, nil, nil)
	if err != E_DOMAIN {
		t.Errorf("gridDiskDistancesSafe should return E_DOMAIN for invalid k, got %v", err)
	}
}

func TestMaxGridDiskSize_invalid(t *testing.T) {
	t.Parallel()

	var sz int64
	err := maxGridDiskSize(-1, &sz)
	if err != E_DOMAIN {
		t.Errorf("maxGridDiskSize should return E_DOMAIN for negative k, got %v", err)
	}
}

func TestMaxGridDiskSize_large(t *testing.T) {
	t.Parallel()

	var sz int64
	err := maxGridDiskSize(26755, &sz)
	if err != E_SUCCESS {
		t.Fatalf("maxGridDiskSize failed: %v", err)
	}
	if sz != 2147570341 {
		t.Errorf("large (> 32 bit signed int) k should work: expected 2147570341, got %d", sz)
	}
}

func TestMaxGridDiskSize_numCells(t *testing.T) {
	t.Parallel()

	var sz int64
	prev := int64(0)
	var maxCells int64
	maxCells, err := getNumCells(15)
	if err != E_SUCCESS {
		t.Fatalf("getNumCells failed: %v", err)
	}

	// 13780510 will produce values above max
	for k := int32(13780510 - 100); k < 13780510+100; k++ {
		err = maxGridDiskSize(k, &sz)
		if err != E_SUCCESS {
			t.Fatalf("maxGridDiskSize failed for k=%d: %v", k, err)
		}
		if sz > maxCells {
			t.Errorf("maxGridDiskSize does not produce estimates above the number of grid cells: k=%d, sz=%d, max=%d", k, sz, maxCells)
		}
		if prev > sz {
			t.Errorf("maxGridDiskSize is not monotonically increasing: prev=%d, current=%d at k=%d", prev, sz, k)
		}
		prev = sz
	}

	err = maxGridDiskSize(INT32_MAX, &sz)
	if err != E_SUCCESS {
		t.Fatalf("maxGridDiskSize failed for INT32_MAX: %v", err)
	}
	if sz != maxCells {
		t.Errorf("maxGridDiskSize of INT32_MAX should produce valid result: expected %d, got %d", maxCells, sz)
	}
}
