//go:build cgo && c2go && h3v450

// area.c is new in H3 4.5.0 (docs/sync/4.4.0-to-4.5.0.md §5.2): the
// relocated cellAreaRads2/Km2/M2 plus geoLoop/geoPolygon/geoMultiPolygon
// area helpers. Compiled only in the h3v450 harness configuration.
#include "area.c"
