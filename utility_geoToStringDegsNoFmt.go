package h3

import "fmt"

// geoToStringDegsNoFmt converts a LatLng coordinate to string format in degrees without parentheses.
// Returns a string in the format "lat lng" with 9 decimal places.
// Ported from H3 C: utility.c::geoToStringDegsNoFmt
func geoToStringDegsNoFmt(p *LatLng) string {
	return fmt.Sprintf("%.9f %.9f", p.Lat.Deg(), p.Lng.Deg())
}
