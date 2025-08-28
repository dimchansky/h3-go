package h3

// initVertexGraph initializes a new VertexGraph.
// Allocates memory for buckets if numBuckets > 0, otherwise sets buckets to nil.
// Sets the graph's metadata fields (numBuckets, size, resolution).
// Ported from H3 C: vertexGraph.c::initVertexGraph
func initVertexGraph(graph *VertexGraph, numBuckets int32, res int32) {
	if numBuckets > 0 {
		graph.Buckets = make([]*VertexNode, numBuckets)
	} else {
		graph.Buckets = nil
	}
	graph.NumBuckets = numBuckets
	graph.Size = 0
	graph.Res = res
}
