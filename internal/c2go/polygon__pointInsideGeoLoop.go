package c2go

import "math"

// pointInsideGeoLoop determines if coord is inside loop, using bbox for fast rejection.
// Ported from polygonAlgos.h::pointInside (GeoLoop variant).
// Ported from H3 C: polygon.c::pointInsideGeoLoop
func pointInsideGeoLoop(loop []LatLng, bbox *BBox, coord *LatLng) bool {
	if !bboxContains(bbox, coord) {
		return false
	}
	isTrans := bboxIsTransmeridian(bbox)
	contains := false
	const dblEps = 2.220446049250313e-16
	lat := coord.Lat
	lng := coord.Lng
	if isTrans && lng < 0 {
		lng += 2 * math.Pi
	}
	for i := 0; i < len(loop); i++ {
		a := loop[i]
		b := loop[(i+1)%len(loop)]
		if a.Lat > b.Lat {
			a, b = b, a
		}
		if lat == a.Lat || lat == b.Lat {
			lat += dblEps
		}
		if lat < a.Lat || lat > b.Lat {
			continue
		}
		aLng := a.Lng
		bLng := b.Lng
		if isTrans && aLng < 0 {
			aLng += 2 * math.Pi
		}
		if isTrans && bLng < 0 {
			bLng += 2 * math.Pi
		}
		if aLng == lng || bLng == lng {
			lng -= dblEps
		}
		ratio := (lat - a.Lat) / (b.Lat - a.Lat)
		testLng := aLng + (bLng-aLng)*ratio
		if isTrans && testLng < 0 {
			testLng += 2 * math.Pi
		}
		if testLng > lng {
			contains = !contains
		}
	}
	return contains
}
