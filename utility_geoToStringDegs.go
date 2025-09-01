package h3

import "fmt"

// geoToStringDegs converts a LatLng coordinate to string format in degrees.
// Returns a string in the format "(lat, lng)" with 9 decimal places.
// Ported from H3 C: utility.c::geoToStringDegs
func geoToStringDegs(p *LatLng) string {
	return fmt.Sprintf("(%.9f, %.9f)", p.Lat.Deg(), p.Lng.Deg())
}
