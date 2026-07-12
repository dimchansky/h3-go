package uberbench

import (
	"fmt"
	"slices"
	"testing"

	pure "github.com/dimchansky/h3-go"
	uber "github.com/uber/h3-go/v4"
)

// This file keeps the before/after example in
// docs/migration-from-uber-h3-go.md executable: coverageUber is the
// "before" (binding) version, coveragePure the "after" (this library)
// version, and the test asserts they produce identical output for the
// same input. If an API change breaks either version, this fails before
// the guide goes stale.

// coverageUber is the guide's "before" snippet.
func coverageUber(cellStr string, fence uber.GeoPolygon) ([]string, error) {
	c := uber.CellFromString(cellStr)
	if !c.IsValid() {
		return nil, fmt.Errorf("bad cell %q", cellStr)
	}
	disk, err := c.GridDisk(1)
	if err != nil {
		return nil, err
	}
	area, err := uber.PolygonToCells(fence, c.Resolution())
	if err != nil {
		return nil, err
	}
	compact, err := uber.CompactCells(append(area, disk...))
	if err != nil {
		return nil, err
	}
	out := make([]string, len(compact))
	for i, cc := range compact {
		out[i] = uber.CellToString(cc)
	}
	return out, nil
}

// coveragePure is the guide's "after" snippet.
func coveragePure(cellStr string, fence pure.GeoPolygon) ([]string, error) {
	c, err := pure.ParseCell(cellStr) // parse + validate in one step
	if err != nil {
		return nil, fmt.Errorf("bad cell %q: %w", cellStr, err)
	}
	disk, err := c.GridDisk(1)
	if err != nil {
		return nil, err
	}
	area, err := pure.PolygonToCells(fence, c.Resolution())
	if err != nil {
		return nil, err
	}
	compact, err := pure.CompactCells(append(area, disk...))
	if err != nil {
		return nil, err
	}
	out := make([]string, len(compact))
	for i, cc := range compact {
		out[i] = cc.String()
	}
	return out, nil
}

func TestMigrationExampleEquivalence(t *testing.T) {
	// A res-9 cell in Manhattan: far from the SF fence, so the k=1 disk
	// and the polyfill are disjoint (CompactCells requires unique input).
	cell, err := pure.LatLngToCell(pure.LatLngDegs(40.7128, -74.0060), 9)
	if err != nil {
		t.Fatal(err)
	}
	cellStr := cell.String()

	got, err := coveragePure(cellStr, sfPolygonPure)
	if err != nil {
		t.Fatalf("pure coverage: %v", err)
	}
	want, err := coverageUber(cellStr, sfPolygonUber)
	if err != nil {
		t.Fatalf("uber coverage: %v", err)
	}
	slices.Sort(got)
	slices.Sort(want)
	if !slices.Equal(got, want) {
		t.Fatalf("migration example outputs differ: pure %d strings, uber %d strings", len(got), len(want))
	}
	if len(got) == 0 {
		t.Fatal("migration example produced no cells")
	}

	// Both versions must also agree on rejecting bad input.
	if _, err := coveragePure("not-a-cell", sfPolygonPure); err == nil {
		t.Fatal("pure coverage accepted an invalid cell string")
	}
	if _, err := coverageUber("not-a-cell", sfPolygonUber); err == nil {
		t.Fatal("uber coverage accepted an invalid cell string")
	}
}
