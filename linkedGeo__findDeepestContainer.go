package h3

// findDeepestContainer finds the most deeply nested container from a list of nested polygons.
// Given a list of nested containers, this function finds the one most deeply nested.
// It does this by checking each polygon to see how many containers it has within the list,
// and selecting the polygon with the maximum count (most deeply nested).
//
// Technical details:
// - Returns the first polygon if only one candidate or empty list
// - For multiple polygons, counts containers for each using countContainers
// - Selects the polygon with the highest container count (most deeply nested)
// - Used in polygon normalization to find parent polygons for holes
//
// Ported from H3 C: linkedGeo.c::findDeepestContainer (static function)
func findDeepestContainer(polygons []*LinkedGeoPolygon, bboxes []*BBox) *LinkedGeoPolygon {
	if len(polygons) != len(bboxes) {
		panic("findDeepestContainer: polygons and bboxes must have same length")
	}

	polygonCount := len(polygons)

	// Set the initial return value to the first candidate
	var parent *LinkedGeoPolygon
	if polygonCount > 0 {
		parent = polygons[0]
	}

	// If we have multiple polygons, they must be nested inside each other.
	// Find the innermost polygon by taking the one with the most containers
	// in the list.
	if polygonCount > 1 {
		maxCnt := int32(-1)
		for i := 0; i < polygonCount; i++ {
			// Count how many containers this polygon has
			count := countContainers(polygons[i].First, polygons, bboxes)
			if count > maxCnt {
				parent = polygons[i]
				maxCnt = count
			}
		}
	}

	return parent
}
