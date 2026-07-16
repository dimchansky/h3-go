package h3_test

import (
	"errors"
	"strconv"
	"testing"

	h3 "github.com/dimchansky/h3-go"
)

// TestDocumentedErrorContracts pins the guaranteed error conditions and
// special cases stated in the public GoDoc that are not covered by other
// tests or examples.
func TestDocumentedErrorContracts(t *testing.T) {
	t.Parallel()

	cell, err := h3.ParseCell("8928308280fffff") // resolution 9
	if err != nil {
		t.Fatalf("ParseCell: %v", err)
	}

	// Cell.Parent: same resolution returns the cell itself; a finer
	// resolution is a guaranteed ErrResolutionMismatch.
	if same, err := cell.Parent(9); err != nil || same != cell {
		t.Errorf("Parent(same res) = %v, %v; want the cell itself", same, err)
	}
	if _, err := cell.Parent(10); !errors.Is(err, h3.ErrResolutionMismatch) {
		t.Errorf("Parent(finer res) error = %v, want ErrResolutionMismatch", err)
	}

	// Cell.CenterChild: same resolution returns the cell itself; a coarser
	// resolution is a guaranteed ErrResolutionDomain.
	if same, err := cell.CenterChild(9); err != nil || same != cell {
		t.Errorf("CenterChild(same res) = %v, %v; want the cell itself", same, err)
	}
	if _, err := cell.CenterChild(5); !errors.Is(err, h3.ErrResolutionDomain) {
		t.Errorf("CenterChild(coarser res) error = %v, want ErrResolutionDomain", err)
	}

	// Cell.ChildAtPos: an out-of-range position is a guaranteed ErrDomain.
	if _, err := cell.ChildAtPos(int64(1)<<40, 10); !errors.Is(err, h3.ErrDomain) {
		t.Errorf("ChildAtPos(out of range) error = %v, want ErrDomain", err)
	}

	// UncompactCells family: input finer than the target resolution is a
	// guaranteed ErrResolutionMismatch.
	if _, err := h3.UncompactCells([]h3.Cell{cell}, 5); !errors.Is(err, h3.ErrResolutionMismatch) {
		t.Errorf("UncompactCells(finer input) error = %v, want ErrResolutionMismatch", err)
	}
	if _, err := h3.UncompactCellsSize([]h3.Cell{cell}, 5); !errors.Is(err, h3.ErrResolutionMismatch) {
		t.Errorf("UncompactCellsSize(finer input) error = %v, want ErrResolutionMismatch", err)
	}

	// ParseCell: >64-bit values wrap strconv.ErrRange, not an Err* sentinel.
	if _, err := h3.ParseCell("1ffffffffffffffff"); !errors.Is(err, strconv.ErrRange) {
		t.Errorf("ParseCell(>64 bit) error = %v, want wrapped strconv.ErrRange", err)
	}

	// Cell.IsNeighbor: a cell is never its own neighbor — (false, nil).
	if ok, err := cell.IsNeighbor(cell); ok || err != nil {
		t.Errorf("IsNeighbor(self) = %v, %v; want false, nil", ok, err)
	}
}

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
