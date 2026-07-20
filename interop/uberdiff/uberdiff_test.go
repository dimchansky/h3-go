// Package uberdiff cross-checks this pure-Go H3 implementation against the
// official cgo binding github.com/uber/h3-go (which wraps the H3 C library).
//
// This module is intentionally separate from the root module so the library
// keeps zero dependencies; running these tests requires cgo and network
// access (make test-uberdiff). See README.md in this directory for what
// this suite does and does not prove, and how it relates to the direct C
// parity suite in the root module.
package uberdiff

import (
	"math"
	"math/rand/v2"
	"slices"
	"testing"

	pure "github.com/dimchansky/h3-go"
	uber "github.com/uber/h3-go/v4"
)

const numRandom = 500

// randomLatLngs yields deterministic pseudo-random coordinates covering the
// globe, including polar and antimeridian neighborhoods.
func randomLatLngs(n int) [][2]float64 {
	rng := rand.New(rand.NewPCG(0xd1ce, 0x5eed))
	out := make([][2]float64, 0, n+4)
	out = append(out, [2]float64{89.9, 0}, [2]float64{-89.9, 45}, [2]float64{0, 179.9}, [2]float64{0, -179.9})
	for range n {
		out = append(out, [2]float64{rng.Float64()*180 - 90, rng.Float64()*360 - 180})
	}
	return out
}

func TestLatLngToCellParity(t *testing.T) {
	for _, ll := range randomLatLngs(numRandom) {
		for _, res := range []int{0, 4, 9, 15} {
			got, err := pure.LatLngToCell(pure.LatLngDegs(ll[0], ll[1]), res)
			if err != nil {
				t.Fatalf("pure LatLngToCell(%v, %d): %v", ll, res, err)
			}
			want, err := uber.LatLngToCell(uber.NewLatLng(ll[0], ll[1]), res)
			if err != nil {
				t.Fatalf("uber LatLngToCell(%v, %d): %v", ll, res, err)
			}
			if uint64(got) != uint64(want) {
				t.Fatalf("LatLngToCell(%v, %d): pure %v != uber %v", ll, res, got, want)
			}
		}
	}
}

func TestCellRoundTripParity(t *testing.T) {
	for _, ll := range randomLatLngs(numRandom) {
		c, err := pure.LatLngToCell(pure.LatLngDegs(ll[0], ll[1]), 9)
		if err != nil {
			t.Fatal(err)
		}
		uc := uber.Cell(c)

		// Centroid.
		gotLL, err := c.LatLng()
		if err != nil {
			t.Fatal(err)
		}
		wantLL, err := uc.LatLng()
		if err != nil {
			t.Fatal(err)
		}
		if math.Abs(gotLL.Lat.Deg()-wantLL.Lat) > 1e-10 || math.Abs(gotLL.Lng.Deg()-wantLL.Lng) > 1e-10 {
			t.Fatalf("LatLng(%v): pure %v,%v != uber %v,%v", c,
				gotLL.Lat.Deg(), gotLL.Lng.Deg(), wantLL.Lat, wantLL.Lng)
		}

		// Boundary.
		gotB, err := c.Boundary()
		if err != nil {
			t.Fatal(err)
		}
		wantB, err := uc.Boundary()
		if err != nil {
			t.Fatal(err)
		}
		if gotB.Len() != len(wantB) {
			t.Fatalf("Boundary(%v): pure %d verts != uber %d", c, gotB.Len(), len(wantB))
		}
		for i, v := range gotB.Verts() {
			if math.Abs(v.Lat.Deg()-wantB[i].Lat) > 1e-10 || math.Abs(v.Lng.Deg()-wantB[i].Lng) > 1e-10 {
				t.Fatalf("Boundary(%v)[%d]: pure %v != uber %v", c, i, v, wantB[i])
			}
		}

		// Inspection.
		if c.IsPentagon() != uc.IsPentagon() || c.Resolution() != uc.Resolution() ||
			c.BaseCellNumber() != uc.BaseCellNumber() || c.IsValid() != uc.IsValid() {
			t.Fatalf("inspection mismatch for %v", c)
		}
	}
}

func TestHierarchyParity(t *testing.T) {
	for _, ll := range randomLatLngs(200) {
		c, _ := pure.LatLngToCell(pure.LatLngDegs(ll[0], ll[1]), 9)
		uc := uber.Cell(c)

		gotP, err := c.Parent(5)
		if err != nil {
			t.Fatal(err)
		}
		wantP, err := uc.Parent(5)
		if err != nil {
			t.Fatal(err)
		}
		if uint64(gotP) != uint64(wantP) {
			t.Fatalf("Parent(%v): %v != %v", c, gotP, wantP)
		}

		gotCh, err := c.Children(10)
		if err != nil {
			t.Fatal(err)
		}
		wantCh, err := uc.Children(10)
		if err != nil {
			t.Fatal(err)
		}
		if len(gotCh) != len(wantCh) {
			t.Fatalf("Children(%v): %d != %d", c, len(gotCh), len(wantCh))
		}
		for i := range gotCh {
			if uint64(gotCh[i]) != uint64(wantCh[i]) {
				t.Fatalf("Children(%v)[%d]: %v != %v", c, i, gotCh[i], wantCh[i])
			}
		}
	}
}

