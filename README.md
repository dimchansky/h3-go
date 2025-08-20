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
- External **C oracle CLI** (built from H3 v4.3.0) invoked by Go tests via `exec.Command` — keeps this module pure Go while validating against the reference.
- Fuzz tests for reversible transforms; microbenchmarks for hot paths.

## Development
- Go ≥ 1.22 recommended.
- Linting: `go vet`, `golangci-lint` (incl. `staticcheck`). 
- CI runs build, tests (with `-race` for tests), lints, and benches (as configured).

## License & attribution
- Apache-2.0.
- This work reimplements algorithms from Uber’s H3; see `NOTICE` for attribution.
