// Tests ported from H3 v4.5.0: src/apps/testapps/testCellsToMultiPoly.c.
package h3

import (
	"math"
	"testing"
)

func relativeDiff(a, b float64) float64 {
	diff := math.Abs(a - b)
	denom := math.Max(math.Abs(a), math.Abs(b))

	// If both numbers are extremely close to zero, fall back to absolute diff
	if denom < dblEpsilon {
		return diff
	}

	return diff / denom
}

func getLoopArea(t *testing.T, loop GeoLoop) float64 {
	t.Helper()
	area, err := geoLoopAreaRads2(loop)
	if err != eSuccess {
		t.Fatalf("geoLoopAreaRads2: %v", err)
	}

	if area < 0 {
		t.Error("Area should be nonnegative")
	}
	if area >= 4.0*math.Pi {
		t.Error("Area should be less than entire globe")
	}
	if len(loop) < 3 {
		t.Error("Every loop should have at least 3 vertexes.")
	}

	return area
}

func getOuterLoopArea(t *testing.T, poly GeoPolygon) float64 {
	t.Helper()
	return getLoopArea(t, poly.GeoLoop)
}

func checkPoly(t *testing.T, poly GeoPolygon) {
	t.Helper()
	area, err := geoPolygonAreaRads2(poly)
	if err != eSuccess {
		t.Fatalf("geoPolygonAreaRads2: %v", err)
	}

	if area < 0 {
		t.Error("Area should be nonnegative")
	}
	if area >= 4.0*math.Pi {
		t.Error("Area should be less than entire globe")
	}
	if area > getOuterLoopArea(t, poly) {
		t.Error("Total area should be less than the outer loop area")
	}

	// The outer ring and holes should be ordered in 'increasing' order; that
	// is, since the holes are oriented clockwise, they will naively enclose
	// more area than the outer ring, which is oriented counterclockwise.
	if len(poly.Holes) > 0 {
		if getLoopArea(t, poly.GeoLoop) > getLoopArea(t, poly.Holes[0]) {
			t.Error("Outer loop should have 'less' area than first hole.")
		}
	}
	for i := 0; i < len(poly.Holes)-1; i++ {
		if getLoopArea(t, poly.Holes[i]) > getLoopArea(t, poly.Holes[i+1]) {
			t.Error("Polygon holes should be listed in 'increasing' order.")
		}
	}
}

func getMpoly(t *testing.T, cells []h3Index, numCells int64) geoMultiPolygon {
	t.Helper()
	const relTol = 1e-8
	var mpoly geoMultiPolygon
	if err := cellsToMultiPolygon(cells, numCells, &mpoly); err != eSuccess {
		t.Fatalf("cellsToMultiPolygon: %v", err)
	}

	for i := int32(0); i < mpoly.NumPolygons-1; i++ {
		if getOuterLoopArea(t, mpoly.Polygons[i]) < getOuterLoopArea(t, mpoly.Polygons[i+1]) {
			t.Error("Polygons should ordered by area enclosed by outer loop; decreasing.")
		}
	}

	for i := int32(0); i < mpoly.NumPolygons; i++ {
		checkPoly(t, mpoly.Polygons[i])
	}

	{
		// Check that area matches sum of areas of cells
		polyArea, _ := geoMultiPolygonAreaRads2(mpoly)

		var k adder
		for i := int64(0); i < numCells; i++ {
			temp, err := cellAreaRads2(cells[i])
			if err != eSuccess {
				t.Fatalf("cellAreaRads2: %v", err)
			}
			kadd(&k, temp)
		}
		cellArea := k.sum
		if relativeDiff(cellArea, polyArea) > relTol {
			t.Error("Polygon area should match summing area of cells")
		}
	}

	return mpoly
}

func checkCell(t *testing.T, cell h3Index) {
	t.Helper()
	mpoly := getMpoly(t, []h3Index{cell}, 1)
	if mpoly.NumPolygons != 1 {
		t.Error("Exactly one polygon.")
	}
	if len(mpoly.Polygons[0].Holes) != 0 {
		t.Error("Individual cells should have zero holes.")
	}
	if len(mpoly.Polygons[0].GeoLoop) < 5 {
		t.Error("Individual cells should have at least 5 verices")
	}
	if len(mpoly.Polygons[0].GeoLoop) > 10 {
		t.Error("Individual cells should never have more than 10 verices")
	}
	destroyGeoMultiPolygon(&mpoly)
}

