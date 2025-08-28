package h3

import "math"

// lineHexEstimate returns an estimated number of hexagons that trace
// the cartesian-projected line.
// Ported from H3 C: bbox.c::lineHexEstimate
func lineHexEstimate(origin *LatLng, destination *LatLng, res int32, out *int64) H3Error {
	// Get the area of the pentagon as the maximally-distorted area possible
	var pentagons = make([]H3Index, NUM_PENTAGONS)
	pentagonsErr := getPentagons(res, pentagons)
	if pentagonsErr != E_SUCCESS {
		return pentagonsErr
	}
	pentagonRadiusKm := _hexRadiusKm(pentagons[0])

	dist := greatCircleDistanceKm(origin, destination)
	distCeil := math.Ceil(dist / (2 * pentagonRadiusKm))
	if math.IsInf(distCeil, 0) || math.IsNaN(distCeil) {
		return E_FAILED
	}
	estimate := int64(distCeil)
	if estimate == 0 {
		estimate = 1
	}
	*out = estimate
	return E_SUCCESS
}
