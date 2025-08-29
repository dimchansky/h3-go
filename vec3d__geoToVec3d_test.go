// Tests ported from testVec3dInternal.c
package h3

import (
	"math"
	"testing"
)

func Test_geoToVec3d(t *testing.T) {
	t.Parallel()
	
	origin := &Vec3d{0, 0, 0}

	// Test: Geo point is on the unit sphere
	c1 := &LatLng{0, 0}
	var p1 Vec3d
	_geoToVec3d(c1, &p1)
	if math.Abs(_pointSquareDist(origin, &p1)-1) >= EPSILON_RAD {
		t.Error("Geo point should be on the unit sphere")
	}

	// Test: Geo point is on another axis
	c2 := &LatLng{math.Pi / 2, 0}
	var p2 Vec3d
	_geoToVec3d(c2, &p2)
	if math.Abs(_pointSquareDist(&p1, &p2)-2) >= EPSILON_RAD {
		t.Error("Geo point should be on another axis")
	}

	// Test: Geo point is the other side of the sphere
	c3 := &LatLng{math.Pi, 0}
	var p3 Vec3d
	_geoToVec3d(c3, &p3)
	if math.Abs(_pointSquareDist(&p1, &p3)-4) >= EPSILON_RAD {
		t.Error("Geo point should be the other side of the sphere")
	}
}