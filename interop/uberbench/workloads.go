package uberbench

import (
	"fmt"
	"math/bits"

	pure "github.com/dimchansky/h3-go"
	uber "github.com/uber/h3-go/v4"
)

// serviceN is the number of input points in the service-style workload.
const serviceN = 256

// The service workload models a common point-enrichment pattern: index an
// incoming coordinate, look at the immediate neighborhood, and aggregate at
// a coarser resolution. Per point: latLngToCell(res 9) -> gridDisk(k=1) ->
// cellToParent(res 7) for every disk cell. The returned checksum folds in
// every produced index so the equivalence test can pin that both
// implementations do identical work.

func serviceWorkloadPureAlloc() uint64 {
	var acc uint64
	for i := range serviceN {
		c, err := pure.LatLngToCell(llsPure[i], benchRes)
		if err != nil {
			panic(err)
		}
		disk, err := c.GridDisk(1)
		if err != nil {
			panic(err)
		}
		for _, d := range disk {
			p, err := d.Parent(7)
			if err != nil {
				panic(err)
			}
			acc = bits.RotateLeft64(acc, 7) ^ uint64(p)
		}
	}
	return acc
}

// serviceWorkloadPureWarm is the same workload on the buffer-reuse path;
// the caller owns buf across calls (warm reuse). Returns the possibly-grown
// buffer alongside the checksum.
func serviceWorkloadPureWarm(buf []pure.Cell) ([]pure.Cell, uint64) {
	var acc uint64
	for i := range serviceN {
		c, err := pure.LatLngToCell(llsPure[i], benchRes)
		if err != nil {
			panic(err)
		}
		buf, err = c.AppendGridDisk(buf[:0], 1)
		if err != nil {
			panic(err)
		}
		for _, d := range buf {
			p, err := d.Parent(7)
			if err != nil {
				panic(err)
			}
			acc = bits.RotateLeft64(acc, 7) ^ uint64(p)
		}
	}
	return buf, acc
}

func serviceWorkloadUber() uint64 {
	var acc uint64
	for i := range serviceN {
		c, err := uber.LatLngToCell(llsUber[i], benchRes)
		if err != nil {
			panic(err)
		}
		disk, err := c.GridDisk(1)
		if err != nil {
			panic(err)
		}
		for _, d := range disk {
			p, err := d.Parent(7)
			if err != nil {
				panic(err)
			}
			acc = bits.RotateLeft64(acc, 7) ^ uint64(p)
		}
	}
	return acc
}

// MemResult is what a memory workload hands back to the memprobe command:
// a checksum (defeats dead-code elimination) and whatever must stay
// reachable while the process-level measurement is taken.
type MemResult struct {
	Checksum uint64
	Retained any
}

// MemWorkload is one process-level memory measurement scenario, implemented
// once per library over identical inputs.
type MemWorkload struct {
	Name        string
	Description string
	Pure        func(iters int) MemResult
	Uber        func(iters int) MemResult
}

func checksumPure(cells []pure.Cell) uint64 {
	var acc uint64
	for _, c := range cells {
		acc = bits.RotateLeft64(acc, 7) ^ uint64(c)
	}
	return acc
}

func checksumUber(cells []uber.Cell) uint64 {
	var acc uint64
	for _, c := range cells {
		acc = bits.RotateLeft64(acc, 7) ^ uint64(c)
	}
	return acc
}

// res2Cell is the fixed ancestor for the children-deep workload
// (res 2 -> res 8 is 7^6 = 117,649 cells).
var res2Cell = func() pure.Cell {
	p, err := res4Cell.Parent(2)
	if err != nil {
		panic(err)
	}
	return p
}()

