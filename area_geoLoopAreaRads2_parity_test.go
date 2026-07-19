//go:build cgo && c2go && h3v450

package h3

import (
	"math"
	"testing"
)

// Parity for the 4.5.0 area implementation: geoLoopAreaRads2 over the
// ported testGeoLoopArea fixtures and cellAreaRads2 over every base
// cell's center child at several resolutions. Error codes compare
// exactly; the areas admit ~1e-14 (upstream's own test tolerance):
// each cagnoli term runs through sin/cos, where Go's math library and
// the platform libm differ by an ulp on some inputs, and a loop
// accumulates up to a dozen such terms — the same ~1e-15 noise the
// discovery record measured between the 4.4.0 and 4.5.0 algorithms
// (§16). Bit-exact area parity is not attainable across libms.

func areaClose(a, b float64) bool {
	if a == b {
		return true
	}
	diff := math.Abs(a - b)
	scale := 1.0
	if s := math.Abs(a); s > scale {
		scale = s
	}
	return diff <= 1e-14*scale
}

func Test_geoLoopAreaRads2_parity(t *testing.T) {
	loops := [][]LatLng{
		{{Lat: Rad(math.Pi / 2)}, {Lat: Rad(0)}, {Lat: Rad(0), Lng: Rad(math.Pi / 2)}},
		{{Lat: Rad(0), Lng: Rad(math.Pi / 2)}, {Lat: Rad(0)}, {Lat: Rad(math.Pi / 2)}},
		{{Lat: Rad(math.Pi / 2)}, {Lat: Rad(0)}, {Lat: Rad(-math.Pi / 2)}, {Lat: Rad(0), Lng: Rad(math.Pi / 2)}},
		{{Lat: Rad(0), Lng: Rad(-math.Pi)}, {Lat: Rad(0), Lng: Rad(-math.Pi / 2)}, {Lat: Rad(0)}, {Lat: Rad(0), Lng: Rad(math.Pi / 2)}},
		{{Lat: Rad(math.Pi / 2)}, {Lat: Rad(0), Lng: Rad(-math.Pi / 2)}},
		{{Lat: Rad(0)}},
		nil,
	}
	for i, verts := range loops {
		goOut, goErr := geoLoopAreaRads2(GeoLoop(verts))
		cOut, cErr := geoLoopAreaRads2C(GeoLoop(verts))
		if goErr != cErr || !areaClose(goOut, cOut) {
			t.Errorf("loop %d: Go=(%v,%v) C=(%v,%v)", i, goOut, goErr, cOut, cErr)
		}
	}
}

func Test_cellAreaRads2_parity_450(t *testing.T) {
	for base := int32(0); base < 122; base++ {
		var cell h3Index
		setH3Index(&cell, 0, base, 0)
		for _, res := range []int32{0, 1, 2, 5, 9, 15} {
			child, err := cellToCenterChild(cell, res)
			if err != eSuccess {
				t.Fatalf("cellToCenterChild(base %d, res %d): %v", base, res, err)
			}
			goOut, goErr := cellAreaRads2(child)
			cOut, cErr := cellAreaRads2C(child)
			if goErr != cErr || !areaClose(goOut, cOut) {
				t.Fatalf("cellAreaRads2(%x): Go=(%v,%v) C=(%v,%v)",
					uint64(child), goOut, goErr, cOut, cErr)
			}
		}
	}

	// Error path: invalid cell.
	_, goErr := cellAreaRads2(0x7fffffffffffffff)
	_, cErr := cellAreaRads2C(0x7fffffffffffffff)
	if goErr != cErr {
		t.Errorf("cellAreaRads2(invalid): Go=%v C=%v", goErr, cErr)
	}
}

func Test_kadd_parity(t *testing.T) {
	// Pure arithmetic: bit-exact, including the compensation term.
	terms := []float64{
		1.0, 1e-16, -1e-16, 0.1, -4 * 3.141592653589793,
		2.2440497074541694, 1e300, -1e300, 5e-324,
	}
	var goA adder
	cA := adder{}
	for _, x := range terms {
		kadd(&goA, x)
		cA = kaddC(cA, x)
		if goA != cA {
			t.Fatalf("kadd(%v): Go=%+v C=%+v", x, goA, cA)
		}
	}
}

func Test_cagnoli_parity(t *testing.T) {
	// Trig-dependent: compares with vec3UlpClose (libm differences;
	// measured ≤1 ulp on these cases); see the file comment.
	pts := []LatLng{
		{Lat: Rad(0.5), Lng: Rad(-1.3)},
		{Lat: Rad(0.6), Lng: Rad(-1.2)},
		{Lat: Rad(-0.4), Lng: Rad(0.8)},
		{Lat: Rad(1.5), Lng: Rad(3.1)},
		{Lat: Rad(0), Lng: Rad(0)},
	}
	for _, x := range pts {
		for _, y := range pts {
			goT := _cagnoli(x, y)
			cT := _cagnoliC(x, y)
			if !vec3UlpClose(goT, cT) {
				t.Errorf("cagnoli(%v,%v): Go=%v C=%v", x, y, goT, cT)
			}
		}
	}
}

func Test_geoPolygonAreaRads2_parity(t *testing.T) {
	ccwTriangle := GeoLoop{
		{Lat: Rad(math.Pi / 2)}, {Lat: Rad(0)}, {Lat: Rad(0), Lng: Rad(math.Pi / 2)},
	}
	cwSmallHole := GeoLoop{
		{Lat: Rad(0.2), Lng: Rad(0.2)}, {Lat: Rad(0.3), Lng: Rad(0.25)}, {Lat: Rad(0.25), Lng: Rad(0.3)},
	}
	polys := []GeoPolygon{
		{GeoLoop: ccwTriangle},
		{GeoLoop: ccwTriangle, Holes: []GeoLoop{cwSmallHole}},
		{GeoLoop: ccwTriangle, Holes: []GeoLoop{cwSmallHole, cwSmallHole}},
		{},
	}
	for i, poly := range polys {
		goOut, goErr := geoPolygonAreaRads2(poly)
		cOut, cErr := geoPolygonAreaRads2C(poly)
		if goErr != cErr || !areaClose(goOut, cOut) {
			t.Errorf("poly %d: Go=(%v,%v) C=(%v,%v)", i, goOut, goErr, cOut, cErr)
		}
	}

	mpolys := []geoMultiPolygon{
		{},
		{NumPolygons: 1, Polygons: polys[:1]},
		{NumPolygons: 3, Polygons: polys[:3]},
	}
	for i, mp := range mpolys {
		goOut, goErr := geoMultiPolygonAreaRads2(mp)
		cOut, cErr := geoMultiPolygonAreaRads2C(mp)
		if goErr != cErr || !areaClose(goOut, cOut) {
			t.Errorf("mpoly %d: Go=(%v,%v) C=(%v,%v)", i, goOut, goErr, cOut, cErr)
		}
	}
}
