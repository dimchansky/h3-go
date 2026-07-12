package cli

import (
	"io"

	h3 "github.com/dimchansky/h3-go"
)

func regionCommands() []command {
	file := func() optionSpec { return opt("-i --file", false) }
	format := func() optionSpec { return opt("-f --format", false) }
	res := func() optionSpec { return opt("-r --resolution", true) }
	return []command{
		{name: "polygonToCells", description: "Converts a polygon into covering cells", options: []optionSpec{file(), opt("-p --polygon", false), res(), format()}, run: runPolygonToCells},
		{name: "maxPolygonToCellsSize", description: "Returns the maximum polygon cover size", options: []optionSpec{file(), opt("-p --polygon", false), res()}, run: runMaxPolygonToCellsSize},
		{name: "cellsToMultiPolygon", description: "Returns a polygon for a set of cells", options: []optionSpec{file(), opt("-c --cells", false), format()}, run: runCellsToMultiPolygon},
	}
}

func polygonInput(env environment, p parsedArgs) (h3.GeoPolygon, error) {
	data, err := readSource(env, p, "-i", "-p")
	if err != nil {
		return h3.GeoPolygon{}, err
	}
	return parsePolygon(data)
}

func runPolygonToCells(env environment, p parsedArgs) error {
	polygon, err := polygonInput(env, p)
	if err != nil {
		return err
	}
	res, err := p.integer("-r")
	if err != nil {
		return h3.ErrFailed
	}
	cells, err := h3.PolygonToCells(polygon, res)
	if err != nil {
		return err
	}
	return writeCells(env.out, cells, formatValue(p))
}

func runMaxPolygonToCellsSize(env environment, p parsedArgs) error {
	polygon, err := polygonInput(env, p)
	if err != nil {
		return err
	}
	res, err := p.integer("-r")
	if err != nil {
		return h3.ErrFailed
	}
	size, err := h3.MaxPolygonToCellsSize(polygon, res)
	if err != nil {
		return err
	}
	writeln(env.out, size)
	return nil
}

func runCellsToMultiPolygon(env environment, p parsedArgs) error {
	data, err := readSource(env, p, "-i", "-c")
	if err != nil {
		return err
	}
	polygons, err := h3.CellsToMultiPolygon(parseCells(data))
	if err != nil {
		return err
	}
	switch formatValue(p) {
	case "json":
		writeMultiPolygonJSON(env.out, polygons)
	case "wkt":
		writeMultiPolygonWKT(env.out, polygons)
	default:
		return h3.ErrFailed
	}
	return nil
}

func writeMultiPolygonJSON(w io.Writer, polygons []h3.GeoPolygon) {
	writeText(w, "[")
	for pi, polygon := range polygons {
		if pi > 0 {
			writeText(w, ", ")
		}
		writeText(w, "[")
		loops := append([]h3.GeoLoop{polygon.GeoLoop}, polygon.Holes...)
		for li, loop := range loops {
			if li > 0 {
				writeText(w, ", ")
			}
			writeText(w, "[")
			for vi, ll := range loop {
				if vi > 0 {
					writeText(w, ", ")
				}
				writef(w, "[%.6f, %.6f]", ll.Lat.Deg(), ll.Lng.Deg())
			}
			writeText(w, "]")
		}
		writeText(w, "]")
	}
	writeln(w, "]")
}

func writeMultiPolygonWKT(w io.Writer, polygons []h3.GeoPolygon) {
	writeText(w, "MULTIPOLYGON (")
	for pi, polygon := range polygons {
		if pi > 0 {
			writeText(w, ", ")
		}
		writeText(w, "(")
		loops := append([]h3.GeoLoop{polygon.GeoLoop}, polygon.Holes...)
		for li, loop := range loops {
			if li > 0 {
				writeText(w, ", ")
			}
			writeText(w, "(")
			for vi, ll := range loop {
				if vi > 0 {
					writeText(w, ", ")
				}
				writef(w, "%.6f %.6f", ll.Lng.Deg(), ll.Lat.Deg())
			}
			writeText(w, ")")
		}
		writeText(w, ")")
	}
	writeln(w, ")")
}

