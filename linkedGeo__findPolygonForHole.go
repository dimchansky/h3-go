package h3

// findPolygonForHole finds the polygon to which a given hole should be allocated.
// Given a hole loop and a list of potential parent polygons, this function identifies
// which polygon should contain the hole as its inner loop. The algorithm finds all
// polygons that contain the hole's first vertex, then selects the most deeply nested
// one using findDeepestContainer.
//
// Technical details:
// - Early exit if no polygons provided
// - Tests each polygon to see if it contains the hole's first vertex
// - Builds candidate list of containing polygons
// - Uses findDeepestContainer to select most deeply nested candidate
// - Returns nil if no suitable parent polygon is found
// - Used in polygon normalization to assign holes to their proper parent polygons
//
// Ported from H3 C: linkedGeo.c::findPolygonForHole (static function).
func findPolygonForHole(loop *linkedGeoLoop, polygon *linkedGeoPolygon, bboxes []*bbox) *linkedGeoPolygon {
	polygonCount := len(bboxes)

	// Early exit with no polygons
	if polygonCount == 0 {
		return nil
	}

	// Early exit if loop is nil or has no coordinates
	if loop == nil || loop.First == nil {
		return nil
	}

	// Initialize arrays for candidate loops and their bounding boxes
	var candidates []*linkedGeoPolygon
	var candidateBBoxes []*bbox

	// Find all polygons that contain the loop
	index := 0
	currentPolygon := polygon

	for currentPolygon != nil && index < polygonCount {
		// We are guaranteed not to overlap, so just test the first point
		if currentPolygon.First != nil &&
			pointInsideLinkedGeoLoop(currentPolygon.First, bboxes[index], &loop.First.Vertex) {
			candidates = append(candidates, currentPolygon)
			candidateBBoxes = append(candidateBBoxes, bboxes[index])
		}
		currentPolygon = currentPolygon.Next
		index++
	}

	// The most deeply nested container is the immediate parent
	parent := findDeepestContainer(candidates, candidateBBoxes)

	return parent
}
