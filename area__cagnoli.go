package h3

import "math"

// _cagnoli is the Cagnoli contribution for edge arc x to y, following
// d3-geo's area implementation:
// https://github.com/d3/d3-geo/blob/8c53a90ae70c94bace73ecb02f2c792c649c86ba/src/area.js#L51-L70
//
// (Named cagnoli in C; the underscore marks it as a file-local helper of
// area.c. The explicit float64 conversion in the atan2 argument forces
// the product to round before the addition, defeating arm64 FMA fusion;
// see vec3Dot.)
// Ported from H3 C: area.c::cagnoli.
func _cagnoli(x, y LatLng) float64 {
	xLat := x.Lat.Rad()/2.0 + math.Pi/4.0
	yLat := y.Lat.Rad()/2.0 + math.Pi/4.0

	sa := math.Sin(xLat) * math.Sin(yLat)
	ca := math.Cos(xLat) * math.Cos(yLat)

	d := y.Lng.Rad() - x.Lng.Rad()
	sd := math.Sin(d)
	cd := math.Cos(d)

	return -2.0 * math.Atan2(sa*sd, float64(sa*cd)+ca)
}
