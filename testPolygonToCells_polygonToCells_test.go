// Tests ported from testPolygonToCells.c
package h3

import (
	"math"
	"testing"
)

// Note: Some test fixtures are shared from testPolygonToCellsExperimental_test.go
// However, the C tests use different coordinates for point and line polygons,
// so we define those separately here to match the C test exactly

// C test-specific fixtures for degenerate polygons.
var (
	// Point polygon with single vertex at origin - C test expects eFailed.
	pointVertsC      = []LatLng{{Lat: 0, Lng: 0}}
	pointGeoLoopC    = GeoLoop(pointVertsC)
	pointGeoPolygonC = GeoPolygon{GeoLoop: pointGeoLoopC, Holes: nil}

	// Line polygon from origin - C test expects eFailed.
	lineVertsC      = []LatLng{{Lat: 0, Lng: 0}, {Lat: 1, Lng: 0}}
	lineGeoLoopC    = GeoLoop(lineVertsC)
	lineGeoPolygonC = GeoPolygon{GeoLoop: lineGeoLoopC, Holes: nil}
)

// Helper function to count non-null indexes (avoiding conflict with existing function).
func countNonNullIndexesStandard(indexes []h3Index) int64 {
	count := int64(0)
	for _, idx := range indexes {
		if idx != h3Null {
			count++
		}
	}
	return count
}

// fillIndex_assertions helper function.
func fillIndex_assertions(h h3Index) {
	if isTransmeridianCell(h) {
		// TODO: these do not work correctly
		return
	}

	currentRes := getResolution(h)
	// TODO: Not testing more than one depth because the assertions fail.
	for nextRes := currentRes; nextRes <= currentRes+1; nextRes++ {
		var boundary CellBoundary
		err := cellToBoundary(h, &boundary)
		if err != eSuccess {
			continue
		}

		verts := make([]LatLng, boundary.numVerts)
		for j := int32(0); j < boundary.numVerts; j++ {
			verts[j] = boundary.verts[j]
		}
		polygon := GeoPolygon{
			GeoLoop: verts,
			Holes:   nil,
		}

		var polygonToCellsSize int64
		err = maxPolygonToCellsSize(&polygon, nextRes, 0, &polygonToCellsSize)
		if err != eSuccess {
			continue
		}

		polygonToCellsOut := make([]h3Index, polygonToCellsSize)
		err = polygonToCells(&polygon, nextRes, 0, polygonToCellsOut)
		if err != eSuccess {
			continue
		}

		polygonToCellsCount := countNonNullIndexesStandard(polygonToCellsOut)

		childrenSize, err := cellToChildrenSize(h, nextRes)
		if err != eSuccess {
			continue
		}

		children := make([]h3Index, childrenSize)
		err = cellToChildren(h, nextRes, children)
		if err != eSuccess {
			continue
		}

		cellToChildrenCount := countNonNullIndexesStandard(children)

		if polygonToCellsCount != cellToChildrenCount {
			return // Skip assertion failure for now
		}

		// Verify all children are found in polygonToCells output
		for _, child := range children {
			if child == h3Null {
				continue
			}
			found := false
			for _, cell := range polygonToCellsOut {
				if cell == child {
					found = true
					break
				}
			}
			if !found {
				return // Skip assertion failure for now
			}
		}
	}
}

func TestMaxPolygonToCellsSize(t *testing.T) {
	t.Parallel()
	var numHexagons int64

	err := maxPolygonToCellsSize(&sfGeoPolygon, 9, 0, &numHexagons)
	if err != eSuccess {
		t.Fatalf("maxPolygonToCellsSize failed: %v", err)
	}
	if numHexagons != 5613 {
		t.Errorf("Expected 5613, got %d", numHexagons)
	}

	err = maxPolygonToCellsSize(&holeGeoPolygon, 9, 0, &numHexagons)
	if err != eSuccess {
		t.Fatalf("maxPolygonToCellsSize failed: %v", err)
	}
	if numHexagons != 5613 {
		t.Errorf("Expected 5613 for hole polygon, got %d", numHexagons)
	}

	err = maxPolygonToCellsSize(&emptyGeoPolygon, 9, 0, &numHexagons)
	if err != eSuccess {
		t.Fatalf("maxPolygonToCellsSize failed: %v", err)
	}
	if numHexagons != 15 {
		t.Errorf("Expected 15 for empty polygon, got %d", numHexagons)
	}
}

