// Tests ported from testBBoxInternal.c
package h3

import (
	"math"
	"testing"
)

func assertBBoxFromGeoLoop(t *testing.T, geoloop []LatLng, expected *bbox, inside *LatLng, outside *LatLng) {
	t.Helper()
	var result bbox

	bboxFromGeoLoop(geoloop, &result)

	if !bboxEquals(&result, expected) {
		t.Errorf("Got unexpected bbox: got %+v, expected %+v", result, expected)
	}
	if !bboxContains(&result, inside) {
		t.Errorf("Does not contain expected inside point %+v", inside)
	}
	if bboxContains(&result, outside) {
		t.Errorf("Contains expected outside point %+v (should not)", outside)
	}
}

func assertBBox(t *testing.T, bb *bbox, expected *bbox) {
	t.Helper()
	actualNE := LatLng{Lat: bb.North, Lng: bb.East}
	expectedNE := LatLng{Lat: expected.North, Lng: expected.East}
	if !geoAlmostEqual(&actualNE, &expectedNE) {
		t.Errorf("NE corner does not match: got %+v, expected %+v", actualNE, expectedNE)
	}

	actualSW := LatLng{Lat: bb.South, Lng: bb.West}
	expectedSW := LatLng{Lat: expected.South, Lng: expected.West}
	if !geoAlmostEqual(&actualSW, &expectedSW) {
		t.Errorf("SW corner does not match: got %+v, expected %+v", actualSW, expectedSW)
	}
}

func TestPosLatPosLng(t *testing.T) {
	t.Parallel()
	verts := []LatLng{{Rad(0.8), Rad(0.3)}, {Rad(0.7), Rad(0.6)}, {Rad(1.1), Rad(0.7)}, {Rad(1.0), Rad(0.2)}}
	expected := bbox{Rad(1.1), Rad(0.7), Rad(0.7), Rad(0.2)}
	inside := LatLng{Rad(0.9), Rad(0.4)}
	outside := LatLng{Rad(0.0), Rad(0.0)}
	assertBBoxFromGeoLoop(t, verts, &expected, &inside, &outside)
}

func TestNegLatPosLng(t *testing.T) {
	t.Parallel()
	verts := []LatLng{{Rad(-0.3), Rad(0.6)}, {Rad(-0.4), Rad(0.9)}, {Rad(-0.2), Rad(0.8)}, {Rad(-0.1), Rad(0.6)}}
	expected := bbox{Rad(-0.1), Rad(-0.4), Rad(0.9), Rad(0.6)}
	inside := LatLng{Rad(-0.3), Rad(0.8)}
	outside := LatLng{Rad(0.0), Rad(0.0)}
	assertBBoxFromGeoLoop(t, verts, &expected, &inside, &outside)
}

func TestPosLatNegLng(t *testing.T) {
	t.Parallel()
	verts := []LatLng{{Rad(0.7), Rad(-1.4)}, {Rad(0.8), Rad(-0.9)}, {Rad(1.0), Rad(-0.8)}, {Rad(1.1), Rad(-1.3)}}
	expected := bbox{Rad(1.1), Rad(0.7), Rad(-0.8), Rad(-1.4)}
	inside := LatLng{Rad(0.9), Rad(-1.0)}
	outside := LatLng{Rad(0.0), Rad(0.0)}
	assertBBoxFromGeoLoop(t, verts, &expected, &inside, &outside)
}

func TestNegLatNegLng(t *testing.T) {
	t.Parallel()
	verts := []LatLng{
		{Rad(-0.4), Rad(-1.4)}, {Rad(-0.3), Rad(-1.1)}, {Rad(-0.1), Rad(-1.2)}, {Rad(-0.2), Rad(-1.4)}}
	expected := bbox{Rad(-0.1), Rad(-0.4), Rad(-1.1), Rad(-1.4)}
	inside := LatLng{Rad(-0.3), Rad(-1.2)}
	outside := LatLng{Rad(0.0), Rad(0.0)}
	assertBBoxFromGeoLoop(t, verts, &expected, &inside, &outside)
}

