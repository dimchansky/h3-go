package h3

// _setGeoRads sets the components of a LatLng in radians.
// Ported from H3 C: latLng.c::_setGeoRads
func _setGeoRads(p *LatLng, latRads, lngRads float64) {
	p.Lat = Rad(latRads)
	p.Lng = Rad(lngRads)
}
