package uberbench

import (
	"math"
	"slices"
	"testing"

	pure "github.com/dimchansky/h3-go"
	uber "github.com/uber/h3-go/v4"
)

// The tests in this file gate the benchmarks: every operation pairing that
// bench_test.go times is first shown to produce semantically equivalent
// results on the identical benchmark datasets. Comparing the speed of two
// functions that disagree would be meaningless.
//
// Tolerances:
//   - cells, edges, vertexes: exact (same uint64 bit patterns);
//   - coordinates: 1e-9 degrees (~0.1 mm; both libraries do the same
//     radians->degrees multiply, observed agreement is ~1e-13);
//   - areas, lengths, distances: 1e-8 relative. Bit-exactness is not
//     expected for the trigonometry-heavy measurement functions: the C
//     compiler may contract multiply-adds into FMAs differently than the
//     Go compiler, and the area/length formulas amplify the resulting
//     last-ulp differences through cancellation (observed maximum is
//     ~2e-9 relative — sub-micrometer at res-9 scales; the root module's
//     own C parity suite compares these functions with tolerances for the
//     same reason).
const (
	degTol = 1e-9
	relTol = 1e-8
)

func relClose(a, b float64) bool {
	if a == b {
		return true
	}
	den := math.Max(math.Abs(a), math.Abs(b))
	return math.Abs(a-b) <= relTol*den
}

func sortedU64s[T ~uint64](in []T) []uint64 {
	out := make([]uint64, len(in))
	for i, v := range in {
		out[i] = uint64(v)
	}
	slices.Sort(out)
	return out
}

func sortedI64s[T ~int64](in []T) []uint64 {
	out := make([]uint64, len(in))
	for i, v := range in {
		out[i] = uint64(v)
	}
	slices.Sort(out)
	return out
}

// pruneZeros drops H3_NULL slots from an uber result. The pure-Go library
// prunes null slots at the public boundary (docs/DEVIATIONS.md item 5);
// uber/h3-go prunes in most collection APIs but keeps zeros inside
// GridDiskDistances rings.
func pruneZeros(in []uber.Cell) []uber.Cell {
	out := make([]uber.Cell, 0, len(in))
	for _, c := range in {
		if c != 0 {
			out = append(out, c)
		}
	}
	return out
}

func TestEquivalenceLatLngToCell(t *testing.T) {
	for _, res := range []int{0, 5, benchRes, 15} {
		for i := range lls {
			got, err := pure.LatLngToCell(llsPure[i], res)
			if err != nil {
				t.Fatalf("pure LatLngToCell(%v, %d): %v", lls[i], res, err)
			}
			want, err := uber.LatLngToCell(llsUber[i], res)
			if err != nil {
				t.Fatalf("uber LatLngToCell(%v, %d): %v", lls[i], res, err)
			}
			if uint64(got) != uint64(want) {
				t.Fatalf("LatLngToCell(%v, %d): pure %v != uber %v", lls[i], res, got, want)
			}
		}
	}
}

func TestEquivalenceCellToLatLng(t *testing.T) {
	for i, c := range cellsPure9 {
		g, err := c.LatLng()
		if err != nil {
			t.Fatalf("pure LatLng(%v): %v", c, err)
		}
		u, err := cellsUber9[i].LatLng()
		if err != nil {
			t.Fatalf("uber LatLng(%v): %v", c, err)
		}
		if math.Abs(g.Lat.Deg()-u.Lat) > degTol || math.Abs(g.Lng.Deg()-u.Lng) > degTol {
			t.Fatalf("LatLng(%v): pure (%v, %v) != uber (%v, %v)",
				c, g.Lat.Deg(), g.Lng.Deg(), u.Lat, u.Lng)
		}
	}
}

func TestEquivalenceCellToBoundary(t *testing.T) {
	for i, c := range cellsPure9 {
		pb, err := c.Boundary()
		if err != nil {
			t.Fatalf("pure Boundary(%v): %v", c, err)
		}
		ub, err := cellsUber9[i].Boundary()
		if err != nil {
			t.Fatalf("uber Boundary(%v): %v", c, err)
		}
		if pb.Len() != len(ub) {
			t.Fatalf("Boundary(%v): pure %d verts != uber %d verts", c, pb.Len(), len(ub))
		}
		for j := range ub {
			pv := pb.At(j)
			if math.Abs(pv.Lat.Deg()-ub[j].Lat) > degTol || math.Abs(pv.Lng.Deg()-ub[j].Lng) > degTol {
				t.Fatalf("Boundary(%v) vert %d: pure (%v, %v) != uber (%v, %v)",
					c, j, pv.Lat.Deg(), pv.Lng.Deg(), ub[j].Lat, ub[j].Lng)
			}
		}
	}
}

