package h3

import (
	"math"
)

// _hashVertex returns an integer hash for a lat/lng point, at a precision
// determined by the current hexagon resolution.
// Simple hash: Take the sum of the lat and lng with a precision level
// determined by the resolution, converted to int, modulo bucket count.
// Ported from H3 C: vertexGraph.c::_hashVertex
func _hashVertex(vertex *LatLng, res int32, numBuckets int32) uint32 {
	return uint32(math.Mod((vertex.Lat+vertex.Lng).Abs().Rad()*math.Pow(10, float64(15-res)), float64(numBuckets)))
}
