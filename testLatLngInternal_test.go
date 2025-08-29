// Tests ported from testLatLngInternal.c

package h3

import (
	"math"
	"testing"
)

func TestGeoAlmostEqualThreshold(t *testing.T) {
	t.Parallel()

	// same point
	a := LatLng{Lat: Rad(15), Lng: Rad(10)}
	b := LatLng{Lat: Rad(15), Lng: Rad(10)}
	if !geoAlmostEqualThreshold(&a, &b, 2.2204460492503131e-16) { // DBL_EPSILON
		t.Error("same point should be equal")
	}

	// differences under threshold
	b.Lat = Rad(15.00001)
	b.Lng = Rad(10.00002)
	if !geoAlmostEqualThreshold(&a, &b, 0.0001) {
		t.Error("differences under threshold should be equal")
	}

	// lat over threshold
	b.Lat = Rad(15.00001)
	b.Lng = Rad(10)
	if geoAlmostEqualThreshold(&a, &b, 0.000001) {
		t.Error("lat over threshold should not be equal")
	}

	// lng over threshold
	b.Lat = Rad(15)
	b.Lng = Rad(10.00001)
	if geoAlmostEqualThreshold(&a, &b, 0.000001) {
		t.Error("lng over threshold should not be equal")
	}
}

func TestConstrainLatLng(t *testing.T) {
	t.Parallel()

	// Test constrainLat
	if constrainLat(0) != 0 {
		t.Error("lat 0 should remain 0")
	}
	if constrainLat(1) != 1 {
		t.Error("lat 1 should remain 1")
	}
	if constrainLat(math.Pi/2) != math.Pi/2 {
		t.Error("lat pi/2 should remain pi/2")
	}
	if constrainLat(math.Pi) != 0 {
		t.Error("lat pi should become 0")
	}
	if constrainLat(math.Pi+1) != 1 {
		t.Error("lat pi+1 should become 1")
	}
	if constrainLat(2*math.Pi+1) != 1 {
		t.Error("lat 2pi+1 should become 1")
	}

	// Test constrainLng
	if constrainLng(Rad(0)) != Rad(0) {
		t.Error("lng 0 should remain 0")
	}
	if constrainLng(Rad(1)) != Rad(1) {
		t.Error("lng 1 should remain 1")
	}
	if constrainLng(Rad(math.Pi)) != Rad(math.Pi) {
		t.Error("lng pi should remain pi")
	}
	if constrainLng(Rad(2*math.Pi)) != Rad(0) {
		t.Error("lng 2pi should become 0")
	}
	if constrainLng(Rad(3*math.Pi)) != Rad(math.Pi) {
		t.Error("lng 3pi should become pi")
	}
	if constrainLng(Rad(4*math.Pi)) != Rad(0) {
		t.Error("lng 4pi should become 0")
	}
}

func TestGeoAzDistanceRadsNoop(t *testing.T) {
	t.Parallel()

	start := LatLng{Lat: Rad(15), Lng: Rad(10)}
	var out LatLng
	expected := LatLng{Lat: Rad(15), Lng: Rad(10)}

	_geoAzDistanceRads(&start, 0, 0, &out)
	if !geoAlmostEqual(&expected, &out) {
		t.Error("0 distance should produce same point")
	}
}

func TestGeoAzDistanceRadsDueNorthSouth(t *testing.T) {
	t.Parallel()

	var start, out, expected LatLng

	// Due north to north pole
	setGeoDegs(&start, 45, 1)
	setGeoDegs(&expected, 90, 0)
	_geoAzDistanceRads(&start, 0, degsToRads(45), &out)
	if !geoAlmostEqual(&expected, &out) {
		t.Error("due north to north pole should produce north pole")
	}

	// Due north to south pole, which doesn't get wrapped correctly
	setGeoDegs(&start, 45, 1)
	setGeoDegs(&expected, 270, 1)
	_geoAzDistanceRads(&start, 0, degsToRads(45+180), &out)
	if !geoAlmostEqual(&expected, &out) {
		t.Error("due north to south pole should produce south pole")
	}

	// Due south to south pole
	setGeoDegs(&start, -45, 2)
	setGeoDegs(&expected, -90, 0)
	_geoAzDistanceRads(&start, degsToRads(180), degsToRads(45), &out)
	if !geoAlmostEqual(&expected, &out) {
		t.Error("due south to south pole should produce south pole")
	}

	// Due north to non-pole
	setGeoDegs(&start, -45, 10)
	setGeoDegs(&expected, -10, 10)
	_geoAzDistanceRads(&start, 0, degsToRads(35), &out)
	if !geoAlmostEqual(&expected, &out) {
		t.Error("due north should produce expected result")
	}
}

func TestGeoAzDistanceRadsPoleToPole(t *testing.T) {
	t.Parallel()

	var start, out, expected LatLng

	// Azimuth doesn't really matter in this case. Any azimuth from the
	// north pole is south, any azimuth from the south pole is north.

	setGeoDegs(&start, 90, 0)
	setGeoDegs(&expected, -90, 0)
	_geoAzDistanceRads(&start, degsToRads(12), degsToRads(180), &out)
	if !geoAlmostEqual(&expected, &out) {
		t.Error("some direction to south pole should produce south pole")
	}

	setGeoDegs(&start, -90, 0)
	setGeoDegs(&expected, 90, 0)
	_geoAzDistanceRads(&start, degsToRads(34), degsToRads(180), &out)
	if !geoAlmostEqual(&expected, &out) {
		t.Error("some direction to north pole should produce north pole")
	}
}

func TestGeoAzDistanceRadsInvertible(t *testing.T) {
	t.Parallel()

	var start, out LatLng
	setGeoDegs(&start, 15, 10)

	azimuth := degsToRads(20)
	degrees180 := degsToRads(180)
	distance := degsToRads(15)

	_geoAzDistanceRads(&start, azimuth, distance, &out)
	if math.Abs(greatCircleDistanceRads(&start, &out)-distance) >= EPSILON_RAD {
		t.Error("moved distance should be as expected")
	}

	start2 := out
	_geoAzDistanceRads(&start2, azimuth+degrees180, distance, &out)
	// TODO: Epsilon is relatively large
	if greatCircleDistanceRads(&start, &out) >= 0.01 {
		t.Error("should be able to move back to origin")
	}
}