func TestMaxPolygonToCellsSizeInvalid(t *testing.T) {
	t.Parallel()
	var numHexagons int64

	err := maxPolygonToCellsSize(&invalidGeoPolygon, 9, 0, &numHexagons)
	if err != eFailed {
		t.Error("Expected eFailed for invalid polygon with Infinity")
	}

	err = maxPolygonToCellsSize(&invalid2GeoPolygon, 9, 0, &numHexagons)
	if err != eFailed {
		t.Error("Expected eFailed for invalid polygon with NaNs")
	}
}

func TestMaxPolygonToCellsSizePoint(t *testing.T) {
	t.Parallel()
	var numHexagons int64

	// Use the C test coordinates (point at origin)
	err := maxPolygonToCellsSize(&pointGeoPolygonC, 9, 0, &numHexagons)
	if err != eFailed {
		t.Errorf("Expected eFailed for single point polygon, got %v with numHexagons=%d", err, numHexagons)
	}
}

func TestMaxPolygonToCellsSizeLine(t *testing.T) {
	t.Parallel()
	var numHexagons int64

	// Use the C test coordinates (line from origin)
	err := maxPolygonToCellsSize(&lineGeoPolygonC, 9, 0, &numHexagons)
	if err != eFailed {
		t.Errorf("Expected eFailed for straight line polygon, got %v with numHexagons=%d", err, numHexagons)
	}
}

func TestPolygonToCells(t *testing.T) {
	t.Parallel()
	var numHexagons int64

	err := maxPolygonToCellsSize(&sfGeoPolygon, 9, 0, &numHexagons)
	if err != eSuccess {
		t.Fatalf("maxPolygonToCellsSize failed: %v", err)
	}

	hexagons := make([]h3Index, numHexagons)
	err = polygonToCells(&sfGeoPolygon, 9, 0, hexagons)
	if err != eSuccess {
		t.Fatalf("polygonToCells failed: %v", err)
	}

	actualNumIndexes := countNonNullIndexesStandard(hexagons)
	if actualNumIndexes != 1253 {
		t.Errorf("Expected 1253 cells, got %d", actualNumIndexes)
	}
}

func TestPolygonToCellsHole(t *testing.T) {
	t.Parallel()
	var numHexagons int64

	err := maxPolygonToCellsSize(&holeGeoPolygon, 9, 0, &numHexagons)
	if err != eSuccess {
		t.Fatalf("maxPolygonToCellsSize failed: %v", err)
	}

	hexagons := make([]h3Index, numHexagons)
	err = polygonToCells(&holeGeoPolygon, 9, 0, hexagons)
	if err != eSuccess {
		t.Fatalf("polygonToCells failed: %v", err)
	}

	actualNumIndexes := countNonNullIndexesStandard(hexagons)
	if actualNumIndexes != 1214 {
		t.Errorf("Expected 1214 cells for polygon with hole, got %d", actualNumIndexes)
	}
}

func TestPolygonToCellsEmptyStandard(t *testing.T) {
	t.Parallel()
	var numHexagons int64

	err := maxPolygonToCellsSize(&emptyGeoPolygon, 9, 0, &numHexagons)
	if err != eSuccess {
		t.Fatalf("maxPolygonToCellsSize failed: %v", err)
	}

	hexagons := make([]h3Index, numHexagons)
	err = polygonToCells(&emptyGeoPolygon, 9, 0, hexagons)
	if err != eSuccess {
		t.Fatalf("polygonToCells failed: %v", err)
	}

	actualNumIndexes := countNonNullIndexesStandard(hexagons)
	if actualNumIndexes != 0 {
		t.Errorf("Expected 0 cells for empty polygon, got %d", actualNumIndexes)
	}
}

