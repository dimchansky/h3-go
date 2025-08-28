package h3

// getHexagonAreaAvgM2 gets the average area in square meters of H3 hexagons at a given resolution.
// Ported from H3 C: latLng.c::getHexagonAreaAvgM2
func getHexagonAreaAvgM2(res int32, out *float64) H3Error {
	areas := [16]float64{
		4.357449416078390e+12, 6.097884417941339e+11, 8.680178039899731e+10,
		1.239343465508818e+10, 1.770347654491309e+09, 2.529038581819452e+08,
		3.612906216441250e+07, 5.161293359717198e+06, 7.373275975944188e+05,
		1.053325134272069e+05, 1.504750190766437e+04, 2.149643129451882e+03,
		3.070918756316063e+02, 4.387026794728301e+01, 6.267181135324322e+00,
		8.953115907605802e-01,
	}
	if res < 0 || res > MAX_H3_RES {
		return E_RES_DOMAIN
	}
	*out = areas[res]
	return E_SUCCESS
}
