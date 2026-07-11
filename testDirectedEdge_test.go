// Tests ported from H3 v4.4.0: src/apps/testapps/testDirectedEdge.c.
package h3

import (
	"testing"
)

// Fixtures.
var sfGeo = LatLng{Lat: 0.659966917655, Lng: -2.1364398519396}

func TestAreNeighborCells(t *testing.T) {
	t.Parallel()

	var sf h3Index
	err := latLngToCell(&sfGeo, 9, &sf)
	if err != eSuccess {
		t.Fatalf("latLngToCell failed: %v", err)
	}

	ring := make([]h3Index, 7)
	err = gridRingUnsafe(sf, 1, ring)
	if err != eSuccess {
		t.Fatalf("gridRingUnsafe failed: %v", err)
	}

	isNeighbor, err := areNeighborCells(sf, sf)
	if err != eSuccess {
		t.Fatalf("areNeighborCells failed: %v", err)
	}
	if isNeighbor {
		t.Error("an index should not neighbor itself")
	}

	neighbors := 0
	var maxNeighborsSize int64
	err = maxGridDiskSize(1, &maxNeighborsSize)
	if err != eSuccess {
		t.Fatalf("maxGridDiskSize failed: %v", err)
	}
	for i := int64(0); i < maxNeighborsSize && i < int64(len(ring)); i++ {
		if ring[i] != 0 {
			isNeighbor, err := areNeighborCells(sf, ring[i])
			if err != eSuccess {
				t.Fatalf("areNeighborCells failed: %v", err)
			}
			if isNeighbor {
				neighbors++
			}
		}
	}
	if neighbors != 6 {
		t.Errorf("Expected 6 neighbors from a k-ring of 1, got %d", neighbors)
	}

	largerRing := make([]h3Index, 19)
	err = gridRingUnsafe(sf, 2, largerRing)
	if err != eSuccess {
		t.Fatalf("gridRingUnsafe failed: %v", err)
	}

	neighbors = 0
	err = maxGridDiskSize(2, &maxNeighborsSize)
	if err != eSuccess {
		t.Fatalf("maxGridDiskSize failed: %v", err)
	}
	for i := int64(0); i < maxNeighborsSize && i < int64(len(largerRing)); i++ {
		if largerRing[i] != 0 {
			isNeighbor, err := areNeighborCells(sf, largerRing[i])
			if err != eSuccess {
				t.Fatalf("areNeighborCells failed: %v", err)
			}
			if isNeighbor {
				neighbors++
			}
		}
	}
	if neighbors != 0 {
		t.Errorf("Expected 0 neighbors from a k-ring of 2, got %d", neighbors)
	}

	// Test invalid cells
	sfBroken := sf
	sfBroken = (sfBroken & ^h3Index(0xF<<59)) | (h3Index(h3DirectededgeMode) << 59)
	_, err = areNeighborCells(sf, sfBroken)
	if err != eCellInvalid {
		t.Error("Expected eCellInvalid for broken H3Indexes")
	}
	_, err = areNeighborCells(sfBroken, sf)
	if err != eCellInvalid {
		t.Error("Expected eCellInvalid for broken H3Indexes (reversed)")
	}

	var sfBigger h3Index
	err = latLngToCell(&sfGeo, 7, &sfBigger)
	if err != eSuccess {
		t.Fatalf("latLngToCell failed: %v", err)
	}
	_, err = areNeighborCells(sf, sfBigger)
	if err != eResMismatch {
		t.Error("Expected eResMismatch for hexagons of different resolution")
	}

	if len(ring) >= 3 {
		isNeighbor, err := areNeighborCells(ring[2], ring[1])
		if err != eSuccess {
			t.Fatalf("areNeighborCells failed: %v", err)
		}
		if !isNeighbor {
			t.Error("Hexagons in a ring should be neighbors")
		}
	}
}

