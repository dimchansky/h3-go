# h3 (pure Go)

**Status:** WIP — goal is a **pure-Go** reimplementation of Uber’s H3 (tag **v4.3.0**) with Go-first, allocation-aware APIs. No cgo, no external dependencies.

- Reference impl (ground truth): https://github.com/uber/h3 (v4.3.0)
- API docs: https://h3geo.org/docs/api/indexing
- Go bindings for inspiration (do not depend on): https://github.com/uber/h3-go

## Design goals
- Behavioral equivalence with H3 C v4.3.0.
- **Zero/low allocations** via `dst []T` buffer parameters across collection-returning APIs.
- Deterministic outputs (documented order).
- Concurrency-safe: no mutable global state post-init.
- Readable, well-documented Go with tests and benchmarks.

## API patterns
Functions that return slices accept an optional destination buffer:
```go
// Reuses dst capacity when sufficient; may allocate otherwise.
cells, err := KRing(buf[:0], origin, k)
```

## Roadmap
See [TODO.md](./TODO.md) for the live plan, dependency breakdown, and acceptance criteria.

## Testing strategy
- Table-driven tests mirroring H3 C behavior.
- External **C oracle CLI** (built from H3 v4.3.0) invoked by Go tests — keeps this module pure Go while validating against the reference.
- Fuzz tests for reversible transforms; microbenchmarks for hot paths.

## C-to-Go Conversion (internal/c2go)
To make future ports easier to audit against the C source, there is a dedicated workspace under `internal/c2go` for near line-by-line conversions of H3 C functions.

- Structure:
  - One Go file per function: `<cfile>__<function>.go` (unexported, original name where possible).
  - One cgo interop file per C module: `<cfile>_cgo.go` — includes the original C file by name and exposes small C wrappers for parity tests.
  - Parity tests: `<cfile>__<function>_parity_test.go` with build tag `c2go` (compares Go vs C). Reserve `<cfile>__<function>_test.go` for plain-Go tests later.

- No version pin in code: the cgo interop includes headers/sources by name only. Include directories are passed at build time via `CGO_CPPFLAGS` from the Makefile.

- How to run parity tests locally:
  1) Build the H3 C sources used as reference (once): `make ref`
  2) Run c2go tests (requires cgo toolchain):
     - `make test-c2go` (uses `H3VER=4.3.0` by default)
     - `make H3VER=4.4.0 test-c2go` to test against a different checked-out version under `testref/`

- Porting algorithm:
  - Pick a function from `testref/h3-<ver>/src/h3lib/lib/<cfile>.c` with minimal dependencies.
  - Create `internal/c2go/<cfile>__<function>.go` with a faithful Go translation and the same unexported name.
  - If it depends on other C helpers, add a `TODO:` with the exact C symbol names, then port those recursively.
  - In `internal/c2go/<cfile>_cgo.go`, add tiny C wrappers that call the original C functions; keep C wrapper names distinct (e.g., `_ipow_c_wrapper`) and expose Go helpers with a `C` suffix (e.g., `_ipowC`).
  - Add `internal/c2go/<cfile>__<function>_parity_test.go` with `//go:build c2go` that compares Go vs the C wrapper across representative inputs. Keep `<cfile>__<function>_test.go` for future pure-Go tests.
  - Run `make test-c2go` (or `H3VER=... make test-c2go`).

Notes:
- Library code remains pure Go; cgo is used only inside `internal/c2go/*_cgo.go` and only compiled with the `c2go` tag.
- Do not modify anything under `testref/` except via `make ref` (which downloads/builds the C sources). Version is controlled by `H3VER` at test time.

### Oracle-backed tests
- Build oracle once: `make ref` (produces `testref/h3ref`).
- Run parity tests: `make test-oracle` (or `make test-all` for both normal and oracle).
- Controls (env):
  - `ORACLE_MAX`: cap randomized/exhaustive cases in parity tests (default `200`).
  - `ORACLE_SEED`: seed for randomized generators (default `1337`; set `0` for time-based).
  - `ORACLE_PATH`: optional absolute path to `h3ref` binary (defaults to `./testref/h3ref`).

Examples:
```bash
# Quick run with defaults
make test-oracle

# Heavier sweep locally
ORACLE_MAX=2000 ORACLE_SEED=42 make test-oracle

# Use a custom-built oracle binary
ORACLE_PATH=/abs/path/to/h3ref make test-oracle

### c2go parity tests
```bash
# Build the reference C sources once (downloads to testref/)
make ref

# Run c2go tests against default H3 version (4.3.0)
make test-c2go

# Override H3 version if another source tree exists under testref/
make H3VER=4.4.0 test-c2go
```
```

## Development
- Go ≥ 1.22 recommended.
- Linting: `go vet`, `golangci-lint` (incl. `staticcheck`). 
- CI runs build, tests (with `-race` for tests), lints, and benches (as configured).

## License & attribution
- Apache-2.0.
- This work reimplements algorithms from Uber's H3; see [NOTICE](./NOTICE) for attribution.
