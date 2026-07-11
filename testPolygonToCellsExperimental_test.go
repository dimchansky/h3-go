// Tests ported from testPolygonToCellsExperimental.c
package h3

import (
	"math"
	"testing"
)

// Test fixtures - ported from C.
var (
	sfVerts = []LatLng{
		{Lat: 0.659966917655, Lng: -2.1364398519396},
		{Lat: 0.6595011102219, Lng: -2.1359434279405},
		{Lat: 0.6583348114025, Lng: -2.1354884206045},
		{Lat: 0.6581220034068, Lng: -2.1382437718946},
		{Lat: 0.6594479998527, Lng: -2.1384597563896},
		{Lat: 0.6599990002976, Lng: -2.1376771158464},
	}
	sfGeoLoop    = sfVerts
	sfGeoPolygon = GeoPolygon{GeoLoop: sfGeoLoop, Holes: nil}

	holeVerts = []LatLng{
		{Lat: 0.6595072188743, Lng: -2.1371053983433},
		{Lat: 0.6591482046471, Lng: -2.1373141048153},
		{Lat: 0.6592295020837, Lng: -2.1365222838402},
	}
	holeGeoLoop    = holeVerts
	holeGeoPolygon = GeoPolygon{GeoLoop: sfGeoLoop, Holes: []GeoLoop{holeGeoLoop}}

	emptyVerts = []LatLng{
		{Lat: 0.659966917655, Lng: -2.1364398519394},
		{Lat: 0.659966917656, Lng: -2.1364398519395},
		{Lat: 0.659966917657, Lng: -2.1364398519396},
	}
	emptyGeoLoop    = emptyVerts
	emptyGeoPolygon = GeoPolygon{GeoLoop: emptyGeoLoop, Holes: nil}

	invalidVerts = []LatLng{
		{Lat: Rad(math.Inf(1)), Lng: Rad(math.Inf(1))},
		{Lat: Rad(math.Inf(-1)), Lng: Rad(math.Inf(-1))},
	}
	invalidGeoLoop    = invalidVerts
	invalidGeoPolygon = GeoPolygon{GeoLoop: invalidGeoLoop, Holes: nil}

	outOfBoundsVert = []LatLng{
		{Lat: Rad(-2000), Lng: Rad(-2000)},
	}
	outOfBoundsVertGeoLoop    = outOfBoundsVert
	outOfBoundsVertGeoPolygon = GeoPolygon{GeoLoop: outOfBoundsVertGeoLoop, Holes: nil}

	invalid2Verts = []LatLng{
		{Lat: Rad(math.NaN()), Lng: Rad(math.NaN())},
		{Lat: Rad(math.NaN()), Lng: Rad(math.NaN())},
	}
	invalid2GeoLoop    = invalid2Verts
	invalid2GeoPolygon = GeoPolygon{GeoLoop: invalid2GeoLoop, Holes: nil}

	nullGeoLoop    = []LatLng{}
	nullGeoPolygon = GeoPolygon{GeoLoop: nullGeoLoop, Holes: nil}

	pointVerts = []LatLng{
		{Lat: 0.6595072188743, Lng: -2.1371053983433},
	}
	pointGeoLoop    = pointVerts
	pointGeoPolygon = GeoPolygon{GeoLoop: pointGeoLoop, Holes: nil}

	lineVerts = []LatLng{
		{Lat: 0.6595072188743, Lng: -2.1371053983433},
		{Lat: 0.6591482046471, Lng: -2.1373141048153},
	}
	lineGeoLoop    = lineVerts
	lineGeoPolygon = GeoPolygon{GeoLoop: lineGeoLoop, Holes: nil}

	nullHoleGeoPolygon  = GeoPolygon{GeoLoop: sfGeoLoop, Holes: []GeoLoop{nullGeoLoop}}
	pointHoleGeoPolygon = GeoPolygon{GeoLoop: sfGeoLoop, Holes: []GeoLoop{pointGeoLoop}}
	lineHoleGeoPolygon  = GeoPolygon{GeoLoop: sfGeoLoop, Holes: []GeoLoop{lineGeoLoop}}
)

// isTransmeridianCell returns true if the cell crosses the meridian.
func isTransmeridianCell(h h3Index) bool {
	var boundary CellBoundary
	err := cellToBoundary(h, &boundary)
	if err != eSuccess {
		return false
	}

	minLng := math.Pi
	maxLng := -math.Pi
	for i := int32(0); i < boundary.NumVerts; i++ {
		if boundary.Verts[i].Lng.Rad() < minLng {
			minLng = boundary.Verts[i].Lng.Rad()
		}
		if boundary.Verts[i].Lng.Rad() > maxLng {
			maxLng = boundary.Verts[i].Lng.Rad()
		}
	}

	return maxLng-minLng > math.Pi-(math.Pi/4)
}

// countNonNullIndexesWithSize counts non-null H3 indexes in a slice up to size.
func countNonNullIndexesWithSize(indexes []h3Index, size int64) int64 {
	var count int64
	for i := int64(0); i < size && i < int64(len(indexes)); i++ {
		if indexes[i] != h3Null {
			count++
		}
	}
	return count
}

func TestPolygonToCells_ZeroSize(t *testing.T) {
	t.Parallel()
	err := polygonToCellsExperimental(&sfGeoPolygon, 9, uint32(ContainmentCenter), 0, nil)
	if err != eMemoryBounds {
		t.Errorf("Expected eMemoryBounds for zero size, got %v", err)
	}
}

func TestPolygonToCells_CenterContainment(t *testing.T) {
	t.Parallel()
	numHexagons, err := maxPolygonToCellsSizeExperimental(&sfGeoPolygon, 9, uint32(ContainmentCenter))
	if err != eSuccess {
		t.Fatalf("maxPolygonToCellsSizeExperimental failed: %v", err)
	}

	hexagons := make([]h3Index, numHexagons)
	err = polygonToCellsExperimental(&sfGeoPolygon, 9, uint32(ContainmentCenter), numHexagons, hexagons)
	if err != eSuccess {
		t.Fatalf("polygonToCellsExperimental failed: %v", err)
	}

	actualNumIndexes := countNonNullIndexesWithSize(hexagons, numHexagons)
	expectedCount := int64(1253)
	if actualNumIndexes != expectedCount {
		t.Errorf("Expected %d indexes, got %d", expectedCount, actualNumIndexes)
	}
}

func TestPolygonToCells_FullContainment(t *testing.T) {
	t.Parallel()
	numHexagons, err := maxPolygonToCellsSizeExperimental(&sfGeoPolygon, 9, uint32(ContainmentFull))
	if err != eSuccess {
		t.Fatalf("maxPolygonToCellsSizeExperimental failed: %v", err)
	}

	hexagons := make([]h3Index, numHexagons)
	err = polygonToCellsExperimental(&sfGeoPolygon, 9, uint32(ContainmentFull), numHexagons, hexagons)
	if err != eSuccess {
		t.Fatalf("polygonToCellsExperimental failed: %v", err)
	}

	actualNumIndexes := countNonNullIndexesWithSize(hexagons, numHexagons)
	expectedCount := int64(1175)
	if actualNumIndexes != expectedCount {
		t.Errorf("Expected %d indexes, got %d (full containment mode)", expectedCount, actualNumIndexes)
	}
}

