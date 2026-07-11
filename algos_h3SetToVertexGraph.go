package h3

// h3SetToVertexGraph creates a vertex graph from a set of hexagons.
// It builds a graph where each edge of the hexagons is represented as a vertex node.
// Edges shared between adjacent hexagons are removed, leaving only the boundary edges.
// The graph is used to generate polygon outlines from sets of hexagons.
// Ported from H3 C: algos.c::h3SetToVertexGraph.
func h3SetToVertexGraph(h3Set []h3Index, numHexes int32, graph *vertexGraph) h3Error {
	if numHexes < 1 {
		// We still need to init the graph, or calls to destroyVertexGraph will fail
		initVertexGraph(graph, 0, 0)
		return eSuccess
	}

	res := getResolution(h3Set[0])
	const minBuckets = 6
	// TODO: Better way to calculate/guess?
	numBuckets := numHexes
	if numBuckets < minBuckets {
		numBuckets = minBuckets
	}

	initVertexGraph(graph, numBuckets, res)

	// Iterate through every hexagon
	for i := int32(0); i < numHexes; i++ {
		var vertices CellBoundary
		boundaryErr := cellToBoundary(h3Set[i], &vertices)
		if boundaryErr != eSuccess {
			// Destroy vertex graph as caller will not know to do so.
			destroyVertexGraph(graph)
			return boundaryErr
		}

		// Iterate through every edge
		for j := int32(0); j < vertices.numVerts; j++ {
			fromVtx := &vertices.verts[j]
			toVtx := &vertices.verts[(j+1)%vertices.numVerts]

			// If we've seen this edge already, it will be reversed
			edge := findNodeForEdge(graph, toVtx, fromVtx)
			if edge != nil {
				// If we've seen it, drop it. No edge is shared by more than 2
				// hexagons, so we'll never see it again.
				removeVertexNode(graph, edge)
			} else {
				// Add a new node for this edge
				addVertexNode(graph, fromVtx, toVtx)
			}
		}
	}

	return eSuccess
}
