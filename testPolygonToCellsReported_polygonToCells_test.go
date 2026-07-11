// Tests ported from testPolygonToCellsReported.c
package h3

import (
	"testing"
)

func TestPolygonToCells_entireWorld(t *testing.T) {
	t.Parallel()

	// Test for entire world coverage using two polygons
	// https://github.com/uber/h3-js/issues/76#issuecomment-561204505
	worldVerts1 := []LatLng{
		{Lat: Angle(-mPi2), Lng: Angle(-mPi)},
		{Lat: Angle(mPi2), Lng: Angle(-mPi)},
		{Lat: Angle(mPi2), Lng: Angle(0)},
		{Lat: Angle(-mPi2), Lng: Angle(0)},
	}
	worldGeoPolygon1 := GeoPolygon{GeoLoop: worldVerts1, Holes: nil}

	worldVerts2 := []LatLng{
		{Lat: Angle(-mPi2), Lng: Angle(0)},
		{Lat: Angle(mPi2), Lng: Angle(0)},
		{Lat: Angle(mPi2), Lng: Angle(mPi)},
		{Lat: Angle(-mPi2), Lng: Angle(mPi)},
	}
	worldGeoPolygon2 := GeoPolygon{GeoLoop: worldVerts2, Holes: nil}

	for res := int32(0); res < 3; res++ {
		// Process first polygon
		var polygonToCellsSize1 int64
		err := maxPolygonToCellsSize(&worldGeoPolygon1, res, 0, &polygonToCellsSize1)
		if err != eSuccess {
			t.Fatalf("maxPolygonToCellsSize failed for polygon 1 at res %d: %v", res, err)
		}
		polygonToCellsOut1 := make([]h3Index, polygonToCellsSize1)

		err = polygonToCells(&worldGeoPolygon1, res, 0, polygonToCellsOut1)
		if err != eSuccess {
			t.Fatalf("polygonToCells failed for polygon 1 at res %d: %v", res, err)
		}
		actualNumIndexes1 := countNonNullIndexes(polygonToCellsOut1)

		// Process second polygon
		var polygonToCellsSize2 int64
		err = maxPolygonToCellsSize(&worldGeoPolygon2, res, 0, &polygonToCellsSize2)
		if err != eSuccess {
			t.Fatalf("maxPolygonToCellsSize failed for polygon 2 at res %d: %v", res, err)
		}
		polygonToCellsOut2 := make([]h3Index, polygonToCellsSize2)

		err = polygonToCells(&worldGeoPolygon2, res, 0, polygonToCellsOut2)
		if err != eSuccess {
			t.Fatalf("polygonToCells failed for polygon 2 at res %d: %v", res, err)
		}
		actualNumIndexes2 := countNonNullIndexes(polygonToCellsOut2)

		// Get expected total world cells
		expectedTotalWorld, err := getNumCells(res)
		if err != eSuccess {
			t.Fatalf("getNumCells failed for res %d: %v", res, err)
		}

		if actualNumIndexes1+actualNumIndexes2 != expectedTotalWorld {
			t.Errorf("Total world cells mismatch at res %d: got %d+%d=%d, expected %d",
				res, actualNumIndexes1, actualNumIndexes2, actualNumIndexes1+actualNumIndexes2, expectedTotalWorld)
		}

		// Check that sets are disjoint
		indexSet := make(map[h3Index]bool)
		for _, idx := range polygonToCellsOut1 {
			if idx != h3Null {
				indexSet[idx] = true
			}
		}

		for _, idx := range polygonToCellsOut2 {
			if idx != h3Null {
				if indexSet[idx] {
					t.Errorf("Index 0x%x found in both polygon results at res %d - sets should be disjoint", idx, res)
				}
			}
		}
	}
}

func TestPolygonToCells_h3js_67(t *testing.T) {
	t.Parallel()

	// https://github.com/uber/h3-js/issues/67
	east := degsToRads(-56.25)
	north := degsToRads(-33.13755119234615)
	south := degsToRads(-34.30714385628804)
	west := degsToRads(-57.65625)

	testVerts := []LatLng{
		{Lat: Angle(north), Lng: Angle(east)},
		{Lat: Angle(south), Lng: Angle(east)},
		{Lat: Angle(south), Lng: Angle(west)},
		{Lat: Angle(north), Lng: Angle(west)},
	}
	testPolygon := GeoPolygon{
		GeoLoop: testVerts,
		Holes:   nil,
	}

	res := int32(7)
	var numHexagons int64
	err := maxPolygonToCellsSize(&testPolygon, res, 0, &numHexagons)
	if err != eSuccess {
		t.Fatalf("maxPolygonToCellsSize failed: %v", err)
	}
	hexagons := make([]h3Index, numHexagons)

	err = polygonToCells(&testPolygon, res, 0, hexagons)
	if err != eSuccess {
		t.Fatalf("polygonToCells failed: %v", err)
	}
	actualNumIndexes := countNonNullIndexes(hexagons)

	if actualNumIndexes != 4499 {
		t.Errorf("Expected 4499 polygonToCells, got %d (h3-js#67)", actualNumIndexes)
	}
}

