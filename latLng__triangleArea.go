package h3

// triangleArea computes area of a spherical triangle given its vertices.
// Ported from H3 C: latLng.c::triangleArea
func triangleArea(a, b, c *LatLng) float64 {
	return triangleEdgeLengthsToArea(
		greatCircleDistanceRads(a, b),
		greatCircleDistanceRads(b, c),
		greatCircleDistanceRads(c, a),
	)
}
