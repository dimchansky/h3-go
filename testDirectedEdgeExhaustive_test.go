// Tests ported from testDirectedEdgeExhaustive.c
package h3

import (
	"testing"
)

// Helper function to iterate all indexes at a given resolution.
func iterateAllIndexesAtResForDirectedEdge(t *testing.T, res int32, testFunc func(t *testing.T, h3 h3Index)) {
	t.Helper()

	// Get all base cells
	baseCells := make([]h3Index, numBaseCells)
	if err := getRes0Cells(baseCells); err != eSuccess {
		t.Fatalf("Failed to get res 0 cells: %v", err)
	}

	if res == 0 {
		// For resolution 0, just test the base cells
		for _, cell := range baseCells {
			testFunc(t, cell)
		}
		return
	}

	// For higher resolutions, get children of each base cell
	for _, baseCell := range baseCells {
		childrenSize, err := cellToChildrenSize(baseCell, res)
		if err != eSuccess {
			continue // Some cells might not have children at certain resolutions
		}

		children := make([]h3Index, childrenSize)
		if err := cellToChildren(baseCell, res, children); err != eSuccess {
			continue
		}

		for _, child := range children {
			if child != h3Null {
				testFunc(t, child)
			}
		}
	}
}

// Helper function to iterate base cell indexes at a specific resolution for directed edge tests.
func iterateBaseCellIndexesAtResForDirectedEdge(t *testing.T, res int32, testFunc func(t *testing.T, h3 h3Index), baseCell int32) {
	t.Helper()

	// Create base cell index
	baseCellIndex := h3Index(h3Init)
	baseCellIndex = setMode(baseCellIndex, h3CellMode)
	baseCellIndex = setBaseCell(baseCellIndex, baseCell)

	if res == 0 {
		testFunc(t, baseCellIndex)
		return
	}

	// Get children at the specified resolution
	childrenSize, err := cellToChildrenSize(baseCellIndex, res)
	if err != eSuccess {
		t.Fatalf("Failed to get children size for base cell %d at res %d: %v", baseCell, res, err)
	}

	children := make([]h3Index, childrenSize)
	if err := cellToChildren(baseCellIndex, res, children); err != eSuccess {
		t.Fatalf("Failed to get children for base cell %d at res %d: %v", baseCell, res, err)
	}

	for _, child := range children {
		if child != h3Null {
			testFunc(t, child)
		}
	}
}

func directedEdge_correctness_assertions(t *testing.T, h3 h3Index) {
	t.Helper()

	edges := make([]h3Index, 6)
	pentagon := isPentagon(h3)
	if err := originToDirectedEdges(h3, edges); err != eSuccess {
		t.Errorf("originToDirectedEdges failed for cell %#016x: %v", h3, err)
		return
	}

	for i := 0; i < 6; i++ {
		if pentagon && i == 0 {
			if edges[i] != h3Null {
				t.Errorf("Expected null edge for pentagon at position 0, got %#016x", edges[i])
			}
			continue
		}

		if !isValidDirectedEdge(edges[i]) {
			t.Errorf("Edge %#016x is not valid directed edge", edges[i])
			continue
		}

		origin, err := getDirectedEdgeOrigin(edges[i])
		if err != eSuccess {
			t.Errorf("getDirectedEdgeOrigin failed for edge %#016x: %v", edges[i], err)
			continue
		}
		if origin != h3 {
			t.Errorf("Origin mismatch for edge %#016x: got %#016x, expected %#016x", edges[i], origin, h3)
		}

		var destination h3Index
		if err := getDirectedEdgeDestination(edges[i], &destination); err != eSuccess {
			t.Errorf("getDirectedEdgeDestination failed for edge %#016x: %v", edges[i], err)
			continue
		}

		isNeighbor, err := areNeighborCells(h3, destination)
		if err != eSuccess {
			t.Errorf("areNeighborCells failed for cells %#016x and %#016x: %v", h3, destination, err)
			continue
		}
		if !isNeighbor {
			t.Errorf("Destination %#016x is not a neighbor of origin %#016x", destination, h3)
		}
	}
}

