package h3

// getHexagonEdgeLengthAvgKm gets the average edge length in kilometers of H3 hexagons at a given resolution.
// Ported from H3 C: latLng.c::getHexagonEdgeLengthAvgKm.
func getHexagonEdgeLengthAvgKm(res int32, out *float64) h3Error {
	lens := [16]float64{
		1281.256011, 483.0568391, 182.5129565, 68.97922179,
		26.07175968, 9.854090990, 3.724532667, 1.406475763,
		0.531414010, 0.200786148, 0.075863783, 0.028663897,
		0.010830188, 0.004092010, 0.001546100, 0.000584169,
	}
	if res < 0 || res > maxH3Res {
		return eResDomain
	}
	*out = lens[res]
	return eSuccess
}
