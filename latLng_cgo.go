//go:build cgo && c2go

package h3

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
func _constrainLngC(lng Angle) Angle {
	result := float64(C.constrainLng_c_wrapper(C.double(lng.Rad())))
	return Rad(result)
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
	ca.lat = C.double(a.Lat.Rad())
	ca.lng = C.double(a.Lng.Rad())
	cb.lat = C.double(b.Lat.Rad())
	cb.lng = C.double(b.Lng.Rad())
	return float64(C._geoAzimuthRads(&ca, &cb))
}

// _geoAzDistanceRadsC calls the original C internal implementation.
func _geoAzDistanceRadsC(p1 LatLng, az, distance float64) LatLng {
	var c1, c2 C.LatLng
	c1.lat = C.double(p1.Lat.Rad())
	c1.lng = C.double(p1.Lng.Rad())
	C._geoAzDistanceRads(&c1, C.double(az), C.double(distance), &c2)
	return LatLng{Lat: Rad(float64(c2.lat)), Lng: Rad(float64(c2.lng))}
}

// greatCircleDistanceRadsC calls the original C implementation.
func greatCircleDistanceRadsC(a, b LatLng) float64 {
	var ca, cb C.LatLng
	ca.lat = C.double(a.Lat.Rad())
	ca.lng = C.double(a.Lng.Rad())
	cb.lat = C.double(b.Lat.Rad())
	cb.lng = C.double(b.Lng.Rad())
	return float64(C.greatCircleDistanceRads(&ca, &cb))
}

// greatCircleDistanceKmC calls the original C implementation.
func greatCircleDistanceKmC(a, b LatLng) float64 {
	var ca, cb C.LatLng
	ca.lat = C.double(a.Lat.Rad())
	ca.lng = C.double(a.Lng.Rad())
	cb.lat = C.double(b.Lat.Rad())
	cb.lng = C.double(b.Lng.Rad())
	return float64(C.greatCircleDistanceKm(&ca, &cb))
}

// greatCircleDistanceMC calls the original C implementation.
func greatCircleDistanceMC(a, b LatLng) float64 {
	var ca, cb C.LatLng
	ca.lat = C.double(a.Lat.Rad())
	ca.lng = C.double(a.Lng.Rad())
	cb.lat = C.double(b.Lat.Rad())
	cb.lng = C.double(b.Lng.Rad())
	return float64(C.greatCircleDistanceM(&ca, &cb))
}

// normalizeLngC calls the original C implementation.
func normalizeLngC(lng float64, normalization longitudeNormalization) float64 {
	return float64(C.normalizeLng(C.double(lng), C.LongitudeNormalization(normalization)))
}

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

// _setGeoRadsC calls the original C implementation.
func _setGeoRadsC(p *LatLng, latRads, lngRads float64) {
	var cp C.LatLng
	cp.lat = C.double(p.Lat.Rad())
	cp.lng = C.double(p.Lng.Rad())
	C._setGeoRads(&cp, C.double(latRads), C.double(lngRads))
	p.Lat = Rad(float64(cp.lat))
	p.Lng = Rad(float64(cp.lng))
}

// setGeoDegsC calls the original C implementation.
func setGeoDegsC(p *LatLng, latDegs, lngDegs float64) {
	var cp C.LatLng
	cp.lat = C.double(p.Lat.Rad())
	cp.lng = C.double(p.Lng.Rad())
	C.setGeoDegs(&cp, C.double(latDegs), C.double(lngDegs))
	p.Lat = Rad(float64(cp.lat))
	p.Lng = Rad(float64(cp.lng))
}

// geoAlmostEqualThresholdC calls the original C function using plain doubles.
func geoAlmostEqualThresholdC(a, b LatLng, threshold float64) bool {
	var ca, cb C.LatLng
	ca.lat = C.double(a.Lat.Rad())
	ca.lng = C.double(a.Lng.Rad())
	cb.lat = C.double(b.Lat.Rad())
	cb.lng = C.double(b.Lng.Rad())
	// Convert C.bool to int via C helper for toolchain compatibility
	return C.h3_bool_to_int(C.geoAlmostEqualThreshold(&ca, &cb, C.double(threshold))) != 0
}

// geoAlmostEqualC calls the original C function using plain doubles.
func geoAlmostEqualC(a, b LatLng) bool {
	var ca, cb C.LatLng
	ca.lat = C.double(a.Lat.Rad())
	ca.lng = C.double(a.Lng.Rad())
	cb.lat = C.double(b.Lat.Rad())
	cb.lng = C.double(b.Lng.Rad())
	return C.h3_bool_to_int(C.geoAlmostEqual(&ca, &cb)) != 0
}

// cellAreaRads2C calls the original C implementation.
func cellAreaRads2C(cell h3Index) (float64, h3Error) {
	var out C.double
	err := h3Error(C.cellAreaRads2(C.H3Index(cell), &out))
	return float64(out), err
}

// cellAreaKm2C calls the original C implementation.
func cellAreaKm2C(cell h3Index) (float64, h3Error) {
	var out C.double
	err := h3Error(C.cellAreaKm2(C.H3Index(cell), &out))
	return float64(out), err
}

// cellAreaM2C calls the original C implementation.
func cellAreaM2C(cell h3Index) (float64, h3Error) {
	var out C.double
	err := h3Error(C.cellAreaM2(C.H3Index(cell), &out))
	return float64(out), err
}

// getNumCellsC calls the original C implementation.
func getNumCellsC(res int32) (int64, h3Error) {
	var out C.int64_t
	err := h3Error(C.getNumCells(C.int(res), &out))
	return int64(out), err
}

// edgeLengthRadsC calls the original C implementation.
func edgeLengthRadsC(edge h3Index) (float64, h3Error) {
	var length C.double
	err := h3Error(C.edgeLengthRads(C.H3Index(edge), &length))
	return float64(length), err
}

// edgeLengthKmC calls the original C implementation.
func edgeLengthKmC(edge h3Index) (float64, h3Error) {
	var length C.double
	err := h3Error(C.edgeLengthKm(C.H3Index(edge), &length))
	return float64(length), err
}

// edgeLengthMC calls the original C implementation.
func edgeLengthMC(edge h3Index) (float64, h3Error) {
	var length C.double
	err := h3Error(C.edgeLengthM(C.H3Index(edge), &length))
	return float64(length), err
}

// getHexagonAreaAvgKm2C calls the original C implementation.
func getHexagonAreaAvgKm2C(res int32) (float64, h3Error) {
	var out C.double
	err := h3Error(C.getHexagonAreaAvgKm2(C.int(res), &out))
	return float64(out), err
}

// getHexagonAreaAvgM2C calls the original C implementation.
func getHexagonAreaAvgM2C(res int32) (float64, h3Error) {
	var out C.double
	err := h3Error(C.getHexagonAreaAvgM2(C.int(res), &out))
	return float64(out), err
}

// getHexagonEdgeLengthAvgKmC calls the original C implementation.
func getHexagonEdgeLengthAvgKmC(res int32) (float64, h3Error) {
	var out C.double
	err := h3Error(C.getHexagonEdgeLengthAvgKm(C.int(res), &out))
	return float64(out), err
}

// getHexagonEdgeLengthAvgMC calls the original C implementation.
func getHexagonEdgeLengthAvgMC(res int32) (float64, h3Error) {
	var out C.double
	err := h3Error(C.getHexagonEdgeLengthAvgM(C.int(res), &out))
	return float64(out), err
}
