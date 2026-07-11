package h3

import (
	"math"
)

// bboxHexEstimate returns an estimated number of hexagons that fit
// within the cartesian-projected bounding box.
// Ported from H3 C: bbox.c::bboxHexEstimate.
func bboxHexEstimate(bbox *bbox, res int32, out *int64) h3Error {
	// Get the area of the pentagon as the maximally-distorted area possible
	var pentagons = make([]h3Index, numPentagons)
	pentagonsErr := getPentagons(res, pentagons)
	if pentagonsErr != eSuccess {
		return pentagonsErr
	}
	pentagonRadiusKm := _hexRadiusKm(pentagons[0])

	// Area of a regular hexagon is 3/2*sqrt(3) * r * r
	// The pentagon has the most distortion (smallest edges) and shares its
	// edges with hexagons, so the most-distorted hexagons have this area,
	// shrunk by 20% off chance that the bounding box perfectly bounds a
	// pentagon.
	pentagonAreaKm2 := 0.8 * (2.59807621135 * pentagonRadiusKm * pentagonRadiusKm)

	// Then get the area of the bounding box of the geoloop in question
	var p1, p2 LatLng
	p1.Lat = bbox.North
	p1.Lng = bbox.East
	p2.Lat = bbox.South
	p2.Lng = bbox.West
	d := greatCircleDistanceKm(&p1, &p2)
	lngDiff := (p1.Lng - p2.Lng).Abs().Rad()
	latDiff := (p1.Lat - p2.Lat).Abs().Rad()
	if lngDiff == 0 || latDiff == 0 {
		return eFailed
	}
	length := math.Max(lngDiff, latDiff)
	width := math.Min(lngDiff, latDiff)
	ratio := length / width
	// Derived constant based on: https://math.stackexchange.com/a/1921940
	// Clamped to 3 as higher values tend to rapidly drag the estimate to zero.
	a := d * d / math.Min(3.0, ratio)

	// Divide the two to get an estimate of the number of hexagons needed
	estimateDouble := math.Ceil(a / pentagonAreaKm2)
	if math.IsInf(estimateDouble, 0) || math.IsNaN(estimateDouble) {
		return eFailed
	}
	estimate := int64(estimateDouble)
	if estimate == 0 {
		estimate = 1
	}
	*out = estimate
	return eSuccess
}
