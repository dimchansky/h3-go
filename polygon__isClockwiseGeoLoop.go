package h3

// isClockwiseGeoLoopNormalized determines whether the winding order of a given loop is clockwise,
// with normalization for loops crossing the antimeridian.
// Ported from H3 C: polygonAlgos.h::GENERIC_LOOP_ALGO(isClockwiseNormalized) -> isClockwiseGeoLoopNormalized.
func isClockwiseGeoLoopNormalized(loop GeoLoop, isTransmeridian bool) bool {
	sum := 0.0
	for i := 0; i < len(loop); i++ {
		a := loop[i]
		b := loop[(i+1)%len(loop)]

		// If we identify a transmeridian arc (> 180 degrees longitude),
		// start over with the transmeridian flag set
		if !isTransmeridian && (a.Lng-b.Lng).Abs() > Pi {
			return isClockwiseGeoLoopNormalized(loop, true)
		}

		sum += (normalizeLngTransmeridian(b.Lng.Rad(), isTransmeridian) -
			normalizeLngTransmeridian(a.Lng.Rad(), isTransmeridian)) *
			(b.Lat.Rad() + a.Lat.Rad())
	}

	return sum > 0
}

// isClockwiseGeoLoop determines whether the winding order of a given loop is clockwise.
// In GeoJSON, clockwise loops are always inner loops (holes).
// This function uses the shoelace formula to calculate the signed area and determine
// orientation, with proper handling of transmeridian polygons.
// Ported from H3 C: polygonAlgos.h::GENERIC_LOOP_ALGO(isClockwise) -> isClockwiseGeoLoop.
func isClockwiseGeoLoop(loop GeoLoop) bool {
	return isClockwiseGeoLoopNormalized(loop, false)
}
