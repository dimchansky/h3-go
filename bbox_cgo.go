//go:build cgo && c2go

package h3

/*
#include <stdint.h>
#include <stdbool.h>
#include "h3api.h"
#include "latLng.h"
#include "bbox.h"

// Prototype for scaleBBox
extern void scaleBBox(BBox* bbox, double scale);

// Prototype for _hexRadiusKm
extern double _hexRadiusKm(H3Index h3Index);

// Prototype for lineHexEstimate
extern H3Error lineHexEstimate(const LatLng* origin, const LatLng* destination, int res, int64_t* out);

// Prototype for bboxHexEstimate
extern H3Error bboxHexEstimate(const BBox* bbox, int res, int64_t* out);
*/
import "C"

func toCBBox(b BBox) C.BBox {
	var cb C.BBox
	cb.north = C.double(b.North.Rad())
	cb.south = C.double(b.South.Rad())
	cb.east = C.double(b.East.Rad())
	cb.west = C.double(b.West.Rad())
	return cb
}

// bboxIsTransmeridianC calls the original C implementation.
func bboxIsTransmeridianC(b BBox) bool {
	cb := toCBBox(b)
	if C.bboxIsTransmeridian(&cb) {
		return true
	} else {
		return false
	}
}

// bboxWidthRadsC calls the original C implementation.
func bboxWidthRadsC(b BBox) float64 {
	cb := toCBBox(b)
	return float64(C.bboxWidthRads(&cb))
}

// bboxHeightRadsC calls the original C implementation.
func bboxHeightRadsC(b BBox) float64 {
	cb := toCBBox(b)
	return float64(C.bboxHeightRads(&cb))
}

// bboxEqualsC calls the original C implementation.
func bboxEqualsC(b1, b2 BBox) bool {
	cb1 := toCBBox(b1)
	cb2 := toCBBox(b2)
	if C.bboxEquals(&cb1, &cb2) {
		return true
	} else {
		return false
	}
}

// bboxCenterC calls the original C implementation.
func bboxCenterC(b BBox) LatLng {
	cb := toCBBox(b)
	var center C.LatLng
	C.bboxCenter(&cb, &center)
	return LatLng{Lat: Rad(float64(center.lat)), Lng: Rad(float64(center.lng))}
}

// bboxContainsC calls the original C implementation.
func bboxContainsC(b BBox, p LatLng) bool {
	cb := toCBBox(b)
	var cp C.LatLng
	cp.lat = C.double(p.Lat.Rad())
	cp.lng = C.double(p.Lng.Rad())
	if C.bboxContains(&cb, &cp) {
		return true
	} else {
		return false
	}
}

// bboxOverlapsBBoxC calls the original C implementation.
func bboxOverlapsBBoxC(a, b BBox) bool {
	ca := toCBBox(a)
	cb := toCBBox(b)
	if C.bboxOverlapsBBox(&ca, &cb) {
		return true
	} else {
		return false
	}
}

// bboxContainsBBoxC calls the original C implementation.
func bboxContainsBBoxC(a, b BBox) bool {
	ca := toCBBox(a)
	cb := toCBBox(b)
	if C.bboxContainsBBox(&ca, &cb) {
		return true
	} else {
		return false
	}
}

// scaleBBoxC calls the original C implementation.
func scaleBBoxC(bbox *BBox, scale float64) {
	cb := toCBBox(*bbox)
	C.scaleBBox(&cb, C.double(scale))
	// Convert back to Go struct
	bbox.North = Rad(float64(cb.north))
	bbox.South = Rad(float64(cb.south))
	bbox.East = Rad(float64(cb.east))
	bbox.West = Rad(float64(cb.west))
}

// _hexRadiusKmC calls the original C implementation.
func _hexRadiusKmC(h3Index H3Index) float64 {
	return float64(C._hexRadiusKm(C.H3Index(h3Index)))
}

// lineHexEstimateC calls the original C implementation.
func lineHexEstimateC(origin *LatLng, destination *LatLng, res int32, out *int64) H3Error {
	var cOrigin C.LatLng
	cOrigin.lat = C.double(origin.Lat)
	cOrigin.lng = C.double(origin.Lng)

	var cDest C.LatLng
	cDest.lat = C.double(destination.Lat)
	cDest.lng = C.double(destination.Lng)

	var cOut C.int64_t
	err := C.lineHexEstimate(&cOrigin, &cDest, C.int(res), &cOut)
	*out = int64(cOut)
	return H3Error(err)
}

// bboxHexEstimateC calls the original C implementation.
func bboxHexEstimateC(bbox *BBox, res int32, out *int64) H3Error {
	cb := toCBBox(*bbox)
	var cOut C.int64_t
	err := C.bboxHexEstimate(&cb, C.int(res), &cOut)
	*out = int64(cOut)
	return H3Error(err)
}
