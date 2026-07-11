// Tests ported from H3 v4.4.0: src/apps/testapps/testLinkedGeoInternal.c.
package h3

import "testing"

func TestLinkedGeoInternal(t *testing.T) {
	t.Parallel()

	// Fixtures - convert from degrees to radians
	var vertex1, vertex2, vertex3, vertex4 LatLng
	setGeoDegs(&vertex1, 87.372002166, 166.160981117)
	setGeoDegs(&vertex2, 87.370101364, 166.160184306)
	setGeoDegs(&vertex3, 87.369088356, 166.196239997)
	setGeoDegs(&vertex4, 87.369975080, 166.233115768)

	t.Run("createLinkedGeo", func(t *testing.T) {
		t.Parallel()

		polygon := &linkedGeoPolygon{}

		loop := addNewLinkedLoop(polygon)
		if loop == nil {
			t.Fatal("Loop created should not be nil")
		}

		coord := addLinkedCoord(loop, &vertex1)
		if coord == nil {
			t.Fatal("Coord created should not be nil")
		}

		coord = addLinkedCoord(loop, &vertex2)
		if coord == nil {
			t.Fatal("Coord created should not be nil")
		}

		coord = addLinkedCoord(loop, &vertex3)
		if coord == nil {
			t.Fatal("Coord created should not be nil")
		}

		loop = addNewLinkedLoop(polygon)
		if loop == nil {
			t.Fatal("Loop created should not be nil")
		}

		coord = addLinkedCoord(loop, &vertex2)
		if coord == nil {
			t.Fatal("Coord created should not be nil")
		}

		coord = addLinkedCoord(loop, &vertex4)
		if coord == nil {
			t.Fatal("Coord created should not be nil")
		}

		if count := countLinkedPolygons(polygon); count != 1 {
			t.Errorf("Expected polygon count to be 1, got %d", count)
		}

		if count := countLinkedLoops(polygon); count != 2 {
			t.Errorf("Expected loop count to be 2, got %d", count)
		}

		if count := countLinkedCoords(polygon.First); count != 3 {
			t.Errorf("Expected coord count 1 to be 3, got %d", count)
		}

		if count := countLinkedCoords(polygon.Last); count != 2 {
			t.Errorf("Expected coord count 2 to be 2, got %d", count)
		}

		nextPolygon := addNewLinkedPolygon(polygon)
		if nextPolygon == nil {
			t.Fatal("polygon created should not be nil")
		}

		if count := countLinkedPolygons(polygon); count != 2 {
			t.Errorf("Expected polygon count to be 2, got %d", count)
		}

		destroyLinkedMultiPolygon(polygon)
	})
}
