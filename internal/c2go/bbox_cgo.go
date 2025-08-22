//go:build cgo && c2go

package c2go

/*
#include <stdint.h>
#include <stdbool.h>
#include "h3api.h"
#include "latLng.h"
#include "bbox.h"
*/
import "C"

func toCBBox(b BBox) C.BBox {
	var cb C.BBox
	cb.north = C.double(b.North)
	cb.south = C.double(b.South)
	cb.east = C.double(b.East)
	cb.west = C.double(b.West)
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
	return LatLng{Lat: float64(center.lat), Lng: float64(center.lng)}
}

// bboxContainsC calls the original C implementation.
func bboxContainsC(b BBox, p LatLng) bool {
	cb := toCBBox(b)
	var cp C.LatLng
	cp.lat = C.double(p.Lat)
	cp.lng = C.double(p.Lng)
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
