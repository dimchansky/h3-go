package h3

// cellAreaKm2 computes the area of an H3 cell in kilometers^2.
//
// This function converts the result from cellAreaRads2 (radians^2) to
// kilometers^2 by multiplying by the square of Earth's radius in kilometers.
// Uses WGS84 authalic radius: 6371.007180918475 km.
// Ported from H3 C: latLng.c::cellAreaKm2.
func cellAreaKm2(cell h3Index) (float64, h3Error) {
	areaRads2, err := cellAreaRads2(cell)
	if err != eSuccess {
		return 0.0, err
	}
	return areaRads2 * earthRadiusKm * earthRadiusKm, eSuccess
}
