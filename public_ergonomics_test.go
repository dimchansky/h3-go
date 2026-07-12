package h3

import (
	"errors"
	"slices"
	"testing"
)

func TestImmediateHierarchy(t *testing.T) {
	t.Parallel()

	parent, err := sfCellRes9.ImmediateParent()
	if err != nil {
		t.Fatal(err)
	}
	wantParent, err := sfCellRes9.Parent(8)
	if err != nil {
		t.Fatal(err)
	}
	if parent != wantParent {
		t.Fatalf("ImmediateParent = %v, want %v", parent, wantParent)
	}

	children, err := sfCellRes9.ImmediateChildren()
	if err != nil {
		t.Fatal(err)
	}
	wantChildren, err := sfCellRes9.Children(10)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(children, wantChildren) {
		t.Fatalf("ImmediateChildren = %v, want %v", children, wantChildren)
	}

	pentagons, err := Pentagons(9)
	if err != nil {
		t.Fatal(err)
	}
	pentChildren, err := pentagons[0].ImmediateChildren()
	if err != nil {
		t.Fatal(err)
	}
	wantPentChildren, err := pentagons[0].Children(10)
	if err != nil {
		t.Fatal(err)
	}
	if len(pentChildren) != 6 || !slices.Equal(pentChildren, wantPentChildren) {
		t.Fatalf("pentagon ImmediateChildren = %v, want %v", pentChildren, wantPentChildren)
	}

	if _, err := Res0Cells()[0].ImmediateParent(); !errors.Is(err, ErrResolutionDomain) {
		t.Errorf("resolution-0 ImmediateParent: got %v, want ErrResolutionDomain", err)
	}
	if _, err := sfCellRes15.ImmediateChildren(); !errors.Is(err, ErrResolutionDomain) {
		t.Errorf("resolution-15 ImmediateChildren: got %v, want ErrResolutionDomain", err)
	}
	buf := []Cell{sfCellRes9}
	if got, err := sfCellRes15.AppendImmediateChildren(buf); !errors.Is(err, ErrResolutionDomain) || !slices.Equal(got, buf) {
		t.Errorf("resolution-15 AppendImmediateChildren = %v, %v; want unchanged dst and ErrResolutionDomain", got, err)
	}
}

func TestImmediateChildrenAllBaseCellsAndResolutions(t *testing.T) {
	t.Parallel()

	for res := 0; res < MaxResolution; res++ {
		for _, base := range Res0Cells() {
			parent, err := base.CenterChild(res)
			if err != nil {
				t.Fatalf("base %v center child at res %d: %v", base, res, err)
			}
			got, err := parent.ImmediateChildren()
			if err != nil {
				t.Fatalf("res %d, %v: %v", res, parent, err)
			}
			want, err := parent.Children(res + 1)
			if err != nil {
				t.Fatalf("res %d, %v composed: %v", res, parent, err)
			}
			if !slices.Equal(got, want) {
				t.Fatalf("res %d, %v: got %v, want %v", res, parent, got, want)
			}
		}
	}
}

func TestTypedIsValidIndexAndIndexDigits(t *testing.T) {
	t.Parallel()

	edges, err := sfCellRes9.DirectedEdges()
	if err != nil {
		t.Fatal(err)
	}
	vertexes, err := sfCellRes9.Vertexes()
	if err != nil {
		t.Fatal(err)
	}
	if !IsValidIndex(sfCellRes9) || !IsValidIndex(edges[0]) || !IsValidIndex(vertexes[0]) {
		t.Fatal("typed valid indexes must pass IsValidIndex")
	}
	// IsValidIndex is deliberately mode-agnostic: the bit pattern remains a
	// valid edge even when carried in a Cell value.
	if Cell(edges[0]).IsValid() || !IsValidIndex(Cell(edges[0])) {
		t.Fatal("IsValidIndex must validate any mode, not the Go wrapper type")
	}

	origin, err := edges[0].Origin()
	if err != nil {
		t.Fatal(err)
	}
	for res := 1; res <= MaxResolution; res++ {
		got, err := edges[0].IndexDigit(res)
		if err != nil {
			t.Fatal(err)
		}
		want, err := origin.IndexDigit(res)
		if err != nil {
			t.Fatal(err)
		}
		if got != want {
			t.Errorf("edge digit %d = %d, want origin digit %d", res, got, want)
		}
		got, err = vertexes[0].IndexDigit(res)
		if err != nil {
			t.Fatal(err)
		}
		if want := int(h3GetIndexDigit(h3Index(vertexes[0]), int32(res))); got != want {
			t.Errorf("vertex digit %d = %d, want stored digit %d", res, got, want)
		}
	}
	if _, err := edges[0].IndexDigit(0); !errors.Is(err, ErrResolutionDomain) {
		t.Errorf("edge IndexDigit(0): got %v, want ErrResolutionDomain", err)
	}
	if _, err := vertexes[0].IndexDigit(MaxResolution + 1); !errors.Is(err, ErrResolutionDomain) {
		t.Errorf("vertex IndexDigit(16): got %v, want ErrResolutionDomain", err)
	}
}

