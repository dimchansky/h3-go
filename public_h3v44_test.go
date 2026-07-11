package h3

import (
	"errors"
	"testing"
)

func TestIndexDigit(t *testing.T) {
	t.Parallel()

	// Reconstruct the SF res-9 cell from its digits.
	digits := make([]int, 9)
	for r := 1; r <= 9; r++ {
		d, err := sfCellRes9.IndexDigit(r)
		if err != nil {
			t.Fatal(err)
		}
		if d < 0 || d > 6 {
			t.Fatalf("digit %d out of range: %d", r, d)
		}
		digits[r-1] = d
	}
	// Digits beyond the resolution are 7 for valid cells.
	d10, err := sfCellRes9.IndexDigit(10)
	if err != nil || d10 != 7 {
		t.Errorf("IndexDigit(10) = %d (%v), want 7", d10, err)
	}
	if _, err := sfCellRes9.IndexDigit(0); !errors.Is(err, ErrResolutionDomain) {
		t.Errorf("IndexDigit(0): got %v, want ErrResolutionDomain", err)
	}

	back, err := ConstructCell(9, sfCellRes9.BaseCellNumber(), digits)
	if err != nil {
		t.Fatal(err)
	}
	if back != sfCellRes9 {
		t.Fatalf("ConstructCell round trip = %v, want %v", back, sfCellRes9)
	}
}

func TestConstructCell(t *testing.T) {
	t.Parallel()

	// res 0: digits may be nil; must equal the base cell.
	c, err := ConstructCell(0, 20, nil)
	if err != nil {
		t.Fatal(err)
	}
	if c != Res0Cells()[20] {
		t.Fatalf("ConstructCell(0, 20) = %v", c)
	}

	if _, err := ConstructCell(0, NumBaseCells, nil); !errors.Is(err, ErrBaseCellDomain) {
		t.Errorf("bad base cell: got %v, want ErrBaseCellDomain", err)
	}
	if _, err := ConstructCell(-1, 0, nil); !errors.Is(err, ErrResolutionDomain) {
		t.Errorf("bad res: got %v, want ErrResolutionDomain", err)
	}
	if _, err := ConstructCell(1, 0, []int{7}); !errors.Is(err, ErrDigitDomain) {
		t.Errorf("digit 7: got %v, want ErrDigitDomain", err)
	}
	if _, err := ConstructCell(1, 0, nil); !errors.Is(err, ErrDigitDomain) {
		t.Errorf("missing digits: got %v, want ErrDigitDomain", err)
	}
	// Base cell 4 is a pentagon: k-axis digit (1) is a deleted subsequence.
	if _, err := ConstructCell(1, 4, []int{1}); !errors.Is(err, ErrDeletedDigit) {
		t.Errorf("pentagon k-digit: got %v, want ErrDeletedDigit", err)
	}
	// Center digit on a pentagon stays valid.
	pc, err := ConstructCell(1, 4, []int{0})
	if err != nil || !pc.IsPentagon() {
		t.Errorf("pentagon center child: %v (%v)", pc, err)
	}

	// Every constructed cell must be valid, and reconstruct from digits.
	for _, res := range []int{1, 5, 15} {
		orig, err := LatLngToCell(LatLngDegs(sfLatDegs, sfLngDegs), res)
		if err != nil {
			t.Fatal(err)
		}
		digits := make([]int, res)
		for r := 1; r <= res; r++ {
			digits[r-1], _ = orig.IndexDigit(r)
		}
		back, err := ConstructCell(res, orig.BaseCellNumber(), digits)
		if err != nil || back != orig {
			t.Errorf("res %d: reconstruct = %v (%v), want %v", res, back, err, orig)
		}
	}
}

func TestIsValidIndex(t *testing.T) {
	t.Parallel()

	if !IsValidIndex(uint64(sfCellRes9)) {
		t.Error("valid cell must be a valid index")
	}
	edges, err := sfCellRes9.DirectedEdges()
	if err != nil {
		t.Fatal(err)
	}
	if !IsValidIndex(uint64(edges[0])) {
		t.Error("valid edge must be a valid index")
	}
	verts, err := sfCellRes9.Vertexes()
	if err != nil {
		t.Fatal(err)
	}
	if !IsValidIndex(uint64(verts[0])) {
		t.Error("valid vertex must be a valid index")
	}
	if IsValidIndex(0) || IsValidIndex(^uint64(0)) {
		t.Error("zero/all-ones must not be valid indexes")
	}
	// Corrupted index (high bit set), from upstream testH3Index.c isValidIndex.
	if IsValidIndex(uint64(sfCellRes9) | 1<<63) {
		t.Error("high-bit-corrupted index must not be valid")
	}
}