func TestAreNeighborCells_invalid(t *testing.T) {
	t.Parallel()

	// Invalid index digit in origin is rejected.
	var origin h3Index
	setH3Index(&origin, 5, 0, int32(centerDigit))
	dest := origin
	origin = h3SetIndexDigit(origin, 5, int32(invalidDigit))
	dest = h3SetIndexDigit(dest, 5, int32(jkAxesDigit))

	if _, err := areNeighborCells(origin, dest); err != eCellInvalid {
		t.Errorf("Invalid index digit origin: expected eCellInvalid, got %v", err)
	}

	// Invalid k subsequence origin is rejected (pentagon base cell 4).
	setH3Index(&origin, 5, 4, int32(centerDigit))
	dest = origin
	origin = h3SetIndexDigit(origin, 5, int32(kAxesDigit))
	dest = h3SetIndexDigit(dest, 5, int32(ikAxesDigit))

	if _, err := areNeighborCells(origin, dest); err != eCellInvalid {
		t.Errorf("Invalid k subsequence origin: expected eCellInvalid, got %v", err)
	}

	// Invalid k subsequence destination is rejected.
	setH3Index(&origin, 5, 4, int32(centerDigit))
	dest = origin
	origin = h3SetIndexDigit(origin, 5, int32(ikAxesDigit))
	dest = h3SetIndexDigit(dest, 5, int32(kAxesDigit))

	if _, err := areNeighborCells(origin, dest); err != eCellInvalid {
		t.Errorf("Invalid k subsequence destination: expected eCellInvalid, got %v", err)
	}
}

func Test_debugIsValidCell(t *testing.T) {
	// Create the problematic cell from our parity test
	var origin h3Index
	setH3Index(&origin, 5, 0, int32(centerDigit))

	// Set invalid digit (invalidDigit = 7) at resolution 5
	origin = h3SetIndexDigit(origin, 5, int32(invalidDigit))

	t.Logf("Debug isValidCell for cell: %x", origin)
	t.Logf("  Mode: %d", getMode(origin))
	t.Logf("  Resolution: %d", getResolution(origin))
	t.Logf("  BaseCell: %d", getBaseCell(origin))

	// Test each validation step
	t.Logf("  _hasGoodTopBits: %v", _hasGoodTopBits(origin))
	t.Logf("  BaseCell < numBaseCells: %v (%d < %d)", getBaseCell(origin) < numBaseCells, getBaseCell(origin), numBaseCells)
	t.Logf("  _hasAny7UptoRes(h, res): %v", _hasAny7UptoRes(origin, getResolution(origin)))
	t.Logf("  _hasAll7AfterRes(h, res): %v", _hasAll7AfterRes(origin, getResolution(origin)))
	t.Logf("  _hasDeletedSubsequence(h, bc): %v", _hasDeletedSubsequence(origin, getBaseCell(origin)))

	t.Logf("  Overall isValidCell: %v", isValidCell(origin))

	// Let's also check what digits are actually in the cell
	res := getResolution(origin)
	t.Logf("  Digits from 1 to %d:", res)
	for r := int32(1); r <= res; r++ {
		digit := h3GetIndexDigit(origin, r)
		t.Logf("    Digit at res %d: %d", r, digit)
	}
}

