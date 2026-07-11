// Tests ported from testVertexExhaustive.c
package h3

import (
	"testing"
)

// Helper function to iterate all indexes at a given resolution.
func iterateAllIndexesAtRes(t *testing.T, res int32, testFunc func(t *testing.T, h3 H3Index)) {
	t.Helper()

	// Get all base cells
	baseCells := make([]H3Index, NUM_BASE_CELLS)
	if err := getRes0Cells(baseCells); err != E_SUCCESS {
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
		if err != E_SUCCESS {
			continue // Some cells might not have children at certain resolutions
		}

		children := make([]H3Index, childrenSize)
		if err := cellToChildren(baseCell, res, children); err != E_SUCCESS {
			continue
		}

		for _, child := range children {
			if child != H3_NULL {
				testFunc(t, child)
			}
		}
	}
}

// Helper function to iterate base cell indexes at a specific resolution.
func iterateBaseCellIndexesAtRes(t *testing.T, res int32, testFunc func(t *testing.T, h3 H3Index), baseCell int32) {
	t.Helper()

	// Create base cell index
	baseCellIndex := H3Index(H3_INIT)
	baseCellIndex = setMode(baseCellIndex, H3_CELL_MODE)
	baseCellIndex = setBaseCell(baseCellIndex, baseCell)

	if res == 0 {
		testFunc(t, baseCellIndex)
		return
	}

	// Get children at the specified resolution
	childrenSize, err := cellToChildrenSize(baseCellIndex, res)
	if err != E_SUCCESS {
		t.Fatalf("Failed to get children size for base cell %d at res %d: %v", baseCell, res, err)
	}

	children := make([]H3Index, childrenSize)
	if err := cellToChildren(baseCellIndex, res, children); err != E_SUCCESS {
		t.Fatalf("Failed to get children for base cell %d at res %d: %v", baseCell, res, err)
	}

	for _, child := range children {
		if child != H3_NULL {
			testFunc(t, child)
		}
	}
}

func directionForVertexNum_symmetry_assertions(t *testing.T, h3 H3Index) {
	t.Helper()

	numVerts := int32(NUM_HEX_VERTS)
	isPent := isPentagon(h3)
	if isPent {
		numVerts = NUM_PENT_VERTS
	}

	for i := int32(0); i < numVerts; i++ {
		dir := directionForVertexNum(h3, i)
		vertexNum := vertexNumForDirection(h3, dir)
		if vertexNum != i {
			t.Errorf("Direction symmetry failed for cell %#016x (res=%d, %s):\n"+
				"  vertexNum=%d -> direction=%d -> vertexNum=%d (expected %d)",
				h3, getResolution(h3),
				map[bool]string{true: "pentagon", false: "hexagon"}[isPent],
				i, dir, vertexNum, i)
		}
	}
}

func cellToVertex_point_assertions(t *testing.T, h3 H3Index) {
	t.Helper()

	var gb CellBoundary
	if err := cellToBoundary(h3, &gb); err != E_SUCCESS {
		t.Skipf("Failed to get cell boundary for h3=%#016x: %v", h3, err)
		return
	}

	numVerts := int32(NUM_HEX_VERTS)
	isPent := isPentagon(h3)
	if isPent {
		numVerts = NUM_PENT_VERTS
	}

	// This test won't work if there are distortion vertexes in the boundary
	if numVerts < int32(gb.NumVerts) {
		return
	}

	var coord LatLng
	for i := int32(0); i < numVerts; i++ {
		var vertex H3Index
		if err := cellToVertex(h3, i, &vertex); err != E_SUCCESS {
			t.Errorf("cellToVertex failed for cell %#016x (res=%d, %s), vertexNum=%d: error=%v",
				h3, getResolution(h3),
				map[bool]string{true: "pentagon", false: "hexagon"}[isPent],
				i, err)
			continue
		}

		if err := vertexToLatLng(vertex, &coord); err != E_SUCCESS {
			t.Errorf("vertexToLatLng failed for vertex %#016x (from cell %#016x, vertexNum=%d): error=%v",
				vertex, h3, i, err)
			continue
		}

		almostEqual := geoAlmostEqualThreshold(&gb.Verts[i], &coord, 0.000001)
		if !almostEqual {
			t.Errorf("Vertex coordinates mismatch for cell %#016x (res=%d, %s), vertexNum=%d:\n"+
				"  vertex index: %#016x\n"+
				"  expected (boundary): lat=%.9f, lng=%.9f\n"+
				"  got (from vertex):   lat=%.9f, lng=%.9f\n"+
				"  difference: lat=%.9e, lng=%.9e",
				h3, getResolution(h3),
				map[bool]string{true: "pentagon", false: "hexagon"}[isPent],
				i, vertex,
				gb.Verts[i].Lat.Deg(), gb.Verts[i].Lng.Deg(),
				coord.Lat.Deg(), coord.Lng.Deg(),
				(gb.Verts[i].Lat - coord.Lat).Deg(),
				(gb.Verts[i].Lng - coord.Lng).Deg())
		}
	}
}