func directedEdge_boundary_assertions(t *testing.T, h3 h3Index) {
	t.Helper()

	edges := make([]h3Index, 6)
	if err := originToDirectedEdges(h3, edges); err != eSuccess {
		t.Errorf("originToDirectedEdges failed for cell %#016x: %v", h3, err)
		return
	}

	var destination h3Index
	var revEdge h3Index
	var edgeBoundary CellBoundary
	var revEdgeBoundary CellBoundary

	for i := 0; i < 6; i++ {
		if edges[i] == h3Null {
			continue
		}

		if err := directedEdgeToBoundary(edges[i], &edgeBoundary); err != eSuccess {
			t.Errorf("directedEdgeToBoundary failed for edge %#016x: %v", edges[i], err)
			continue
		}

		if err := getDirectedEdgeDestination(edges[i], &destination); err != eSuccess {
			t.Errorf("getDirectedEdgeDestination failed for edge %#016x: %v", edges[i], err)
			continue
		}

		var err h3Error
		revEdge, err = cellsToDirectedEdge(destination, h3)
		if err != eSuccess {
			t.Errorf("cellsToDirectedEdge failed for cells %#016x -> %#016x: %v", destination, h3, err)
			continue
		}

		if err := directedEdgeToBoundary(revEdge, &revEdgeBoundary); err != eSuccess {
			t.Errorf("directedEdgeToBoundary failed for reverse edge %#016x: %v", revEdge, err)
			continue
		}

		if edgeBoundary.NumVerts != revEdgeBoundary.NumVerts {
			t.Errorf("NumVerts mismatch for edge %#016x: edge=%d, reverse=%d",
				edges[i], edgeBoundary.NumVerts, revEdgeBoundary.NumVerts)
			continue
		}

		for j := int32(0); j < edgeBoundary.NumVerts; j++ {
			almostEqual := geoAlmostEqualThreshold(
				&edgeBoundary.Verts[j],
				&revEdgeBoundary.Verts[revEdgeBoundary.NumVerts-1-j],
				0.000001)
			if !almostEqual {
				t.Errorf("Vertex mismatch for edge %#016x at position %d", edges[i], j)
			}
		}
	}
}

func TestDirectedEdge_correctness(t *testing.T) {
	t.Parallel()
	iterateAllIndexesAtResForDirectedEdge(t, 0, directedEdge_correctness_assertions)
	iterateAllIndexesAtResForDirectedEdge(t, 1, directedEdge_correctness_assertions)
	iterateAllIndexesAtResForDirectedEdge(t, 2, directedEdge_correctness_assertions)
	iterateAllIndexesAtResForDirectedEdge(t, 3, directedEdge_correctness_assertions)
	iterateAllIndexesAtResForDirectedEdge(t, 4, directedEdge_correctness_assertions)
}

func TestDirectedEdge_boundary(t *testing.T) {
	t.Parallel()
	iterateAllIndexesAtResForDirectedEdge(t, 0, directedEdge_boundary_assertions)
	iterateAllIndexesAtResForDirectedEdge(t, 1, directedEdge_boundary_assertions)
	iterateAllIndexesAtResForDirectedEdge(t, 2, directedEdge_boundary_assertions)
	iterateAllIndexesAtResForDirectedEdge(t, 3, directedEdge_boundary_assertions)
	iterateAllIndexesAtResForDirectedEdge(t, 4, directedEdge_boundary_assertions)
	// Res 5: normal base cell
	iterateBaseCellIndexesAtResForDirectedEdge(t, 5, directedEdge_boundary_assertions, 0)
	// Res 5: pentagon base cell
	iterateBaseCellIndexesAtResForDirectedEdge(t, 5, directedEdge_boundary_assertions, 14)
	// Res 5: polar pentagon base cell
	iterateBaseCellIndexesAtResForDirectedEdge(t, 5, directedEdge_boundary_assertions, 117)
	// Res 6: Test one pentagon just to check for new edge cases
	iterateBaseCellIndexesAtResForDirectedEdge(t, 6, directedEdge_boundary_assertions, 14)
}
