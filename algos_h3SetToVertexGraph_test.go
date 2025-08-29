// Tests ported from testH3SetToVertexGraphInternal.c
package h3

import (
	"testing"
)

func TestH3SetToVertexGraphInternal_empty(t *testing.T) {
	t.Parallel()
	var graph VertexGraph

	err := h3SetToVertexGraph(nil, 0, &graph)
	if err != E_SUCCESS {
		t.Errorf("h3SetToVertexGraph failed with error: %v", err)
	}

	if graph.Size != 0 {
		t.Errorf("Expected graph size 0, got %d (No edges added to graph)", graph.Size)
	}

	destroyVertexGraph(&graph)
}

func TestH3SetToVertexGraphInternal_singleHex(t *testing.T) {
	t.Parallel()
	var graph VertexGraph
	set := []H3Index{0x890dab6220bffff}
	numHexes := int32(len(set))

	err := h3SetToVertexGraph(set, numHexes, &graph)
	if err != E_SUCCESS {
		t.Errorf("h3SetToVertexGraph failed with error: %v", err)
	}

	if graph.Size != 6 {
		t.Errorf("Expected graph size 6, got %d (All edges of one hex added to graph)", graph.Size)
	}

	destroyVertexGraph(&graph)
}

func TestH3SetToVertexGraphInternal_nonContiguous2(t *testing.T) {
	t.Parallel()
	var graph VertexGraph
	set := []H3Index{0x8928308291bffff, 0x89283082943ffff}
	numHexes := int32(len(set))

	err := h3SetToVertexGraph(set, numHexes, &graph)
	if err != E_SUCCESS {
		t.Errorf("h3SetToVertexGraph failed with error: %v", err)
	}

	if graph.Size != 12 {
		t.Errorf("Expected graph size 12, got %d (All edges of two non-contiguous hexes added to graph)", graph.Size)
	}

	destroyVertexGraph(&graph)
}

func TestH3SetToVertexGraphInternal_contiguous2(t *testing.T) {
	t.Parallel()
	var graph VertexGraph
	set := []H3Index{0x8928308291bffff, 0x89283082957ffff}
	numHexes := int32(len(set))

	err := h3SetToVertexGraph(set, numHexes, &graph)
	if err != E_SUCCESS {
		t.Errorf("h3SetToVertexGraph failed with error: %v", err)
	}

	if graph.Size != 10 {
		t.Errorf("Expected graph size 10, got %d (All edges except 2 shared added to graph)", graph.Size)
	}

	destroyVertexGraph(&graph)
}

func TestH3SetToVertexGraphInternal_contiguous2distorted(t *testing.T) {
	t.Parallel()
	var graph VertexGraph
	set := []H3Index{0x894cc5365afffff, 0x894cc536537ffff}
	numHexes := int32(len(set))

	err := h3SetToVertexGraph(set, numHexes, &graph)
	if err != E_SUCCESS {
		t.Errorf("h3SetToVertexGraph failed with error: %v", err)
	}

	if graph.Size != 12 {
		t.Errorf("Expected graph size 12, got %d (All edges except 2 shared added to graph)", graph.Size)
	}

	destroyVertexGraph(&graph)
}

func TestH3SetToVertexGraphInternal_contiguous3(t *testing.T) {
	t.Parallel()
	var graph VertexGraph
	set := []H3Index{0x8928308288bffff, 0x892830828d7ffff, 0x8928308289bffff}
	numHexes := int32(len(set))

	err := h3SetToVertexGraph(set, numHexes, &graph)
	if err != E_SUCCESS {
		t.Errorf("h3SetToVertexGraph failed with error: %v", err)
	}

	expectedSize := int32(3 * 4) // 3 * 4 = 12
	if graph.Size != expectedSize {
		t.Errorf("Expected graph size %d, got %d (All edges except 6 shared added to graph)", expectedSize, graph.Size)
	}

	destroyVertexGraph(&graph)
}

func TestH3SetToVertexGraphInternal_hole(t *testing.T) {
	t.Parallel()
	var graph VertexGraph
	set := []H3Index{
		0x892830828c7ffff, 0x892830828d7ffff,
		0x8928308289bffff, 0x89283082813ffff,
		0x8928308288fffff, 0x89283082883ffff,
	}
	numHexes := int32(len(set))

	err := h3SetToVertexGraph(set, numHexes, &graph)
	if err != E_SUCCESS {
		t.Errorf("h3SetToVertexGraph failed with error: %v", err)
	}

	expectedSize := int32((6 * 3) + 6) // (6 * 3) + 6 = 24
	if graph.Size != expectedSize {
		t.Errorf("Expected graph size %d, got %d (All outer edges and inner hole edges added to graph)", expectedSize, graph.Size)
	}

	destroyVertexGraph(&graph)
}
