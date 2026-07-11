// Tests ported from testCellToBBoxExhaustive.c
package h3

import (
	"testing"
)

// cellBBoxAssertions tests that a cell's bounding box contains all its vertices.
func cellBBoxAssertions(t *testing.T, h3 h3Index) {
	bbox, err := cellToBBox(h3, false)
	if err != eSuccess {
		t.Fatalf("cellToBBox failed for cell %x: %v", h3, err)
	}

	var verts CellBoundary
	err = cellToBoundary(h3, &verts)
	if err != eSuccess {
		t.Fatalf("cellToBoundary failed for cell %x: %v", h3, err)
	}

	for j := int32(0); j < verts.NumVerts; j++ {
		if !bboxContains(&bbox, &verts.Verts[j]) {
			t.Errorf("bbox does not contain cell vertex %d for cell %x", j, h3)
			t.Logf("Cell: %x", h3)
			t.Logf("bbox: North=%f, South=%f, East=%f, West=%f",
				bbox.North.Deg(), bbox.South.Deg(), bbox.East.Deg(), bbox.West.Deg())
			t.Logf("Vertex: Lat=%f, Lng=%f",
				verts.Verts[j].Lat.Deg(), verts.Verts[j].Lng.Deg())
		}
	}
}

// childBBoxAssertions tests that a cell's child-covering bounding box contains all vertices of its children.
func childBBoxAssertions(t *testing.T, h3 h3Index) {
	parentRes := getResolution(h3)

	bbox, err := cellToBBox(h3, true)
	if err != eSuccess {
		t.Fatalf("cellToBBox with coverChildren=true failed for cell %x: %v", h3, err)
	}

	for resolutionOffset := int32(0); resolutionOffset < 5; resolutionOffset++ {
		// Test whether all verts of all children are inside the bbox
		childRes := parentRes + resolutionOffset

		numChildren, err := cellToChildrenSize(h3, childRes)
		if err != eSuccess {
			t.Fatalf("cellToChildrenSize failed for cell %x at resolution %d: %v", h3, childRes, err)
		}

		children := make([]h3Index, numChildren)
		err = cellToChildren(h3, childRes, children)
		if err != eSuccess {
			t.Fatalf("cellToChildren failed for cell %x at resolution %d: %v", h3, childRes, err)
		}

		for i := range children {
			var childVerts CellBoundary
			err = cellToBoundary(children[i], &childVerts)
			if err != eSuccess {
				t.Fatalf("cellToBoundary failed for child cell %x: %v", children[i], err)
			}

			for j := int32(0); j < childVerts.NumVerts; j++ {
				if !bboxContains(&bbox, &childVerts.Verts[j]) {
					t.Errorf("bbox does not contain child vertex %d for parent %x, child %x", j, h3, children[i])
					t.Logf("Parent: %x", h3)
					t.Logf("bbox: North=%f, South=%f, East=%f, West=%f",
						bbox.North.Deg(), bbox.South.Deg(), bbox.East.Deg(), bbox.West.Deg())
					t.Logf("Child: %x", children[i])
					t.Logf("Vertex: Lat=%f, Lng=%f",
						childVerts.Verts[j].Lat.Deg(), childVerts.Verts[j].Lng.Deg())
				}
			}
		}
	}
}

func TestCellBBox_correctness(t *testing.T) {
	t.Parallel()

	// Test resolution 0
	_iterateAllIndexesAtRes(0, func(h3 h3Index) {
		cellBBoxAssertions(t, h3)
	})

	// Test resolution 1
	_iterateAllIndexesAtRes(1, func(h3 h3Index) {
		cellBBoxAssertions(t, h3)
	})

	// Test resolution 2
	_iterateAllIndexesAtRes(2, func(h3 h3Index) {
		cellBBoxAssertions(t, h3)
	})
}

func TestChildBBox_correctness(t *testing.T) {
	t.Parallel()

	// Test resolution 0
	_iterateAllIndexesAtRes(0, func(h3 h3Index) {
		childBBoxAssertions(t, h3)
	})

	// Test resolution 1
	_iterateAllIndexesAtRes(1, func(h3 h3Index) {
		childBBoxAssertions(t, h3)
	})

	// Test resolution 2
	_iterateAllIndexesAtRes(2, func(h3 h3Index) {
		childBBoxAssertions(t, h3)
	})
}
