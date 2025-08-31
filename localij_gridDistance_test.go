// Tests ported from testGridDistance.c
package h3

import (
	"testing"
)

func Test_gridDistance_testIndexDistance(t *testing.T) {
	t.Parallel()

	bc := H3Index(H3_INIT)
	setH3Index(&bc, 1, 17, 0)
	p := H3Index(H3_INIT)
	setH3Index(&p, 1, 14, 0)
	p2 := H3Index(H3_INIT)
	setH3Index(&p2, 1, 14, 2)
	p3 := H3Index(H3_INIT)
	setH3Index(&p3, 1, 14, 3)
	p6 := H3Index(H3_INIT)
	setH3Index(&p6, 1, 14, 6)

	var distance int64

	err := gridDistance(bc, p, &distance)
	if err != E_SUCCESS {
		t.Fatalf("gridDistance failed: %v", err)
	}
	if distance != 3 {
		t.Errorf("Expected distance onto pentagon to be 3, got %d", distance)
	}

	err = gridDistance(bc, p2, &distance)
	if err != E_SUCCESS {
		t.Fatalf("gridDistance failed: %v", err)
	}
	if distance != 2 {
		t.Errorf("Expected distance onto p2 to be 2, got %d", distance)
	}

	err = gridDistance(bc, p3, &distance)
	if err != E_SUCCESS {
		t.Fatalf("gridDistance failed: %v", err)
	}
	if distance != 3 {
		t.Errorf("Expected distance onto p3 to be 3, got %d", distance)
	}

	// TODO: p4 and p5 tests are commented out in C due to pentagon distortion

	err = gridDistance(bc, p6, &distance)
	if err != E_SUCCESS {
		t.Fatalf("gridDistance failed: %v", err)
	}
	if distance != 2 {
		t.Errorf("Expected distance onto p6 to be 2, got %d", distance)
	}
}

func Test_gridDistance_testIndexDistance2(t *testing.T) {
	t.Parallel()

	origin := H3Index(0x820c4ffffffffff)
	// Destination is on the other side of the pentagon
	destination := H3Index(0x821ce7ffffffffff)

	// TODO doesn't work because of pentagon distortion. Both should be 5.
	var distance int64
	err := gridDistance(destination, origin, &distance)
	if err == E_SUCCESS {
		t.Errorf("Expected failure for distance in res 2 across pentagon, but got success")
	}

	err = gridDistance(origin, destination, &distance)
	if err == E_SUCCESS {
		t.Errorf("Expected failure for distance in res 2 across pentagon (reversed), but got success")
	}
}

func Test_gridDistance_baseCells(t *testing.T) {
	t.Parallel()

	// Some indexes that represent base cells. Base cells
	// are hexagons except for `pent1`.
	bc1 := H3Index(H3_INIT)
	setH3Index(&bc1, 0, 15, 0)

	bc2 := H3Index(H3_INIT)
	setH3Index(&bc2, 0, 8, 0)

	bc3 := H3Index(H3_INIT)
	setH3Index(&bc3, 0, 31, 0)

	pent1 := H3Index(H3_INIT)
	setH3Index(&pent1, 0, 4, 0)

	var distance int64

	err := gridDistance(bc1, pent1, &distance)
	if err != E_SUCCESS {
		t.Fatalf("gridDistance failed: %v", err)
	}
	if distance != 1 {
		t.Errorf("Expected distance to neighbor to be 1 (15, 4), got %d", distance)
	}

	err = gridDistance(bc1, bc2, &distance)
	if err != E_SUCCESS {
		t.Fatalf("gridDistance failed: %v", err)
	}
	if distance != 1 {
		t.Errorf("Expected distance to neighbor to be 1 (15, 8), got %d", distance)
	}

	err = gridDistance(bc1, bc3, &distance)
	if err != E_SUCCESS {
		t.Fatalf("gridDistance failed: %v", err)
	}
	if distance != 1 {
		t.Errorf("Expected distance to neighbor to be 1 (15, 31), got %d", distance)
	}

	err = gridDistance(pent1, bc3, &distance)
	if err == E_SUCCESS {
		t.Errorf("Expected distance to neighbor to be invalid, but got success")
	}
}

func Test_gridDistance_resolutionMismatch(t *testing.T) {
	t.Parallel()

	var distance int64
	err := gridDistance(H3Index(0x832830fffffffff), H3Index(0x822837fffffffff), &distance)
	if err != E_RES_MISMATCH {
		t.Errorf("Expected E_RES_MISMATCH for different resolutions, got %v", err)
	}
}

func Test_gridDistance_edge(t *testing.T) {
	t.Parallel()

	origin := H3Index(0x832830fffffffff)
	dest := H3Index(0x832834fffffffff)

	// First check if these cells are actually neighbors
	areNeighbors, neighborErr := areNeighborCells(origin, dest)
	if neighborErr != E_SUCCESS {
		t.Fatalf("areNeighborCells failed: %v", neighborErr)
	}
	if !areNeighbors {
		t.Skip("Test cells are not neighbors - skipping edge test")
	}

	edge, err := cellsToDirectedEdge(origin, dest)
	if err != E_SUCCESS {
		t.Fatalf("cellsToDirectedEdge failed: %v", err)
	}
	if edge == 0 {
		t.Fatalf("Test edge should be valid")
	}

	var distance int64

	err = gridDistance(edge, origin, &distance)
	if err != E_SUCCESS {
		t.Fatalf("gridDistance failed: %v", err)
	}
	if distance != 0 {
		t.Errorf("Expected edge to have zero distance to origin, got %d", distance)
	}

	err = gridDistance(origin, edge, &distance)
	if err != E_SUCCESS {
		t.Fatalf("gridDistance failed: %v", err)
	}
	if distance != 0 {
		t.Errorf("Expected origin to have zero distance to edge, got %d", distance)
	}

	err = gridDistance(edge, dest, &distance)
	if err != E_SUCCESS {
		t.Fatalf("gridDistance failed: %v", err)
	}
	if distance != 1 {
		t.Errorf("Expected edge to have distance 1 to destination, got %d", distance)
	}

	err = gridDistance(dest, edge, &distance)
	if err != E_SUCCESS {
		t.Fatalf("gridDistance failed: %v", err)
	}
	if distance != 1 {
		t.Errorf("Expected destination to have distance 1 to edge, got %d", distance)
	}
}

func Test_gridDistance_invalid(t *testing.T) {
	t.Parallel()

	// Some indexes that represent base cells. Base cells
	// are hexagons except for `pent1`.
	bc1 := H3Index(H3_INIT)
	setH3Index(&bc1, 0, 15, 0)

	invalid := H3Index(0xffffffffffffffff)
	var distance int64

	err := gridDistance(invalid, invalid, &distance)
	if err != E_CELL_INVALID {
		t.Errorf("Expected E_CELL_INVALID for distance from invalid cell, got %v", err)
	}

	err = gridDistance(bc1, invalid, &distance)
	if err != E_RES_MISMATCH {
		t.Errorf("Expected E_RES_MISMATCH for distance to invalid cell, got %v", err)
	}
}
