// Tests ported from testPolygonInternal.c
package h3

import (
	"math"
	"testing"
)

func TestPointInsideGeoLoop(t *testing.T) {
	t.Parallel()

	// Fixture from C code
	sfVerts := []LatLng{
		{0.659966917655, -2.1364398519396}, {0.6595011102219, -2.1359434279405},
		{0.6583348114025, -2.1354884206045}, {0.6581220034068, -2.1382437718946},
		{0.6594479998527, -2.1384597563896}, {0.6599990002976, -2.1376771158464},
	}

	inside := LatLng{0.659, -2.136}
	somewhere := LatLng{1, 2}

	var bbox BBox
	bboxFromGeoLoop(sfVerts, &bbox)

	// For exact points on the polygon, we bias west and south
	result := pointInsideGeoLoop(sfVerts, &bbox, &sfVerts[0])
	if result {
		t.Error("does not contain exact vertex 0")
	}

	result = pointInsideGeoLoop(sfVerts, &bbox, &sfVerts[3])
	if !result {
		t.Error("contains exact vertex 3")
	}

	result = pointInsideGeoLoop(sfVerts, &bbox, &inside)
	if !result {
		t.Error("contains point inside")
	}

	result = pointInsideGeoLoop(sfVerts, &bbox, &somewhere)
	if result {
		t.Error("contains somewhere else")
	}
}

func TestPointInsideGeoLoopCornerCases(t *testing.T) {
	t.Parallel()

	verts := []LatLng{{0, 0}, {1, 0}, {1, 1}, {0, 1}}

	var bbox BBox
	bboxFromGeoLoop(verts, &bbox)

	point := LatLng{0, 0}

	// Test corners. For exact points on the polygon, we bias west and
	// north, so only the southeast corner is contained.
	result := pointInsideGeoLoop(verts, &bbox, &point)
	if result {
		t.Error("does not contain sw corner")
	}

	point.Lat = 1
	result = pointInsideGeoLoop(verts, &bbox, &point)
	if result {
		t.Error("does not contain nw corner")
	}

	point.Lng = 1
	result = pointInsideGeoLoop(verts, &bbox, &point)
	if result {
		t.Error("does not contain ne corner")
	}

	point.Lat = 0
	result = pointInsideGeoLoop(verts, &bbox, &point)
	if !result {
		t.Error("contains se corner")
	}
}

func TestPointInsideGeoLoopEdgeCases(t *testing.T) {
	t.Parallel()

	verts := []LatLng{{0, 0}, {1, 0}, {1, 1}, {0, 1}}

	var bbox BBox
	bboxFromGeoLoop(verts, &bbox)

	// Test edges. Only points on south and east edges are contained.
	point := LatLng{0.5, 0}
	result := pointInsideGeoLoop(verts, &bbox, &point)
	if result {
		t.Error("does not contain point on west edge")
	}

	point = LatLng{1, 0.5}
	result = pointInsideGeoLoop(verts, &bbox, &point)
	if result {
		t.Error("does not contain point on north edge")
	}

	point = LatLng{0.5, 1}
	result = pointInsideGeoLoop(verts, &bbox, &point)
	if !result {
		t.Error("contains point on east edge")
	}

	point = LatLng{0, 0.5}
	result = pointInsideGeoLoop(verts, &bbox, &point)
	if !result {
		t.Error("contains point on south edge")
	}
}

func TestPointInsideGeoLoopExtraEdgeCase(t *testing.T) {
	t.Parallel()

	// This is a carefully crafted shape + point to hit an otherwise
	// missed branch in coverage
	verts := []LatLng{{0, 0}, {1, 0.5}, {0, 1}}

	var bbox BBox
	bboxFromGeoLoop(verts, &bbox)

	point := LatLng{0.5, 0.5}
	result := pointInsideGeoLoop(verts, &bbox, &point)
	if !result {
		t.Error("contains inside point matching longitude of a vertex")
	}
}

func TestPointInsideGeoLoopTransmeridian(t *testing.T) {
	t.Parallel()

	verts := []LatLng{
		{0.01, -math.Pi + 0.01},
		{0.01, math.Pi - 0.01},
		{-0.01, math.Pi - 0.01},
		{-0.01, -math.Pi + 0.01},
	}

	eastPoint := LatLng{0.001, -math.Pi + 0.001}
	eastPointOutside := LatLng{0.001, -math.Pi + 0.1}
	westPoint := LatLng{0.001, math.Pi - 0.001}
	westPointOutside := LatLng{0.001, math.Pi - 0.1}

	var bbox BBox
	bboxFromGeoLoop(verts, &bbox)

	result := pointInsideGeoLoop(verts, &bbox, &westPoint)
	if !result {
		t.Error("contains point to the west of the antimeridian")
	}

	result = pointInsideGeoLoop(verts, &bbox, &eastPoint)
	if !result {
		t.Error("contains point to the east of the antimeridian")
	}

	result = pointInsideGeoLoop(verts, &bbox, &westPointOutside)
	if result {
		t.Error("does not contain outside point to the west of the antimeridian")
	}

	result = pointInsideGeoLoop(verts, &bbox, &eastPointOutside)
	if result {
		t.Error("does not contain outside point to the east of the antimeridian")
	}
}

