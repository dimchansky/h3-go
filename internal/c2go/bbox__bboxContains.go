package c2go

// bboxContains reports whether bbox contains the given point. Port of bbox.c::bboxContains
// Ported from H3 C: bbox.c::bboxContains
func bboxContains(b *BBox, p *LatLng) bool {
	if !(p.Lat >= b.South && p.Lat <= b.North) {
		return false
	}
	if bboxIsTransmeridian(b) {
		return p.Lng >= b.West || p.Lng <= b.East
	}
	return p.Lng >= b.West && p.Lng <= b.East
}