func TestPolygonToCells_Overlapping(t *testing.T) {
	t.Parallel()
	numHexagons, err := maxPolygonToCellsSizeExperimental(&sfGeoPolygon, 9, uint32(ContainmentOverlapping))
	if err != eSuccess {
		t.Fatalf("maxPolygonToCellsSizeExperimental failed: %v", err)
	}

	hexagons := make([]h3Index, numHexagons)
	err = polygonToCellsExperimental(&sfGeoPolygon, 9, uint32(ContainmentOverlapping), numHexagons, hexagons)
	if err != eSuccess {
		t.Fatalf("polygonToCellsExperimental failed: %v", err)
	}

	actualNumIndexes := countNonNullIndexesWithSize(hexagons, numHexagons)
	expectedCount := int64(1334)
	if actualNumIndexes != expectedCount {
		t.Errorf("Expected %d indexes, got %d (overlapping mode)", expectedCount, actualNumIndexes)
	}
}

func TestPolygonToCells_OverlappingBBox(t *testing.T) {
	t.Parallel()
	numHexagons, err := maxPolygonToCellsSizeExperimental(&sfGeoPolygon, 9, uint32(ContainmentOverlappingBBox))
	if err != eSuccess {
		t.Fatalf("maxPolygonToCellsSizeExperimental failed: %v", err)
	}

	hexagons := make([]h3Index, numHexagons)
	err = polygonToCellsExperimental(&sfGeoPolygon, 9, uint32(ContainmentOverlappingBBox), numHexagons, hexagons)
	if err != eSuccess {
		t.Fatalf("polygonToCellsExperimental failed: %v", err)
	}

	actualNumIndexes := countNonNullIndexesWithSize(hexagons, numHexagons)
	expectedCount := int64(1416)
	if actualNumIndexes != expectedCount {
		t.Errorf("Expected %d indexes, got %d (overlapping bbox mode)", expectedCount, actualNumIndexes)
	}
}

func TestPolygonToCellsHole_CenterContainment(t *testing.T) {
	t.Parallel()
	numHexagons, err := maxPolygonToCellsSizeExperimental(&holeGeoPolygon, 9, uint32(ContainmentCenter))
	if err != eSuccess {
		t.Fatalf("maxPolygonToCellsSizeExperimental failed: %v", err)
	}

	hexagons := make([]h3Index, numHexagons)
	err = polygonToCellsExperimental(&holeGeoPolygon, 9, uint32(ContainmentCenter), numHexagons, hexagons)
	if err != eSuccess {
		t.Fatalf("polygonToCellsExperimental failed: %v", err)
	}

	actualNumIndexes := countNonNullIndexesWithSize(hexagons, numHexagons)
	expectedCount := int64(1214)
	if actualNumIndexes != expectedCount {
		t.Errorf("Expected %d indexes, got %d (hole)", expectedCount, actualNumIndexes)
	}
}

func TestPolygonToCellsHole_FullContainment(t *testing.T) {
	t.Parallel()
	numHexagons, err := maxPolygonToCellsSizeExperimental(&holeGeoPolygon, 9, uint32(ContainmentFull))
	if err != eSuccess {
		t.Fatalf("maxPolygonToCellsSizeExperimental failed: %v", err)
	}

	hexagons := make([]h3Index, numHexagons)
	err = polygonToCellsExperimental(&holeGeoPolygon, 9, uint32(ContainmentFull), numHexagons, hexagons)
	if err != eSuccess {
		t.Fatalf("polygonToCellsExperimental failed: %v", err)
	}

	actualNumIndexes := countNonNullIndexesWithSize(hexagons, numHexagons)
	expectedCount := int64(1118)
	if actualNumIndexes != expectedCount {
		t.Errorf("Expected %d indexes, got %d (hole, full containment mode)", expectedCount, actualNumIndexes)
	}
}

func TestPolygonToCellsHole_Overlapping(t *testing.T) {
	t.Parallel()
	numHexagons, err := maxPolygonToCellsSizeExperimental(&holeGeoPolygon, 9, uint32(ContainmentOverlapping))
	if err != eSuccess {
		t.Fatalf("maxPolygonToCellsSizeExperimental failed: %v", err)
	}

	hexagons := make([]h3Index, numHexagons)
	err = polygonToCellsExperimental(&holeGeoPolygon, 9, uint32(ContainmentOverlapping), numHexagons, hexagons)
	if err != eSuccess {
		t.Fatalf("polygonToCellsExperimental failed: %v", err)
	}

	actualNumIndexes := countNonNullIndexesWithSize(hexagons, numHexagons)
	expectedCount := int64(1311)
	if actualNumIndexes != expectedCount {
		t.Errorf("Expected %d indexes, got %d (hole, overlapping mode)", expectedCount, actualNumIndexes)
	}
}

func TestPolygonToCellsHole_OverlappingBBox(t *testing.T) {
	t.Parallel()
	numHexagons, err := maxPolygonToCellsSizeExperimental(&holeGeoPolygon, 9, uint32(ContainmentOverlappingBBox))
	if err != eSuccess {
		t.Fatalf("maxPolygonToCellsSizeExperimental failed: %v", err)
	}

	hexagons := make([]h3Index, numHexagons)
	err = polygonToCellsExperimental(&holeGeoPolygon, 9, uint32(ContainmentOverlappingBBox), numHexagons, hexagons)
	if err != eSuccess {
		t.Fatalf("polygonToCellsExperimental failed: %v", err)
	}

	actualNumIndexes := countNonNullIndexesWithSize(hexagons, numHexagons)
	expectedCount := int64(1403)
	if actualNumIndexes != expectedCount {
		t.Errorf("Expected %d indexes, got %d (hole, overlapping bbox mode)", expectedCount, actualNumIndexes)
	}
}

