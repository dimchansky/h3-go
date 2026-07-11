package h3

// radsToDegs converts radians to decimal degrees.
// Ported from H3 C: latLng.c::H3_EXPORT(radsToDegs).
func radsToDegs(radians float64) float64 { return radians * (180.0 / 3.14159265358979323846) }
