package h3_test

import (
	"testing"

	h3 "github.com/dimchansky/h3-go"
)

// TestChildrenNotGeometricallyContained pins the documented
// logical-vs-geometric hierarchy contract (Cell.Parent, Cell.Children,
// package docs): the parent/child relationship is index-hierarchical, and a
// child's boundary is not required to lie within the parent's boundary.
// Filling the parent's own boundary polygon in ContainmentFull mode at the
// children's resolution captures only a subset of the children.
func TestChildrenNotGeometricallyContained(t *testing.T) {
	t.Parallel()

	parent, err := h3.ParseCell("85283473fffffff")
	if err != nil {
		t.Fatalf("ParseCell: %v", err)
	}
	children, err := parent.Children(6)
	if err != nil {
		t.Fatalf("Children: %v", err)
	}
	if len(children) != 7 {
		t.Fatalf("Children returned %d cells, want 7", len(children))
	}

	boundary, err := parent.Boundary()
	if err != nil {
		t.Fatalf("Boundary: %v", err)
	}
	fully, err := h3.PolygonToCellsExperimental(
		h3.GeoPolygon{GeoLoop: h3.GeoLoop(boundary.Verts())}, 6, h3.ContainmentFull)
	if err != nil {
		t.Fatalf("PolygonToCellsExperimental: %v", err)
	}
	if len(fully) >= len(children) {
		t.Fatalf("all %d children fully contained in the parent boundary; "+
			"expected fewer (logical, not geometric, containment)", len(fully))
	}
}
