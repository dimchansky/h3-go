package cli

// This file holds the input decoding and output rendering shared by several
// commands. Rendering matches upstream printf calls digit for digit — the
// scenario suite and the C differential tests compare output textually, so
// precision and separators here are contract, not style.

import (
	"encoding/json"
	"io"
	"os"
	"strconv"
	"strings"

	h3 "github.com/dimchansky/h3-go"
)

// rawHex parses a hex index argument permissively: parse errors yield 0,
// which then fails IsValid/domain checks downstream. This mirrors the C
// CLI's unchecked sscanf("%" PRIx64) — rejecting bad hex here would change
// which error (and exit code) the user sees.
func rawHex(s string) uint64 {
	v, _ := strconv.ParseUint(strings.TrimPrefix(s, "0x"), 16, 64)
	return v
}

// formatValue returns the normalized -f/--format value; upstream's default
// everywhere is json. Runners reject values their upstream counterpart does
// not accept.
func formatValue(p parsedArgs) string {
	if !p.has("-f") || p.get("-f") == "" {
		return "json"
	}
	return strings.ToLower(p.get("-f"))
}

// readSource fetches command input from exactly one of two mutually
// exclusive options: a file option (where the value "--" selects stdin) or
// a direct-value option. The diagnostics are exact upstream wording.
func readSource(env environment, p parsedArgs, fileKey, directKey string) ([]byte, error) {
	if p.has(fileKey) == p.has(directKey) {
		return nil, failDirect(env.errOut, "You must provide exactly one input source")
	}
	if p.has(directKey) {
		return []byte(p.get(directKey)), nil
	}
	if p.get(fileKey) == "--" {
		return io.ReadAll(env.in)
	}
	data, err := os.ReadFile(p.get(fileKey))
	if os.IsNotExist(err) {
		return nil, failDirect(env.errOut, "The specified file does not exist.")
	}
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, failDirect(env.errOut, "The specified file is empty.")
	}
	return data, nil
}

// parseCells extracts 15-hex-digit cell tokens from free-form text, exactly
// reproducing readCellsFromFile in the upstream CLI — including its
// 1500-byte streaming window. The C reader scans a fixed buffer, stops 15
// bytes short of the end (a token that would straddle the boundary is
// re-scanned after the next refill), and can silently skip malformed or
// boundary-straddling sequences. Those quirks are observable, so this is a
// quirk-for-quirk port. H3 4.5.0 fixed one genuine bug here — stale bytes
// from the previous chunk survived a short final read and were re-scanned
// as phantom cells (upstream's multipolygon fixture 5 used to depend on
// it) — and this port mirrors the fix: the stale region is zeroed after
// each refill, exactly like the upstream memset.
func parseCells(data []byte) []h3.Cell {
	const bufferSize = 1500
	const bufferSizeLessCell = bufferSize - 15
	buffer := make([]byte, bufferSize)
	readOffset := copy(buffer, data)
	out := make([]h3.Cell, 0, 128)
	for {
		bufferOffset := 0
		lastGoodOffset := 0
		for bufferOffset < bufferSizeLessCell {
			start := bufferOffset
			for start < bufferSize && (buffer[start] == ' ' || buffer[start] == '\t' || buffer[start] == '\r' || buffer[start] == '\n' || buffer[start] == '\v' || buffer[start] == '\f') {
				start++
			}
			end := start
			for end < bufferSize && isHex(buffer[end]) {
				end++
			}
			scanLen := end - bufferOffset
			if scanLen != 15 {
				bufferOffset++
				continue
			}
			value, _ := strconv.ParseUint(string(buffer[start:end]), 16, 64)
			out = append(out, h3.Cell(value))
			bufferOffset += 16
			lastGoodOffset = bufferOffset
		}
		if readOffset >= len(data) {
			break
		}
		if lastGoodOffset < bufferSizeLessCell {
			lastGoodOffset = bufferSizeLessCell
		}
		preserved := copy(buffer, buffer[lastGoodOffset:])
		n := copy(buffer[preserved:], data[readOffset:])
		readOffset += n
		// When the refill is short (final chunk), stale data from the
		// previous read remains in the buffer tail. The scan has no
		// concept of "valid length" — it will keep parsing whatever is
		// in the buffer — so zero the entire stale region (the H3 4.5.0
		// scanner fix; upstream h3.c does the same with memset).
		for i := preserved + n; i < bufferSize; i++ {
			buffer[i] = 0
		}
	}
	return out
}

func isHex(c byte) bool {
	return c >= '0' && c <= '9' || c >= 'a' && c <= 'f' || c >= 'A' && c <= 'F'
}

// parsePolygon decodes the upstream CLI's polygon JSON: arrays of
// [lat, lng] pairs (latitude first, unlike GeoJSON), at any nesting depth.
// The first loop found is the outer boundary; the rest are holes.
func parsePolygon(data []byte) (h3.GeoPolygon, error) {
	var value any
	if err := json.Unmarshal(data, &value); err != nil {
		return h3.GeoPolygon{}, h3.ErrFailed
	}
	var loops [][][]float64
	collectCoordinateLoops(value, &loops)
	if len(loops) == 0 {
		return h3.GeoPolygon{}, h3.ErrFailed
	}
	return loopsToPolygon(loops)
}

