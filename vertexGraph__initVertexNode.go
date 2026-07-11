package h3

// _initVertexNode initializes a VertexNode with the provided from and to vertices.
// Sets the node's from and to LatLng values and ensures the next pointer is nil.
// Ported from H3 C: vertexGraph.c::_initVertexNode.
//
//nolint:unused // ported from H3 C for parity completeness; exercised by cgo && c2go parity tests
func _initVertexNode(node *VertexNode, fromVtx *LatLng, toVtx *LatLng) {
	node.From = *fromVtx
	node.To = *toVtx
	node.Next = nil
}
