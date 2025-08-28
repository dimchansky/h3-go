package h3

// cellBoundaryInsidePolygon checks if a cell boundary is completely contained by a polygon.
// Ported from polygon.c::cellBoundaryInsidePolygon
// Ported from H3 C: polygon.c::cellBoundaryInsidePolygon
func cellBoundaryInsidePolygon(poly GeoPolygon, bboxes []BBox, boundary *CellBoundary, boundaryBBox *BBox) bool {
	// First test a single point (first vertex). Fails fast via bboxContains.
	if !pointInsidePolygon(poly, bboxes, &boundary.Verts[0]) {
		return false
	}
	// If outer loop crossings exist, not contained
	if cellBoundaryCrossesGeoLoop(poly.Geoloop, &bboxes[0], boundary, boundaryBBox) {
		return false
	}
	// Convert boundary to loop for point-inside checks
	boundaryLoop := GeoLoop(boundary.Verts)
	// Check hole intersections or containment
	for i := 0; i < len(poly.Holes); i++ {
		hole := poly.Holes[i]
		if len(hole) > 0 && (pointInsideGeoLoop(boundaryLoop, boundaryBBox, &hole[0]) ||
			cellBoundaryCrossesGeoLoop(hole, &bboxes[i+1], boundary, boundaryBBox)) {
			return false
		}
	}
	return true
}
