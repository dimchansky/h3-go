package h3

// cellBoundaryCrossesGeoLoop reports whether any segment of the boundary crosses the loop.
// Ported from polygon.c::cellBoundaryCrossesGeoLoop
// Ported from H3 C: polygon.c::cellBoundaryCrossesGeoLoop.
func cellBoundaryCrossesGeoLoop(geoloop GeoLoop, loopBBox *bbox, boundary *CellBoundary, boundaryBBox *bbox) bool {
	if !bboxOverlapsBBox(loopBBox, boundaryBBox) {
		return false
	}
	var loopNorm, boundNorm longitudeNormalization
	bboxNormalization(loopBBox, boundaryBBox, &loopNorm, &boundNorm)
	// Normalize boundary longitudes into a stack copy (C copies the struct).
	normalBoundary := *boundary
	for i := int32(0); i < normalBoundary.numVerts; i++ {
		normalBoundary.verts[i].Lng = normalizeLng(normalBoundary.verts[i].Lng, boundNorm)
	}
	normalBoundaryBBox := bbox{
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
		for j := int32(0); j < normalBoundary.numVerts; j++ {
			a := normalBoundary.verts[j]
			b := normalBoundary.verts[(j+1)%normalBoundary.numVerts]
			if lineCrossesLine(&loop1, &loop2, &a, &b) {
				return true
			}
		}
	}
	return false
}
