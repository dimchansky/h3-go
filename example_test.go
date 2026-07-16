package h3_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	h3 "github.com/dimchansky/h3-go"
)

func ExampleParseCell() {
	// Upper/lowercase hex digits, an optional "0x"/"0X" prefix, and leading
	// zeros are all accepted.
	a, _ := h3.ParseCell("0x8928308280FFFFF")
	b, _ := h3.ParseCell("8928308280fffff")
	fmt.Println(a == b)

	// A well-formed string that is not a valid cell fails with the
	// ErrCellInvalid sentinel...
	_, err := h3.ParseCell("ffffffffffffffff")
	fmt.Println(errors.Is(err, h3.ErrCellInvalid))

	// ...while malformed text fails with a wrapped strconv error instead of
	// an Err* sentinel.
	_, err = h3.ParseCell("not-hex")
	fmt.Println(errors.Is(err, strconv.ErrSyntax), errors.Is(err, h3.ErrCellInvalid))
	// Output:
	// true
	// true
	// true false
}

func ExampleCell_MarshalText() {
	type record struct {
		Cell h3.Cell `json:"cell"`
	}
	in := record{Cell: h3.Cell(0x8928308280fffff)}
	data, err := json.Marshal(in)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))

	var out record
	if err := json.Unmarshal(data, &out); err != nil {
		panic(err)
	}
	fmt.Println(out.Cell == in.Cell)

	// MarshalText never validates — even the zero Cell marshals (as "0") —
	// while UnmarshalText is the validating direction and rejects it.
	zero, _ := json.Marshal(record{})
	fmt.Println(string(zero))
	err = json.Unmarshal(zero, &out)
	fmt.Println(errors.Is(err, h3.ErrCellInvalid))
	// Output:
	// {"cell":"8928308280fffff"}
	// true
	// {"cell":"0"}
	// true
}

func ExampleCell_Boundary() {
	cell, _ := h3.ParseCell("8928308280fffff")
	boundary, err := cell.Boundary()
	if err != nil {
		panic(err)
	}
	// Iterate with Len/At: a boundary holds the cell's topological vertices
	// (6 for a hexagon, 5 for a pentagon) plus any distortion vertices where
	// it crosses icosahedron faces, so the count varies by cell.
	for i := 0; i < boundary.Len(); i++ {
		v := boundary.At(i)
		fmt.Printf("%.4f, %.4f\n", v.Lat.Deg(), v.Lng.Deg())
	}
	// Output:
	// 37.7752, -122.4172
	// 37.7769, -122.4161
	// 37.7784, -122.4174
	// 37.7782, -122.4197
	// 37.7765, -122.4208
	// 37.7750, -122.4195
}

func ExampleLatLngToCell() {
	cell, err := h3.LatLngToCell(h3.LatLngDegs(37.775938728915946, -122.41795063018799), 9)
	if err != nil {
		panic(err)
	}
	fmt.Println(cell)
	// Output: 8928308280fffff
}

func ExampleCell_LatLng() {
	cell, _ := h3.ParseCell("8928308280fffff")
	center, err := cell.LatLng()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%.4f, %.4f\n", center.Lat.Deg(), center.Lng.Deg())
	// Output: 37.7767, -122.4185
}

func ExampleCell_GridDisk() {
	cell, _ := h3.ParseCell("8928308280fffff")
	disk, err := cell.GridDisk(1)
	if err != nil {
		panic(err)
	}
	fmt.Println(len(disk), "cells")
	// Output: 7 cells
}

func ExampleCell_AppendGridDisk() {
	cell, _ := h3.ParseCell("8928308280fffff")
	// Reuse one buffer across many queries: zero allocations once warm.
	buf := make([]h3.Cell, 0, 64)
	for k := 1; k <= 3; k++ {
		disk, err := cell.AppendGridDisk(buf[:0], k)
		if err != nil {
			panic(err)
		}
		fmt.Println(len(disk))
	}

	// On error the destination comes back with its original length and
	// elements — no partial results are observable.
	disk, _ := cell.AppendGridDisk(buf[:0], 1)
	rolled, err := cell.AppendGridDisk(disk, -1)
	fmt.Println(errors.Is(err, h3.ErrDomain), len(rolled) == len(disk))
	// Output:
	// 7
	// 19
	// 37
	// true true
}

func ExampleCell_GridDistance() {
	a, _ := h3.ParseCell("8928308280fffff")
	b, _ := h3.ParseCell("8928308280bffff")
	d, err := a.GridDistance(b)
	if err != nil {
		panic(err)
	}
	fmt.Println("grid moves:", d)

	// Cells of different resolutions fail with ErrResolutionMismatch.
	coarser, _ := a.Parent(5)
	_, err = a.GridDistance(coarser)
	fmt.Println(errors.Is(err, h3.ErrResolutionMismatch))
	// Output:
	// grid moves: 1
	// true
}

