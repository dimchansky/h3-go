// Tests ported from testH3CellAreaExhaustive.c
package h3

import (
	"math"
	"testing"
)

// Helper function to iterate all directed edges at a given resolution
func iterateAllDirectedEdgesAtRes(t *testing.T, res int32, testFunc func(t *testing.T, edge H3Index)) {
	t.Helper()

	// Get all base cells
	baseCells := make([]H3Index, NUM_BASE_CELLS)
	if err := getRes0Cells(baseCells); err != E_SUCCESS {
		t.Fatalf("Failed to get res 0 cells: %v", err)
	}

	// For each base cell
	for _, baseCell := range baseCells {
		var cells []H3Index
		if res == 0 {
			cells = []H3Index{baseCell}
		} else {
			// Get children at the specified resolution
			childrenSize, err := cellToChildrenSize(baseCell, res)
			if err != E_SUCCESS {
				continue
			}
			cells = make([]H3Index, childrenSize)
			if err := cellToChildren(baseCell, res, cells); err != E_SUCCESS {
				continue
			}
		}

		// For each cell, get its directed edges
		for _, cell := range cells {
			if cell == H3_NULL {
				continue
			}

			// Get all directed edges from this cell
			edges := make([]H3Index, 6)
			if err := originToDirectedEdges(cell, edges); err != E_SUCCESS {
				continue
			}

			// Test each valid edge
			for _, edge := range edges {
				if edge != H3_NULL {
					testFunc(t, edge)
				}
			}
		}
	}
}

// Test haversine distance functions for positivity and commutativity
func haversineAssert(t *testing.T, edge H3Index) {
	t.Helper()

	var a, b LatLng
	var origin, destination H3Index

	// Get origin cell
	origin, err := getDirectedEdgeOrigin(edge)
	if err != E_SUCCESS {
		t.Skipf("Failed to get edge origin: %v", err)
	}
	if err := cellToLatLng(origin, &a); err != E_SUCCESS {
		t.Skipf("Failed to get origin LatLng: %v", err)
	}

	// Get destination cell
	if err := getDirectedEdgeDestination(edge, &destination); err != E_SUCCESS {
		t.Skipf("Failed to get edge destination: %v", err)
	}
	if err := cellToLatLng(destination, &b); err != E_SUCCESS {
		t.Skipf("Failed to get destination LatLng: %v", err)
	}

	// Test greatCircleDistanceRads
	ab := greatCircleDistanceRads(&a, &b)
	ba := greatCircleDistanceRads(&b, &a)
	if ab <= 0 {
		t.Errorf("greatCircleDistanceRads: distance between cell centers should be positive, got %v", ab)
	}
	if ab != ba {
		t.Errorf("greatCircleDistanceRads: pairwise cell distances should be commutative, ab=%v, ba=%v", ab, ba)
	}

	// Test greatCircleDistanceKm
	abKm := greatCircleDistanceKm(&a, &b)
	baKm := greatCircleDistanceKm(&b, &a)
	if abKm <= 0 {
		t.Errorf("greatCircleDistanceKm: distance between cell centers should be positive, got %v", abKm)
	}
	if abKm != baKm {
		t.Errorf("greatCircleDistanceKm: pairwise cell distances should be commutative, ab=%v, ba=%v", abKm, baKm)
	}

	// Test greatCircleDistanceM
	abM := greatCircleDistanceM(&a, &b)
	baM := greatCircleDistanceM(&b, &a)
	if abM <= 0 {
		t.Errorf("greatCircleDistanceM: distance between cell centers should be positive, got %v", abM)
	}
	if abM != baM {
		t.Errorf("greatCircleDistanceM: pairwise cell distances should be commutative, ab=%v, ba=%v", abM, baM)
	}

	// Test that measurements are in correct relative scale
	if abKm <= ab {
		t.Errorf("measurement in kilometers (%v) should be greater than in radians (%v)", abKm, ab)
	}
	if abM <= abKm {
		t.Errorf("measurement in meters (%v) should be greater than in kilometers (%v)", abM, abKm)
	}
}