func TestPolygonToCellsHoleParentIssue(t *testing.T) {
	t.Parallel()
	// This checks a specific issue where the bounding box of the parent
	// cell fully contains the hole.
	outer := []LatLng{
		{Lat: 0.7774570821346158, Lng: 0.19441847890170674},
		{Lat: 0.7528853613617879, Lng: 0.19441847890170674},
		{Lat: 0.7528853613617879, Lng: 0.23497118026107888},
		{Lat: 0.7774570821346158, Lng: 0.23497118026107888},
	}
	sanMarino := []LatLng{
		{Lat: 0.7662242554877188, Lng: 0.21790879024779208},
		{Lat: 0.7660964275733029, Lng: 0.21688101821117023},
		{Lat: 0.7668029019479251, Lng: 0.21636628570817204},
		{Lat: 0.7676380769015895, Lng: 0.21713838446266925},
		{Lat: 0.7677659048160054, Lng: 0.21823092566783267},
		{Lat: 0.7671241996099247, Lng: 0.2184218123281233},
		{Lat: 0.7662242554877188, Lng: 0.21790879024779208},
	}
	polygon := GeoPolygon{
		GeoLoop: outer,
		Holes:   []GeoLoop{sanMarino},
	}

	numHexagons, err := maxPolygonToCellsSizeExperimental(&polygon, 6, uint32(ContainmentCenter))
	if err != eSuccess {
		t.Fatalf("maxPolygonToCellsSizeExperimental failed: %v", err)
	}

	hexagons := make([]h3Index, numHexagons)
	err = polygonToCellsExperimental(&polygon, 6, uint32(ContainmentCenter), numHexagons, hexagons)
	if err != eSuccess {
		t.Fatalf("polygonToCellsExperimental failed: %v", err)
	}

	// This is the cell inside San Marino (i.e. inside the hole)
	holeCell := h3Index(0x861ea3cefffffff)

	found := false
	for i := int64(0); i < numHexagons; i++ {
		if hexagons[i] == holeCell {
			found = true
			break
		}
	}

	if found {
		t.Error("Should not include cell in hole")
	}
}

func TestPolygonToCellsEmpty(t *testing.T) {
	t.Parallel()
	numHexagons, err := maxPolygonToCellsSizeExperimental(&emptyGeoPolygon, 9, uint32(ContainmentCenter))
	if err != eSuccess {
		t.Fatalf("maxPolygonToCellsSizeExperimental failed: %v", err)
	}

	hexagons := make([]h3Index, numHexagons)
	err = polygonToCellsExperimental(&emptyGeoPolygon, 9, uint32(ContainmentCenter), numHexagons, hexagons)
	if err != eSuccess {
		t.Fatalf("polygonToCellsExperimental failed: %v", err)
	}

	actualNumIndexes := countNonNullIndexesWithSize(hexagons, numHexagons)
	expectedCount := int64(0)
	if actualNumIndexes != expectedCount {
		t.Errorf("Expected %d indexes, got %d (empty)", expectedCount, actualNumIndexes)
	}
}

func TestPolygonToCellsContainsPolygon(t *testing.T) {
	t.Parallel()
	numHexagons, err := maxPolygonToCellsSizeExperimental(&sfGeoPolygon, 4, uint32(ContainmentCenter))
	if err != eSuccess {
		t.Fatalf("maxPolygonToCellsSizeExperimental failed: %v", err)
	}

	hexagons := make([]h3Index, numHexagons)
	err = polygonToCellsExperimental(&sfGeoPolygon, 4, uint32(ContainmentCenter), numHexagons, hexagons)
	if err != eSuccess {
		t.Fatalf("polygonToCellsExperimental failed: %v", err)
	}

	actualNumIndexes := countNonNullIndexesWithSize(hexagons, numHexagons)
	expectedCount := int64(0)
	if actualNumIndexes != expectedCount {
		t.Errorf("Expected %d indexes, got %d", expectedCount, actualNumIndexes)
	}
}

func TestPolygonToCellsContainsPolygon_CenterContainment(t *testing.T) {
	t.Parallel()
	// Contains the center point of a res 4 cell
	centerVerts := []LatLng{
		{Lat: 0.6595645, Lng: -2.1353315},
		{Lat: 0.6595645, Lng: -2.1353314},
		{Lat: 0.6595644, Lng: -2.1353314},
		{Lat: 0.6595644, Lng: -2.1353314265},
	}
	centerGeoPolygon := GeoPolygon{GeoLoop: centerVerts, Holes: nil}

	numHexagons, err := maxPolygonToCellsSizeExperimental(&centerGeoPolygon, 4, uint32(ContainmentCenter))
	if err != eSuccess {
		t.Fatalf("maxPolygonToCellsSizeExperimental failed: %v", err)
	}

	hexagons := make([]h3Index, numHexagons)
	err = polygonToCellsExperimental(&centerGeoPolygon, 4, uint32(ContainmentCenter), numHexagons, hexagons)
	if err != eSuccess {
		t.Fatalf("polygonToCellsExperimental failed: %v", err)
	}

	actualNumIndexes := countNonNullIndexesWithSize(hexagons, numHexagons)
	expectedCount := int64(1)
	if actualNumIndexes != expectedCount {
		t.Errorf("Expected %d indexes, got %d", expectedCount, actualNumIndexes)
	}

	expectedHexagon := h3Index(0x8428309ffffffff)
	if hexagons[0] != expectedHexagon {
		t.Errorf("Expected hexagon %x, got %x", expectedHexagon, hexagons[0])
	}
}

func TestPolygonToCellsContainsPolygon_FullContainment(t *testing.T) {
	t.Parallel()
	numHexagons, err := maxPolygonToCellsSizeExperimental(&sfGeoPolygon, 4, uint32(ContainmentFull))
	if err != eSuccess {
		t.Fatalf("maxPolygonToCellsSizeExperimental failed: %v", err)
	}

	hexagons := make([]h3Index, numHexagons)
	err = polygonToCellsExperimental(&sfGeoPolygon, 4, uint32(ContainmentFull), numHexagons, hexagons)
	if err != eSuccess {
		t.Fatalf("polygonToCellsExperimental failed: %v", err)
	}

	actualNumIndexes := countNonNullIndexesWithSize(hexagons, numHexagons)
	expectedCount := int64(0)
	if actualNumIndexes != expectedCount {
		t.Errorf("Expected %d indexes, got %d (full containment mode)", expectedCount, actualNumIndexes)
	}
}

func TestPolygonToCellsContainsPolygon_Overlapping(t *testing.T) {
	t.Parallel()
	numHexagons, err := maxPolygonToCellsSizeExperimental(&sfGeoPolygon, 4, uint32(ContainmentOverlapping))
	if err != eSuccess {
		t.Fatalf("maxPolygonToCellsSizeExperimental failed: %v", err)
	}

	hexagons := make([]h3Index, numHexagons)
	err = polygonToCellsExperimental(&sfGeoPolygon, 4, uint32(ContainmentOverlapping), numHexagons, hexagons)
	if err != eSuccess {
		t.Fatalf("polygonToCellsExperimental failed: %v", err)
	}

	actualNumIndexes := countNonNullIndexesWithSize(hexagons, numHexagons)
	expectedCount := int64(1)
	if actualNumIndexes != expectedCount {
		t.Errorf("Expected %d indexes, got %d (overlapping mode)", expectedCount, actualNumIndexes)
	}

	expectedHexagon := h3Index(0x8428309ffffffff)
	if hexagons[0] != expectedHexagon {
		t.Errorf("Expected hexagon %x, got %x", expectedHexagon, hexagons[0])
	}
}

