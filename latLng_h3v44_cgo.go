//go:build cgo && c2go && !h3v450

package h3

// Wrappers for latLng.c internals that exist only in the H3 4.4.0 tree:
// triangleArea and triangleEdgeLengthsToArea were removed in 4.5.0 with
// the area refactor (docs/sync/4.4.0-to-4.5.0.md §5.2). This file and its
// parity tests retire with I-B; docs/sync/h3v450-exclusion-inventory.md
// tracks the exclusion. (The _geoAzimuthRads/_geoAzDistanceRads wrappers
// that also lived here retired with I-A.)

/*
// Interop wrapper only; the original C sources are compiled via separate
// build-tagged C shim files in this package (see h3lib_*.c with //go:build c2go).
#include "h3api.h"
#include "latLng.h"

// Forward declarations for non-exported latLng.c functions used in tests
double triangleEdgeLengthsToArea(double a, double b, double c);
double triangleArea(const LatLng* a, const LatLng* b, const LatLng* c);
*/
import "C"

// triangleEdgeLengthsToAreaC calls the original C implementation.
func triangleEdgeLengthsToAreaC(a, b, c float64) float64 {
	return float64(C.triangleEdgeLengthsToArea(C.double(a), C.double(b), C.double(c)))
}

// triangleAreaC calls the original C implementation.
func triangleAreaC(a, b, c LatLng) float64 {
	var ca, cb, cc C.LatLng
	ca.lat = C.double(a.Lat.Rad())
	ca.lng = C.double(a.Lng.Rad())
	cb.lat = C.double(b.Lat.Rad())
	cb.lng = C.double(b.Lng.Rad())
	cc.lat = C.double(c.Lat.Rad())
	cc.lng = C.double(c.Lng.Rad())
	return float64(C.triangleArea(&ca, &cb, &cc))
}
