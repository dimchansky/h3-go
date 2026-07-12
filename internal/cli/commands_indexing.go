package cli

// Indexing and inspection commands (cellToLatLng ... getIcosahedronFaces).

import (
	"strconv"
	"strings"

	h3 "github.com/dimchansky/h3-go"
)

func indexingCommands() []command {
	cell := func() optionSpec { return opt("-c --cell", true) }
	format := func() optionSpec { return opt("-f --format", false) }
	resolution := func(required bool) optionSpec { return opt("-r --resolution", required) }
	return []command{
		{name: "cellToLatLng", description: "Convert an H3Cell to a coordinate", options: []optionSpec{cell(), format()}, run: runCellToLatLng},
		{name: "latLngToCell", description: "Convert degrees latitude/longitude coordinate to an H3 cell", options: []optionSpec{resolution(true), opt("--lat --latitude", true), opt("--lng --longitude", true), format()}, run: runLatLngToCell},
		{name: "cellToBoundary", description: "Convert an H3 cell to a polygon defining its boundary", options: []optionSpec{cell(), format()}, run: runCellToBoundary},
		{name: "getResolution", description: "Extracts the resolution (0 - 15) from the H3 cell", options: []optionSpec{cell()}, run: runGetResolution},
		{name: "getBaseCellNumber", description: "Extracts the base cell number (0 - 121) from the H3 cell", options: []optionSpec{cell()}, run: runGetBaseCellNumber},
		{name: "getIndexDigit", description: "Extracts the indexing digit (0 - 7) from the H3 cell", options: []optionSpec{cell(), opt("-r --res", true)}, run: runGetIndexDigit},
		{name: "constructCell", description: "Construct an H3 cell from resolution, base cell, and digits", options: []optionSpec{resolution(false), opt("-b --baseCell", true), opt("-d --digits", false), format()}, run: runConstructCell},
		{name: "stringToInt", description: "Converts an H3 index in string form to integer form", options: []optionSpec{cell()}, run: runStringToInt},
		{name: "intToString", description: "Converts an H3 index in int form to string form", options: []optionSpec{cell()}, run: runIntToString},
		{name: "isValidCell", description: "Checks if the provided H3 index is actually valid", options: []optionSpec{cell(), format()}, run: runIsValidCell},
		{name: "isResClassIII", description: "Checks if the provided H3 index has a Class III orientation", options: []optionSpec{cell(), format()}, run: runIsResClassIII},
		{name: "isPentagon", description: "Checks if the provided H3 index is a pentagon instead of a hexagon", options: []optionSpec{cell(), format()}, run: runIsPentagon},
		{name: "getIcosahedronFaces", description: "Returns the icosahedron face numbers that the H3 index intersects", options: []optionSpec{cell(), format()}, run: runGetIcosahedronFaces},
	}
}

