# Changelog

All notable changes to this project are documented here. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/); versions follow
[Semantic Versioning](https://semver.org). Pre-1.0, minor versions may
contain breaking changes (called out explicitly below). The module version
and the H3 Core compatibility target are independent axes; the full
versioning and release policy is [docs/versioning.md](docs/versioning.md).

## [Unreleased]

Migration of the behavioral compatibility target from **H3 C v4.4.0** to
**H3 C v4.5.0**, executed per the reviewed
[sync record](docs/sync/4.4.0-to-4.5.0.md) (discovery issue #25,
implementation issues #27–#36). `VersionMajor/Minor/Patch` now report
4.5.0.

### Added

- `DirectedEdge.Reverse() (DirectedEdge, error)` — the H3 4.5.0
  `reverseDirectedEdge` function (79 public C functions are now covered,
  up from 78). Never allocates; matches the uber/h3-go v4.5.0 method
  shape.

### Changed

- **`CellsToMultiPolygon` / the multipolygon pipeline** now implements H3
  4.5.0's arc-cancellation algorithm (upstream replaced the vertex-graph
  algorithm). Observable differences, inherited from upstream: loop
  start-vertex/extraction order differs from 4.4.0; invalid cells now
  fail with `ErrCellInvalid` (previously `ErrFailed`); duplicate inputs
  fail with `ErrDuplicateInput`; the full-globe input returns the eight
  octant polygons. The public contract documents guaranteed error
  ordering.
- **Internal Vec3 indexing refactor** (upstream 4.5.0): `latLngToCell`
  and related indexing paths route through 3D-vector geometry; discrete
  results are unchanged and parity-locked, floating-point behavior is
  compared against the C oracle with measured tolerances.
- **CLI at the 4.5.0 contract** (172 scenarios, up from 170): the
  cell-input scanner adopts the 4.5.0 stale-buffer fix (no more phantom
  cells at 1500-byte chunk boundaries — upstream's uber/h3#1124/#1125),
  `edgeLengthM` prints at `%.8f` (was `%.10f`), and the two new
  `readCellsFromFile` scenarios pin the fixed scanner (127 cells). The
  `--version` banner reports `h3 4.5.0`.
- `gridPathCells` adopts upstream 4.5.0's bidirectional interpolation
  (endpoint-exact paths; ported and parity-locked), and
  `destroyLinkedMultiPolygon` semantics (idempotent head-zeroing) are
  mirrored in the internal linked-geo lifecycle.
- Interop/comparison baseline moved to **uber/h3-go v4.5.0** (was
  v4.4.1): differential and equivalence suites rerun green; the
  comparison matrix gains the `reverseDirectedEdge` row. Committed
  benchmark results remain dated snapshots of the v4.4.0-vs-v4.4.1 runs
  until the suite is rerun.

### Compatibility

- No Go API removals or renames; the only signature-visible addition is
  `DirectedEdge.Reverse`. Behavior differences from v0.3.0 are exactly
  the upstream 4.5.0 behavioral changes listed above.
- Deviations from C remain limited to those listed in
  [docs/DEVIATIONS.md](docs/DEVIATIONS.md); the cgo parity suite
  (213 files) now compares against pristine H3 C v4.5.0 sources.

## [v0.3.0] — 2026-07-15

Feature release and first public release. The behavioral compatibility
target is unchanged: **H3 C v4.4.0**. No breaking Go API changes — every
API addition below is purely additive.

### Added

- Zero-allocation `Cell.AppendVertexes(dst []Vertex)`: appends the cell's
  6 (hexagon) or 5 (pentagon) topological vertexes to a caller-owned
  buffer with true append semantics; a capacity of 6 always suffices for
  an allocation-free warm path, and on error `dst` is returned unchanged.
  `Cell.Vertexes` now delegates to it with identical results, errors, and
  allocation profile (one 48 B result slice); an allocation study recorded
  warm reuse as the only observable win, and a post-change interleaved
  benchmark comparison confirmed the delegation is within noise.
- API ergonomics additions: generic typed `IsValidIndex`,
  `DirectedEdge.IndexDigit`, `Vertex.IndexDigit`, `Cell.ImmediateParent`,
  `Cell.ImmediateChildren`, zero-allocation `Cell.AppendImmediateChildren`,
  `Cell.GridDiskDistancesGrouped`, and `NumIcosahedronFaces` — covered by
  public tests, allocation assertions, equivalence checks against
  uber/h3-go, runnable examples, and focused benchmarks.
- A pure-Go `h3` executable compatible with all 63 commands and all 170
  registered CLI scenarios from H3 C v4.4.0 — file/stdin workflows,
  JSON/WKT/newline formats, upstream exit codes — plus semantic CLI,
  scenario, fixture, and defining-source inventories with a strict
  upstream drift gate, opt-in differential tests against the C `h3_bin`,
  and cross-platform builds.
- Evidence-based comparison with the official cgo binding uber/h3-go: a
  function-by-function coverage matrix (`docs/comparison-uber-h3-go.md`),
  a migration guide (`docs/migration-from-uber-h3-go.md`, its before/after
  example kept executable by a test), a comparative benchmark module
  `interop/uberbench` with equivalence-gated benchmarks and a
  process-level memory probe (`cmd/memprobe`), committed per-environment
  artifacts under `docs/benchmarks/`, and a manual/monthly `benchmarks`
  workflow.
- Validation depth: the Nightly workflow (schedule, manual dispatch, and
  `v*` release tags) replays the full upstream fixture suites — 526,546
  golden conversion/boundary records — and the fuzz rotation grew from
  three to six targets (the seventh is seed-corpus-only pending an
  investigation of pathologically slow in-domain polygon inputs, tracked
  in issue #3 with a preserved reproducer under
  `testdata/fuzz-findings/`).
- Release engineering for the first public release: a single authoritative
  release builder (`tools/releasepack`, `make release-dist`) producing
  bit-reproducible archives verified by an independent CI rebuild and
  runtime-smoked on architecture-proven runners; an aggregate
  `CI / required` merge gate (`tools/cirequired`) with an explicit
  truth table; pinned release toolchain and SHA-pinned actions; pinned
  vulnerability (`govulncheck`) and full-history secret (`gitleaks`) scans
  in the tag-triggered release gate; Dependabot for actions and interop
  modules; a standalone README shipped inside every release archive.
- Maintenance tooling and gates: `tools/docscheck` Markdown link/anchor
  checker (`make check-docs`, run on every push/PR including docs-only);
  `tools/upstreamdiff` symbol-level C-tree diff (mandatory in upstream
  syncs); `tools/layoutinventory` with the freshness-gated
  `docs/file-layer-inventory.csv` and `make check-layout`;
  `tools/benchdocs` and `tools/ubercompare` generating and drift-gating
  the performance and comparison documentation.

### Changed

- Documentation overhaul: rewritten README (repository map, CLI section,
  per-audience documentation paths), new `docs/README.md` and
  `tools/README.md` indexes, READMEs for `cmd/h3` and `interop/uberdiff`,
  package documentation throughout `internal/cli`, and corrected stale
  statements (test-inventory dispositions, DEVIATIONS parity target,
  architecture-document status, two broken README anchors). Maintainer
  guidance is tracked in `AGENTS.md`/`CLAUDE.md`; public-readiness
  documents (CONTRIBUTING, SECURITY, this changelog, `docs/releasing.md`,
  issue templates, `CODE_OF_CONDUCT.md`) are in place.
- Performance documentation presents separate Apple M1 Max darwin/arm64
  and shared-runner linux/amd64 excerpts (including platform reversals and
  Go-heap allocation metrics), generated and drift-gated by
  `tools/benchdocs`. The complete darwin/arm64 artifacts were re-measured
  with Go 1.26.5; linux/amd64 measurements are unchanged.
- CI: tiered pipeline (`docs/ci-policy.md`) — docs-only changes skip Go
  jobs, fast checks signal in ~3 minutes, race runs on PRs/nightly/tags,
  heavy suites nightly/on-demand; version-agnostic parity include paths
  (no H3 version hardcoded in code); in-flight runs cancelled by newer
  commits.
- Lint policy documented in `docs/lint-policy.md`: style-tier
  gocritic/revive exclusions are path-scoped to the mechanically ported
  tier instead of disabled globally; two stale global exclusions were
  removed.
- The uberdiff differential module was re-pinned from uber/h3-go v4.2.2 to
  v4.4.1 (vendoring H3 C v4.4.1, behaviorally identical to this library's
  v4.4.0 target).
- Repository layout review (DR-008, `docs/repository-layout-review.md`):
  the flat single-package layout reaffirmed with compile-probe and
  benchmark evidence; the filename-casing outlier
  (`h3Index_getBaseCellNumber.go`) renamed to the standard `h3index_`
  prefix. No API or behavior changes.
- `go.mod` now retracts the stale 2019 pseudo-version
  `v0.0.0-20191115084349-4ab39a97d2af` — a pre-rewrite codebase previously
  cached by the Go module proxy that is not part of this line; `go`
  tooling will warn anyone still pinned to it.
- The project copyright (`Copyright 2019-2026 Dmitrij Koniajev`,
  reflecting the project's 2019 origin) is now recorded consistently in
  `LICENSE` (the Apache License 2.0 application notice) and in `NOTICE`,
  which opens with the project notice ahead of the clearly distinguished
  Uber Technologies H3 attribution.

### Fixed

- **4.4.0 sync completion**: the v0.2.0 sync ported the implementation
  delta but missed the upstream *test* delta. Now closed: ported the new
  `testConstructCell.c` and `testIndexDigits.c` suites and the 4.4.0
  additions to five existing test files; full symbol-level audit in
  `docs/sync/4.3.0-to-4.4.0.md` (all remaining 4.4.0 implementation
  changes verified as release-behavior no-ops).
- Linux parity build fixes (portable `C.int64_t`, `-lm`, no section-GC
  linker flags) and platform-deterministic parity tests — avoiding C
  undefined behavior in `gridRing(k<0)` and `clz(0)` with zeroed C output
  buffers; the `gridRing` negative-k UB is an upstream-reportable finding.

## [v0.2.0] — 2026-07-11

Upstream sync: **H3 C 4.3.0 → 4.4.0** (first execution of the documented
sync workflow).

### Added
- `Cell.IndexDigit`, `ConstructCell`, `IsValidIndex` — the three functions
  added in H3 C 4.4.0.
- Error sentinels for the new 4.4.0 codes: `ErrIndexInvalid`,
  `ErrBaseCellDomain`, `ErrDigitDomain`, `ErrDeletedDigit`.

### Changed
- Parity target and `Version*` constants now report H3 4.4.0.
- `describeH3Error` (internal) follows the 4.4.0 bounds-check form.

## [v0.1.0] — 2026-07-11

First complete public API. All 75 public functions of H3 C v4.3.0 are
represented through an idiomatic, strongly typed Go surface:

### Added
- Typed indexes `Cell`, `DirectedEdge`, `Vertex` (distinct `uint64` types)
  with `String`, validated `ParseCell`/`ParseDirectedEdge`/`ParseVertex`,
  and `encoding.TextMarshaler`/`TextUnmarshaler` (JSON-ready).
- `Angle`-based geographic types (`LatLng`, `CellBoundary`); degree/radian
  confusion is uncompilable.
- Error sentinels mirroring the H3 C error codes, matched with `errors.Is`.
- Full API surface: indexing/inspection, hierarchy (parents, children,
  compaction), traversal (grid disks/rings/paths, local IJ), directed edges,
  vertexes, regions (polyfill, multipolygon outlines), and measurement
  (areas, lengths, great-circle distances).
- Zero-allocation `Append*` variants for every collection API and
  `iter.Seq` iterators (`Cell.ChildrenSeq`, `CellsAtRes`,
  `PolygonToCellsExperimentalSeq`).
- Structural guarantees, all CI-enforced: pure safe Go (no cgo, no `unsafe`
  in any normal build), C-API completeness gate, golden API-surface lock.
- Validation: 227-file cgo parity suite against pristine upstream C sources,
  ported upstream unit tests, fuzz targets, allocation assertions, and a
  differential suite against the official uber/h3-go binding.

[Unreleased]: https://github.com/dimchansky/h3-go/compare/v0.3.0...HEAD
[v0.3.0]: https://github.com/dimchansky/h3-go/compare/v0.2.0...v0.3.0
[v0.2.0]: https://github.com/dimchansky/h3-go/compare/v0.1.0...v0.2.0
[v0.1.0]: https://github.com/dimchansky/h3-go/releases/tag/v0.1.0
