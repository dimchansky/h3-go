# Maintenance tools

Small, dependency-free `package main` commands that keep the port, its
inventories, and its documentation verifiably in sync with upstream H3 C.
None of them ship to library or CLI users. Each has a package comment with
full details (`go doc ./tools/<name>`) and accepts `-h`.

| Tool | Purpose | Typical command | Used by |
|---|---|---|---|
| [apiinventory](apiinventory) | Map the H3 C public API to the Go port; verify completeness | `make api-inventory` / `make check-api` | CI (`api-gates`, nightly), upstream syncs |
| [testinventory](testinventory) | Verify every upstream test-ecosystem entry has a reviewed disposition | `make check-test-inventory` | CI (`api-gates`, nightly), upstream syncs |
| [cliinventory](cliinventory) | Discover the upstream CLI contract; verify the committed CLI registries against it | `make check-cli-inventory` | CI (`api-gates`, nightly), upstream syncs |
| [upstreamdiff](upstreamdiff) | Symbol-level diff of two upstream H3 trees, mapped to the Go port | `make upstream-diff FROM=4.3.0 TO=4.4.0` | Upstream syncs (manual, mandatory) |
| [docscheck](docscheck) | Verify relative Markdown links and #anchors | `make check-docs` | CI (`docs` job) |
| [benchdocs](benchdocs) | Generate and verify the README scorecard and complete benchmark comparison from committed artifacts | `make gen-benchdocs` / `make check-benchdocs` | CI (`docs` job), benchmark refreshes |
| [ubercompare](ubercompare) | Generate and verify the uber/h3-go comparison-matrix tables | `make gen-ubercompare` / `make check-ubercompare` | CI (`docs` job), binding/H3 release updates |
| [layoutinventory](layoutinventory) | Classify every root source file into its architectural layer; verify none is unclassifiable | `go run ./tools/layoutinventory > docs/file-layer-inventory.csv` | Layout discoverability ([docs/repository-layout-review.md](../docs/repository-layout-review.md)) |
| [unexport](unexport) | **Historical** one-time migration sweep (Phase 2 unexport) | `go run ./tools/unexport` (dry run) | Nothing — kept as a migration record |

All of them exit non-zero on failure and, except for the explicitly marked
modes below, only read the repository.

## apiinventory

Parses `h3api.h.in` from the pristine upstream tree under `testref/` and the
root package's `// Ported from H3 C: <file>::<name>` attribution comments,
then cross-references the two.

- **Default mode** writes the CSV published as
  [docs/c-api-inventory.csv](../docs/c-api-inventory.csv) to stdout
  (columns: `c_function,c_signature,go_file,go_func,go_signature,is_c_public,notes`)
  and a summary to stderr.
- **`-verify`** exits 1 unless every C public function is ported *and*
  referenced by an `H3 C API:` doc line (or listed in the reviewed omissions
  table inside `main.go`).
- Flags: `-repo`, `-h3ver` (default 4.4.0), `-header` (explicit header
  path), `-verify`. Requires `make -C testref h3-source` first.
- CI keeps the committed CSV current: the `api-gates` job regenerates it and
  fails on `git diff`.

## testinventory

Discovers every upstream test-ecosystem entry — named `TEST(...)` cases, CLI
registrations, fuzzers, benchmarks, filters, helpers, support sources,
fixtures, and build definitions — fingerprints each with SHA-256, and checks
the reviewed registry
[docs/upstream-test-inventory.csv](../docs/upstream-test-inventory.csv).

- **Default mode** prints a report (counts by kind, dispositions, unreviewed
  cases, stale rows, integrity problems).
- **`-verify`** exits 1 on any unreviewed/stale/missing/invalid row;
  referenced Go tests must actually exist.
- **`-init`** prints skeleton CSV rows for unreviewed cases (used when
  adopting a new upstream release).
- Flags: `-h3ver`, `-upstream`, `-repo`, `-registry`, `-verify`, `-init`.
  Read-only; never edits Go code or the registry.

## cliinventory

Extracts the CLI contract from the upstream tree (registered subcommands in
`h3.c`, every `add_h3_cli_test(...)` scenario, referenced fixtures, defining
sources) and checks it against the four committed registries
(`docs/cli-*.csv`).

