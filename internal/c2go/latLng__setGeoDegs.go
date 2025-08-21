package c2go

// setGeoDegs sets the components of a LatLng in decimal degrees (converted to radians).
// Ported from H3 C: latLng.c::setGeoDegs
func setGeoDegs(p *LatLng, latDegs, lngDegs float64) {
    _setGeoRads(p, degsToRads(latDegs), degsToRads(lngDegs))
}

