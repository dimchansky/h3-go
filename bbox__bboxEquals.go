package h3

// bboxEquals reports strict equality of bounding boxes. Port of bbox.c::bboxEquals
// Ported from H3 C: bbox.c::bboxEquals.
func bboxEquals(b1, b2 *bbox) bool {
	return b1.North == b2.North && b1.South == b2.South && b1.East == b2.East && b1.West == b2.West
}
