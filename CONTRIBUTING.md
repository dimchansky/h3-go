# Contributing

Thanks for your interest! This project has an unusual structure — a public
Go API layered over a mechanically ported C implementation — and a few hard
rules that keep it maintainable against upstream H3 releases. Reading this
page first will save you review round-trips.

## Ground rules (CI-enforced)

1. **Production code is pure, safe Go.** No `unsafe`, no cgo, no
   dependencies in any file selected by a normal build. cgo exists only in
   the parity harness behind `//go:build cgo && c2go`. `make check-unsafe`
   is the gate; introducing production `unsafe` requires a new reviewed
   decision record (see DR-007 in
   [docs/public-api-architecture.md](docs/public-api-architecture.md)).
2. **Ported code stays traceable to C.** Files named
   `<cfile>_<function>.go` mirror one C function each (a double underscore —
   `<cfile>__<name>.go` — marks a C `_`-prefixed static helper) and keep C
   names, C-shaped bodies, and a `// Ported from H3 C: <file>::<name>`
   attribution. Do not restructure them for style — style-tier lint checks
   are excluded for exactly these files, and only these; the rationale and
   scope are documented in [docs/lint-policy.md](docs/lint-policy.md).
   C `int` maps to Go `int32`, C `int64_t` to `int64` (overflow parity).
3. **Public wrappers carry an `H3 C API: <name>` doc line** naming their C
   counterpart. `make check-api` fails if a C public function is neither
   referenced nor listed in the omissions table
   (`tools/apiinventory/main.go`).
4. **The exported surface is locked.** Any intentional API change must
   regenerate the golden file:
   `UPDATE_API_SURFACE=1 go test -run TestAPISurface .`
5. **Intentional divergences from C live in
   [docs/DEVIATIONS.md](docs/DEVIATIONS.md).** If your change makes Go
   behavior differ from C on purpose, document it there; if it isn't listed
   there, parity with C is a requirement.
6. **Every root file belongs to exactly one architectural layer.** The
   flat root package deliberately mixes the public API and the ported layer
   (a package split is impossible without breaking the zero-copy alias and
   method architecture — DR-008 in
   [docs/repository-layout-review.md](docs/repository-layout-review.md)),
   so layer identity is enforced by naming: public API files use topical
   underscore-free names (`cell.go`, `traversal.go`); ported files use
   `<cfile>_<name>.go`; the parity harness uses `*_cgo.go` / `h3lib_*.c` /
   `*_parity_test.go`; test files follow the layer of what they exercise.
   `make check-layout` fails on any file that matches no rule, and the CI
   `api-gates` job fails if the generated per-file map
   ([docs/file-layer-inventory.csv](docs/file-layer-inventory.csv),
   `make layout-inventory`) is stale.

## Development workflow

```sh
make test                    # pure-Go tests (CGO_ENABLED=0) — needs only Go
go test -race ./...          # race detector
make lint                    # gofmt -s, go vet, golangci-lint, smrcptr
make check-unsafe            # no-unsafe gate
make check-layout            # file-layer gate (+ make layout-inventory to regenerate)
make check-docs              # Markdown link/anchor gate
make bench                   # benchmarks with allocation stats
```

Parity validation (optional locally, mandatory in CI; needs a C toolchain
and network):

```sh
make -C testref h3-source    # download the upstream H3 sources (never vendored)
make test-c2go               # full parity suite vs the original C objects
make test-c2go TEST=Test_gridDisk_parity   # single parity test
make check-api               # C-API completeness gate
make api-inventory           # regenerate docs/c-api-inventory.csv
make check-test-inventory    # case-level upstream test completeness gate
make test-upstream-fixtures  # 526,546 golden conversion/boundary records
make test-cli               # all 170 upstream CLI scenarios
make test-cli-process       # actual binary, pipes, stderr, exit status
make check-cli-inventory    # CLI semantic/source/fixture drift gate
make test-cli-diff          # build and compare with upstream C h3_bin
```

Differential testing and benchmarking against the official cgo binding:

```sh
make test-uberdiff           # separate module in interop/uberdiff
make test-uberbench          # benchmark-pairing equivalence (interop/uberbench)
make bench-uber              # full comparative benchmark + memory suite
make gen-ubercompare         # regenerate the comparison-matrix tables
make check-ubercompare       # matrix drift gate (runs in CI's docs job)
```

The comparison matrix, migration guide, and benchmark artifacts are
documented in [docs/README.md](docs/README.md#comparison-with-the-official-cgo-binding).

## Porting a C function (upstream syncs)

The full workflow is
[docs/public-api-architecture.md §10](docs/public-api-architecture.md#10-upstream-synchronization-workflow);
every tool referenced below is documented in [tools/README.md](tools/README.md),
and [docs/README.md](docs/README.md) indexes the registries they maintain.
The short version:

1. Fetch the new version: `make -C testref H3_VERSION=<ver> h3-source`.
2. Run the **symbol-level diff** — this step is mandatory, an API check alone
   does not prove implementation equivalence:
   `make upstream-diff FROM=<old> TO=<ver>`. Review every changed/added/
   removed symbol AND every changed upstream test file to a recorded
   disposition in `docs/sync/<old>-to-<ver>.md`.
   Then `go run ./tools/apiinventory -h3ver <ver> -verify` for the public-API
   completeness view.
3. Port each changed C function in its own file, preserving the attribution
   comment. Never hardcode an H3 version in code — include paths come from
   `make test-c2go` (`H3VER=` selects the tree).
4. Add/extend a `<cfile>_cgo.go` wrapper and a `*_parity_test.go`
   (both `//go:build cgo && c2go`), then `make test-c2go H3VER=<ver>`.
5. Audit every changed upstream test-ecosystem entry, port meaningful
   behavior, update `docs/upstream-test-inventory.csv`, and run
   `make check-test-inventory`. The scope includes named `TEST` cases,
   input-driven executables, CLI registrations, fuzzers, benchmarks,
   filters, support sources, fixtures, and build definitions; see
   [docs/ported-c-tests.md](docs/ported-c-tests.md).
   Run `make check-cli-inventory H3VER=<ver>` for every upgrade; changes to
   CMake target naming, parser sources, `h3.c`, CLI tests, or fixtures require
   reviewed updates to the inventories in `docs/cli-*.csv`, followed by
   `make test-cli-diff H3VER=<ver>`.
6. Add the public wrapper with its `H3 C API:` line, tests (including an
   allocation assertion if it returns collections), and regenerate the
   inventory + API surface.

## Pull requests

- Keep commits focused; describe *what* changed and *why* (see `git log`
  for the house style — prefixes by area: `api:`, `c2go:`, `cli:`, `tests:`,
  `tools:`, `docs:`, `ci:`). Keep messages free of tool-attribution noise
  (no `Co-Authored-By:` bot trailers or AI-assistant references).
- Run `make fmt lint test check-unsafe` before pushing; run the parity suite
  if you touched ported code or the harness. CI is tiered
  ([docs/ci-policy.md](docs/ci-policy.md)): docs-only changes skip Go jobs,
  and the race detector runs on PRs/nightly/tags rather than every push.
- Benchmark deltas are expected for performance-related changes
  (`make bench`), and allocation assertions must keep passing — new
  convenience APIs must not add allocations to existing paths.
