package h3

// _vertexGraphToLinkedGeo creates a linkedGeoPolygon from a vertex graph.
// It walks the graph extracting edges into closed loops. As edges are consumed
// from the graph (removed via removeVertexNode), connected edges form loops.
// Each loop becomes a linkedGeoLoop in the output linkedGeoPolygon. The function
// continues until all edges have been consumed from the graph.
// Ported from H3 C: algos.c::_vertexGraphToLinkedGeo.
func _vertexGraphToLinkedGeo(graph *vertexGraph, out *linkedGeoPolygon) {
	// Initialize the output polygon to empty
	*out = linkedGeoPolygon{}

	var loop *linkedGeoLoop
	var edge *vertexNode
	var nextVtx LatLng

	// Find the next unused entry point
	for edge = firstVertexNode(graph); edge != nil; edge = firstVertexNode(graph) {
		loop = addNewLinkedLoop(out)
		// Walk the graph to get the outline
		for {
			addLinkedCoord(loop, &edge.From)
			nextVtx = edge.To
			// Remove frees the node, so we can't use edge after this
			removeVertexNode(graph, edge)
			edge = findNodeForVertex(graph, &nextVtx)
			if edge == nil {
				break
			}
		}
	}
}