func cellToVertex_uniqueness_assertions(t *testing.T, h3 H3Index) {
	t.Helper()

	var originVerts [NUM_HEX_VERTS]H3Index
	if err := cellToVertexes(h3, &originVerts); err != E_SUCCESS {
		t.Skipf("Failed to get vertexes for cell %#016x (res=%d): %v", h3, getResolution(h3), err)
		return
	}

	isPent := isPentagon(h3)
	for v1 := 0; v1 < NUM_HEX_VERTS-1; v1++ {
		for v2 := v1 + 1; v2 < NUM_HEX_VERTS; v2++ {
			if originVerts[v1] != H3_NULL && originVerts[v2] != H3_NULL && originVerts[v1] == originVerts[v2] {
				t.Errorf("Duplicate vertex found for cell %#016x (res=%d, %s):\n"+
					"  vertex[%d] = %#016x\n"+
					"  vertex[%d] = %#016x (duplicate)",
					h3, getResolution(h3),
					map[bool]string{true: "pentagon", false: "hexagon"}[isPent],
					v1, originVerts[v1], v2, originVerts[v2])
			}
		}
	}
}

func cellToVertex_validity_assertions(t *testing.T, h3 H3Index) {
	t.Helper()

	var verts [NUM_HEX_VERTS]H3Index
	if err := cellToVertexes(h3, &verts); err != E_SUCCESS {
		t.Skipf("Failed to get vertexes for cell %#016x (res=%d): %v", h3, getResolution(h3), err)
		return
	}

	isPent := isPentagon(h3)
	maxVerts := NUM_HEX_VERTS
	if isPent {
		maxVerts = NUM_PENT_VERTS
	}

	for i := 0; i < NUM_HEX_VERTS; i++ {
		if verts[i] != H3_NULL {
			if !isValidVertex(verts[i]) {
				t.Errorf("Invalid vertex for cell %#016x (res=%d, %s):\n"+
					"  position: %d (of %d vertices)\n"+
					"  vertex index: %#016x\n"+
					"  vertex mode: %#x (expected %#x)\n"+
					"  vertex reserved bits: %d",
					h3, getResolution(h3),
					map[bool]string{true: "pentagon", false: "hexagon"}[isPent],
					i, maxVerts, verts[i],
					getMode(verts[i]), H3_VERTEX_MODE,
					getReservedBits(verts[i]))
			}
		}
	}
}

func cellToVertex_neighbor_assertions(t *testing.T, h3 H3Index) {
	t.Helper()

	neighbors := make([]H3Index, 7)
	var originVerts [NUM_HEX_VERTS]H3Index
	var neighborVerts [NUM_HEX_VERTS]H3Index

	if err := gridDisk(h3, 1, neighbors); err != E_SUCCESS {
		t.Skipf("Failed to get neighbors for cell %#016x (res=%d): %v", h3, getResolution(h3), err)
		return
	}

	if err := cellToVertexes(h3, &originVerts); err != E_SUCCESS {
		t.Skipf("Failed to get vertexes for cell %#016x (res=%d): %v", h3, getResolution(h3), err)
		return
	}

	isPent := isPentagon(h3)
	for i := 0; i < 7; i++ {
		neighbor := neighbors[i]
		if neighbor == H3_NULL || neighbor == h3 {
			continue
		}

		if err := cellToVertexes(neighbor, &neighborVerts); err != E_SUCCESS {
			continue
		}

		// Calculate the set intersection
		intersection := 0
		var sharedVerts []H3Index
		for v1 := 0; v1 < NUM_HEX_VERTS; v1++ {
			for v2 := 0; v2 < NUM_HEX_VERTS; v2++ {
				if neighborVerts[v1] != H3_NULL && originVerts[v2] != H3_NULL && neighborVerts[v1] == originVerts[v2] {
					intersection++
					sharedVerts = append(sharedVerts, neighborVerts[v1])
				}
			}
		}

		if intersection != 2 {
			neighborIsPent := isPentagon(neighbor)
			t.Errorf("Vertex sharing mismatch between cells:\n"+
				"  origin: %#016x (res=%d, %s)\n"+
				"  neighbor[%d]: %#016x (res=%d, %s)\n"+
				"  expected shared vertices: 2\n"+
				"  actual shared vertices: %d\n"+
				"  shared vertex indices: %v",
				h3, getResolution(h3),
				map[bool]string{true: "pentagon", false: "hexagon"}[isPent],
				i, neighbor, getResolution(neighbor),
				map[bool]string{true: "pentagon", false: "hexagon"}[neighborIsPent],
				intersection, sharedVerts)
		}
	}
}