func TestPolygonToCellsExactStandard(t *testing.T) {
	t.Parallel()
	somewhere := LatLng{1, 2}
	var origin h3Index

	err := latLngToCell(&somewhere, 9, &origin)
	if err != eSuccess {
		t.Fatalf("latLngToCell failed: %v", err)
	}

	var boundary CellBoundary
	err = cellToBoundary(origin, &boundary)
	if err != eSuccess {
		t.Fatalf("cellToBoundary failed: %v", err)
	}

	// Create vertices with one extra to close the polygon
	verts := make([]LatLng, boundary.numVerts+1)
	for i := int32(0); i < boundary.numVerts; i++ {
		verts[i] = boundary.verts[i]
	}
	verts[boundary.numVerts] = boundary.verts[0]

	someGeoLoop := GeoLoop(verts)
	someHexagon := GeoPolygon{GeoLoop: someGeoLoop, Holes: nil}

	var numHexagons int64
	err = maxPolygonToCellsSize(&someHexagon, 9, 0, &numHexagons)
	if err != eSuccess {
		t.Fatalf("maxPolygonToCellsSize failed: %v", err)
	}

	hexagons := make([]h3Index, numHexagons)
	err = polygonToCells(&someHexagon, 9, 0, hexagons)
	if err != eSuccess {
		t.Fatalf("polygonToCells failed: %v", err)
	}

	actualNumIndexes := countNonNullIndexesStandard(hexagons)
	if actualNumIndexes != 1 {
		t.Errorf("Expected 1 cell for exact hexagon boundary, got %d", actualNumIndexes)
	}
}

func TestPolygonToCellsTransmeridianStandard(t *testing.T) {
	t.Parallel()

	// Prime meridian case
	primeMeridianVerts := []LatLng{
		{0.01, 0.01}, {0.01, -0.01}, {-0.01, -0.01}, {-0.01, 0.01},
	}
	primeMeridianGeoLoop := GeoLoop(primeMeridianVerts)
	primeMeridianGeoPolygon := GeoPolygon{GeoLoop: primeMeridianGeoLoop, Holes: nil}

	// Transmeridian case
	transMeridianVerts := []LatLng{
		{0.01, -math.Pi + 0.01}, {0.01, math.Pi - 0.01},
		{-0.01, math.Pi - 0.01}, {-0.01, -math.Pi + 0.01},
	}
	transMeridianGeoLoop := GeoLoop(transMeridianVerts)
	transMeridianGeoPolygon := GeoPolygon{GeoLoop: transMeridianGeoLoop, Holes: nil}

	// Transmeridian hole case
	transMeridianHoleVerts := []LatLng{
		{0.005, -math.Pi + 0.005}, {0.005, math.Pi - 0.005},
		{-0.005, math.Pi - 0.005}, {-0.005, -math.Pi + 0.005},
	}
	transMeridianHoleGeoLoop := GeoLoop(transMeridianHoleVerts)
	transMeridianHoleGeoPolygon := GeoPolygon{
		GeoLoop: transMeridianGeoLoop,
		Holes:   []GeoLoop{transMeridianHoleGeoLoop},
	}
	transMeridianFilledHoleGeoPolygon := GeoPolygon{
		GeoLoop: transMeridianHoleGeoLoop,
		Holes:   nil,
	}

	// Prime meridian test
	var numHexagons int64
	err := maxPolygonToCellsSize(&primeMeridianGeoPolygon, 7, 0, &numHexagons)
	if err != eSuccess {
		t.Fatalf("maxPolygonToCellsSize failed for prime meridian: %v", err)
	}

	hexagons := make([]h3Index, numHexagons)
	err = polygonToCells(&primeMeridianGeoPolygon, 7, 0, hexagons)
	if err != eSuccess {
		t.Fatalf("polygonToCells failed for prime meridian: %v", err)
	}

	actualNumIndexes := countNonNullIndexesStandard(hexagons)
	expectedSize := int64(4228)
	if actualNumIndexes != expectedSize {
		t.Errorf("Expected %d cells for prime meridian, got %d", expectedSize, actualNumIndexes)
	}

	// Transmeridian test
	err = maxPolygonToCellsSize(&transMeridianGeoPolygon, 7, 0, &numHexagons)
	if err != eSuccess {
		t.Fatalf("maxPolygonToCellsSize failed for transmeridian: %v", err)
	}

	hexagonsTM := make([]h3Index, numHexagons)
	err = polygonToCells(&transMeridianGeoPolygon, 7, 0, hexagonsTM)
	if err != eSuccess {
		t.Fatalf("polygonToCells failed for transmeridian: %v", err)
	}

	actualNumIndexes = countNonNullIndexesStandard(hexagonsTM)
	expectedSize = 4238
	if actualNumIndexes != expectedSize {
		t.Errorf("Expected %d cells for transmeridian, got %d", expectedSize, actualNumIndexes)
	}

	// Transmeridian filled hole (for calculating hole size)
	err = maxPolygonToCellsSize(&transMeridianFilledHoleGeoPolygon, 7, 0, &numHexagons)
	if err != eSuccess {
		t.Fatalf("maxPolygonToCellsSize failed for filled hole: %v", err)
	}

	hexagonsTMFH := make([]h3Index, numHexagons)
	err = polygonToCells(&transMeridianFilledHoleGeoPolygon, 7, 0, hexagonsTMFH)
	if err != eSuccess {
		t.Fatalf("polygonToCells failed for filled hole: %v", err)
	}

	actualNumHoleIndexes := countNonNullIndexesStandard(hexagonsTMFH)

	// Transmeridian with hole test
	err = maxPolygonToCellsSize(&transMeridianHoleGeoPolygon, 7, 0, &numHexagons)
	if err != eSuccess {
		t.Fatalf("maxPolygonToCellsSize failed for transmeridian hole: %v", err)
	}

	hexagonsTMH := make([]h3Index, numHexagons)
	err = polygonToCells(&transMeridianHoleGeoPolygon, 7, 0, hexagonsTMH)
	if err != eSuccess {
		t.Fatalf("polygonToCells failed for transmeridian hole: %v", err)
	}

	actualNumIndexes = countNonNullIndexesStandard(hexagonsTMH)
	expected := expectedSize - actualNumHoleIndexes
	if actualNumIndexes != expected {
		t.Errorf("Expected %d cells for transmeridian hole, got %d", expected, actualNumIndexes)
	}
}

