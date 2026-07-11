package h3

// findNodeForVertex finds a Vertex node starting at the given vertex.
// Searches the graph for any edge starting from the specified fromVtx vertex.
// This is a convenience wrapper around findNodeForEdge with toVtx set to nil.
// Returns the first matching vertexNode or nil if not found.
// Ported from H3 C: vertexGraph.c::findNodeForVertex.
func findNodeForVertex(graph *vertexGraph, fromVtx *LatLng) *vertexNode {
	return findNodeForEdge(graph, fromVtx, nil)
}