func TestGridDiskParity(t *testing.T) {
	for _, ll := range randomLatLngs(200) {
		c, _ := pure.LatLngToCell(pure.LatLngDegs(ll[0], ll[1]), 7)
		uc := uber.Cell(c)
		for _, k := range []int{1, 3} {
			got, err := c.GridDisk(k)
			if err != nil {
				t.Fatal(err)
			}
			want, err := uc.GridDisk(k)
			if err != nil {
				t.Fatal(err)
			}
			gotU := make([]uint64, len(got))
			for i, x := range got {
				gotU[i] = uint64(x)
			}
			wantU := make([]uint64, len(want))
			for i, x := range want {
				wantU[i] = uint64(x)
			}
			slices.Sort(gotU)
			slices.Sort(wantU)
			if !slices.Equal(gotU, wantU) {
				t.Fatalf("GridDisk(%v, %d): sets differ (%d vs %d cells)", c, k, len(gotU), len(wantU))
			}
		}
	}
}

func TestCompactParity(t *testing.T) {
	for _, ll := range randomLatLngs(100) {
		c, _ := pure.LatLngToCell(pure.LatLngDegs(ll[0], ll[1]), 4)
		children, err := c.Children(6)
		if err != nil {
			t.Fatal(err)
		}
		got, err := pure.CompactCells(children)
		if err != nil {
			t.Fatal(err)
		}
		uberIn := make([]uber.Cell, len(children))
		for i, x := range children {
			uberIn[i] = uber.Cell(x)
		}
		want, err := uber.CompactCells(uberIn)
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != len(want) || uint64(got[0]) != uint64(want[0]) {
			t.Fatalf("CompactCells: %v vs %v", got, want)
		}
	}
}

// areaRelTol bounds the pure-Go vs cgo-binding relative difference in
// CellAreaKm2. Areas are the one scalar this suite compares with a
// tolerance rather than exactly, and this is not a claim of bit-exact
// equality: exact area equality is not portable across libm/compiler
// implementations. Both sides run the same H3 C v4.5.0 boundary-loop
// area algorithm, but this suite compares end to end — each library
// derives the cell boundary through its own trig (pure Go `math` here,
// platform libm in the binding) and then sums per-edge Cagnoli/atan2
// terms — so the tiny boundary-vertex differences are amplified in the
// area, and the more so as areas shrink with resolution.
//
// Measured across the full deterministic input set on linux/amd64 and
// darwin/arm64, the res-8 cells this test compares reach at most
// ~2.1e-9 relative (cell 8803263523fffff on linux/amd64; the per-
// platform, per-resolution measurement is in the corrective PR that set
// this constant). areaRelTol = 1e-8 is ~4.7x that measured maximum — a
// deliberate, evidence-based margin that still stays far below any
// meaningful algorithm/constant/version skew (area scales with the
// square of the Earth radius, so even a ~5e-9 relative radius or
// constant error is caught at 1e-8; a genuine algorithm or version
// change moves areas by orders of magnitude more). The root cgo parity
// suite (make test-c2go) is the correctness anchor for the algorithm
// itself, but even there cell-area comparisons are tolerance-based for
// the same reason (latLng__cellAreaKm2_parity_test.go uses an absolute
// km^2 tolerance; area_geoLoopAreaRads2_parity_test.go a ~1e-14 relative
// one on identical input loops).
const areaRelTol = 1e-8

func TestMetricsParity(t *testing.T) {
	cells := make([]pure.Cell, 0, 205)
	for _, ll := range randomLatLngs(200) {
		c, err := pure.LatLngToCell(pure.LatLngDegs(ll[0], ll[1]), 8)
		if err != nil {
			t.Fatal(err)
		}
		cells = append(cells, c)
	}
	// Explicit regression pin: 8803263523fffff maximized the res-8
	// cross-libm area difference (~2.11e-9 relative on linux/amd64) and
	// deterministically failed the previous 1e-9 tolerance in Nightly, so
	// its coverage must not depend on the pseudo-random sequence.
	cells = append(cells, pure.Cell(0x8803263523fffff))

	for _, c := range cells {
		gotA, err := c.AreaKm2()
		if err != nil {
			t.Fatal(err)
		}
		wantA, err := uber.CellAreaKm2(uber.Cell(c))
		if err != nil {
			t.Fatal(err)
		}
		if rel := math.Abs(gotA-wantA) / wantA; rel > areaRelTol {
			t.Fatalf("AreaKm2(%x): pure %.17g != uber %.17g (rel %.3e > %.0e)",
				uint64(c), gotA, wantA, rel, areaRelTol)
		}
	}
}
