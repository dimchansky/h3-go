package h3_test

import (
	"fmt"

	h3 "github.com/dimchansky/h3-go"
)

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
	// Output:
	// 7
	// 19
	// 37
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