func TestEquivalenceHierarchy(t *testing.T) {
	for i, c := range cellsPure9 {
		u := cellsUber9[i]

		pp, err := c.Parent(4)
		if err != nil {
			t.Fatalf("pure Parent: %v", err)
		}
		up, err := u.Parent(4)
		if err != nil {
			t.Fatalf("uber Parent: %v", err)
		}
		if uint64(pp) != uint64(up) {
			t.Fatalf("Parent(%v, 4): pure %v != uber %v", c, pp, up)
		}

		pcc, err := c.CenterChild(11)
		if err != nil {
			t.Fatalf("pure CenterChild: %v", err)
		}
		ucc, err := u.CenterChild(11)
		if err != nil {
			t.Fatalf("uber CenterChild: %v", err)
		}
		if uint64(pcc) != uint64(ucc) {
			t.Fatalf("CenterChild(%v, 11): pure %v != uber %v", c, pcc, ucc)
		}

		ppos, err := c.ChildPos(4)
		if err != nil {
			t.Fatalf("pure ChildPos: %v", err)
		}
		upos, err := u.ChildPos(4)
		if err != nil {
			t.Fatalf("uber ChildPos: %v", err)
		}
		if ppos != int64(upos) {
			t.Fatalf("ChildPos(%v, 4): pure %d != uber %d", c, ppos, upos)
		}

		pac, err := pp.ChildAtPos(ppos, benchRes)
		if err != nil {
			t.Fatalf("pure ChildAtPos: %v", err)
		}
		uac, err := up.ChildPosToCell(upos, benchRes)
		if err != nil {
			t.Fatalf("uber ChildPosToCell: %v", err)
		}
		if uint64(pac) != uint64(uac) || uint64(pac) != uint64(c) {
			t.Fatalf("ChildAtPos round-trip(%v): pure %v, uber %v", c, pac, uac)
		}
	}
}

func TestEquivalenceChildren(t *testing.T) {
	parents := []pure.Cell{res4Cell, cellsPure5[4], pentagons9[0]}
	for _, p := range parents {
		for depth := 1; depth <= 3; depth++ {
			res := p.Resolution() + depth
			pc, err := p.Children(res)
			if err != nil {
				t.Fatalf("pure Children(%v, %d): %v", p, res, err)
			}
			uc, err := uber.Cell(uint64(p)).Children(res)
			if err != nil {
				t.Fatalf("uber Children(%v, %d): %v", p, res, err)
			}
			// Children order is canonical in both libraries.
			if len(pc) != len(uc) {
				t.Fatalf("Children(%v, %d): pure %d cells != uber %d cells", p, res, len(pc), len(uc))
			}
			for i := range pc {
				if uint64(pc[i]) != uint64(uc[i]) {
					t.Fatalf("Children(%v, %d)[%d]: pure %v != uber %v", p, res, i, pc[i], uc[i])
				}
			}
		}
	}
}

func TestEquivalenceGridDisk(t *testing.T) {
	origins := append(slices.Clone(cellsPure9[:64]), pentagons9...)
	for _, c := range origins {
		u := uber.Cell(uint64(c))
		for _, k := range []int{1, 5, 20} {
			pd, err := c.GridDisk(k)
			if err != nil {
				t.Fatalf("pure GridDisk(%v, %d): %v", c, k, err)
			}
			ud, err := u.GridDisk(k)
			if err != nil {
				t.Fatalf("uber GridDisk(%v, %d): %v", c, k, err)
			}
			// Both libraries prune null slots; order is unspecified.
			if !slices.Equal(sortedU64s(pd), sortedI64s(ud)) {
				t.Fatalf("GridDisk(%v, %d): pure and uber disagree (len %d vs %d)", c, k, len(pd), len(ud))
			}
		}
	}
}

