package cli

// Grid traversal commands (gridDisk ... localIjToCell) and hierarchy
// commands (cellToParent ... uncompactCells) — two registry groups that
// share the cell/resolution input helpers below.

import (
	"io"
	"os"

	h3 "github.com/dimchansky/h3-go"
)

func gridCommands() []command {
	cell := func() optionSpec { return opt("-c --cell", true) }
	format := func() optionSpec { return opt("-f --format", false) }
	k := func() optionSpec { return opt("-k", true) }
	origin := func() optionSpec { return opt("-o --origin", true) }
	destination := func() optionSpec { return opt("-d --destination", true) }
	return []command{
		{name: "gridDisk", description: "Returns H3 cells within k steps of the origin", options: []optionSpec{cell(), k(), format()}, run: runGridDisk},
		{name: "gridDiskDistances", description: "Returns H3 cells grouped by grid distance", options: []optionSpec{cell(), k(), format()}, run: runGridDiskDistances},
		{name: "gridRing", description: "Returns H3 cells exactly k steps from the origin", options: []optionSpec{cell(), k(), format()}, run: runGridRing},
		{name: "gridPathCells", description: "Returns a path between two H3 cells", options: []optionSpec{origin(), destination(), format()}, run: runGridPath},
		{name: "gridDistance", description: "Returns the grid distance between two H3 cells", options: []optionSpec{origin(), destination()}, run: runGridDistance},
		{name: "cellToLocalIj", description: "Returns local IJ coordinates for a cell", options: []optionSpec{cell(), origin(), format()}, run: runCellToLocalIJ},
		{name: "localIjToCell", description: "Returns a cell for local IJ coordinates", options: []optionSpec{origin(), opt("-i", true), opt("-j", true), format()}, run: runLocalIJToCell},
	}
}

func runGridDisk(env environment, p parsedArgs) error {
	k, err := p.integer("-k")
	if err != nil {
		return h3.ErrFailed
	}
	cells, err := h3.Cell(rawHex(p.get("-c"))).GridDisk(k)
	if err != nil {
		return err
	}
	return writeCells(env.out, cells, formatValue(p))
}

// runGridDiskDistances groups the disk by distance ring, as upstream does:
// json nests one array per k, newline separates rings with a blank line.
func runGridDiskDistances(env environment, p parsedArgs) error {
	k, err := p.integer("-k")
	if err != nil {
		return h3.ErrFailed
	}
	cells, distances, err := h3.Cell(rawHex(p.get("-c"))).GridDiskDistances(k)
	if err != nil {
		return err
	}
	format := formatValue(p)
	if format == "json" {
		writeText(env.out, "[")
		for distance := 0; distance <= k; distance++ {
			if distance > 0 {
				writeText(env.out, ", ")
			}
			writeText(env.out, "[")
			written := 0
			for i, cell := range cells {
				if int(distances[i]) != distance {
					continue
				}
				if written > 0 {
					writeText(env.out, ", ")
				}
				writef(env.out, "\"%s\"", cell)
				written++
			}
			writeText(env.out, "]")
		}
		writeln(env.out, "]")
		return nil
	}
	if format == "newline" {
		for distance := 0; distance <= k; distance++ {
			for i, cell := range cells {
				if int(distances[i]) == distance {
					writeln(env.out, cell)
				}
			}
			if distance < k {
				writeln(env.out)
			}
		}
		return nil
	}
	return h3.ErrFailed
}

func runGridRing(env environment, p parsedArgs) error {
	k, err := p.integer("-k")
	if err != nil {
		return h3.ErrFailed
	}
	cells, err := h3.Cell(rawHex(p.get("-c"))).GridRing(k)
	if err != nil {
		return err
	}
	return writeCells(env.out, cells, formatValue(p))
}

func runGridPath(env environment, p parsedArgs) error {
	cells, err := h3.Cell(rawHex(p.get("-o"))).GridPath(h3.Cell(rawHex(p.get("-d"))))
	if err != nil {
		return err
	}
	return writeCells(env.out, cells, formatValue(p))
}

func runGridDistance(env environment, p parsedArgs) error {
	distance, err := h3.Cell(rawHex(p.get("-o"))).GridDistance(h3.Cell(rawHex(p.get("-d"))))
	if err != nil {
		return err
	}
	// Hex, not decimal: upstream prints the distance with PRIx64.
	writef(env.out, "%x\n", distance)
	return nil
}

func runCellToLocalIJ(env environment, p parsedArgs) error {
	ij, err := h3.CellToLocalIJ(h3.Cell(rawHex(p.get("-o"))), h3.Cell(rawHex(p.get("-c"))))
	if err != nil {
		return err
	}
	switch formatValue(p) {
	case "json":
		writef(env.out, "[%d, %d]\n", ij.I, ij.J)
	case "newline":
		writef(env.out, "%d\n%d\n", ij.I, ij.J)
	default:
		return h3.ErrFailed
	}
	return nil
}

func runLocalIJToCell(env environment, p parsedArgs) error {
	i, err := p.integer("-i")
	if err != nil {
		return h3.ErrFailed
	}
	j, err := p.integer("-j")
	if err != nil {
		return h3.ErrFailed
	}
	cell, err := h3.LocalIJToCell(h3.Cell(rawHex(p.get("-o"))), h3.CoordIJ{I: int32(i), J: int32(j)})
	if err != nil {
		return err
	}
	return writeCell(env.out, cell, formatValue(p))
}