// MemWorkloads is the process-level memory comparison matrix run by
// cmd/memprobe. Every scenario exists in both implementations over the
// identical deterministic inputs from fixtures.go.
var MemWorkloads = []MemWorkload{
	{
		Name:        "polyfill-large",
		Description: "PolygonToCells(SF polygon, res 11) — ~61k cells per call; last result retained",
		Pure: func(iters int) MemResult {
			var last []pure.Cell
			var acc uint64
			for range iters {
				cells, err := pure.PolygonToCells(sfPolygonPure, 11)
				if err != nil {
					panic(err)
				}
				acc ^= checksumPure(cells)
				last = cells
			}
			return MemResult{Checksum: acc ^ uint64(len(last)), Retained: last}
		},
		Uber: func(iters int) MemResult {
			var last []uber.Cell
			var acc uint64
			for range iters {
				cells, err := uber.PolygonToCells(sfPolygonUber, 11)
				if err != nil {
					panic(err)
				}
				acc ^= checksumUber(cells)
				last = cells
			}
			return MemResult{Checksum: acc ^ uint64(len(last)), Retained: last}
		},
	},
	{
		Name:        "children-deep",
		Description: "Children(res 2 cell, res 8) — 117,649 cells per call; last result retained",
		Pure: func(iters int) MemResult {
			var last []pure.Cell
			var acc uint64
			for range iters {
				cells, err := res2Cell.Children(8)
				if err != nil {
					panic(err)
				}
				acc ^= checksumPure(cells)
				last = cells
			}
			return MemResult{Checksum: acc, Retained: last}
		},
		Uber: func(iters int) MemResult {
			u := uber.Cell(uint64(res2Cell))
			var last []uber.Cell
			var acc uint64
			for range iters {
				cells, err := u.Children(8)
				if err != nil {
					panic(err)
				}
				acc ^= checksumUber(cells)
				last = cells
			}
			return MemResult{Checksum: acc, Retained: last}
		},
	},
	{
		Name:        "uncompact-res0-to-5",
		Description: "UncompactCells(all 122 res-0 cells, res 5) — 2,050,854 cells per call; last result retained",
		Pure: func(iters int) MemResult {
			base := pure.Res0Cells()
			var last []pure.Cell
			var acc uint64
			for range iters {
				cells, err := pure.UncompactCells(base, 5)
				if err != nil {
					panic(err)
				}
				acc ^= checksumPure(cells)
				last = cells
			}
			return MemResult{Checksum: acc, Retained: last}
		},
		Uber: func(iters int) MemResult {
			base, err := uber.Res0Cells()
			if err != nil {
				panic(err)
			}
			var last []uber.Cell
			var acc uint64
			for range iters {
				cells, err := uber.UncompactCells(base, 5)
				if err != nil {
					panic(err)
				}
				acc ^= checksumUber(cells)
				last = cells
			}
			return MemResult{Checksum: acc, Retained: last}
		},
	},
	{
		Name:        "compact-large",
		Description: "CompactCells over the 2,050,854-cell res-5 set (input built once, outside the loop)",
		Pure: func(iters int) MemResult {
			input, err := pure.UncompactCells(pure.Res0Cells(), 5)
			if err != nil {
				panic(err)
			}
			var last []pure.Cell
			var acc uint64
			for range iters {
				cells, err := pure.CompactCells(input)
				if err != nil {
					panic(err)
				}
				acc ^= checksumPure(cells)
				last = cells
			}
			return MemResult{Checksum: acc, Retained: [2]any{input, last}}
		},
		Uber: func(iters int) MemResult {
			base, err := uber.Res0Cells()
			if err != nil {
				panic(err)
			}
			input, err := uber.UncompactCells(base, 5)
			if err != nil {
				panic(err)
			}
			var last []uber.Cell
			var acc uint64
			for range iters {
				cells, err := uber.CompactCells(input)
				if err != nil {
					panic(err)
				}
				acc ^= checksumUber(cells)
				last = cells
			}
			return MemResult{Checksum: acc, Retained: [2]any{input, last}}
		},
	},
	{
		Name:        "scalar-1m",
		Description: "1,000,000 latLngToCell(res 9) calls cycling the coordinate dataset; nothing retained",
		Pure: func(iters int) MemResult {
			var acc uint64
			for range iters {
				for i := range 1_000_000 {
					c, err := pure.LatLngToCell(llsPure[i%len(llsPure)], benchRes)
					if err != nil {
						panic(err)
					}
					acc = bits.RotateLeft64(acc, 7) ^ uint64(c)
				}
			}
			return MemResult{Checksum: acc}
		},
		Uber: func(iters int) MemResult {
			var acc uint64
			for range iters {
				for i := range 1_000_000 {
					c, err := uber.LatLngToCell(llsUber[i%len(llsUber)], benchRes)
					if err != nil {
						panic(err)
					}
					acc = bits.RotateLeft64(acc, 7) ^ uint64(c)
				}
			}
			return MemResult{Checksum: acc}
		},
	},
	{
		Name:        "multipolygon-sf9",
		Description: "CellsToMultiPolygon over the 1253-cell SF res-9 polyfill; last result retained",
		Pure: func(iters int) MemResult {
			var last []pure.GeoPolygon
			for range iters {
				mp, err := pure.CellsToMultiPolygon(sfCells9)
				if err != nil {
					panic(err)
				}
				last = mp
			}
			return MemResult{Checksum: uint64(len(last)), Retained: last}
		},
		Uber: func(iters int) MemResult {
			var last []uber.GeoPolygon
			for range iters {
				mp, err := uber.CellsToMultiPolygon(sfCells9Uber)
				if err != nil {
					panic(err)
				}
				last = mp
			}
			return MemResult{Checksum: uint64(len(last)), Retained: last}
		},
	},
	{
		Name:        "retained-polyfill",
		Description: "200 x PolygonToCells(SF, res 9), all results retained — steady-state retention shape",
		Pure: func(iters int) MemResult {
			var all [][]pure.Cell
			var acc uint64
			for range iters {
				for range 200 {
					cells, err := pure.PolygonToCells(sfPolygonPure, 9)
					if err != nil {
						panic(err)
					}
					acc ^= checksumPure(cells)
					all = append(all, cells)
				}
			}
			return MemResult{Checksum: acc, Retained: all}
		},
		Uber: func(iters int) MemResult {
			var all [][]uber.Cell
			var acc uint64
			for range iters {
				for range 200 {
					cells, err := uber.PolygonToCells(sfPolygonUber, 9)
					if err != nil {
						panic(err)
					}
					acc ^= checksumUber(cells)
					all = append(all, cells)
				}
			}
			return MemResult{Checksum: acc, Retained: all}
		},
	},
}

// MemWorkloadByName returns the named workload or an error listing the
// valid names.
func MemWorkloadByName(name string) (*MemWorkload, error) {
	for i := range MemWorkloads {
		if MemWorkloads[i].Name == name {
			return &MemWorkloads[i], nil
		}
	}
	names := make([]string, len(MemWorkloads))
	for i, w := range MemWorkloads {
		names[i] = w.Name
	}
	return nil, fmt.Errorf("unknown workload %q (valid: %v)", name, names)
}