func TestCellsToDirectedEdgeAndFriends(t *testing.T) {
	t.Parallel()

	var sf h3Index
	err := latLngToCell(&sfGeo, 9, &sf)
	if err != eSuccess {
		t.Fatalf("latLngToCell failed: %v", err)
	}

	ring := make([]h3Index, 7)
	err = gridRingUnsafe(sf, 1, ring)
	if err != eSuccess {
		t.Fatalf("gridRingUnsafe failed: %v", err)
	}
	if len(ring) == 0 {
		t.Fatal("Ring should not be empty")
	}

	var sf2 h3Index
	for _, cell := range ring {
		if cell != 0 && cell != sf {
			sf2 = cell
			break
		}
	}
	if sf2 == 0 {
		t.Fatal("Could not find valid neighbor cell")
	}

	edge, err := cellsToDirectedEdge(sf, sf2)
	if err != eSuccess {
		t.Fatalf("cellsToDirectedEdge failed: %v", err)
	}

	edgeOrigin, err := getDirectedEdgeOrigin(edge)
	if err != eSuccess {
		t.Fatalf("getDirectedEdgeOrigin failed: %v", err)
	}
	if sf != edgeOrigin {
		t.Error("Should be able to retrieve the origin from the edge")
	}

	var edgeDestination h3Index
	err = getDirectedEdgeDestination(edge, &edgeDestination)
	if err != eSuccess {
		t.Fatalf("getDirectedEdgeDestination failed: %v", err)
	}
	if sf2 != edgeDestination {
		t.Error("Should be able to retrieve the destination from the edge")
	}

	originDestination := make([]h3Index, 2)
	err = directedEdgeToCells(edge, originDestination)
	if err != eSuccess {
		t.Fatalf("directedEdgeToCells failed: %v", err)
	}
	if originDestination[0] != sf {
		t.Error("Expected origin first in the pair request")
	}
	if originDestination[1] != sf2 {
		t.Error("Expected destination last in the pair request")
	}

	err = directedEdgeToCells(0, originDestination)
	if err != eDirEdgeInvalid {
		t.Error("Expected eDirEdgeInvalid for invalid edges")
	}

	// Create invalid edge
	var invalidEdge h3Index
	setH3Index(&invalidEdge, 1, 4, 0)
	invalidEdge = (invalidEdge & ^h3Index(0x7<<56)) | (h3Index(invalidDigit) << 56)       // Set reserved bits
	invalidEdge = (invalidEdge & ^h3Index(0xF<<59)) | (h3Index(h3DirectededgeMode) << 59) // Set mode
	err = directedEdgeToCells(invalidEdge, originDestination)
	if err == eSuccess {
		t.Error("Expected error for invalid edge")
	}

	largerRing := make([]h3Index, 19)
	err = gridRingUnsafe(sf, 2, largerRing)
	if err != eSuccess {
		t.Fatalf("gridRingUnsafe failed: %v", err)
	}

	var sf3 h3Index
	for _, cell := range largerRing {
		if cell != 0 {
			sf3 = cell
			break
		}
	}
	if sf3 == 0 {
		t.Fatal("Could not find valid cell in larger ring")
	}

	_, err = cellsToDirectedEdge(sf, sf3)
	if err != eNotNeighbors {
		t.Error("Expected eNotNeighbors for non-neighbors")
	}
}

func TestGetDirectedEdgeOriginBadInput(t *testing.T) {
	t.Parallel()

	hexagon := h3Index(0x891ea6d6533ffff)

	_, err := getDirectedEdgeOrigin(hexagon)
	if err != eDirEdgeInvalid {
		t.Error("Expected eDirEdgeInvalid for hexagon index")
	}

	_, err = getDirectedEdgeOrigin(0)
	if err != eDirEdgeInvalid {
		t.Error("Expected eDirEdgeInvalid for null index")
	}
}

func TestGetDirectedEdgeOriginBadInput2(t *testing.T) {
	t.Parallel()

	var sf h3Index
	err := latLngToCell(&sfGeo, 9, &sf)
	if err != eSuccess {
		t.Fatalf("latLngToCell failed: %v", err)
	}

	ring := make([]h3Index, 7)
	err = gridRingUnsafe(sf, 1, ring)
	if err != eSuccess {
		t.Fatalf("gridRingUnsafe failed: %v", err)
	}

	var sf2 h3Index
	for _, cell := range ring {
		if cell != 0 && cell != sf {
			sf2 = cell
			break
		}
	}
	if sf2 == 0 {
		t.Fatal("Could not find valid neighbor cell")
	}

	edge, err := cellsToDirectedEdge(sf, sf2)
	if err != eSuccess {
		t.Fatalf("cellsToDirectedEdge failed: %v", err)
	}

	// Set invalid reserved bits
	edge = (edge & ^h3Index(0x7<<56)) | (h3Index(invalidDigit) << 56)
	var out h3Index
	err = getDirectedEdgeDestination(edge, &out)
	if err != eFailed {
		t.Error("Expected eFailed for invalid directed edge")
	}
}