func TestEquivalenceGridDiskDistances(t *testing.T) {
	origins := append(slices.Clone(cellsPure9[:64]), pentagons9...)
	const k = 5
	for _, c := range origins {
		pd, pdist, err := c.GridDiskDistances(k)
		if err != nil {
			t.Fatalf("pure GridDiskDistances(%v, %d): %v", c, k, err)
		}
		urings, err := uber.Cell(uint64(c)).GridDiskDistances(k)
		if err != nil {
			t.Fatalf("uber GridDiskDistances(%v, %d): %v", c, k, err)
		}
		// Group the pure result by distance so it is shaped like uber's
		// rings. uber keeps H3_NULL slots inside rings for pentagon-affected
		// disks; the pure library prunes them, so drop zeros before
		// comparing.
		prings := make([][]pure.Cell, k+1)
		for i, cell := range pd {
			d := pdist[i]
			prings[d] = append(prings[d], cell)
		}
		if len(urings) != k+1 {
			t.Fatalf("uber GridDiskDistances(%v, %d): %d rings, want %d", c, k, len(urings), k+1)
		}
		for d := 0; d <= k; d++ {
			uring := pruneZeros(urings[d])
			if !slices.Equal(sortedU64s(prings[d]), sortedI64s(uring)) {
				t.Fatalf("GridDiskDistances(%v, %d) ring %d: pure %d cells != uber %d cells",
					c, k, d, len(prings[d]), len(uring))
			}
		}
	}
}

func TestEquivalenceGridRing(t *testing.T) {
	origins := append(slices.Clone(cellsPure9[:64]), pentagons9...)
	for _, c := range origins {
		for _, k := range []int{1, 5} {
			pr, err := c.GridRing(k)
			if err != nil {
				t.Fatalf("pure GridRing(%v, %d): %v", c, k, err)
			}
			ur, err := uber.Cell(uint64(c)).GridRing(k)
			if err != nil {
				t.Fatalf("uber GridRing(%v, %d): %v", c, k, err)
			}
			if !slices.Equal(sortedU64s(pr), sortedI64s(ur)) {
				t.Fatalf("GridRing(%v, %d): pure and uber disagree (len %d vs %d)", c, k, len(pr), len(ur))
			}
		}
	}
}

func TestEquivalenceGridDisksUnsafe(t *testing.T) {
	origins := hexCells9[:64]
	originsUber := hexCellsUber9[:64]
	const k = 2
	flat, err := pure.GridDisksUnsafe(origins, k)
	if err != nil {
		t.Fatalf("pure GridDisksUnsafe: %v", err)
	}
	nested, err := uber.GridDisksUnsafe(originsUber, k)
	if err != nil {
		t.Fatalf("uber GridDisksUnsafe: %v", err)
	}
	if len(nested) != len(origins) {
		t.Fatalf("uber GridDisksUnsafe: %d disks, want %d", len(nested), len(origins))
	}
	size, err := pure.MaxGridDiskSize(k)
	if err != nil {
		t.Fatalf("MaxGridDiskSize: %v", err)
	}
	group := int(size)
	if len(flat) != group*len(origins) {
		t.Fatalf("pure GridDisksUnsafe: flat len %d, want %d", len(flat), group*len(origins))
	}
	for i := range origins {
		// pure returns fixed-size groups (unpruned, C layout); uber prunes
		// zeros per disk. Normalize both to sorted non-null sets.
		pdisk := make([]pure.Cell, 0, group)
		for _, c := range flat[i*group : (i+1)*group] {
			if c != 0 {
				pdisk = append(pdisk, c)
			}
		}
		if !slices.Equal(sortedU64s(pdisk), sortedI64s(nested[i])) {
			t.Fatalf("GridDisksUnsafe disk %d: pure and uber disagree", i)
		}
	}
}

