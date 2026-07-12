package uberbench

import (
	"testing"

	pure "github.com/dimchansky/h3-go"
	uber "github.com/uber/h3-go/v4"
)

// Comparative benchmarks. Naming convention:
//
//	Benchmark<Operation>[/<case>=<v>]/impl=<impl>
//
// where impl is one of:
//
//	pure       this library, convenience (allocating) API
//	pure-cold  this library, Append* form with a nil destination
//	           (allocates like the convenience form; shows the two are the
//	           same path)
//	pure-warm  this library, Append* form reusing a persistent buffer
//	           (the zero-allocation path)
//	uber       github.com/uber/h3-go/v4 (cgo binding; its APIs allocate per
//	           call — the binding has no buffer-reuse form, so `uber` is
//	           also its most efficient supported usage; where the binding
//	           has a batch API, e.g. GridDisksUnsafe, it is benchmarked
//	           too)
//
// Compare with:
//
//	benchstat -col /impl <output>
//
// Every pairing is validated for semantic equivalence by
// equivalence_test.go before these numbers mean anything. Where the two
// libraries return differently shaped results, the benchmark comments say
// so. Inputs cycle deterministic datasets from fixtures.go; the cycling
// index arithmetic is identical for both implementations.
//
// The loops deliberately use the classic b.N + package-level sink pattern
// rather than b.Loop(): the benchmarks cycle datasets by loop index, and
// keeping one uniform pattern across every pairing matters more here than
// b.Loop's convenience. Do not modernize half of them.

var (
	sinkCell   pure.Cell
	sinkUCell  uber.Cell
	sinkCells  []pure.Cell
	sinkUCells []uber.Cell
	sinkURings [][]uber.Cell
	sinkI32s   []int32
	sinkPolys  []pure.GeoPolygon
	sinkUPolys []uber.GeoPolygon
	sinkEdges  []pure.DirectedEdge
	sinkUEdges []uber.DirectedEdge
	sinkVerts  []pure.Vertex
	sinkUVerts []uber.Vertex
	sinkInt    int
	sinkF64    float64
	sinkBool   bool
	sinkStr    string
	sinkU64    uint64
	sinkLL     pure.LatLng
	sinkULL    uber.LatLng
)

const llMask = numFixedLatLngs - 1 // fixtures length is a power of two

// --- Scalar, cgo-boundary-sensitive operations ---------------------------

func BenchmarkLatLngToCell(b *testing.B) {
	b.Run("res=9/impl=pure", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			c, err := pure.LatLngToCell(llsPure[i&llMask], benchRes)
			if err != nil {
				b.Fatal(err)
			}
			sinkCell = c
		}
	})
	b.Run("res=9/impl=uber", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			c, err := uber.LatLngToCell(llsUber[i&llMask], benchRes)
			if err != nil {
				b.Fatal(err)
			}
			sinkUCell = c
		}
	})
}

func BenchmarkCellToLatLng(b *testing.B) {
	b.Run("res=9/impl=pure", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			g, err := cellsPure9[i&llMask].LatLng()
			if err != nil {
				b.Fatal(err)
			}
			sinkLL = g
		}
	})
	b.Run("res=9/impl=uber", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			g, err := cellsUber9[i&llMask].LatLng()
			if err != nil {
				b.Fatal(err)
			}
			sinkULL = g
		}
	})
}

// CellToBoundary result shapes differ: this library returns a fixed-size
// value (no heap allocation); the binding returns a freshly allocated
// []LatLng.
func BenchmarkCellToBoundary(b *testing.B) {
	b.Run("res=9/impl=pure", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			bd, err := cellsPure9[i&llMask].Boundary()
			if err != nil {
				b.Fatal(err)
			}
			sinkInt = bd.Len()
		}
	})
	b.Run("res=9/impl=uber", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			bd, err := cellsUber9[i&llMask].Boundary()
			if err != nil {
				b.Fatal(err)
			}
			sinkInt = len(bd)
		}
	})
}

func BenchmarkCellToParent(b *testing.B) {
	b.Run("res=9to7/impl=pure", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			p, err := cellsPure9[i&llMask].Parent(7)
			if err != nil {
				b.Fatal(err)
			}
			sinkCell = p
		}
	})
	b.Run("res=9to7/impl=uber", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			p, err := cellsUber9[i&llMask].Parent(7)
			if err != nil {
				b.Fatal(err)
			}
			sinkUCell = p
		}
	})
}

