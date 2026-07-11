// Tests ported from testVertexGraphInternal.c
package h3

import (
	"testing"
)

// Test fixtures - mirrors the C test fixtures.
func getTestVertices() (center, vertex1, vertex2, vertex3, vertex4, vertex5, vertex6 LatLng) {
	setGeoDegs(&center, 37.77362016769341, -122.41673772517154)
	setGeoDegs(&vertex1, 87.372002166, 166.160981117)
	setGeoDegs(&vertex2, 87.370101364, 166.160184306)
	setGeoDegs(&vertex3, 87.369088356, 166.196239997)
	setGeoDegs(&vertex4, 87.369975080, 166.233115768)
	setGeoDegs(&vertex5, 0, 0)
	setGeoDegs(&vertex6, -10, -10)
	return
}

func TestMakeVertexGraph(t *testing.T) {
	t.Parallel()
	var graph vertexGraph
	initVertexGraph(&graph, 10, 9)

	if graph.NumBuckets != 10 {
		t.Errorf("numBuckets not set correctly: got %d, want 10", graph.NumBuckets)
	}
	if graph.Size != 0 {
		t.Errorf("size not initialized to 0: got %d", graph.Size)
	}

	destroyVertexGraph(&graph)
}

func TestVertexHash(t *testing.T) {
	t.Parallel()
	center, _, _, _, _, _, _ := getTestVertices()
	numBuckets := int32(1000)

	for res := 0; res < 11; res++ {
		var centerIndex h3Index
		err := latLngToCell(&center, int32(res), &centerIndex)
		if err != eSuccess {
			t.Fatalf("latLngToCell failed for res %d: %v", res, err)
		}

		var boundary CellBoundary
		err = cellToBoundary(centerIndex, &boundary)
		if err != eSuccess {
			t.Fatalf("cellToBoundary failed: %v", err)
		}

		numVerts := int(boundary.numVerts)
		for i := 0; i < numVerts; i++ {
			hash1 := _hashVertex(&boundary.verts[i], int32(res), numBuckets)
			hash2 := _hashVertex(&boundary.verts[(i+1)%numVerts], int32(res), numBuckets)

			if hash1 == hash2 {
				t.Errorf("Hashes must not be equal at res %d, vertex %d: both are %d", res, i, hash1)
			}
		}
	}
}

func TestVertexHashNegative(t *testing.T) {
	t.Parallel()
	_, _, _, _, _, vertex5, vertex6 := getTestVertices()
	numBuckets := int32(10)

	hash5 := _hashVertex(&vertex5, 5, numBuckets)
	if hash5 >= uint32(numBuckets) {
		t.Errorf("zero vertex hash out of bounds: %d >= %d", hash5, numBuckets)
	}

	hash6 := _hashVertex(&vertex6, 5, numBuckets)
	if hash6 >= uint32(numBuckets) {
		t.Errorf("negative coordinates vertex hash out of bounds: %d >= %d", hash6, numBuckets)
	}
}

func TestAddVertexNode(t *testing.T) {
	t.Parallel()
	_, vertex1, vertex2, vertex3, vertex4, _, _ := getTestVertices()
	var graph vertexGraph
	initVertexGraph(&graph, 10, 9)

	// Basic add
	addedNode := addVertexNode(&graph, &vertex1, &vertex2)
	node := findNodeForEdge(&graph, &vertex1, &vertex2)
	if node == nil {
		t.Error("Node not found after add")
	}
	if node != addedNode {
		t.Error("Wrong node found")
	}
	if graph.Size != 1 {
		t.Errorf("Graph size not incremented: got %d, want 1", graph.Size)
	}

	// Collision add
	addedNode = addVertexNode(&graph, &vertex1, &vertex3)
	node = findNodeForEdge(&graph, &vertex1, &vertex3)
	if node == nil {
		t.Error("Node not found after hash collision")
	}
	if node != addedNode {
		t.Error("Wrong node found after collision")
	}
	if graph.Size != 2 {
		t.Errorf("Graph size not incremented: got %d, want 2", graph.Size)
	}

	// Collision add #2
	addedNode = addVertexNode(&graph, &vertex1, &vertex4)
	node = findNodeForEdge(&graph, &vertex1, &vertex4)
	if node == nil {
		t.Error("Node not found after 2nd hash collision")
	}
	if node != addedNode {
		t.Error("Wrong node found after 2nd collision")
	}
	if graph.Size != 3 {
		t.Errorf("Graph size not incremented: got %d, want 3", graph.Size)
	}

	// Exact match no-op
	oldNode := findNodeForEdge(&graph, &vertex1, &vertex2)
	addedNode = addVertexNode(&graph, &vertex1, &vertex2)
	if findNodeForEdge(&graph, &vertex1, &vertex2) != oldNode {
		t.Error("Exact match changed existing node")
	}
	if addedNode != oldNode {
		t.Error("Old node not returned for duplicate")
	}
	if graph.Size != 3 {
		t.Errorf("Graph size changed on duplicate: got %d, want 3", graph.Size)
	}

	destroyVertexGraph(&graph)
}