// Helper function to create a linked loop like the C code.
func createLinkedLoop(verts []LatLng) *LinkedGeoLoop {
	loop := &LinkedGeoLoop{}
	for i := range verts {
		addLinkedCoord(loop, &verts[i])
	}
	return loop
}

func TestPointInsideLinkedGeoLoop(t *testing.T) {
	t.Parallel()

	sfVerts := []LatLng{
		{0.659966917655, -2.1364398519396}, {0.6595011102219, -2.1359434279405},
		{0.6583348114025, -2.1354884206045}, {0.6581220034068, -2.1382437718946},
		{0.6594479998527, -2.1384597563896}, {0.6599990002976, -2.1376771158464},
	}

	somewhere := LatLng{1, 2}
	inside := LatLng{0.659, -2.136}

	loop := createLinkedLoop(sfVerts)
	defer destroyLinkedGeoLoop(loop)

	var bbox BBox
	bboxFromLinkedGeoLoop(loop, &bbox)

	result := pointInsideLinkedGeoLoop(loop, &bbox, &inside)
	if !result {
		t.Error("contains exact4")
	}

	result = pointInsideLinkedGeoLoop(loop, &bbox, &somewhere)
	if result {
		t.Error("contains somewhere else")
	}
}

func TestBboxFromGeoLoop(t *testing.T) {
	t.Parallel()

	verts := []LatLng{{0.8, 0.3}, {0.7, 0.6}, {1.1, 0.7}, {1.0, 0.2}}

	expected := BBox{1.1, 0.7, 0.7, 0.2}

	var result BBox
	bboxFromGeoLoop(verts, &result)

	if !bboxEquals(&result, &expected) {
		t.Errorf("Got expected bbox: expected %+v, got %+v", expected, result)
	}
}

func TestBboxFromGeoLoopTransmeridian(t *testing.T) {
	t.Parallel()

	verts := []LatLng{
		{0.1, -math.Pi + 0.1}, {0.1, math.Pi - 0.1},
		{0.05, math.Pi - 0.2}, {-0.1, math.Pi - 0.1},
		{-0.1, -math.Pi + 0.1}, {-0.05, -math.Pi + 0.2},
	}

	expected := BBox{0.1, -0.1, -math.Pi + 0.2, math.Pi - 0.2}

	var result BBox
	bboxFromGeoLoop(verts, &result)

	if !bboxEquals(&result, &expected) {
		t.Error("Got expected transmeridian bbox")
	}
}

func TestBboxFromGeoLoopNoVertices(t *testing.T) {
	t.Parallel()

	var verts []LatLng

	expected := BBox{0.0, 0.0, 0.0, 0.0}

	var result BBox
	bboxFromGeoLoop(verts, &result)

	if !bboxEquals(&result, &expected) {
		t.Error("Got expected bbox")
	}
}

func TestBboxesFromGeoPolygon(t *testing.T) {
	t.Parallel()

	verts := []LatLng{{0.8, 0.3}, {0.7, 0.6}, {1.1, 0.7}, {1.0, 0.2}}
	geoloop := GeoLoop(verts)
	polygon := GeoPolygon{GeoLoop: geoloop}

	expected := BBox{1.1, 0.7, 0.7, 0.2}

	result := make([]BBox, 1)
	bboxesFromGeoPolygon(&polygon, result)

	if !bboxEquals(&result[0], &expected) {
		t.Error("Got expected bbox")
	}
}

func TestBboxesFromGeoPolygonHole(t *testing.T) {
	t.Parallel()

	verts := []LatLng{{0.8, 0.3}, {0.7, 0.6}, {1.1, 0.7}, {1.0, 0.2}}
	geoloop := GeoLoop(verts)

	// not a real hole, but doesn't matter for the test
	holeVerts := []LatLng{{0.9, 0.3}, {0.9, 0.5}, {1.0, 0.7}, {0.9, 0.3}}
	holeGeoLoop := GeoLoop(holeVerts)

	polygon := GeoPolygon{
		GeoLoop: geoloop,
		Holes:   []GeoLoop{holeGeoLoop},
	}

	expected := BBox{1.1, 0.7, 0.7, 0.2}
	expectedHole := BBox{1.0, 0.9, 0.7, 0.3}

	result := make([]BBox, 2)
	bboxesFromGeoPolygon(&polygon, result)

	if !bboxEquals(&result[0], &expected) {
		t.Error("Got expected bbox")
	}
	if !bboxEquals(&result[1], &expectedHole) {
		t.Error("Got expected hole bbox")
	}
}

