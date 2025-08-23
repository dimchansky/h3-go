package c2go

import "math"

// bboxCenter returns the center of a bbox.
// Handles transmeridian bboxes by adding 2π to east longitude before averaging.
// Ported from H3 C: bbox.c::bboxCenter
func bboxCenter(b *BBox) LatLng {
	center := LatLng{}
	center.Lat = (b.North + b.South) * 0.5
	east := b.East
	if bboxIsTransmeridian(b) {
		east += 2 * math.Pi
	}
	center.Lng = constrainLng((east + b.West) * 0.5)
	return center
}
