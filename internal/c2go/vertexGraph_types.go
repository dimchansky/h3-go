package c2go

// VertexNode mirrors the C struct used in vertexGraph.h
type VertexNode struct {
	From LatLng
	To   LatLng
	Next *VertexNode
}

// VertexGraph mirrors the C struct used in vertexGraph.h
type VertexGraph struct {
	Buckets    []*VertexNode
	NumBuckets int
	Size       int
	Res        int
}
