package h3

// isClockwiseLinkedGeoLoopNormalized determines clockwise winding order with normalization
// for loops crossing the antimeridian. This is a helper function that handles the core
// algorithm for winding order detection.
// Ported from H3 C: polygonAlgos.h::GENERIC_LOOP_ALGO(isClockwiseNormalized) -> isClockwiseLinkedGeoLoopNormalized.
func isClockwiseLinkedGeoLoopNormalized(loop *linkedGeoLoop, isTransmeridian bool) bool {
	sum := 0.0
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

		// If we identify a transmeridian arc (> 180 degrees longitude),
		// start over with the transmeridian flag set
		if !isTransmeridian && (a.Lng-b.Lng).Abs() > Pi {
			return isClockwiseLinkedGeoLoopNormalized(loop, true)
		}

		sum += (normalizeLngTransmeridian(b.Lng.Rad(), isTransmeridian) - normalizeLngTransmeridian(a.Lng.Rad(), isTransmeridian)) * (b.Lat.Rad() + a.Lat.Rad())
	}

	return sum > 0
}

// isClockwiseLinkedGeoLoop determines whether the winding order of a linked geographic loop
// is clockwise. In GeoJSON, clockwise loops are always inner loops (holes).
// This function uses the shoelace formula to calculate the signed area and determine
// orientation, with proper handling of transmeridian polygons.
// Ported from H3 C: polygonAlgos.h::GENERIC_LOOP_ALGO(isClockwise) -> isClockwiseLinkedGeoLoop.
func isClockwiseLinkedGeoLoop(loop *linkedGeoLoop) bool {
	return isClockwiseLinkedGeoLoopNormalized(loop, false)
}
