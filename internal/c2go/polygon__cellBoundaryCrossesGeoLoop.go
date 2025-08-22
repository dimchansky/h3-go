package c2go

// cellBoundaryCrossesGeoLoop reports whether any segment of the boundary crosses the loop.
// Ported from polygon.c::cellBoundaryCrossesGeoLoop
func cellBoundaryCrossesGeoLoop(geoloop GeoLoop, loopBBox BBox, boundary CellBoundary, boundaryBBox BBox) bool {
    if !bboxOverlapsBBox(loopBBox, boundaryBBox) {
        return false
    }
    loopNorm, boundNorm := bboxNormalization(loopBBox, boundaryBBox)
    // Normalize boundary longitudes
    normalBoundary := CellBoundary{NumVerts: boundary.NumVerts, Verts: make([]LatLng, boundary.NumVerts)}
    copy(normalBoundary.Verts, boundary.Verts)
    for i := 0; i < normalBoundary.NumVerts; i++ {
        normalBoundary.Verts[i].Lng = normalizeLng(normalBoundary.Verts[i].Lng, boundNorm)
    }
    normalBoundaryBBox := BBox{
        North: boundaryBBox.North,
        South: boundaryBBox.South,
        East:  normalizeLng(boundaryBBox.East, boundNorm),
        West:  normalizeLng(boundaryBBox.West, boundNorm),
    }

    for i := 0; i < len(geoloop); i++ {
        loop1 := geoloop[i]
        loop2 := geoloop[(i+1)%len(geoloop)]
        loop1.Lng = normalizeLng(loop1.Lng, loopNorm)
        loop2.Lng = normalizeLng(loop2.Lng, loopNorm)
        // quick bbox reject
        if (loop1.Lat >= normalBoundaryBBox.North && loop2.Lat >= normalBoundaryBBox.North) ||
            (loop1.Lat <= normalBoundaryBBox.South && loop2.Lat <= normalBoundaryBBox.South) ||
            (loop1.Lng <= normalBoundaryBBox.West && loop2.Lng <= normalBoundaryBBox.West) ||
            (loop1.Lng >= normalBoundaryBBox.East && loop2.Lng >= normalBoundaryBBox.East) {
            continue
        }
        for j := 0; j < normalBoundary.NumVerts; j++ {
            a := normalBoundary.Verts[j]
            b := normalBoundary.Verts[(j+1)%normalBoundary.NumVerts]
            if lineCrossesLine(loop1, loop2, a, b) {
                return true
            }
        }
    }
    return false
}

