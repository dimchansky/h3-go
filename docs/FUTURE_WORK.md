# Future work backlog

Durable context for improvements that were deliberately **not** implemented
during the public-API build-out (Phases 0–7 of
[public-api-architecture.md](./public-api-architecture.md), commits
`b2e4870`…`a69a804`, tags `v0.1.0`/`v0.2.0`). Each entry records enough
background to be actionable in a fresh session without re-deriving the
reasoning.

Standing constraints that apply to everything below:

- **Production code stays pure, safe Go** — no `unsafe`, no cgo (DR-007,
  enforced by `make check-unsafe` + depguard). Any feature that cannot be
  built under that constraint is out of scope by definition.
- **Convenience must not weaken the efficient core.** The allocation-efficient
  APIs (`Append*` forms, `iter.Seq` iterators, value-type `CellBoundary`) are
  the primary surface; new conveniences layer *on top* of them and must never
  force extra copies or allocations into the existing paths.
- **Geometry is `Angle`-based (radians inside).** `LatLng{Lat, Lng Angle}` is
  shared verbatim between the public API and the ported C-shaped layer —
  that identity is what makes polygon inputs and boundary outputs zero-copy
  (DR-003). Any serialization feature must convert at the edge, never by
  changing the internal representation.

The upstream-compatible `h3` CLI, its 170 scenarios, differential harness,
and release cross-builds are complete. Remaining CLI work is limited to
optional shell completion; it should only be added if it remains dependency
free and cannot alter upstream-compatible parsing or help.

---

## 1. JSON / GeoJSON support for geometry types

### User problem and use cases

Users routinely need to (a) persist or transmit cells and geometries through
JSON APIs, (b) render H3 geometry on maps (Leaflet, Mapbox, deck.gl) that
consume GeoJSON, and (c) ingest polygons from GeoJSON files for
`PolygonToCells`. Today the index types (`Cell`, `DirectedEdge`, `Vertex`)
marshal as canonical hex via `encoding.TextMarshaler` — that part is done —
but `LatLng`, `Angle`, `CellBoundary`, `GeoLoop`, and `GeoPolygon` have **no**
JSON support, and users must hand-roll conversions.

### Why it was deliberately excluded (v0.x decision, §12-Q5)

Any default marshaling of geometry silently embeds three contestable
conventions, and getting them wrong mis-places data on maps by design:

| Axis | This library (internal) | GeoJSON (RFC 7946) |
|---|---|---|
| Unit | radians (inside `Angle`) | degrees |
| Field order | `Lat`, `Lng` struct fields | `[longitude, latitude]` arrays |
| Ring rules | C loop order, implicit closure | outer CCW, holes CW, explicitly closed rings |

A `MarshalJSON` on `LatLng` that emits `{"lat": 0.659, "lng": -2.136}`
(radians) looks plausible and is wrong for almost every consumer; one that
emits degrees breaks round-tripping with anything that assumed the struct's
radian semantics. Refusing a default was the safe call for v0.

### Design questions and trade-offs

1. **Where does the conversion live?** On the core types (methods) or in a
   separate `geojson` sub-package? Methods are discoverable but bake one
   convention into the core forever; a sub-package keeps the core neutral and
   can be versioned/replaced independently.
2. **Geometry model**: GeoJSON `Polygon` vs `MultiPolygon` vs `Feature`
   (with the H3 index in `properties`)? Cells naturally map to `Feature`s;
   `CellsToMultiPolygon` output naturally maps to `MultiPolygon`.
3. **Ring orientation and closure**: H3 boundaries are CCW and unclosed;
   GeoJSON wants closed rings and prescribes winding. The encoder must append
   the closing vertex and (for holes from `CellsToMultiPolygon`) verify or
   fix winding.
4. **Antimeridian**: RFC 7946 §3.1.9 says geometries SHOULD be cut at the
   antimeridian. Cells crossing ±180° (and the poles) need either cutting
   (correct, complex) or documented non-cutting (what most H3 bindings do).
5. **Precision**: how many decimal digits to emit (GeoJSON convention is ~6–7;
   full float64 round-trip needs 17).

### Possible API shapes

Recommended shape — a leaf sub-package, core types untouched:

```go
// package h3/geojson (own file tree, still zero external dependencies:
// encoding/json only)
func CellBoundaryPolygon(b *h3.CellBoundary) Polygon          // closed ring, degrees
func MultiPolygon(polys []h3.GeoPolygon) MultiPolygonGeom     // from CellsToMultiPolygon
func CellFeature(c h3.Cell) (Feature, error)                  // boundary + {"h3": "89283..."}
func ParsePolygon(data []byte) (h3.GeoPolygon, error)         // degrees -> Angle at the edge

// Explicit, unambiguous DTO — no magic marshaling on h3.LatLng itself:
type Position [2]float64 // [lng, lat], degrees — the GeoJSON convention, named so
```

Rejected shape (do not resurrect without revisiting §12-Q5):
`func (ll LatLng) MarshalJSON()` on the core type — whichever unit/order it
picks, it is a silent trap for the other half of users, and it freezes the
choice into the core package's compatibility surface.

An intermediate option if a sub-package feels heavy: explicit conversion
methods on core types (`ll.GeoJSONPosition() [2]float64`) that make the
convention visible at the call site, still without `MarshalJSON`.

### Allocation / ownership / concurrency / safety

- Marshaling inherently allocates (the JSON bytes); that is (A)-class,
  acceptable. Conversion `Angle→degrees` is per-scalar, no copies of backing
  arrays needed except the output structures the user asked for.
- `ParsePolygon` should build the `GeoLoop`/`Holes` slices exactly once at
  the right size; the result is owned by the caller, no retained state.
- Pure Go throughout (`encoding/json`); no concurrency concerns (stateless
  functions).

### Evidence needed before building

User demand, not profiling: issues/requests for GeoJSON, or friction writing
the conversion by hand. Prototype against a real consumer (Leaflet/deck.gl
rendering of `CellsToMultiPolygon` output) to validate winding/closure
decisions before freezing names.

### Tests and benchmarks

- Round-trip: `ParsePolygon(Marshal(MultiPolygon(x)))` set-equality on
  `PolygonToCells` results.
- Golden GeoJSON fixtures validated against an external checker (e.g.
  `geojson.io`-conformant linter output committed as testdata).
- Antimeridian/pole cells (transmeridian fixtures already exist in the
  ported tests — reuse them).
- Winding tests for holes from donut `CellsToMultiPolygon` output.
- Allocation budget test: encoding N cells allocates O(output), not O(N²).

### Compatibility

Fully backward-compatible if shipped as a sub-package (pure addition). This
is also **the main open decision to settle before v1.0.0** — not because v1
must include it, but because v1 must commit to *not* having marshaling on the
core geometry types (aliasing that decision into the frozen surface).

### Recommended direction

Leaf sub-package `geojson` with explicit DTO types and edge conversion;
`Feature`-oriented helpers for cells; document non-cutting at the
antimeridian in v1 of the sub-package (cutting as a later option). Do not
add `MarshalJSON` to core geometry types.

---

## 2. Grouped convenience variant of GridDiskDistances — resolved

