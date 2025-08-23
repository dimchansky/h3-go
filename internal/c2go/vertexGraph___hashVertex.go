package c2go

import (
	"math"
)

// _hashVertex returns an integer hash for a lat/lng point, at a precision
// determined by the current hexagon resolution.
// Simple hash: Take the sum of the lat and lng with a precision level
// determined by the resolution, converted to int, modulo bucket count.
// Ported from H3 C: vertexGraph.c::_hashVertex
func _hashVertex(vertex *LatLng, res int, numBuckets int) uint32 {
	return uint32(math.Mod(math.Abs((vertex.Lat+vertex.Lng)*math.Pow(10, float64(15-res))), float64(numBuckets)))
}
