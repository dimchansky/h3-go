// Tests ported from testH3Api.c
package h3

import (
	"math"
	"testing"
)

func Test_latLngToCell_res(t *testing.T) {
	t.Parallel()
	var h H3Index
	anywhere := LatLng{Lat: 0, Lng: 0}

	// Test resolution below 0 is invalid
	err := latLngToCell(&anywhere, -1, &h)
	if err != E_RES_DOMAIN {
		t.Errorf("expected E_RES_DOMAIN for resolution -1, got %v", err)
	}

	// Test resolution above 15 is invalid
	err = latLngToCell(&anywhere, 16, &h)
	if err != E_RES_DOMAIN {
		t.Errorf("expected E_RES_DOMAIN for resolution 16, got %v", err)
	}
}

func Test_latLngToCell_coord(t *testing.T) {
	t.Parallel()
	var h H3Index

	// Test invalid latitude (NaN)
	invalidLat := LatLng{Lat: Angle(math.NaN()), Lng: 0}
	err := latLngToCell(&invalidLat, 1, &h)
	if err != E_LATLNG_DOMAIN {
		t.Errorf("expected E_LATLNG_DOMAIN for invalid latitude, got %v", err)
	}

	// Test invalid longitude (NaN)
	invalidLng := LatLng{Lat: 0, Lng: Angle(math.NaN())}
	err = latLngToCell(&invalidLng, 1, &h)
	if err != E_LATLNG_DOMAIN {
		t.Errorf("expected E_LATLNG_DOMAIN for invalid longitude, got %v", err)
	}

	// Test coordinates with infinity
	invalidLatLng := LatLng{Lat: Angle(math.Inf(1)), Lng: Angle(math.Inf(-1))}
	err = latLngToCell(&invalidLatLng, 1, &h)
	if err != E_LATLNG_DOMAIN {
		t.Errorf("expected E_LATLNG_DOMAIN for coordinates with infinity, got %v", err)
	}
}

func Test_cellToBoundary_classIIIEdgeVertex(t *testing.T) {
	t.Parallel()
	// Bug test for https://github.com/uber/h3/issues/45
	hexes := []H3Index{
		0x894cc5349b7ffff, 0x894cc534d97ffff, 0x894cc53682bffff,
		0x894cc536b17ffff, 0x894cc53688bffff, 0x894cead92cbffff,
		0x894cc536537ffff, 0x894cc5acbabffff, 0x894cc536597ffff,
	}

	var b CellBoundary
	for i, hex := range hexes {
		err := cellToBoundary(hex, &b)
		if err != E_SUCCESS {
			t.Errorf("cellToBoundary failed for hex %d: %v", i, err)
			continue
		}
		if b.NumVerts != 7 {
			t.Errorf("expected 7 vertices for hex %d, got %d", i, b.NumVerts)
		}
	}
}

func Test_cellToBoundary_classIIIEdgeVertex_exact(t *testing.T) {
	t.Parallel()
	// Bug test for https://github.com/uber/h3/issues/45
	h3, err := stringToH3("894cc536537ffff")
	if err != 0 {
		t.Fatalf("stringToH3 failed: %v", err)
	}

	var boundary CellBoundary
	err2 := cellToBoundary(h3, &boundary)
	if err2 != E_SUCCESS {
		t.Fatalf("cellToBoundary failed: %v", err2)
	}

	expected := CellBoundary{
		NumVerts: 7,
		Verts: []LatLng{
			{Lat: Deg(18.043333154), Lng: Deg(-66.27836523500002)},
			{Lat: Deg(18.042238363), Lng: Deg(-66.27929062800001)},
			{Lat: Deg(18.040818259), Lng: Deg(-66.27854193899998)},
			{Lat: Deg(18.040492975), Lng: Deg(-66.27686786700002)},
			{Lat: Deg(18.041040385), Lng: Deg(-66.27640518300001)},
			{Lat: Deg(18.041757122), Lng: Deg(-66.27596711500001)},
			{Lat: Deg(18.043007860), Lng: Deg(-66.27669118199998)},
		},
	}

	if boundary.NumVerts != expected.NumVerts {
		t.Errorf("expected %d vertices, got %d", expected.NumVerts, boundary.NumVerts)
	}

	// Check each vertex with tolerance
	tolerance := 1e-9
	for i := 0; i < int(boundary.NumVerts) && i < len(boundary.Verts); i++ {
		latDiff := math.Abs(boundary.Verts[i].Lat.Deg() - expected.Verts[i].Lat.Deg())
		lngDiff := math.Abs(boundary.Verts[i].Lng.Deg() - expected.Verts[i].Lng.Deg())

		if latDiff > tolerance {
			t.Errorf("vertex %d latitude mismatch: expected %f, got %f, diff %f",
				i, expected.Verts[i].Lat.Deg(), boundary.Verts[i].Lat.Deg(), latDiff)
		}
		if lngDiff > tolerance {
			t.Errorf("vertex %d longitude mismatch: expected %f, got %f, diff %f",
				i, expected.Verts[i].Lng.Deg(), boundary.Verts[i].Lng.Deg(), lngDiff)
		}
	}
}

