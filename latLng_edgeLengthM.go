package h3

// edgeLengthM returns the length of a directed edge in meters.
// This function calls edgeLengthKm and converts the result to meters
// by multiplying by 1000.
// Ported from H3 C: latLng.c::H3_EXPORT(edgeLengthM).
func edgeLengthM(edge h3Index, length *float64) h3Error {
	err := edgeLengthKm(edge, length)
	*length = *length * 1000
	return err
}
