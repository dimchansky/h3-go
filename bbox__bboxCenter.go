package h3

// bboxCenter returns the center of a bbox.
// Handles transmeridian bboxes by adding 2π to east longitude before averaging.
// Ported from H3 C: bbox.c::bboxCenter
func bboxCenter(b *BBox) LatLng {
	center := LatLng{}
	// Average latitude directly with Angle operations
	center.Lat = (b.North + b.South).Mul(0.5)

	east := b.East
	if bboxIsTransmeridian(b) {
		east += TwoPi
	}
	// Average longitude and wrap to [-π, π]
	center.Lng = constrainLng((east + b.West).Mul(0.5))
	return center
}
