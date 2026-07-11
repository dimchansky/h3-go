package h3

import (
	"encoding"
	"fmt"
	"strconv"
	"strings"
)

// index constrains the three public H3 index types for the shared parse,
// format, and marshal helpers.
type index interface {
	Cell | DirectedEdge | Vertex
}

// appendIndex appends the canonical lowercase-hex form (no "0x" prefix,
// matching C h3ToString) of i to dst.
func appendIndex[T index](dst []byte, i T) []byte {
	return strconv.AppendUint(dst, uint64(i), 16)
}

// parseIndex parses the hexadecimal string form of an H3 index (an optional
// "0x"/"0X" prefix is accepted) and validates it with the supplied predicate.
// Unlike C stringToH3, syntax errors are reported rather than swallowed.
func parseIndex[T index](s string, valid func(T) bool, invalidErr error) (T, error) {
	hexStr := strings.TrimPrefix(strings.TrimPrefix(s, "0x"), "0X")
	v, err := strconv.ParseUint(hexStr, 16, 64)
	if err != nil {
		return 0, fmt.Errorf("h3: invalid index string %q: %w", s, err)
	}
	i := T(v)
	if !valid(i) {
		return 0, invalidErr
	}
	return i, nil
}

// ParseCell parses the hexadecimal string form of a cell index (as produced
// by Cell.String; an optional "0x" prefix is accepted) and validates it.
//
// H3 C API: stringToH3 (+ isValidCell).
func ParseCell(s string) (Cell, error) {
	return parseIndex[Cell](s, Cell.IsValid, ErrCellInvalid)
}

// ParseDirectedEdge parses the hexadecimal string form of a directed edge
// index and validates it.
//
// H3 C API: stringToH3 (+ isValidDirectedEdge).
func ParseDirectedEdge(s string) (DirectedEdge, error) {
	return parseIndex[DirectedEdge](s, func(e DirectedEdge) bool { return isValidDirectedEdge(h3Index(e)) }, ErrDirectedEdgeInvalid)
}

// ParseVertex parses the hexadecimal string form of a vertex index and
// validates it.
//
// H3 C API: stringToH3 (+ isValidVertex).
func ParseVertex(s string) (Vertex, error) {
	return parseIndex[Vertex](s, func(v Vertex) bool { return isValidVertex(h3Index(v)) }, ErrVertexInvalid)
}

// String returns the canonical lowercase-hex form of the cell index, without
// a "0x" prefix (e.g. "8928308280fffff").
//
// H3 C API: h3ToString.
func (c Cell) String() string {
	var buf [16]byte
	return string(appendIndex(buf[:0], c))
}

// String returns the canonical lowercase-hex form of the directed edge index.
//
// H3 C API: h3ToString.
func (e DirectedEdge) String() string {
	var buf [16]byte
	return string(appendIndex(buf[:0], e))
}

// String returns the canonical lowercase-hex form of the vertex index.
//
// H3 C API: h3ToString.
func (v Vertex) String() string {
	var buf [16]byte
	return string(appendIndex(buf[:0], v))
}

// MarshalText implements encoding.TextMarshaler using the canonical hex form.
func (c Cell) MarshalText() ([]byte, error) { return appendIndex(nil, c), nil }

// UnmarshalText implements encoding.TextUnmarshaler; it accepts the forms
// ParseCell accepts and validates the index.
func (c *Cell) UnmarshalText(text []byte) error {
	parsed, err := ParseCell(string(text))
	if err != nil {
		return err
	}
	*c = parsed
	return nil
}

// MarshalText implements encoding.TextMarshaler using the canonical hex form.
func (e DirectedEdge) MarshalText() ([]byte, error) { return appendIndex(nil, e), nil }

// UnmarshalText implements encoding.TextUnmarshaler; it accepts the forms
// ParseDirectedEdge accepts and validates the index.
func (e *DirectedEdge) UnmarshalText(text []byte) error {
	parsed, err := ParseDirectedEdge(string(text))
	if err != nil {
		return err
	}
	*e = parsed
	return nil
}

// MarshalText implements encoding.TextMarshaler using the canonical hex form.
func (v Vertex) MarshalText() ([]byte, error) { return appendIndex(nil, v), nil }

// UnmarshalText implements encoding.TextUnmarshaler; it accepts the forms
// ParseVertex accepts and validates the index.
func (v *Vertex) UnmarshalText(text []byte) error {
	parsed, err := ParseVertex(string(text))
	if err != nil {
		return err
	}
	*v = parsed
	return nil
}

var (
	_ fmt.Stringer             = Cell(0)
	_ encoding.TextMarshaler   = Cell(0)
	_ encoding.TextUnmarshaler = (*Cell)(nil)
	_ fmt.Stringer             = DirectedEdge(0)
	_ encoding.TextMarshaler   = DirectedEdge(0)
	_ encoding.TextUnmarshaler = (*DirectedEdge)(nil)
	_ fmt.Stringer             = Vertex(0)
	_ encoding.TextMarshaler   = Vertex(0)
	_ encoding.TextUnmarshaler = (*Vertex)(nil)
)