func TestAddVertexNodeDupe(t *testing.T) {
	t.Parallel()
	_, vertex1, vertex2, _, _, _, _ := getTestVertices()
	var graph vertexGraph
	initVertexGraph(&graph, 10, 9)

	// Basic add
	addedNode := addVertexNode(&graph, &vertex1, &vertex2)
	node := findNodeForEdge(&graph, &vertex1, &vertex2)
	if node == nil {
		t.Error("Node not found")
	}
	if node != addedNode {
		t.Error("Wrong node found")
	}
	if graph.Size != 1 {
		t.Errorf("Graph size not incremented: got %d, want 1", graph.Size)
	}

	// Dupe add
	addedNode = addVertexNode(&graph, &vertex1, &vertex2)
	if node != addedNode {
		t.Error("addVertexNode did not return the original node")
	}
	if graph.Size != 1 {
		t.Errorf("Graph size incremented on duplicate: got %d, want 1", graph.Size)
	}

	destroyVertexGraph(&graph)
}

func TestFindNodeForEdge(t *testing.T) {
	t.Parallel()
	// Basic lookup tested in TestAddVertexNode, only test failures here
	_, vertex1, vertex2, vertex3, vertex4, _, _ := getTestVertices()
	var graph vertexGraph
	initVertexGraph(&graph, 10, 9)

	// Empty graph
	node := findNodeForEdge(&graph, &vertex1, &vertex2)
	if node != nil {
		t.Error("Node lookup should fail for empty graph")
	}

	addVertexNode(&graph, &vertex1, &vertex2)

	// Different hash
	node = findNodeForEdge(&graph, &vertex3, &vertex2)
	if node != nil {
		t.Error("Node lookup should fail for different hash")
	}

	// Hash collision
	node = findNodeForEdge(&graph, &vertex1, &vertex3)
	if node != nil {
		t.Error("Node lookup should fail for hash collision")
	}

	addVertexNode(&graph, &vertex1, &vertex4)

	// Hash collision, list iteration
	node = findNodeForEdge(&graph, &vertex1, &vertex3)
	if node != nil {
		t.Error("Node lookup should fail for collision w/iteration")
	}

	destroyVertexGraph(&graph)
}

func TestFindNodeForVertex(t *testing.T) {
	t.Parallel()
	_, vertex1, vertex2, vertex3, _, _, _ := getTestVertices()
	var graph vertexGraph
	initVertexGraph(&graph, 10, 9)

	// Empty graph
	node := findNodeForVertex(&graph, &vertex1)
	if node != nil {
		t.Error("Node lookup should fail for empty graph")
	}

	addVertexNode(&graph, &vertex1, &vertex2)

	node = findNodeForVertex(&graph, &vertex1)
	if node == nil {
		t.Error("Node lookup should succeed for correct node")
	}

	node = findNodeForVertex(&graph, &vertex3)
	if node != nil {
		t.Error("Node lookup should fail for different node")
	}

	destroyVertexGraph(&graph)
}

