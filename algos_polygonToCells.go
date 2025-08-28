package h3

// polygonToCells takes a given GeoPolygon and preallocated memory and fills it with
// the hexagons that are contained by the GeoPolygon.
//
// This implementation traces the GeoPolygon geoloop(s) in cartesian space with
// hexagons, tests them and their neighbors to be contained by the geoloop(s),
// and then any newly found hexagons are used to test again until no new
// hexagons are found.
//
// The algorithm follows a point-in-polygon approach for center points to ensure
// that adjacent polygons with zero overlap have zero overlapping hexagons.
// Ported from H3 C: algos.c::polygonToCells
func polygonToCells(geoPolygon *GeoPolygon, res int32, flags uint32, out []H3Index) H3Error {
	flagErr := validatePolygonFlags(flags)
	if flagErr != E_SUCCESS {
		return flagErr
	}

	// Get the bounding boxes for the polygon and any holes
	bboxes := make([]BBox, len(geoPolygon.Holes)+1)
	bboxesFromGeoPolygon(geoPolygon, bboxes)

	// Get the estimated number of hexagons and allocate temporary memory
	var numHexagons int64
	numHexagonsError := maxPolygonToCellsSize(geoPolygon, res, flags, &numHexagons)
	if numHexagonsError != E_SUCCESS {
		return numHexagonsError
	}

	search := make([]H3Index, numHexagons)
	found := make([]H3Index, numHexagons)

	// Clear out array
	for i := range out {
		out[i] = H3_NULL
	}

	// Some metadata for tracking the state of the search and found memory blocks
	var numSearchHexes int64
	var numFoundHexes int64

	// 1. Trace the hexagons along the polygon defining the outer geoloop and
	// add them to the search hash. The hexagon containing the geoloop point
	// may or may not be contained by the geoloop (as the hexagon's center
	// point may be outside of the boundary.)
	geoloop := geoPolygon.GeoLoop
	edgeHexError := _getEdgeHexagons(geoloop, numHexagons, res,
		&numSearchHexes, search, found)
	if edgeHexError != E_SUCCESS {
		return edgeHexError
	}

	// 2. Iterate over all holes, trace the polygons defining the holes with
	// hexagons and add to only the search hash. We're going to temporarily use
	// the `found` hash to use for dedupe purposes and then re-zero it once
	// we're done here, otherwise we'd have to scan the whole set on each insert
	// to make sure there's no duplicates, which is very inefficient.
	for i := 0; i < len(geoPolygon.Holes); i++ {
		hole := geoPolygon.Holes[i]
		edgeHexError = _getEdgeHexagons(hole, numHexagons, res, &numSearchHexes,
			search, found)
		if edgeHexError != E_SUCCESS {
			return edgeHexError
		}
	}

	// 3. Re-zero the found hash so it can be used in the main loop below
	for i := int64(0); i < numHexagons; i++ {
		found[i] = H3_NULL
	}

	// 4. Begin main loop. While the search hash is not empty do the following
	for numSearchHexes > 0 {
		// Iterate through all hexagons in the current search hash, then loop
		// through all neighbors and test Point-in-Poly, if point-in-poly
		// succeeds, add to out and found hashes if not already there.
		var currentSearchNum int64
		var i int64
		for currentSearchNum < numSearchHexes {
			ring := make([]H3Index, MAX_ONE_RING_SIZE)
			searchHex := search[i]
			gridDiskErr := gridDisk(searchHex, 1, ring)
			if gridDiskErr != E_SUCCESS {
				return gridDiskErr
			}

			for j := 0; j < MAX_ONE_RING_SIZE; j++ {
				if ring[j] == H3_NULL {
					continue // Skip if this was a pentagon and only had 5 neighbors
				}

				hex := ring[j]

				// A simple hash to store the hexagon, or move to another place
				// if needed. This MUST be done before the point-in-poly check
				// since that's far more expensive
				loc := int64(hex % H3Index(numHexagons))
				var loopCount int64
				for out[loc] != H3_NULL {
					// If this branch is reached, we have exceeded the maximum
					// number of hexagons possible
					if loopCount > numHexagons {
						return E_FAILED
					}
					if out[loc] == hex {
						break // Skip duplicates found
					}
					loc = (loc + 1) % numHexagons
					loopCount++
				}
				if out[loc] == hex {
					continue // Skip this hex, already exists in the out hash
				}

				// Check if the hexagon is in the polygon or not
				var hexCenter LatLng
				cellToLatLngErr := cellToLatLng(hex, &hexCenter)
				if cellToLatLngErr != E_SUCCESS {
					return cellToLatLngErr
				}

				// If not, skip
				if !pointInsidePolygon(*geoPolygon, bboxes, &hexCenter) {
					continue
				}

				// Otherwise set it in the output array
				out[loc] = hex

				// Set the hexagon in the found hash
				found[numFoundHexes] = hex
				numFoundHexes++
			}
			currentSearchNum++
			i++
		}

		// Swap the search and found pointers, copy the found hex count to the
		// search hex count, and zero everything related to the found memory.
		search, found = found, search
		for j := int64(0); j < numSearchHexes; j++ {
			found[j] = H3_NULL
		}
		numSearchHexes = numFoundHexes
		numFoundHexes = 0
		// Repeat until no new hexagons are found
	}

	return E_SUCCESS
}
