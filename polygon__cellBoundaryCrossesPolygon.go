package h3

// cellBoundaryCrossesPolygon reports whether any part of a cell boundary crosses a polygon.
// Ported from polygon.c::cellBoundaryCrossesPolygon
// Ported from H3 C: polygon.c::cellBoundaryCrossesPolygon.
func cellBoundaryCrossesPolygon(poly GeoPolygon, bboxes []bbox, boundary *CellBoundary, boundaryBBox *bbox) bool {
	if cellBoundaryCrossesGeoLoop(poly.GeoLoop, &bboxes[0], boundary, boundaryBBox) {
		return true
	}
	for i := 0; i < len(poly.Holes); i++ {
		if cellBoundaryCrossesGeoLoop(poly.Holes[i], &bboxes[i+1], boundary, boundaryBBox) {
			return true
		}
	}
	return false
}
