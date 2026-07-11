# Changelog

All notable changes to this project are documented here. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/); versions follow
[Semantic Versioning](https://semver.org). Pre-1.0, minor versions may
contain breaking changes (called out explicitly below).

## [Unreleased]

- CI: version-agnostic parity/include paths (no H3 version hardcoded in
  code); platform-independent allocation assertions; Linux parity build
  fixes (portable `C.int64_t`, `-lm`, no section-GC linker flags) and
  platform-deterministic parity tests (avoids C UB in `gridRing(k<0)` and
  `clz(0)`, zeroed C output buffers). Discovered upstream-reportable UB in
  C `gridRing` with negative k.
- Public-readiness documentation: README overhaul, CONTRIBUTING, SECURITY,
  changelog.

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

[Unreleased]: https://github.com/dimchansky/h3-go/compare/v0.2.0...HEAD
[v0.2.0]: https://github.com/dimchansky/h3-go/compare/v0.1.0...v0.2.0
[v0.1.0]: https://github.com/dimchansky/h3-go/releases/tag/v0.1.0
