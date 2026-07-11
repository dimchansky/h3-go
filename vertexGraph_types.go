package h3

// vertexNode mirrors the C struct used in vertexGraph.h.
type vertexNode struct {
	From LatLng
	To   LatLng
	Next *vertexNode
}

// vertexGraph mirrors the C struct used in vertexGraph.h.
type vertexGraph struct {
	Buckets    []*vertexNode
	NumBuckets int32
	Size       int32
	Res        int32
}
