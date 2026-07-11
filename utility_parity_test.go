//go:build cgo && c2go

package h3

import (
	"testing"
)

func Test_h3Print_parity(t *testing.T) {
	testCases := []H3Index{
		0x85283473fffffff,
		0x8a2a1072b59ffff,
		0x8f734e64992abb7f,
		0,
		0xffffffffffffffff,
	}

	for _, h := range testCases {
		// We can't directly compare print output, but we can verify the functions run without crashing
		h3Print(h)
		h3PrintC(h)

		h3Println(h)
		h3PrintlnC(h)
	}
}

func Test_coordIjkPrint_parity(t *testing.T) {
	testCases := []*CoordIJK{
		{I: 0, J: 0, K: 0},
		{I: 1, J: 0, K: 0},
		{I: 0, J: 1, K: 0},
		{I: 0, J: 0, K: 1},
		{I: -1, J: 2, K: -3},
		{I: 100, J: -50, K: 75},
	}

	for _, c := range testCases {
		// We can't directly compare print output, but we can verify the functions run without crashing
		coordIjkPrint(c)
		coordIjkPrintC(c)
	}
}

func Test_geoToStringRads_parity(t *testing.T) {
	testCases := []*LatLng{
		{Lat: Rad(0.0), Lng: Rad(0.0)},
		{Lat: Rad(1.0), Lng: Rad(2.0)},
		{Lat: Rad(-0.5), Lng: Rad(1.5)},
		{Lat: Rad(1.5708), Lng: Rad(3.1416)},   // ~π/2, π
		{Lat: Rad(-1.5708), Lng: Rad(-3.1416)}, // ~-π/2, -π
	}

	for _, p := range testCases {
		goResult := geoToStringRads(p)
		cResult := geoToStringRadsC(p)

		if goResult != cResult {
			t.Errorf("geoToStringRads mismatch: Go=%s, C=%s", goResult, cResult)
		}
	}
}

func Test_geoToStringDegs_parity(t *testing.T) {
	testCases := []*LatLng{
		{Lat: Rad(0.0), Lng: Rad(0.0)},
		{Lat: Rad(1.0), Lng: Rad(2.0)},
		{Lat: Rad(-0.5), Lng: Rad(1.5)},
		{Lat: Rad(1.5708), Lng: Rad(3.1416)},   // ~π/2, π
		{Lat: Rad(-1.5708), Lng: Rad(-3.1416)}, // ~-π/2, -π
	}

	for _, p := range testCases {
		goResult := geoToStringDegs(p)
		cResult := geoToStringDegsC(p)

		if goResult != cResult {
			t.Errorf("geoToStringDegs mismatch: Go=%s, C=%s", goResult, cResult)
		}
	}
}

func Test_geoToStringDegsNoFmt_parity(t *testing.T) {
	testCases := []*LatLng{
		{Lat: Rad(0.0), Lng: Rad(0.0)},
		{Lat: Rad(1.0), Lng: Rad(2.0)},
		{Lat: Rad(-0.5), Lng: Rad(1.5)},
		{Lat: Rad(1.5708), Lng: Rad(3.1416)},   // ~π/2, π
		{Lat: Rad(-1.5708), Lng: Rad(-3.1416)}, // ~-π/2, -π
	}

	for _, p := range testCases {
		goResult := geoToStringDegsNoFmt(p)
		cResult := geoToStringDegsNoFmtC(p)

		if goResult != cResult {
			t.Errorf("geoToStringDegsNoFmt mismatch: Go=%s, C=%s", goResult, cResult)
		}
	}
}

func Test_geoPrint_parity(t *testing.T) {
	testCases := []*LatLng{
		{Lat: Rad(0.0), Lng: Rad(0.0)},
		{Lat: Rad(1.0), Lng: Rad(2.0)},
		{Lat: Rad(-0.5), Lng: Rad(1.5)},
	}

	for _, p := range testCases {
		// We can't directly compare print output, but we can verify the functions run without crashing
		geoPrint(p)
		geoPrintC(p)

		geoPrintln(p)
		geoPrintlnC(p)

		geoPrintNoFmt(p)
		geoPrintNoFmtC(p)

		geoPrintlnNoFmt(p)
		geoPrintlnNoFmtC(p)
	}
}

func Test_cellBoundaryPrint_parity(t *testing.T) {
	testCases := []*CellBoundary{
		{
			NumVerts: 3,
			Verts: []LatLng{
				{Lat: Rad(0.0), Lng: Rad(0.0)},
				{Lat: Rad(1.0), Lng: Rad(1.0)},
				{Lat: Rad(2.0), Lng: Rad(2.0)},
			},
		},
		{
			NumVerts: 6,
			Verts: []LatLng{
				{Lat: Rad(0.1), Lng: Rad(0.1)},
				{Lat: Rad(0.2), Lng: Rad(0.2)},
				{Lat: Rad(0.3), Lng: Rad(0.3)},
				{Lat: Rad(0.4), Lng: Rad(0.4)},
				{Lat: Rad(0.5), Lng: Rad(0.5)},
				{Lat: Rad(0.6), Lng: Rad(0.6)},
			},
		},
	}

	for _, b := range testCases {
		// We can't directly compare print output, but we can verify the functions run without crashing
		cellBoundaryPrint(b)
		cellBoundaryPrintC(b)

		cellBoundaryPrintln(b)
		cellBoundaryPrintlnC(b)
	}
}

func Test_bboxPrint_parity(t *testing.T) {
	testCases := []*BBox{
		{
			North: Rad(1.0),
			South: Rad(0.0),
			East:  Rad(2.0),
			West:  Rad(1.0),
		},
		{
			North: Rad(0.5),
			South: Rad(-0.5),
			East:  Rad(1.5),
			West:  Rad(-1.5),
		},
	}

	for _, bbox := range testCases {
		// We can't directly compare print output, but we can verify the functions run without crashing
		bboxPrint(bbox)
		bboxPrintC(bbox)

		bboxPrintln(bbox)
		bboxPrintlnC(bbox)
	}
}