func ExampleCellToLocalIJ() {
	origin, _ := h3.ParseCell("8928308280fffff")
	cell, _ := h3.ParseCell("8928308280bffff")

	// IJ coordinates are only meaningful relative to their origin, and the
	// coordinate space may change between H3 versions — don't persist them.
	ij, err := h3.CellToLocalIJ(origin, cell)
	if err != nil {
		panic(err)
	}

	// Within one H3 version and origin, LocalIJToCell inverts the mapping.
	back, err := h3.LocalIJToCell(origin, ij)
	if err != nil {
		panic(err)
	}
	fmt.Println(back == cell)
	// Output: true
}

func ExampleCell_Parent() {
	cell, _ := h3.ParseCell("8928308280fffff")
	parent, err := cell.Parent(5)
	if err != nil {
		panic(err)
	}
	fmt.Println(parent)
	// Output: 85283083fffffff
}

func ExampleCell_ChildrenSeq() {
	parent, _ := h3.ParseCell("85283473fffffff")
	n := 0
	for range parent.ChildrenSeq(7) {
		n++
	}
	fmt.Println(n, "children")
	// Output: 49 children
}

func ExampleCell_ImmediateChildren() {
	parent, _ := h3.ParseCell("8928308280fffff")
	children, err := parent.ImmediateChildren()
	if err != nil {
		panic(err)
	}
	fmt.Println(len(children), "immediate children")
	// Output: 7 immediate children
}

func ExampleCell_ChildPos() {
	cell, _ := h3.ParseCell("8928308280fffff") // resolution 9
	// A cell's position within the canonical child order of its ancestor
	// round-trips through ChildAtPos.
	pos, err := cell.ChildPos(5)
	if err != nil {
		panic(err)
	}
	parent, _ := cell.Parent(5)
	back, err := parent.ChildAtPos(pos, 9)
	if err != nil {
		panic(err)
	}
	fmt.Println(pos, back == cell)
	// Output: 1718 true
}

func ExampleCompactCells() {
	parent, _ := h3.ParseCell("85283473fffffff")
	// Valid input: cells of one resolution, no duplicates.
	children, _ := parent.Children(6)
	compacted, err := h3.CompactCells(children)
	if err != nil {
		panic(err)
	}
	fmt.Println(len(children), "->", len(compacted), compacted[0] == parent)

	// UncompactCells reverses the compaction at the chosen resolution.
	uncompacted, err := h3.UncompactCells(compacted, 6)
	if err != nil {
		panic(err)
	}
	fmt.Println(len(uncompacted))
	// Output:
	// 7 -> 1 true
	// 7
}

func ExampleCell_GridDiskDistancesGrouped() {
	cell, _ := h3.ParseCell("8928308280fffff")
	rings, err := cell.GridDiskDistancesGrouped(2)
	if err != nil {
		panic(err)
	}
	for distance, ring := range rings {
		fmt.Printf("distance %d: %d cells\n", distance, len(ring))
	}
	// Output:
	// distance 0: 1 cells
	// distance 1: 6 cells
	// distance 2: 12 cells
}

func ExamplePolygonToCells() {
	// Loops are implicitly closed (the first vertex is not repeated) and
	// use radians-backed LatLng values; holes are given structurally.
	outer := h3.GeoLoop{
		h3.LatLngDegs(37.813, -122.408),
		h3.LatLngDegs(37.782, -122.386),
		h3.LatLngDegs(37.708, -122.390),
		h3.LatLngDegs(37.708, -122.507),
		h3.LatLngDegs(37.784, -122.511),
	}
	hole := h3.GeoLoop{
		h3.LatLngDegs(37.790, -122.490),
		h3.LatLngDegs(37.790, -122.410),
		h3.LatLngDegs(37.720, -122.410),
		h3.LatLngDegs(37.720, -122.490),
	}
	cells, err := h3.PolygonToCells(h3.GeoPolygon{GeoLoop: outer}, 7)
	if err != nil {
		panic(err)
	}
	holed, err := h3.PolygonToCells(h3.GeoPolygon{GeoLoop: outer, Holes: []h3.GeoLoop{hole}}, 7)
	if err != nil {
		panic(err)
	}
	fmt.Println(len(cells), "cells;", len(holed), "outside the hole")
	// Output: 20 cells; 10 outside the hole
}

func ExamplePolygonToCells_antimeridian() {
	// Loops crossing the antimeridian are handled automatically; longitudes
	// need not be pre-unwrapped.
	polygon := h3.GeoPolygon{GeoLoop: h3.GeoLoop{
		h3.LatLngDegs(0.5, 179.5),
		h3.LatLngDegs(0.5, -179.5),
		h3.LatLngDegs(-0.5, -179.5),
		h3.LatLngDegs(-0.5, 179.5),
	}}
	cells, err := h3.PolygonToCells(polygon, 4)
	if err != nil {
		panic(err)
	}
	fmt.Println(len(cells), "cells")
	// Output: 8 cells
}

