package h3

// edgeLengthKm returns the length of a directed edge in kilometers.
// This function calls edgeLengthRads and converts the result to kilometers
// by multiplying by Earth's radius.
// Ported from H3 C: latLng.c::H3_EXPORT(edgeLengthKm).
func edgeLengthKm(edge h3Index, length *float64) h3Error {
	err := edgeLengthRads(edge, length)
	*length = *length * earthRadiusKm
	return err
}
