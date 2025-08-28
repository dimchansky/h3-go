package h3

// degsToRads converts decimal degrees to radians.
// Ported from H3 C: latLng.c::H3_EXPORT(degsToRads)
func degsToRads(degrees float64) float64 { return degrees * (3.14159265358979323846 / 180.0) }