- **`-verify`** exits 1 on any command/scenario/fixture/source drift,
  including hash changes of the defining sources.
- **`-emit-cases` / `-emit-fixtures` / `-emit-sources`** print discovered
  CSVs to stdout (used to seed or refresh the registries during review).
- **`-update-ecosystem-inventory` / `-update-contract-metadata`** are the
  two *file-modifying* modes: they rewrite
  `docs/upstream-test-inventory.csv` / `docs/cli-contract.csv` in place
  (used once per upstream sync, then reviewed in the diff).
- Flags: `-upstream`, `-registry`, `-contract`, `-fixtures`, `-sources`,
  plus the mode flags above.

## upstreamdiff

Compares two upstream H3 source trees at the *symbol* level (functions,
tables, macros, types — not files) and maps every changed symbol to the Go
port via attribution comments. This is the mandatory first review step of an
upstream sync: an API-surface check alone cannot prove implementation
equivalence.

- Output: a Markdown report (changed/added/removed symbols with their Go
  files, plus changed upstream test files) — commit the reviewed result as
  `docs/sync/<old>-to-<new>.md`.
- **`-strict`** exits 1 if any changed library symbol has no Go mapping.
- Flags: `-from`, `-to` (required), `-repo`, `-ported-tests`, `-strict`.
  Both trees must exist under `testref/`
  (`make -C testref H3_VERSION=<ver> h3-source`).

## docscheck

Checks all Markdown files (excluding downloaded/generated trees): relative
links must resolve to existing files/directories, and `#fragment` links into
Markdown files must match a GitHub-generated heading anchor. Fenced code
blocks and inline code are ignored. Run with `make check-docs`; CI runs it
on every push/PR, including docs-only changes.

## ubercompare

Maintains the comparison matrix against the official cgo binding. The
curated data lives in
[docs/comparison-uber-h3-go.csv](../docs/comparison-uber-h3-go.csv) (one
row per public C function); the tool renders it into the generated tables
of [docs/comparison-uber-h3-go.md](../docs/comparison-uber-h3-go.md) and
cross-checks it — entirely offline — against the committed inventories.

- **Default mode** prints the generated Markdown tables to stdout.
- **`-write`** rewrites the marked generated section of the comparison doc
  (`make gen-ubercompare`) — an explicitly file-modifying mode.
- **`-verify`** exits 1 on drift (`make check-ubercompare`, run by the CI
  `docs` job): every matrix row must be a public C function from
  [docs/c-api-inventory.csv](../docs/c-api-inventory.csv) and vice versa;
  every `this_api` symbol must exist in
  [docs/api-surface.txt](../docs/api-surface.txt); statuses must come from
  the fixed vocabulary; the doc tables must match the CSV.
- Flags: `-repo`, `-write`, `-verify`.
- The binding's side of the matrix needs the uber/h3-go dependency and is
  verified by `TestMappingSymbolsExist` in
  [interop/uberbench](../interop/uberbench/README.md) instead.

## layoutinventory

Classifies all root source files (`*.go`, `h3lib_*.c`) into the nine
architectural layers of the flat `h3` package — public API, ported
implementation, ported public types, the three parity-harness kinds, and
the three test kinds — using file evidence only (package clause, build
tags, attribution comments, exported declarations, filename shape). The
layer taxonomy and its rationale live in
[docs/repository-layout-review.md](../docs/repository-layout-review.md)
(DR-008).

- **Default mode** writes the CSV published as
  [docs/file-layer-inventory.csv](../docs/file-layer-inventory.csv) to
  stdout (columns:
  `file,layer,package,build_tags,attributions,h3_c_api_refs,exported_decls`)
  and a per-layer summary to stderr.
- **`-verify`** exits 1 if any root source file matches no layer rule —
  the guard that keeps the classification from rotting as files are added.
- Flags: `-repo`, `-verify`. Runs without `testref/`.

## unexport (historical)

The one-time mechanical sweep that unexported the accidentally exported
C-style identifiers during the public-API build-out
(docs/public-api-architecture.md §7, Phase 2). It is **not wired into any
Makefile target or workflow** and exists only as the reviewable record of
that migration (`docs/DEVIATIONS.md` item 11 references it). Dry run by
default; `-apply` rewrites the root `*.go` files and is not something you
should need again.
