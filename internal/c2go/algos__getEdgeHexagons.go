package c2go

// _getEdgeHexagons takes a given geoloop ring and traces it with hexagons,
// updating the search and found arrays. This is used for determining the
// initial hexagon set for polygon operations.
// Ported from H3 C: algos.c::_getEdgeHexagons
func _getEdgeHexagons(geoloop []LatLng, numHexagons int64, res int32,
	numSearchHexes *int64, search []H3Index, found []H3Index) H3Error {

	for i := 0; i < len(geoloop); i++ {
		origin := geoloop[i]
		var destination LatLng
		if i == len(geoloop)-1 {
			destination = geoloop[0]
		} else {
			destination = geoloop[i+1]
		}

		var numHexesEstimate int64
		estimateErr := lineHexEstimate(&origin, &destination, res, &numHexesEstimate)
		if estimateErr != E_SUCCESS {
			return estimateErr
		}

		for j := int64(0); j < numHexesEstimate; j++ {
			var interpolate LatLng
			invNumHexesEst := 1.0 / float64(numHexesEstimate)
			interpolate.Lat = origin.Lat.Mul(float64(numHexesEstimate-j)*invNumHexesEst) + (destination.Lat.Mul(float64(j) * invNumHexesEst))
			interpolate.Lng = origin.Lng.Mul(float64(numHexesEstimate-j)*invNumHexesEst) + (destination.Lng.Mul(float64(j) * invNumHexesEst))

			var pointHex H3Index
			err := latLngToCell(&interpolate, res, &pointHex)
			if err != E_SUCCESS {
				return err
			}

			// A simple hash to store the hexagon, or move to another place if needed
			loc := int64(pointHex % H3Index(numHexagons))
			loopCount := int64(0)
			for found[loc] != 0 {
				// If this conditional is reached, the found memory block is
				// too small for the given polygon. This should not happen.
				if loopCount > numHexagons {
					return E_FAILED
				}
				if found[loc] == pointHex {
					break // At least two points of the geoloop index to the same cell
				}
				loc = (loc + 1) % numHexagons
				loopCount++
			}
			if found[loc] == pointHex {
				continue // Skip this hex, already exists in the found hash
			}
			// Otherwise, set it in the found hash for now
			found[loc] = pointHex

			search[*numSearchHexes] = pointHex
			(*numSearchHexes)++
		}
	}
	return E_SUCCESS
}