func TestPolygonToCellsContainsPolygon_OverlappingBBox(t *testing.T) {
	t.Parallel()
	numHexagons, err := maxPolygonToCellsSizeExperimental(&sfGeoPolygon, 4, uint32(ContainmentOverlappingBBox))
	if err != eSuccess {
		t.Fatalf("maxPolygonToCellsSizeExperimental failed: %v", err)
	}

	hexagons := make([]h3Index, numHexagons)
	err = polygonToCellsExperimental(&sfGeoPolygon, 4, uint32(ContainmentOverlappingBBox), numHexagons, hexagons)
	if err != eSuccess {
		t.Fatalf("polygonToCellsExperimental failed: %v", err)
	}

	actualNumIndexes := countNonNullIndexesWithSize(hexagons, numHexagons)
	expectedCount := int64(5)
	if actualNumIndexes != expectedCount {
		t.Errorf("Expected %d indexes, got %d (overlapping bbox mode)", expectedCount, actualNumIndexes)
	}

	expectedHexagon := h3Index(0x8428309ffffffff)
	if hexagons[0] != expectedHexagon {
		t.Errorf("Expected hexagon %x, got %x", expectedHexagon, hexagons[0])
	}
}

func TestPolygonToCellsExact(t *testing.T) {
	t.Parallel()
	somewhere := LatLng{Lat: 1, Lng: 2}
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

	verts := make([]LatLng, boundary.NumVerts+1)
	for i := int32(0); i < boundary.NumVerts; i++ {
		verts[i] = boundary.Verts[i]
	}
	verts[boundary.NumVerts] = boundary.Verts[0]

	someHexagon := GeoPolygon{
		GeoLoop: verts,
		Holes:   nil,
	}

	numHexagons, err := maxPolygonToCellsSizeExperimental(&someHexagon, 9, uint32(ContainmentCenter))
	if err != eSuccess {
		t.Fatalf("maxPolygonToCellsSizeExperimental failed: %v", err)
	}

	hexagons := make([]h3Index, numHexagons)

	// Test center containment
	err = polygonToCellsExperimental(&someHexagon, 9, uint32(ContainmentCenter), numHexagons, hexagons)
	if err != eSuccess {
		t.Fatalf("polygonToCellsExperimental failed: %v", err)
	}
	actualNumIndexes := countNonNullIndexesWithSize(hexagons, numHexagons)
	if actualNumIndexes != 1 {
		t.Errorf("Expected 1 index for center containment, got %d", actualNumIndexes)
	}

	// Test full containment
	err = polygonToCellsExperimental(&someHexagon, 9, uint32(ContainmentFull), numHexagons, hexagons)
	if err != eSuccess {
		t.Fatalf("polygonToCellsExperimental failed: %v", err)
	}
	actualNumIndexes = countNonNullIndexesWithSize(hexagons, numHexagons)
	if actualNumIndexes != 1 {
		t.Errorf("Expected 1 index for full containment, got %d", actualNumIndexes)
	}

	// Test overlapping bbox containment
	err = polygonToCellsExperimental(&someHexagon, 9, uint32(ContainmentOverlappingBBox), numHexagons, hexagons)
	if err != eSuccess {
		t.Fatalf("polygonToCellsExperimental failed: %v", err)
	}
	actualNumIndexes = countNonNullIndexesWithSize(hexagons, numHexagons)
	// Overlapping bbox is very rough, so we get a couple of overlaps from
	// non-neighboring cells
	if actualNumIndexes != 9 {
		t.Errorf("Expected 9 indexes for overlapping bbox containment, got %d", actualNumIndexes)
	}
}

func TestPolygonToCellsTransmeridian(t *testing.T) {
	t.Parallel()
	primeMeridianVerts := []LatLng{
		{Lat: 0.01, Lng: 0.01},
		{Lat: 0.01, Lng: -0.01},
		{Lat: -0.01, Lng: -0.01},
		{Lat: -0.01, Lng: 0.01},
	}
	primeMeridianGeoPolygon := GeoPolygon{GeoLoop: primeMeridianVerts, Holes: nil}

	transMeridianVerts := []LatLng{
		{Lat: 0.01, Lng: -math.Pi + 0.01},
		{Lat: 0.01, Lng: math.Pi - 0.01},
		{Lat: -0.01, Lng: math.Pi - 0.01},
		{Lat: -0.01, Lng: -math.Pi + 0.01},
	}
	transMeridianGeoPolygon := GeoPolygon{GeoLoop: transMeridianVerts, Holes: nil}

	transMeridianHoleVerts := []LatLng{
		{Lat: 0.005, Lng: -math.Pi + 0.005},
		{Lat: 0.005, Lng: math.Pi - 0.005},
		{Lat: -0.005, Lng: math.Pi - 0.005},
		{Lat: -0.005, Lng: -math.Pi + 0.005},
	}
	transMeridianHoleGeoPolygon := GeoPolygon{
		GeoLoop: transMeridianVerts,
		Holes:   []GeoLoop{transMeridianHoleVerts},
	}
	transMeridianFilledHoleGeoPolygon := GeoPolygon{
		GeoLoop: transMeridianHoleVerts,
		Holes:   nil,
	}

	// Prime meridian case
	expectedSize := int64(4228)
	numHexagons, err := maxPolygonToCellsSizeExperimental(&primeMeridianGeoPolygon, 7, uint32(ContainmentCenter))
	if err != eSuccess {
		t.Fatalf("maxPolygonToCellsSizeExperimental failed: %v", err)
	}

	hexagons := make([]h3Index, numHexagons)
	err = polygonToCellsExperimental(&primeMeridianGeoPolygon, 7, uint32(ContainmentCenter), numHexagons, hexagons)
	if err != eSuccess {
		t.Fatalf("polygonToCellsExperimental failed: %v", err)
	}

	actualNumIndexes := countNonNullIndexesWithSize(hexagons, numHexagons)
	if actualNumIndexes != expectedSize {
		t.Errorf("Expected %d indexes, got %d (prime meridian)", expectedSize, actualNumIndexes)
	}

	// Transmeridian case
	// This doesn't exactly match the prime meridian count because of slight
	// differences in hex size and grid offset between the two cases
	expectedSize = int64(4238)
	numHexagons, err = maxPolygonToCellsSizeExperimental(&transMeridianGeoPolygon, 7, uint32(ContainmentCenter))
	if err != eSuccess {
		t.Fatalf("maxPolygonToCellsSizeExperimental failed: %v", err)
	}

	hexagonsTM := make([]h3Index, numHexagons)
	err = polygonToCellsExperimental(&transMeridianGeoPolygon, 7, uint32(ContainmentCenter), numHexagons, hexagonsTM)
	if err != eSuccess {
		t.Fatalf("polygonToCellsExperimental failed: %v", err)
	}

	actualNumIndexes = countNonNullIndexesWithSize(hexagonsTM, numHexagons)
	if actualNumIndexes != expectedSize {
		t.Errorf("Expected %d indexes, got %d (transmeridian)", expectedSize, actualNumIndexes)
	}

	// Transmeridian filled hole case -- only needed for calculating hole size
	numHexagons, err = maxPolygonToCellsSizeExperimental(&transMeridianFilledHoleGeoPolygon, 7, uint32(ContainmentCenter))
	if err != eSuccess {
		t.Fatalf("maxPolygonToCellsSizeExperimental failed: %v", err)
	}

	hexagonsTMFH := make([]h3Index, numHexagons)
	err = polygonToCellsExperimental(&transMeridianFilledHoleGeoPolygon, 7, uint32(ContainmentCenter), numHexagons, hexagonsTMFH)
	if err != eSuccess {
		t.Fatalf("polygonToCellsExperimental failed: %v", err)
	}

	actualNumHoleIndexes := countNonNullIndexesWithSize(hexagonsTMFH, numHexagons)

	// Transmeridian hole case
	numHexagons, err = maxPolygonToCellsSizeExperimental(&transMeridianHoleGeoPolygon, 7, uint32(ContainmentCenter))
	if err != eSuccess {
		t.Fatalf("maxPolygonToCellsSizeExperimental failed: %v", err)
	}

	hexagonsTMH := make([]h3Index, numHexagons)
	err = polygonToCellsExperimental(&transMeridianHoleGeoPolygon, 7, uint32(ContainmentCenter), numHexagons, hexagonsTMH)
	if err != eSuccess {
		t.Fatalf("polygonToCellsExperimental failed: %v", err)
	}

	actualNumIndexes = countNonNullIndexesWithSize(hexagonsTMH, numHexagons)
	expectedCount := expectedSize - actualNumHoleIndexes
	if actualNumIndexes != expectedCount {
		t.Errorf("Expected %d indexes, got %d (transmeridian hole)", expectedCount, actualNumIndexes)
	}
}

