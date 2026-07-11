package h3

// originToDirectedEdges provides all of the directed edges from the current h3Index.
//
// This function generates the six possible directed edges from a hexagon cell,
// with special handling for pentagon cells which have only 5 valid edges.
// For pentagons, the K-axis direction (index 0) is set to h3Null.
//
// The edges slice must have length at least 6. Results are stored in canonical order:
// - edges[0]: direction 1 (or h3Null for pentagons)
// - edges[1]: direction 2
// - edges[2]: direction 3
// - edges[3]: direction 4
// - edges[4]: direction 5
// - edges[5]: direction 6
//
// Ported from H3 C: directedEdge.c::originToDirectedEdges.
func originToDirectedEdges(origin h3Index, edges []h3Index) h3Error {
	// Determine if the origin is a pentagon and special treatment needed.
	isPent := isPentagon(origin)

	// This is actually quite simple. Just modify the bits of the origin
	// slightly for each direction, except the 'k' direction in pentagons,
	// which is zeroed.
	for i := int32(0); i < 6; i++ {
		if isPent && i == 0 {
			edges[i] = h3Null
		} else {
			edges[i] = origin
			edges[i] = setMode(edges[i], h3DirectededgeMode)
			edges[i] = setReservedBits(edges[i], i+1)
		}
	}
	return eSuccess
}
