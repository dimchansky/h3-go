package h3

// cellAreaM2 computes the area of an H3 cell in meters^2.
// Ported from H3 C: area.c::cellAreaM2.
func cellAreaM2(cell h3Index) (float64, h3Error) {
	out, err := cellAreaKm2(cell)
	if err == eSuccess {
		out *= 1000 * 1000
	}
	return out, err
}