func TestPolygonToCellsTransmeridianComplex(t *testing.T) {
	t.Parallel()
	// This polygon is "complex" in that it has > 4 vertices - this
	// tests for a bug that was taking the max and min longitude as
	// the bounds for transmeridian polygons
	verts := []LatLng{
		{Lat: 0.1, Lng: -math.Pi + 0.00001},
		{Lat: 0.1, Lng: math.Pi - 0.00001},
		{Lat: 0.05, Lng: math.Pi - 0.2},
		{Lat: -0.1, Lng: math.Pi - 0.00001},
		{Lat: -0.1, Lng: -math.Pi + 0.00001},
		{Lat: -0.05, Lng: -math.Pi + 0.2},
	}
	polygon := GeoPolygon{GeoLoop: verts, Holes: nil}

	numHexagons, err := maxPolygonToCellsSizeExperimental(&polygon, 4, uint32(ContainmentCenter))
	if err != eSuccess {
		t.Fatalf("maxPolygonToCellsSizeExperimental failed: %v", err)
	}

	hexagons := make([]h3Index, numHexagons)
	err = polygonToCellsExperimental(&polygon, 4, uint32(ContainmentCenter), numHexagons, hexagons)
	if err != eSuccess {
		t.Fatalf("polygonToCellsExperimental failed: %v", err)
	}

	actualNumIndexes := countNonNullIndexesWithSize(hexagons, numHexagons)
	expectedCount := int64(1204)
	if actualNumIndexes != expectedCount {
		t.Errorf("Expected %d indexes, got %d (complex transmeridian)", expectedCount, actualNumIndexes)
	}
}

func TestPolygonToCellsPentagon(t *testing.T) {
	t.Parallel()
	// Get a pentagon cell
	var pentagon h3Index
	setH3Index(&pentagon, 9, 24, 0)

	var coord LatLng
	err := cellToLatLng(pentagon, &coord)
	if err != eSuccess {
		t.Fatalf("cellToLatLng failed: %v", err)
	}

	// Length of half an edge of the polygon, in radians
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
	polygon := GeoPolygon{
		GeoLoop: verts,
		Holes:   nil,
	}

	numHexagons, err := maxPolygonToCellsSizeExperimental(&polygon, 9, uint32(ContainmentCenter))
	if err != eSuccess {
		t.Fatalf("maxPolygonToCellsSizeExperimental failed: %v", err)
	}

	hexagons := make([]h3Index, numHexagons)
	err = polygonToCellsExperimental(&polygon, 9, uint32(ContainmentCenter), numHexagons, hexagons)
	if err != eSuccess {
		t.Fatalf("polygonToCellsExperimental failed: %v", err)
	}

	found := 0
	numPentagons := 0
	for i := 0; i < int(numHexagons); i++ {
		if hexagons[i] != 0 {
			found++
		}
		if isPentagon(hexagons[i]) {
			numPentagons++
		}
	}

	if found != 1 {
		t.Errorf("Expected 1 index found, got %d", found)
	}
	if numPentagons != 1 {
		t.Errorf("Expected 1 pentagon found, got %d", numPentagons)
	}
}

func TestPolygonToCellsNullPolygon(t *testing.T) {
	t.Parallel()
	for res := int32(0); res <= maxH3Res; res++ {
		for flags := uint32(0); flags < uint32(ContainmentInvalid); flags++ {
			numHexagons, err := maxPolygonToCellsSizeExperimental(&nullGeoPolygon, res, flags)
			if err != eSuccess {
				t.Fatalf("maxPolygonToCellsSizeExperimental failed: %v", err)
			}
			if numHexagons != 0 {
				t.Errorf("Expected 0 estimated size, got %d", numHexagons)
			}

			var hexagons []h3Index
			if numHexagons > 0 {
				hexagons = make([]h3Index, numHexagons)
				err = polygonToCellsExperimental(&nullGeoPolygon, res, flags, numHexagons, hexagons)
				if err != eSuccess {
					t.Fatalf("polygonToCellsExperimental failed: %v", err)
				}
			}

			actualNumIndexes := countNonNullIndexesWithSize(hexagons, numHexagons)
			if actualNumIndexes != 0 {
				t.Errorf("Expected 0 polygonToCells size, got %d", actualNumIndexes)
			}
		}
	}
}

func TestPolygonToCellsPointPolygon_CenterContainment(t *testing.T) {
	t.Parallel()
	for res := int32(0); res <= maxH3Res; res++ {
		numHexagons, err := maxPolygonToCellsSizeExperimental(&pointGeoPolygon, res, uint32(ContainmentCenter))
		if err != eSuccess {
			t.Fatalf("maxPolygonToCellsSizeExperimental failed: %v", err)
		}
		if numHexagons < 1 || numHexagons > 5 {
			t.Errorf("Expected estimated size between 1 and 5, got %d", numHexagons)
		}

		hexagons := make([]h3Index, numHexagons)
		err = polygonToCellsExperimental(&pointGeoPolygon, res, uint32(ContainmentCenter), numHexagons, hexagons)
		if err != eSuccess {
			t.Fatalf("polygonToCellsExperimental failed: %v", err)
		}

		actualNumIndexes := countNonNullIndexesWithSize(hexagons, numHexagons)
		if actualNumIndexes != 0 {
			t.Errorf("Expected 0 polygonToCells size, got %d", actualNumIndexes)
		}
	}
}