func checkGlobalPoly(t *testing.T, mpoly geoMultiPolygon) {
	t.Helper()
	if mpoly.NumPolygons != 8 {
		t.Error("Expecting 8 polygons")
	}
	for i := 0; i < 8; i++ {
		if len(mpoly.Polygons[i].Holes) != 0 || mpoly.Polygons[i].Holes != nil ||
			len(mpoly.Polygons[i].GeoLoop) != 3 {
			t.Error("Expecting each polygon is a triangle")
		}
	}

	area, _ := geoMultiPolygonAreaRads2(mpoly)
	// Upstream asserts area == 4π exactly; the Cagnoli edge terms run
	// through libm trig, where Go's math library and the platform libm
	// legitimately differ by an ulp on some inputs (see
	// area_geoLoopAreaRads2_parity_test.go — bit-exact area parity is
	// not attainable across libms), so this port admits exactly one ulp
	// of 4π. Measured: Go returns the immediate successor of 4π here.
	ulp4pi := math.Nextafter(4*math.Pi, math.Inf(1)) - 4*math.Pi
	if math.Abs(area-4*math.Pi) > ulp4pi {
		t.Errorf("Exact area expected (±1 ulp): got %v, want %v", area, 4*math.Pi)
	}
}

func TestCellsToMultiPoly_three_polygons(t *testing.T) {
	t.Parallel()

	const relTol = 1e-15
	// Results in 3 polygons: 0 holes, 1 hole, 3 holes
	cells := []h3Index{
		0x8027fffffffffff, 0x802bfffffffffff, 0x804dfffffffffff,
		0x8067fffffffffff, 0x806dfffffffffff, 0x8049fffffffffff,
		0x805ffffffffffff, 0x8057fffffffffff, 0x807dfffffffffff,
		0x80a5fffffffffff, 0x80a9fffffffffff, 0x808bfffffffffff,
		0x801bfffffffffff, 0x8035fffffffffff, 0x803ffffffffffff,
		0x8053fffffffffff, 0x8043fffffffffff, 0x8021fffffffffff,
		0x8011fffffffffff, 0x801ffffffffffff, 0x8097fffffffffff,
	}

	mpoly := getMpoly(t, cells, int64(len(cells)))
	if mpoly.NumPolygons != 3 {
		t.Error("expecting 3 polygons")
	}

	if len(mpoly.Polygons[0].Holes) != 3 {
		t.Error("3 holes in first poly")
	}
	if len(mpoly.Polygons[1].Holes) != 1 {
		t.Error("1 hole in second poly")
	}
	if len(mpoly.Polygons[2].Holes) != 0 {
		t.Error("0 holes in third poly")
	}

	expected := 2.2440497074541694
	area, _ := geoMultiPolygonAreaRads2(mpoly)
	if relativeDiff(area, expected) >= relTol {
		t.Errorf("Expected area: got %v, want %v", area, expected)
	}

	destroyGeoMultiPolygon(&mpoly)
}

func TestCellsToMultiPoly_cells_at_res(t *testing.T) {
	t.Parallel()

	_iterateAllIndexesAtRes(0, func(cell h3Index) { checkCell(t, cell) })
	_iterateAllIndexesAtRes(1, func(cell h3Index) { checkCell(t, cell) })
	_iterateAllIndexesAtRes(2, func(cell h3Index) { checkCell(t, cell) })
}

func TestCellsToMultiPoly_res15_hex(t *testing.T) {
	t.Parallel()

	// 0x8f754e64992d6d8 is a res 15 *hex*
	checkCell(t, 0x8f754e64992d6d8)
}

func TestCellsToMultiPoly_all_pentagons(t *testing.T) {
	t.Parallel()

	// check all pentagons at all resolutions
	var cells [12]h3Index
	for res := int32(0); res <= 15; res++ {
		if err := getPentagons(res, cells[:]); err != eSuccess {
			t.Fatalf("getPentagons(%d): %v", res, err)
		}
		for i := 0; i < 12; i++ {
			checkCell(t, cells[i])
		}
	}
}

