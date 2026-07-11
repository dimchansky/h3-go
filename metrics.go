package h3

// AreaRads2 returns the exact spherical area of the cell in square radians.
//
// H3 C API: cellAreaRads2.
func (c Cell) AreaRads2() (float64, error) {
	a, errC := cellAreaRads2(c)
	if errC != eSuccess {
		return 0, toErr(errC)
	}
	return a, nil
}

// AreaKm2 returns the exact spherical area of the cell in square kilometers.
//
// H3 C API: cellAreaKm2.
func (c Cell) AreaKm2() (float64, error) {
	a, errC := cellAreaKm2(c)
	if errC != eSuccess {
		return 0, toErr(errC)
	}
	return a, nil
}

// AreaM2 returns the exact spherical area of the cell in square meters.
//
// H3 C API: cellAreaM2.
func (c Cell) AreaM2() (float64, error) {
	a, errC := cellAreaM2(c)
	if errC != eSuccess {
		return 0, toErr(errC)
	}
	return a, nil
}

// LengthRads returns the exact spherical length of the edge in radians.
//
// H3 C API: edgeLengthRads.
func (e DirectedEdge) LengthRads() (float64, error) {
	var l float64
	if errC := edgeLengthRads(h3Index(e), &l); errC != eSuccess {
		return 0, toErr(errC)
	}
	return l, nil
}

// LengthKm returns the exact spherical length of the edge in kilometers.
//
// H3 C API: edgeLengthKm.
func (e DirectedEdge) LengthKm() (float64, error) {
	var l float64
	if errC := edgeLengthKm(h3Index(e), &l); errC != eSuccess {
		return 0, toErr(errC)
	}
	return l, nil
}

// LengthM returns the exact spherical length of the edge in meters.
//
// H3 C API: edgeLengthM.
func (e DirectedEdge) LengthM() (float64, error) {
	var l float64
	if errC := edgeLengthM(h3Index(e), &l); errC != eSuccess {
		return 0, toErr(errC)
	}
	return l, nil
}

// HexagonAreaAvgKm2 returns the average hexagon area in square kilometers at
// the given resolution (pentagons excluded).
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

// HexagonAreaAvgM2 returns the average hexagon area in square meters at the
// given resolution (pentagons excluded).
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
// kilometers at the given resolution (pentagons excluded).
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

// HexagonEdgeLengthAvgM returns the average hexagon edge length in meters at
// the given resolution (pentagons excluded).
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

// GreatCircleDistanceRads returns the great-circle distance between the two
// coordinates in radians.
//
// H3 C API: greatCircleDistanceRads.
func GreatCircleDistanceRads(a, b LatLng) float64 { return greatCircleDistanceRads(&a, &b) }

// GreatCircleDistanceKm returns the great-circle distance between the two
// coordinates in kilometers.
//
// H3 C API: greatCircleDistanceKm.
func GreatCircleDistanceKm(a, b LatLng) float64 { return greatCircleDistanceKm(&a, &b) }

// GreatCircleDistanceM returns the great-circle distance between the two
// coordinates in meters.
//
// H3 C API: greatCircleDistanceM.
func GreatCircleDistanceM(a, b LatLng) float64 { return greatCircleDistanceM(&a, &b) }
