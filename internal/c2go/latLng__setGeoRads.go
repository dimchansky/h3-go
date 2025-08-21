package c2go

// _setGeoRads sets the components of a LatLng in radians.
// Ported from H3 C: latLng.c::_setGeoRads
func _setGeoRads(p *LatLng, latRads, lngRads float64) {
    p.Lat = latRads
    p.Lng = lngRads
}