func BenchmarkGridDistance(b *testing.B) {
	n := len(pathPairs)
	b.Run("impl=pure", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			pair := pathPairs[i%n]
			d, err := pair[0].GridDistance(pair[1])
			if err != nil {
				b.Fatal(err)
			}
			sinkInt = d
		}
	})
	b.Run("impl=uber", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			pair := pathPairs[i%n]
			d, err := uber.Cell(uint64(pair[0])).GridDistance(uber.Cell(uint64(pair[1])))
			if err != nil {
				b.Fatal(err)
			}
			sinkInt = d
		}
	})
}

func BenchmarkIsNeighbor(b *testing.B) {
	n := len(neighborPairs)
	b.Run("impl=pure", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			pair := neighborPairs[i%n]
			v, err := pair[0].IsNeighbor(pair[1])
			if err != nil {
				b.Fatal(err)
			}
			sinkBool = v
		}
	})
	b.Run("impl=uber", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			pair := neighborPairs[i%n]
			v, err := uber.Cell(uint64(pair[0])).IsNeighbor(uber.Cell(uint64(pair[1])))
			if err != nil {
				b.Fatal(err)
			}
			sinkBool = v
		}
	})
}

func BenchmarkCellArea(b *testing.B) {
	b.Run("unit=km2/impl=pure", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			a, err := cellsPure9[i&llMask].AreaKm2()
			if err != nil {
				b.Fatal(err)
			}
			sinkF64 = a
		}
	})
	b.Run("unit=km2/impl=uber", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			a, err := uber.CellAreaKm2(cellsUber9[i&llMask])
			if err != nil {
				b.Fatal(err)
			}
			sinkF64 = a
		}
	})
}

// ParseCell semantics differ: this library validates the parsed index
// (mode + digit checks); the binding's CellFromString swallows syntax
// errors and returns 0. The uber loop therefore adds the IsValid call a
// correct caller needs, making the two sides semantically equivalent
// (asserted in TestEquivalenceStrings).
func BenchmarkParseCell(b *testing.B) {
	n := len(cellStrings)
	b.Run("impl=pure", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			c, err := pure.ParseCell(cellStrings[i%n])
			if err != nil {
				b.Fatal(err)
			}
			sinkCell = c
		}
	})
	b.Run("impl=uber", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			c := uber.CellFromString(cellStrings[i%n])
			if !c.IsValid() {
				b.Fatal("invalid cell")
			}
			sinkUCell = c
		}
	})
}

func BenchmarkCellToString(b *testing.B) {
	b.Run("impl=pure", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sinkStr = cellsPure9[i&llMask].String()
		}
	})
	b.Run("impl=uber", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sinkStr = cellsUber9[i&llMask].String()
		}
	})
}

// validityInputs mixes valid cells with corrupted ones so both branches of
// the validity check are exercised.
var (
	validityPure = func() []pure.Cell {
		out := make([]pure.Cell, 0, 2*len(cellsPure9))
		for _, c := range cellsPure9 {
			out = append(out, c, c^0x7) // corrupt the res-9 digit bits
		}
		return out
	}()
	validityUber = toUberCells(validityPure)
)

func BenchmarkIsValidCell(b *testing.B) {
	n := len(validityPure)
	b.Run("impl=pure", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sinkBool = validityPure[i%n].IsValid()
		}
	})
	b.Run("impl=uber", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sinkBool = validityUber[i%n].IsValid()
		}
	})
}

func BenchmarkResolution(b *testing.B) {
	b.Run("impl=pure", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sinkInt = cellsPure9[i&llMask].Resolution()
		}
	})
	b.Run("impl=uber", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sinkInt = cellsUber9[i&llMask].Resolution()
		}
	})
}

// --- Fixed-size geometry: directed edges and vertexes --------------------

var (
	edgesPure = func() []pure.DirectedEdge {
		out := make([]pure.DirectedEdge, len(hexCells9))
		for i, c := range hexCells9 {
			es, err := c.DirectedEdges()
			if err != nil {
				panic(err)
			}
			out[i] = es[0]
		}
		return out
	}()
	edgesUber = func() []uber.DirectedEdge {
		out := make([]uber.DirectedEdge, len(edgesPure))
		for i, e := range edgesPure {
			out[i] = uber.DirectedEdge(uint64(e))
		}
		return out
	}()
	vertsPure = func() []pure.Vertex {
		out := make([]pure.Vertex, len(hexCells9))
		for i, c := range hexCells9 {
			vs, err := c.Vertexes()
			if err != nil {
				panic(err)
			}
			out[i] = vs[0]
		}
		return out
	}()
	vertsUber = func() []uber.Vertex {
		out := make([]uber.Vertex, len(vertsPure))
		for i, v := range vertsPure {
			out[i] = uber.Vertex(uint64(v))
		}
		return out
	}()
)

