package h3

// firstVertexNode returns the first vertex node found in the graph.
// Iterates through buckets sequentially until finding a non-empty bucket,
// then returns the first node in that bucket. Returns nil if no nodes exist.
// Ported from H3 C: vertexGraph.c::firstVertexNode
func firstVertexNode(graph *VertexGraph) *VertexNode {
	var node *VertexNode
	currentIndex := int32(0)

	for node == nil {
		if currentIndex < graph.NumBuckets {
			// find the first node in the next bucket
			node = graph.Buckets[currentIndex]
		} else {
			// end of iteration
			return nil
		}
		currentIndex++
	}
	return node
}
