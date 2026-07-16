package h3

// AreaRads2 returns the exact spherical area of the cell in square
// radians, computed by triangulating the cell's actual boundary (including
// any distortion vertices). For the resolution-wide hexagon average, see
// HexagonAreaAvgKm2.
//
// H3 C API: cellAreaRads2.
func (c Cell) AreaRads2() (float64, error) {
	a, errC := cellAreaRads2(c)
	if errC != eSuccess {
		return 0, toErr(errC)
	}
	return a, nil
}

// AreaKm2 returns the exact spherical area of the cell in square
// kilometers, computed by triangulating the cell's actual boundary
// (including any distortion vertices). The result is a spherical
// approximation on the WGS84 authalic sphere (radius 6371.007180918475
// km), not an ellipsoidal geodesic value. For the resolution-wide hexagon
// average, see HexagonAreaAvgKm2.
//
// H3 C API: cellAreaKm2.
func (c Cell) AreaKm2() (float64, error) {
	a, errC := cellAreaKm2(c)
	if errC != eSuccess {
		return 0, toErr(errC)
	}
	return a, nil
}

// AreaM2 returns the exact spherical area of the cell in square meters;
// see AreaKm2 for the computation method and earth model.
//
// H3 C API: cellAreaM2.
func (c Cell) AreaM2() (float64, error) {
	a, errC := cellAreaM2(c)
	if errC != eSuccess {
		return 0, toErr(errC)
	}
	return a, nil
}

// LengthRads returns the exact spherical length of the edge in radians,
// computed by summing great-circle segments along the edge's boundary
// (including any distortion vertex). For the resolution-wide hexagon
// average, see HexagonEdgeLengthAvgKm.
//
// H3 C API: edgeLengthRads.
func (e DirectedEdge) LengthRads() (float64, error) {
	var l float64
	if errC := edgeLengthRads(h3Index(e), &l); errC != eSuccess {
		return 0, toErr(errC)
	}
	return l, nil
}

// LengthKm returns the exact spherical length of the edge in kilometers,
// computed by summing great-circle segments along the edge's boundary
// (including any distortion vertex). The result is a spherical
// approximation on the WGS84 authalic sphere (radius 6371.007180918475
// km), not an ellipsoidal geodesic value. For the resolution-wide hexagon
// average, see HexagonEdgeLengthAvgKm.
//
// H3 C API: edgeLengthKm.
func (e DirectedEdge) LengthKm() (float64, error) {
	var l float64
	if errC := edgeLengthKm(h3Index(e), &l); errC != eSuccess {
		return 0, toErr(errC)
	}
	return l, nil
}

// LengthM returns the exact spherical length of the edge in meters; see
// LengthKm for the computation method and earth model.
//
// H3 C API: edgeLengthM.
func (e DirectedEdge) LengthM() (float64, error) {
	var l float64
	if errC := edgeLengthM(h3Index(e), &l); errC != eSuccess {
		return 0, toErr(errC)
	}
	return l, nil
}

// HexagonAreaAvgKm2 returns the average hexagon area in square kilometers
// at the given resolution (pentagons excluded). The average assumes the
// same spherical earth model as Cell.AreaKm2 (WGS84 authalic sphere,
// radius 6371.007180918475 km); for the exact area of a specific cell, use
// Cell.AreaKm2.
//
// H3 C API: getHexagonAreaAvgKm2.
func HexagonAreaAvgKm2(res int) (float64, error) {
	if err := checkRes(res); err != nil {
		return 0, err
	}
	var a float64
	if errC := getHexagonAreaAvgKm2(int32(res), &a); errC != eSuccess {
		return 0, toErr(errC)
	}
	return a, nil
}

// HexagonAreaAvgM2 returns the average hexagon area in square meters at
// the given resolution (pentagons excluded); see HexagonAreaAvgKm2. For
// the exact area of a specific cell, use Cell.AreaM2.
//
// H3 C API: getHexagonAreaAvgM2.
func HexagonAreaAvgM2(res int) (float64, error) {
	if err := checkRes(res); err != nil {
		return 0, err
	}
	var a float64
	if errC := getHexagonAreaAvgM2(int32(res), &a); errC != eSuccess {
		return 0, toErr(errC)
	}
	return a, nil
}

// HexagonEdgeLengthAvgKm returns the average hexagon edge length in
// kilometers at the given resolution (pentagons excluded). The average
// assumes the same spherical earth model as DirectedEdge.LengthKm (WGS84
// authalic sphere, radius 6371.007180918475 km); for the exact length of a
// specific edge, use DirectedEdge.LengthKm.
//
// H3 C API: getHexagonEdgeLengthAvgKm.
func HexagonEdgeLengthAvgKm(res int) (float64, error) {
	if err := checkRes(res); err != nil {
		return 0, err
	}
	var l float64
	if errC := getHexagonEdgeLengthAvgKm(int32(res), &l); errC != eSuccess {
		return 0, toErr(errC)
	}
	return l, nil
}

// HexagonEdgeLengthAvgM returns the average hexagon edge length in meters
// at the given resolution (pentagons excluded); see HexagonEdgeLengthAvgKm.
// For the exact length of a specific edge, use DirectedEdge.LengthM.
//
// H3 C API: getHexagonEdgeLengthAvgM.
func HexagonEdgeLengthAvgM(res int) (float64, error) {
	if err := checkRes(res); err != nil {
		return 0, err
	}
	var l float64
	if errC := getHexagonEdgeLengthAvgM(int32(res), &l); errC != eSuccess {
		return 0, toErr(errC)
	}
	return l, nil
}

// GreatCircleDistanceRads returns the great-circle distance between the
// two coordinates in radians, computed with the haversine formula.
//
// H3 C API: greatCircleDistanceRads.
func GreatCircleDistanceRads(a, b LatLng) float64 { return greatCircleDistanceRads(&a, &b) }

// GreatCircleDistanceKm returns the great-circle distance between the two
// coordinates in kilometers, computed with the haversine formula on the
// WGS84 authalic sphere (radius 6371.007180918475 km) — a spherical
// approximation, not an ellipsoidal geodesic distance.
//
// H3 C API: greatCircleDistanceKm.
func GreatCircleDistanceKm(a, b LatLng) float64 { return greatCircleDistanceKm(&a, &b) }

// GreatCircleDistanceM returns the great-circle distance between the two
// coordinates in meters; see GreatCircleDistanceKm for the formula and
// earth model.
//
// H3 C API: greatCircleDistanceM.
func GreatCircleDistanceM(a, b LatLng) float64 { return greatCircleDistanceM(&a, &b) }