func TestEquivalenceGridDistanceAndNeighbors(t *testing.T) {
	for _, pair := range neighborPairs {
		a, b := pair[0], pair[1]
		ua, ub := uber.Cell(uint64(a)), uber.Cell(uint64(b))

		pd, err := a.GridDistance(b)
		if err != nil {
			t.Fatalf("pure GridDistance: %v", err)
		}
		ud, err := ua.GridDistance(ub)
		if err != nil {
			t.Fatalf("uber GridDistance: %v", err)
		}
		if pd != ud || pd != 1 {
			t.Fatalf("GridDistance(%v, %v): pure %d, uber %d, want 1", a, b, pd, ud)
		}

		pn, err := a.IsNeighbor(b)
		if err != nil {
			t.Fatalf("pure IsNeighbor: %v", err)
		}
		un, err := ua.IsNeighbor(ub)
		if err != nil {
			t.Fatalf("uber IsNeighbor: %v", err)
		}
		if pn != un || !pn {
			t.Fatalf("IsNeighbor(%v, %v): pure %v, uber %v", a, b, pn, un)
		}
	}
}

func TestEquivalenceGridPath(t *testing.T) {
	if len(pathPairs) == 0 {
		t.Fatal("no path pairs in dataset")
	}
	for _, pair := range pathPairs {
		pp, err := pair[0].GridPath(pair[1])
		if err != nil {
			t.Fatalf("pure GridPath: %v", err)
		}
		up, err := uber.Cell(uint64(pair[0])).GridPath(uber.Cell(uint64(pair[1])))
		if err != nil {
			t.Fatalf("uber GridPath: %v", err)
		}
		if len(pp) != len(up) {
			t.Fatalf("GridPath(%v, %v): pure len %d != uber len %d", pair[0], pair[1], len(pp), len(up))
		}
		for i := range pp {
			if uint64(pp[i]) != uint64(up[i]) {
				t.Fatalf("GridPath(%v, %v)[%d]: pure %v != uber %v", pair[0], pair[1], i, pp[i], up[i])
			}
		}
	}
}

func TestEquivalenceLocalIJ(t *testing.T) {
	for _, pair := range neighborPairs[:64] {
		origin, c := pair[0], pair[1]
		pij, err := pure.CellToLocalIJ(origin, c)
		if err != nil {
			t.Fatalf("pure CellToLocalIJ: %v", err)
		}
		uij, err := uber.CellToLocalIJ(uber.Cell(uint64(origin)), uber.Cell(uint64(c)))
		if err != nil {
			t.Fatalf("uber CellToLocalIJ: %v", err)
		}
		if int(pij.I) != uij.I || int(pij.J) != uij.J {
			t.Fatalf("CellToLocalIJ(%v, %v): pure %+v != uber %+v", origin, c, pij, uij)
		}
		back, err := pure.LocalIJToCell(origin, pij)
		if err != nil {
			t.Fatalf("pure LocalIJToCell: %v", err)
		}
		uback, err := uber.LocalIJToCell(uber.Cell(uint64(origin)), uij)
		if err != nil {
			t.Fatalf("uber LocalIJToCell: %v", err)
		}
		if uint64(back) != uint64(uback) || uint64(back) != uint64(c) {
			t.Fatalf("LocalIJToCell round-trip: pure %v, uber %v, want %v", back, uback, c)
		}
	}
}

func TestEquivalencePolygonToCells(t *testing.T) {
	cases := []struct {
		name string
		pp   pure.GeoPolygon
		up   uber.GeoPolygon
		res  int
	}{
		{"sf/res=7", sfPolygonPure, sfPolygonUber, 7},
		{"sf/res=9", sfPolygonPure, sfPolygonUber, 9},
		{"sf/res=11", sfPolygonPure, sfPolygonUber, 11},
		{"sf-hole/res=9", sfHolePolygonPure, sfHolePolygonUber, 9},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			pc, err := pure.PolygonToCells(tc.pp, tc.res)
			if err != nil {
				t.Fatalf("pure PolygonToCells: %v", err)
			}
			uc, err := uber.PolygonToCells(tc.up, tc.res)
			if err != nil {
				t.Fatalf("uber PolygonToCells: %v", err)
			}
			if !slices.Equal(sortedU64s(pc), sortedI64s(uc)) {
				t.Fatalf("PolygonToCells: pure %d cells != uber %d cells", len(pc), len(uc))
			}
		})
	}
}