func BenchmarkDirectedEdges(b *testing.B) {
	n := len(hexCells9)
	b.Run("impl=pure", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			es, err := hexCells9[i%n].DirectedEdges()
			if err != nil {
				b.Fatal(err)
			}
			sinkEdges = es
		}
	})
	b.Run("impl=uber", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			es, err := hexCellsUber9[i%n].DirectedEdges()
			if err != nil {
				b.Fatal(err)
			}
			sinkUEdges = es
		}
	})
}

func BenchmarkDirectedEdgeBoundary(b *testing.B) {
	n := len(edgesPure)
	b.Run("impl=pure", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			bd, err := edgesPure[i%n].Boundary()
			if err != nil {
				b.Fatal(err)
			}
			sinkInt = bd.Len()
		}
	})
	b.Run("impl=uber", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			bd, err := edgesUber[i%n].Boundary()
			if err != nil {
				b.Fatal(err)
			}
			sinkInt = len(bd)
		}
	})
}

func BenchmarkVertexes(b *testing.B) {
	n := len(hexCells9)
	b.Run("impl=pure", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			vs, err := hexCells9[i%n].Vertexes()
			if err != nil {
				b.Fatal(err)
			}
			sinkVerts = vs
		}
	})
	b.Run("impl=uber", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			vs, err := hexCellsUber9[i%n].Vertexes()
			if err != nil {
				b.Fatal(err)
			}
			sinkUVerts = vs
		}
	})
}

func BenchmarkVertexLatLng(b *testing.B) {
	n := len(vertsPure)
	b.Run("impl=pure", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			g, err := vertsPure[i%n].LatLng()
			if err != nil {
				b.Fatal(err)
			}
			sinkLL = g
		}
	})
	b.Run("impl=uber", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			g, err := vertsUber[i%n].LatLng()
			if err != nil {
				b.Fatal(err)
			}
			sinkULL = g
		}
	})
}

// --- Collection operations -----------------------------------------------

func BenchmarkChildren(b *testing.B) {
	for _, depth := range []int{1, 3, 5} {
		res := res4Cell.Resolution() + depth
		name := map[int]string{1: "depth=1", 3: "depth=3", 5: "depth=5"}[depth]

		b.Run(name+"/impl=pure", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				cs, err := res4Cell.Children(res)
				if err != nil {
					b.Fatal(err)
				}
				sinkCells = cs
			}
		})
		b.Run(name+"/impl=pure-cold", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				cs, err := res4Cell.AppendChildren(nil, res)
				if err != nil {
					b.Fatal(err)
				}
				sinkCells = cs
			}
		})
		b.Run(name+"/impl=pure-warm", func(b *testing.B) {
			var buf []pure.Cell
			for i := 0; i < b.N; i++ {
				cs, err := res4Cell.AppendChildren(buf[:0], res)
				if err != nil {
					b.Fatal(err)
				}
				buf = cs
			}
			sinkCells = buf
		})
		b.Run(name+"/impl=uber", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				cs, err := res4CellUber.Children(res)
				if err != nil {
					b.Fatal(err)
				}
				sinkUCells = cs
			}
		})
	}
}

func BenchmarkGridDisk(b *testing.B) {
	n := len(hexCells9)
	for _, k := range []int{1, 5, 20} {
		name := map[int]string{1: "k=1", 5: "k=5", 20: "k=20"}[k]

		b.Run(name+"/impl=pure", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				d, err := hexCells9[i%n].GridDisk(k)
				if err != nil {
					b.Fatal(err)
				}
				sinkCells = d
			}
		})
		b.Run(name+"/impl=pure-cold", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				d, err := hexCells9[i%n].AppendGridDisk(nil, k)
				if err != nil {
					b.Fatal(err)
				}
				sinkCells = d
			}
		})
		b.Run(name+"/impl=pure-warm", func(b *testing.B) {
			var buf []pure.Cell
			for i := 0; i < b.N; i++ {
				d, err := hexCells9[i%n].AppendGridDisk(buf[:0], k)
				if err != nil {
					b.Fatal(err)
				}
				buf = d
			}
			sinkCells = buf
		})
		b.Run(name+"/impl=uber", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				d, err := hexCellsUber9[i%n].GridDisk(k)
				if err != nil {
					b.Fatal(err)
				}
				sinkUCells = d
			}
		})
	}
}

