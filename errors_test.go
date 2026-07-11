package h3

import (
	"errors"
	"strings"
	"testing"
)

func TestToErrMapping(t *testing.T) {
	t.Parallel()

	cases := []struct {
		code H3Error
		want error
	}{
		{E_SUCCESS, nil},
		{E_FAILED, ErrFailed},
		{E_DOMAIN, ErrDomain},
		{E_LATLNG_DOMAIN, ErrLatLngDomain},
		{E_RES_DOMAIN, ErrResolutionDomain},
		{E_CELL_INVALID, ErrCellInvalid},
		{E_DIR_EDGE_INVALID, ErrDirectedEdgeInvalid},
		{E_UNDIR_EDGE_INVALID, ErrUndirectedEdgeInvalid},
		{E_VERTEX_INVALID, ErrVertexInvalid},
		{E_PENTAGON, ErrPentagon},
		{E_DUPLICATE_INPUT, ErrDuplicateInput},
		{E_NOT_NEIGHBORS, ErrNotNeighbors},
		{E_RES_MISMATCH, ErrResolutionMismatch},
		{E_MEMORY_ALLOC, ErrMemoryAlloc},
		{E_MEMORY_BOUNDS, ErrMemoryBounds},
		{E_OPTION_INVALID, ErrOptionInvalid},
	}
	for _, c := range cases {
		got := toErr(c.code)
		if c.want == nil {
			if got != nil {
				t.Errorf("toErr(%d) = %v, want nil", c.code, got)
			}
			continue
		}
		if !errors.Is(got, c.want) {
			t.Errorf("toErr(%d) = %v, want %v", c.code, got, c.want)
		}
	}

	// Unknown codes map to ErrFailed.
	if got := toErr(H3Error(42)); !errors.Is(got, ErrFailed) {
		t.Errorf("toErr(42) = %v, want ErrFailed", got)
	}
	if got := toErr(H3Error(16)); !errors.Is(got, ErrFailed) {
		t.Errorf("toErr(16) = %v, want ErrFailed", got)
	}
}

func TestSentinelMessages(t *testing.T) {
	t.Parallel()

	// Every sentinel carries the C describeH3Error text with an "h3: " prefix.
	for code := H3Error(1); code <= 15; code++ {
		err := toErr(code)
		if err == nil {
			t.Fatalf("toErr(%d) = nil", code)
		}
		want := "h3: " + describeH3Error(code)
		if err.Error() != want {
			t.Errorf("toErr(%d).Error() = %q, want %q", code, err.Error(), want)
		}
	}
}

func TestSentinelsDistinct(t *testing.T) {
	t.Parallel()

	seen := map[error]H3Error{}
	for code := H3Error(1); code <= 15; code++ {
		err := toErr(code)
		if prev, dup := seen[err]; dup {
			t.Errorf("codes %d and %d map to the same sentinel %v", prev, code, err)
		}
		seen[err] = code
	}
}

func TestToErrDoesNotAllocate(t *testing.T) {
	allocs := testing.AllocsPerRun(100, func() {
		_ = toErr(E_PENTAGON)
		_ = toErr(E_SUCCESS)
	})
	if allocs != 0 {
		t.Errorf("toErr allocates %v times per run, want 0", allocs)
	}
}

// TestAliasTypeIdentity pins the architectural keystone: H3Index is an alias
// of Cell, so slices of the two are the very same type and pass through the
// ported layer with zero conversion (docs/public-api-architecture.md DR-003).
func TestAliasTypeIdentity(t *testing.T) {
	t.Parallel()

	cells := []Cell{0x8f2830828052d25}
	// Passing []Cell where []H3Index is expected compiles only because the
	// two are the same type; the write must be visible through both.
	takesIndexes := func(idx []H3Index) { idx[0] = 0 }
	takesIndexes(cells)
	if cells[0] != 0 {
		t.Fatal("[]H3Index and []Cell must share storage")
	}

	// Scalar conversions between the index types are representation-preserving.
	raw := uint64(0x8f2830828052d25)
	if uint64(Cell(raw)) != raw || uint64(DirectedEdge(raw)) != raw || uint64(Vertex(raw)) != raw {
		t.Fatal("index type conversions must preserve the raw representation")
	}
}

func TestCuratedConstants(t *testing.T) {
	t.Parallel()

	if MaxResolution != 15 {
		t.Errorf("MaxResolution = %d, want 15", MaxResolution)
	}
	if NumBaseCells != 122 || NumRes0Cells != 122 {
		t.Errorf("NumBaseCells/NumRes0Cells = %d/%d, want 122/122", NumBaseCells, NumRes0Cells)
	}
	if NumPentagons != 12 {
		t.Errorf("NumPentagons = %d, want 12", NumPentagons)
	}
	if MaxCellBoundaryVerts != 10 {
		t.Errorf("MaxCellBoundaryVerts = %d, want 10", MaxCellBoundaryVerts)
	}
	v := [3]int{VersionMajor, VersionMinor, VersionPatch}
	if v != [3]int{4, 3, 0} {
		t.Errorf("Version = %v, want [4 3 0]", v)
	}

	// The error message text must never look like a Go-style wrapped chain.
	if strings.Contains(ErrPentagon.Error(), ":") && !strings.HasPrefix(ErrPentagon.Error(), "h3: ") {
		t.Errorf("unexpected sentinel format: %q", ErrPentagon.Error())
	}
}