func TestEquivalencePolygonToCellsExperimental(t *testing.T) {
	modes := []struct {
		name string
		pm   pure.ContainmentMode
		um   uber.ContainmentMode
	}{
		{"center", pure.ContainmentCenter, uber.ContainmentCenter},
		{"full", pure.ContainmentFull, uber.ContainmentFull},
		{"overlapping", pure.ContainmentOverlapping, uber.ContainmentOverlapping},
	}
	for _, m := range modes {
		t.Run(m.name, func(t *testing.T) {
			pc, err := pure.PolygonToCellsExperimental(sfHolePolygonPure, 9, m.pm)
			if err != nil {
				t.Fatalf("pure PolygonToCellsExperimental: %v", err)
			}
			uc, err := uber.PolygonToCellsExperimental(sfHolePolygonUber, 9, m.um)
			if err != nil {
				t.Fatalf("uber PolygonToCellsExperimental: %v", err)
			}
			if !slices.Equal(sortedU64s(pc), sortedI64s(uc)) {
				t.Fatalf("PolygonToCellsExperimental(%s): pure %d cells != uber %d cells", m.name, len(pc), len(uc))
			}
		})
	}
}

func TestEquivalenceCompactUncompact(t *testing.T) {
	pcomp, err := pure.CompactCells(sfCells9)
	if err != nil {
		t.Fatalf("pure CompactCells: %v", err)
	}
	ucomp, err := uber.CompactCells(sfCells9Uber)
	if err != nil {
		t.Fatalf("uber CompactCells: %v", err)
	}
	if !slices.Equal(sortedU64s(pcomp), sortedI64s(ucomp)) {
		t.Fatalf("CompactCells: pure %d cells != uber %d cells", len(pcomp), len(ucomp))
	}

	punc, err := pure.UncompactCells(sfCompacted, 9)
	if err != nil {
		t.Fatalf("pure UncompactCells: %v", err)
	}
	uunc, err := uber.UncompactCells(sfCompactedUber, 9)
	if err != nil {
		t.Fatalf("uber UncompactCells: %v", err)
	}
	if !slices.Equal(sortedU64s(punc), sortedI64s(uunc)) {
		t.Fatalf("UncompactCells: pure %d cells != uber %d cells", len(punc), len(uunc))
	}
	if !slices.Equal(sortedU64s(punc), sortedU64s(sfCells9)) {
		t.Fatal("UncompactCells round-trip does not reproduce the input set")
	}

	// The expansion workload benchmarked as Uncompact/res=4to9.
	pexp, err := pure.UncompactCells([]pure.Cell{res4Cell}, 9)
	if err != nil {
		t.Fatalf("pure UncompactCells(res4): %v", err)
	}
	uexp, err := uber.UncompactCells([]uber.Cell{res4CellUber}, 9)
	if err != nil {
		t.Fatalf("uber UncompactCells(res4): %v", err)
	}
	if !slices.Equal(sortedU64s(pexp), sortedI64s(uexp)) {
		t.Fatal("UncompactCells(res4, 9): pure and uber disagree")
	}
}

func TestEquivalenceCellsToMultiPolygon(t *testing.T) {
	for _, tc := range []struct {
		name  string
		cells []pure.Cell
	}{
		{"disk331", diskCells331},
		{"sf-compacted-uncompacted", sfCells9},
	} {
		t.Run(tc.name, func(t *testing.T) {
			pmp, err := pure.CellsToMultiPolygon(tc.cells)
			if err != nil {
				t.Fatalf("pure CellsToMultiPolygon: %v", err)
			}
			ump, err := uber.CellsToMultiPolygon(toUberCells(tc.cells))
			if err != nil {
				t.Fatalf("uber CellsToMultiPolygon: %v", err)
			}
			if len(pmp) != len(ump) {
				t.Fatalf("CellsToMultiPolygon: pure %d polygons != uber %d", len(pmp), len(ump))
			}
			for i := range pmp {
				ploops := append([]pure.GeoLoop{pmp[i].GeoLoop}, pmp[i].Holes...)
				uloops := append([]uber.GeoLoop{ump[i].GeoLoop}, ump[i].Holes...)
				if len(ploops) != len(uloops) {
					t.Fatalf("polygon %d: pure %d loops != uber %d loops", i, len(ploops), len(uloops))
				}
				for j := range ploops {
					if len(ploops[j]) != len(uloops[j]) {
						t.Fatalf("polygon %d loop %d: pure %d verts != uber %d verts",
							i, j, len(ploops[j]), len(uloops[j]))
					}
					for v := range ploops[j] {
						if math.Abs(ploops[j][v].Lat.Deg()-uloops[j][v].Lat) > degTol ||
							math.Abs(ploops[j][v].Lng.Deg()-uloops[j][v].Lng) > degTol {
							t.Fatalf("polygon %d loop %d vert %d differs", i, j, v)
						}
					}
				}
			}
		})
	}
}