func TestBboxFromLinkedGeoLoop(t *testing.T) {
	t.Parallel()

	verts := []LatLng{{0.8, 0.3}, {0.7, 0.6}, {1.1, 0.7}, {1.0, 0.2}}

	loop := createLinkedLoop(verts)
	defer destroyLinkedGeoLoop(loop)

	expected := BBox{1.1, 0.7, 0.7, 0.2}

	var result BBox
	bboxFromLinkedGeoLoop(loop, &result)

	if !bboxEquals(&result, &expected) {
		t.Error("Got expected bbox")
	}
}

func TestBboxFromLinkedGeoLoopNoVertices(t *testing.T) {
	t.Parallel()

	loop := &LinkedGeoLoop{}

	expected := BBox{0.0, 0.0, 0.0, 0.0}

	var result BBox
	bboxFromLinkedGeoLoop(loop, &result)

	if !bboxEquals(&result, &expected) {
		t.Error("Got expected bbox")
	}

	destroyLinkedGeoLoop(loop)
}

func TestIsClockwiseGeoLoop(t *testing.T) {
	t.Parallel()

	verts := []LatLng{{0, 0}, {0.1, 0.1}, {0, 0.1}}
	geoloop := GeoLoop(verts)

	result := isClockwiseGeoLoop(geoloop)
	if !result {
		t.Error("Got true for clockwise geoloop")
	}
}

func TestIsClockwiseLinkedGeoLoop(t *testing.T) {
	t.Parallel()

	verts := []LatLng{{0.1, 0.1}, {0.2, 0.2}, {0.1, 0.2}}
	loop := createLinkedLoop(verts)
	defer destroyLinkedGeoLoop(loop)

	result := isClockwiseLinkedGeoLoop(loop)
	if !result {
		t.Error("Got true for clockwise loop")
	}
}

func TestIsNotClockwiseLinkedGeoLoop(t *testing.T) {
	t.Parallel()

	verts := []LatLng{{0, 0}, {0, 0.4}, {0.4, 0.4}, {0.4, 0}}
	loop := createLinkedLoop(verts)
	defer destroyLinkedGeoLoop(loop)

	result := isClockwiseLinkedGeoLoop(loop)
	if result {
		t.Error("Got false for counter-clockwise loop")
	}
}

func TestIsClockwiseGeoLoopTransmeridian(t *testing.T) {
	t.Parallel()

	verts := []LatLng{
		{0.4, math.Pi - 0.1},
		{0.4, -math.Pi + 0.1},
		{-0.4, -math.Pi + 0.1},
		{-0.4, math.Pi - 0.1},
	}
	geoloop := GeoLoop(verts)

	result := isClockwiseGeoLoop(geoloop)
	if !result {
		t.Error("Got true for clockwise geoloop")
	}
}

func TestIsClockwiseLinkedGeoLoopTransmeridian(t *testing.T) {
	t.Parallel()

	verts := []LatLng{
		{0.4, math.Pi - 0.1},
		{0.4, -math.Pi + 0.1},
		{-0.4, -math.Pi + 0.1},
		{-0.4, math.Pi - 0.1},
	}
	loop := createLinkedLoop(verts)
	defer destroyLinkedGeoLoop(loop)

	result := isClockwiseLinkedGeoLoop(loop)
	if !result {
		t.Error("Got true for clockwise transmeridian loop")
	}
}

func TestIsNotClockwiseLinkedGeoLoopTransmeridian(t *testing.T) {
	t.Parallel()

	verts := []LatLng{
		{0.4, math.Pi - 0.1},
		{-0.4, math.Pi - 0.1},
		{-0.4, -math.Pi + 0.1},
		{0.4, -math.Pi + 0.1},
	}
	loop := createLinkedLoop(verts)
	defer destroyLinkedGeoLoop(loop)

	result := isClockwiseLinkedGeoLoop(loop)
	if result {
		t.Error("Got false for counter-clockwise transmeridian loop")
	}
}

func TestNormalizeMultiPolygonSingle(t *testing.T) {
	t.Parallel()

	verts := []LatLng{{0, 0}, {0, 1}, {1, 1}}

	outer := createLinkedLoop(verts)

	polygon := LinkedGeoPolygon{}
	addLinkedLoop(&polygon, outer)

	err := normalizeMultiPolygon(&polygon)
	if err != E_SUCCESS {
		t.Errorf("Expected success, got %v", err)
	}

	if countLinkedPolygons(&polygon) != 1 {
		t.Error("Polygon count correct")
	}
	if countLinkedLoops(&polygon) != 1 {
		t.Error("Loop count correct")
	}
	if polygon.First != outer {
		t.Error("Got expected loop")
	}

	destroyLinkedMultiPolygon(&polygon)
}

