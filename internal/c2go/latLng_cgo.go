//go:build cgo

package c2go

/*
// Interop wrapper only; the original C sources are compiled via separate
// build-tagged C shim files in this package (see h3lib_*.c with //go:build c2go).
#include <stdint.h>
#include <stdbool.h>
#include "h3api.h"
#include "latLng.h"
// Normalize C bool to int for cgo comparisons when needed (toolchain-safe)
static int h3_bool_to_int(_Bool b) { return b ? 1 : 0; }

// Forward declarations for non-exported latLng.c functions used in tests
double triangleEdgeLengthsToArea(double a, double b, double c);
double triangleArea(const LatLng* a, const LatLng* b, const LatLng* c);

// Forward declarations for wrappers
double _posAngleRads(double);
double constrainLng(double);
double constrainLat(double);
double degsToRads(double);
double radsToDegs(double);

static double _posAngleRads_c_wrapper(double rads) { return _posAngleRads(rads); }
static double constrainLng_c_wrapper(double lng) { return constrainLng(lng); }
static double constrainLat_c_wrapper(double lat) { return constrainLat(lat); }
static double degsToRads_c_wrapper(double degrees) { return degsToRads(degrees); }
static double radsToDegs_c_wrapper(double radians) { return radsToDegs(radians); }
*/
import "C"

// _posAngleRadsC invokes the original C implementation for parity tests.
func _posAngleRadsC(rads float64) float64 {
	return float64(C._posAngleRads_c_wrapper(C.double(rads)))
}

// _constrainLngC invokes the original C implementation for parity tests.
func _constrainLngC(lng float64) float64 {
	return float64(C.constrainLng_c_wrapper(C.double(lng)))
}

// _constrainLatC invokes the original C implementation for parity tests.
func _constrainLatC(lat float64) float64 {
	return float64(C.constrainLat_c_wrapper(C.double(lat)))
}

// degsToRadsC invokes the original C implementation for parity tests.
func degsToRadsC(degrees float64) float64 {
	return float64(C.degsToRads_c_wrapper(C.double(degrees)))
}

// radsToDegsC invokes the original C implementation for parity tests.
func radsToDegsC(radians float64) float64 {
	return float64(C.radsToDegs_c_wrapper(C.double(radians)))
}

// _geoAzimuthRadsC calls the original C internal implementation.
func _geoAzimuthRadsC(a, b LatLng) float64 {
	var ca, cb C.LatLng
	ca.lat = C.double(a.Lat)
	ca.lng = C.double(a.Lng)
	cb.lat = C.double(b.Lat)
	cb.lng = C.double(b.Lng)
	return float64(C._geoAzimuthRads(&ca, &cb))
}

// _geoAzDistanceRadsC calls the original C internal implementation.
func _geoAzDistanceRadsC(p1 LatLng, az, distance float64) LatLng {
	var c1, c2 C.LatLng
	c1.lat = C.double(p1.Lat)
	c1.lng = C.double(p1.Lng)
	C._geoAzDistanceRads(&c1, C.double(az), C.double(distance), &c2)
	return LatLng{Lat: float64(c2.lat), Lng: float64(c2.lng)}
}

// greatCircleDistanceRadsC calls the original C implementation.
func greatCircleDistanceRadsC(a, b LatLng) float64 {
	var ca, cb C.LatLng
	ca.lat = C.double(a.Lat)
	ca.lng = C.double(a.Lng)
	cb.lat = C.double(b.Lat)
	cb.lng = C.double(b.Lng)
	return float64(C.greatCircleDistanceRads(&ca, &cb))
}

// greatCircleDistanceKmC calls the original C implementation.
func greatCircleDistanceKmC(a, b LatLng) float64 {
	var ca, cb C.LatLng
	ca.lat = C.double(a.Lat)
	ca.lng = C.double(a.Lng)
	cb.lat = C.double(b.Lat)
	cb.lng = C.double(b.Lng)
	return float64(C.greatCircleDistanceKm(&ca, &cb))
}

// greatCircleDistanceMC calls the original C implementation.
func greatCircleDistanceMC(a, b LatLng) float64 {
	var ca, cb C.LatLng
	ca.lat = C.double(a.Lat)
	ca.lng = C.double(a.Lng)
	cb.lat = C.double(b.Lat)
	cb.lng = C.double(b.Lng)
	return float64(C.greatCircleDistanceM(&ca, &cb))
}

// normalizeLngC calls the original C implementation.
func normalizeLngC(lng float64, normalization LongitudeNormalization) float64 {
	return float64(C.normalizeLng(C.double(lng), C.LongitudeNormalization(normalization)))
}

// triangleEdgeLengthsToAreaC calls the original C implementation.
func triangleEdgeLengthsToAreaC(a, b, c float64) float64 {
	return float64(C.triangleEdgeLengthsToArea(C.double(a), C.double(b), C.double(c)))
}

// triangleAreaC calls the original C implementation.
func triangleAreaC(a, b, c LatLng) float64 {
	var ca, cb, cc C.LatLng
	ca.lat = C.double(a.Lat)
	ca.lng = C.double(a.Lng)
	cb.lat = C.double(b.Lat)
	cb.lng = C.double(b.Lng)
	cc.lat = C.double(c.Lat)
	cc.lng = C.double(c.Lng)
	return float64(C.triangleArea(&ca, &cb, &cc))
}

// _setGeoRadsC calls the original C implementation.
func _setGeoRadsC(p *LatLng, latRads, lngRads float64) {
	var cp C.LatLng
	cp.lat = C.double(p.Lat)
	cp.lng = C.double(p.Lng)
	C._setGeoRads(&cp, C.double(latRads), C.double(lngRads))
	p.Lat = float64(cp.lat)
	p.Lng = float64(cp.lng)
}

// setGeoDegsC calls the original C implementation.
func setGeoDegsC(p *LatLng, latDegs, lngDegs float64) {
	var cp C.LatLng
	cp.lat = C.double(p.Lat)
	cp.lng = C.double(p.Lng)
	C.setGeoDegs(&cp, C.double(latDegs), C.double(lngDegs))
	p.Lat = float64(cp.lat)
	p.Lng = float64(cp.lng)
}

// geoAlmostEqualThresholdC calls the original C function using plain doubles.
func geoAlmostEqualThresholdC(a, b LatLng, threshold float64) bool {
	var ca, cb C.LatLng
	ca.lat = C.double(a.Lat)
	ca.lng = C.double(a.Lng)
	cb.lat = C.double(b.Lat)
	cb.lng = C.double(b.Lng)
	// Convert C.bool to int via C helper for toolchain compatibility
	return C.h3_bool_to_int(C.geoAlmostEqualThreshold(&ca, &cb, C.double(threshold))) != 0
}

// geoAlmostEqualC calls the original C function using plain doubles.
func geoAlmostEqualC(a, b LatLng) bool {
	var ca, cb C.LatLng
	ca.lat = C.double(a.Lat)
	ca.lng = C.double(a.Lng)
	cb.lat = C.double(b.Lat)
	cb.lng = C.double(b.Lng)
	return C.h3_bool_to_int(C.geoAlmostEqual(&ca, &cb)) != 0
}
