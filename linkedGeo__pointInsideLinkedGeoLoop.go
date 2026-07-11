package h3

import (
	"math"
)

const (
	// dblEpsilon represents the smallest positive floating-point number
	// such that 1.0 + dblEpsilon != 1.0.
	dblEpsilon = 2.220446049250313e-16
)

// normalizeLngTransmeridian normalizes longitude for transmeridian arcs
// This is the NORMALIZE_LNG macro from polygonAlgos.h.
func normalizeLngTransmeridian(lng float64, isTransmeridian bool) float64 {
	if isTransmeridian && lng < 0 {
		return lng + 2*math.Pi
	}
	return lng
}

// pointInsideLinkedGeoLoop implements point-in-polygon testing using ray casting.
// This is the core algorithm for determining if a coordinate lies within a polygon
// defined by a linked list of coordinates. Uses ray casting with proper handling
// of edge cases and transmeridian polygons.
// Ported from H3 C: polygonAlgos.h::GENERIC_LOOP_ALGO(pointInside) -> pointInsideLinkedGeoLoop.
func pointInsideLinkedGeoLoop(loop *linkedGeoLoop, bbox *bbox, coord *LatLng) bool {
	// fail fast if we're outside the bounding box
	if !bboxContains(bbox, coord) {
		return false
	}
	isTransmeridian := bboxIsTransmeridian(bbox)
	contains := false

	lat := coord.Lat.Rad()
	lng := normalizeLngTransmeridian(coord.Lng.Rad(), isTransmeridian)

	var a, b LatLng

	// Initialize iteration variables (INIT_ITERATION)
	var currentCoord *linkedLatLng = nil
	var nextCoord *linkedLatLng

	for {
		// ITERATE(loop, a, b) macro expansion:
		// currentCoord = GET_NEXT_COORD(loop, currentCoord)
		if currentCoord == nil {
			currentCoord = loop.First
		} else {
			currentCoord = currentCoord.Next
		}

		// if (currentCoord == NULL) break;
		if currentCoord == nil {
			break
		}

		// vertexA = currentCoord->vertex;
		a = currentCoord.Vertex

		// nextCoord = GET_NEXT_COORD(loop, currentCoord->next);
		if currentCoord.Next == nil {
			nextCoord = loop.First
		} else {
			nextCoord = currentCoord.Next
		}

		// vertexB = nextCoord->vertex
		b = nextCoord.Vertex

		// Ray casting algo requires the second point to always be higher
		// than the first, so swap if needed
		if a.Lat.Rad() > b.Lat.Rad() {
			a, b = b, a
		}

		// If the latitude matches exactly, we'll hit an edge case where
		// the ray passes through the vertex twice on successive segment
		// checks. To avoid this, adjust the latitude northward if needed.
		//
		// NOTE: This currently means that a point at the north pole cannot
		// be contained in any polygon. This is acceptable in current usage,
		// because the point we test in this function at present is always
		// a cell center or vertex, and no cell has a center or vertex on the
		// north pole. If we need to expand this algo to more generic uses we
		// might need to handle this edge case.
		if lat == a.Lat.Rad() || lat == b.Lat.Rad() {
			lat += dblEpsilon
		}

		// If we're totally above or below the latitude ranges, the test
		// ray cannot intersect the line segment, so let's move on
		if lat < a.Lat.Rad() || lat > b.Lat.Rad() {
			continue
		}

		aLng := normalizeLngTransmeridian(a.Lng.Rad(), isTransmeridian)
		bLng := normalizeLngTransmeridian(b.Lng.Rad(), isTransmeridian)

		// Rays are cast in the longitudinal direction, in case a point
		// exactly matches, to decide tiebreakers, bias westerly
		if aLng == lng || bLng == lng {
			lng -= dblEpsilon
		}

		// For the latitude of the point, compute the longitude of the
		// point that lies on the line segment defined by a and b
		// This is done by computing the percent above a the lat is,
		// and traversing the same percent in the longitudinal direction
		// of a to b
		ratio := (lat - a.Lat.Rad()) / (b.Lat.Rad() - a.Lat.Rad())
		testLng := normalizeLngTransmeridian(aLng+(bLng-aLng)*ratio, isTransmeridian)

		// Intersection of the ray
		if testLng > lng {
			contains = !contains
		}
	}

	return contains
}
