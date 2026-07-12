// Package uberbench compares this pure-Go H3 implementation against the
// official cgo binding github.com/uber/h3-go on semantically equivalent
// operations over identical, deterministic datasets.
//
// The package has three parts:
//
//   - fixtures.go (this file): deterministic input datasets shared by the
//     equivalence tests, the benchmarks, and the memprobe command;
//   - equivalence_test.go: correctness checks that both libraries return
//     semantically equivalent results on every benchmarked pairing — these
//     gate the benchmarks, which are meaningless if the outputs differ;
//   - bench_test.go: the comparative benchmarks themselves.
//
// This is a separate Go module so the root library keeps zero dependencies
// and zero cgo; see README.md for methodology, what the numbers do and do
// not show, and how to run everything (make bench-uber).
package uberbench

import (
	"math"
	"math/rand/v2"

	pure "github.com/dimchansky/h3-go"
	uber "github.com/uber/h3-go/v4"
)

// benchRes is the primary resolution for scalar benchmarks (~0.1 km² cells,
// a typical service resolution).
const benchRes = 9

// numFixedLatLngs is the size of the shared coordinate dataset. Benchmarks
// cycle through it so results are not specific to a single location.
const numFixedLatLngs = 1024

// fixedLatLngs returns a deterministic, globe-covering coordinate set in
// degrees: fixed polar and antimeridian probes followed by PCG-seeded
// pseudo-random points. The seed is fixed so every run, machine, and
// implementation sees the identical inputs.
func fixedLatLngs(n int) [][2]float64 {
	rng := rand.New(rand.NewPCG(0xbe9c, 0x0dda7a))
	out := make([][2]float64, 0, n)
	out = append(out,
		[2]float64{89.9, 0},      // near north pole
		[2]float64{-89.9, 45},    // near south pole
		[2]float64{0, 179.9},     // antimeridian, east side
		[2]float64{0, -179.9},    // antimeridian, west side
		[2]float64{37.7, -122.4}, // ordinary mid-latitude city
	)
	for len(out) < n {
		out = append(out, [2]float64{rng.Float64()*180 - 90, rng.Float64()*360 - 180})
	}
	return out
}

// sfVerts is the San Francisco test polygon from the upstream H3 C test
// suite (radians there; converted to degrees once here so both libraries
// receive bit-identical degree inputs).
var sfVerts = [][2]float64{
	{0.659966917655, -2.1364398519396},
	{0.6595011102219, -2.1359434279405},
	{0.6583348114025, -2.1354884206045},
	{0.6581220034068, -2.1382437718946},
	{0.6594479998527, -2.1384597563896},
	{0.6599990002976, -2.1376771158464},
}

// sfHoleVerts is the upstream "SF hole" loop, inside sfVerts.
var sfHoleVerts = [][2]float64{
	{0.6595072188743, -2.1371053983433},
	{0.6591482046471, -2.1373141048153},
	{0.6592295020837, -2.1365222838402},
}

func radsToDegsPair(v [2]float64) (lat, lng float64) {
	return v[0] * 180 / math.Pi, v[1] * 180 / math.Pi
}

func pureLoop(verts [][2]float64) pure.GeoLoop {
	loop := make(pure.GeoLoop, len(verts))
	for i, v := range verts {
		lat, lng := radsToDegsPair(v)
		loop[i] = pure.LatLngDegs(lat, lng)
	}
	return loop
}

func uberLoop(verts [][2]float64) uber.GeoLoop {
	loop := make(uber.GeoLoop, len(verts))
	for i, v := range verts {
		lat, lng := radsToDegsPair(v)
		loop[i] = uber.NewLatLng(lat, lng)
	}
	return loop
}