func TestPolygonToCellsPointPolygon_FullContainment(t *testing.T) {
	t.Parallel()
	for res := int32(0); res <= maxH3Res; res++ {
		numHexagons, err := maxPolygonToCellsSizeExperimental(&pointGeoPolygon, res, uint32(ContainmentFull))
		if err != eSuccess {
			t.Fatalf("maxPolygonToCellsSizeExperimental failed: %v", err)
		}
		if numHexagons < 1 || numHexagons > 5 {
			t.Errorf("Expected estimated size between 1 and 5, got %d", numHexagons)
		}

		hexagons := make([]h3Index, numHexagons)
		err = polygonToCellsExperimental(&pointGeoPolygon, res, uint32(ContainmentFull), numHexagons, hexagons)
		if err != eSuccess {
			t.Fatalf("polygonToCellsExperimental failed: %v", err)
		}

		actualNumIndexes := countNonNullIndexesWithSize(hexagons, numHexagons)
		if actualNumIndexes != 0 {
			t.Errorf("Expected 0 polygonToCells size, got %d", actualNumIndexes)
		}
	}
}

func TestPolygonToCellsPointPolygon_Overlapping(t *testing.T) {
	t.Parallel()
	for res := int32(0); res <= maxH3Res; res++ {
		numHexagons, err := maxPolygonToCellsSizeExperimental(&pointGeoPolygon, res, uint32(ContainmentOverlapping))
		if err != eSuccess {
			t.Fatalf("maxPolygonToCellsSizeExperimental failed: %v", err)
		}
		if numHexagons < 1 || numHexagons > 5 {
			t.Errorf("Expected estimated size between 1 and 5, got %d", numHexagons)
		}

		hexagons := make([]h3Index, numHexagons)
		err = polygonToCellsExperimental(&pointGeoPolygon, res, uint32(ContainmentOverlapping), numHexagons, hexagons)
		if err != eSuccess {
			t.Fatalf("polygonToCellsExperimental failed: %v", err)
		}

		actualNumIndexes := countNonNullIndexesWithSize(hexagons, numHexagons)
		if actualNumIndexes != 1 {
			t.Errorf("Expected 1 polygonToCells size, got %d", actualNumIndexes)
		}
	}
}

func TestPolygonToCellsPointPolygon_OverlappingBBox(t *testing.T) {
	t.Parallel()
	for res := int32(0); res <= maxH3Res; res++ {
		numHexagons, err := maxPolygonToCellsSizeExperimental(&pointGeoPolygon, res, uint32(ContainmentOverlappingBBox))
		if err != eSuccess {
			t.Fatalf("maxPolygonToCellsSizeExperimental failed: %v", err)
		}
		if numHexagons < 1 || numHexagons > 5 {
			t.Errorf("Expected estimated size between 1 and 5, got %d", numHexagons)
		}

		hexagons := make([]h3Index, numHexagons)
		err = polygonToCellsExperimental(&pointGeoPolygon, res, uint32(ContainmentOverlappingBBox), numHexagons, hexagons)
		if err != eSuccess {
			t.Fatalf("polygonToCellsExperimental failed: %v", err)
		}

		actualNumIndexes := countNonNullIndexesWithSize(hexagons, numHexagons)
		if actualNumIndexes < 1 || actualNumIndexes > 5 {
			t.Errorf("Expected polygonToCells size between 1 and 5, got %d", actualNumIndexes)
		}
	}
}

func TestPolygonToCellsOutOfBoundsPolygon(t *testing.T) {
	t.Parallel()
	for res := int32(0); res <= maxH3Res; res++ {
		for flags := uint32(0); flags < uint32(ContainmentInvalid); flags++ {
			numHexagons, err := maxPolygonToCellsSizeExperimental(&outOfBoundsVertGeoPolygon, res, flags)
			if err != eSuccess {
				t.Fatalf("maxPolygonToCellsSizeExperimental failed: %v", err)
			}
			if numHexagons != 0 {
				t.Errorf("Expected 0 estimated size, got %d", numHexagons)
			}

			// Note: We're allocating more memory than the estimate to test
			// for out-of-bounds writes here
			hexagons := make([]h3Index, 10)
			err = polygonToCellsExperimental(&outOfBoundsVertGeoPolygon, res, flags, numHexagons, hexagons)
			if err != eSuccess {
				t.Fatalf("polygonToCellsExperimental failed: %v", err)
			}

			actualNumIndexes := countNonNullIndexesWithSize(hexagons, numHexagons)
			if actualNumIndexes != 0 {
				t.Errorf("Expected 0 polygonToCells size, got %d", actualNumIndexes)
			}
		}
	}
}

func TestPolygonToCellsLinePolygon_CenterContainment(t *testing.T) {
	t.Parallel()
	numHexagons, err := maxPolygonToCellsSizeExperimental(&lineGeoPolygon, 9, uint32(ContainmentCenter))
	if err != eSuccess {
		t.Fatalf("maxPolygonToCellsSizeExperimental failed: %v", err)
	}

	hexagons := make([]h3Index, numHexagons)
	err = polygonToCellsExperimental(&lineGeoPolygon, 9, uint32(ContainmentCenter), numHexagons, hexagons)
	if err != eSuccess {
		t.Fatalf("polygonToCellsExperimental failed: %v", err)
	}

	actualNumIndexes := countNonNullIndexesWithSize(hexagons, numHexagons)
	if actualNumIndexes != 0 {
		t.Errorf("Expected 0 polygonToCells size, got %d", actualNumIndexes)
	}
}

func TestPolygonToCellsLinePolygon_FullContainment(t *testing.T) {
	t.Parallel()
	numHexagons, err := maxPolygonToCellsSizeExperimental(&lineGeoPolygon, 9, uint32(ContainmentFull))
	if err != eSuccess {
		t.Fatalf("maxPolygonToCellsSizeExperimental failed: %v", err)
	}

	hexagons := make([]h3Index, numHexagons)
	err = polygonToCellsExperimental(&lineGeoPolygon, 9, uint32(ContainmentFull), numHexagons, hexagons)
	if err != eSuccess {
		t.Fatalf("polygonToCellsExperimental failed: %v", err)
	}

	actualNumIndexes := countNonNullIndexesWithSize(hexagons, numHexagons)
	if actualNumIndexes != 0 {
		t.Errorf("Expected 0 polygonToCells size, got %d", actualNumIndexes)
	}
}

func TestPolygonToCellsLinePolygon_Overlapping(t *testing.T) {
	t.Parallel()
	numHexagons, err := maxPolygonToCellsSizeExperimental(&lineGeoPolygon, 9, uint32(ContainmentOverlapping))
	if err != eSuccess {
		t.Fatalf("maxPolygonToCellsSizeExperimental failed: %v", err)
	}

	hexagons := make([]h3Index, numHexagons)
	err = polygonToCellsExperimental(&lineGeoPolygon, 9, uint32(ContainmentOverlapping), numHexagons, hexagons)
	if err != eSuccess {
		t.Fatalf("polygonToCellsExperimental failed: %v", err)
	}

	actualNumIndexes := countNonNullIndexesWithSize(hexagons, numHexagons)
	if actualNumIndexes != 9 {
		t.Errorf("Expected 9 polygonToCells size, got %d", actualNumIndexes)
	}
}