func TestNormalizeMultiPolygonTwoOuterLoops(t *testing.T) {
	t.Parallel()

	verts1 := []LatLng{{0, 0}, {0, 1}, {1, 1}}
	outer1 := createLinkedLoop(verts1)

	verts2 := []LatLng{{2, 2}, {2, 3}, {3, 3}}
	outer2 := createLinkedLoop(verts2)

	polygon := LinkedGeoPolygon{}
	addLinkedLoop(&polygon, outer1)
	addLinkedLoop(&polygon, outer2)

	err := normalizeMultiPolygon(&polygon)
	if err != E_SUCCESS {
		t.Errorf("Expected success, got %v", err)
	}

	if countLinkedPolygons(&polygon) != 2 {
		t.Error("Polygon count correct")
	}
	if countLinkedLoops(&polygon) != 1 {
		t.Error("Loop count on first polygon correct")
	}
	if countLinkedLoops(polygon.Next) != 1 {
		t.Error("Loop count on second polygon correct")
	}

	destroyLinkedMultiPolygon(&polygon)
}

func TestNormalizeMultiPolygonOneHole(t *testing.T) {
	t.Parallel()

	verts := []LatLng{{0, 0}, {0, 3}, {3, 3}, {3, 0}}
	outer := createLinkedLoop(verts)

	verts2 := []LatLng{{1, 1}, {2, 2}, {1, 2}}
	inner := createLinkedLoop(verts2)

	polygon := LinkedGeoPolygon{}
	addLinkedLoop(&polygon, inner)
	addLinkedLoop(&polygon, outer)

	err := normalizeMultiPolygon(&polygon)
	if err != E_SUCCESS {
		t.Errorf("Expected success, got %v", err)
	}

	if countLinkedPolygons(&polygon) != 1 {
		t.Error("Polygon count correct")
	}
	if countLinkedLoops(&polygon) != 2 {
		t.Error("Loop count on first polygon correct")
	}
	if polygon.First != outer {
		t.Error("Got expected outer loop")
	}
	if polygon.First.Next != inner {
		t.Error("Got expected inner loop")
	}

	destroyLinkedMultiPolygon(&polygon)
}

func TestNormalizeMultiPolygonTwoHoles(t *testing.T) {
	t.Parallel()

	verts := []LatLng{{0, 0}, {0, 0.4}, {0.4, 0.4}, {0.4, 0}}
	outer := createLinkedLoop(verts)

	verts2 := []LatLng{{0.1, 0.1}, {0.2, 0.2}, {0.1, 0.2}}
	inner1 := createLinkedLoop(verts2)

	verts3 := []LatLng{{0.2, 0.2}, {0.3, 0.3}, {0.2, 0.3}}
	inner2 := createLinkedLoop(verts3)

	polygon := LinkedGeoPolygon{}
	addLinkedLoop(&polygon, inner2)
	addLinkedLoop(&polygon, outer)
	addLinkedLoop(&polygon, inner1)

	err := normalizeMultiPolygon(&polygon)
	if err != E_SUCCESS {
		t.Errorf("Expected success, got %v", err)
	}

	if countLinkedPolygons(&polygon) != 1 {
		t.Error("Polygon count correct for 2 holes")
	}
	if polygon.First != outer {
		t.Error("Got expected outer loop")
	}
	if countLinkedLoops(&polygon) != 3 {
		t.Error("Loop count on first polygon correct")
	}

	destroyLinkedMultiPolygon(&polygon)
}

func TestNormalizeMultiPolygonTwoDonuts(t *testing.T) {
	t.Parallel()

	verts := []LatLng{{0, 0}, {0, 3}, {3, 3}, {3, 0}}
	outer := createLinkedLoop(verts)

	verts2 := []LatLng{{1, 1}, {2, 2}, {1, 2}}
	inner := createLinkedLoop(verts2)

	verts3 := []LatLng{{0, 0}, {0, -3}, {-3, -3}, {-3, 0}}
	outer2 := createLinkedLoop(verts3)

	verts4 := []LatLng{{-1, -1}, {-2, -2}, {-1, -2}}
	inner2 := createLinkedLoop(verts4)

	polygon := LinkedGeoPolygon{}
	addLinkedLoop(&polygon, inner2)
	addLinkedLoop(&polygon, inner)
	addLinkedLoop(&polygon, outer)
	addLinkedLoop(&polygon, outer2)

	err := normalizeMultiPolygon(&polygon)
	if err != E_SUCCESS {
		t.Errorf("Expected success, got %v", err)
	}

	if countLinkedPolygons(&polygon) != 2 {
		t.Error("Polygon count correct")
	}
	if countLinkedLoops(&polygon) != 2 {
		t.Error("Loop count on first polygon correct")
	}
	if countLinkedCoords(polygon.First) != 4 {
		t.Error("Got expected outer loop")
	}
	if countLinkedCoords(polygon.First.Next) != 3 {
		t.Error("Got expected inner loop")
	}
	if countLinkedLoops(polygon.Next) != 2 {
		t.Error("Loop count on second polygon correct")
	}
	if countLinkedCoords(polygon.Next.First) != 4 {
		t.Error("Got expected outer loop")
	}
	if countLinkedCoords(polygon.Next.First.Next) != 3 {
		t.Error("Got expected inner loop")
	}

	destroyLinkedMultiPolygon(&polygon)
}