func runCellToLatLng(env environment, p parsedArgs) error {
	cell := h3.Cell(rawHex(p.get("-c")))
	if !cell.IsValid() {
		return h3.ErrCellInvalid
	}
	ll, err := cell.LatLng()
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

func runLatLngToCell(env environment, p parsedArgs) error {
	res, err := p.integer("-r")
	if err != nil {
		return h3.ErrFailed
	}
	lat, err := p.float("--lat")
	if err != nil {
		return h3.ErrFailed
	}
	lng, err := p.float("--lng")
	if err != nil {
		return h3.ErrFailed
	}
	cell, err := h3.LatLngToCell(h3.LatLngDegs(lat, lng), res)
	if err != nil {
		return err
	}
	return writeCell(env.out, cell, formatValue(p))
}

func runCellToBoundary(env environment, p parsedArgs) error {
	cell := h3.Cell(rawHex(p.get("-c")))
	if !cell.IsValid() {
		return h3.ErrCellInvalid
	}
	boundary, err := cell.Boundary()
	if err != nil {
		return err
	}
	return writeBoundary(env.out, boundary, formatValue(p), false)
}

func runGetResolution(env environment, p parsedArgs) error {
	raw := rawHex(p.get("-c"))
	if !h3.IsValidIndex(raw) {
		return h3.ErrIndexInvalid
	}
	writeln(env.out, h3.Cell(raw).Resolution())
	return nil
}

func runGetBaseCellNumber(env environment, p parsedArgs) error {
	raw := rawHex(p.get("-c"))
	if !h3.IsValidIndex(raw) {
		return h3.ErrIndexInvalid
	}
	writeln(env.out, h3.Cell(raw).BaseCellNumber())
	return nil
}

func runGetIndexDigit(env environment, p parsedArgs) error {
	raw := rawHex(p.get("-c"))
	if !h3.IsValidIndex(raw) {
		return h3.ErrIndexInvalid
	}
	res, err := p.integer("-r")
	if err != nil {
		return h3.ErrFailed
	}
	digit, err := h3.Cell(raw).IndexDigit(res)
	if err != nil {
		return err
	}
	writeln(env.out, digit)
	return nil
}

// runConstructCell mirrors upstream's argument semantics: -d is a
// comma-separated list of single digits 0–6; the resolution defaults to the
// number of digits and, when -r is given explicitly, must equal it.
func runConstructCell(env environment, p parsedArgs) error {
	var digits []int
	if p.has("-d") {
		for _, digit := range strings.Split(p.get("-d"), ",") {
			if len(digit) != 1 || digit[0] < '0' || digit[0] > '6' {
				return h3.ErrDigitDomain
			}
			digits = append(digits, int(digit[0]-'0'))
		}
	}
	res := len(digits)
	if p.has("-r") {
		var err error
		res, err = p.integer("-r")
		if err != nil {
			return h3.ErrFailed
		}
	}
	if res < 0 || res > h3.MaxResolution {
		return h3.ErrResolutionDomain
	}
	if len(digits) != res {
		return h3.ErrDigitDomain
	}
	base, err := p.integer("-b")
	if err != nil {
		return h3.ErrFailed
	}
	cell, err := h3.ConstructCell(res, base, digits)
	if err != nil {
		return err
	}
	return writeCell(env.out, cell, formatValue(p))
}

func runStringToInt(env environment, p parsedArgs) error {
	raw, err := strconv.ParseUint(strings.TrimPrefix(p.get("-c"), "0x"), 16, 64)
	if err != nil {
		return h3.ErrFailed
	}
	writeln(env.out, raw)
	return nil
}

func runIntToString(env environment, p parsedArgs) error {
	raw, err := strconv.ParseUint(p.get("-c"), 10, 64)
	if err != nil {
		return h3.ErrFailed
	}
	writef(env.out, "%x\n", raw)
	return nil
}

func runIsValidCell(env environment, p parsedArgs) error {
	return writeBool(env.out, h3.Cell(rawHex(p.get("-c"))).IsValid(), formatValue(p))
}

func runIsResClassIII(env environment, p parsedArgs) error {
	raw := rawHex(p.get("-c"))
	if !h3.IsValidIndex(raw) {
		return h3.ErrIndexInvalid
	}
	return writeBool(env.out, h3.Cell(raw).IsResClassIII(), formatValue(p))
}

func runIsPentagon(env environment, p parsedArgs) error {
	raw := rawHex(p.get("-c"))
	if !h3.IsValidIndex(raw) {
		return h3.ErrIndexInvalid
	}
	return writeBool(env.out, h3.Cell(raw).IsPentagon(), formatValue(p))
}

func runGetIcosahedronFaces(env environment, p parsedArgs) error {
	faces, err := h3.Cell(rawHex(p.get("-c"))).IcosahedronFaces()
	if err != nil {
		return err
	}
	switch formatValue(p) {
	case "json":
		writeText(env.out, "[")
		for i, face := range faces {
			if i > 0 {
				writeText(env.out, ", ")
			}
			writeText(env.out, face)
		}
		writeln(env.out, "]")
	case "newline":
		for _, face := range faces {
			writeln(env.out, face)
		}
	default:
		return h3.ErrFailed
	}
	return nil
}