func TestPolygonToCellsLinePolygon_OverlappingBBox(t *testing.T) {
	t.Parallel()
	numHexagons, err := maxPolygonToCellsSizeExperimental(&lineGeoPolygon, 9, uint32(ContainmentOverlappingBBox))
	if err != eSuccess {
		t.Fatalf("maxPolygonToCellsSizeExperimental failed: %v", err)
	}

	hexagons := make([]h3Index, numHexagons)
	err = polygonToCellsExperimental(&lineGeoPolygon, 9, uint32(ContainmentOverlappingBBox), numHexagons, hexagons)
	if err != eSuccess {
		t.Fatalf("polygonToCellsExperimental failed: %v", err)
	}

	actualNumIndexes := countNonNullIndexesWithSize(hexagons, numHexagons)
	if actualNumIndexes != 21 {
		t.Errorf("Expected 21 polygonToCells size, got %d", actualNumIndexes)
	}
}

func TestPolygonToCellsNullHole_CenterContainment(t *testing.T) {
	t.Parallel()
	numHexagons, err := maxPolygonToCellsSizeExperimental(&nullHoleGeoPolygon, 9, uint32(ContainmentCenter))
	if err != eSuccess {
		t.Fatalf("maxPolygonToCellsSizeExperimental failed: %v", err)
	}

	hexagons := make([]h3Index, numHexagons)
	err = polygonToCellsExperimental(&nullHoleGeoPolygon, 9, uint32(ContainmentCenter), numHexagons, hexagons)
	if err != eSuccess {
		t.Fatalf("polygonToCellsExperimental failed: %v", err)
	}

	actualNumIndexes := countNonNullIndexesWithSize(hexagons, numHexagons)
	// Same as without the hole
	if actualNumIndexes != 1253 {
		t.Errorf("Expected 1253 polygonToCells size (null hole), got %d", actualNumIndexes)
	}
}

func TestPolygonToCellsNullHole_FullContainment(t *testing.T) {
	t.Parallel()
	numHexagons, err := maxPolygonToCellsSizeExperimental(&nullHoleGeoPolygon, 9, uint32(ContainmentFull))
	if err != eSuccess {
		t.Fatalf("maxPolygonToCellsSizeExperimental failed: %v", err)
	}

	hexagons := make([]h3Index, numHexagons)
	err = polygonToCellsExperimental(&nullHoleGeoPolygon, 9, uint32(ContainmentFull), numHexagons, hexagons)
	if err != eSuccess {
		t.Fatalf("polygonToCellsExperimental failed: %v", err)
	}

	actualNumIndexes := countNonNullIndexesWithSize(hexagons, numHexagons)
	// Same as without the hole
	if actualNumIndexes != 1175 {
		t.Errorf("Expected 1175 polygonToCells size (null hole), got %d", actualNumIndexes)
	}
}

func TestPolygonToCellsNullHole_Overlapping(t *testing.T) {
	t.Parallel()
	numHexagons, err := maxPolygonToCellsSizeExperimental(&nullHoleGeoPolygon, 9, uint32(ContainmentOverlapping))
	if err != eSuccess {
		t.Fatalf("maxPolygonToCellsSizeExperimental failed: %v", err)
	}

	hexagons := make([]h3Index, numHexagons)
	err = polygonToCellsExperimental(&nullHoleGeoPolygon, 9, uint32(ContainmentOverlapping), numHexagons, hexagons)
	if err != eSuccess {
		t.Fatalf("polygonToCellsExperimental failed: %v", err)
	}

	actualNumIndexes := countNonNullIndexesWithSize(hexagons, numHexagons)
	// Same as without the hole
	if actualNumIndexes != 1334 {
		t.Errorf("Expected 1334 polygonToCells size (null hole), got %d", actualNumIndexes)
	}
}

func TestPolygonToCellsNullHole_OverlappingBBox(t *testing.T) {
	t.Parallel()
	numHexagons, err := maxPolygonToCellsSizeExperimental(&nullHoleGeoPolygon, 9, uint32(ContainmentOverlappingBBox))
	if err != eSuccess {
		t.Fatalf("maxPolygonToCellsSizeExperimental failed: %v", err)
	}

	hexagons := make([]h3Index, numHexagons)
	err = polygonToCellsExperimental(&nullHoleGeoPolygon, 9, uint32(ContainmentOverlappingBBox), numHexagons, hexagons)
	if err != eSuccess {
		t.Fatalf("polygonToCellsExperimental failed: %v", err)
	}

	actualNumIndexes := countNonNullIndexesWithSize(hexagons, numHexagons)
	// Same as without the hole
	if actualNumIndexes != 1416 {
		t.Errorf("Expected 1416 polygonToCells size (null hole), got %d", actualNumIndexes)
	}
}

func TestPolygonToCellsPointHole_CenterContainment(t *testing.T) {
	t.Parallel()
	numHexagons, err := maxPolygonToCellsSizeExperimental(&pointHoleGeoPolygon, 9, uint32(ContainmentCenter))
	if err != eSuccess {
		t.Fatalf("maxPolygonToCellsSizeExperimental failed: %v", err)
	}

	hexagons := make([]h3Index, numHexagons)
	err = polygonToCellsExperimental(&pointHoleGeoPolygon, 9, uint32(ContainmentCenter), numHexagons, hexagons)
	if err != eSuccess {
		t.Fatalf("polygonToCellsExperimental failed: %v", err)
	}

	actualNumIndexes := countNonNullIndexesWithSize(hexagons, numHexagons)
	// Same as without the hole
	if actualNumIndexes != 1253 {
		t.Errorf("Expected 1253 polygonToCells size (point hole), got %d", actualNumIndexes)
	}
}

func TestPolygonToCellsPointHole_FullContainment(t *testing.T) {
	t.Parallel()
	numHexagons, err := maxPolygonToCellsSizeExperimental(&pointHoleGeoPolygon, 9, uint32(ContainmentFull))
	if err != eSuccess {
		t.Fatalf("maxPolygonToCellsSizeExperimental failed: %v", err)
	}

	hexagons := make([]h3Index, numHexagons)
	err = polygonToCellsExperimental(&pointHoleGeoPolygon, 9, uint32(ContainmentFull), numHexagons, hexagons)
	if err != eSuccess {
		t.Fatalf("polygonToCellsExperimental failed: %v", err)
	}

	actualNumIndexes := countNonNullIndexesWithSize(hexagons, numHexagons)
	// We expect that the cell containing the hole is not included
	if actualNumIndexes != 1175-1 {
		t.Errorf("Expected %d polygonToCells size (point hole), got %d", 1175-1, actualNumIndexes)
	}
}

