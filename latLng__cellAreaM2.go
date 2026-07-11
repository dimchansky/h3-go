package h3

// cellAreaM2 computes the area of an H3 cell in square meters.
//
// This function converts the result from cellAreaKm2 (kilometers^2) to
// square meters by multiplying by 1,000,000 (1000 * 1000).
// Ported from H3 C: latLng.c::cellAreaM2.
func cellAreaM2(cell h3Index) (float64, h3Error) {
	areaKm2, err := cellAreaKm2(cell)
	if err != eSuccess {
		return 0.0, err
	}
	return areaKm2 * 1000 * 1000, eSuccess
}
