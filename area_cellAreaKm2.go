package h3

// cellAreaKm2 computes the area of an H3 cell in kilometers^2.
// Ported from H3 C: area.c::cellAreaKm2.
func cellAreaKm2(cell h3Index) (float64, h3Error) {
	out, err := cellAreaRads2(cell)
	if err == eSuccess {
		out *= earthRadiusKm * earthRadiusKm
	}
	return out, err
}