func TestAroundZeroZero(t *testing.T) {
	t.Parallel()
	verts := []LatLng{{Rad(0.4), Rad(-0.4)}, {Rad(0.4), Rad(0.4)}, {Rad(-0.4), Rad(0.4)}, {Rad(-0.4), Rad(-0.4)}}
	expected := bbox{Rad(0.4), Rad(-0.4), Rad(0.4), Rad(-0.4)}
	inside := LatLng{Rad(-0.1), Rad(-0.1)}
	outside := LatLng{Rad(1.0), Rad(-1.0)}
	assertBBoxFromGeoLoop(t, verts, &expected, &inside, &outside)
}

func TestTransmeridian(t *testing.T) {
	t.Parallel()
	verts := []LatLng{{Rad(0.4), Pi - Rad(0.1)},
		{Rad(0.4), -Pi + Rad(0.1)},
		{Rad(-0.4), -Pi + Rad(0.1)},
		{Rad(-0.4), Pi - Rad(0.1)}}
	expected := bbox{Rad(0.4), Rad(-0.4), -Pi + Rad(0.1), Pi - Rad(0.1)}
	insideOnMeridian := LatLng{Rad(-0.1), Pi}
	outside := LatLng{Rad(1.0), Pi - Rad(0.5)}
	assertBBoxFromGeoLoop(t, verts, &expected, &insideOnMeridian, &outside)

	westInside := LatLng{Rad(0.1), Pi - Rad(0.05)}
	if !bboxContains(&expected, &westInside) {
		t.Error("Does not contain expected west inside point")
	}
	eastInside := LatLng{Rad(0.1), -Pi + Rad(0.05)}
	if !bboxContains(&expected, &eastInside) {
		t.Error("Does not contain expected east inside point")
	}

	westOutside := LatLng{Rad(0.1), Pi - Rad(0.5)}
	if bboxContains(&expected, &westOutside) {
		t.Error("Contains expected west outside point")
	}
	eastOutside := LatLng{Rad(0.1), -Pi + Rad(0.5)}
	if bboxContains(&expected, &eastOutside) {
		t.Error("Contains expected east outside point")
	}
}

func TestEdgeOnNorthPole(t *testing.T) {
	t.Parallel()
	verts := []LatLng{{PiOver2 - Rad(0.1), Rad(0.1)},
		{PiOver2 - Rad(0.1), Rad(0.8)},
		{PiOver2, Rad(0.8)},
		{PiOver2, Rad(0.1)}}
	expected := bbox{PiOver2, PiOver2 - Rad(0.1), Rad(0.8), Rad(0.1)}
	inside := LatLng{PiOver2 - Rad(0.01), Rad(0.4)}
	outside := LatLng{PiOver2, Rad(0.9)}
	assertBBoxFromGeoLoop(t, verts, &expected, &inside, &outside)
}

func TestEdgeOnSouthPole(t *testing.T) {
	t.Parallel()
	verts := []LatLng{{-PiOver2 + Rad(0.1), Rad(0.1)},
		{-PiOver2 + Rad(0.1), Rad(0.8)},
		{-PiOver2, Rad(0.8)},
		{-PiOver2, Rad(0.1)}}
	expected := bbox{-PiOver2 + Rad(0.1), -PiOver2, Rad(0.8), Rad(0.1)}
	inside := LatLng{-PiOver2 + Rad(0.01), Rad(0.4)}
	outside := LatLng{-PiOver2, Rad(0.9)}
	assertBBoxFromGeoLoop(t, verts, &expected, &inside, &outside)
}

func TestContainsEdges(t *testing.T) {
	t.Parallel()
	bb := bbox{Rad(0.1), Rad(-0.1), Rad(0.2), Rad(-0.2)}
	points := []LatLng{
		{Rad(0.1), Rad(0.2)}, {Rad(0.1), Rad(0.0)}, {Rad(0.1), Rad(-0.2)}, {Rad(0.0), Rad(0.2)},
		{Rad(-0.1), Rad(0.2)}, {Rad(-0.1), Rad(0.0)}, {Rad(-0.1), Rad(-0.2)}, {Rad(0.0), Rad(-0.2)},
	}

	for i, point := range points {
		if !bboxContains(&bb, &point) {
			t.Errorf("Does not contain edge point %d: %+v", i, point)
		}
	}
}

