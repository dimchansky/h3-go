package h3

// removeVertexNode removes a node from the graph. The input node will be freed,
// and should not be used after removal.
// Returns 0 on success, 1 on failure (node not found).
// The function searches for the specified node in the appropriate hash bucket
// and removes it from the linked list structure. Once removed, the node is
// no longer valid and the graph size is decremented.
// Ported from H3 C: vertexGraph.c::removeVertexNode.
func removeVertexNode(graph *vertexGraph, node *vertexNode) int32 {
	// Determine location using hash of the node's from vertex
	index := _hashVertex(&node.From, graph.Res, graph.NumBuckets)
	currentNode := graph.Buckets[index]
	found := false

	if currentNode != nil {
		if currentNode == node {
			// Node is the first in the bucket
			graph.Buckets[index] = node.Next
			found = true
		} else {
			// Look through the list to find the node
			for !found && currentNode.Next != nil {
				if currentNode.Next == node {
					// Splice the node out of the list
					currentNode.Next = node.Next
					found = true
				} else {
					currentNode = currentNode.Next
				}
			}
		}
	}

	if found {
		// In Go, we don't need to explicitly free memory as the GC handles it
		// Just remove references and decrement size
		node.Next = nil // Clear the reference to help GC
		graph.Size--
		return 0 // Success
	}

	// Failed to find the node
	return 1
}