func TestCellsToMultiPoly_issue_1049(t *testing.T) {
	t.Parallel()

	// from https://github.com/uber/h3/issues/1049
	cells := []h3Index{
		0x827487fffffffff, 0x82748ffffffffff, 0x827497fffffffff,
		0x82749ffffffffff, 0x8274affffffffff, 0x8274c7fffffffff,
		0x8274cffffffffff, 0x8274d7fffffffff, 0x8274e7fffffffff,
		0x8274effffffffff, 0x8274f7fffffffff, 0x82754ffffffffff,
		0x827c07fffffffff, 0x827c27fffffffff, 0x827c2ffffffffff,
		0x827c37fffffffff, 0x827d87fffffffff, 0x827d8ffffffffff,
		0x827d97fffffffff, 0x827d9ffffffffff, 0x827da7fffffffff,
		0x827daffffffffff, 0x82801ffffffffff, 0x8280a7fffffffff,
		0x8280affffffffff, 0x8280b7fffffffff, 0x828197fffffffff,
		0x82819ffffffffff, 0x8281a7fffffffff, 0x8281b7fffffffff,
		0x828207fffffffff, 0x82820ffffffffff, 0x828227fffffffff,
		0x82822ffffffffff, 0x8282e7fffffffff, 0x828307fffffffff,
		0x82830ffffffffff, 0x82831ffffffffff, 0x82832ffffffffff,
		0x828347fffffffff, 0x82834ffffffffff, 0x828357fffffffff,
		0x82835ffffffffff, 0x828367fffffffff, 0x828377fffffffff,
		0x82a447fffffffff, 0x82a457fffffffff, 0x82a45ffffffffff,
		0x82a467fffffffff, 0x82a46ffffffffff, 0x82a477fffffffff,
		0x82a4c7fffffffff, 0x82a4cffffffffff, 0x82a4d7fffffffff,
		0x82a4e7fffffffff, 0x82a4effffffffff, 0x82a4f7fffffffff,
		0x82a547fffffffff, 0x82a54ffffffffff, 0x82a557fffffffff,
		0x82a55ffffffffff, 0x82a567fffffffff, 0x82a577fffffffff,
		0x82a837fffffffff, 0x82a897fffffffff, 0x82a8a7fffffffff,
		0x82a8b7fffffffff, 0x82a917fffffffff, 0x82a927fffffffff,
		0x82a937fffffffff, 0x82a987fffffffff, 0x82a98ffffffffff,
		0x82a997fffffffff, 0x82a99ffffffffff, 0x82a9a7fffffffff,
		0x82a9affffffffff, 0x82ac47fffffffff, 0x82ac57fffffffff,
		0x82ac5ffffffffff, 0x82ac67fffffffff, 0x82ac6ffffffffff,
		0x82ac77fffffffff, 0x82ad47fffffffff, 0x82ad4ffffffffff,
		0x82ad57fffffffff, 0x82ad5ffffffffff, 0x82ad67fffffffff,
		0x82ad77fffffffff, 0x82c207fffffffff, 0x82c217fffffffff,
		0x82c227fffffffff, 0x82c237fffffffff, 0x82c287fffffffff,
		0x82c28ffffffffff, 0x82c29ffffffffff, 0x82c2a7fffffffff,
		0x82c2affffffffff, 0x82c2b7fffffffff, 0x82c307fffffffff,
		0x82c317fffffffff, 0x82c31ffffffffff, 0x82c337fffffffff,
		0x82cfb7fffffffff, 0x82d0c7fffffffff, 0x82d0d7fffffffff,
		0x82d0dffffffffff, 0x82d0e7fffffffff, 0x82d0f7fffffffff,
		0x82d147fffffffff, 0x82d157fffffffff, 0x82d15ffffffffff,
		0x82d167fffffffff, 0x82d177fffffffff, 0x82d187fffffffff,
		0x82d18ffffffffff, 0x82d197fffffffff, 0x82d19ffffffffff,
		0x82d1a7fffffffff, 0x82d1affffffffff, 0x82dc47fffffffff,
		0x82dc57fffffffff, 0x82dc5ffffffffff, 0x82dc67fffffffff,
		0x82dc6ffffffffff, 0x82dc77fffffffff, 0x82dcc7fffffffff,
		0x82dccffffffffff, 0x82dcd7fffffffff, 0x82dce7fffffffff,
		0x82dceffffffffff, 0x82dcf7fffffffff, 0x82dd1ffffffffff,
		0x82dd47fffffffff, 0x82dd4ffffffffff, 0x82dd57fffffffff,
		0x82dd5ffffffffff, 0x82dd6ffffffffff, 0x82dd87fffffffff,
		0x82dd8ffffffffff, 0x82dd97fffffffff, 0x82dd9ffffffffff,
		0x82ddaffffffffff, 0x82ddb7fffffffff, 0x82dec7fffffffff,
		0x82decffffffffff, 0x82ded7fffffffff, 0x82dee7fffffffff,
		0x82deeffffffffff, 0x82def7fffffffff, 0x82df0ffffffffff,
		0x82df1ffffffffff, 0x82df47fffffffff, 0x82df4ffffffffff,
		0x82df57fffffffff, 0x82df5ffffffffff, 0x82df77fffffffff,
		0x82df8ffffffffff, 0x82df9ffffffffff, 0x82e6c7fffffffff,
		0x82e6cffffffffff, 0x82e6d7fffffffff, 0x82e6dffffffffff,
		0x82e6effffffffff, 0x82e6f7fffffffff,
	}

	mpoly := getMpoly(t, cells, int64(len(cells)))

	if mpoly.NumPolygons != 12 {
		t.Error("expecting 12 polygons")
	}

	for i := int32(0); i < mpoly.NumPolygons; i++ {
		if len(mpoly.Polygons[i].Holes) != 0 {
			t.Error("expecting 0 holes")
		}
	}

	destroyGeoMultiPolygon(&mpoly)
}