func ExampleCellsToMultiPolygon() {
	// The hollow ring around a cell outlines as one polygon with one hole.
	center, _ := h3.ParseCell("8928308280fffff")
	ring, err := center.GridRing(1)
	if err != nil {
		panic(err)
	}
	polys, err := h3.CellsToMultiPolygon(ring)
	if err != nil {
		panic(err)
	}
	p := polys[0]
	fmt.Println(len(polys), "polygon;", len(p.GeoLoop), "outer vertices;", len(p.Holes), "hole")

	// The structure is GeoJSON-like but not GeoJSON: to export, convert
	// radians to degrees, swap to [lng, lat] order, and close each ring by
	// repeating its first position.
	ringCoords := make([][2]float64, 0, len(p.GeoLoop)+1)
	for _, v := range p.GeoLoop {
		ringCoords = append(ringCoords, [2]float64{v.Lng.Deg(), v.Lat.Deg()})
	}
	ringCoords = append(ringCoords, ringCoords[0]) // explicit GeoJSON closure
	fmt.Println("closed:", len(ringCoords), "positions; first == last:", ringCoords[0] == ringCoords[len(ringCoords)-1])
	// Output:
	// 1 polygon; 18 outer vertices; 1 hole
	// closed: 19 positions; first == last: true
}

func ExamplePolygonToCellsExperimentalSeq() {
	polygon := h3.GeoPolygon{GeoLoop: h3.GeoLoop{
		h3.LatLngDegs(37.813, -122.408),
		h3.LatLngDegs(37.782, -122.386),
		h3.LatLngDegs(37.708, -122.390),
		h3.LatLngDegs(37.708, -122.507),
		h3.LatLngDegs(37.784, -122.511),
	}}
	// Validation happens up front; iteration itself cannot fail.
	seq, err := h3.PolygonToCellsExperimentalSeq(polygon, 7, h3.ContainmentCenter)
	if err != nil {
		panic(err)
	}
	total := 0
	for range seq {
		total++
	}
	// The sequence is re-runnable, and breaking early is safe.
	first := 0
	for range seq {
		first++
		if first == 5 {
			break
		}
	}
	fmt.Println(total, "cells; stopped early at", first)
	// Output: 20 cells; stopped early at 5
}

func ExampleCell_DirectedEdges() {
	cell, _ := h3.ParseCell("8928308280fffff")
	edges, err := cell.DirectedEdges()
	if err != nil {
		panic(err)
	}
	dest, err := edges[0].Destination()
	if err != nil {
		panic(err)
	}
	ok, _ := cell.IsNeighbor(dest)
	fmt.Println(len(edges), "edges; first leads to a neighbor:", ok)
	// Output: 6 edges; first leads to a neighbor: true
}

func ExampleCell_Vertex() {
	// A topological corner is shared by three cells; the canonical owner
	// makes the Vertex value identical no matter which cell derives it.
	cell, _ := h3.ParseCell("8928308280fffff")
	v0, err := cell.Vertex(0)
	if err != nil {
		panic(err)
	}
	shared := 0
	ring, _ := cell.GridRing(1)
	for _, neighbor := range ring {
		vs, err := neighbor.Vertexes()
		if err != nil {
			panic(err)
		}
		for _, v := range vs {
			if v == v0 {
				shared++
			}
		}
	}
	fmt.Println("neighbors deriving the same Vertex value:", shared)
	// Output: neighbors deriving the same Vertex value: 2
}

func ExampleDirectedEdge() {
	origin, _ := h3.ParseCell("8928308280fffff")
	edges, err := origin.DirectedEdges()
	if err != nil {
		panic(err)
	}
	edge := edges[0]

	// Cells returns origin and destination, in that order.
	o, d, err := edge.Cells()
	if err != nil {
		panic(err)
	}
	fmt.Println("origin matches:", o == origin)

	// The reverse edge runs from the destination back to the origin.
	back, err := d.DirectedEdgeTo(o)
	if err != nil {
		panic(err)
	}
	ro, _ := back.Origin()
	fmt.Println("reverse origin is destination:", ro == d)

	// An edge boundary holds the two topological endpoints, plus one
	// distortion vertex when the edge crosses an icosahedron face.
	b, err := edge.Boundary()
	if err != nil {
		panic(err)
	}
	fmt.Println("boundary vertices:", b.Len())
	// Output:
	// origin matches: true
	// reverse origin is destination: true
	// boundary vertices: 2
}

func ExampleCell_AreaKm2() {
	cell, _ := h3.ParseCell("8928308280fffff")
	// The exact spherical area of this specific cell...
	exact, err := cell.AreaKm2()
	if err != nil {
		panic(err)
	}
	// ...versus the average hexagon area at its resolution.
	avg, err := h3.HexagonAreaAvgKm2(9)
	if err != nil {
		panic(err)
	}
	fmt.Printf("exact %.6f km², average %.6f km²\n", exact, avg)
	// Output: exact 0.109398 km², average 0.105333 km²
}

func ExampleCell_AppendVertexes() {
	cell, _ := h3.ParseCell("8928308280fffff")
	// Reuse one buffer across many cells: capacity 6 always suffices, so
	// the loop is allocation-free once warm.
	buf := make([]h3.Vertex, 0, 6)
	cells, _ := cell.GridDisk(1)
	for _, c := range cells {
		verts, err := c.AppendVertexes(buf[:0])
		if err != nil {
			panic(err)
		}
		buf = verts[:0]
		fmt.Print(len(verts), " ")
	}
	// Output: 6 6 6 6 6 6 6
}