Implemented as `Cell.GridDiskDistancesGrouped` after the uber/h3-go
migration guide demonstrated repeated reshaping friction. The flat
`GridDiskDistances`/`AppendGridDiskDistances` APIs remain the efficient core;
the grouped form partitions their cell backing array in place, returns
len-limited ring slices, prunes pentagon null holes, and uses three
allocations for common radii. Grouped Safe/Unsafe variants remain deferred
until a distinct use case appears. The full decision and measurements are in
[public-api-ergonomics-review.md](public-api-ergonomics-review.md#5-grouped-griddiskdistances--implement).

---

## 3. Reusable scratch/workspace buffers for polygon operations

### User problem and use cases

High-throughput polyfill services (tiling pipelines, geofencing engines)
call `PolygonToCells`/`AppendPolygonToCells` in tight loops. The *result*
buffer is already reusable via `AppendPolygonToCells`, but each call still
performs algorithm-internal allocations. Measured at v0.2.0
(`BenchmarkAppendPolygonToCells`, 1253-cell SF polygon, warm result buffer):
**3 allocations / ~152 KB per call** — the search/found hash arrays and the
per-polygon bbox slice inside the ported `polygonToCells` (the C original
heap-allocates the same arrays). `AppendCompactCells` similarly retains 3
internal working arrays (again matching C's mallocs), and
`polygonToCells` also uses a small per-search ring buffer that C keeps on
the stack.

### Why it was deliberately excluded

- The allocations mirror the C implementation's own mallocs — the port is
  not *worse* than C; there is no regression to fix.
- A workspace API adds real costs: lifecycle rules (who resets it?),
  concurrency rules (one workspace per goroutine, never shared), retained-
  memory growth (a workspace sized by the largest polygon ever seen pins
  that memory), and permanent API surface. Per the architecture's §5.7 this
  was explicitly deferred: "could later be offered a scratch buffer via an
  options struct if profiling demands it — explicitly out of scope now."
- Amortized over the ~1 ms compute of a 1253-cell fill, 3 allocations are
  noise for most users; only sustained-throughput workloads could tell the
  difference, and none has been measured yet.

### Design questions and trade-offs

1. **Workspace object vs pool vs variadic option.** A `sync.Pool` hidden
   inside the package trades determinism for convenience and *retains
   memory invisibly* — poor fit for a library that advertises explicit
   allocation control. An explicit workspace type keeps ownership visible.
2. **Scope**: polygon fill only, or also compact/uncompact and the gridDisk
   internal distance scratch? Start narrow (polygon fill is the heavyweight).
3. **Interaction with ported code**: threading a workspace through
   `polygonToCells` means touching ported function signatures — a
   traceability cost. The workspace parameters must be *additive* shims
   (e.g. a thin variant that the attribution-carrying function delegates
   to), never a rewrite of the C-shaped body.
4. **Growth policy**: grow-only (simple, pins peak memory) vs shrink
   heuristics (complex). Grow-only with a documented `Reset()`/re-make
   escape hatch is the honest option.

### Possible API shape (sketch only — do not commit prematurely)

```go
// PolyfillWorkspace holds reusable internal buffers for polygon-to-cells
// operations. Not safe for concurrent use; use one per goroutine.
// The zero value is ready to use. Buffers grow to the largest polygon
// processed and are retained until the workspace is garbage collected.
type PolyfillWorkspace struct {
    search, found []Cell  // internal hash arrays
    bboxes        []bbox  // note: internal type — fields stay unexported
}

func (w *PolyfillWorkspace) AppendPolygonToCells(dst []Cell, p GeoPolygon, res int) ([]Cell, error)
```

Alternative considered: `AppendPolygonToCells(dst, p, res, WithWorkspace(w))`
variadic options — heavier machinery for one option; only worth it if more
options (e.g. flags) accumulate.

### Allocation / copying / concurrency / ownership / safety

- Target: 0 allocations on the warm path (result buffer via `dst`, internals
  via workspace). No copying changes; no `unsafe` needed — the buffers are
  ordinary slices of internal types.
- **Not** concurrency-safe by design; the doc comment must say "one per
  goroutine" and tests should include a `-race` misuse test only if we decide
  to detect misuse (probably not — document instead).
- Ownership: caller owns the workspace; the library never retains a
  reference past the call (must be asserted in review — retaining would be a
  correctness bug, not just a perf one).

### Evidence needed before building (the gate)

A profile of a realistic sustained workload (thousands of polyfills/second)
showing GC pressure or allocation time from *these specific* buffers —
e.g. `pprof` alloc_space attributing a meaningful share to
`polygonToCells`'s `make` calls, or benchmark deltas >5–10% from a
prototype. Without that, the API-complexity cost wins. The cheap first step
if polyfill throughput ever matters: move the small per-search ring buffer
(`maxOneRingSize`) to a stack array — C keeps it on the stack, it is
bounded, and it needs no API change at all.

### Tests and benchmarks

- Alloc assertion: warm workspace path == 0 allocs (AllocsPerRun).
- Equivalence: workspace results byte-identical to the plain path across the
  ported polygon test corpus (SF polygon, holes, transmeridian, pentagon
  fixtures).
- Reuse across differently-sized polygons (grow, then smaller input).
- Benchmark: plain vs workspace at several polygon sizes; report B/op.

### Compatibility

Purely additive; backward-compatible. The existing `AppendPolygonToCells`
stays untouched.

### Recommended direction

Do nothing until profiling evidence exists. If it arrives: explicit
grow-only workspace struct (zero value usable, per-goroutine), polygon fill
first, stack-ring micro-fix independently of the workspace. Reject the
hidden `sync.Pool`.

---

## Other follow-ups discovered during Phases 0–7

### Planned / promising (no blockers, do when convenient)

- **Report the `gridRing` negative-k UB upstream.** C `gridRing` executes
  `memset(out, 0, 6*k*sizeof(H3Index))` after `gridRingUnsafe` rejects
  k < 0 — a negative size wrapped to a huge `size_t` (no-op under macOS's
  end-pointer memset, segfault under glibc). Found by this repo's parity
  suite on Linux CI; details in DEVIATIONS.md item 9. The Go port is
  unaffected (bounded slice zeroing).
- **Report the `maxPolygonToCellsSizeExperimental` timeout pathology
  upstream.** On loops with huge-magnitude latitudes (~1e287 rad — in-domain
  for upstream's own `fuzzerPolygonToCellsExperimental`, which feeds raw
  doubles as radians), the estimator's rough bbox area
  `height * width / cos(min(|north|, |south|)) * R²` comes out *negative*
  (cosine of a huge angle can be negative), which defeats the
  resolution-coarsening loop, so the size estimate scans every cell of the
  planet at the requested resolution. One such input costs ~27 s in H3 C
  4.4.0 (~15 s in this port) — past OSS-Fuzz's default 25 s single-input
  budget. Found by this repo's fuzz port and analyzed in issue #3; the
  reproducer and full analysis live in
  [testdata/fuzz-findings/README.md](../testdata/fuzz-findings/README.md).
  A `fabs()` around the area (or clamping negative estimates) restores the
  coarsening and would bound the scan; behavior parity keeps this port
  matching C until upstream decides.
- **Repository publication and release polish** — largely **done** for
  v0.3.0: the complete release procedure lives in
  [releasing.md](releasing.md) (reproducible archives published as
  immutable GitHub Release assets, notes per [versioning.md](versioning.md),
  CLI docs linking the Releases page); the Code of Conduct, issue
  templates, and CODEOWNERS are in place. Remaining, inherently
  post-public items (tracked in releasing.md and the release checklist):
  homepage → pkg.go.dev after the new version renders there; anonymous
  badge verification; Go Report Card badge once the service can access the
  repository; the social-preview image upload.
  - Evaluate **OpenSSF Scorecard** (workflow + badge) as a post-public
    repository-security improvement: it grades pinned dependencies, token
    permissions, branch protection, and fuzzing — most of which this
    repository already satisfies — at the cost of one more scheduled
    workflow. Decide once the repository has been public for a while.
  - **Developer ID signing and Apple notarization for the macOS release
    assets.** Today the darwin binaries are ad-hoc linker-signed but not
    Developer ID-signed or notarized, so Gatekeeper can reject
    browser-downloaded copies (the archive README documents the safe
    per-binary workaround). Doing this properly requires: an Apple
    Developer Program membership; secure handling of the signing and
    notarization credentials; signing + notarization steps in the release
    pipeline; Gatekeeper testing on clean macOS systems; and a deliberate
    design for the current byte-for-byte reproducibility guarantee —
    signatures embed secure timestamps, so the design must either keep
    reproducible unsigned build artifacts alongside separately signed
    distribution artifacts, or redefine (and re-verify) reproducibility at
    the pre-signing stage while treating signed assets as
    provenance-attested outputs. Which design wins is an open decision for
    that future work. References:
    [Developer ID](https://developer.apple.com/support/developer-id/),
    [notarizing macOS software](https://developer.apple.com/documentation/security/notarizing-macos-software-before-distribution),
    [Gatekeeper/user guidance](https://support.apple.com/en-us/102445).
- **uberdiff extensions**: the benchmark-comparison half is **done** —
  [interop/uberbench](../interop/uberbench/README.md) benchmarks every
  operation category against the binding with equivalence gating, plus
  process-level memory probes; results under
  [docs/benchmarks/](benchmarks/README.md), matrix and docs under
  [docs/comparison-uber-h3-go.md](comparison-uber-h3-go.md). Remaining:
  absorb Uber's unreleased pure-Go port (`x/h3go`, present in uber/h3-go
  master, still untagged as of 2026-07) into the same differential and
  benchmark harnesses — cgo-free — once it ships.
- **Error-code accessor (§12-Q7)**: `func Code(err error) (int, bool)`
  returning the numeric H3 code for cross-language stability. Additive;
  implement on first request.

### Ideas requiring profiling or user demand (gated)

- Items 1–3 above (GeoJSON: demand; grouped distances: demand; workspaces:
  profiling).
- **Stack ring buffer in `polygonToCells`** (see item 3) — smallest possible
  perf follow-up, no API change, C-fidelity improvement.
- **Filename convention cleanup (§12-Q12)**: 49 of the original 75 C-public
  functions live in double-underscore files
  (e.g. `algos__gridDisk.go`). The apiinventory tool already compensates, so
  this is cosmetic churn — only worth doing in a quiet moment, as a pure
  `git mv` commit with no content changes. New ports (e.g. the 4.4.0 trio)
  already follow the correct convention. *Partially done*: the lone casing
  outlier (`h3Index_getBaseCellNumber.go` → `h3index_…`) was renamed in the
  DR-008 phase-3 commit
  ([repository-layout-review.md](repository-layout-review.md) §8); the mass
  `__` renames remain deliberately not worth it (ibid. §5-A).
- **Internal tidy**: the ported `h3ToString` returns `(string, uint32)` and
  uses `fmt.Sprintf` — unused by the public path (which formats via
  `strconv`); could be aligned with the C signature during a future sync.
- **`AppendDirectedEdges` / `AppendIcosahedronFaces`**: `Cell.AppendVertexes`
  landed in 2026-07 after an allocation study of the `Cell.Vertexes` path
  (warm-buffer reuse was the only measurable win — 0 allocs/op; exact
  sizing and staging-copy removal measured as noise). The other two
  fixed-size collection APIs can adopt the same pattern for family
  symmetry. Gated on demand: each is a copy of the same ≤6-element
  template, but each also grows the locked API surface, and the study
  showed the win is confined to warm-buffer loops (48 B, 1 alloc per call
  otherwise). Implement on first request or first profile that shows
  either API in a hot loop.

### Intentionally rejected designs (do not revisit without new evidence)

Recorded here so they are not re-proposed from scratch; full rationale in
[public-api-architecture.md](./public-api-architecture.md) (decision records)
and [DEVIATIONS.md](./DEVIATIONS.md):

- `internal/` package split for the ported layer (DR-001: breaks the
  zero-copy alias, methods, and white-box parity tests). Re-investigated
  in depth 2026-07 against a full `internal/h3` proposal and rejected
  again as **DR-008** with probe evidence (methods cannot attach to
  aliases of non-local types; `[]Cell` is not convertible to a
  differently defined index slice — measured 3.7–4.8× boundary tax plus a
  warm-path allocation; shared-type packages import-cycle; the alias
  facade hides every method from godoc while relocating, not separating,
  the mix): see
  [repository-layout-review.md](repository-layout-review.md). Layer
  discoverability is handled instead by `make check-layout` and the
  generated [file-layer-inventory.csv](file-layer-inventory.csv).
- `unsafe` slice reinterpretation (`castSlice`) — obsoleted by the
  `h3Index = Cell` alias; DR-007 requires a new reviewed decision record,
  benchmarks, and proof no safe design suffices before any production
  `unsafe`.
- Public umbrella `Index` value type (DR-002) — forces conversions
  everywhere. The exported `Index` *constraint* used by generic
  `IsValidIndex` is compile-time-only and does not create that problem.
- Degrees-based `float64` `LatLng` (DR-003/§12-Q4) — would force O(n)
  convert-copies on every polygon/boundary crossing the API boundary.
- Dual package-function + method forms for every operation (§12-Q10) — one
  obvious form each; uber's duplication was judged a wart.
- `MarshalJSON` directly on core geometry types (§12-Q5; see item 1).
- Sorting/ordering guarantees beyond C's documented contracts (§12-Q11).

### Release considerations before v1.0.0

v1 freezes the public surface (`docs/api-surface.txt` is the inventory of
what gets frozen). Versioning mechanics — tag format, the independent H3
Core target, release titles, and the release-note outline — are fixed in
[versioning.md](versioning.md). Before tagging:

1. CI green on GitHub (all jobs), including the api-gates job against a
   fresh upstream download.
2. Settle item 1's *boundary* decision — specifically, commit to "no JSON
   marshaling on core geometry types" (the sub-package can come later).
3. Decide whether `Code(err)` (§12-Q7) is in or out — additive either way,
   but cheap to include if wanted.
4. Ideally ride through one more upstream release (4.5.x) with the §10 sync
   workflow to confirm the 4.4.0 rehearsal wasn't a one-off; the
   `getIndexDigit` macro/function collision showed each sync can surface a
   naming surprise.
5. Re-run the full matrix on both supported Go versions and refresh the
   benchmark numbers in the README if they are quoted there.