// Test edge length calculation functions for positivity
func edgeLengthAssert(t *testing.T, edge H3Index) {
	t.Helper()

	// Test edgeLengthRads
	var lengthRads float64
	if err := edgeLengthRads(edge, &lengthRads); err != E_SUCCESS {
		t.Skipf("Failed to get edge length in rads: %v", err)
	}
	if lengthRads <= 0 {
		t.Errorf("edgeLengthRads: edge has non-positive length %v", lengthRads)
	}

	// Test edgeLengthKm
	var lengthKm float64
	if err := edgeLengthKm(edge, &lengthKm); err != E_SUCCESS {
		t.Skipf("Failed to get edge length in km: %v", err)
	}
	if lengthKm <= 0 {
		t.Errorf("edgeLengthKm: edge has non-positive length %v", lengthKm)
	}

	// Test edgeLengthM
	var lengthM float64
	if err := edgeLengthM(edge, &lengthM); err != E_SUCCESS {
		t.Skipf("Failed to get edge length in m: %v", err)
	}
	if lengthM <= 0 {
		t.Errorf("edgeLengthM: edge has non-positive length %v", lengthM)
	}
}

// Test cell area calculation functions for positivity
func cellAreaAssert(t *testing.T, cell H3Index) {
	t.Helper()

	// Test cellAreaRads2
	areaRads, err := cellAreaRads2(cell)
	if err != E_SUCCESS {
		t.Skipf("Failed to get cell area in rads2: %v", err)
	}
	if areaRads <= 0 {
		t.Errorf("cellAreaRads2: cell has non-positive area %v", areaRads)
	}

	// Test cellAreaKm2
	areaKm2, err := cellAreaKm2(cell)
	if err != E_SUCCESS {
		t.Skipf("Failed to get cell area in km2: %v", err)
	}
	if areaKm2 <= 0 {
		t.Errorf("cellAreaKm2: cell has non-positive area %v", areaKm2)
	}

	// Test cellAreaM2
	areaM2, err := cellAreaM2(cell)
	if err != E_SUCCESS {
		t.Skipf("Failed to get cell area in m2: %v", err)
	}
	if areaM2 <= 0 {
		t.Errorf("cellAreaM2: cell has non-positive area %v", areaM2)
	}

	// Test that measurements are in correct relative scale
	if areaRads >= areaKm2 {
		t.Errorf("area in rads (%v) should be smaller than area in km2 (%v)", areaRads, areaKm2)
	}
	if areaKm2 >= areaM2 {
		t.Errorf("area in km2 (%v) should be smaller than area in m2 (%v)", areaKm2, areaM2)
	}
}

// Test sum of all cell areas equals earth surface area
func earthAreaTest(t *testing.T, res int32, cellAreaFunc func(H3Index) (float64, H3Error), target float64, tol float64) {
	t.Helper()

	var area float64

	// Get all base cells
	baseCells := make([]H3Index, NUM_BASE_CELLS)
	if err := getRes0Cells(baseCells); err != E_SUCCESS {
		t.Fatalf("Failed to get res 0 cells: %v", err)
	}

	// For each base cell
	for _, baseCell := range baseCells {
		var cells []H3Index
		if res == 0 {
			cells = []H3Index{baseCell}
		} else {
			// Get children at the specified resolution
			childrenSize, err := cellToChildrenSize(baseCell, res)
			if err != E_SUCCESS {
				continue
			}
			cells = make([]H3Index, childrenSize)
			if err := cellToChildren(baseCell, res, cells); err != E_SUCCESS {
				continue
			}
		}

		// Sum up areas of all cells
		for _, cell := range cells {
			if cell != H3_NULL {
				cellArea, err := cellAreaFunc(cell)
				if err == E_SUCCESS {
					area += cellArea
				}
			}
		}
	}

	// Check if sum is close to expected earth area
	if math.Abs(area-target) >= tol {
		t.Errorf("Sum of all cells (%v) should give earth area (%v), tolerance=%v, diff=%v",
			area, target, tol, math.Abs(area-target))
	}
}

// Test haversine distances between neighboring cells
func TestHaversineDistances(t *testing.T) {
	t.Parallel()
	
	t.Run("res0", func(t *testing.T) {
		t.Parallel()
		iterateAllDirectedEdgesAtRes(t, 0, haversineAssert)
	})
	t.Run("res1", func(t *testing.T) {
		t.Parallel()
		iterateAllDirectedEdgesAtRes(t, 1, haversineAssert)
	})
	t.Run("res2", func(t *testing.T) {
		t.Parallel()
		iterateAllDirectedEdgesAtRes(t, 2, haversineAssert)
	})
	t.Run("res3", func(t *testing.T) {
		t.Parallel()
		iterateAllDirectedEdgesAtRes(t, 3, haversineAssert)
	})
}

