// Tests ported from H3 v4.5.0: src/apps/testapps/testVec3dInternal.c.
// (The v4.4.0 suite tested _pointSquareDist/_geoToVec3d; 4.5.0 renamed
// them to the by-value vec3DistSq/latLngToVec3 with the same fixtures and
// added the two vec3Normalize edge cases.)

package h3

import (
	"math"
	"testing"
)

func TestVec3DistSq(t *testing.T) {
	t.Parallel()

	v1 := vec3d{0, 0, 0}
	v2 := vec3d{1, 0, 0}
	v3 := vec3d{0, 1, 1}
	v4 := vec3d{1, 1, 1}
	v5 := vec3d{1, 1, 2}

	const dblEpsilon = 2.220446049250313e-16
	if math.Abs(vec3DistSq(v1, v1)) >= dblEpsilon {
		t.Error("distance to self is 0")
	}
	if math.Abs(vec3DistSq(v1, v2)-1) >= dblEpsilon {
		t.Error("distance to <1,0,0> is 1")
	}
	if math.Abs(vec3DistSq(v1, v3)-2) >= dblEpsilon {
		t.Error("distance to <0,1,1> is 2")
	}
	if math.Abs(vec3DistSq(v1, v4)-3) >= dblEpsilon {
		t.Error("distance to <1,1,1> is 3")
	}
	if math.Abs(vec3DistSq(v1, v5)-6) >= dblEpsilon {
		t.Error("distance to <1,1,2> is 6")
	}
}

func TestVec3NormalizeSmallNonzero(t *testing.T) {
	t.Parallel()

	// 1e-163 squared underflows to 0, so norm == 0.
	// vec3Normalize should produce the zero vector.
	v := vec3d{1e-163, 0, 0}

	if v.X == 0.0 {
		t.Error("vector is nonzero")
	}
	if vec3Norm(v) != 0.0 {
		t.Error("norm underflows to zero")
	}

	vec3Normalize(&v)
	if v.X != 0.0 || v.Y != 0.0 || v.Z != 0.0 {
		t.Error("underflowed vector normalizes to zero")
	}
}

func TestVec3NormalizeDblEpsilonHalf(t *testing.T) {
	t.Parallel()

	const dblEpsilon = 2.220446049250313e-16
	// DBL_EPSILON/2 is small but normalizes fine.
	v := vec3d{dblEpsilon / 2.0, 0, 0}

	if vec3Norm(v) >= dblEpsilon {
		t.Error("norm is small but nonzero")
	}

	vec3Normalize(&v)
	if math.Abs(v.X-1.0) >= dblEpsilon || v.Y != 0 || v.Z != 0 {
		t.Error("still normalizable to unit vector")
	}
}

func TestLatLngToVec3(t *testing.T) {
	t.Parallel()

	origin := vec3d{}

	c1 := LatLng{Lat: Rad(0), Lng: Rad(0)}
	p1 := latLngToVec3(c1)
	if math.Abs(vec3DistSq(origin, p1)-1) >= epsilonRad {
		t.Error("Geo point is on the unit sphere")
	}

	c2 := LatLng{Lat: Rad(math.Pi / 2), Lng: Rad(0)}
	p2 := latLngToVec3(c2)
	if math.Abs(vec3DistSq(p1, p2)-2) >= epsilonRad {
		t.Error("Geo point is on another axis")
	}

	c3 := LatLng{Lat: Rad(math.Pi), Lng: Rad(0)}
	p3 := latLngToVec3(c3)
	if math.Abs(vec3DistSq(p1, p3)-4) >= epsilonRad {
		t.Error("Geo point is the other side of the sphere")
	}
}