func TestGetDirectedEdgeDestination(t *testing.T) {
	t.Parallel()

	hexagon := h3Index(0x891ea6d6533ffff)

	var out h3Index
	err := getDirectedEdgeDestination(hexagon, &out)
	if err != eDirEdgeInvalid {
		t.Error("Expected eDirEdgeInvalid for hexagon index")
	}

	err = getDirectedEdgeDestination(0, &out)
	if err != eDirEdgeInvalid {
		t.Error("Expected eDirEdgeInvalid for null index")
	}
}

func TestCellsToDirectedEdgeFromPentagon(t *testing.T) {
	t.Parallel()

	for res := int32(0); res <= maxH3Res; res++ {
		pentagons := make([]h3Index, numPentagons)
		err := getPentagons(res, pentagons)
		if err != eSuccess {
			t.Fatalf("getPentagons failed for res %d: %v", res, err)
		}

		for _, pentagon := range pentagons {
			if pentagon == 0 {
				continue
			}

			ring := make([]h3Index, 7)
			err := gridDisk(pentagon, 1, ring)
			if err != eSuccess {
				t.Fatalf("gridDisk failed: %v", err)
			}

			for _, neighbor := range ring {
				if neighbor == pentagon || neighbor == 0 {
					continue
				}

				edge, err := cellsToDirectedEdge(pentagon, neighbor)
				if err != eSuccess {
					t.Fatalf("cellsToDirectedEdge failed for pentagon to neighbor: %v", err)
				}
				if !isValidDirectedEdge(edge) {
					t.Error("Pentagon-to-neighbor should be a valid edge")
				}

				edge, err = cellsToDirectedEdge(neighbor, pentagon)
				if err != eSuccess {
					t.Fatalf("cellsToDirectedEdge failed for neighbor to pentagon: %v", err)
				}
				if !isValidDirectedEdge(edge) {
					t.Error("Neighbor-to-pentagon should be a valid edge")
				}
			}
		}
	}
}

func TestIsValidDirectedEdge(t *testing.T) {
	t.Parallel()

	var sf h3Index
	err := latLngToCell(&sfGeo, 9, &sf)
	if err != eSuccess {
		t.Fatalf("latLngToCell failed: %v", err)
	}

	ring := make([]h3Index, 7)
	err = gridRingUnsafe(sf, 1, ring)
	if err != eSuccess {
		t.Fatalf("gridRingUnsafe failed: %v", err)
	}

	var sf2 h3Index
	for _, cell := range ring {
		if cell != 0 && cell != sf {
			sf2 = cell
			break
		}
	}
	if sf2 == 0 {
		t.Fatal("Could not find valid neighbor cell")
	}

	edge, err := cellsToDirectedEdge(sf, sf2)
	if err != eSuccess {
		t.Fatalf("cellsToDirectedEdge failed: %v", err)
	}

	if !isValidDirectedEdge(edge) {
		t.Error("Edges should validate correctly")
	}
	if isValidDirectedEdge(sf) {
		t.Error("Hexagons should not validate")
	}

	// Test undirected edge
	undirectedEdge := edge
	undirectedEdge = (undirectedEdge & ^h3Index(0xF<<59)) | (h3Index(h3EdgeMode) << 59)
	if isValidDirectedEdge(undirectedEdge) {
		t.Error("Undirected edges should not validate")
	}

	// Test hexagon with reserved bits
	hexagonWithReserved := sf
	hexagonWithReserved = (hexagonWithReserved & ^h3Index(0x7<<56)) | (h3Index(1) << 56)
	if isValidDirectedEdge(hexagonWithReserved) {
		t.Error("Hexagons with reserved bits should not validate")
	}

	// Test fake edge (no edge specified)
	fakeEdge := sf
	fakeEdge = (fakeEdge & ^h3Index(0xF<<59)) | (h3Index(h3DirectededgeMode) << 59)
	if isValidDirectedEdge(fakeEdge) {
		t.Error("Edges without an edge specified should not work")
	}

	// Test invalid edge
	invalidEdge := sf
	invalidEdge = (invalidEdge & ^h3Index(0xF<<59)) | (h3Index(h3DirectededgeMode) << 59)
	invalidEdge = (invalidEdge & ^h3Index(0x7<<56)) | (h3Index(invalidDigit) << 56)
	if isValidDirectedEdge(invalidEdge) {
		t.Error("Edges with an invalid edge specified should not work")
	}

	// Test pentagonal edge
	pentagon := h3Index(0x821c07fffffffff)
	goodPentagonalEdge := pentagon
	goodPentagonalEdge = (goodPentagonalEdge & ^h3Index(0xF<<59)) | (h3Index(h3DirectededgeMode) << 59)
	goodPentagonalEdge = (goodPentagonalEdge & ^h3Index(0x7<<56)) | (h3Index(2) << 56)
	if !isValidDirectedEdge(goodPentagonalEdge) {
		t.Error("Pentagonal edge should validate")
	}

	// Test bad pentagonal edge (missing edge)
	badPentagonalEdge := goodPentagonalEdge
	badPentagonalEdge = (badPentagonalEdge & ^h3Index(0x7<<56)) | (h3Index(1) << 56)
	if isValidDirectedEdge(badPentagonalEdge) {
		t.Error("Missing pentagonal edge should not validate")
	}

	// Test high bit edge
	highBitEdge := edge
	highBitEdge = (highBitEdge & ^h3Index(1<<63)) | (h3Index(1) << 63)
	if isValidDirectedEdge(highBitEdge) {
		t.Error("High bit set edge should not validate")
	}
}