func TestEquivalenceEdgesAndVertexes(t *testing.T) {
	cells := append(slices.Clone(hexCells9[:64]), pentagons9...)
	for _, c := range cells {
		u := uber.Cell(uint64(c))

		pe, err := c.DirectedEdges()
		if err != nil {
			t.Fatalf("pure DirectedEdges: %v", err)
		}
		ue, err := u.DirectedEdges()
		if err != nil {
			t.Fatalf("uber DirectedEdges: %v", err)
		}
		if !slices.Equal(sortedU64s(pe), sortedI64s(ue)) {
			t.Fatalf("DirectedEdges(%v): pure %d != uber %d", c, len(pe), len(ue))
		}

		for _, e := range pe {
			ueid := uber.DirectedEdge(uint64(e))

			po, err := e.Origin()
			if err != nil {
				t.Fatalf("pure Origin: %v", err)
			}
			uo, err := ueid.Origin()
			if err != nil {
				t.Fatalf("uber Origin: %v", err)
			}
			if uint64(po) != uint64(uo) {
				t.Fatalf("Origin(%v): pure %v != uber %v", e, po, uo)
			}

			pb, err := e.Boundary()
			if err != nil {
				t.Fatalf("pure edge Boundary: %v", err)
			}
			ub, err := ueid.Boundary()
			if err != nil {
				t.Fatalf("uber edge Boundary: %v", err)
			}
			if pb.Len() != len(ub) {
				t.Fatalf("edge Boundary(%v): pure %d verts != uber %d", e, pb.Len(), len(ub))
			}
			for j := range ub {
				pv := pb.At(j)
				if math.Abs(pv.Lat.Deg()-ub[j].Lat) > degTol || math.Abs(pv.Lng.Deg()-ub[j].Lng) > degTol {
					t.Fatalf("edge Boundary(%v) vert %d differs", e, j)
				}
			}
		}

		pv, err := c.Vertexes()
		if err != nil {
			t.Fatalf("pure Vertexes: %v", err)
		}
		uv, err := u.Vertexes()
		if err != nil {
			t.Fatalf("uber Vertexes: %v", err)
		}
		if !slices.Equal(sortedU64s(pv), sortedI64s(uv)) {
			t.Fatalf("Vertexes(%v): pure %d != uber %d", c, len(pv), len(uv))
		}

		for _, vtx := range pv {
			pg, err := vtx.LatLng()
			if err != nil {
				t.Fatalf("pure vertex LatLng: %v", err)
			}
			ug, err := uber.Vertex(uint64(vtx)).LatLng()
			if err != nil {
				t.Fatalf("uber vertex LatLng: %v", err)
			}
			if math.Abs(pg.Lat.Deg()-ug.Lat) > degTol || math.Abs(pg.Lng.Deg()-ug.Lng) > degTol {
				t.Fatalf("vertex LatLng(%v) differs", vtx)
			}
		}
	}
}

func TestEquivalenceMetrics(t *testing.T) {
	for i, c := range cellsPure9[:128] {
		u := cellsUber9[i]

		pa, err := c.AreaKm2()
		if err != nil {
			t.Fatalf("pure AreaKm2: %v", err)
		}
		ua, err := uber.CellAreaKm2(u)
		if err != nil {
			t.Fatalf("uber CellAreaKm2: %v", err)
		}
		if !relClose(pa, ua) {
			t.Fatalf("AreaKm2(%v): pure %v != uber %v", c, pa, ua)
		}
	}

	for _, pair := range neighborPairs[:64] {
		e, err := pair[0].DirectedEdgeTo(pair[1])
		if err != nil {
			t.Fatalf("pure DirectedEdgeTo: %v", err)
		}
		ue, err := uber.Cell(uint64(pair[0])).DirectedEdge(uber.Cell(uint64(pair[1])))
		if err != nil {
			t.Fatalf("uber DirectedEdge: %v", err)
		}
		if uint64(e) != uint64(ue) {
			t.Fatalf("DirectedEdgeTo: pure %v != uber %v", e, ue)
		}
		pl, err := e.LengthM()
		if err != nil {
			t.Fatalf("pure LengthM: %v", err)
		}
		ul, err := uber.EdgeLengthM(ue)
		if err != nil {
			t.Fatalf("uber EdgeLengthM: %v", err)
		}
		if !relClose(pl, ul) {
			t.Fatalf("LengthM(%v): pure %v != uber %v", e, pl, ul)
		}
	}

	for i := 0; i+1 < 64; i += 2 {
		pd := pure.GreatCircleDistanceKm(llsPure[i], llsPure[i+1])
		ud := uber.GreatCircleDistanceKm(llsUber[i], llsUber[i+1])
		if !relClose(pd, ud) {
			t.Fatalf("GreatCircleDistanceKm(%v, %v): pure %v != uber %v", lls[i], lls[i+1], pd, ud)
		}
	}
}