// Shared fixtures. Everything is derived deterministically from
// fixedLatLngs and sfVerts. The equivalence tests assert that both
// libraries agree on all derived data, so it does not matter which
// implementation builds a fixture.
var (
	lls = fixedLatLngs(numFixedLatLngs)

	llsPure = func() []pure.LatLng {
		out := make([]pure.LatLng, len(lls))
		for i, ll := range lls {
			out[i] = pure.LatLngDegs(ll[0], ll[1])
		}
		return out
	}()

	llsUber = func() []uber.LatLng {
		out := make([]uber.LatLng, len(lls))
		for i, ll := range lls {
			out[i] = uber.NewLatLng(ll[0], ll[1])
		}
		return out
	}()

	cellsPure5  = mustCellsPure(5)
	cellsPure9  = mustCellsPure(benchRes)
	cellsPure15 = mustCellsPure(15)

	cellsUber5  = toUberCells(cellsPure5)
	cellsUber9  = toUberCells(cellsPure9)
	cellsUber15 = toUberCells(cellsPure15)

	// cellStrings are the canonical string forms of cellsPure9, for the
	// parse benchmarks.
	cellStrings = func() []string {
		out := make([]string, len(cellsPure9))
		for i, c := range cellsPure9 {
			out[i] = c.String()
		}
		return out
	}()

	// hexCells9 excludes pentagons and pentagon-adjacent cells so that
	// benchmarks of the *Unsafe variants and of edge/vertex null-slot
	// semantics run on inputs where both libraries return identical
	// shapes. Random coordinates essentially never hit pentagons, but the
	// fixed polar probes can be near one, so filter deterministically.
	hexCells9 = func() []pure.Cell {
		out := make([]pure.Cell, 0, len(cellsPure9))
		for _, c := range cellsPure9 {
			if c.IsPentagon() {
				continue
			}
			disk, err := c.GridDiskUnsafe(2)
			if err != nil { // pentagon within k=2
				continue
			}
			_ = disk
			out = append(out, c)
		}
		return out
	}()

	hexCellsUber9 = toUberCells(hexCells9)

	// neighborPairs are adjacent cell pairs (grid distance 1).
	neighborPairs = func() [][2]pure.Cell {
		out := make([][2]pure.Cell, 0, len(hexCells9))
		for _, c := range hexCells9 {
			ring, err := c.GridRing(1)
			if err != nil || len(ring) == 0 {
				continue
			}
			out = append(out, [2]pure.Cell{c, ring[0]})
		}
		return out
	}()

	// pathPairs are cell pairs a few dozen grid steps apart for which
	// gridPathCells succeeds. Distance ~0.05° keeps both endpoints on the
	// same icosahedron face for most inputs; failing pairs are skipped
	// deterministically (both libraries fail identically; asserted in the
	// equivalence tests).
	pathPairs = func() [][2]pure.Cell {
		out := make([][2]pure.Cell, 0, 256)
		for i, ll := range lls {
			if len(out) == 256 {
				break
			}
			a := cellsPure9[i]
			b, err := pure.LatLngToCell(pure.LatLngDegs(ll[0]+0.05, ll[1]+0.05), benchRes)
			if err != nil {
				continue
			}
			if _, err := a.GridPath(b); err != nil {
				continue
			}
			out = append(out, [2]pure.Cell{a, b})
		}
		return out
	}()

	// pentagons9 are the twelve res-9 pentagons (edge-case dataset for
	// equivalence tests).
	pentagons9 = func() []pure.Cell {
		p, err := pure.Pentagons(benchRes)
		if err != nil {
			panic(err)
		}
		return p
	}()

	sfPolygonPure = pure.GeoPolygon{GeoLoop: pureLoop(sfVerts)}
	sfPolygonUber = uber.GeoPolygon{GeoLoop: uberLoop(sfVerts)}

	sfHolePolygonPure = pure.GeoPolygon{
		GeoLoop: pureLoop(sfVerts),
		Holes:   []pure.GeoLoop{pureLoop(sfHoleVerts)},
	}
	sfHolePolygonUber = uber.GeoPolygon{
		GeoLoop: uberLoop(sfVerts),
		Holes:   []uber.GeoLoop{uberLoop(sfHoleVerts)},
	}

	// sfCells9 is the res-9 polyfill of the SF polygon (1253 cells): the
	// realistic mixed input for compact and multi-polygon workloads.
	sfCells9 = func() []pure.Cell {
		cells, err := pure.PolygonToCells(sfPolygonPure, 9)
		if err != nil {
			panic(err)
		}
		return cells
	}()

	sfCells9Uber = toUberCells(sfCells9)

	// sfCompacted is CompactCells(sfCells9), the input for uncompact.
	sfCompacted = func() []pure.Cell {
		cells, err := pure.CompactCells(sfCells9)
		if err != nil {
			panic(err)
		}
		return cells
	}()

	sfCompactedUber = toUberCells(sfCompacted)

	// diskCells331 is a contiguous k=10 disk (331 cells) around a fixed
	// cell: the input for CellsToMultiPolygon.
	diskCells331 = func() []pure.Cell {
		cells, err := cellsPure9[4].GridDisk(10) // fixed mid-latitude city cell
		if err != nil {
			panic(err)
		}
		return cells
	}()

	diskCells331Uber = toUberCells(diskCells331)

	// res4Cell is the fixed parent for children/uncompact expansion
	// workloads (res 4 -> res 9 is 7^5 = 16807 cells).
	res4Cell = func() pure.Cell {
		p, err := cellsPure9[4].Parent(4)
		if err != nil {
			panic(err)
		}
		return p
	}()

	res4CellUber = uber.Cell(uint64(res4Cell))
)

func mustCellsPure(res int) []pure.Cell {
	out := make([]pure.Cell, len(llsPure))
	for i, g := range llsPure {
		c, err := pure.LatLngToCell(g, res)
		if err != nil {
			panic(err)
		}
		out[i] = c
	}
	return out
}

// toUberCells converts cells between the two libraries. Both encode the
// same H3 bit layout; valid indexes never set the high bit, so the
// uint64 -> int64 conversion is loss-free (asserted by equivalence tests).
func toUberCells(in []pure.Cell) []uber.Cell {
	out := make([]uber.Cell, len(in))
	for i, c := range in {
		out[i] = uber.Cell(uint64(c))
	}
	return out
}

func toPureCells(in []uber.Cell) []pure.Cell {
	out := make([]pure.Cell, len(in))
	for i, c := range in {
		out[i] = pure.Cell(uint64(c))
	}
	return out
}