func TestContainsEdgesTransmeridian(t *testing.T) {
	t.Parallel()
	bb := bbox{Rad(0.1), Rad(-0.1), -Pi + Rad(0.2), Pi - Rad(0.2)}
	points := []LatLng{
		{Rad(0.1), -Pi + Rad(0.2)}, {Rad(0.1), Pi}, {Rad(0.1), Pi - Rad(0.2)},
		{Rad(0.0), -Pi + Rad(0.2)}, {Rad(-0.1), -Pi + Rad(0.2)}, {Rad(-0.1), Pi},
		{Rad(-0.1), Pi - Rad(0.2)}, {Rad(0.0), Pi - Rad(0.2)},
	}

	for i, point := range points {
		if !bboxContains(&bb, &point) {
			t.Errorf("Does not contain transmeridian edge point %d: %+v", i, point)
		}
	}
}

func TestBboxOverlapsBBox(t *testing.T) {
	t.Parallel()
	a := bbox{Rad(1.0), Rad(0.0), Rad(1.0), Rad(0.0)}

	b1 := bbox{Rad(1.0), Rad(0.0), Rad(-1.0), Rad(-1.5)}
	if bboxOverlapsBBox(&a, &b1) {
		t.Error("Should have no intersection to the west")
	}
	if bboxOverlapsBBox(&b1, &a) {
		t.Error("Should have no intersection to the west, reverse")
	}

	b2 := bbox{Rad(1.0), Rad(0.0), Rad(2.0), Rad(1.5)}
	if bboxOverlapsBBox(&a, &b2) {
		t.Error("Should have no intersection to the east")
	}
	if bboxOverlapsBBox(&b2, &a) {
		t.Error("Should have no intersection to the east, reverse")
	}

	b3 := bbox{Rad(-1.0), Rad(-1.5), Rad(1.0), Rad(0.0)}
	if bboxOverlapsBBox(&a, &b3) {
		t.Error("Should have no intersection to the south")
	}
	if bboxOverlapsBBox(&b3, &a) {
		t.Error("Should have no intersection to the south, reverse")
	}

	b4 := bbox{Rad(2.0), Rad(1.5), Rad(1.0), Rad(0.0)}
	if bboxOverlapsBBox(&a, &b4) {
		t.Error("Should have no intersection to the north")
	}
	if bboxOverlapsBBox(&b4, &a) {
		t.Error("Should have no intersection to the north, reverse")
	}

	b5 := bbox{Rad(1.0), Rad(0.0), Rad(0.5), Rad(-1.5)}
	if !bboxOverlapsBBox(&a, &b5) {
		t.Error("Should have intersection to the west")
	}

	b6 := bbox{Rad(1.0), Rad(0.0), Rad(2.0), Rad(0.5)}
	if !bboxOverlapsBBox(&a, &b6) {
		t.Error("Should have intersection to the east")
	}

	b7 := bbox{Rad(0.5), Rad(-1.5), Rad(1.0), Rad(0.0)}
	if !bboxOverlapsBBox(&a, &b7) {
		t.Error("Should have intersection to the south")
	}

	b8 := bbox{Rad(2.0), Rad(0.5), Rad(1.0), Rad(0.0)}
	if !bboxOverlapsBBox(&a, &b8) {
		t.Error("Should have intersection to the north")
	}

	b9 := bbox{Rad(1.5), Rad(-0.5), Rad(1.5), Rad(-0.5)}
	if !bboxOverlapsBBox(&a, &b9) {
		t.Error("Should have intersection, b contains a")
	}

	b10 := bbox{Rad(0.5), Rad(0.25), Rad(0.5), Rad(0.25)}
	if !bboxOverlapsBBox(&a, &b10) {
		t.Error("Should have intersection, a contains b")
	}

	b11 := bbox{Rad(1.0), Rad(0.0), Rad(1.0), Rad(0.0)}
	if !bboxOverlapsBBox(&a, &b11) {
		t.Error("Should have intersection, a equals b")
	}
}