func TestEquivalenceStrings(t *testing.T) {
	for i, c := range cellsPure9 {
		if c.String() != cellsUber9[i].String() {
			t.Fatalf("String(%v): pure %q != uber %q", uint64(c), c.String(), cellsUber9[i].String())
		}
	}
	for _, s := range cellStrings {
		pc, err := pure.ParseCell(s)
		if err != nil {
			t.Fatalf("pure ParseCell(%q): %v", s, err)
		}
		uc := uber.CellFromString(s)
		if !uc.IsValid() {
			t.Fatalf("uber CellFromString(%q): invalid", s)
		}
		if uint64(pc) != uint64(uc) {
			t.Fatalf("parse(%q): pure %v != uber %v", s, pc, uc)
		}
	}
	// Invalid inputs: pure reports an error; uber returns a zero index that
	// fails IsValid. The benchmarked composite (parse + validate) must agree
	// on acceptance.
	for _, s := range []string{"", "zzz", "0x", "ffffffffffffffff1", "8928308280fffff9999"} {
		_, perr := pure.ParseCell(s)
		uc := uber.CellFromString(s)
		if (perr == nil) != uc.IsValid() {
			t.Fatalf("parse acceptance mismatch for %q: pure err=%v, uber valid=%v", s, perr, uc.IsValid())
		}
	}
}

func TestEquivalenceValidity(t *testing.T) {
	inputs := make([]uint64, 0, len(cellsPure9)+6)
	for _, c := range cellsPure9 {
		inputs = append(inputs, uint64(c))
	}
	inputs = append(inputs, 0, 1, 0x7fffffffffffffff,
		uint64(cellsPure9[0])|1<<63, // high bit set
		0x2928308280fffff,           // wrong mode bits
		uint64(cellsPure9[0])^0x7,   // corrupted digit
	)
	for _, raw := range inputs {
		pv := pure.Cell(raw).IsValid()
		uv := uber.Cell(raw).IsValid()
		if pv != uv {
			t.Fatalf("IsValid(%#x): pure %v != uber %v", raw, pv, uv)
		}
	}
}

// TestEquivalenceServiceWorkload pins the aggregate result of the
// service-style batch benchmark so both implementations demonstrably do the
// same work.
func TestEquivalenceServiceWorkload(t *testing.T) {
	p := serviceWorkloadPureAlloc()
	u := serviceWorkloadUber()
	if p != u {
		t.Fatalf("service workload checksum: pure %d != uber %d", p, u)
	}
	if _, warm := serviceWorkloadPureWarm(make([]pure.Cell, 0, 8)); warm != p {
		t.Fatalf("service workload checksum (warm): %d != %d", warm, p)
	}
}

// TestEquivalenceMemWorkloads pins that every process-level memory workload
// (cmd/memprobe) computes the identical checksum in both implementations.
func TestEquivalenceMemWorkloads(t *testing.T) {
	for _, w := range MemWorkloads {
		if w.Name == "scalar-1m" && testing.Short() {
			continue
		}
		t.Run(w.Name, func(t *testing.T) {
			p := w.Pure(1)
			u := w.Uber(1)
			if p.Checksum != u.Checksum {
				t.Fatalf("checksum: pure %#x != uber %#x", p.Checksum, u.Checksum)
			}
		})
	}
}