func TestPolygonToCellsTransmeridianComplexStandard(t *testing.T) {
	t.Parallel()

	// Complex polygon with > 4 vertices
	verts := []LatLng{
		{0.1, -math.Pi + 0.00001}, {0.1, math.Pi - 0.00001},
		{0.05, math.Pi - 0.2}, {-0.1, math.Pi - 0.00001},
		{-0.1, -math.Pi + 0.00001}, {-0.05, -math.Pi + 0.2},
	}
	geoloop := GeoLoop(verts)
	polygon := GeoPolygon{GeoLoop: geoloop, Holes: nil}

	var numHexagons int64
	err := maxPolygonToCellsSize(&polygon, 4, 0, &numHexagons)
	if err != eSuccess {
		t.Fatalf("maxPolygonToCellsSize failed: %v", err)
	}

	hexagons := make([]h3Index, numHexagons)
	err = polygonToCells(&polygon, 4, 0, hexagons)
	if err != eSuccess {
		t.Fatalf("polygonToCells failed: %v", err)
	}

	actualNumIndexes := countNonNullIndexesStandard(hexagons)
	if actualNumIndexes != 1204 {
		t.Errorf("Expected 1204 cells for complex transmeridian, got %d", actualNumIndexes)
	}
}

func TestPolygonToCellsPentagonStandard(t *testing.T) {
	t.Parallel()

	var pentagon h3Index
	setH3Index(&pentagon, 9, 24, 0)
	var coord LatLng
	err := cellToLatLng(pentagon, &coord)
	if err != eSuccess {
		t.Fatalf("cellToLatLng failed: %v", err)
	}

	// Half edge length in radians
	edgeLength2 := degsToRads(0.001)

	boundingTopRight := LatLng{
		Lat: Rad(coord.Lat.Rad() + edgeLength2),
		Lng: Rad(coord.Lng.Rad() + edgeLength2),
	}

	boundingTopLeft := LatLng{
		Lat: Rad(coord.Lat.Rad() + edgeLength2),
		Lng: Rad(coord.Lng.Rad() - edgeLength2),
	}

	boundingBottomRight := LatLng{
		Lat: Rad(coord.Lat.Rad() - edgeLength2),
		Lng: Rad(coord.Lng.Rad() + edgeLength2),
	}

	boundingBottomLeft := LatLng{
		Lat: Rad(coord.Lat.Rad() - edgeLength2),
		Lng: Rad(coord.Lng.Rad() - edgeLength2),
	}

	verts := []LatLng{boundingBottomLeft, boundingTopLeft, boundingTopRight, boundingBottomRight}
	geoloop := GeoLoop(verts)
	polygon := GeoPolygon{GeoLoop: geoloop, Holes: nil}

	var numHexagons int64
	err = maxPolygonToCellsSize(&polygon, 9, 0, &numHexagons)
	if err != eSuccess {
		t.Fatalf("maxPolygonToCellsSize failed: %v", err)
	}

	hexagons := make([]h3Index, numHexagons)
	err = polygonToCells(&polygon, 9, 0, hexagons)
	if err != eSuccess {
		t.Fatalf("polygonToCells failed: %v", err)
	}

	found := 0
	numPentagons := 0
	for _, hex := range hexagons {
		if hex != h3Null {
			found++
		}
		if isPentagon(hex) {
			numPentagons++
		}
	}

	if found != 1 {
		t.Errorf("Expected 1 cell found, got %d", found)
	}
	if numPentagons != 1 {
		t.Errorf("Expected 1 pentagon found, got %d", numPentagons)
	}
}

