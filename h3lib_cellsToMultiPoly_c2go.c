//go:build cgo && c2go && h3v450

// cellsToMultiPoly.c is new in H3 4.5.0 (docs/sync/4.4.0-to-4.5.0.md
// §5.2): the arc-based cellsToMultiPolygon machinery that the rewritten
// algos.c::cellsToLinkedMultiPolygon delegates to. Compiled only in the
// h3v450 harness configuration.
#include "cellsToMultiPoly.c"
