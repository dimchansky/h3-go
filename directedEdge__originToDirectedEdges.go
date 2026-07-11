package h3

// originToDirectedEdges provides all of the directed edges from the current H3Index.
//
// This function generates the six possible directed edges from a hexagon cell,
// with special handling for pentagon cells which have only 5 valid edges.
// For pentagons, the K-axis direction (index 0) is set to H3_NULL.
//
// The edges slice must have length at least 6. Results are stored in canonical order:
// - edges[0]: Direction 1 (or H3_NULL for pentagons)
// - edges[1]: Direction 2
// - edges[2]: Direction 3
// - edges[3]: Direction 4
// - edges[4]: Direction 5
// - edges[5]: Direction 6
//
// Ported from H3 C: directedEdge.c::originToDirectedEdges.
func originToDirectedEdges(origin H3Index, edges []H3Index) H3Error {
	// Determine if the origin is a pentagon and special treatment needed.
	isPent := isPentagon(origin)

	// This is actually quite simple. Just modify the bits of the origin
	// slightly for each direction, except the 'k' direction in pentagons,
	// which is zeroed.
	for i := int32(0); i < 6; i++ {
		if isPent && i == 0 {
			edges[i] = H3_NULL
		} else {
			edges[i] = origin
			edges[i] = setMode(edges[i], H3_DIRECTEDEDGE_MODE)
			edges[i] = setReservedBits(edges[i], i+1)
		}
	}
	return E_SUCCESS
}
