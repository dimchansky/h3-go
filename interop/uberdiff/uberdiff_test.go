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

func TestMetricsParity(t *testing.T) {
	for _, ll := range randomLatLngs(200) {
		c, _ := pure.LatLngToCell(pure.LatLngDegs(ll[0], ll[1]), 8)
		uc := uber.Cell(c)
		gotA, err := c.AreaKm2()
		if err != nil {
			t.Fatal(err)
		}
		wantA, err := uber.CellAreaKm2(uc)
		if err != nil {
			t.Fatal(err)
		}
		// This module pins exact index equality but only near-equality for
		// areas: both sides now implement the H3 C v4.5.0 algorithm, but
		// the binding computes through cgo-compiled C while this library
		// is pure Go, and floating-point codegen/libm differences reach
		// ~2e-12 relative near pentagons (measured across the pentagon
		// 2-disks of res 0-6 at the v4.5.0/v4.5.0 pairing). Exact v4.5.0
		// equality is enforced by the cgo parity suite in the root module.
		if math.Abs(gotA-wantA)/wantA > 1e-9 {
			t.Fatalf("AreaKm2(%v): %v != %v", c, gotA, wantA)
		}
	}
}
