package cli

// Vertex commands (cellToVertex ... isValidVertex) and the miscellaneous
// group: unit conversions, area/length metrics, global enumerations,
// great-circle distances, and describeH3Error. The metric*Command and
// greatCircleCommand factories fold the many single-value commands that
// differ only in the wrapped library call.

import h3 "github.com/dimchansky/h3-go"

func vertexCommands() []command {
	cell := func() optionSpec { return opt("-c --cell", true) }
	format := func() optionSpec { return opt("-f --format", false) }
	return []command{
		{name: "cellToVertex", description: "Returns a vertex for a cell", options: []optionSpec{cell(), opt("-v --vertex", true), format()}, run: runCellToVertex},
		{name: "cellToVertexes", description: "Returns all vertexes for a cell", options: []optionSpec{cell(), format()}, run: runCellToVertexes},
		{name: "vertexToLatLng", description: "Returns the lat, lng pair for a vertex", options: []optionSpec{cell(), format()}, run: runVertexToLatLng},
		{name: "isValidVertex", description: "Checks whether an index is a vertex", options: []optionSpec{cell(), format()}, run: runIsValidVertex},
	}
}

func runCellToVertex(env environment, p parsedArgs) error {
	cell := h3.Cell(rawHex(p.get("-c")))
	if !cell.IsValid() {
		return h3.ErrCellInvalid
	}
	n, err := p.integer("-v")
	if err != nil {
		return h3.ErrFailed
	}
	vertex, err := cell.Vertex(n)
	if err != nil {
		return err
	}
	return writeCell(env.out, h3.Cell(vertex), formatValue(p))
}

func runCellToVertexes(env environment, p parsedArgs) error {
	cell := h3.Cell(rawHex(p.get("-c")))
	if !cell.IsValid() {
		return h3.ErrCellInvalid
	}
	vertices, err := cell.Vertexes()
	if err != nil {
		return err
	}
	return writeCells(env.out, cellsFromVertices(vertices), formatValue(p))
}

func runVertexToLatLng(env environment, p parsedArgs) error {
	vertex := h3.Vertex(rawHex(p.get("-c")))
	if !vertex.IsValid() {
		return h3.ErrVertexInvalid
	}
	ll, err := vertex.LatLng()
	if err != nil {
		return err
	}
	switch formatValue(p) {
	case "json":
		writef(env.out, "[%.10f, %.10f]\n", ll.Lat.Deg(), ll.Lng.Deg())
	case "wkt":
		writef(env.out, "POINT(%.10f %.10f)\n", ll.Lng.Deg(), ll.Lat.Deg())
	case "newline":
		writef(env.out, "%.10f\n%.10f\n", ll.Lat.Deg(), ll.Lng.Deg())
	default:
		return h3.ErrFailed
	}
	return nil
}

func runIsValidVertex(env environment, p parsedArgs) error {
	return writeBool(env.out, h3.Vertex(rawHex(p.get("-c"))).IsValid(), formatValue(p))
}

func miscCommands() []command {
	cell := func() optionSpec { return opt("-c --cell", true) }
	res := func() optionSpec { return opt("-r --resolution", true) }
	format := func() optionSpec { return opt("-f --format", false) }
	commands := []command{
		{name: "degsToRads", description: "Converts degrees to radians", options: []optionSpec{opt("-d --degree", true)}, run: runDegsToRads},
		{name: "radsToDegs", description: "Converts radians to degrees", options: []optionSpec{opt("-r --radian", true)}, run: runRadsToDegs},
		metricResolutionCommand("getHexagonAreaAvgKm2", res(), h3.HexagonAreaAvgKm2),
		metricResolutionCommand("getHexagonAreaAvgM2", res(), h3.HexagonAreaAvgM2),
		metricCellCommand("cellAreaRads2", cell(), func(c h3.Cell) (float64, error) { return c.AreaRads2() }),
		metricCellCommand("cellAreaKm2", cell(), func(c h3.Cell) (float64, error) { return c.AreaKm2() }),
		metricCellCommand("cellAreaM2", cell(), func(c h3.Cell) (float64, error) { return c.AreaM2() }),
		metricResolutionCommand("getHexagonEdgeLengthAvgKm", res(), h3.HexagonEdgeLengthAvgKm),
		metricResolutionCommand("getHexagonEdgeLengthAvgM", res(), h3.HexagonEdgeLengthAvgM),
		metricEdgeCommand("edgeLengthRads", cell(), "%.10f\n", func(e h3.DirectedEdge) (float64, error) { return e.LengthRads() }),
		metricEdgeCommand("edgeLengthKm", cell(), "%.10f\n", func(e h3.DirectedEdge) (float64, error) { return e.LengthKm() }),
		metricEdgeCommand("edgeLengthM", cell(), "%.8f\n", func(e h3.DirectedEdge) (float64, error) { return e.LengthM() }),
		{name: "getNumCells", description: "Returns the number of cells at a resolution", options: []optionSpec{res()}, run: runGetNumCells},
		{name: "getRes0Cells", description: "Returns all resolution 0 cells", options: []optionSpec{format()}, run: runGetRes0Cells},
		{name: "getPentagons", description: "Returns all pentagons at a resolution", options: []optionSpec{res(), format()}, run: runGetPentagons},
		{name: "pentagonCount", description: "Returns 12", run: func(env environment, _ parsedArgs) error { writeln(env.out, h3.NumPentagons); return nil }},
		greatCircleCommand("greatCircleDistanceRads", h3.GreatCircleDistanceRads),
		greatCircleCommand("greatCircleDistanceKm", h3.GreatCircleDistanceKm),
		greatCircleCommand("greatCircleDistanceM", h3.GreatCircleDistanceM),
		{name: "describeH3Error", description: "Returns the description of an H3 error", options: []optionSpec{opt("-e --error", true)}, run: runDescribeH3Error},
	}
	return commands
}

