package h3

// NewLatLng returns the geographic coordinate with the given latitude and
// longitude angles.
func NewLatLng(lat, lng Angle) LatLng { return LatLng{Lat: lat, Lng: lng} }

// LatLngDegs returns the geographic coordinate for latitude and longitude
// given in degrees.
func LatLngDegs(latDegs, lngDegs float64) LatLng {
	return LatLng{Lat: Deg(latDegs), Lng: Deg(lngDegs)}
}

// Cell returns the cell containing the coordinate at the given resolution.
// It has the same contract as LatLngToCell: non-finite coordinates fail
// with ErrLatLngDomain, res outside 0..MaxResolution fails with
// ErrResolutionDomain, and finite out-of-range coordinates are not
// validated (the result is unspecified).
//
// H3 C API: latLngToCell.
func (g LatLng) Cell(res int) (Cell, error) { return LatLngToCell(g, res) }
