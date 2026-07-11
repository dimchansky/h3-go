package h3

import "fmt"

// geoToStringRads converts a LatLng coordinate to string format in radians.
// Returns a string in the format "(lat, lng)" with 4 decimal places.
// Ported from H3 C: utility.c::geoToStringRads.
//
//nolint:unused // ported from H3 C for parity completeness; exercised by cgo && c2go parity tests
func geoToStringRads(p *LatLng) string {
	return fmt.Sprintf("(%.4f, %.4f)", p.Lat.Rad(), p.Lng.Rad())
}
