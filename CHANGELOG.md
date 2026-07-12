# Changelog

All notable changes to this project are documented here. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/); versions follow
[Semantic Versioning](https://semver.org). Pre-1.0, minor versions may
contain breaking changes (called out explicitly below).

## [Unreleased]

- Evidence-based comparison with the official cgo binding uber/h3-go:
  a function-by-function coverage matrix (`docs/comparison-uber-h3-go.md`,
  generated from a curated CSV by the new `tools/ubercompare`, drift-gated
  in CI), a migration guide (`docs/migration-from-uber-h3-go.md`, with its
  before/after example kept executable by a test), and a new comparative
  benchmark module `interop/uberbench` — equivalence-gated benchmarks
  across scalar/geometry/collection/batch workloads plus a process-level
  memory probe (`cmd/memprobe`), with committed per-environment artifacts
  under `docs/benchmarks/` and a manual/monthly `benchmarks` workflow.
  The README's "Why this library" and "Performance" sections now cite
  these measured results. The uberdiff differential module was re-pinned
  from uber/h3-go v4.2.2 to v4.4.1 (vendoring H3 C v4.4.1, behaviorally
  identical to this library's v4.4.0 target).
- CI: the Nightly workflow (schedule, manual dispatch, and `v*` release
  tags) now replays the full upstream fixture suites — 526,546 golden
  conversion/boundary records — as a step of the `parity` job, reusing its
  downloaded reference tree.
- Maintainer guidance is now tracked: a one-page `AGENTS.md` quick reference
  (with `CLAUDE.md` as a pointer to it) replaces the previously gitignored
  local agent files; the durable rules they carried live in
  `CONTRIBUTING.md` (`__` static-helper naming, commit-trailer policy).
- Lint policy documented in `docs/lint-policy.md`: the style-tier
  gocritic/revive exclusions are now path-scoped to mechanically ported
  files instead of disabled globally (new idiomatic code gets the full rule
  set), and two stale global exclusions (staticcheck SA1019, govet
  fieldalignment-in-tests) were removed after verifying nothing fires.

- Documentation overhaul: rewritten README (repository map, CLI section,
  per-audience documentation paths), new indexes `docs/README.md` and
  `tools/README.md`, READMEs for `cmd/h3` and `interop/uberdiff`, package
  documentation and compatibility-invariant comments throughout
  `internal/cli`, and `-h` usage summaries for all tools. Corrected stale
  statements (test-inventory dispositions, DEVIATIONS parity target,
  architecture-document status, two broken README anchors).
- New `tools/docscheck` Markdown link/anchor checker with `make check-docs`;
  CI now runs it on every push/PR — including docs-only changes, which
  previously ran no checks at all.

- Added a pure-Go `h3` executable compatible with all 63 commands and all 170
  registered CLI scenarios from H3 C v4.4.0, including file/stdin workflows,
  JSON/WKT/newline formats, upstream exit codes, process tests, and opt-in
  differential tests against `h3_bin`.
- Added semantic CLI, scenario, fixture, and defining-source inventories plus
  a strict upstream drift gate; CI now builds/tests the CLI and nightly/tag
  validation performs C differential and cross-platform builds.

- **4.4.0 sync completion**: the v0.2.0 sync ported the implementation delta
  but missed the upstream *test* delta. Now closed: ported the new
  `testConstructCell.c` and `testIndexDigits.c` suites and the 4.4.0
  additions to five existing test files; full symbol-level audit recorded in
  `docs/sync/4.3.0-to-4.4.0.md` (all implementation changes verified —
  the remaining 4.4.0 code changes were release-behavior no-ops).
- New `tools/upstreamdiff` (`make upstream-diff FROM=... TO=...`): symbol-level
  C tree diff mapped to the Go port via attribution comments; upstream-sync
  documentation now makes this review mandatory.
- CI: tiered pipeline (docs/ci-policy.md) — docs-only changes skip Go jobs,
  fast checks signal in ~3 minutes, race runs on PRs/nightly/tags, heavy
  suites nightly/on-demand; in-flight runs are cancelled by newer commits.
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
