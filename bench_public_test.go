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
