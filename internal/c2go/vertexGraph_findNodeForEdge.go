package c2go

// findNodeForEdge finds the Vertex node for a given edge, if it exists.
// Searches the graph for an edge from fromVtx to toVtx. If toVtx is nil,
// matches any edge starting from fromVtx (wildcard search).
// Returns the matching VertexNode or nil if not found.
// Ported from H3 C: vertexGraph.c::findNodeForEdge
func findNodeForEdge(graph *VertexGraph, fromVtx *LatLng, toVtx *LatLng) *VertexNode {
	// Determine location using hash function
	index := _hashVertex(fromVtx, graph.Res, graph.NumBuckets)

	// Check whether there's an existing node in that spot
	node := graph.Buckets[index]
	if node != nil {
		// Look through the list and see if we find the edge
		for node != nil {
			if geoAlmostEqual(&node.From, fromVtx) &&
				(toVtx == nil || geoAlmostEqual(&node.To, toVtx)) {
				return node
			}
			node = node.Next
		}
	}

	// Iteration lookup fail
	return nil
}