// collectCoordinateLoops walks arbitrarily nested JSON arrays and collects
// every array whose elements are all [number, number] pairs. Accepting any
// nesting depth matches upstream, which tolerates [[...]], [[[...]]], and
// GeoJSON-like [[[[...]]]] wrappers alike.
func collectCoordinateLoops(value any, loops *[][][]float64) {
	array, ok := value.([]any)
	if !ok {
		return
	}
	loop := make([][]float64, len(array))
	for i, item := range array {
		point, ok := item.([]any)
		if !ok || len(point) != 2 {
			loop = nil
			break
		}
		lat, latOK := point[0].(float64)
		lng, lngOK := point[1].(float64)
		if !latOK || !lngOK {
			loop = nil
			break
		}
		loop[i] = []float64{lat, lng}
	}
	if len(loop) > 0 {
		*loops = append(*loops, loop)
		return
	}
	for _, item := range array {
		collectCoordinateLoops(item, loops)
	}
}

func loopsToPolygon(loops [][][]float64) (h3.GeoPolygon, error) {
	if len(loops) == 0 {
		return h3.GeoPolygon{}, nil
	}
	convert := func(points [][]float64) (h3.GeoLoop, error) {
		out := make(h3.GeoLoop, len(points))
		for i, point := range points {
			if len(point) != 2 {
				return nil, h3.ErrFailed
			}
			out[i] = h3.LatLngDegs(point[0], point[1])
		}
		return out, nil
	}
	outer, err := convert(loops[0])
	if err != nil {
		return h3.GeoPolygon{}, err
	}
	polygon := h3.GeoPolygon{GeoLoop: outer}
	for _, raw := range loops[1:] {
		hole, err := convert(raw)
		if err != nil {
			return h3.GeoPolygon{}, err
		}
		polygon.Holes = append(polygon.Holes, hole)
	}
	return polygon, nil
}

// writeBool renders a boolean result: json prints true/false, numeric
// prints 1/0 (only the commands whose upstream counterpart offers numeric
// pass it through).
func writeBool(w io.Writer, value bool, format string) error {
	switch format {
	case "json":
		writeln(w, value)
	case "numeric":
		if value {
			writeln(w, 1)
		} else {
			writeln(w, 0)
		}
	default:
		return h3.ErrFailed
	}
	return nil
}

func writeCell(w io.Writer, cell h3.Cell, format string) error {
	switch format {
	case "json":
		writef(w, "\"%s\"\n", cell)
	case "newline":
		writeln(w, cell)
	default:
		return h3.ErrFailed
	}
	return nil
}

// writeCells renders a cell list: json as [ "hex", ... ] with upstream's
// exact spacing, newline as one cell per line.
func writeCells(w io.Writer, cells []h3.Cell, format string) error {
	switch format {
	case "json":
		writeText(w, "[ ")
		for i, cell := range cells {
			if i != 0 {
				writeText(w, ", ")
			}
			writef(w, "\"%s\"", cell)
		}
		writeln(w, " ]")
	case "newline":
		for _, cell := range cells {
			writeln(w, cell)
		}
	default:
		return h3.ErrFailed
	}
	return nil
}

func cellsFromVertices(vertices []h3.Vertex) []h3.Cell {
	out := make([]h3.Cell, len(vertices))
	for i, vertex := range vertices {
		out[i] = h3.Cell(vertex)
	}
	return out
}

// writeBoundary renders a cell or edge boundary. Coordinates print at
// upstream's %.10f precision; json/newline emit lat,lng order while WKT
// emits lng lat. For cells the WKT ring is explicitly closed by repeating
// the first vertex (POLYGON), for edges it stays open (LINESTRING).
func writeBoundary(w io.Writer, boundary h3.CellBoundary, format string, edge bool) error {
	verts := boundary.Verts()
	switch format {
	case "json":
		writeText(w, "[")
		for i, ll := range verts {
			if i > 0 {
				writeText(w, ", ")
			}
			writef(w, "[%.10f, %.10f]", ll.Lat.Deg(), ll.Lng.Deg())
		}
		writeln(w, "]")
	case "newline":
		for _, ll := range verts {
			writef(w, "%.10f\n%.10f\n", ll.Lat.Deg(), ll.Lng.Deg())
		}
	case "wkt":
		if edge {
			writeText(w, "LINESTRING (")
		} else {
			writeText(w, "POLYGON((")
		}
		for i, ll := range verts {
			if i > 0 {
				writeText(w, ", ")
			}
			writef(w, "%.10f %.10f", ll.Lng.Deg(), ll.Lat.Deg())
		}
		if !edge && len(verts) > 0 {
			writef(w, ", %.10f %.10f", verts[0].Lng.Deg(), verts[0].Lat.Deg())
		}
		if edge {
			writeln(w, ")")
		} else {
			writeln(w, "))")
		}
	default:
		return h3.ErrFailed
	}
	return nil
}