func TestBboxOverlapsBBoxTransmeridian(t *testing.T) {
	t.Parallel()
	a := bbox{Rad(1.0), Rad(0.0), -Pi + Rad(0.5), Pi - Rad(0.5)}

	b1 := bbox{Rad(1.0), Rad(0.0), Pi - Rad(0.7), Pi - Rad(0.9)}
	if bboxOverlapsBBox(&a, &b1) {
		t.Error("Should have no intersection to the west")
	}
	if bboxOverlapsBBox(&b1, &a) {
		t.Error("Should have no intersection to the west, reverse")
	}

	b2 := bbox{Rad(1.0), Rad(0.0), -Pi + Rad(0.9), -Pi + Rad(0.7)}
	if bboxOverlapsBBox(&a, &b2) {
		t.Error("Should have no intersection to the east")
	}
	if bboxOverlapsBBox(&b2, &a) {
		t.Error("Should have no intersection to the east")
	}

	b3 := bbox{Rad(1.0), Rad(0.0), Pi - Rad(0.4), Pi - Rad(0.9)}
	if !bboxOverlapsBBox(&a, &b3) {
		t.Error("Should have intersection to the west")
	}
	if !bboxOverlapsBBox(&b3, &a) {
		t.Error("Should have intersection to the west, reverse")
	}

	b4 := bbox{Rad(1.0), Rad(0.0), -Pi + Rad(0.9), -Pi + Rad(0.4)}
	if !bboxOverlapsBBox(&a, &b4) {
		t.Error("Should have intersection to the east")
	}
	if !bboxOverlapsBBox(&b4, &a) {
		t.Error("Should have intersection to the east, reverse")
	}

	b5 := bbox{Rad(1.0), Rad(0.0), -Pi + Rad(0.4), Pi - Rad(0.4)}
	if !bboxOverlapsBBox(&a, &b5) {
		t.Error("Should have intersection, a contains b")
	}
	if !bboxOverlapsBBox(&b5, &a) {
		t.Error("Should have intersection, a contains b, reverse")
	}

	b6 := bbox{Rad(1.0), Rad(0.0), -Pi + Rad(0.6), Pi - Rad(0.6)}
	if !bboxOverlapsBBox(&a, &b6) {
		t.Error("Should have intersection, b contains a")
	}
	if !bboxOverlapsBBox(&b6, &a) {
		t.Error("Should have intersection, b contains a, reverse")
	}

	b7 := bbox{Rad(1.0), Rad(0.0), -Pi + Rad(0.5), Pi - Rad(0.5)}
	if !bboxOverlapsBBox(&a, &b7) {
		t.Error("Should have intersection, a equals b")
	}

	b8 := bbox{Rad(1.0), Rad(0.0), -Pi + Rad(0.9), Pi - Rad(0.4)}
	if !bboxOverlapsBBox(&a, &b8) {
		t.Error("Should have intersection, transmeridian to the east")
	}
	if !bboxOverlapsBBox(&b8, &a) {
		t.Error("Should have intersection, transmeridian to the east, reverse")
	}

	b9 := bbox{Rad(1.0), Rad(0.0), -Pi + Rad(0.4), Pi - Rad(0.9)}
	if !bboxOverlapsBBox(&a, &b9) {
		t.Error("Should have intersection, transmeridian to the west")
	}
	if !bboxOverlapsBBox(&b9, &a) {
		t.Error("Should have intersection, transmeridian to the west, reverse")
	}
}

func TestBboxCenterBasicQuandrants(t *testing.T) {
	t.Parallel()

	bbox1 := bbox{Rad(1.0), Rad(0.8), Rad(1.0), Rad(0.8)}
	expected1 := LatLng{Rad(0.9), Rad(0.9)}
	center := bboxCenter(&bbox1)
	if !geoAlmostEqual(&center, &expected1) {
		t.Errorf("pos/pos not as expected: got %+v, expected %+v", center, expected1)
	}

	bbox2 := bbox{Rad(-0.8), Rad(-1.0), Rad(1.0), Rad(0.8)}
	expected2 := LatLng{Rad(-0.9), Rad(0.9)}
	center = bboxCenter(&bbox2)
	if !geoAlmostEqual(&center, &expected2) {
		t.Errorf("neg/pos not as expected: got %+v, expected %+v", center, expected2)
	}

	bbox3 := bbox{Rad(1.0), Rad(0.8), Rad(-0.8), Rad(-1.0)}
	expected3 := LatLng{Rad(0.9), Rad(-0.9)}
	center = bboxCenter(&bbox3)
	if !geoAlmostEqual(&center, &expected3) {
		t.Errorf("pos/neg not as expected: got %+v, expected %+v", center, expected3)
	}

	bbox4 := bbox{Rad(-0.8), Rad(-1.0), Rad(-0.8), Rad(-1.0)}
	expected4 := LatLng{Rad(-0.9), Rad(-0.9)}
	center = bboxCenter(&bbox4)
	if !geoAlmostEqual(&center, &expected4) {
		t.Errorf("neg/neg not as expected: got %+v, expected %+v", center, expected4)
	}

	bbox5 := bbox{Rad(0.8), Rad(-0.8), Rad(1.0), Rad(-1.0)}
	expected5 := LatLng{Rad(0.0), Rad(0.0)}
	center = bboxCenter(&bbox5)
	if !geoAlmostEqual(&center, &expected5) {
		t.Errorf("around origin not as expected: got %+v, expected %+v", center, expected5)
	}
}

