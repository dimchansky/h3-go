package c2go

// edgeLengthKm returns the length of a directed edge in kilometers.
// This function calls edgeLengthRads and converts the result to kilometers
// by multiplying by Earth's radius.
// Ported from H3 C: latLng.c::H3_EXPORT(edgeLengthKm)
func edgeLengthKm(edge H3Index, length *float64) H3Error {
	err := edgeLengthRads(edge, length)
	*length = *length * EARTH_RADIUS_KM
	return err
}
