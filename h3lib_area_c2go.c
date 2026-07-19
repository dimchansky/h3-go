//go:build cgo && c2go

// area.c is new in H3 4.5.0 (docs/sync/4.4.0-to-4.5.0.md §5.2): the
// relocated cellAreaRads2/Km2/M2 plus geoLoop/geoPolygon/geoMultiPolygon
// area helpers.
#include "area.c"

// Test-only wrapper exposing the file-static cagnoli edge term to the
// parity harness, in the same translation unit as the static it calls.
double h3goTest_cagnoli(LatLng x, LatLng y) { return cagnoli(x, y); }