func TestNormalizeMultiPolygonNestedDonuts(t *testing.T) {
	t.Parallel()

	verts := []LatLng{{0.2, 0.2}, {0.2, -0.2}, {-0.2, -0.2}, {-0.2, 0.2}}
	outer := createLinkedLoop(verts)

	verts2 := []LatLng{{0.1, 0.1}, {-0.1, 0.1}, {-0.1, -0.1}, {0.1, -0.1}}
	inner := createLinkedLoop(verts2)

	verts3 := []LatLng{{0.6, 0.6}, {0.6, -0.6}, {-0.6, -0.6}, {-0.6, 0.6}}
	outerBig := createLinkedLoop(verts3)

	verts4 := []LatLng{{0.5, 0.5}, {-0.5, 0.5}, {-0.5, -0.5}, {0.5, -0.5}}
	innerBig := createLinkedLoop(verts4)

	polygon := LinkedGeoPolygon{}
	addLinkedLoop(&polygon, inner)
	addLinkedLoop(&polygon, outerBig)
	addLinkedLoop(&polygon, innerBig)
	addLinkedLoop(&polygon, outer)

	err := normalizeMultiPolygon(&polygon)
	if err != E_SUCCESS {
		t.Errorf("Expected success, got %v", err)
	}

	if countLinkedPolygons(&polygon) != 2 {
		t.Error("Polygon count correct")
	}
	if countLinkedLoops(&polygon) != 2 {
		t.Error("Loop count on first polygon correct")
	}
	if polygon.First != outerBig {
		t.Error("Got expected outer loop")
	}
	if polygon.First.Next != innerBig {
		t.Error("Got expected inner loop")
	}
	if countLinkedLoops(polygon.Next) != 2 {
		t.Error("Loop count on second polygon correct")
	}
	if polygon.Next.First != outer {
		t.Error("Got expected outer loop")
	}
	if polygon.Next.First.Next != inner {
		t.Error("Got expected inner loop")
	}

	destroyLinkedMultiPolygon(&polygon)
}

func TestNormalizeMultiPolygonNoOuterLoops(t *testing.T) {
	t.Parallel()

	verts1 := []LatLng{{0, 0}, {1, 1}, {0, 1}}
	outer1 := createLinkedLoop(verts1)

	verts2 := []LatLng{{2, 2}, {3, 3}, {2, 3}}
	outer2 := createLinkedLoop(verts2)

	polygon := LinkedGeoPolygon{}
	addLinkedLoop(&polygon, outer1)
	addLinkedLoop(&polygon, outer2)

	err := normalizeMultiPolygon(&polygon)
	if err != E_FAILED {
		t.Error("Expected error code returned")
	}

	if countLinkedPolygons(&polygon) != 1 {
		t.Error("Polygon count correct")
	}
	if countLinkedLoops(&polygon) != 0 {
		t.Error("Loop count as expected with invalid input")
	}

	destroyLinkedMultiPolygon(&polygon)
}

func TestNormalizeMultiPolygonAlreadyNormalized(t *testing.T) {
	t.Parallel()

	verts1 := []LatLng{{0, 0}, {0, 1}, {1, 1}}
	outer1 := createLinkedLoop(verts1)

	verts2 := []LatLng{{2, 2}, {2, 3}, {3, 3}}
	outer2 := createLinkedLoop(verts2)

	polygon := LinkedGeoPolygon{}
	addLinkedLoop(&polygon, outer1)
	next := addNewLinkedPolygon(&polygon)
	addLinkedLoop(next, outer2)

	// Should be a no-op
	err := normalizeMultiPolygon(&polygon)
	if err != E_FAILED {
		t.Error("Expected error code returned")
	}

	if countLinkedPolygons(&polygon) != 2 {
		t.Error("Polygon count correct")
	}
	if countLinkedLoops(&polygon) != 1 {
		t.Error("Loop count on first polygon correct")
	}
	if polygon.First != outer1 {
		t.Error("Got expected outer loop")
	}
	if countLinkedLoops(polygon.Next) != 1 {
		t.Error("Loop count on second polygon correct")
	}
	if polygon.Next.First != outer2 {
		t.Error("Got expected outer loop")
	}

	destroyLinkedMultiPolygon(&polygon)
}

