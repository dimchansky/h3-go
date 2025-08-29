package h3

// bboxToCellBoundary converts a bounding box to a cell boundary with 4 vertices
// in counter-clockwise order.
// Ported from H3 C: bbox.c::bboxToCellBoundary
func bboxToCellBoundary(bbox *BBox) CellBoundary {
	// Convert bbox to cell boundary, CCW vertex order
	return CellBoundary{
		NumVerts: 4,
		Verts: []LatLng{
			{Lat: bbox.North, Lng: bbox.East},
			{Lat: bbox.North, Lng: bbox.West},
			{Lat: bbox.South, Lng: bbox.West},
			{Lat: bbox.South, Lng: bbox.East},
		},
	}
}