func TestOriginToDirectedEdges(t *testing.T) {
	t.Parallel()

	var sf h3Index
	err := latLngToCell(&sfGeo, 9, &sf)
	if err != eSuccess {
		t.Fatalf("latLngToCell failed: %v", err)
	}

	edges := make([]h3Index, 6)
	err = originToDirectedEdges(sf, edges)
	if err != eSuccess {
		t.Fatalf("originToDirectedEdges failed: %v", err)
	}

	for i, edge := range edges {
		if !isValidDirectedEdge(edge) {
			t.Errorf("Edge %d should be valid", i)
		}

		origin, err := getDirectedEdgeOrigin(edge)
		if err != eSuccess {
			t.Fatalf("getDirectedEdgeOrigin failed: %v", err)
		}
		if sf != origin {
			t.Errorf("Edge %d origin should be correct", i)
		}

		var destination h3Index
		err = getDirectedEdgeDestination(edge, &destination)
		if err != eSuccess {
			t.Fatalf("getDirectedEdgeDestination failed: %v", err)
		}
		if sf == destination {
			t.Errorf("Edge %d destination should not be origin", i)
		}
	}
}

func TestGetH3DirectedEdgesFromPentagon(t *testing.T) {
	t.Parallel()

	pentagon := h3Index(0x821c07fffffffff)
	edges := make([]h3Index, 6)
	err := originToDirectedEdges(pentagon, edges)
	if err != eSuccess {
		t.Fatalf("originToDirectedEdges failed: %v", err)
	}

	missingEdgeCount := 0
	for i, edge := range edges {
		if edge == 0 {
			missingEdgeCount++
		} else {
			if !isValidDirectedEdge(edge) {
				t.Errorf("Edge %d should be valid", i)
			}

			origin, err := getDirectedEdgeOrigin(edge)
			if err != eSuccess {
				t.Fatalf("getDirectedEdgeOrigin failed: %v", err)
			}
			if pentagon != origin {
				t.Errorf("Edge %d origin should be correct", i)
			}

			var destination h3Index
			err = getDirectedEdgeDestination(edge, &destination)
			if err != eSuccess {
				t.Fatalf("getDirectedEdgeDestination failed: %v", err)
			}
			if pentagon == destination {
				t.Errorf("Edge %d destination should not be origin", i)
			}
		}
	}

	if missingEdgeCount != 1 {
		t.Errorf("Expected 1 missing edge for pentagon, got %d", missingEdgeCount)
	}
}

