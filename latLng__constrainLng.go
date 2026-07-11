package h3

// constrainLng makes sure longitudes are in the proper bounds.
// Ported from H3 C: latLng.c::constrainLng.
func constrainLng(lng Angle) Angle {
	for lng > Pi {
		lng = lng - TwoPi
	}
	for lng < -Pi {
		lng = lng + TwoPi
	}
	return lng
}
