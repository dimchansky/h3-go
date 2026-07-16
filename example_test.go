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
	polygon := h3.GeoPolygon{GeoLoop: h3.GeoLoop{
		h3.LatLngDegs(37.813, -122.408),
		h3.LatLngDegs(37.782, -122.386),
		h3.LatLngDegs(37.708, -122.390),
		h3.LatLngDegs(37.708, -122.507),
		h3.LatLngDegs(37.784, -122.511),
	}}
	cells, err := h3.PolygonToCells(polygon, 7)
	if err != nil {
		panic(err)
	}
	fmt.Println(len(cells), "cells")
	// Output: 20 cells
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