func edgeCommands() []command {
	cell := func() optionSpec { return opt("-c --cell", true) }
	format := func() optionSpec { return opt("-f --format", false) }
	origin := func() optionSpec { return opt("-o --origin", true) }
	dest := func() optionSpec { return opt("-d --destination", true) }
	return []command{
		{name: "areNeighborCells", description: "Determines whether two cells are neighbors", options: []optionSpec{origin(), dest(), format()}, run: runAreNeighborCells},
		{name: "cellsToDirectedEdge", description: "Returns a directed edge between neighboring cells", options: []optionSpec{origin(), dest(), format()}, run: runCellsToDirectedEdge},
		{name: "isValidDirectedEdge", description: "Checks whether an index is a directed edge", options: []optionSpec{cell(), format()}, run: runIsValidDirectedEdge},
		{name: "getDirectedEdgeOrigin", description: "Returns a directed edge origin", options: []optionSpec{cell(), format()}, run: runGetDirectedEdgeOrigin},
		{name: "getDirectedEdgeDestination", description: "Returns a directed edge destination", options: []optionSpec{cell(), format()}, run: runGetDirectedEdgeDestination},
		{name: "directedEdgeToCells", description: "Returns the origin and destination cells", options: []optionSpec{cell(), format()}, run: runDirectedEdgeToCells},
		{name: "originToDirectedEdges", description: "Returns all directed edges from a cell", options: []optionSpec{cell(), format()}, run: runOriginToDirectedEdges},
		{name: "directedEdgeToBoundary", description: "Returns a directed edge boundary", options: []optionSpec{cell(), format()}, run: runDirectedEdgeToBoundary},
	}
}

func edgeFrom(p parsedArgs) h3.DirectedEdge { return h3.DirectedEdge(rawHex(p.get("-c"))) }

func runAreNeighborCells(env environment, p parsedArgs) error {
	ok, err := h3.Cell(rawHex(p.get("-o"))).IsNeighbor(h3.Cell(rawHex(p.get("-d"))))
	if err != nil {
		return err
	}
	return writeBool(env.out, ok, formatValue(p))
}

func runCellsToDirectedEdge(env environment, p parsedArgs) error {
	edge, err := h3.Cell(rawHex(p.get("-o"))).DirectedEdgeTo(h3.Cell(rawHex(p.get("-d"))))
	if err != nil {
		return err
	}
	return writeCell(env.out, h3.Cell(edge), formatValue(p))
}

func runIsValidDirectedEdge(env environment, p parsedArgs) error {
	return writeBool(env.out, edgeFrom(p).IsValid(), formatValue(p))
}

func runGetDirectedEdgeOrigin(env environment, p parsedArgs) error {
	cell, err := edgeFrom(p).Origin()
	if err != nil {
		return err
	}
	return writeCell(env.out, cell, formatValue(p))
}

func runGetDirectedEdgeDestination(env environment, p parsedArgs) error {
	cell, err := edgeFrom(p).Destination()
	if err != nil {
		return err
	}
	return writeCell(env.out, cell, formatValue(p))
}

func runDirectedEdgeToCells(env environment, p parsedArgs) error {
	origin, destination, err := edgeFrom(p).Cells()
	if err != nil {
		return err
	}
	if formatValue(p) == "json" {
		writef(env.out, "[\"%s\", \"%s\"]\n", origin, destination)
		return nil
	}
	if formatValue(p) == "newline" {
		writeln(env.out, origin)
		writeln(env.out, destination)
		return nil
	}
	return h3.ErrFailed
}

func runOriginToDirectedEdges(env environment, p parsedArgs) error {
	cell := h3.Cell(rawHex(p.get("-c")))
	if !cell.IsValid() {
		return h3.ErrCellInvalid
	}
	edges, err := cell.DirectedEdges()
	if err != nil {
		return err
	}
	if formatValue(p) == "json" {
		writeText(env.out, "[")
		for i, edge := range edges {
			if i > 0 {
				writeText(env.out, ", ")
			}
			writef(env.out, "\"%s\"", edge)
		}
		writeln(env.out, "]")
		return nil
	}
	if formatValue(p) == "newline" {
		for _, edge := range edges {
			writeln(env.out, edge)
		}
		return nil
	}
	return h3.ErrFailed
}

func runDirectedEdgeToBoundary(env environment, p parsedArgs) error {
	boundary, err := edgeFrom(p).Boundary()
	if err != nil {
		return err
	}
	return writeBoundary(env.out, boundary, formatValue(p), true)
}
