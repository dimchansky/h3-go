# Public API ergonomics review (v0.x, pre-v1.0.0)

A systematic review of every API-shape difference between this library and
the official cgo binding [uber/h3-go v4.4.1](https://github.com/uber/h3-go),
using the [migration guide](migration-from-uber-h3-go.md) and the
[function-by-function matrix](comparison-uber-h3-go.md#function-by-function-matrix)
as the friction inventory. The question asked of each difference is not
"does uber do it differently?" but "is the difference *earning* something —
and if not, can an additive API remove the friction without weakening the
design invariants?"

Reviewed: 2026-07-12, against uber/h3-go **v4.4.1** (`h3.go` in the module
cache) and this library at the current commit. Related prior decisions:
[public-api-architecture.md](public-api-architecture.md) (DR-002, DR-003,
§12-Q5/Q7/Q8/Q10/Q11) and [FUTURE_WORK.md](FUTURE_WORK.md) (items 2 and
"IndexDigit for edges/vertexes", both resolved by this review).

Decisions use four classifications:

- **implement** — additive API lands now (see the CHANGELOG entry);
- **keep** — the current design is deliberately better; the migration guide
  documents the adaptation;
- **defer** — plausibly useful, but not until demand or evidence appears;
- **reject** — adopting uber's shape would make this library worse.

Invariants that were treated as non-negotiable throughout (from
[FUTURE_WORK.md](FUTURE_WORK.md) standing constraints and the architecture
document): pure, safe Go; `Angle` unit safety; validated typed parsing;
sentinel error consistency; zero-allocation `Append*` and value
`CellBoundary` forms as the primary surface; one obvious form per
operation; explicit failure instead of silent invalid zero values.

## Summary table

| Candidate | Current API | Uber API | Options considered | Performance impact | Compatibility impact | Decision |
|---|---|---|---|---|---|---|
| Immediate parent | `c.Parent(c.Resolution()-1)` | `c.ImmediateParent()` | method sugar; keep composition | none (delegates to `Parent`) | additive | **implement** `Cell.ImmediateParent` |
| Immediate children | `c.Children(c.Resolution()+1)` | `c.ImmediateChildren()` | delegate; specialized `makeDirectChild` loop; keep composition | specialized path measured (see §1); `Append` form 0-alloc warm | additive | **implement** `Cell.ImmediateChildren` + `Cell.AppendImmediateChildren` |
| `NumImmediateChildren` | `c.NumChildren(res+1)` | — (uber has no equivalent) | dedicated method; document counts | avoids `_ipow` — trivial | additive | **defer** (count is 7/6, documented; buffer cap 7 always suffices) |
| `ImmediateChildrenSeq` | `c.ChildrenSeq(res+1)` | — | dedicated iterator | none | additive | **defer** (≤7 elements; existing Seq composes) |
| Generic index validation | `IsValidIndex(uint64)` | `IsValidIndex[T Index]`, exported `Index` | generic over exported constraint incl. typed indexes, `uint64`, and legacy inferred `int`; keep raw-only; both forms | zero (compile-time, inlined) | source-compatible signature generalization | **implement** generic `IsValidIndex[T Index]` + exported `Index` constraint |
| Generic `IndexDigit[T]` helper | `Cell.IndexDigit` only | methods on all three types (via internal generic) | methods on `DirectedEdge`/`Vertex`; generic free function | none | additive | **implement** the two methods; **reject** the free function (one obvious form) |
| `CoordIJ` field types | `{I, J int32}` | `{I, J int}` | switch to `int`; internal `coordIJ32` + boundary conversion; keep `int32` | conversion layer would add copies | breaking if changed | **keep** `int32` (see §3 — uber's `int` silently truncates at its cgo boundary) |
| `NumIcosaFaces` | not exported (`numIcosaFaces` internal) | `NumIcosaFaces = 20` | export under uber's name; export under this library's naming; keep literal-20 advice | none | additive | **implement** as `NumIcosahedronFaces` (naming follows `Cell.IcosahedronFaces`) |
| Grouped disk distances | flat `([]Cell, []int32)` | `[][]Cell` rings (the only form) | grouped convenience atop flat core (FUTURE_WORK §2 design); keep flat-only | 3 allocs, counting-sort partition; flat/`Append` path untouched | additive | **implement** `Cell.GridDiskDistancesGrouped` |
| `GridDisksUnsafe` shape | flat `[]Cell`, stride `MaxGridDiskSize(k)` | `[][]Cell` pruned per origin | reshape; add grouped wrapper | grouped wrapper would add origins+1 allocs | additive if wrapped | **keep** flat (exact C layout, 1 alloc); slicing recipe documented |
| `DirectedEdge.Cells` | `(origin, destination Cell, err)` | `([]Cell, error)` | tuple; slice; both | tuple is 0-alloc, slice is 1-alloc | breaking if changed | **keep** tuple |
| Free-function duals | method-only (one form) | ~17 operations have both forms | add duals; keep one form | none | additive but duplicative | **reject** duals (§12-Q10 stands) |
| Unvalidated parse helpers | `ParseCell` etc. (validated, sentinel errors) | `CellFromString`/`IndexFromString` swallow errors, return 0 | compatibility aliases | none | additive but trap-restoring | **reject** aliases |
| `InvalidH3Index` constant | zero values documented per type | `InvalidH3Index = 0` | export the constant | none | additive | **reject** (duplicate spelling of the zero-value idiom) |
| `DegsToRads`/`RadsToDegs` | `Deg`/`Rad`/`.Deg()`/`.Rad()` + `RadPerDeg`/`DegPerRad` | multiply-by-constant | rename constants to uber's names | none | additive | **keep** (constants exist under accurate names; `Angle` is the real API) |
| `CellBoundary` iteration | `Len`/`At`/`Verts` (zero-alloc value) | `[]LatLng` slice | add `iter.Seq[LatLng]`; make it a slice | slice form costs 1 alloc/call | breaking if changed | **keep**; `range b.Verts()` already iterates without allocation |
| `NumCells` shape | `(int64, error)` | `int`, **panics** on bad res | match uber; keep | none | breaking if changed | **keep** (explicit error beats panic; counts are int64 in C) |
| `ChildPos`/`ChildAtPos` types | `int64` positions, `ChildAtPos` name | `int` positions, `ChildPosToCell` name | rename/re-type | none | breaking if changed | **keep** (C `int64_t` parity; `ChildAtPos` states inverse relation) |
| `Res0Cells` error shape | `[]Cell` (cannot fail) | `([]Cell, error)` (never fires) | add error back | none | breaking if changed | **keep** |
| Containment naming | `ContainmentOverlappingBBox` | `ContainmentOverlappingBbox` | match uber's casing | none | breaking if changed | **keep** (Go initialism convention) |
| Index base types | `uint64` | `int64` | match uber | none | breaking if changed | **keep** `uint64` (indexes are bit patterns; loss-free either way) |
| `LatLng` representation | `{Lat, Lng Angle}` | `{Lat, Lng float64}` degrees | match uber | degrees form forces convert-copies (DR-003) | breaking if changed | **keep** (the core safety feature) |
| `LatLng.String()` | none (fields render via `Angle.String`) | 5-digit degrees string | add `String()` | none | additive | **reject** (`fmt` already prints `{37.775939° -122.417951°}` via the field Stringer) |
| `PolygonToCellsExperimental` cap | no cap parameter | optional variadic `maxNumCellsReturn` | add variadic cap | none | additive | **keep** absence (`PolygonToCellsExperimentalSeq` stops early without a magic argument) |
| Error spelling | `ErrResolutionMismatch` | `ErrRsolutionMismatch` (sic) | match typo | none | breaking if changed | **keep** |
| Compatibility adapter package | none | n/a | `compat` package | n/a | new frozen surface | **reject** (rationale in the [migration guide](migration-from-uber-h3-go.md#why-there-is-no-compatibility-adapter-package)) |

Everything classified **keep**/**reject** above remains documented as a
migration adaptation in the [migration guide](migration-from-uber-h3-go.md);
everything **implemented** removes a row from it.

## 1. Immediate hierarchy methods — implement

Uber's implementations are one-line delegations
(`h3.go:877-879`, `899-901` in v4.4.1):

```go
func (c Cell) ImmediateParent() (Cell, error)  { return c.Parent(c.Resolution() - 1) }
func (c Cell) ImmediateChildren() ([]Cell, error) { return c.Children(c.Resolution() + 1) }
```

Edge semantics (verified against uber's vendored C):

- **resolution-0 `ImmediateParent`** → `Parent(-1)` → C `cellToParent`
  rejects `parentRes < 0` → `ErrResolutionDomain` with a zero `Cell`;
- **resolution-15 `ImmediateChildren`** → `Children(16)` → C
  `cellToChildrenSize` rejects `childRes > MAX_H3_RES` →
  `ErrResolutionDomain` with a nil slice.

This library's composed forms produce the **identical sentinel errors**
(`checkRes` catches the out-of-range resolution one layer earlier), so the
new methods are well-defined at both boundaries with no new error shape:
they return errors exactly like `Parent`/`Children` do. Both methods are
defined for valid cells; like every H3 bit-level operation (and like C),
feeding garbage bits produces garbage.

Design decisions:

- `ImmediateParent` **delegates** to `Parent(Resolution()-1)`. The parent
  kernel is a handful of bit operations; there is nothing to specialize.
- `AppendImmediateChildren` exists because *every* collection API here has
  an `Append*` form — the invariant, not an option. `ImmediateChildren`
  delegates to it.
- The children path has a real specialization opportunity: the count is
  always 7 (hexagon) or 6 (pentagon), so the generic
  `cellToChildrenSize`/iterator machinery (`_ipow`, `_zeroIndexDigits`,
  skip-digit bookkeeping) can be replaced by a loop over the ported
  `makeDirectChild` primitive (pure bit ops, digits 0..6, skipping the
  deleted K-axis digit 1 for pentagons — exactly the sequence the child
  iterator emits at `res+1`). Both forms were implemented and benchmarked
  before choosing (results in the table below; exhaustive-equality tests
  pin the specialized path to `Children(res+1)` output across every base
  cell's center descendant at every parent resolution).
- `NumImmediateChildren` and `ImmediateChildrenSeq` are **deferred**: the
  count is a documented constant-pair (7/6, so a cap-7 buffer always
  suffices for the `Append` form), and `ChildrenSeq(res+1)` already
  streams. Both are cheap to add later if demand appears.

Measured locally (darwin/arm64, Apple M1 Max, Go 1.25.4, `-count=5`;
representative medians rounded from `bench_public_test.go`):

| Benchmark | ns/op | B/op | allocs/op |
|---|---|---|---|
| `Children(res+1)` composed (hexagon res 9) | 55.4 | 64 | 1 |
| `ImmediateChildren` specialized | 41.3 | 64 | 1 |
| `AppendImmediateChildren` warm, specialized | 20.4 | 0 | 0 |

`ImmediateParent` and the composed parent call both measured about 4.2
ns/op with no allocations; it is ergonomic sugar, not a performance API.

## 2. Generic `IsValidIndex` — implement

Uber (`h3.go:176-181`, `1242-1246`) exports the constraint and the generic:

```go
Index interface { Cell | DirectedEdge | Vertex }

func IsValidIndex[T Index](index T) bool {
    return C.isValidIndex(C.H3Index(index)) == 1
}
```

The generic parameter is **cosmetic**: whatever T is, the single C
`isValidIndex` runs — `isValidCell(h) || isValidDirectedEdge(h) ||
isValidVertex(h)`. `IsValidIndex(Cell(x))` is true when x is a valid
*vertex*. It answers "is this bit pattern a structurally valid H3 index of
*any* mode", never "does this value match its Go type's mode".

This library had `IsValidIndex(raw uint64)` — same check, but every typed
call site pays a visible `uint64(...)` conversion, and DR-002's rejection
of a public umbrella `Index` **type** left no way to write user generics
over the three index types.

Decision: make `IsValidIndex` generic over an **exported constraint** that
includes the raw form:

```go
type Index interface { Cell | DirectedEdge | Vertex | uint64 | int }

func IsValidIndex[T Index](index T) bool
```

- **Not a revision of DR-002.** DR-002 rejected an umbrella *value type*
  (forces conversions, blurs modes at runtime). A union constraint cannot
  be a value type — the compiler forbids `var x Index` — so none of
  DR-002's costs apply; it is compile-time machinery only.
- **`uint64` stays in the union** deliberately, where uber has only the
  three types. The function's entire purpose is a mode-agnostic check on
  raw bits — the natural input when bits arrive from storage or the wire
  *before* the caller knows which typed form to convert to. Forcing
  `IsValidIndex(Cell(raw))` to ask "is this maybe a vertex?" would be
  worse than the conversion it removes. This also keeps the change
  **source-compatible**: every existing `IsValidIndex(uint64(idx))` call
  still compiles (T inferred as `uint64`). `int` is included only because
  an untyped call such as the previously valid `IsValidIndex(0)` defaults
  to `int` during generic inference; negative values simply validate false.
- **Naming risk addressed in docs**: the doc comment states explicitly
  that this is any-mode structural validation and points to
  `Cell.IsValid`/`DirectedEdge.IsValid`/`Vertex.IsValid` for mode-specific
  validity. (Uber's version has the same false-friend potential and does
  not warn.)
- **Zero overhead, no reflection/interfaces/`unsafe`**: conversion to the
  raw bits is compile-time-specialized. Typed and explicit-`uint64` calls
  both measured about 4.13 ns/op and 0 allocations; allocation assertions
  keep that property locked.
- The internal three-type constraint (`format.go`'s `index`) stays
  internal and separate: parse/format helpers must *return* typed values,
  so `uint64` does not belong in that type set.

A generic `IndexDigit[T]` free function was considered alongside and
**rejected**: methods are this library's one obvious form (§4).

## 3. `CoordIJ` field types — keep `int32`

Full analysis, since this is the one candidate where uber's shape is more
"Go-natural" on the surface:

- **Where the type lives.** The public `CoordIJ{I, J int32}`
  (`h3api_types.go`) *is* the type the mechanical port computes with:
  `cellToLocalIj` writes into `*CoordIJ`; `localIjToCell` reads it;
  `ijkToIj`/`ijToIjk` convert to the internal `coordIJK{I, J, K int32}`.
  There is no boundary conversion today; the two public functions
  (`CellToLocalIJ`, `LocalIJToCell`) are zero-copy pass-throughs.
- **What uber actually does with `int`.** Its `toCPtr` narrows Go `int`
  → `C.int` (int32) **silently**: `CoordIJ{I: 1 << 40}` truncates on the
  way into C and returns whatever the truncated bits mean. The `int`
  fields don't give users a larger domain — they give users a way to
  *express* out-of-domain values that then corrupt silently. C's `int` is
  32-bit on every platform H3 supports; `int32` here is the honest
  spelling of the actual domain.
- **Overflow parity.** The local-IJ algorithms overflow-check in terms of
  C `int` arithmetic (`ijkNormalizeCouldOverflow`, ported 1:1). A 64-bit
  public field type would either (a) silently truncate like uber, (b) add
  a new range-error surface C doesn't have (a deviation requiring a
  DEVIATIONS.md entry and parity-suite carve-outs), or (c) push a
  conversion+validation layer between the public type and a new internal
  `coordIJ32` — cost and drift risk with no functional gain.
- **32-bit platforms.** `int32` behaves identically everywhere; `int`
  would at least be consistent (32-bit `int` = `int32`) but only by
  accident of platform.
- **Ergonomics in practice.** Struct literals with constants
  (`CoordIJ{I: 5, J: -3}`) compile unchanged — untyped constants fit
  `int32`. The friction is limited to variables of type `int`, which need
  an explicit `int32(...)` — and that conversion is exactly where a
  thoughtful caller should consider range. No constructors/accessors are
  needed for a two-field value struct.
- **Serialization.** `encoding/json` renders `int32` and `int`
  identically; no wire impact.
- **Breaking-change calculus.** Switching to `int` is possible pre-v1 but
  buys familiarity only, at the price of one of: silent truncation, a
  behavioral deviation from C, or a permanent conversion layer. All three
  are worse than a compile-time visible `int32(...)`.

**Decision: keep.** The doc comment on `CoordIJ` now also states the
practical consequence (the full I/J domain of H3's local-IJ space is
exactly what the type can represent).

## 4. `IndexDigit` on `DirectedEdge` and `Vertex` — implement

C 4.4.0's `getIndexDigit` is mode-agnostic bit extraction; uber exposes it
on all three types through one internal generic. For a directed edge the
digits are the **origin cell's** digits (the edge number lives in mode/
reserved bits); for a vertex, the **owner cell's** digits — well-defined,
cheap, and occasionally exactly what an index-structure tool needs.
[FUTURE_WORK.md](FUTURE_WORK.md) already classified this as "trivial if a
use case shows up"; the migration matrix is the use case (uber exposes it,
and the guide previously had to say "no direct equivalent").

Implemented as two methods mirroring `Cell.IndexDigit`'s contract
(`checkRes` + ported `getIndexDigit`, `ErrResolutionDomain` for res 0 or
out of range). The generic free-function shape was rejected: this library
exposes operations as methods on the type they inspect, one obvious form.

## 5. Grouped `GridDiskDistances` — implement

[FUTURE_WORK.md item 2](FUTURE_WORK.md#2-grouped-convenience-variant-of-griddiskdistances--resolved)
recorded the full design and gated it on demand ("implement when …
migration friction from uber/h3-go shows up"). The migration guide's
reshaping snippet *is* that friction — every migrating `[][]Cell` consumer
hand-rolls the same grouping loop.

Implemented per the recorded design, layered on the flat zero-copy core:

```go
func (c Cell) GridDiskDistancesGrouped(k int) ([][]Cell, error)
```

- counting-sort partition of the flat result into **one backing array**
  plus one headers slice (3 allocations total including the flat core's
  own result storage; the scratch counts live on the stack for the
  overwhelmingly common small k);
- `result[d]` holds exactly the cells at distance d; rings have
  `len == cap` (three-index sub-slicing), so an `append` to one ring can
  never bleed into the next;
- no order within a ring is promised (matching C);
- **semantic upgrade over uber, documented**: the binding's rings retain
  `H3_NULL` zero slots in ring 0 on pentagon-affected disks; here nulls
  are pruned before grouping (the flat core already prunes), so every
  element of every ring is a real cell;
- the flat `GridDiskDistances`/`AppendGridDiskDistances` remain the
  efficient primary form; no `Append` variant of the grouped form exists
  because a `[][]Cell` result cannot be meaningfully caller-buffered —
  callers with allocation budgets use the flat form (this is the
  documented escape hatch, per the FUTURE_WORK design's own reasoning).
- Safe/Unsafe grouped variants: **deferred** until asked, as the recorded
  design recommends (the Unsafe form is already ring-ordered, making
  manual grouping nearly free).

Benchmarked against the hand-grouping loop from the migration guide (same
output shape): at k=5 the grouped API measured about 2.00 µs, 1,296 B and 3
allocations versus 2.44 µs, 3,040 B and 30 allocations for manual
append-per-ring grouping. See `bench_public_test.go`; equivalence with
uber's `GridDiskDistances`
(modulo its null-retention wart) is asserted in `interop/uberbench`'s
equivalence suite and exercised differentially in `interop/uberdiff`.

## 6. `NumIcosahedronFaces` — implement

The internal `numIcosaFaces = 20` becomes public. Named
`NumIcosahedronFaces` for consistency with `Cell.IcosahedronFaces` (this
library expands abbreviations: `MaxCellBndryVerts` → `MaxCellBoundaryVerts`
set the precedent). Face numbers are `0..NumIcosahedronFaces-1`; the
migration guide maps uber's `NumIcosaFaces` to it and drops the
"use literal 20" advice.

## 7. Kept designs — why each uber shape was not adopted

- **`DirectedEdge.Cells() (origin, destination Cell, err)`** — the tuple
  is zero-allocation and self-labeling; uber's `[]Cell{origin, dest}`
  costs an allocation and turns position into meaning. Destructuring at
  migration is a one-line change.
- **Flat `GridDisksUnsafe`** — the single flat buffer with stride
  `MaxGridDiskSize(k)` is the exact C memory contract (unpruned zero
  slots preserved), 1 allocation regardless of origin count. Uber's
  pruned `[][]Cell` hides the C layout and allocates per origin. The
  per-origin slicing recipe stays in the migration guide.
- **Method-only, one obvious form (§12-Q10)** — uber ships ~17 operations
  in both free-function and method form; the duplication doubles the
  searchable surface without adding capability. The one dual kept here is
  `LatLngToCell`/`LatLng.Cell` (both directions read naturally).
- **Validated parsing only** — `CellFromString`-style aliases that swallow
  syntax errors and return 0 would reintroduce the exact trap `ParseCell`
  exists to close. Migrating callers *delete* code (the manual `IsValid`
  check).
- **`NumCells (int64, error)`** — uber's `int` return panics on a bad
  resolution (array index out of range on its `pow7` table). An error
  return is strictly safer, and `int64` is the C-honest width.
- **`Res0Cells() []Cell`** — the operation cannot fail (constant-sized
  output, infallible C path); uber's error return never fires and every
  caller writes a dead branch.
- **`ChildPos`/`ChildAtPos` with `int64`** — C's positions are `int64_t`
  (res-15 child counts overflow int32 and, philosophically, `int`);
  `ChildAtPos` names the inverse relationship more precisely than uber's
  `ChildPosToCell`.
- **`CellBoundary` as an opaque fixed-array value** — the zero-allocation
  boundary is a headline capability (`Boundary()` does no heap work at
  all); a `[]LatLng` result costs an allocation per call by construction.
  `range b.Verts()` iterates with zero overhead, so no iterator method is
  needed either.
- **`Angle`-typed `LatLng`** — DR-003; degree/radian confusion becomes a
  compile error, and the radians-inside representation is what makes
  polygon inputs and boundary outputs zero-copy against the ported layer.
- **`uint64` index types** — H3 indexes are bit patterns; unsigned is the
  honest carrier and matches C's `uint64_t`. Valid indexes never set the
  high bit, so uber-interop conversions are loss-free both ways.
- **`ContainmentOverlappingBBox`, `ErrResolutionMismatch`,
  `MaxCellBoundaryVerts`** — Go initialism convention, a fixed typo, and
  an expanded abbreviation; all mechanical renames at migration.
- **No `InvalidH3Index`** — the zero value already is the invalid index,
  every index type documents it, and `IsValid` is the real check; a named
  constant would be a second spelling of the same idiom.
- **No variadic result cap on `PolygonToCellsExperimental`** — the
  `...Seq` form stops early at any count without a magic positional
  argument, and the `Append`+`Max*Size` pair covers pre-sizing.
- **No compatibility adapter package** — full rationale in the
  [migration guide](migration-from-uber-h3-go.md#why-there-is-no-compatibility-adapter-package);
  unchanged by this review.

## 8. Deferred

- `NumImmediateChildren`, `ImmediateChildrenSeq` (§1) — add on demand.
- Grouped Safe/Unsafe disk-distance variants (§5) — add on demand.
- `Code(err) (int, bool)` error-code accessor (§12-Q7) — no migration
  friction (uber has no equivalent); implement on first request.
- GeoJSON support — unchanged;
  [FUTURE_WORK.md item 1](FUTURE_WORK.md#1-json--geojson-support-for-geometry-types)
  governs it. Not a uber-migration concern (uber has none either).
- `DirectedEdge.Reverse` and other H3 4.5.0 additions — belong to the next
  upstream sync, not to ergonomics.

## 9. Open API questions to settle before v1.0.0

Unchanged from [FUTURE_WORK.md](FUTURE_WORK.md#release-considerations-before-v100),
plus one addition from this review:

- whether the exported `Index` constraint should gain future generic
  helpers (e.g. a generic `Resolution[T Index]`) or stay
  single-purpose — current position: stay single-purpose until a second
  compelling generic exists.