func TestGridDiskDistancesGrouped(t *testing.T) {
	t.Parallel()

	rings, err := sfCellRes9.GridDiskDistancesGrouped(2)
	if err != nil {
		t.Fatal(err)
	}
	if len(rings) != 3 || len(rings[0]) != 1 || len(rings[1]) != 6 || len(rings[2]) != 12 {
		t.Fatalf("ring lengths = %v, want [1 6 12]", []int{len(rings[0]), len(rings[1]), len(rings[2])})
	}
	for distance, ring := range rings {
		if cap(ring) != len(ring) {
			t.Errorf("ring %d cap = %d, want len %d", distance, cap(ring), len(ring))
		}
		for _, cell := range ring {
			got, err := sfCellRes9.GridDistance(cell)
			if err != nil || got != distance {
				t.Errorf("ring %d cell %v: GridDistance = %d, %v", distance, cell, got, err)
			}
		}
	}
	firstOuter := rings[1][0]
	_ = append(rings[0], sfCellRes15)
	if rings[1][0] != firstOuter {
		t.Fatal("appending to one ring overwrote the next ring")
	}

	pentagons, err := Pentagons(4)
	if err != nil {
		t.Fatal(err)
	}
	pentRings, err := pentagons[0].GridDiskDistancesGrouped(1)
	if err != nil {
		t.Fatal(err)
	}
	if len(pentRings) != 2 || len(pentRings[0]) != 1 || len(pentRings[1]) != 5 {
		t.Fatalf("pentagon ring lengths = %d/%d, want 1/5", len(pentRings[0]), len(pentRings[1]))
	}
	for _, ring := range pentRings {
		for _, cell := range ring {
			if cell == 0 {
				t.Fatal("grouped pentagon disk retained H3_NULL")
			}
		}
	}
	if _, err := sfCellRes9.GridDiskDistancesGrouped(-1); !errors.Is(err, ErrDomain) {
		t.Errorf("negative k: got %v, want ErrDomain", err)
	}
}

func TestErgonomicAPIAllocations(t *testing.T) {
	childrenBuf := make([]Cell, 0, 7)
	if got := testing.AllocsPerRun(200, func() {
		children, err := sfCellRes9.AppendImmediateChildren(childrenBuf[:0])
		if err != nil || len(children) != 7 {
			t.Fatalf("AppendImmediateChildren = %d, %v", len(children), err)
		}
	}); got != 0 {
		t.Errorf("AppendImmediateChildren warm allocs = %g, want 0", got)
	}
	if got := testing.AllocsPerRun(200, func() {
		if !IsValidIndex(sfCellRes9) {
			t.Fatal("valid cell rejected")
		}
	}); got != 0 {
		t.Errorf("typed IsValidIndex allocs = %g, want 0", got)
	}
	if got := testing.AllocsPerRun(100, func() {
		rings, err := sfCellRes9.GridDiskDistancesGrouped(2)
		if err != nil || len(rings) != 3 {
			t.Fatalf("GridDiskDistancesGrouped = %d, %v", len(rings), err)
		}
	}); got > 3 {
		t.Errorf("GridDiskDistancesGrouped allocs = %g, want at most 3", got)
	}
}

func TestNumIcosahedronFaces(t *testing.T) {
	if NumIcosahedronFaces != 20 {
		t.Fatalf("NumIcosahedronFaces = %d, want 20", NumIcosahedronFaces)
	}
}