func TestPolygonToCellsPointHole_Overlapping(t *testing.T) {
	t.Parallel()
	numHexagons, err := maxPolygonToCellsSizeExperimental(&pointHoleGeoPolygon, 9, uint32(ContainmentOverlapping))
	if err != eSuccess {
		t.Fatalf("maxPolygonToCellsSizeExperimental failed: %v", err)
	}

	hexagons := make([]h3Index, numHexagons)
	err = polygonToCellsExperimental(&pointHoleGeoPolygon, 9, uint32(ContainmentOverlapping), numHexagons, hexagons)
	if err != eSuccess {
		t.Fatalf("polygonToCellsExperimental failed: %v", err)
	}

	actualNumIndexes := countNonNullIndexesWithSize(hexagons, numHexagons)
	// Same as without the hole
	if actualNumIndexes != 1334 {
		t.Errorf("Expected 1334 polygonToCells size (point hole), got %d", actualNumIndexes)
	}
}

func TestPolygonToCellsPointHole_OverlappingBBox(t *testing.T) {
	t.Parallel()
	numHexagons, err := maxPolygonToCellsSizeExperimental(&pointHoleGeoPolygon, 9, uint32(ContainmentOverlappingBBox))
	if err != eSuccess {
		t.Fatalf("maxPolygonToCellsSizeExperimental failed: %v", err)
	}

	hexagons := make([]h3Index, numHexagons)
	err = polygonToCellsExperimental(&pointHoleGeoPolygon, 9, uint32(ContainmentOverlappingBBox), numHexagons, hexagons)
	if err != eSuccess {
		t.Fatalf("polygonToCellsExperimental failed: %v", err)
	}

	actualNumIndexes := countNonNullIndexesWithSize(hexagons, numHexagons)
	// Same as without the hole
	if actualNumIndexes != 1416 {
		t.Errorf("Expected 1416 polygonToCells size (point hole), got %d", actualNumIndexes)
	}
}

func TestPolygonToCellsLineHole_CenterContainment(t *testing.T) {
	t.Parallel()
	numHexagons, err := maxPolygonToCellsSizeExperimental(&lineHoleGeoPolygon, 9, uint32(ContainmentCenter))
	if err != eSuccess {
		t.Fatalf("maxPolygonToCellsSizeExperimental failed: %v", err)
	}

	hexagons := make([]h3Index, numHexagons)
	err = polygonToCellsExperimental(&lineHoleGeoPolygon, 9, uint32(ContainmentCenter), numHexagons, hexagons)
	if err != eSuccess {
		t.Fatalf("polygonToCellsExperimental failed: %v", err)
	}

	actualNumIndexes := countNonNullIndexesWithSize(hexagons, numHexagons)
	// Same as without the hole
	if actualNumIndexes != 1253 {
		t.Errorf("Expected 1253 polygonToCells size (line hole), got %d", actualNumIndexes)
	}
}

func TestPolygonToCellsLineHole_FullContainment(t *testing.T) {
	t.Parallel()
	numHexagons, err := maxPolygonToCellsSizeExperimental(&lineHoleGeoPolygon, 9, uint32(ContainmentFull))
	if err != eSuccess {
		t.Fatalf("maxPolygonToCellsSizeExperimental failed: %v", err)
	}

	hexagons := make([]h3Index, numHexagons)
	err = polygonToCellsExperimental(&lineHoleGeoPolygon, 9, uint32(ContainmentFull), numHexagons, hexagons)
	if err != eSuccess {
		t.Fatalf("polygonToCellsExperimental failed: %v", err)
	}

	actualNumIndexes := countNonNullIndexesWithSize(hexagons, numHexagons)
	// We expect that the cells intersecting the line are not included
	if actualNumIndexes != 1175-9 {
		t.Errorf("Expected %d polygonToCells size (line hole), got %d", 1175-9, actualNumIndexes)
	}
}

func TestPolygonToCellsLineHole_Overlapping(t *testing.T) {
	t.Parallel()
	numHexagons, err := maxPolygonToCellsSizeExperimental(&lineHoleGeoPolygon, 9, uint32(ContainmentOverlapping))
	if err != eSuccess {
		t.Fatalf("maxPolygonToCellsSizeExperimental failed: %v", err)
	}

	hexagons := make([]h3Index, numHexagons)
	err = polygonToCellsExperimental(&lineHoleGeoPolygon, 9, uint32(ContainmentOverlapping), numHexagons, hexagons)
	if err != eSuccess {
		t.Fatalf("polygonToCellsExperimental failed: %v", err)
	}

	actualNumIndexes := countNonNullIndexesWithSize(hexagons, numHexagons)
	// Same as without the hole
	if actualNumIndexes != 1334 {
		t.Errorf("Expected 1334 polygonToCells size (line hole), got %d", actualNumIndexes)
	}
}

func TestPolygonToCellsLineHole_OverlappingBBox(t *testing.T) {
	t.Parallel()
	numHexagons, err := maxPolygonToCellsSizeExperimental(&lineHoleGeoPolygon, 9, uint32(ContainmentOverlappingBBox))
	if err != eSuccess {
		t.Fatalf("maxPolygonToCellsSizeExperimental failed: %v", err)
	}

	hexagons := make([]h3Index, numHexagons)
	err = polygonToCellsExperimental(&lineHoleGeoPolygon, 9, uint32(ContainmentOverlappingBBox), numHexagons, hexagons)
	if err != eSuccess {
		t.Fatalf("polygonToCellsExperimental failed: %v", err)
	}

	actualNumIndexes := countNonNullIndexesWithSize(hexagons, numHexagons)
	// Same as without the hole
	if actualNumIndexes != 1416 {
		t.Errorf("Expected 1416 polygonToCells size (line hole), got %d", actualNumIndexes)
	}
}

func TestInvalidFlags(t *testing.T) {
	t.Parallel()
	for flags := uint32(uint32(ContainmentInvalid)); flags <= 32; flags++ {
		_, err := maxPolygonToCellsSizeExperimental(&sfGeoPolygon, 9, flags)
		if err != eOptionInvalid {
			t.Errorf("Expected eOptionInvalid for invalid flags %d in maxPolygonToCellsSizeExperimental, got %v", flags, err)
		}
	}

	numHexagons, err := maxPolygonToCellsSizeExperimental(&sfGeoPolygon, 9, uint32(ContainmentCenter))
	if err != eSuccess {
		t.Fatalf("maxPolygonToCellsSizeExperimental failed: %v", err)
	}

	hexagons := make([]h3Index, numHexagons)
	for flags := uint32(uint32(ContainmentInvalid)); flags <= 32; flags++ {
		err := polygonToCellsExperimental(&sfGeoPolygon, 9, flags, numHexagons, hexagons)
		if err != eOptionInvalid {
			t.Errorf("Expected eOptionInvalid for invalid flags %d in polygonToCellsExperimental, got %v", flags, err)
		}
	}
}

// Note: The C fillIndex test uses iterateAllIndexesAtRes which is not available in Go yet.
// This test would need to be implemented once that function is ported.
