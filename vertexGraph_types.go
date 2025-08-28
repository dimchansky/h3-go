package h3

// VertexNode mirrors the C struct used in vertexGraph.h
type VertexNode struct {
	From LatLng
	To   LatLng
	Next *VertexNode
}

// VertexGraph mirrors the C struct used in vertexGraph.h
type VertexGraph struct {
	Buckets    []*VertexNode
	NumBuckets int32
	Size       int32
	Res        int32
}