func hierarchyCommands() []command {
	cell := func() optionSpec { return opt("-c --cell", true) }
	res := func() optionSpec { return opt("-r --resolution", true) }
	format := func() optionSpec { return opt("-f --format", false) }
	file := func() optionSpec { return opt("-i --file", false) }
	cells := func() optionSpec { return opt("-c --cells", false) }
	return []command{
		{name: "cellToParent", description: "Returns the parent cell", options: []optionSpec{cell(), res(), format()}, run: runCellToParent},
		{name: "cellToChildren", description: "Returns child cells", options: []optionSpec{cell(), res(), format()}, run: runCellToChildren},
		{name: "cellToChildrenSize", description: "Returns the number of child cells", options: []optionSpec{cell(), res()}, run: runCellToChildrenSize},
		{name: "cellToCenterChild", description: "Returns the center child cell", options: []optionSpec{cell(), res(), format()}, run: runCellToCenterChild},
		{name: "cellToChildPos", description: "Returns the position of a child", options: []optionSpec{cell(), res()}, run: runCellToChildPos},
		{name: "childPosToCell", description: "Returns a child at a position", options: []optionSpec{cell(), res(), opt("-p --position", true), format()}, run: runChildPosToCell},
		{name: "compactCells", description: "Compacts the provided set of cells", options: []optionSpec{file(), cells(), format()}, run: runCompactCells},
		{name: "uncompactCells", description: "Uncompacts the provided set of cells", options: []optionSpec{file(), cells(), res(), format()}, run: runUncompactCells},
	}
}

// hierarchyInput decodes the -c cell / -r resolution pair shared by most
// hierarchy commands; a malformed resolution is the returned error.
func hierarchyInput(p parsedArgs) (h3.Cell, int, error) {
	res, err := p.integer("-r")
	return h3.Cell(rawHex(p.get("-c"))), res, err
}

func runCellToParent(env environment, p parsedArgs) error {
	cell, res, err := hierarchyInput(p)
	if err != nil {
		return h3.ErrFailed
	}
	out, err := cell.Parent(res)
	if err != nil {
		return err
	}
	return writeCell(env.out, out, formatValue(p))
}

func runCellToChildren(env environment, p parsedArgs) error {
	cell, res, err := hierarchyInput(p)
	if err != nil {
		return h3.ErrFailed
	}
	out, err := cell.Children(res)
	if err != nil {
		return err
	}
	return writeCells(env.out, out, formatValue(p))
}

func runCellToChildrenSize(env environment, p parsedArgs) error {
	cell, res, err := hierarchyInput(p)
	if err != nil {
		return h3.ErrFailed
	}
	n, err := cell.NumChildren(res)
	if err != nil {
		return err
	}
	writeln(env.out, n)
	return nil
}

func runCellToCenterChild(env environment, p parsedArgs) error {
	cell, res, err := hierarchyInput(p)
	if err != nil {
		return h3.ErrFailed
	}
	out, err := cell.CenterChild(res)
	if err != nil {
		return err
	}
	return writeCell(env.out, out, formatValue(p))
}

func runCellToChildPos(env environment, p parsedArgs) error {
	cell, res, err := hierarchyInput(p)
	if err != nil {
		return h3.ErrFailed
	}
	pos, err := cell.ChildPos(res)
	if err != nil {
		return err
	}
	writeln(env.out, pos)
	return nil
}

func runChildPosToCell(env environment, p parsedArgs) error {
	cell, res, err := hierarchyInput(p)
	if err != nil {
		return h3.ErrFailed
	}
	pos, err := p.int64("-p")
	if err != nil {
		return h3.ErrFailed
	}
	out, err := cell.ChildAtPos(pos, res)
	if err != nil {
		return err
	}
	return writeCell(env.out, out, formatValue(p))
}

// readCellCommandInput reads a cell set for compactCells/uncompactCells
// from -i (file, or stdin via "--") or inline -c. Unlike readSource, giving
// both options is not an error here — -i wins — matching the upstream
// implementations of these two commands. Any read failure reports the
// upstream "does not exist" wording; that is the only diagnostic C emits.
func readCellCommandInput(env environment, p parsedArgs) ([]h3.Cell, error) {
	if !p.has("-i") && !p.has("-c") {
		return nil, failDirect(env.errOut, "You must provide either a file to read from or a set of cells")
	}
	var data []byte
	var err error
	if p.has("-i") {
		if p.get("-i") == "--" {
			data, err = io.ReadAll(env.in)
		} else {
			data, err = os.ReadFile(p.get("-i"))
		}
	} else {
		data = []byte(p.get("-c"))
	}
	if err != nil {
		return nil, failDirect(env.errOut, "The specified file does not exist.")
	}
	if len(data) == 0 {
		return nil, failDirect(env.errOut, "The specified file is empty.")
	}
	return parseCells(data), nil
}

func runCompactCells(env environment, p parsedArgs) error {
	cells, err := readCellCommandInput(env, p)
	if err != nil {
		return err
	}
	out, err := h3.CompactCells(cells)
	if err != nil {
		return err
	}
	return writeCells(env.out, out, formatValue(p))
}

func runUncompactCells(env environment, p parsedArgs) error {
	cells, err := readCellCommandInput(env, p)
	if err != nil {
		return err
	}
	res, err := p.integer("-r")
	if err != nil {
		return h3.ErrFailed
	}
	out, err := h3.UncompactCells(cells, res)
	if err != nil {
		return err
	}
	return writeCells(env.out, out, formatValue(p))
}
