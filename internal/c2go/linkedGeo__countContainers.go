package c2go

// countContainers counts the number of polygons containing a given loop.
// It checks each polygon to see if the loop's first vertex is inside that polygon.
// A polygon is considered a container if:
// 1. The loop being tested is not the polygon's first loop (to avoid self-containment)
// 2. The loop's first vertex is inside the polygon's first loop
// This is used in normalizeMultiPolygon to find parent polygons for holes.
// Ported from H3 C: linkedGeo.c::countContainers (static function)
func countContainers(loop *LinkedGeoLoop, polygons []*LinkedGeoPolygon, bboxes []*BBox) int {
	if len(polygons) != len(bboxes) {
		panic("countContainers: polygons and bboxes must have same length")
	}
	
	// If loop is nil or has no coordinates, it can't be contained
	if loop == nil || loop.First == nil {
		return 0
	}
	
	containerCount := 0
	for i := 0; i < len(polygons); i++ {
		// Check that this isn't the same loop (avoid self-containment)
		// and that the loop's first vertex is inside the polygon
		if loop != polygons[i].First &&
			polygons[i] != nil &&
			polygons[i].First != nil &&
			pointInsideLinkedGeoLoop(polygons[i].First, bboxes[i], &loop.First.Vertex) {
			containerCount++
		}
	}
	return containerCount
}