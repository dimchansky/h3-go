package c2go

// bboxCenter returns the center of a bbox. Port of bbox.c::bboxCenter
func bboxCenter(b BBox) LatLng {
	center := LatLng{}
	center.Lat = (b.North + b.South) * 0.5
	east := b.East
	if bboxIsTransmeridian(b) {
		east += 2 * 3.14159265358979323846
	}
	center.Lng = constrainLng((east + b.West) * 0.5)
	return center
}
