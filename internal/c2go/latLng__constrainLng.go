package c2go

import "github.com/dimchansky/h3-go/angle"

// constrainLng makes sure longitudes are in the proper bounds.
// Ported from H3 C: latLng.c::constrainLng
func constrainLng(lng angle.Angle) angle.Angle {
	for lng > angle.Pi {
		lng = lng - angle.TwoPi
	}
	for lng < -angle.Pi {
		lng = lng + angle.TwoPi
	}
	return lng
}
