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