func TestDirectedEdgeToBoundary(t *testing.T) {
	t.Parallel()

	expectedVertices := [][]int{
		{3, 4}, {1, 2}, {2, 3},
		{5, 0}, {4, 5}, {0, 1},
	}

	for res := int32(0); res <= maxH3Res; res++ {
		var sf h3Index
		err := latLngToCell(&sfGeo, res, &sf)
		if err != eSuccess {
			t.Fatalf("latLngToCell failed: %v", err)
		}

		var boundary CellBoundary
		err = cellToBoundary(sf, &boundary)
		if err != eSuccess {
			t.Fatalf("cellToBoundary failed: %v", err)
		}

		edges := make([]h3Index, 6)
		err = originToDirectedEdges(sf, edges)
		if err != eSuccess {
			t.Fatalf("originToDirectedEdges failed: %v", err)
		}

		for i, edge := range edges {
			var edgeBoundary CellBoundary
			err := directedEdgeToBoundary(edge, &edgeBoundary)
			if err != eSuccess {
				t.Fatalf("directedEdgeToBoundary failed: %v", err)
			}

			if edgeBoundary.numVerts != 2 {
				t.Errorf("Expected 2 vertices for edge %d, got %d", i, edgeBoundary.numVerts)
				continue
			}

			for j := int32(0); j < edgeBoundary.numVerts; j++ {
				expectedIdx := expectedVertices[i][j]
				if !geoAlmostEqual(&edgeBoundary.verts[j], &boundary.verts[expectedIdx]) {
					t.Errorf("Edge %d vertex %d does not match expected", i, j)
				}
			}
		}
	}
}

func TestDirectedEdgeToBoundaryPentagonClassIII(t *testing.T) {
	t.Parallel()

	expectedVertices := [][]int{
		{-1, -1, -1}, {2, 3, 4}, {4, 5, 6},
		{8, 9, 0}, {6, 7, 8}, {0, 1, 2},
	}

	for res := int32(1); res <= maxH3Res; res += 2 { // Only odd resolutions (Class III)
		var pentagon h3Index
		setH3Index(&pentagon, res, 24, 0)
		var boundary CellBoundary
		err := cellToBoundary(pentagon, &boundary)
		if err != eSuccess {
			t.Fatalf("cellToBoundary failed: %v", err)
		}

		edges := make([]h3Index, 6)
		err = originToDirectedEdges(pentagon, edges)
		if err != eSuccess {
			t.Fatalf("originToDirectedEdges failed: %v", err)
		}

		missingEdgeCount := 0
		for i, edge := range edges {
			if edge == 0 {
				missingEdgeCount++
				continue
			}

			var edgeBoundary CellBoundary
			err := directedEdgeToBoundary(edge, &edgeBoundary)
			if err != eSuccess {
				t.Fatalf("directedEdgeToBoundary failed: %v", err)
			}

			if edgeBoundary.numVerts != 3 {
				t.Errorf("Expected 3 vertices for Class III pentagon edge %d, got %d", i, edgeBoundary.numVerts)
				continue
			}

			for j := int32(0); j < edgeBoundary.numVerts; j++ {
				expectedIdx := expectedVertices[i][j]
				if expectedIdx >= 0 && !geoAlmostEqual(&edgeBoundary.verts[j], &boundary.verts[expectedIdx]) {
					t.Errorf("Pentagon edge %d vertex %d does not match expected", i, j)
				}
			}
		}

		if missingEdgeCount != 1 {
			t.Errorf("Expected 1 missing edge for pentagon, got %d", missingEdgeCount)
		}
	}
}