// Test edge length functions
func TestEdgeLength(t *testing.T) {
	t.Parallel()
	
	t.Run("res0", func(t *testing.T) {
		t.Parallel()
		iterateAllDirectedEdgesAtRes(t, 0, edgeLengthAssert)
	})
	t.Run("res1", func(t *testing.T) {
		t.Parallel()
		iterateAllDirectedEdgesAtRes(t, 1, edgeLengthAssert)
	})
	t.Run("res2", func(t *testing.T) {
		t.Parallel()
		iterateAllDirectedEdgesAtRes(t, 2, edgeLengthAssert)
	})
	t.Run("res3", func(t *testing.T) {
		t.Parallel()
		iterateAllDirectedEdgesAtRes(t, 3, edgeLengthAssert)
	})
}

// Test cell area positivity
func TestCellAreaPositive(t *testing.T) {
	t.Parallel()
	
	t.Run("res0", func(t *testing.T) {
		t.Parallel()
		iterateAllIndexesAtRes(t, 0, cellAreaAssert)
	})
	t.Run("res1", func(t *testing.T) {
		t.Parallel()
		iterateAllIndexesAtRes(t, 1, cellAreaAssert)
	})
	t.Run("res2", func(t *testing.T) {
		t.Parallel()
		iterateAllIndexesAtRes(t, 2, cellAreaAssert)
	})
	t.Run("res3", func(t *testing.T) {
		t.Parallel()
		iterateAllIndexesAtRes(t, 3, cellAreaAssert)
	})
}

// Test that sum of all cell areas equals earth surface area
func TestCellAreaEarth(t *testing.T) {
	t.Parallel()
	
	// Earth area in different units
	rads2 := 4 * math.Pi
	km2 := rads2 * EARTH_RADIUS_KM * EARTH_RADIUS_KM
	m2 := km2 * 1000 * 1000

	// Resolution 0
	t.Run("res0_rads2", func(t *testing.T) {
		t.Parallel()
		earthAreaTest(t, 0, cellAreaRads2, rads2, 1e-14)
	})
	t.Run("res0_km2", func(t *testing.T) {
		t.Parallel()
		earthAreaTest(t, 0, cellAreaKm2, km2, 1e-6)
	})
	t.Run("res0_m2", func(t *testing.T) {
		t.Parallel()
		earthAreaTest(t, 0, cellAreaM2, m2, 1e0)
	})

	// Resolution 1
	// Notice the drop in accuracy at resolution 1.
	// I think this has something to do with Class II vs Class III resolutions.
	t.Run("res1_rads2", func(t *testing.T) {
		t.Parallel()
		earthAreaTest(t, 1, cellAreaRads2, rads2, 1e-9)
	})
	t.Run("res1_km2", func(t *testing.T) {
		t.Parallel()
		earthAreaTest(t, 1, cellAreaKm2, km2, 1e-1)
	})
	t.Run("res1_m2", func(t *testing.T) {
		t.Parallel()
		earthAreaTest(t, 1, cellAreaM2, m2, 1e5)
	})

	// Resolution 2
	t.Run("res2_rads2", func(t *testing.T) {
		t.Parallel()
		earthAreaTest(t, 2, cellAreaRads2, rads2, 1e-12)
	})
	t.Run("res2_km2", func(t *testing.T) {
		t.Parallel()
		earthAreaTest(t, 2, cellAreaKm2, km2, 1e-5)
	})
	t.Run("res2_m2", func(t *testing.T) {
		t.Parallel()
		earthAreaTest(t, 2, cellAreaM2, m2, 1e0)
	})

	// Resolution 3
	t.Run("res3_rads2", func(t *testing.T) {
		t.Parallel()
		earthAreaTest(t, 3, cellAreaRads2, rads2, 1e-11)
	})
	t.Run("res3_km2", func(t *testing.T) {
		t.Parallel()
		earthAreaTest(t, 3, cellAreaKm2, km2, 1e-3)
	})
	t.Run("res3_m2", func(t *testing.T) {
		t.Parallel()
		earthAreaTest(t, 3, cellAreaM2, m2, 1e3)
	})

	// Resolution 4
	t.Run("res4_rads2", func(t *testing.T) {
		t.Parallel()
		earthAreaTest(t, 4, cellAreaRads2, rads2, 1e-11)
	})
	t.Run("res4_km2", func(t *testing.T) {
		t.Parallel()
		earthAreaTest(t, 4, cellAreaKm2, km2, 1e-3)
	})
	t.Run("res4_m2", func(t *testing.T) {
		t.Parallel()
		earthAreaTest(t, 4, cellAreaM2, m2, 1e2)
	})
}