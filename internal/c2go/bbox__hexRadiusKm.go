package c2go

// _hexRadiusKm returns the radius of a given hexagon in Km.
// There is probably a cheaper way to determine the radius of a
// hexagon, but this way is conceptually simple.
// Ported from H3 C: bbox.c::_hexRadiusKm
func _hexRadiusKm(h3Index H3Index) float64 {
	var h3Center LatLng
	var h3Boundary CellBoundary
	cellToLatLng(h3Index, &h3Center)
	cellToBoundary(h3Index, &h3Boundary)
	return greatCircleDistanceKm(&h3Center, &h3Boundary.Verts[0])
}
