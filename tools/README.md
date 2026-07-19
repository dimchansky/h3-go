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
| [doclinkcheck](doclinkcheck) | Verify `[Symbol]` doc links in the root package's GoDoc resolve | `make check-docs` | CI (`docs` job) |
| [benchdocs](benchdocs) | Generate and verify the README scorecard and complete benchmark comparison from committed artifacts | `make gen-benchdocs` / `make check-benchdocs` | CI (`docs` job), benchmark refreshes |
| [ubercompare](ubercompare) | Generate and verify the uber/h3-go comparison-matrix tables and the C→Go API map | `make gen-ubercompare` / `make check-ubercompare` | CI (`docs` job), binding/H3 release updates |
| [layoutinventory](layoutinventory) | Classify every root source file into its architectural layer; verify none is unclassifiable | `make layout-inventory` / `make check-layout` | CI (`fast`, `api-gates`), layout discoverability ([docs/repository-layout-review.md](../docs/repository-layout-review.md)) |
| [cirequired](cirequired) | Evaluate the final CI job results against the required-check truth table | `go run ./tools/cirequired` (CI-only) | CI (`required` job — the "CI / required" aggregate merge gate) |
| [releasepack](releasepack) | Single authoritative release build: invariant checks, hermetic cross-builds, deterministic archives, SHA256SUMS | `make release-dist VERSION=vX.Y.Z OUT=<dir>` | release-builds workflow (`build` + `verify-reproducible`), the local release procedure (docs/releasing.md) |
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
- Flags: `-repo`, `-h3ver` (default 4.5.0), `-header` (explicit header
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
- **`-comments`** additionally diffs every extracted symbol's *leading*
  documentation comment as an independent dimension (public functions,
  internal helpers, tables/constants, types, macros). The symbol table
  gains a Change column (`body` / `comment-only` / `body+comment`) and a
  derived "Leading-comment-only changes" section lists symbols whose body
  is unchanged. Default mode ignores leading comments entirely (its output
  is unchanged by this flag's existence); comments *inside* bodies always
  count as body changes. This is the documentation-drift review step of
  the sync workflow (public GoDoc mirrors upstream contract comments) —
  see CONTRIBUTING.md. `-strict` is unaffected by `-comments`.
- Flags: `-from`, `-to` (required), `-repo`, `-ported-tests`, `-strict`,
  `-comments`. Both trees must exist under `testref/`
  (`make -C testref H3_VERSION=<ver> h3-source`).

## docscheck

Checks all Markdown files (excluding downloaded/generated trees): relative
links must resolve to existing files/directories, and `#fragment` links into
Markdown files must match a GitHub-generated heading anchor. Fenced code
blocks and inline code are ignored. Run with `make check-docs`; CI runs it
on every push/PR, including docs-only changes.

## doclinkcheck

Checks that every `[Symbol]` / `[Type.Method]` doc-link candidate in the
root package's doc comments resolves to a declared package symbol.
`go/doc/comment` leaves unresolvable candidates as plain text (no `DocLink`
node is emitted), so the tool scans the syntactic candidates itself:
bracketed exported identifiers in package, declaration, spec, and
struct-field doc comments. Candidates not starting with an uppercase letter
(numeric ranges, GeoJSON `[lng, lat]`) and package-qualified references are
ignored; test files are skipped. `-dir` selects another package directory
(used by its own fixture test). Run with `make check-docs`; CI runs it on
every push/PR, including docs-only changes.

## ubercompare

Maintains the comparison matrix against the official cgo binding and the
simplified C→Go API map projected from it. The curated data lives in
[docs/comparison-uber-h3-go.csv](../docs/comparison-uber-h3-go.csv) (one
row per public C function); the tool renders it into the generated tables
of [docs/comparison-uber-h3-go.md](../docs/comparison-uber-h3-go.md) and
[docs/api-map.md](../docs/api-map.md) (C function → idiomatic Go API →
additive `Append*`/`*Seq`/grouped forms) and cross-checks it — entirely
offline — against the committed inventories.

- **Default mode** prints the generated Markdown tables to stdout.
- **`-write`** rewrites the marked generated sections of both documents
  (`make gen-ubercompare`) — an explicitly file-modifying mode.
- **`-verify`** exits 1 on drift (`make check-ubercompare`, run by the CI
  `docs` job): every matrix row must be a public C function from
  [docs/c-api-inventory.csv](../docs/c-api-inventory.csv) and vice versa;
  every `this_api` symbol must exist in
  [docs/api-surface.txt](../docs/api-surface.txt); statuses must come from
  the fixed vocabulary; both documents' generated tables must match the
  CSV.
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
  and a per-layer summary to stderr (`make layout-inventory`).
- **`-verify`** exits 1 if any root source file matches no layer rule —
  the guard that keeps the classification from rotting as files are added
  (`make check-layout`, run by the CI `fast` job; the `api-gates` job also
  fails if the committed CSV is stale).
- Flags: `-repo`, `-verify`. Runs without `testref/`.

## cirequired

The logic behind the `CI / required` aggregate status check — the single
context branch protection requires. The `required` job in
[.github/workflows/ci.yml](../.github/workflows/ci.yml) runs with
`if: always()`, passes the final `needs` results plus the event name and the
docs-only classifier verdict through environment variables
(`CIREQUIRED_NEEDS`/`CIREQUIRED_EVENT`/`CIREQUIRED_CODE`), and delegates the
pass/fail decision to this tool.

- Passes only when the results match the truth-table row for the event and
  classifier verdict exactly: docs-only changes must *skip* the Go jobs,
  code changes must *run* them (`race` is a PR-only merge gate, so it is
  expected `skipped` on pushes).
- Fails on unknown events, invalid classifier values, malformed or
  duplicate-key JSON, unknown/missing jobs, and any failure, cancellation,
  or unexpected skip — a broken job condition fails the gate instead of
  slipping a change past it. There is deliberately no generic
  "success or skipped" acceptance.
- The gated job set is defined once in `main.go` (`gatedJobs`); change it
  and the workflow's `needs:` list together.
- Table-driven tests cover every truth-table row and every rejection path
  (`go test ./tools/cirequired`).

## releasepack

The single authoritative release-build implementation behind
`make release-dist VERSION=<tag> OUT=<empty dir>`. The release-builds
workflow's `build` job, its independent `verify-reproducible` job, and the
local release procedure (docs/releasing.md) all run exactly this code, and
their archives must be byte-identical.

- **Preconditions (reject, don't warn):** canonical tag syntax (`vX.Y.Z` or
  `vX.Y.Z-rc.N` per [docs/versioning.md](../docs/versioning.md)); clean
  worktree and index; the tag exists and points at `HEAD`; `go env
  GOVERSION` equals the pinned release toolchain (`requiredGoVersion` in
  `main.go`, kept in lockstep with the workflow's `RELEASE_GO_VERSION`);
  empty output directory (pass one **outside** the repository).
- **Hermetic builds:** six `CGO_ENABLED=0 go build -trimpath` cross-builds
  with `-ldflags "-s -w -buildid= -X …internal/cli.buildVersion=<tag>"`;
  `GOTOOLCHAIN=local`, `GOENV=off`, `GOWORK=off`, empty
  `GOFLAGS`/`GOEXPERIMENT`, `TZ=UTC`, `LC_ALL=C`, and per-target
  `GOAMD64=v1`/`GOARM64=v8.0` are set explicitly — host settings cannot
  leak into the binaries. `SOURCE_DATE_EPOCH` derives from the tagged
  commit, never from the environment.
- **Postconditions:** every binary's `go version -m` must report the module
  path, the tagged commit, `vcs.modified=false`, the target GOOS/GOARCH,
  and no host paths; the host-runnable binary is executed and must print
  `h3 4.5.0 (<tag>)`.
- **Deterministic archives:** sorted entries, uid/gid 0, modes 0755/0644,
  all mtimes = `SOURCE_DATE_EPOCH`, zero gzip MTIME/name, extra-field-free
  zip entries with fixed DOS timestamps; plus a sha256sum-compatible
  `SHA256SUMS`. Structural tests inspect headers, ordering, ownership,
  timestamps, and manifest contents, and prove run-twice byte-identity
  (`go test ./tools/releasepack`).

## unexport (historical)

The one-time mechanical sweep that unexported the accidentally exported
C-style identifiers during the public-API build-out
(docs/public-api-architecture.md §7, Phase 2). It is **not wired into any
Makefile target or workflow** and exists only as the reviewable record of
that migration (`docs/DEVIATIONS.md` item 11 references it). Dry run by
default; `-apply` rewrites the root `*.go` files and is not something you
should need again.
