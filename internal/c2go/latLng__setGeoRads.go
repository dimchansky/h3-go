package c2go

import "github.com/dimchansky/h3-go/angle"

// _setGeoRads sets the components of a LatLng in radians.
// Ported from H3 C: latLng.c::_setGeoRads
func _setGeoRads(p *LatLng, latRads, lngRads float64) {
	p.Lat = angle.Rad(latRads)
	p.Lng = angle.Rad(lngRads)
}