func TestCellsToMultiPoly_equator_cells(t *testing.T) {
	t.Parallel()

	// A "global polygon" example.
	cells := []h3Index{
		0x81807ffffffffff, 0x817efffffffffff, 0x81723ffffffffff,
		0x817ebffffffffff, 0x817c3ffffffffff, 0x817e3ffffffffff,
		0x817a3ffffffffff, 0x8166fffffffffff, 0x8172bffffffffff,
		0x816afffffffffff, 0x81933ffffffffff, 0x8168fffffffffff,
		0x8188fffffffffff, 0x81853ffffffffff, 0x817f7ffffffffff,
		0x8180bffffffffff, 0x81783ffffffffff, 0x81743ffffffffff,
		0x8170bffffffffff, 0x8173bffffffffff, 0x8179bffffffffff,
		0x817cbffffffffff, 0x8188bffffffffff, 0x81857ffffffffff,
		0x816f7ffffffffff, 0x8177bffffffffff, 0x81617ffffffffff,
		0x816f3ffffffffff, 0x8174bffffffffff, 0x8180fffffffffff,
		0x817a7ffffffffff, 0x81767ffffffffff, 0x81757ffffffffff,
		0x81957ffffffffff, 0x81787ffffffffff, 0x81847ffffffffff,
		0x81653ffffffffff, 0x817bbffffffffff, 0x816cfffffffffff,
		0x816abffffffffff, 0x815f3ffffffffff, 0x817c7ffffffffff,
		0x8168bffffffffff, 0x818cbffffffffff, 0x818cfffffffffff,
		0x818afffffffffff, 0x8174fffffffffff, 0x8172fffffffffff,
		0x8170fffffffffff, 0x816fbffffffffff, 0x81657ffffffffff,
		0x816c7ffffffffff, 0x8186bffffffffff, 0x81763ffffffffff,
		0x818a7ffffffffff, 0x8186fffffffffff, 0x81707ffffffffff,
		0x8182bffffffffff, 0x818f3ffffffffff, 0x8182fffffffffff,
	}
	mpoly := getMpoly(t, cells, int64(len(cells)))

	if mpoly.NumPolygons != 1 {
		t.Error("expecting 1 polygon")
	}
	if len(mpoly.Polygons[0].Holes) != 1 {
		t.Error("expecting 1 hole")
	}

	destroyGeoMultiPolygon(&mpoly)
}