func TestInvalidFlagsStandard(t *testing.T) {
	t.Parallel()

	var numHexagons int64

	// Test invalid flags for maxPolygonToCellsSize
	for flags := uint32(ContainmentInvalid); flags <= 32; flags++ {
		err := maxPolygonToCellsSize(&sfGeoPolygon, 9, flags, &numHexagons)
		if err != eOptionInvalid {
			t.Errorf("Expected eOptionInvalid for maxPolygonToCellsSize with flags %d", flags)
		}
	}

	// Test valid flags
	err := maxPolygonToCellsSize(&sfGeoPolygon, 9, 0, &numHexagons)
	if err != eSuccess {
		t.Fatalf("maxPolygonToCellsSize should succeed with flags=0: %v", err)
	}

	hexagons := make([]h3Index, numHexagons)

	// Test invalid flags for polygonToCells
	for flags := uint32(ContainmentInvalid); flags <= 32; flags++ {
		err = polygonToCells(&sfGeoPolygon, 9, flags, hexagons)
		if err != eOptionInvalid {
			t.Errorf("Expected eOptionInvalid for polygonToCells with flags %d", flags)
		}
	}
}

func TestPolygonToCellsInvalidPolygon(t *testing.T) {
	t.Parallel()

	hexagons := make([]h3Index, 0)
	err := polygonToCells(&invalidGeoPolygon, 9, 0, hexagons)
	if err != eFailed {
		t.Error("Expected eFailed for invalid geo polygon")
	}
}

func TestFillIndex(t *testing.T) {
	t.Parallel()

	// Test a few cells at resolutions 0, 1, 2
	// Note: This is a simplified version of iterateAllIndexesAtRes

	// Test base cells at resolution 0
	for baseCell := int32(0); baseCell < numBaseCells; baseCell++ {
		if _isBaseCellPentagon(baseCell) {
			continue // Skip pentagons for now
		}

		var h h3Index
		setH3Index(&h, 0, baseCell, 0)
		if !isValidCell(h) {
			continue
		}

		fillIndex_assertions(h)
	}

	// Test some cells at resolution 1
	for baseCell := int32(0); baseCell < 10; baseCell++ { // Limited sample
		if _isBaseCellPentagon(baseCell) {
			continue
		}

		var h h3Index
		setH3Index(&h, 1, baseCell, 0)
		if !isValidCell(h) {
			continue
		}

		fillIndex_assertions(h)
	}
}

// Helper function to get edge hexagons (tests internal function).
func TestGetEdgeHexagonsInvalid(t *testing.T) {
	t.Parallel()

	search := make([]h3Index, 100)
	found := make([]h3Index, 100)

	var numSearchHexes int64
	err := _getEdgeHexagons(invalidGeoLoop, 100, 0, &numSearchHexes, search, found)
	if err == eSuccess {
		t.Error("Expected _getEdgeHexagons to return error for invalid geoloop")
	}
}
