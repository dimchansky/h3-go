package h3

import "math"

// maxPolygonToCellsSizeExperimental returns the maximum number of cells that could
// result from the polygonToCellsExperimental algorithm with the given polygon and
// parameters.
//
// This function provides an upper bound for memory allocation purposes. The actual
// number of cells may be significantly smaller than this estimate.
//
// Ported from H3 C: polyfill.c::maxPolygonToCellsSizeExperimental
func maxPolygonToCellsSizeExperimental(polygon *GeoPolygon, res int32, flags uint32) (int64, H3Error) {
	// Special case: 0-vertex polygon
	if len(polygon.GeoLoop) == 0 {
		return 0, E_SUCCESS
	}

	// Initialize the iterator without stepping, so we can adjust the res and
	// flags (after they are validated by the initialization) before we start
	iter := initIterPolygonCompact(polygon, res, flags)

	if iter.Error != E_SUCCESS {
		return 0, iter.Error
	}

	// Ignore the requested flags and use the faster overlapping-bbox mode
	iter.flags = uint32(CONTAINMENT_OVERLAPPING_BBOX)

	// Get a (very) rough area of the polygon bounding box
	polygonBBox := &iter.bboxes[0]
	polygonBBoxAreaKm2 :=
		bboxHeightRads(polygonBBox) * bboxWidthRads(polygonBBox) /
			math.Cos(math.Min(math.Abs(polygonBBox.North.Rad()), math.Abs(polygonBBox.South.Rad()))) *
			EARTH_RADIUS_KM * EARTH_RADIUS_KM

	// Determine the res for the size estimate, based on a (very) rough estimate
	// of the number of cells at various resolutions that would fit in the
	// polygon. All we need here is a general order of magnitude.
	for iter.res > 0 {
		var avgArea float64
		err := getHexagonAreaAvgKm2(iter.res-1, &avgArea)
		if err != E_SUCCESS {
			return 0, err
		}
		if polygonBBoxAreaKm2/avgArea <= float64(MAX_SIZE_CELL_THRESHOLD) {
			break
		}
		iter.res--
	}

	// Now run the polyfill, counting the output in the target res.
	// We have to take the first step outside the loop, to get the first
	// valid output cell
	iterStepPolygonCompact(&iter)

	var out int64 = 0
	for iter.Cell != H3_NULL {
		childrenSize, err := cellToChildrenSize(iter.Cell, res)
		if err != E_SUCCESS {
			return 0, err
		}
		out += childrenSize
		iterStepPolygonCompact(&iter)
	}

	return out, iter.Error
}
