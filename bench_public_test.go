package h3

import "testing"

func BenchmarkLatLngToCell(b *testing.B) {
	ll := LatLngDegs(37.775938728915946, -122.41795063018799)
	b.ReportAllocs()
	for b.Loop() {
		if _, err := LatLngToCell(ll, 9); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCellToLatLng(b *testing.B) {
	c := Cell(0x8928308280fffff)
	b.ReportAllocs()
	for b.Loop() {
		if _, err := c.LatLng(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCellBoundary(b *testing.B) {
	c := Cell(0x8928308280fffff)
	b.ReportAllocs()
	for b.Loop() {
		if _, err := c.Boundary(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkAppendChildren(b *testing.B) {
	parent := Cell(0x8928308280fffff)
	buf := make([]Cell, 0, 7*7*7)
	b.ReportAllocs()
	for b.Loop() {
		out, err := parent.AppendChildren(buf, 12)
		if err != nil {
			b.Fatal(err)
		}
		_ = out
	}
}

func BenchmarkImmediateHierarchy(b *testing.B) {
	cell := Cell(0x8928308280fffff)
	b.Run("ParentComposed", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_, _ = cell.Parent(cell.Resolution() - 1)
		}
	})
	b.Run("ImmediateParent", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_, _ = cell.ImmediateParent()
		}
	})
	b.Run("ChildrenComposed", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_, _ = cell.Children(cell.Resolution() + 1)
		}
	})
	b.Run("ImmediateChildren", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_, _ = cell.ImmediateChildren()
		}
	})
	b.Run("AppendImmediateChildrenWarm", func(b *testing.B) {
		buf := make([]Cell, 0, 7)
		b.ReportAllocs()
		for b.Loop() {
			_, _ = cell.AppendImmediateChildren(buf[:0])
		}
	})
}

func BenchmarkTypedIsValidIndex(b *testing.B) {
	cell := Cell(0x8928308280fffff)
	b.Run("ExplicitUint64", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_ = IsValidIndex(uint64(cell))
		}
	})
	b.Run("Typed", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_ = IsValidIndex(cell)
		}
	})
}

func BenchmarkIndexDigitModes(b *testing.B) {
	cell := Cell(0x8928308280fffff)
	edges, _ := cell.DirectedEdges()
	vertexes, _ := cell.Vertexes()
	b.Run("Cell", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_, _ = cell.IndexDigit(9)
		}
	})
	b.Run("DirectedEdge", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_, _ = edges[0].IndexDigit(9)
		}
	})
	b.Run("Vertex", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_, _ = vertexes[0].IndexDigit(9)
		}
	})
}

func BenchmarkGridDiskDistancesGrouped(b *testing.B) {
	cell := Cell(0x8928308280fffff)
	const k = 5
	b.Run("Grouped", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_, _ = cell.GridDiskDistancesGrouped(k)
		}
	})
	b.Run("ManualAppendPerRing", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			cells, distances, _ := cell.GridDiskDistances(k)
			rings := make([][]Cell, k+1)
			for i, result := range cells {
				rings[distances[i]] = append(rings[distances[i]], result)
			}
			_ = rings
		}
	})
}

func BenchmarkCellString(b *testing.B) {
	c := Cell(0x8928308280fffff)
	b.ReportAllocs()
	for b.Loop() {
		_ = c.String()
	}
}

func BenchmarkParseCell(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		if _, err := ParseCell("8928308280fffff"); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkAppendGridDisk(b *testing.B) {
	c := Cell(0x8928308280fffff)
	buf := make([]Cell, 0, 64)
	b.ReportAllocs()
	for b.Loop() {
		out, err := c.AppendGridDisk(buf, 2)
		if err != nil {
			b.Fatal(err)
		}
		_ = out
	}
}

func BenchmarkAppendGridDiskDistances(b *testing.B) {
	c := Cell(0x8928308280fffff)
	buf := make([]Cell, 0, 64)
	distBuf := make([]int32, 0, 64)
	b.ReportAllocs()
	for b.Loop() {
		cells, dists, err := c.AppendGridDiskDistances(buf, distBuf, 2)
		if err != nil {
			b.Fatal(err)
		}
		_, _ = cells, dists
	}
}

func BenchmarkGridPath(b *testing.B) {
	c := Cell(0x8928308280fffff)
	other, err := LatLngToCell(LatLngDegs(37.8, -122.5), 9)
	if err != nil {
		b.Fatal(err)
	}
	buf := make([]Cell, 0, 128)
	b.ReportAllocs()
	for b.Loop() {
		out, err := c.AppendGridPath(buf, other)
		if err != nil {
			b.Fatal(err)
		}
		_ = out
	}
}

func BenchmarkDirectedEdges(b *testing.B) {
	c := Cell(0x8928308280fffff)
	b.ReportAllocs()
	for b.Loop() {
		if _, err := c.DirectedEdges(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkAppendPolygonToCells(b *testing.B) {
	p := sfTestPolygon()
	buf := make([]Cell, 0, 4096)
	b.ReportAllocs()
	for b.Loop() {
		out, err := AppendPolygonToCells(buf, p, 9)
		if err != nil {
			b.Fatal(err)
		}
		_ = out
	}
}

func BenchmarkAppendCompactCells(b *testing.B) {
	parent := Cell(0x8928308280fffff)
	cells, err := parent.Children(12)
	if err != nil {
		b.Fatal(err)
	}
	buf := make([]Cell, 0, len(cells))
	b.ReportAllocs()
	for b.Loop() {
		out, err := AppendCompactCells(buf, cells)
		if err != nil {
			b.Fatal(err)
		}
		_ = out
	}
}

func BenchmarkCellsToMultiPolygon(b *testing.B) {
	c := Cell(0x8928308280fffff)
	disk, err := c.GridDisk(2)
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	for b.Loop() {
		if _, err := CellsToMultiPolygon(disk); err != nil {
			b.Fatal(err)
		}
	}
}
