package h3

// bboxesFromGeoPolygon creates bounding boxes from a GeoPolygon.
// The function creates one bounding box for the outer loop and one for each hole.
// The bboxes array must be pre-allocated with length 1 + len(polygon.Holes).
// Ported from H3 C: polygon.c::bboxesFromGeoPolygon.
func bboxesFromGeoPolygon(polygon *GeoPolygon, bboxes []BBox) {
	// Create bounding box for the main geoloop
	bboxFromGeoLoop(polygon.GeoLoop, &bboxes[0])

	// Create bounding boxes for each hole
	for i := 0; i < len(polygon.Holes); i++ {
		bboxFromGeoLoop(polygon.Holes[i], &bboxes[i+1])
	}
}
