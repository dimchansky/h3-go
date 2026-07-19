package h3

// addLinkedGeoLoop appends a new LinkedGeoLoop populated from a GeoLoop
// to a LinkedGeoPolygon.
// Ported from H3 C: linkedGeo.c::addLinkedGeoLoop.
func addLinkedGeoLoop(gl GeoLoop, currentPoly *linkedGeoPolygon) h3Error {
	loop := &linkedGeoLoop{}

	if currentPoly.Last != nil {
		currentPoly.Last.Next = loop
	} else {
		currentPoly.First = loop
	}
	currentPoly.Last = loop

	return geoLoopToLinkedGeoLoop(gl, loop)
}
