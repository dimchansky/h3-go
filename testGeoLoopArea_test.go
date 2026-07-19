// Tests ported from H3 v4.5.0: src/apps/testapps/testGeoLoopArea.c.

package h3

import (
	"math"
	"testing"
)

func compareGeoLoopArea(t *testing.T, verts []LatLng, targetArea float64) {
	t.Helper()
	const tol = 1e-14
	loop := GeoLoop(verts)

	out, err := geoLoopAreaRads2(loop)
	if err != eSuccess {
		t.Fatalf("geoLoopAreaRads2: %v", err)
	}
	if math.Abs(out-targetArea) >= tol {
		t.Errorf("area should match: got %v want %v", out, targetArea)
	}
}

func TestGeoLoopAreaTriangleBasic(t *testing.T) {
	t.Parallel()

	// GeoLoop representing a triangle covering 1/8 of the globe, with
	// points ordered according to right-hand rule (counter-clockwise).
	// The globe has an area of 4*pi radians, so this 1/8 triangle piece
	// of the globe should have area pi/2.
	verts := []LatLng{
		{Lat: Rad(math.Pi / 2), Lng: Rad(0.0)},
		{Lat: Rad(0.0), Lng: Rad(0.0)},
		{Lat: Rad(0.0), Lng: Rad(math.Pi / 2)},
	}

	compareGeoLoopArea(t, verts, math.Pi/2)
}

func TestGeoLoopAreaTriangleReversed(t *testing.T) {
	t.Parallel()

	// Reversed (clockwise) order: the GeoLoop represents the whole globe
	// minus the triangle above.
	verts := []LatLng{
		{Lat: Rad(0.0), Lng: Rad(math.Pi / 2)},
		{Lat: Rad(0.0), Lng: Rad(0.0)},
		{Lat: Rad(math.Pi / 2), Lng: Rad(0.0)},
	}

	compareGeoLoopArea(t, verts, 7*math.Pi/2)
}

func TestGeoLoopAreaSlice(t *testing.T) {
	t.Parallel()

	// Two 1/8 triangles stitched along the equator: a 1/4 slice of the
	// globe with vertices at the north and south pole.
	verts := []LatLng{
		{Lat: Rad(math.Pi / 2), Lng: Rad(0.0)},
		{Lat: Rad(0.0), Lng: Rad(0.0)},
		{Lat: Rad(-math.Pi / 2), Lng: Rad(0.0)},
		{Lat: Rad(0.0), Lng: Rad(math.Pi / 2)},
	}

	compareGeoLoopArea(t, verts, math.Pi)
}

func TestGeoLoopAreaSliceReversed(t *testing.T) {
	t.Parallel()

	// 3/4 slice of the globe, from north to south pole, formed by
	// reversing the order of points from the example above.
	verts := []LatLng{
		{Lat: Rad(math.Pi / 2), Lng: Rad(0.0)},
		{Lat: Rad(0.0), Lng: Rad(math.Pi / 2)},
		{Lat: Rad(-math.Pi / 2), Lng: Rad(0.0)},
		{Lat: Rad(0.0), Lng: Rad(0.0)},
	}

	compareGeoLoopArea(t, verts, 3*math.Pi)
}

func TestGeoLoopAreaHemisphereEast(t *testing.T) {
	t.Parallel()

	// Two 1/4 triangles stitched together covering the eastern hemisphere.
	verts := []LatLng{
		{Lat: Rad(math.Pi / 2), Lng: Rad(0.0)},
		{Lat: Rad(0.0), Lng: Rad(0.0)},
		{Lat: Rad(-math.Pi / 2), Lng: Rad(0.0)},
		{Lat: Rad(0.0), Lng: Rad(math.Pi)},
	}

	compareGeoLoopArea(t, verts, 2*math.Pi)
}

func TestGeoLoopAreaHemisphereNorth(t *testing.T) {
	t.Parallel()

	// Four 1/8 triangles stitched together covering the northern
	// hemisphere.
	verts := []LatLng{
		{Lat: Rad(0.0), Lng: Rad(-math.Pi)},
		{Lat: Rad(0.0), Lng: Rad(-math.Pi / 2)},
		{Lat: Rad(0.0), Lng: Rad(0.0)},
		{Lat: Rad(0.0), Lng: Rad(math.Pi / 2)},
	}

	compareGeoLoopArea(t, verts, 2*math.Pi)
}

func TestGeoLoopAreaPercentageSlice(t *testing.T) {
	t.Parallel()

	// Edge arcs between points should be less than 180 degrees: a
	// triangle sweeping t*pi radians along the equator has area t*pi for
	// t < 1 and (2+t)*pi for 1 < t < 2 (the triangle "flips" past the
	// antipodal discontinuity at t == 1, which is skipped).
	for t2 := 0.0; t2 <= 1.2; t2 += 0.01 {
		verts := []LatLng{
			{Lat: Rad(math.Pi / 2), Lng: Rad(0.0)},
			{Lat: Rad(0.0), Lng: Rad(-math.Pi / 2)},
			{Lat: Rad(0.0), Lng: Rad(t2*math.Pi - math.Pi/2)},
		}
		loop := GeoLoop(verts)

		out, err := geoLoopAreaRads2(loop)
		if err != eSuccess {
			t.Fatalf("geoLoopAreaRads2: %v", err)
		}

		const tol = 1e-13
		if t2 < 0.99 {
			if math.Abs(out-t2*math.Pi) > tol {
				t.Errorf("t=%v: expected area %v got %v", t2, t2*math.Pi, out)
			}
		} else if t2 > 1.01 {
			if math.Abs(out-(2+t2)*math.Pi) > tol {
				t.Errorf("t=%v: expected area %v got %v", t2, (2+t2)*math.Pi, out)
			}
		}
	}
}

func TestGeoLoopAreaPercentageSliceLarge(t *testing.T) {
	t.Parallel()

	// A large polygon with t > 1 is still representable: add intermediate
	// vertices so no edge arc exceeds 180 degrees.
	const t2 = 1.2
	verts := []LatLng{
		{Lat: Rad(math.Pi / 2), Lng: Rad(0.0)},
		{Lat: Rad(0.0), Lng: Rad(-math.Pi / 2)},
		{Lat: Rad(0.0), Lng: Rad(0.0)}, // Extra vertex so every angle is < 180 degrees
		{Lat: Rad(0.0), Lng: Rad(t2*math.Pi - math.Pi/2)},
	}

	compareGeoLoopArea(t, verts, t2*math.Pi)
}

func TestGeoLoopAreaDegenerateLoop2(t *testing.T) {
	t.Parallel()

	// geoLoopAreaRads2 works without error on degenerate loops,
	// returning 0 area.
	verts := []LatLng{
		{Lat: Rad(math.Pi / 2), Lng: Rad(0.0)},
		{Lat: Rad(0.0), Lng: Rad(-math.Pi / 2)},
	}
	compareGeoLoopArea(t, verts, 0.0)
}

func TestGeoLoopAreaDegenerateLoop1(t *testing.T) {
	t.Parallel()

	verts := []LatLng{
		{Lat: Rad(0.0), Lng: Rad(0.0)},
	}
	compareGeoLoopArea(t, verts, 0.0)
}

func TestGeoLoopAreaDegenerateLoop0(t *testing.T) {
	t.Parallel()

	compareGeoLoopArea(t, nil, 0.0)
}
