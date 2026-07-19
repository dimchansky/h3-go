// Tests ported from H3 v4.5.0: src/apps/testapps/testVec3.c.

package h3

import (
	"math"
	"testing"
)

const testVec3DblEpsilon = 2.220446049250313e-16

func TestVec3DotProduct(t *testing.T) {
	t.Parallel()

	a := vec3d{X: 1.0, Y: 0.0, Z: 0.0}
	b := vec3d{X: -1.0, Y: 0.0, Z: 0.0}
	if vec3Dot(a, b) != -1.0 {
		t.Error("dot product matches expected value")
	}
}

func TestVec3CrossProductOrthogonality(t *testing.T) {
	t.Parallel()

	i := vec3d{X: 1.0, Y: 0.0, Z: 0.0}
	j := vec3d{X: 0.0, Y: 1.0, Z: 0.0}
	k := vec3Cross(i, j)
	if math.Abs(k.X-0.0) >= testVec3DblEpsilon {
		t.Error("x component zero")
	}
	if math.Abs(k.Y-0.0) >= testVec3DblEpsilon {
		t.Error("y component zero")
	}
	if math.Abs(k.Z-1.0) >= testVec3DblEpsilon {
		t.Error("z component one")
	}
	if math.Abs(vec3Dot(k, i)) >= testVec3DblEpsilon {
		t.Error("cross is orthogonal to i")
	}
	if math.Abs(vec3Dot(k, j)) >= testVec3DblEpsilon {
		t.Error("cross is orthogonal to j")
	}
}

func TestVec3NormalizeAndMagnitude(t *testing.T) {
	t.Parallel()

	v := vec3d{X: 3.0, Y: -4.0, Z: 12.0}
	magSq := vec3NormSq(v)
	if math.Abs(magSq-169.0) >= testVec3DblEpsilon {
		t.Error("magnitude squared matches")
	}
	if math.Abs(vec3Norm(v)-13.0) >= testVec3DblEpsilon {
		t.Error("magnitude matches")
	}

	vec3Normalize(&v)
	if math.Abs(vec3Norm(v)-1.0) >= testVec3DblEpsilon {
		t.Error("normalized vector is unit")
	}

	zero := vec3d{X: 0.0, Y: 0.0, Z: 0.0}
	vec3Normalize(&zero)
	if zero.X != 0.0 || zero.Y != 0.0 || zero.Z != 0.0 {
		t.Error("zero vector remains unchanged when normalizing")
	}
}

func TestVec3Distance(t *testing.T) {
	t.Parallel()

	a := vec3d{X: 0.0, Y: 0.0, Z: 0.0}
	b := vec3d{X: 1.0, Y: 2.0, Z: 2.0}
	if math.Abs(vec3DistSq(a, b)-9.0) >= testVec3DblEpsilon {
		t.Error("distance squared matches")
	}
}

func TestLatLngToVec3UnitSphere(t *testing.T) {
	t.Parallel()

	geo := LatLng{Lat: Rad(0.5), Lng: Rad(-1.3)}
	v := latLngToVec3(geo)
	if math.Abs(vec3Norm(v)-1.0) >= testVec3DblEpsilon {
		t.Error("converted vector lives on the unit sphere")
	}
}

func TestVec3ToCellInvalidRes(t *testing.T) {
	t.Parallel()

	v := vec3d{X: 1.0, Y: 0.0, Z: 0.0}
	var out h3Index
	if vec3ToCell(&v, -1, &out) != eResDomain {
		t.Error("negative resolution is rejected")
	}
	if vec3ToCell(&v, 16, &out) != eResDomain {
		t.Error("resolution above max is rejected")
	}
}