func TestNormalizeMultiPolygonUnassignedHole(t *testing.T) {
	t.Parallel()

	verts := []LatLng{{0, 0}, {0, 1}, {1, 1}, {1, 0}}
	outer := createLinkedLoop(verts)

	verts2 := []LatLng{{2, 2}, {3, 3}, {2, 3}}
	inner := createLinkedLoop(verts2)

	polygon := LinkedGeoPolygon{}
	addLinkedLoop(&polygon, inner)
	addLinkedLoop(&polygon, outer)

	err := normalizeMultiPolygon(&polygon)
	if err != E_FAILED {
		t.Error("Expected error code returned")
	}

	destroyLinkedMultiPolygon(&polygon)
}

func TestLineCrossesLine(t *testing.T) {
	t.Parallel()

	lines1 := []LatLng{{0, 0}, {1, 1}, {0, 1}, {1, 0}}
	result := lineCrossesLine(&lines1[0], &lines1[1], &lines1[2], &lines1[3])
	if !result {
		t.Error("diagonal intersection")
	}

	lines2 := []LatLng{{1, 1}, {0, 0}, {1, 0}, {0, 1}}
	result = lineCrossesLine(&lines2[0], &lines2[1], &lines2[2], &lines2[3])
	if !result {
		t.Error("diagonal intersection, reverse vertexes")
	}

	lines3 := []LatLng{{0.5, 0}, {0.5, 1}, {0, 0.5}, {1, 0.5}}
	result = lineCrossesLine(&lines3[0], &lines3[1], &lines3[2], &lines3[3])
	if !result {
		t.Error("horizontal/vertical intersection")
	}

	lines4 := []LatLng{{0.5, 1}, {0.5, 0}, {1, 0.5}, {0, 0.5}}
	result = lineCrossesLine(&lines4[0], &lines4[1], &lines4[2], &lines4[3])
	if !result {
		t.Error("horizontal/vertical intersection, reverse vertexes")
	}

	lines5 := []LatLng{{0, 0}, {0.4, 0.4}, {0, 1}, {1, 0}}
	result = lineCrossesLine(&lines5[0], &lines5[1], &lines5[2], &lines5[3])
	if result {
		t.Error("diagonal non-intersection, below")
	}

	lines6 := []LatLng{{0.6, 0.6}, {1, 1}, {0, 1}, {1, 0}}
	result = lineCrossesLine(&lines6[0], &lines6[1], &lines6[2], &lines6[3])
	if result {
		t.Error("diagonal non-intersection, above")
	}

	lines7 := []LatLng{{0.5, 0}, {0.5, 1}, {0, 0.5}, {0.4, 0.5}}
	result = lineCrossesLine(&lines7[0], &lines7[1], &lines7[2], &lines7[3])
	if result {
		t.Error("horizontal/vertical non-intersection, below")
	}

	lines8 := []LatLng{{0.5, 0}, {0.5, 1}, {0.6, 0.5}, {1, 0.5}}
	result = lineCrossesLine(&lines8[0], &lines8[1], &lines8[2], &lines8[3])
	if result {
		t.Error("horizontal/vertical non-intersection, above")
	}

	lines9 := []LatLng{{0.5, 0}, {0.5, 0.4}, {0, 0.5}, {1, 0.5}}
	result = lineCrossesLine(&lines9[0], &lines9[1], &lines9[2], &lines9[3])
	if result {
		t.Error("horizontal/vertical non-intersection, left")
	}

	lines10 := []LatLng{{0.5, 0.6}, {0.5, 1}, {0, 0.5}, {1, 0.5}}
	result = lineCrossesLine(&lines10[0], &lines10[1], &lines10[2], &lines10[3])
	if result {
		t.Error("horizontal/vertical non-intersection, right")
	}
}

func TestCellBoundaryInsidePolygonInside(t *testing.T) {
	t.Parallel()

	verts := []LatLng{{0, 0}, {0, 1}, {1, 1}, {1, 0}}
	geoloop := GeoLoop(verts)
	polygon := GeoPolygon{GeoLoop: geoloop}

	bboxes := make([]BBox, 1)
	bboxesFromGeoPolygon(&polygon, bboxes)

	boundary := CellBoundary{
		NumVerts: 4,
		Verts:    []LatLng{{0.6, 0.6}, {0.6, 0.4}, {0.4, 0.4}, {0.4, 0.6}},
	}
	boundaryBBox := BBox{0.6, 0.4, 0.6, 0.4}

	result := cellBoundaryInsidePolygon(polygon, bboxes, &boundary, &boundaryBBox)
	if !result {
		t.Error("simple containment is inside")
	}
}