func TestCellsToMultiPoly_prime_meridian(t *testing.T) {
	t.Parallel()

	// A "global polygon" example.
	cells := []h3Index{
		0x81efbffffffffff, 0x81c07ffffffffff, 0x81d1bffffffffff,
		0x81097ffffffffff, 0x8109bffffffffff, 0x81d0bffffffffff,
		0x81987ffffffffff, 0x81017ffffffffff, 0x81e67ffffffffff,
		0x81ddbffffffffff, 0x81ac7ffffffffff, 0x8158bffffffffff,
		0x81397ffffffffff, 0x81593ffffffffff, 0x81c17ffffffffff,
		0x81827ffffffffff, 0x81197ffffffffff, 0x81eebffffffffff,
		0x81383ffffffffff, 0x81dcbffffffffff, 0x81757ffffffffff,
		0x81093ffffffffff, 0x81073ffffffffff, 0x8159bffffffffff,
		0x81f17ffffffffff, 0x81187ffffffffff, 0x81007ffffffffff,
		0x81997ffffffffff, 0x81753ffffffffff, 0x81033ffffffffff,
		0x81f2bffffffffff, 0x8138bffffffffff,
	}
	mpoly := getMpoly(t, cells, int64(len(cells)))

	if mpoly.NumPolygons != 1 {
		t.Error("expecting 1 polygon")
	}
	if len(mpoly.Polygons[0].Holes) != 0 {
		t.Error("expecting 0 holes")
	}

	destroyGeoMultiPolygon(&mpoly)
}

func TestCellsToMultiPoly_anti_meridian(t *testing.T) {
	t.Parallel()

	// A "global polygon" example.
	cells := []h3Index{
		0x817ebffffffffff, 0x8133bffffffffff, 0x81047ffffffffff,
		0x81f3bffffffffff, 0x81dbbffffffffff, 0x8132bffffffffff,
		0x810cbffffffffff, 0x81bb3ffffffffff, 0x81db3ffffffffff,
		0x81bafffffffffff, 0x81177ffffffffff, 0x817fbffffffffff,
		0x81ba3ffffffffff, 0x815abffffffffff, 0x815bbffffffffff,
		0x81eafffffffffff, 0x81ed7ffffffffff, 0x81057ffffffffff,
		0x819a7ffffffffff, 0x81eabffffffffff, 0x819b7ffffffffff,
		0x81167ffffffffff, 0x81227ffffffffff, 0x8171bffffffffff,
		0x81237ffffffffff, 0x810dbffffffffff, 0x81033ffffffffff,
		0x81f2bffffffffff, 0x8147bffffffffff, 0x81f33ffffffffff,
	}
	mpoly := getMpoly(t, cells, int64(len(cells)))

	if mpoly.NumPolygons != 1 {
		t.Error("expecting 1 polygon")
	}
	if len(mpoly.Polygons[0].Holes) != 0 {
		t.Error("expecting 0 holes")
	}

	destroyGeoMultiPolygon(&mpoly)
}

func TestCellsToMultiPoly_both_meridians(t *testing.T) {
	t.Parallel()

	// A "global polygon" example.
	cells := []h3Index{
		0x81efbffffffffff, 0x81c07ffffffffff, 0x81d1bffffffffff,
		0x81097ffffffffff, 0x817ebffffffffff, 0x8133bffffffffff,
		0x81047ffffffffff, 0x8109bffffffffff, 0x81f3bffffffffff,
		0x81d0bffffffffff, 0x81987ffffffffff, 0x81dbbffffffffff,
		0x81017ffffffffff, 0x81e67ffffffffff, 0x81ddbffffffffff,
		0x8132bffffffffff, 0x810cbffffffffff, 0x81bb3ffffffffff,
		0x81ac7ffffffffff, 0x81db3ffffffffff, 0x8158bffffffffff,
		0x81397ffffffffff, 0x81593ffffffffff, 0x81bafffffffffff,
		0x81177ffffffffff, 0x817fbffffffffff, 0x81ba3ffffffffff,
		0x81c17ffffffffff, 0x815abffffffffff, 0x81827ffffffffff,
		0x815bbffffffffff, 0x81eafffffffffff, 0x81197ffffffffff,
		0x81ed7ffffffffff, 0x81eebffffffffff, 0x81383ffffffffff,
		0x81057ffffffffff, 0x819a7ffffffffff, 0x81dcbffffffffff,
		0x81757ffffffffff, 0x81eabffffffffff, 0x81093ffffffffff,
		0x819b7ffffffffff, 0x81073ffffffffff, 0x8159bffffffffff,
		0x8147bffffffffff, 0x81167ffffffffff, 0x81f17ffffffffff,
		0x8171bffffffffff, 0x81227ffffffffff, 0x81187ffffffffff,
		0x81237ffffffffff, 0x81007ffffffffff, 0x810dbffffffffff,
		0x81997ffffffffff, 0x81753ffffffffff, 0x81033ffffffffff,
		0x81f2bffffffffff, 0x8138bffffffffff, 0x81f33ffffffffff,
	}
	mpoly := getMpoly(t, cells, int64(len(cells)))

	if mpoly.NumPolygons != 1 {
		t.Error("expecting 1 polygon")
	}
	if len(mpoly.Polygons[0].Holes) != 1 {
		t.Error("expecting 1 hole")
	}

	destroyGeoMultiPolygon(&mpoly)
}

