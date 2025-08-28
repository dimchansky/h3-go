package h3

// bboxIsTransmeridian returns whether the bbox crosses the antimeridian.
// Ported from H3 C: bbox.c::bboxIsTransmeridian
func bboxIsTransmeridian(b *BBox) bool { return b.East < b.West }