// metricResolutionCommand builds a command that maps -r through fn and
// prints the value at upstream's %.10f precision. Upstream reuses each
// command's name as its help description for this group, hence
// description: name.
func metricResolutionCommand(name string, spec optionSpec, fn func(int) (float64, error)) command {
	return command{name: name, description: name, options: []optionSpec{spec}, run: func(env environment, p parsedArgs) error {
		res, err := p.integer("-r")
		if err != nil {
			return h3.ErrFailed
		}
		value, err := fn(res)
		if err != nil {
			return err
		}
		writef(env.out, "%.10f\n", value)
		return nil
	}}
}

func metricCellCommand(name string, spec optionSpec, fn func(h3.Cell) (float64, error)) command {
	return command{name: name, description: name, options: []optionSpec{spec}, run: func(env environment, p parsedArgs) error {
		cell := h3.Cell(rawHex(p.get("-c")))
		if !cell.IsValid() {
			return h3.ErrCellInvalid
		}
		value, err := fn(cell)
		if err != nil {
			return err
		}
		writef(env.out, "%.10f\n", value)
		return nil
	}}
}

// metricEdgeCommand takes the printf format explicitly because the
// group is no longer uniform: upstream prints edgeLengthRads/Km at
// %.10lf but edgeLengthM at %.8lf (changed from %.10lf in H3 4.5.0).
func metricEdgeCommand(name string, spec optionSpec, format string, fn func(h3.DirectedEdge) (float64, error)) command {
	return command{name: name, description: name, options: []optionSpec{spec}, run: func(env environment, p parsedArgs) error {
		value, err := fn(h3.DirectedEdge(rawHex(p.get("-c"))))
		if err != nil {
			return err
		}
		writef(env.out, format, value)
		return nil
	}}
}

func runDegsToRads(env environment, p parsedArgs) error {
	v, err := p.float("-d")
	if err != nil {
		return h3.ErrFailed
	}
	writef(env.out, "%.10f\n", h3.Deg(v).Rad())
	return nil
}

func runRadsToDegs(env environment, p parsedArgs) error {
	v, err := p.float("-r")
	if err != nil {
		return h3.ErrFailed
	}
	writef(env.out, "%.10f\n", h3.Rad(v).Deg())
	return nil
}

func runGetNumCells(env environment, p parsedArgs) error {
	res, err := p.integer("-r")
	if err != nil {
		return h3.ErrFailed
	}
	n, err := h3.NumCells(res)
	if err != nil {
		return err
	}
	writeln(env.out, n)
	return nil
}

func runGetRes0Cells(env environment, p parsedArgs) error {
	return writeCells(env.out, h3.Res0Cells(), formatValue(p))
}

func runGetPentagons(env environment, p parsedArgs) error {
	res, err := p.integer("-r")
	if err != nil {
		return h3.ErrFailed
	}
	cells, err := h3.Pentagons(res)
	if err != nil {
		return err
	}
	return writeCells(env.out, cells, formatValue(p))
}

// greatCircleCommand builds the three greatCircleDistance* commands. Their
// input is deliberately parsed with the polygon machinery — upstream does
// the same, treating "[[lat, lng], [lat, lng]]" as a two-vertex loop and
// rejecting anything that does not reduce to exactly one two-point loop.
func greatCircleCommand(name string, fn func(h3.LatLng, h3.LatLng) float64) command {
	return command{name: name, description: name, options: []optionSpec{opt("-i --file", false), opt("-c --coordinates", false)}, run: func(env environment, p parsedArgs) error {
		data, err := readSource(env, p, "-i", "-c")
		if err != nil {
			return err
		}
		polygon, err := parsePolygon(data)
		if err != nil || len(polygon.GeoLoop) != 2 || len(polygon.Holes) != 0 {
			return failDirect(env.errOut, "Only two pairs of coordinates should be provided.")
		}
		writef(env.out, "%.10f\n", fn(polygon.GeoLoop[0], polygon.GeoLoop[1]))
		return nil
	}}
}

func runDescribeH3Error(env environment, p parsedArgs) error {
	code, err := p.integer("-e")
	if err != nil {
		return h3.ErrFailed
	}
	if code == 0 {
		writeln(env.out, "Success")
	} else {
		writeln(env.out, errorDescription(code))
	}
	return nil
}