func TestCellBoundaryInsidePolygonInsideTransmeridianWest(t *testing.T) {
	t.Parallel()

	verts := []LatLng{
		{0, math.Pi - 0.5},
		{0, -math.Pi + 0.5},
		{1, -math.Pi + 0.5},
		{1, math.Pi - 0.5},
	}
	geoloop := GeoLoop(verts)
	polygon := GeoPolygon{GeoLoop: geoloop}

	bboxes := make([]BBox, 1)
	bboxesFromGeoPolygon(&polygon, bboxes)

	boundary := CellBoundary{
		NumVerts: 4,
		Verts: []LatLng{
			{0.6, math.Pi - 0.1},
			{0.6, math.Pi - 0.2},
			{0.4, math.Pi - 0.2},
			{0.4, math.Pi - 0.1},
		},
	}
	boundaryBBox := BBox{0.6, 0.4, 0.6, 0.4}

	result := cellBoundaryInsidePolygon(polygon, bboxes, &boundary, &boundaryBBox)
	if !result {
		t.Error("simple containment is inside, west side of transmeridian")
	}
}

func TestCellBoundaryInsidePolygonInsideTransmeridianEast(t *testing.T) {
	t.Parallel()

	verts := []LatLng{
		{0, math.Pi - 0.5},
		{0, -math.Pi + 0.5},
		{1, -math.Pi + 0.5},
		{1, math.Pi - 0.5},
	}
	geoloop := GeoLoop(verts)
	polygon := GeoPolygon{GeoLoop: geoloop}

	bboxes := make([]BBox, 1)
	bboxesFromGeoPolygon(&polygon, bboxes)

	boundary := CellBoundary{
		NumVerts: 4,
		Verts: []LatLng{
			{0.6, -math.Pi + 0.4},
			{0.6, -math.Pi + 0.2},
			{0.4, -math.Pi + 0.2},
			{0.4, -math.Pi + 0.4},
		},
	}
	boundaryBBox := BBox{0.6, 0.4, 0.6, 0.4}

	result := cellBoundaryInsidePolygon(polygon, bboxes, &boundary, &boundaryBBox)
	if !result {
		t.Error("simple containment is inside, east side of transmeridian")
	}
}

func TestCellBoundaryInsidePolygonInsideWithHole(t *testing.T) {
	t.Parallel()

	verts := []LatLng{{0, 0}, {0, 1}, {1, 1}, {1, 0}}
	geoloop := GeoLoop(verts)

	holeVerts := []LatLng{{0.3, 0.3}, {0.3, 0.1}, {0.1, 0.1}, {0.1, 0.3}}
	holeGeoLoop := GeoLoop(holeVerts)

	polygon := GeoPolygon{
		GeoLoop: geoloop,
		Holes:   []GeoLoop{holeGeoLoop},
	}

	bboxes := make([]BBox, 2)
	bboxesFromGeoPolygon(&polygon, bboxes)

	boundary := CellBoundary{
		NumVerts: 4,
		Verts:    []LatLng{{0.6, 0.6}, {0.6, 0.4}, {0.4, 0.4}, {0.4, 0.6}},
	}
	boundaryBBox := BBox{0.6, 0.4, 0.6, 0.4}

	result := cellBoundaryInsidePolygon(polygon, bboxes, &boundary, &boundaryBBox)
	if !result {
		t.Error("simple containment is inside, with hole")
	}
}

func TestCellBoundaryInsidePolygonNotInside(t *testing.T) {
	t.Parallel()

	verts := []LatLng{{0, 0}, {0, 1}, {1, 1}, {1, 0}}
	geoloop := GeoLoop(verts)
	polygon := GeoPolygon{GeoLoop: geoloop}

	bboxes := make([]BBox, 1)
	bboxesFromGeoPolygon(&polygon, bboxes)

	boundary := CellBoundary{
		NumVerts: 4,
		Verts:    []LatLng{{1.6, 1.6}, {1.6, 1.4}, {1.4, 1.4}, {1.4, 1.6}},
	}
	boundaryBBox := BBox{1.6, 1.4, 1.6, 1.4}

	result := cellBoundaryInsidePolygon(polygon, bboxes, &boundary, &boundaryBBox)
	if result {
		t.Error("fully outside is not inside")
	}
}

func TestCellBoundaryInsidePolygonNotInsideIntersect(t *testing.T) {
	t.Parallel()

	verts := []LatLng{{0, 0}, {0, 1}, {1, 1}, {1, 0}}
	geoloop := GeoLoop(verts)
	polygon := GeoPolygon{GeoLoop: geoloop}

	bboxes := make([]BBox, 1)
	bboxesFromGeoPolygon(&polygon, bboxes)

	boundary := CellBoundary{
		NumVerts: 4,
		Verts:    []LatLng{{0.6, 0.6}, {1.6, 0.4}, {0.4, 0.4}, {0.4, 0.6}},
	}
	boundaryBBox := BBox{1.6, 0.4, 0.6, 0.4}

	result := cellBoundaryInsidePolygon(polygon, bboxes, &boundary, &boundaryBBox)
	if result {
		t.Error("intersecting polygon is not inside")
	}
}