func TestBboxCenterTransmeridian(t *testing.T) {
	t.Parallel()

	bbox1 := bbox{Rad(1.0), Rad(0.8), -Pi + Rad(0.3), Pi - Rad(0.1)}
	expected1 := LatLng{Rad(0.9), -Pi + Rad(0.1)}
	center := bboxCenter(&bbox1)

	if !geoAlmostEqual(&center, &expected1) {
		t.Errorf("skew east not as expected: got %+v, expected %+v", center, expected1)
	}

	bbox2 := bbox{Rad(1.0), Rad(0.8), -Pi + Rad(0.1), Pi - Rad(0.3)}
	expected2 := LatLng{Rad(0.9), Pi - Rad(0.1)}
	center = bboxCenter(&bbox2)
	if !geoAlmostEqual(&center, &expected2) {
		t.Errorf("skew west not as expected: got %+v, expected %+v", center, expected2)
	}

	bbox3 := bbox{Rad(1.0), Rad(0.8), -Pi + Rad(0.1), Pi - Rad(0.1)}
	expected3 := LatLng{Rad(0.9), Pi}
	center = bboxCenter(&bbox3)
	if !geoAlmostEqual(&center, &expected3) {
		t.Errorf("on antimeridian not as expected: got %+v, expected %+v", center, expected3)
	}
}

func TestBboxIsTransmeridian(t *testing.T) {
	t.Parallel()
	bboxNormal := bbox{Rad(1.0), Rad(0.8), Rad(1.0), Rad(0.8)}
	if bboxIsTransmeridian(&bboxNormal) {
		t.Error("Normal bbox should not be transmeridian")
	}

	bboxTransmeridian := bbox{Rad(1.0), Rad(0.8), -Pi + Rad(0.3), Pi - Rad(0.1)}
	if !bboxIsTransmeridian(&bboxTransmeridian) {
		t.Error("Transmeridian bbox should be transmeridian")
	}
}

func TestBboxEquals(t *testing.T) {
	t.Parallel()
	bb := bbox{Rad(1.0), Rad(0.0), Rad(1.0), Rad(0.0)}
	north := bb
	north.North += Rad(0.1)
	south := bb
	south.South += Rad(0.1)
	east := bb
	east.East += Rad(0.1)
	west := bb
	west.West += Rad(0.1)

	if !bboxEquals(&bb, &bb) {
		t.Error("Should equal self")
	}
	if bboxEquals(&bb, &north) {
		t.Error("Should not equal different north")
	}
	if bboxEquals(&bb, &south) {
		t.Error("Should not equal different south")
	}
	if bboxEquals(&bb, &east) {
		t.Error("Should not equal different east")
	}
	if bboxEquals(&bb, &west) {
		t.Error("Should not equal different west")
	}
}

func TestBboxHexEstimate_invalidRes(t *testing.T) {
	t.Parallel()
	var numHexagons int64
	bb := bbox{Rad(1.0), Rad(0.0), Rad(1.0), Rad(0.0)}
	err := bboxHexEstimate(&bb, -1, &numHexagons)
	if err != eResDomain {
		t.Errorf("bboxHexEstimate of invalid resolution should fail with eResDomain, got %v", err)
	}
}

func TestBboxHexEstimate_ratio(t *testing.T) {
	t.Parallel()
	bbox1 := bbox{Rad(0.82294), Rad(0.82273), Rad(0.131671), Rad(0.131668)}
	bbox2 := bbox{Rad(0.131671), Rad(0.131668), Rad(0.82294), Rad(0.82273)}
	var numHexagons1, numHexagons2 int64

	if err := bboxHexEstimate(&bbox1, 15, &numHexagons1); err != eSuccess {
		t.Fatalf("bboxHexEstimate failed: %v", err)
	}
	if err := bboxHexEstimate(&bbox2, 15, &numHexagons2); err != eSuccess {
		t.Fatalf("bboxHexEstimate failed: %v", err)
	}

	diffPercentage := math.Abs(1.0 - float64(numHexagons1)/float64(numHexagons2))

	// numHexagons1 and numHexagons2 cannot be exactly equal because the
	// diameter of the two bboxes is not exactly the same (it's calculated
	// using greatCircleDistanceKm)
	if diffPercentage >= 0.03 {
		t.Errorf("Should be true for bounding boxes with (almost) the same diameter and side ratio, got diffPercentage %f", diffPercentage)
	}
}

