package h3

// addVertexNode adds an edge to the vertex graph.
// Returns a pointer to the new node, or an existing node if the edge already exists.
// The function creates a new vertexNode with the given from/to vertices and adds it
// to the appropriate bucket in the graph's hash table. If an identical edge already
// exists, it returns the existing node instead of creating a duplicate.
// Ported from H3 C: vertexGraph.c::addVertexNode.
func addVertexNode(graph *vertexGraph, fromVtx *LatLng, toVtx *LatLng) *vertexNode {
	// Create the new node
	node := &vertexNode{
		From: *fromVtx,
		To:   *toVtx,
		Next: nil,
	}

	// Determine location using hash function
	index := _hashVertex(fromVtx, graph.Res, graph.NumBuckets)

	// Check whether there's an existing node in that spot
	currentNode := graph.Buckets[index]
	if currentNode == nil {
		// Set bucket to the new node
		graph.Buckets[index] = node
	} else {
		// Find the end of the list
		for {
			// Check if the edge we're adding doesn't already exist
			if geoAlmostEqual(&currentNode.From, fromVtx) &&
				geoAlmostEqual(&currentNode.To, toVtx) {
				// Already exists, return existing node
				return currentNode
			}
			if currentNode.Next == nil {
				break
			}
			currentNode = currentNode.Next
		}
		// Add the new node to the end of the list
		currentNode.Next = node
	}

	graph.Size++
	return node
}
