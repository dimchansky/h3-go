package h3

// destroyVertexGraph destroys a VertexGraph's sub-objects, freeing their memory.
// In the C implementation, this iterates through all nodes and removes them,
// then frees the buckets array. In Go, we clear all references to help the GC
// and set the graph to an empty state.
// Ported from H3 C: vertexGraph.c::destroyVertexGraph
func destroyVertexGraph(graph *VertexGraph) {
	var node *VertexNode
	for {
		node = firstVertexNode(graph)
		if node == nil {
			break
		}
		removeVertexNode(graph, node)
	}
	// Clear the buckets slice
	graph.Buckets = nil
}