// GridDiskDistances result shapes differ: this library returns a flat
// ([]Cell, []int32) pair; the binding returns [][]Cell rings indexed by
// distance. Both compute the same information (asserted in the equivalence
// tests); the nested shape costs the binding extra allocations.
func BenchmarkGridDiskDistances(b *testing.B) {
	n := len(hexCells9)
	const k = 5
	b.Run("k=5/impl=pure", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			cells, dists, err := hexCells9[i%n].GridDiskDistances(k)
			if err != nil {
				b.Fatal(err)
			}
			sinkCells, sinkI32s = cells, dists
		}
	})
	b.Run("k=5/impl=pure-warm", func(b *testing.B) {
		var buf []pure.Cell
		var dbuf []int32
		for i := 0; i < b.N; i++ {
			cells, dists, err := hexCells9[i%n].AppendGridDiskDistances(buf[:0], dbuf[:0], k)
			if err != nil {
				b.Fatal(err)
			}
			buf, dbuf = cells, dists
		}
		sinkCells, sinkI32s = buf, dbuf
	})
	b.Run("k=5/impl=uber", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			rings, err := hexCellsUber9[i%n].GridDiskDistances(k)
			if err != nil {
				b.Fatal(err)
			}
			sinkURings = rings
		}
	})
}

func BenchmarkGridRing(b *testing.B) {
	n := len(hexCells9)
	const k = 5
	b.Run("k=5/impl=pure", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			r, err := hexCells9[i%n].GridRing(k)
			if err != nil {
				b.Fatal(err)
			}
			sinkCells = r
		}
	})
	b.Run("k=5/impl=pure-warm", func(b *testing.B) {
		var buf []pure.Cell
		for i := 0; i < b.N; i++ {
			r, err := hexCells9[i%n].AppendGridRing(buf[:0], k)
			if err != nil {
				b.Fatal(err)
			}
			buf = r
		}
		sinkCells = buf
	})
	b.Run("k=5/impl=uber", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			r, err := hexCellsUber9[i%n].GridRing(k)
			if err != nil {
				b.Fatal(err)
			}
			sinkUCells = r
		}
	})
}

func BenchmarkGridPath(b *testing.B) {
	n := len(pathPairs)
	b.Run("impl=pure", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			pair := pathPairs[i%n]
			p, err := pair[0].GridPath(pair[1])
			if err != nil {
				b.Fatal(err)
			}
			sinkCells = p
		}
	})
	b.Run("impl=pure-warm", func(b *testing.B) {
		var buf []pure.Cell
		for i := 0; i < b.N; i++ {
			pair := pathPairs[i%n]
			p, err := pair[0].AppendGridPath(buf[:0], pair[1])
			if err != nil {
				b.Fatal(err)
			}
			buf = p
		}
		sinkCells = buf
	})
	b.Run("impl=uber", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			pair := pathPairs[i%n]
			p, err := uber.Cell(uint64(pair[0])).GridPath(uber.Cell(uint64(pair[1])))
			if err != nil {
				b.Fatal(err)
			}
			sinkUCells = p
		}
	})
}

// GridDisksUnsafe result shapes differ: this library returns one flat
// fixed-stride buffer (C layout); the binding returns [][]Cell with pruned
// inner slices.
func BenchmarkGridDisksUnsafe(b *testing.B) {
	origins := hexCells9[:64]
	originsUber := hexCellsUber9[:64]
	const k = 2
	b.Run("origins=64/k=2/impl=pure", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			d, err := pure.GridDisksUnsafe(origins, k)
			if err != nil {
				b.Fatal(err)
			}
			sinkCells = d
		}
	})
	b.Run("origins=64/k=2/impl=uber", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			d, err := uber.GridDisksUnsafe(originsUber, k)
			if err != nil {
				b.Fatal(err)
			}
			sinkUCells = d[len(d)-1]
		}
	})
}

