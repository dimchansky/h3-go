package c2go

// getHexagonEdgeLengthAvgM gets the average edge length in meters of H3 hexagons at a given resolution.
// Ported from H3 C: latLng.c::getHexagonEdgeLengthAvgM
func getHexagonEdgeLengthAvgM(res int32, out *float64) H3Error {
	lens := [16]float64{
		1281256.011, 483056.8391, 182512.9565, 68979.22179,
		26071.75968, 9854.090990, 3724.532667, 1406.475763,
		531.4140101, 200.7861476, 75.86378287, 28.66389748,
		10.83018784, 4.092010473, 1.546099657, 0.584168630,
	}
	if res < 0 || res > MAX_H3_RES {
		return E_RES_DOMAIN
	}
	*out = lens[res]
	return E_SUCCESS
}
