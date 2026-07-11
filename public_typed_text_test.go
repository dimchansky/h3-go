package h3

import (
	"errors"
	"testing"
)

func TestDirectedEdgeAndVertexTextRoundTrips(t *testing.T) {
	cell, err := ParseCell("8928308280fffff")
	if err != nil {
		t.Fatal(err)
	}

	edges, err := cell.DirectedEdges()
	if err != nil {
		t.Fatal(err)
	}
	edge := edges[0]
	edgeText, err := edge.MarshalText()
	if err != nil {
		t.Fatal(err)
	}
	var decodedEdge DirectedEdge
	if err := decodedEdge.UnmarshalText(edgeText); err != nil || decodedEdge != edge {
		t.Fatalf("directed edge text round trip = %v, %v; want %v", decodedEdge, err, edge)
	}

	vertex, err := cell.Vertex(0)
	if err != nil {
		t.Fatal(err)
	}
	parsedVertex, err := ParseVertex(vertex.String())
	if err != nil || parsedVertex != vertex {
		t.Fatalf("ParseVertex(%q) = %v, %v; want %v", vertex.String(), parsedVertex, err, vertex)
	}
	vertexText, err := vertex.MarshalText()
	if err != nil {
		t.Fatal(err)
	}
	var decodedVertex Vertex
	if err := decodedVertex.UnmarshalText(vertexText); err != nil || decodedVertex != vertex {
		t.Fatalf("vertex text round trip = %v, %v; want %v", decodedVertex, err, vertex)
	}

	if _, err := ParseVertex(cell.String()); !errors.Is(err, ErrVertexInvalid) {
		t.Fatalf("ParseVertex(cell) = %v; want ErrVertexInvalid", err)
	}
}