func Test_cellToBoundary_coslngConstrain(t *testing.T) {
	t.Parallel()
	// Bug test for https://github.com/uber/h3/issues/212
	h3 := H3Index(0x87dc6d364ffffff)

	var boundary CellBoundary
	err := cellToBoundary(h3, &boundary)
	if err != E_SUCCESS {
		t.Fatalf("cellToBoundary failed: %v", err)
	}

	expected := CellBoundary{
		NumVerts: 6,
		Verts: []LatLng{
			{Lat: Deg(-52.0130533678236091), Lng: Deg(-34.6232931343713091)},
			{Lat: Deg(-52.0041156384652012), Lng: Deg(-34.6096733160584549)},
			{Lat: Deg(-51.9929610229502472), Lng: Deg(-34.6165157145896387)},
			{Lat: Deg(-51.9907410568096608), Lng: Deg(-34.6369680004259877)},
			{Lat: Deg(-51.9996738734672377), Lng: Deg(-34.6505896528323660)},
			{Lat: Deg(-52.0108315681413629), Lng: Deg(-34.6437571897165668)},
		},
	}

	if boundary.NumVerts != expected.NumVerts {
		t.Errorf("expected %d vertices, got %d", expected.NumVerts, boundary.NumVerts)
	}

	// Check each vertex with tolerance
	tolerance := 1e-9
	for i := 0; i < int(boundary.NumVerts) && i < len(boundary.Verts); i++ {
		latDiff := math.Abs(boundary.Verts[i].Lat.Deg() - expected.Verts[i].Lat.Deg())
		lngDiff := math.Abs(boundary.Verts[i].Lng.Deg() - expected.Verts[i].Lng.Deg())

		if latDiff > tolerance {
			t.Errorf("vertex %d latitude mismatch: expected %f, got %f, diff %f",
				i, expected.Verts[i].Lat.Deg(), boundary.Verts[i].Lat.Deg(), latDiff)
		}
		if lngDiff > tolerance {
			t.Errorf("vertex %d longitude mismatch: expected %f, got %f, diff %f",
				i, expected.Verts[i].Lng.Deg(), boundary.Verts[i].Lng.Deg(), lngDiff)
		}
	}
}

func Test_cellToBoundary_failed(t *testing.T) {
	t.Parallel()
	h := H3Index(0x87dc6d364ffffff)
	// Set an invalid base cell (NUM_BASE_CELLS + 1 = 122 + 1 = 123)
	// Need to check how base cells are encoded in Go version
	invalidH := h | (H3Index(123) << H3_BC_OFFSET)

	var gb CellBoundary
	err := cellToBoundary(invalidH, &gb)
	if err != E_CELL_INVALID {
		t.Errorf("expected E_CELL_INVALID for invalid base cell, got %v", err)
	}
}

func Test_cellToLatLngInvalid(t *testing.T) {
	t.Parallel()
	var coord LatLng
	err := cellToLatLng(H3Index(0x7fffffffffffffff), &coord)
	if err != E_CELL_INVALID {
		t.Errorf("expected E_CELL_INVALID for invalid cell, got %v", err)
	}
}

func Test_version(t *testing.T) {
	t.Parallel()
	// Port of C test: t_assert(H3_VERSION_MAJOR >= 0, "major version is set");
	if H3_VERSION_MAJOR < 0 {
		t.Error("major version is set")
	}

	// Port of C test: t_assert(H3_VERSION_MINOR >= 0, "minor version is set");
	if H3_VERSION_MINOR < 0 {
		t.Error("minor version is set")
	}

	// Port of C test: t_assert(H3_VERSION_PATCH >= 0, "patch version is set");
	if H3_VERSION_PATCH < 0 {
		t.Error("patch version is set")
	}
}
