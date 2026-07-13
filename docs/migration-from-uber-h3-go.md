# Migrating from uber/h3-go/v4 — and what you gain

A practical guide for moving existing code from the official cgo binding
[`github.com/uber/h3-go/v4`](https://github.com/uber/h3-go) (pinned
reference: **v4.4.1**, vendoring H3 C v4.4.1) to this library (behavioral
target: **H3 C v4.4.0** — the same H3 behavior; see
[the version rationale](comparison-uber-h3-go.md#versions-compared)).

**Effort in one sentence:** almost everything is a mechanical rename or
re-type that the compiler finds for you; the genuinely behavioral work is
(1) coordinate construction — degrees become typed `Angle` values — and
(2) the few APIs whose result *shape* differs (`CellBoundary`,
`GridDiskDistances`, `DirectedEdge.Cells`, `GridDisksUnsafe`).

The [function-by-function matrix](comparison-uber-h3-go.md#function-by-function-matrix)
classifies all 78 H3 operations as `mechanical`/`adaptation`; this guide
covers the patterns. The realistic example at the bottom is kept
compiling and semantically verified against the binding by
[interop/uberbench/migration_test.go](../interop/uberbench/migration_test.go).

## What you gain after the migration

This is more than an import-path replacement. At the common **H3 4.4**
feature level, this library exposes the complete C API and adds Go-native
ways to avoid materializing or reallocating large results:

| Capability | What becomes available here |
|---|---|
| Complete H3 4.4 surface | All 78 public C functions, including `ConstructCell` and the single-origin `Cell.GridDiskUnsafe`; uber/h3-go v4.4.1 has no equivalent for those two operations. |
| Streaming | `Cell.ChildrenSeq`, `CellsAtRes`, and `PolygonToCellsExperimentalSeq` return `iter.Seq` values and do not materialize the complete result. The binding has no iterator API. |
| Caller-owned buffers | The main variable-size hierarchy, traversal, compaction, and polyfill operations have `Append*` forms; a sufficiently sized warm buffer makes the result path allocation-free where the algorithm itself needs no scratch space. |
| Planning and result shapes | Exported `Max*Size`, `UncompactCellsSize`, `Cell.NumChildren`, and `Cell.GridPathLen` helpers support pre-sizing; `GridDiskDistancesGrouped` offers binding-style rings while the flat form remains available. |
| Stronger Go contracts | Typed `Angle` values prevent degree/radian mix-ups; typed parsers validate the index mode and return sentinel errors instead of silently producing zero. |
| Go-native operations | `CGO_ENABLED=0`, ordinary cross-compilation and profiling, no hidden C heap, plus an upstream-compatible pure-Go `h3` CLI. |

The scope matters: “complete” refers to the shared H3 4.4 C API, not to a
strict superset of every binding-specific convenience. The binding retains
raw `IndexFromString`/`IndexToString` helpers and an optional experimental
polyfill result cap; current uber/h3-go v4.5.0 also tracks newer H3 C
functionality, including `DirectedEdge.Reverse`, which this library will gain
with its H3 4.5 sync. See the
[versioned coverage summary](comparison-uber-h3-go.md#coverage-summary) for
both sides of that trade-off.

## Import path

```go
import h3 "github.com/uber/h3-go/v4"   // before
import h3 "github.com/dimchansky/h3-go" // after
```

With the alias, most call sites keep reading `h3.…`.

## Type mappings

| uber/h3-go | This library | Adaptation |
|---|---|---|
| `Cell int64` | `Cell uint64` | Values are bit-identical (valid indexes never set the high bit). Convert explicitly where you stored raw ints: `h3.Cell(uint64(v))`. |
| `DirectedEdge int64` | `DirectedEdge uint64` | Same. |
| `Vertex int64` | `Vertex uint64` | Same. |
| `Index` (generic constraint) | `Index` | This library's constraint additionally accepts raw `uint64` and legacy integer-literal calls. |
| `LatLng{Lat, Lng float64}` (degrees) | `LatLng{Lat, Lng Angle}` (typed; radians inside) | **The one change you must not skip** — see below. |
| `CellBoundary []LatLng` | `CellBoundary` (opaque value; `Len()`, `At(i)`, `Verts()`) | Replace indexing/`len`/`range` with the accessors. |
| `GeoLoop []LatLng` | `GeoLoop []LatLng` | Same shape; vertices are typed. |
| `GeoPolygon{GeoLoop, Holes}` | `GeoPolygon{GeoLoop, Holes}` | Same shape. |
| `CoordIJ{I, J int}` | `CoordIJ{I, J int32}` | `int32` deliberately preserves C overflow behavior. |
| `ContainmentMode` + `ContainmentOverlappingBbox` | `ContainmentMode` + `ContainmentOverlappingBBox` | Same values; note the `BBox` capitalization. |
| `MaxCellBndryVerts` | `MaxCellBoundaryVerts` | Rename. |
| `DegsToRads` / `RadsToDegs` constants | `Deg(x)` / `Rad(x)` constructors, `.Deg()` / `.Rad()` accessors | Multiplications become typed conversions. |
| `InvalidH3Index = 0` | zero values (`Cell(0)` is invalid) | Compare against `0` or use `IsValid`. |
| `NumIcosaFaces` | `NumIcosahedronFaces` | Rename; face numbers remain `0..19`. |

## Coordinates: degrees become `Angle`

The binding stores degrees in bare `float64` fields; this library wraps
coordinates in the `Angle` type so the unit is part of the type system.
Three call-site patterns cover it:

```go
// construct                       // before                        // after
h3.NewLatLng(37.77, -122.41)       // degrees implied               h3.LatLngDegs(37.77, -122.41)
h3.LatLng{Lat: la, Lng: lo}        // degrees implied               h3.NewLatLng(h3.Deg(la), h3.Deg(lo))

// read back
ll.Lat, ll.Lng                     // float64 degrees               ll.Lat.Deg(), ll.Lng.Deg()

// radians (e.g. your own trig)
ll.Lat * h3.DegsToRads             //                               ll.Lat.Rad()
```

Do **not** port `h3.LatLng{Lat: 37.77, Lng: -122.41}` literally: here a
bare float in that struct is an `Angle` in **radians**, and the compiler
only saves you if you keep the fields typed (it will reject the raw
float literal). Always construct through `LatLngDegs`/`NewLatLng`/`Deg`.

## Function and method mappings

Only renamed or reshaped APIs are listed; everything not listed keeps its
name and shape (modulo `error` returns noted in the next section). The
[full matrix](comparison-uber-h3-go.md#function-by-function-matrix) covers
every operation.

| uber/h3-go | This project | Required adaptation |
|---|---|---|
| `LatLngToCell(ll, res)` / `ll.Cell(res)` | same | construct `ll` with `LatLngDegs` |
| `CellToLatLng(c)` | `c.LatLng()` | method form; read with `.Deg()` |
| `CellToBoundary(c)` / `c.Boundary()` | `c.Boundary()` | result is a value: `b.Len()`, `b.At(i)`, `b.Verts()` |
| `CellFromString(s)` (+ manual `IsValid`) | `ParseCell(s)` | returns `(Cell, error)`; delete the manual validity check |
| `IndexFromString(s)` | `ParseCell` / `ParseDirectedEdge` / `ParseVertex` | pick the typed parser; raw `strconv.ParseUint(s, 16, 64)` if you truly want unvalidated bits |
| `CellToString(c)` / `IndexToString(u)` | `c.String()` | same text |
| `c.ImmediateParent()` | `c.ImmediateParent()` | same |
| `c.ImmediateChildren()` | `c.ImmediateChildren()` | same; `AppendImmediateChildren` reuses a capacity-7 buffer |
| `ChildPosToCell(pos, c, res)` / `c.ChildPosToCell(pos, res)` | `c.ChildAtPos(int64(pos), res)` | rename; position is `int64` |
| `CellToChildPos(c, res)` / `c.ChildPos(res)` | `c.ChildPos(res)` | result is `int64` |
| `NumCells(res)` (`int`, panics on bad res) | `NumCells(res)` (`(int64, error)`) | handle the error; result is `int64` |
| `Res0Cells()` (`([]Cell, error)`) | `Res0Cells()` (`[]Cell`) | drop the error handling |
| `c.DirectedEdge(other)` | `c.DirectedEdgeTo(other)` | rename |
| `e.Cells()` (`[]Cell{origin, dest}`) | `e.Cells()` (`(origin, destination Cell, err)`) | destructure instead of indexing |
| `EdgeLengthKm(e)` (`Rads`/`M`) | `e.LengthKm()` (`LengthRads`/`LengthM`) | method form |
| `CellAreaKm2(c)` (`Rads2`/`M2`) | `c.AreaKm2()` (`AreaRads2`/`AreaM2`) | method form |
| `BaseCellNumber(c)` | `c.BaseCellNumber()` | method form |
| `GridDisk(c, k)` | `c.GridDisk(k)` | method form (free function removed) |
| `GridDiskDistances(c, k)` (`[][]Cell` rings) | `c.GridDiskDistancesGrouped(k)` (`[][]Cell`) | method form; null holes are pruned here |
| `GridDisksUnsafe(origins, k)` (`[][]Cell`) | `GridDisksUnsafe(origins, k)` (`[]Cell`, flat stride `MaxGridDiskSize(k)`) | slice per origin; unpruned zeros preserved (C layout) |
| `GridDistance(a, b)` / `GridPath(a, b)` | `a.GridDistance(b)` / `a.GridPath(b)` | method form |
| `PolygonToCellsExperimental(p, res, mode, cap)` | `PolygonToCellsExperimental(p, res, mode)` | no variadic cap; pre-size via `MaxPolygonToCellsSizeExperimental` + `Append…`, or stream with `PolygonToCellsExperimentalSeq` |
| `IsValidIndex[T](idx)` | `IsValidIndex(idx)` | same typed call; raw `uint64` is also accepted |
| `v.IndexDigit(res)` / `e.IndexDigit(res)` | same | same |

Use the grouped convenience when existing code consumes `[][]Cell`:

```go
rings, err := c.GridDiskDistancesGrouped(k)
```

The original flat `GridDiskDistances` plus its zero-allocation `Append` form
remain available for allocation-sensitive code. Most single-ring consumers
want `c.GridRing(d)` directly.

## Error handling

Both libraries expose sentinel errors mapped from the same C error codes
and matched with `errors.Is`; migration is mostly renames:

- `h3.ErrRsolutionMismatch` (binding's spelling) → `h3.ErrResolutionMismatch`.
- The binding's `UnmarshalText` returns ad-hoc `errors.New("invalid cell
  index")` values; here `UnmarshalText` and `Parse*` return the sentinel
  (`ErrCellInvalid` etc.) wrapped with context, so `errors.Is` works.
- `Res0Cells` loses its (never-firing) error; `NumCells` gains one
  (instead of panicking on a bad resolution).
- Everything else keeps its error shape: fallible operations return
  `(T, error)` in both libraries.

## Parsing and text marshaling

`String()`, `MarshalText`, `UnmarshalText` exist on `Cell`,
`DirectedEdge`, and `Vertex` in both libraries and produce/accept the
identical canonical lowercase-hex text, so **JSON and text round-trips are
wire-compatible** — you can migrate producers and consumers independently.
The behavioral upgrades here: `Parse*` reports syntax errors instead of
returning index 0, and validates the index mode.

## Ordering, pruning, and null slots

Semantics match the binding on everything the equivalence suite covers
(that is, both match H3 C): disks/rings/polyfills are unordered with
`H3_NULL` holes pruned; children are in canonical order; paths are
inclusive and ordered; pentagons yield 5 edges/vertexes. One asymmetry to
re-check if you handled it specially: the binding's `GridDiskDistances*`
rings can contain zero slots near pentagons; the flat form here prunes
them (the `dists` slice tells you the ring).

## Putting the allocation-sensitive upgrades to work

Not required, but the reason many users switch:

```go
buf := make([]h3.Cell, 0, 64)
for _, c := range hot {
    buf, _ = c.AppendGridDisk(buf[:0], 2)   // 0 allocs/op once warm
    consume(buf)
}

for child := range cell.ChildrenSeq(12) {   // stream, never materialize
    visit(child)
}
```

The main variable-size hierarchy, traversal, compaction, and polyfill APIs
have `Append*` forms (pre-size with the exported
`Max*Size`/`NumChildren`/`GridPathLen` helpers); three iterators stream
without materializing. The binding has no equivalent caller-buffer or
iterator forms.

## A realistic before/after

Before (uber/h3-go):

```go
func coverage(cellStr string, fence h3.GeoPolygon /* degrees */) ([]string, error) {
	c := h3.CellFromString(cellStr)
	if !c.IsValid() {
		return nil, fmt.Errorf("bad cell %q", cellStr)
	}
	disk, err := c.GridDisk(1)
	if err != nil {
		return nil, err
	}
	area, err := h3.PolygonToCells(fence, c.Resolution())
	if err != nil {
		return nil, err
	}
	compact, err := h3.CompactCells(append(area, disk...))
	if err != nil {
		return nil, err
	}
	out := make([]string, len(compact))
	for i, cc := range compact {
		out[i] = h3.CellToString(cc)
	}
	return out, nil
}
```

After (this library — same behavior, verified in
[migration_test.go](../interop/uberbench/migration_test.go)):

```go
func coverage(cellStr string, fence h3.GeoPolygon /* LatLngDegs verts */) ([]string, error) {
	c, err := h3.ParseCell(cellStr) // parse + validate in one step
	if err != nil {
		return nil, fmt.Errorf("bad cell %q: %w", cellStr, err)
	}
	disk, err := c.GridDisk(1)
	if err != nil {
		return nil, err
	}
	area, err := h3.PolygonToCells(fence, c.Resolution())
	if err != nil {
		return nil, err
	}
	compact, err := h3.CompactCells(append(area, disk...))
	if err != nil {
		return nil, err
	}
	out := make([]string, len(compact))
	for i, cc := range compact {
		out[i] = cc.String()
	}
	return out, nil
}
```

The diff is three lines: validated parsing replaces
`CellFromString`+`IsValid`, `String()` replaces `CellToString`, and the
polygon's vertices are built with `LatLngDegs` instead of `NewLatLng`.
(Note: `CompactCells` in both libraries requires unique cells — the
`append` is fine when the disk and fence don't overlap, as in the tested
scenario; deduplicate otherwise.)

## Why there is no compatibility adapter package

We considered shipping a `compat` package mirroring the binding's exact
API and decided against it, for now:

- **Maintenance cost**: it would track two moving APIs forever, and its
  surface would need the same golden-lock and parity treatment as the
  main API.
- **API-stability implications**: pre-v1.0, freezing a second public
  surface that mirrors someone else's API multiplies the cost of every
  breaking change.
- **It would obscure the type safety that motivates switching**: an
  adapter must accept bare `float64` degrees and silently re-wrap them,
  reintroducing the unit confusion `Angle` exists to prevent.
- **The migration is compiler-driven**: nearly every change is a rename
  or re-type that fails loudly at build time; this guide plus compiler
  errors have been sufficient in practice for the tested scenarios.

If real-world migrations surface friction this guide cannot cover, an
adapter would be a new public package and therefore a deliberate
architectural decision — proposed as a DR in
[public-api-architecture.md](public-api-architecture.md), not added ad hoc.