func TestCellsToMultiPoly_meridians_and_equator(t *testing.T) {
	t.Parallel()

	// A "global polygon" example.
	cells := []h3Index{
		0x817c3ffffffffff, 0x81047ffffffffff, 0x8188fffffffffff,
		0x817f7ffffffffff, 0x8180bffffffffff, 0x81177ffffffffff,
		0x817fbffffffffff, 0x8188bffffffffff, 0x815bbffffffffff,
		0x81eafffffffffff, 0x816f3ffffffffff, 0x817a7ffffffffff,
		0x819a7ffffffffff, 0x81757ffffffffff, 0x817bbffffffffff,
		0x816cfffffffffff, 0x8168bffffffffff, 0x81237ffffffffff,
		0x818afffffffffff, 0x8172fffffffffff, 0x816fbffffffffff,
		0x81657ffffffffff, 0x81763ffffffffff, 0x818a7ffffffffff,
		0x81eabffffffffff, 0x8138bffffffffff, 0x8182fffffffffff,
		0x81c07ffffffffff, 0x8109bffffffffff, 0x8166fffffffffff,
		0x81987ffffffffff, 0x8172bffffffffff, 0x8168fffffffffff,
		0x81853ffffffffff, 0x810cbffffffffff, 0x81bb3ffffffffff,
		0x81db3ffffffffff, 0x81743ffffffffff, 0x81bafffffffffff,
		0x8179bffffffffff, 0x818f3ffffffffff, 0x81857ffffffffff,
		0x816f7ffffffffff, 0x8177bffffffffff, 0x8174bffffffffff,
		0x81eebffffffffff, 0x81383ffffffffff, 0x81767ffffffffff,
		0x81787ffffffffff, 0x819b7ffffffffff, 0x8159bffffffffff,
		0x8171bffffffffff, 0x818cbffffffffff, 0x818cfffffffffff,
		0x8170fffffffffff, 0x81707ffffffffff, 0x8147bffffffffff,
		0x81167ffffffffff, 0x81f33ffffffffff, 0x817efffffffffff,
		0x81f3bffffffffff, 0x81017ffffffffff, 0x816afffffffffff,
		0x81e67ffffffffff, 0x81ddbffffffffff, 0x8132bffffffffff,
		0x8170bffffffffff, 0x81ba3ffffffffff, 0x81c17ffffffffff,
		0x815abffffffffff, 0x81617ffffffffff, 0x8180fffffffffff,
		0x81dcbffffffffff, 0x81957ffffffffff, 0x81093ffffffffff,
		0x81847ffffffffff, 0x81653ffffffffff, 0x81073ffffffffff,
		0x8174fffffffffff, 0x810dbffffffffff, 0x81997ffffffffff,
		0x816c7ffffffffff, 0x81033ffffffffff, 0x8186bffffffffff,
		0x81f2bffffffffff, 0x81efbffffffffff, 0x81807ffffffffff,
		0x81d1bffffffffff, 0x81097ffffffffff, 0x817ebffffffffff,
		0x81723ffffffffff, 0x8133bffffffffff, 0x817e3ffffffffff,
		0x817a3ffffffffff, 0x81d0bffffffffff, 0x81dbbffffffffff,
		0x81933ffffffffff, 0x81783ffffffffff, 0x81ac7ffffffffff,
		0x8158bffffffffff, 0x81397ffffffffff, 0x81593ffffffffff,
		0x8173bffffffffff, 0x817cbffffffffff, 0x81827ffffffffff,
		0x81197ffffffffff, 0x81ed7ffffffffff, 0x81057ffffffffff,
		0x816abffffffffff, 0x815f3ffffffffff, 0x81f17ffffffffff,
		0x81227ffffffffff, 0x817c7ffffffffff, 0x81007ffffffffff,
		0x81753ffffffffff, 0x8186fffffffffff, 0x8182bffffffffff,
		0x81187ffffffffff,
	}
	mpoly := getMpoly(t, cells, int64(len(cells)))

	if mpoly.NumPolygons != 1 {
		t.Error("expecting 1 polygon")
	}
	if len(mpoly.Polygons[0].Holes) != 3 {
		t.Error("expecting 3 holes")
	}

	destroyGeoMultiPolygon(&mpoly)
}

