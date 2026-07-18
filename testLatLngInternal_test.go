// Tests ported from H3 v4.4.0: src/apps/testapps/testLatLngInternal.c.
// v4.5.0 delta incorporated: the four _geoAzDistanceRads_* cases were
// deleted upstream with the function itself, and constrainLatLng gained
// the -2pi/-3pi constrainLng assertions.

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
	if !geoAlmostEqualThreshold(&a, &b, 2.2204460492503131e-16) { // dblEpsilon
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
	if constrainLng(Rad(-2*math.Pi)) != Rad(0) {
		t.Error("lng -2pi should become 0")
	}
	if constrainLng(Rad(-3*math.Pi)) != Rad(-math.Pi) {
		t.Error("lng -3pi should become -pi")
	}
}