func TestLineHexEstimate_invalidRes(t *testing.T) {
	t.Parallel()
	var numHexagons int64
	origin := LatLng{Rad(0.0), Rad(0.0)}
	destination := LatLng{Rad(1.0), Rad(1.0)}
	err := lineHexEstimate(&origin, &destination, -1, &numHexagons)
	if err != eResDomain {
		t.Errorf("lineHexEstimate of invalid resolution should fail with eResDomain, got %v", err)
	}
}

func TestScaleBBox_noop(t *testing.T) {
	t.Parallel()
	bb := bbox{Rad(1.0), Rad(0.0), Rad(1.0), Rad(0.0)}
	expected := bbox{Rad(1.0), Rad(0.0), Rad(1.0), Rad(0.0)}
	scaleBBox(&bb, 1)
	assertBBox(t, &bb, &expected)
}

func TestScaleBBox_basicGrow(t *testing.T) {
	t.Parallel()
	bb := bbox{Rad(1.0), Rad(0.0), Rad(1.0), Rad(0.0)}
	expected := bbox{Rad(1.5), Rad(-0.5), Rad(1.5), Rad(-0.5)}
	scaleBBox(&bb, 2)
	assertBBox(t, &bb, &expected)
}

func TestScaleBBox_basicShrink(t *testing.T) {
	t.Parallel()
	bb := bbox{Rad(1.0), Rad(0.0), Rad(1.0), Rad(0.0)}
	expected := bbox{Rad(0.75), Rad(0.25), Rad(0.75), Rad(0.25)}
	scaleBBox(&bb, 0.5)
	assertBBox(t, &bb, &expected)
}

func TestScaleBBox_clampNorthSouth(t *testing.T) {
	t.Parallel()
	bb := bbox{PiOver2 * Rad(0.9), -PiOver2 * Rad(0.9), Rad(1.0), Rad(0.0)}
	expected := bbox{PiOver2, -PiOver2, Rad(1.5), Rad(-0.5)}
	scaleBBox(&bb, 2)
	assertBBox(t, &bb, &expected)
}

func TestScaleBBox_clampEastPos(t *testing.T) {
	t.Parallel()
	bb := bbox{Rad(1.0), Rad(0.0), Pi - Rad(0.1), Pi - Rad(1.1)}
	expected := bbox{Rad(1.5), Rad(-0.5), -Pi + Rad(0.4), Pi - Rad(1.6)}
	scaleBBox(&bb, 2)
	assertBBox(t, &bb, &expected)
}

func TestScaleBBox_clampEastNeg(t *testing.T) {
	t.Parallel()
	bb := bbox{Rad(1.5), Rad(-0.5), -Pi + Rad(0.4), Pi - Rad(1.6)}
	expected := bbox{Rad(1.0), Rad(0.0), Pi - Rad(0.1), Pi - Rad(1.1)}
	scaleBBox(&bb, 0.5)
	assertBBox(t, &bb, &expected)
}

func TestScaleBBox_clampWestPos(t *testing.T) {
	t.Parallel()
	bb := bbox{Rad(1.0), Rad(0.0), -Pi + Rad(0.9), Pi - Rad(0.1)}
	expected := bbox{Rad(0.75), Rad(0.25), -Pi + Rad(0.65), -Pi + Rad(0.15)}
	scaleBBox(&bb, 0.5)
	assertBBox(t, &bb, &expected)
}

func TestScaleBBox_clampWestNeg(t *testing.T) {
	t.Parallel()
	bb := bbox{Rad(0.75), Rad(0.25), -Pi + Rad(0.65), -Pi + Rad(0.15)}
	expected := bbox{Rad(1.0), Rad(0.0), -Pi + Rad(0.9), Pi - Rad(0.1)}
	scaleBBox(&bb, 2)
	assertBBox(t, &bb, &expected)
}
