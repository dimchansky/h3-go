package h3

import "math"

// geoLoopAreaRads2 computes the area in radians^2 enclosed by vertices
// in a GeoLoop.
//
// The GeoLoop should represent a simple curve with no self-intersections,
// vertices ordered by the right-hand rule (counter-clockwise interior).
// The loop is closed automatically; edge arcs are interpreted as the
// shortest geodesic (< pi radians). A reversed loop yields 4*pi - a; the
// function does not return min(a, 4*pi - a). Result is in [0, 4*pi].
// Ported from H3 C: area.c::geoLoopAreaRads2.
func geoLoopAreaRads2(loop GeoLoop) (float64, h3Error) {
	// Use `Adder` to improve numerical accuracy of the sum of many Cagnoli
	// terms
	var a adder

	for i := 0; i < len(loop); i++ {
		j := (i + 1) % len(loop)
		kadd(&a, _cagnoli(loop[i], loop[j]))
	}

	// The Cagnoli sum above yields a signed area, with the sign switching
	// with the orientation of the vertices. Since we want our area to always be
	// positive, we normalize into [0, 4*pi] by adding 4*pi when the signed
	// area is negative.
	if a.sum < 0.0 {
		kadd(&a, 4.0*math.Pi)
	}

	return a.sum, eSuccess
}