func TestDirectionForVertexNum_Symmetry(t *testing.T) {
	t.Parallel()

	iterateAllIndexesAtRes(t, 0, directionForVertexNum_symmetry_assertions)
	iterateAllIndexesAtRes(t, 1, directionForVertexNum_symmetry_assertions)
	iterateAllIndexesAtRes(t, 2, directionForVertexNum_symmetry_assertions)
	iterateAllIndexesAtRes(t, 3, directionForVertexNum_symmetry_assertions)
	iterateAllIndexesAtRes(t, 4, directionForVertexNum_symmetry_assertions)
}

func TestCellToVertex_Point(t *testing.T) {
	t.Parallel()

	iterateAllIndexesAtRes(t, 0, cellToVertex_point_assertions)
	iterateAllIndexesAtRes(t, 1, cellToVertex_point_assertions)
	iterateAllIndexesAtRes(t, 2, cellToVertex_point_assertions)
	iterateAllIndexesAtRes(t, 3, cellToVertex_point_assertions)
	iterateAllIndexesAtRes(t, 4, cellToVertex_point_assertions)

	// Res 5: normal base cell
	iterateBaseCellIndexesAtRes(t, 5, cellToVertex_point_assertions, 0)
	// Res 5: pentagon base cell
	iterateBaseCellIndexesAtRes(t, 5, cellToVertex_point_assertions, 14)
	// Res 5: polar pentagon base cell
	iterateBaseCellIndexesAtRes(t, 5, cellToVertex_point_assertions, 117)
}

func TestCellToVertex_Neighbors(t *testing.T) {
	t.Parallel()

	iterateAllIndexesAtRes(t, 0, cellToVertex_neighbor_assertions)
	iterateAllIndexesAtRes(t, 1, cellToVertex_neighbor_assertions)
	iterateAllIndexesAtRes(t, 2, cellToVertex_neighbor_assertions)
	iterateAllIndexesAtRes(t, 3, cellToVertex_neighbor_assertions)
	iterateAllIndexesAtRes(t, 4, cellToVertex_neighbor_assertions)
}

func TestCellToVertex_Uniqueness(t *testing.T) {
	t.Parallel()

	iterateAllIndexesAtRes(t, 0, cellToVertex_uniqueness_assertions)
	iterateAllIndexesAtRes(t, 1, cellToVertex_uniqueness_assertions)
	iterateAllIndexesAtRes(t, 2, cellToVertex_uniqueness_assertions)
	iterateAllIndexesAtRes(t, 3, cellToVertex_uniqueness_assertions)
	iterateAllIndexesAtRes(t, 4, cellToVertex_uniqueness_assertions)
}

func TestCellToVertex_Validity(t *testing.T) {
	t.Parallel()

	iterateAllIndexesAtRes(t, 0, cellToVertex_validity_assertions)
	iterateAllIndexesAtRes(t, 1, cellToVertex_validity_assertions)
	iterateAllIndexesAtRes(t, 2, cellToVertex_validity_assertions)
	iterateAllIndexesAtRes(t, 3, cellToVertex_validity_assertions)
	iterateAllIndexesAtRes(t, 4, cellToVertex_validity_assertions)
}