func BenchmarkCompact(b *testing.B) {
	b.Run("set=sf9/impl=pure", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			cs, err := pure.CompactCells(sfCells9)
			if err != nil {
				b.Fatal(err)
			}
			sinkCells = cs
		}
	})
	b.Run("set=sf9/impl=pure-warm", func(b *testing.B) {
		var buf []pure.Cell
		for i := 0; i < b.N; i++ {
			cs, err := pure.AppendCompactCells(buf[:0], sfCells9)
			if err != nil {
				b.Fatal(err)
			}
			buf = cs
		}
		sinkCells = buf
	})
	b.Run("set=sf9/impl=uber", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			cs, err := uber.CompactCells(sfCells9Uber)
			if err != nil {
				b.Fatal(err)
			}
			sinkUCells = cs
		}
	})
}

func BenchmarkUncompact(b *testing.B) {
	in := []pure.Cell{res4Cell}
	inUber := []uber.Cell{res4CellUber}
	b.Run("res=4to9/impl=pure", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			cs, err := pure.UncompactCells(in, benchRes)
			if err != nil {
				b.Fatal(err)
			}
			sinkCells = cs
		}
	})
	b.Run("res=4to9/impl=pure-warm", func(b *testing.B) {
		var buf []pure.Cell
		for i := 0; i < b.N; i++ {
			cs, err := pure.AppendUncompactCells(buf[:0], in, benchRes)
			if err != nil {
				b.Fatal(err)
			}
			buf = cs
		}
		sinkCells = buf
	})
	b.Run("res=4to9/impl=uber", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			cs, err := uber.UncompactCells(inUber, benchRes)
			if err != nil {
				b.Fatal(err)
			}
			sinkUCells = cs
		}
	})
}

func BenchmarkPolygonToCells(b *testing.B) {
	for _, res := range []int{7, 9, 11} {
		name := map[int]string{7: "poly=sf/res=7", 9: "poly=sf/res=9", 11: "poly=sf/res=11"}[res]

		b.Run(name+"/impl=pure", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				cs, err := pure.PolygonToCells(sfPolygonPure, res)
				if err != nil {
					b.Fatal(err)
				}
				sinkCells = cs
			}
		})
		b.Run(name+"/impl=pure-warm", func(b *testing.B) {
			var buf []pure.Cell
			for i := 0; i < b.N; i++ {
				cs, err := pure.AppendPolygonToCells(buf[:0], sfPolygonPure, res)
				if err != nil {
					b.Fatal(err)
				}
				buf = cs
			}
			sinkCells = buf
		})
		b.Run(name+"/impl=uber", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				cs, err := uber.PolygonToCells(sfPolygonUber, res)
				if err != nil {
					b.Fatal(err)
				}
				sinkUCells = cs
			}
		})
	}
}

func BenchmarkCellsToMultiPolygon(b *testing.B) {
	b.Run("n=331/impl=pure", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			mp, err := pure.CellsToMultiPolygon(diskCells331)
			if err != nil {
				b.Fatal(err)
			}
			sinkPolys = mp
		}
	})
	b.Run("n=331/impl=uber", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			mp, err := uber.CellsToMultiPolygon(diskCells331Uber)
			if err != nil {
				b.Fatal(err)
			}
			sinkUPolys = mp
		}
	})
}

// --- Batch workloads ------------------------------------------------------

func BenchmarkBatchLatLngToCell(b *testing.B) {
	b.Run("n=10000/impl=pure", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var acc uint64
			for j := range 10_000 {
				c, err := pure.LatLngToCell(llsPure[j&llMask], benchRes)
				if err != nil {
					b.Fatal(err)
				}
				acc ^= uint64(c)
			}
			sinkU64 = acc
		}
	})
	b.Run("n=10000/impl=uber", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var acc uint64
			for j := range 10_000 {
				c, err := uber.LatLngToCell(llsUber[j&llMask], benchRes)
				if err != nil {
					b.Fatal(err)
				}
				acc ^= uint64(c)
			}
			sinkU64 = acc
		}
	})
}

func BenchmarkServiceWorkload(b *testing.B) {
	b.Run("pts=256/impl=pure", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sinkU64 = serviceWorkloadPureAlloc()
		}
	})
	b.Run("pts=256/impl=pure-warm", func(b *testing.B) {
		buf := make([]pure.Cell, 0, 8)
		for i := 0; i < b.N; i++ {
			var acc uint64
			buf, acc = serviceWorkloadPureWarm(buf)
			sinkU64 = acc
		}
	})
	b.Run("pts=256/impl=uber", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sinkU64 = serviceWorkloadUber()
		}
	})
}