func TestRemoveVertexNode(t *testing.T) {
	t.Parallel()
	_, vertex1, vertex2, vertex3, vertex4, _, _ := getTestVertices()
	var graph vertexGraph
	initVertexGraph(&graph, 10, 9)

	// Straight removal
	node := addVertexNode(&graph, &vertex1, &vertex2)
	success := removeVertexNode(&graph, node) == 0

	if !success {
		t.Error("Removal failed")
	}
	if findNodeForVertex(&graph, &vertex1) != nil {
		t.Error("Node still found after removal")
	}
	if graph.Size != 0 {
		t.Errorf("Graph size not decremented: got %d, want 0", graph.Size)
	}

	// Remove end of list
	addVertexNode(&graph, &vertex1, &vertex2)
	node = addVertexNode(&graph, &vertex1, &vertex3)
	success = removeVertexNode(&graph, node) == 0

	if !success {
		t.Error("Removal of end node failed")
	}
	if findNodeForEdge(&graph, &vertex1, &vertex3) != nil {
		t.Error("Removed node still found")
	}
	baseNode := findNodeForEdge(&graph, &vertex1, &vertex2)
	if baseNode == nil || baseNode.Next != nil {
		t.Error("Base bucket node incorrectly pointing to removed node")
	}
	if graph.Size != 1 {
		t.Errorf("Graph size not decremented: got %d, want 1", graph.Size)
	}

	// Clean up for next test
	node = findNodeForVertex(&graph, &vertex1)
	if removeVertexNode(&graph, node) != 0 {
		t.Error("Cleanup removal failed")
	}

	// Remove beginning of list
	node = addVertexNode(&graph, &vertex1, &vertex2)
	addVertexNode(&graph, &vertex1, &vertex3)
	success = removeVertexNode(&graph, node) == 0

	if !success {
		t.Error("Removal of beginning node failed")
	}
	if findNodeForEdge(&graph, &vertex1, &vertex2) != nil {
		t.Error("Removed beginning node still found")
	}
	if findNodeForEdge(&graph, &vertex1, &vertex3) == nil {
		t.Error("End of list not found after beginning removal")
	}
	endNode := findNodeForEdge(&graph, &vertex1, &vertex3)
	if endNode == nil || endNode.Next != nil {
		t.Error("New beginning node incorrectly has next pointer")
	}
	if graph.Size != 1 {
		t.Errorf("Graph size not decremented: got %d, want 1", graph.Size)
	}

	// Clean up for next test
	node = findNodeForVertex(&graph, &vertex1)
	if removeVertexNode(&graph, node) != 0 {
		t.Error("Cleanup removal failed")
	}

	// Remove middle of list
	addVertexNode(&graph, &vertex1, &vertex2)
	node = addVertexNode(&graph, &vertex1, &vertex3)
	addVertexNode(&graph, &vertex1, &vertex4)
	success = removeVertexNode(&graph, node) == 0

	if !success {
		t.Error("Removal of middle node failed")
	}
	if findNodeForEdge(&graph, &vertex1, &vertex3) != nil {
		t.Error("Removed middle node still found")
	}
	if findNodeForEdge(&graph, &vertex1, &vertex4) == nil {
		t.Error("End of list not found after middle removal")
	}
	if graph.Size != 2 {
		t.Errorf("Graph size not decremented: got %d, want 2", graph.Size)
	}

	// Remove non-existent node
	// Create a node that's not in the graph
	fakeNode := &vertexNode{}
	success = removeVertexNode(&graph, fakeNode) == 0

	if success {
		t.Error("Removal of non-existent node should fail")
	}
	if graph.Size != 2 {
		t.Errorf("Graph size changed after failed removal: got %d, want 2", graph.Size)
	}

	destroyVertexGraph(&graph)
}

func TestFirstVertexNode(t *testing.T) {
	t.Parallel()
	_, vertex1, vertex2, _, _, _, _ := getTestVertices()
	var graph vertexGraph
	initVertexGraph(&graph, 10, 9)

	node := firstVertexNode(&graph)
	if node != nil {
		t.Error("No node should be found for empty graph")
	}

	addedNode := addVertexNode(&graph, &vertex1, &vertex2)

	node = firstVertexNode(&graph)
	if node != addedNode {
		t.Error("First node not found correctly")
	}

	destroyVertexGraph(&graph)
}

func TestDestroyEmptyVertexGraph(t *testing.T) {
	t.Parallel()
	var graph vertexGraph
	initVertexGraph(&graph, 10, 9)
	destroyVertexGraph(&graph)
	// Test passes if no panic/crash occurs
}

func TestSingleBucketVertexGraph(t *testing.T) {
	t.Parallel()
	_, vertex1, vertex2, vertex3, vertex4, _, _ := getTestVertices()
	var graph vertexGraph
	initVertexGraph(&graph, 1, 9)

	if graph.NumBuckets != 1 {
		t.Errorf("Wrong number of buckets: got %d, want 1", graph.NumBuckets)
	}

	node := firstVertexNode(&graph)
	if node != nil {
		t.Error("No node should be found for empty graph")
	}

	node = addVertexNode(&graph, &vertex1, &vertex2)
	if node == nil {
		t.Error("Node not added")
	}
	if firstVertexNode(&graph) != node {
		t.Error("First node is not the added node")
	}

	addVertexNode(&graph, &vertex2, &vertex3)
	addVertexNode(&graph, &vertex3, &vertex4)
	if firstVertexNode(&graph) != node {
		t.Error("First node changed after adding more nodes")
	}
	if graph.Size != 3 {
		t.Errorf("Graph size not updated correctly: got %d, want 3", graph.Size)
	}

	destroyVertexGraph(&graph)
}
