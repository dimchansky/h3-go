package h3

// getHexagonAreaAvgKm2 gets the average area in square kilometers of H3 hexagons at a given resolution.
// Ported from H3 C: latLng.c::getHexagonAreaAvgKm2
func getHexagonAreaAvgKm2(res int32, out *float64) H3Error {
	areas := [16]float64{
		4.357449416078383e+06, 6.097884417941332e+05, 8.680178039899720e+04,
		1.239343465508816e+04, 1.770347654491307e+03, 2.529038581819449e+02,
		3.612906216441245e+01, 5.161293359717191e+00, 7.373275975944177e-01,
		1.053325134272067e-01, 1.504750190766435e-02, 2.149643129451879e-03,
		3.070918756316060e-04, 4.387026794728296e-05, 6.267181135324313e-06,
		8.953115907605790e-07,
	}
	if res < 0 || res > MAX_H3_RES {
		return E_RES_DOMAIN
	}
	*out = areas[res]
	return E_SUCCESS
}