func TestDirectedEdgeToBoundaryPentagonClassII(t *testing.T) {
	t.Parallel()

	expectedVertices := [][]int{
		{-1, -1}, {1, 2}, {2, 3},
		{4, 0}, {3, 4}, {0, 1},
	}

	for res := int32(0); res <= maxH3Res; res += 2 { // Only even resolutions (Class II)
		var pentagon h3Index
		setH3Index(&pentagon, res, 24, 0)
		var boundary CellBoundary
		err := cellToBoundary(pentagon, &boundary)
		if err != eSuccess {
			t.Fatalf("cellToBoundary failed: %v", err)
		}

		edges := make([]h3Index, 6)
		err = originToDirectedEdges(pentagon, edges)
		if err != eSuccess {
			t.Fatalf("originToDirectedEdges failed: %v", err)
		}

		missingEdgeCount := 0
		for i, edge := range edges {
			if edge == 0 {
				missingEdgeCount++
				continue
			}

			var edgeBoundary CellBoundary
			err := directedEdgeToBoundary(edge, &edgeBoundary)
			if err != eSuccess {
				t.Fatalf("directedEdgeToBoundary failed: %v", err)
			}

			if edgeBoundary.numVerts != 2 {
				t.Errorf("Expected 2 vertices for Class II pentagon edge %d, got %d", i, edgeBoundary.numVerts)
				continue
			}

			for j := int32(0); j < edgeBoundary.numVerts; j++ {
				expectedIdx := expectedVertices[i][j]
				if expectedIdx >= 0 && !geoAlmostEqual(&edgeBoundary.verts[j], &boundary.verts[expectedIdx]) {
					t.Errorf("Pentagon edge %d vertex %d does not match expected", i, j)
				}
			}
		}

		if missingEdgeCount != 1 {
			t.Errorf("Expected 1 missing edge for pentagon, got %d", missingEdgeCount)
		}
	}
}

func TestDirectedEdgeToBoundary_invalid(t *testing.T) {
	t.Parallel()

	var sf h3Index
	err := latLngToCell(&sfGeo, 9, &sf)
	if err != eSuccess {
		t.Fatalf("latLngToCell failed: %v", err)
	}

	// Create invalid edge (no edge direction specified)
	invalidEdge := sf
	invalidEdge = (invalidEdge & ^h3Index(0xF<<59)) | (h3Index(h3DirectededgeMode) << 59)
	var cb CellBoundary
	err = directedEdgeToBoundary(invalidEdge, &cb)
	if err != eDirEdgeInvalid {
		t.Error("Expected eDirEdgeInvalid for invalid edge direction")
	}

	// Create invalid edge with bad base cell and reserved bits
	invalidEdge2 := sf
	invalidEdge2 = (invalidEdge2 & ^h3Index(0x7<<56)) | (h3Index(1) << 56)                  // Set reserved bits
	invalidEdge2 = (invalidEdge2 & ^h3Index(0x7F<<45)) | (h3Index(numBaseCells+1) << 45)    // Set invalid base cell
	invalidEdge2 = (invalidEdge2 & ^h3Index(0xF<<59)) | (h3Index(h3DirectededgeMode) << 59) // Set mode
	err = directedEdgeToBoundary(invalidEdge2, &cb)
	if err == eSuccess {
		t.Error("Expected error for invalid edge indexing digit")
	}
}

func TestEdgeLength_invalid(t *testing.T) {
	t.Parallel()

	// Test invalid edge (null)
	var length float64
	err := edgeLengthRads(0, &length)
	if err != eDirEdgeInvalid {
		t.Error("Expected eDirEdgeInvalid for null edge")
	}

	// Test non-edge (cell)
	zero := LatLng{Lat: 0, Lng: 0}
	var h3 h3Index
	err = latLngToCell(&zero, 0, &h3)
	if err != eSuccess {
		t.Fatalf("latLngToCell failed: %v", err)
	}

	err = edgeLengthRads(h3, &length)
	if err != eDirEdgeInvalid {
		t.Error("Expected eDirEdgeInvalid for non-edge (cell)")
	}
}