func TestCellBoundaryInsidePolygonNotInsideIntersectHole(t *testing.T) {
	t.Parallel()

	verts := []LatLng{{0, 0}, {0, 1}, {1, 1}, {1, 0}}
	geoloop := GeoLoop(verts)

	holeVerts := []LatLng{{0.3, 0.3}, {0.5, 0.5}, {0.1, 0.1}, {0.1, 0.3}}
	holeGeoLoop := GeoLoop(holeVerts)

	polygon := GeoPolygon{
		GeoLoop: geoloop,
		Holes:   []GeoLoop{holeGeoLoop},
	}

	bboxes := make([]BBox, 2)
	bboxesFromGeoPolygon(&polygon, bboxes)

	boundary := CellBoundary{
		NumVerts: 4,
		Verts:    []LatLng{{0.6, 0.6}, {0.6, 0.4}, {0.4, 0.4}, {0.4, 0.6}},
	}
	boundaryBBox := BBox{0.6, 0.4, 0.6, 0.4}

	result := cellBoundaryInsidePolygon(polygon, bboxes, &boundary, &boundaryBBox)
	if result {
		t.Error("not inside with hole intersection")
	}
}

func TestCellBoundaryInsidePolygonNotInsideWithinHole(t *testing.T) {
	t.Parallel()

	verts := []LatLng{{0, 0}, {0, 1}, {1, 1}, {1, 0}}
	geoloop := GeoLoop(verts)

	holeVerts := []LatLng{{0.9, 0.9}, {0.9, 0.1}, {0.1, 0.1}, {0.1, 0.9}}
	holeGeoLoop := GeoLoop(holeVerts)

	polygon := GeoPolygon{
		GeoLoop: geoloop,
		Holes:   []GeoLoop{holeGeoLoop},
	}

	bboxes := make([]BBox, 2)
	bboxesFromGeoPolygon(&polygon, bboxes)

	boundary := CellBoundary{
		NumVerts: 4,
		Verts:    []LatLng{{0.6, 0.6}, {0.6, 0.4}, {0.4, 0.4}, {0.4, 0.6}},
	}
	boundaryBBox := BBox{0.6, 0.4, 0.6, 0.4}

	result := cellBoundaryInsidePolygon(polygon, bboxes, &boundary, &boundaryBBox)
	if result {
		t.Error("not inside when within hole")
	}
}

func TestCellBoundaryInsidePolygonNotInsideContains(t *testing.T) {
	t.Parallel()

	verts := []LatLng{{0.6, 0.6}, {0.6, 0.4}, {0.4, 0.4}, {0.4, 0.6}}
	geoloop := GeoLoop(verts)
	polygon := GeoPolygon{GeoLoop: geoloop}

	bboxes := make([]BBox, 1)
	bboxesFromGeoPolygon(&polygon, bboxes)

	boundary := CellBoundary{
		NumVerts: 4,
		Verts:    []LatLng{{0, 0}, {0, 1}, {1, 1}, {1, 0}},
	}
	boundaryBBox := BBox{0, 1, 0, 1}

	result := cellBoundaryInsidePolygon(polygon, bboxes, &boundary, &boundaryBBox)
	if result {
		t.Error("not inside when it contains outer")
	}
}

func TestCellBoundaryInsidePolygonNotInsideContainsHole(t *testing.T) {
	t.Parallel()

	verts := []LatLng{{0, 0}, {0, 1}, {1, 1}, {1, 0}}
	geoloop := GeoLoop(verts)

	holeVerts := []LatLng{{0.6, 0.6}, {0.6, 0.4}, {0.4, 0.4}, {0.4, 0.6}}
	holeGeoLoop := GeoLoop(holeVerts)

	polygon := GeoPolygon{
		GeoLoop: geoloop,
		Holes:   []GeoLoop{holeGeoLoop},
	}

	bboxes := make([]BBox, 2)
	bboxesFromGeoPolygon(&polygon, bboxes)

	boundary := CellBoundary{
		NumVerts: 4,
		Verts:    []LatLng{{0.9, 0.9}, {0.9, 0.1}, {0.1, 0.1}, {0.1, 0.9}},
	}
	boundaryBBox := BBox{0.9, 0.1, 0.9, 0.1}

	result := cellBoundaryInsidePolygon(polygon, bboxes, &boundary, &boundaryBBox)
	if result {
		t.Error("not inside when it contains hole")
	}
}
