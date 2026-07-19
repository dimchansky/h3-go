// Tests ported from H3 v4.5.0: src/apps/testapps/testCellsToMultiPolyInternal.c.
package h3

import "testing"

func TestCellsToMultiPolyInternal_destroyArcSet_with_arcs(t *testing.T) {
	t.Parallel()

	// Test that destroyArcSet frees memory and sets pointers to NULL
	var arcset arcSet
	arcset.numArcs = 10
	arcset.numBuckets = 100
	arcset.arcs = make([]arc, arcset.numArcs)
	arcset.buckets = make([]*arc, arcset.numBuckets)

	if arcset.arcs == nil {
		t.Error("arcs should be allocated")
	}
	if arcset.buckets == nil {
		t.Error("buckets should be allocated")
	}

	destroyArcSet(&arcset)

	if arcset.arcs != nil {
		t.Error("arcs should be NULL after destroy")
	}
	if arcset.buckets != nil {
		t.Error("buckets should be NULL after destroy")
	}

	// Call again on NULL pointers (should be safe)
	destroyArcSet(&arcset)
}

func TestCellsToMultiPolyInternal_destroySortableLoopSet_with_verts(t *testing.T) {
	t.Parallel()

	// Test with allocated loops and verts
	var loopset sortableLoopSet
	loopset.numLoops = 3
	loopset.sloops = make([]sortableLoop, 3)

	// First loop has verts (exercises positive branch)
	loopset.sloops[0].loop = make(GeoLoop, 5)

	// Second loop has NULL verts (exercises negative branch)
	loopset.sloops[1].loop = nil

	// Third loop has verts (exercises positive branch)
	loopset.sloops[2].loop = make(GeoLoop, 4)

	destroySortableLoopSet(&loopset)

	if loopset.sloops != nil {
		t.Error("sloops should be NULL after destroy")
	}
}

func TestCellsToMultiPolyInternal_destroySortableLoopSet_null(t *testing.T) {
	t.Parallel()

	// Test with NULL sloops (exercises negative branch of outer if)
	var loopset sortableLoopSet
	loopset.numLoops = 0
	loopset.sloops = nil

	destroySortableLoopSet(&loopset)

	if loopset.sloops != nil {
		t.Error("sloops should remain NULL")
	}
}

func TestCellsToMultiPolyInternal_destroySortablePolys_with_holes(t *testing.T) {
	t.Parallel()

	// Test with allocated polygons and holes
	spolys := make([]sortablePoly, 2)

	// First polygon has holes (exercises positive branch)
	spolys[0].poly.Holes = make([]GeoLoop, 2)

	// Second polygon has NULL holes (exercises negative branch)
	spolys[1].poly.Holes = nil

	destroySortablePolys(spolys, 2)
	// spolys is freed in C; can't assert on it
}

func TestCellsToMultiPolyInternal_destroySortablePolys_null(t *testing.T) {
	t.Parallel()

	// Test with NULL spolys (exercises negative branch of outer if)
	destroySortablePolys(nil, 0)
	// Should not crash
}

func TestCellsToMultiPolyInternal_destroySortablePolyVerts_with_verts(t *testing.T) {
	t.Parallel()

	// Test with allocated polygons and outer loop verts
	spolys := make([]sortablePoly, 2)

	// First polygon has verts (exercises positive branch)
	spolys[0].poly.GeoLoop = make(GeoLoop, 6)

	// Second polygon has NULL verts (exercises negative branch)
	spolys[1].poly.GeoLoop = nil

	destroySortablePolyVerts(spolys, 2)
	// spolys is freed in C; can't assert on it
}

func TestCellsToMultiPolyInternal_destroySortablePolyVerts_null(t *testing.T) {
	t.Parallel()

	// Test with NULL spolys (exercises negative branch of outer if)
	destroySortablePolyVerts(nil, 0)
	// Should not crash
}

func TestCellsToMultiPolyInternal_cmp_SortablePoly_equal(t *testing.T) {
	t.Parallel()

	// Test equality branch of cmp_SortablePoly
	var a, b sortablePoly
	a.outerArea = 100.0
	b.outerArea = 100.0

	if result := cmp_SortablePoly(&a, &b); result != 0 {
		t.Errorf("Equal areas should return 0, got %d", result)
	}
}

func TestCellsToMultiPolyInternal_cmp_SortablePoly_descending(t *testing.T) {
	t.Parallel()

	// Test descending order (larger area comes first)
	var a, b sortablePoly

	// a has larger area, should come first (return -1)
	a.outerArea = 200.0
	b.outerArea = 100.0
	if result := cmp_SortablePoly(&a, &b); result != -1 {
		t.Errorf("Larger area should come first, got %d", result)
	}

	// b has larger area, should come first (return 1)
	a.outerArea = 100.0
	b.outerArea = 200.0
	if result := cmp_SortablePoly(&a, &b); result != 1 {
		t.Errorf("Smaller area should come after, got %d", result)
	}
}

func TestCellsToMultiPolyInternal_checkCellsToMultiPolyOverflow_safe(t *testing.T) {
	t.Parallel()

	const hashMultiplier = hashTableMultiplier

	// Test with reasonable number of cells (should succeed)
	if err := checkCellsToMultiPolyOverflow(1000000, hashMultiplier); err != eSuccess {
		t.Errorf("Should succeed for reasonable numCells, got %v", err)
	}

	// Test with zero cells (should succeed)
	if err := checkCellsToMultiPolyOverflow(0, hashMultiplier); err != eSuccess {
		t.Errorf("Should succeed for zero cells, got %v", err)
	}

	// Test with negative cells (should succeed - validated elsewhere)
	if err := checkCellsToMultiPolyOverflow(-1, hashMultiplier); err != eSuccess {
		t.Errorf("Should succeed for negative (check doesn't apply), got %v", err)
	}

	// Test with small and large hash multipliers.
	// Largest allocated array will change, depending on multiplier.
	if err := checkCellsToMultiPolyOverflow(1000000, 1); err != eSuccess {
		t.Errorf("multiplier 1: got %v", err)
	}
	if err := checkCellsToMultiPolyOverflow(1000000, 100); err != eSuccess {
		t.Errorf("multiplier 100: got %v", err)
	}
}

func TestCellsToMultiPolyInternal_checkCellsToMultiPolyOverflow_wouldOverflow(t *testing.T) {
	t.Parallel()

	const hashMultiplier = hashTableMultiplier

	// Test with numCells that would cause size_t overflow (using the
	// pinned C oracle sizes; see the constants next to the port).
	maxBytesPerCell := uint64(6 * hashTableMultiplier * cSizeofArcPtr)
	maxSafeNumCells := cSizeMax / maxBytesPerCell
	overflowNumCells := int64(maxSafeNumCells + 1)

	if err := checkCellsToMultiPolyOverflow(overflowNumCells, hashMultiplier); err != eMemoryBounds {
		t.Errorf("Should return eMemoryBounds when overflow would occur, got %v", err)
	}

	// Also test INT64_MAX directly
	const int64Max = int64(^uint64(0) >> 1)
	if err := checkCellsToMultiPolyOverflow(int64Max, hashMultiplier); err != eMemoryBounds {
		t.Errorf("Should return eMemoryBounds for INT64_MAX, got %v", err)
	}
}
