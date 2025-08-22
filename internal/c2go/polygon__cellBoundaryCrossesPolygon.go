package c2go

// cellBoundaryCrossesPolygon reports whether any part of a cell boundary crosses a polygon.
// Ported from polygon.c::cellBoundaryCrossesPolygon
func cellBoundaryCrossesPolygon(poly GeoPolygon, bboxes []BBox, boundary CellBoundary, boundaryBBox BBox) bool {
    if cellBoundaryCrossesGeoLoop(poly.Geoloop, bboxes[0], boundary, boundaryBBox) {
        return true
    }
    for i := 0; i < len(poly.Holes); i++ {
        if cellBoundaryCrossesGeoLoop(poly.Holes[i], bboxes[i+1], boundary, boundaryBBox) {
            return true
        }
    }
    return false
}

