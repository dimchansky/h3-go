// Tests ported from H3 v4.5.0: src/apps/testapps/testLinkedGeoConvert.c.
package h3

import "testing"

func assertSameLatLng(t *testing.T, a, b LatLng) {
	t.Helper()
	if a.Lat != b.Lat {
		t.Error("lat matches")
	}
	if a.Lng != b.Lng {
		t.Error("lng matches")
	}
}

// Expects vertices in the same order in both loops.
func assertSameLoop(t *testing.T, ll *linkedGeoLoop, gl GeoLoop) {
	t.Helper()
	if countLinkedCoords(ll) != int32(len(gl)) {
		t.Error("vert count matches")
	}
	coord := ll.First
	for i := 0; i < len(gl); i++ {
		if coord == nil {
			t.Fatal("coord exists")
		}
		assertSameLatLng(t, coord.Vertex, gl[i])
		coord = coord.Next
	}
}

// Assumes outer loop and holes in same order.
func assertSamePoly(t *testing.T, lp *linkedGeoPolygon, gp *GeoPolygon) {
	t.Helper()
	expectedLoops := 1 + len(gp.Holes)
	if countLinkedLoops(lp) != int32(expectedLoops) {
		t.Error("loop count matches")
	}

	assertSameLoop(t, lp.First, gp.GeoLoop)

	ll := lp.First.Next
	for h := 0; h < len(gp.Holes); h++ {
		assertSameLoop(t, ll, gp.Holes[h])
		ll = ll.Next
	}
}

// Assumes polygons listed in same order.
func assertSameMultiPoly(t *testing.T, linked *linkedGeoPolygon, mpoly *geoMultiPolygon) {
	t.Helper()
	if countLinkedPolygons(linked) != mpoly.NumPolygons {
		t.Error("polygon count matches")
	}

	lp := linked
	for p := int32(0); p < mpoly.NumPolygons; p++ {
		assertSamePoly(t, lp, &mpoly.Polygons[p])
		lp = lp.Next
	}
}

func TestLinkedGeoConvert_geoMultiPolygonToLinkedAndBack(t *testing.T) {
	t.Parallel()

	// Two polygons: one with 1 hole, and one with no holes.
	cells := []h3Index{
		0x8027fffffffffff, 0x802bfffffffffff, 0x804dfffffffffff,
		0x8067fffffffffff, 0x806dfffffffffff, 0x8049fffffffffff,
		0x8055fffffffffff,
	}
	var mpoly, mpoly2 geoMultiPolygon
	var lpoly linkedGeoPolygon

	if err := cellsToMultiPolygon(cells, int64(len(cells)), &mpoly); err != eSuccess {
		t.Fatalf("cellsToMultiPolygon: %v", err)
	}
	if mpoly.NumPolygons != 2 {
		t.Error("has two polygons")
	}

	if err := geoMultiPolygonToLinkedGeoPolygon(&mpoly, &lpoly); err != eSuccess {
		t.Fatalf("geoMultiPolygonToLinkedGeoPolygon: %v", err)
	}
	if err := linkedGeoPolygonToGeoMultiPolygon(&lpoly, &mpoly2); err != eSuccess {
		t.Fatalf("linkedGeoPolygonToGeoMultiPolygon: %v", err)
	}

	assertSameMultiPoly(t, &lpoly, &mpoly)
	assertSameMultiPoly(t, &lpoly, &mpoly2)

	destroyGeoMultiPolygon(&mpoly)
	destroyGeoMultiPolygon(&mpoly2)
	destroyLinkedMultiPolygon(&lpoly)
}

func TestLinkedGeoConvert_linkedToGeoMultiPolygonRejectsTooFewVerts(t *testing.T) {
	t.Parallel()

	v1 := LatLng{}
	v2 := LatLng{Lat: Rad(1)}

	// A loop with only 2 vertices should be rejected
	var poly1 linkedGeoPolygon
	loop1 := addNewLinkedLoop(&poly1)
	addLinkedCoord(loop1, &v1)
	addLinkedCoord(loop1, &v2)

	var mpoly geoMultiPolygon
	if err := linkedGeoPolygonToGeoMultiPolygon(&poly1, &mpoly); err != eFailed {
		t.Errorf("rejects loop with < 3 verts: got %v", err)
	}

	destroyLinkedMultiPolygon(&poly1)
}

func TestLinkedGeoConvert_linkedToGeoMultiPolygonRejectsEmptyPolygon(t *testing.T) {
	t.Parallel()

	// A polygon node with no loops is rejected (not the same as
	// the empty-chain case, which has no next pointer either)
	var poly1 linkedGeoPolygon
	addNewLinkedPolygon(&poly1) // empty node with no loops

	var mpoly geoMultiPolygon
	if err := linkedGeoPolygonToGeoMultiPolygon(&poly1, &mpoly); err != eFailed {
		t.Errorf("rejects empty polygon node: got %v", err)
	}

	destroyLinkedMultiPolygon(&poly1)
}

func TestLinkedGeoConvert_geoToLinkedMultiPolygonRejectsTooFewVerts(t *testing.T) {
	t.Parallel()

	verts := GeoLoop{{}, {Lat: Rad(1)}}
	poly := GeoPolygon{GeoLoop: verts}
	mpoly := geoMultiPolygon{NumPolygons: 1, Polygons: []GeoPolygon{poly}}

	var out linkedGeoPolygon
	if err := geoMultiPolygonToLinkedGeoPolygon(&mpoly, &out); err != eFailed {
		t.Errorf("rejects geoloop with < 3 verts: got %v", err)
	}
}

func TestLinkedGeoConvert_linkedToGeoMultiPolygonEmpty(t *testing.T) {
	t.Parallel()

	var empty linkedGeoPolygon
	var mpoly geoMultiPolygon
	if err := linkedGeoPolygonToGeoMultiPolygon(&empty, &mpoly); err != eSuccess {
		t.Fatalf("linkedGeoPolygonToGeoMultiPolygon: %v", err)
	}
	if mpoly.NumPolygons != 0 {
		t.Error("0 polygons for empty input")
	}
	if mpoly.Polygons != nil {
		t.Error("NULL polygons for empty input")
	}
}

func TestLinkedGeoConvert_geoToLinkedMultiPolygonEmpty(t *testing.T) {
	t.Parallel()

	mpoly := geoMultiPolygon{NumPolygons: 0, Polygons: nil}
	var out linkedGeoPolygon
	if err := geoMultiPolygonToLinkedGeoPolygon(&mpoly, &out); err != eSuccess {
		t.Fatalf("geoMultiPolygonToLinkedGeoPolygon: %v", err)
	}
	if out.First != nil {
		t.Error("empty linked polygon")
	}
	if out.Next != nil {
		t.Error("no next polygon")
	}
}