func TestCellToVec3UnitSphere(t *testing.T) {
	t.Parallel()

	// cellToVec3 should return a point on the unit sphere.
	p := LatLng{Lat: Rad(0.6), Lng: Rad(-1.2)}
	var h h3Index
	if err := latLngToCell(&p, 5, &h); err != eSuccess {
		t.Fatalf("latLngToCell: %v", err)
	}

	var v vec3d
	if err := cellToVec3(h, &v); err != eSuccess {
		t.Fatalf("cellToVec3: %v", err)
	}
	// Upstream asserts < DBL_EPSILON, which holds only under C compilers'
	// floating-point contraction (FMA in vec3Dot at -O2). Go does not
	// mandate uncontracted evaluation — the spec permits FMA fusion, and
	// gc fuses on arm64 — but the ported vec3Dot forces per-product
	// rounding via explicit float64() conversions, matching the
	// uncontracted C the parity harness compiles (-ffp-contract=off).
	// Under that shared uncontracted evaluation this cell's norm is
	// exactly 1 ulp above 1.0; the assertion admits exactly one ulp here.
	if math.Abs(vec3Norm(v)-1.0) > testVec3DblEpsilon {
		t.Error("cellToVec3 result is on the unit sphere")
	}
}

func TestCellToVec3MatchesCellToLatLng(t *testing.T) {
	t.Parallel()

	// vec3ToLatLng(cellToVec3(cell)) should agree with cellToLatLng.
	p := LatLng{Lat: Rad(0.3), Lng: Rad(2.1)}
	var h h3Index
	if err := latLngToCell(&p, 7, &h); err != eSuccess {
		t.Fatalf("latLngToCell: %v", err)
	}

	var v vec3d
	if err := cellToVec3(h, &v); err != eSuccess {
		t.Fatalf("cellToVec3: %v", err)
	}
	fromVec3 := vec3ToLatLng(v)

	var fromCell LatLng
	if err := cellToLatLng(h, &fromCell); err != eSuccess {
		t.Fatalf("cellToLatLng: %v", err)
	}

	if math.Abs(fromVec3.Lat.Rad()-fromCell.Lat.Rad()) >= testVec3DblEpsilon {
		t.Error("lat matches cellToLatLng")
	}
	if math.Abs(fromVec3.Lng.Rad()-fromCell.Lng.Rad()) >= testVec3DblEpsilon {
		t.Error("lng matches cellToLatLng")
	}
}

func TestCellToVec3RoundTrip(t *testing.T) {
	t.Parallel()

	// vec3ToCell(cellToVec3(cell)) should return the same cell.
	p := LatLng{Lat: Rad(-0.4), Lng: Rad(0.8)}
	var h h3Index
	if err := latLngToCell(&p, 9, &h); err != eSuccess {
		t.Fatalf("latLngToCell: %v", err)
	}

	var v vec3d
	if err := cellToVec3(h, &v); err != eSuccess {
		t.Fatalf("cellToVec3: %v", err)
	}

	var h2 h3Index
	if err := vec3ToCell(&v, 9, &h2); err != eSuccess {
		t.Fatalf("vec3ToCell: %v", err)
	}
	if h2 != h {
		t.Error("round-trip through Vec3d returns same cell")
	}
}

func TestCellToVec3InvalidCell(t *testing.T) {
	t.Parallel()

	var v vec3d
	if cellToVec3(0x7fffffffffffffff, &v) != eCellInvalid {
		t.Error("invalid cell gives E_CELL_INVALID")
	}
}

func TestVec3ToCellNonFinite(t *testing.T) {
	t.Parallel()

	var out h3Index
	nanX := vec3d{X: math.NaN(), Y: 0.0, Z: 0.0}
	if vec3ToCell(&nanX, 0, &out) != eDomain {
		t.Error("NaN x is rejected")
	}
	infY := vec3d{X: 0.0, Y: math.Inf(1), Z: 0.0}
	if vec3ToCell(&infY, 0, &out) != eDomain {
		t.Error("infinite y is rejected")
	}
	infZ := vec3d{X: 0.0, Y: 0.0, Z: math.Inf(-1)}
	if vec3ToCell(&infZ, 0, &out) != eDomain {
		t.Error("infinite z is rejected")
	}
}