func TestCellsToMultiPoly_negative_cells(t *testing.T) {
	t.Parallel()

	var mpoly geoMultiPolygon
	if err := cellsToMultiPolygon(nil, -1, &mpoly); err != eDomain {
		t.Errorf("Can't pass in negative number of cells: got %v, want eDomain", err)
	}
}

func TestCellsToMultiPoly_empty_cells(t *testing.T) {
	t.Parallel()

	mpoly := getMpoly(t, nil, 0)

	if mpoly.NumPolygons != 0 {
		t.Error("expecting 0 polygons")
	}
	if mpoly.Polygons != nil {
		t.Error("Pointer should be NULL")
	}

	destroyGeoMultiPolygon(&mpoly)
}

func TestCellsToMultiPoly_all_cells(t *testing.T) {
	t.Parallel()

	var cells [122]h3Index
	if err := getRes0Cells(cells[:]); err != eSuccess {
		t.Fatalf("getRes0Cells: %v", err)
	}

	// expecting a global multipolygon
	mpoly := getMpoly(t, cells[:], 122)

	checkGlobalPoly(t, mpoly)

	destroyGeoMultiPolygon(&mpoly)
}

func TestCellsToMultiPoly_duplicate_cells(t *testing.T) {
	t.Parallel()

	cells := []h3Index{
		0x81efbffffffffff,
		0x81efbffffffffff,
		0x81efbffffffffff,
	}

	var mpoly geoMultiPolygon
	if err := cellsToMultiPolygon(cells, int64(len(cells)), &mpoly); err != eDuplicateInput {
		t.Errorf("Can't have duplicated cells: got %v, want eDuplicateInput", err)
	}
}

func TestCellsToMultiPoly_multiple_resolutions(t *testing.T) {
	t.Parallel()

	cells := []h3Index{
		0x8027fffffffffff,
		0x81efbffffffffff,
	}

	var mpoly geoMultiPolygon
	if err := cellsToMultiPolygon(cells, int64(len(cells)), &mpoly); err != eResMismatch {
		t.Errorf("Can't have multiple cell resolutions: got %v, want eResMismatch", err)
	}
}

func TestCellsToMultiPoly_invalid_cells(t *testing.T) {
	t.Parallel()

	cells := []h3Index{
		0x8027fffffffffff,
		0x81efbffffffffff,
	}
	cells[1] += 1 // make cell invalid

	var mpoly geoMultiPolygon
	if err := cellsToMultiPolygon(cells, int64(len(cells)), &mpoly); err != eCellInvalid {
		t.Errorf("Can't have invalid cells: got %v, want eCellInvalid", err)
	}
}

func TestCellsToMultiPoly_overflow_check(t *testing.T) {
	t.Parallel()

	// Test an absurdly large numCells returns
	// Use 1000x the number of cells at resolution 15
	numCellsRes15, err := getNumCells(15)
	if err != eSuccess {
		t.Fatalf("getNumCells: %v", err)
	}

	absurdNumCells := numCellsRes15 * 1000

	var mpoly geoMultiPolygon
	if got := cellsToMultiPolygon(nil, absurdNumCells, &mpoly); got != eMemoryBounds {
		t.Errorf("Should return eMemoryBounds for absurdly large numCells: got %v", got)
	}
}