func TestPolygonToCells_h3js_67_2nd(t *testing.T) {
	t.Parallel()

	// 2nd test case from h3-js#67
	east := degsToRads(-57.65625)
	north := degsToRads(-34.30714385628804)
	south := degsToRads(-35.4606699514953)
	west := degsToRads(-59.0625)

	testVerts := []LatLng{
		{Lat: Angle(north), Lng: Angle(east)},
		{Lat: Angle(south), Lng: Angle(east)},
		{Lat: Angle(south), Lng: Angle(west)},
		{Lat: Angle(north), Lng: Angle(west)},
	}
	testPolygon := GeoPolygon{
		GeoLoop: testVerts,
		Holes:   nil,
	}

	res := int32(7)
	var numHexagons int64
	err := maxPolygonToCellsSize(&testPolygon, res, 0, &numHexagons)
	if err != eSuccess {
		t.Fatalf("maxPolygonToCellsSize failed: %v", err)
	}
	hexagons := make([]h3Index, numHexagons)

	err = polygonToCells(&testPolygon, res, 0, hexagons)
	if err != eSuccess {
		t.Fatalf("polygonToCells failed: %v", err)
	}
	actualNumIndexes := countNonNullIndexes(hexagons)

	if actualNumIndexes != 4609 {
		t.Errorf("Expected 4609 polygonToCells, got %d (h3-js#67, 2nd case)", actualNumIndexes)
	}
}

func TestPolygonToCells_h3_136(t *testing.T) {
	t.Parallel()

	// https://github.com/uber/h3/issues/136
	testVerts := []LatLng{
		{Lat: Angle(0.10068990369902957), Lng: Angle(0.8920772174196191)},
		{Lat: Angle(0.10032914690616246), Lng: Angle(0.8915914753447348)},
		{Lat: Angle(0.10033349237998787), Lng: Angle(0.8915860128746426)},
		{Lat: Angle(0.10069496685903621), Lng: Angle(0.8920742194546231)},
	}
	testPolygon := GeoPolygon{
		GeoLoop: testVerts,
		Holes:   nil,
	}

	res := int32(13)
	var numHexagons int64
	err := maxPolygonToCellsSize(&testPolygon, res, 0, &numHexagons)
	if err != eSuccess {
		t.Fatalf("maxPolygonToCellsSize failed: %v", err)
	}
	hexagons := make([]h3Index, numHexagons)

	err = polygonToCells(&testPolygon, res, 0, hexagons)
	if err != eSuccess {
		t.Fatalf("polygonToCells failed: %v", err)
	}
	actualNumIndexes := countNonNullIndexes(hexagons)

	if actualNumIndexes != 4353 {
		t.Errorf("Expected 4353 polygonToCells, got %d", actualNumIndexes)
	}
}

func TestPolygonToCells_h3_595(t *testing.T) {
	t.Parallel()

	// https://github.com/uber/h3/issues/595
	// Note: The C code has the second test incorrectly named as h3_136,
	// but this is actually testing issue 595
	center := h3Index(0x85283473fffffff)
	var centerLatLng LatLng
	err := cellToLatLng(center, &centerLatLng)
	if err != eSuccess {
		t.Fatalf("cellToLatLng failed: %v", err)
	}

	// This polygon should include the center cell. The issue here arises
	// when one of the polygon vertexes is to the east of the index center,
	// with exactly the same latitude
	testVerts := []LatLng{
		{Lat: centerLatLng.Lat, Lng: Angle(-2.121207808248113)},
		{Lat: Angle(0.6565301558937859), Lng: Angle(-2.1281107217935986)},
		{Lat: Angle(0.6515463604919347), Lng: Angle(-2.1345342663428695)},
		{Lat: Angle(0.6466583305904194), Lng: Angle(-2.1276313527973842)},
	}

	testPolygon := GeoPolygon{
		GeoLoop: testVerts,
		Holes:   nil,
	}

	res := int32(5)
	var numHexagons int64
	err = maxPolygonToCellsSize(&testPolygon, res, 0, &numHexagons)
	if err != eSuccess {
		t.Fatalf("maxPolygonToCellsSize failed: %v", err)
	}
	hexagons := make([]h3Index, numHexagons)

	err = polygonToCells(&testPolygon, res, 0, hexagons)
	if err != eSuccess {
		t.Fatalf("polygonToCells failed: %v", err)
	}
	actualNumIndexes := countNonNullIndexes(hexagons)

	if actualNumIndexes != 8 {
		t.Errorf("Expected 8 polygonToCells, got %d", actualNumIndexes)
	}
}
