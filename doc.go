// Package h3 is a pure-Go implementation of Uber's H3 hexagonal hierarchical
// geospatial indexing system, behaviorally equivalent to H3 C v4.3.0.
//
// The production library is safe Go only: no cgo and no unsafe. Behavioral
// equivalence is enforced by an opt-in parity test suite (build tags
// cgo && c2go) that compares every ported function against the original C
// objects compiled from pristine upstream sources.
//
// # Index types
//
// The three H3 index kinds are distinct uint64 types: [Cell] (a hexagon or
// pentagon at some resolution), [DirectedEdge] (a directed connection
// between adjacent cells), and [Vertex] (a topological cell corner). Raw
// conversions like Cell(0x8928308280fffff) are legal but unchecked; use
// IsValid, or parse validated values with [ParseCell], [ParseDirectedEdge],
// and [ParseVertex]. All three implement fmt.Stringer (canonical lowercase
// hex) and encoding.TextMarshaler/TextUnmarshaler, so they work with
// encoding/json out of the box.
//
// # Coordinates
//
// Geographic coordinates use [LatLng], whose fields are [Angle] values
// (stored in radians). Construct them explicitly — [LatLngDegs] for degrees,
// [NewLatLng] with [Deg] or [Rad] angles — so degree/radian mix-ups cannot
// compile. Accessors convert on the way out: ll.Lat.Deg(), ll.Lng.Rad().
//
// # Errors
//
// Operations return sentinel errors matching the H3 C error codes
// ([ErrPentagon], [ErrCellInvalid], [ErrResolutionDomain], ...); match them
// with errors.Is. Pure bit accessors (Resolution, BaseCellNumber, IsValid,
// String) do not fail.
//
// # Allocation control
//
// Collection-returning operations come in two forms: a convenience form that
// allocates its result (GridDisk, Children, PolygonToCells, ...) and an
// Append* form (AppendGridDisk, AppendChildren, AppendPolygonToCells, ...)
// that appends to a caller-provided buffer and allocates nothing when
// capacity suffices. Iterator forms ([Cell.ChildrenSeq], [CellsAtRes],
// [PolygonToCellsExperimentalSeq]) yield cells one at a time with zero
// allocation. [CellBoundary] is a fixed-size value type; obtaining a
// boundary performs no heap allocation.
//
// # Relationship to H3 C
//
// The implementation is a function-by-function port of the C library; every
// public operation's doc comment carries an "H3 C API:" line naming its C
// counterpart, and docs/c-api-inventory.csv maps the entire C API surface.
// All 75 public functions of H3 C v4.3.0 are covered. Intentional behavior
// differences (Go-idiomatic hole pruning, validated parsing, ...) are
// documented in docs/DEVIATIONS.md.
package h3
