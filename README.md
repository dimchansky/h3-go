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
```

## Development
- Go ≥ 1.22 recommended.
- Linting: `go vet`, `golangci-lint` (incl. `staticcheck`). 
- CI runs build, tests (with `-race` for tests), lints, and benches (as configured).

## License & attribution
- Apache-2.0.
- This work reimplements algorithms from Uber’s H3; see `NOTICE` for attribution.
