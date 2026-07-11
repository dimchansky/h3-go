# h3-go — pure Go port of Uber H3

A **pure-Go** reimplementation of [Uber's H3](https://github.com/uber/h3) hexagonal
hierarchical geospatial indexing system (reference version **v4.3.0**). No cgo, no
external dependencies; the production library is **safe Go only** (no `unsafe`).

**Status:** the C implementation layer is complete — all **75/75** public functions of
H3 C 4.3.0 are ported and parity-tested against the original C objects. The idiomatic
public Go API (`Cell`, `DirectedEdge`, `Vertex`, …) is being built per
[docs/public-api-architecture.md](docs/public-api-architecture.md).

## Layout

- **Root package `h3`** — one Go file per ported C function
  (`<cfile>_<function>.go`, `__` marks a C-internal `_`-prefixed helper), each carrying a
  `// Ported from H3 C: <file>::<name>` attribution. Public API files use plain topical
  names (`cell.go`, `traversal.go`, …).
- **Parity harness** — `*_cgo.go` interop wrappers + `h3lib_*_c2go.c` shims, all behind
  `//go:build cgo && c2go`, compile the *original* upstream C sources (downloaded to
  `testref/`, never vendored) and compare Go vs C behavior in-process.
- `tools/apiinventory` — generates [docs/c-api-inventory.csv](docs/c-api-inventory.csv),
  the mechanical C↔Go function mapping used for upstream synchronization.

## Development

```sh
make test              # pure-Go tests (CGO_ENABLED=0)
go test -race ./...    # race detector
make ref               # one-time: download & build the C reference (testref/)
make test-c2go         # cgo parity tests vs original C (H3VER=4.3.0)
make lint              # gofmt + go vet + golangci-lint + smrcptr
make check-unsafe      # DR-007 gate: no unsafe reachable from any normal build
make bench             # benchmarks with allocation stats
```

Requires Go ≥ 1.24. The parity suite additionally needs a C toolchain and network (to
fetch upstream sources).

## Design in one paragraph

The mechanically ported layer keeps C names, C signatures (out-params, caller-sized
buffers), and C integer semantics (`int32` for C `int`) so future upstream releases can
be diffed and merged function-by-function. The public layer exposes defined types
`Cell`/`DirectedEdge`/`Vertex` (`uint64`), with the internal `H3Index` an *alias* of
`Cell` — so `[]Cell` and the ported code's index slices are the same type and every hot
path is zero-copy without `unsafe` or generics. Collection APIs come in an allocating
convenience form and a zero-allocation `Append*` form, plus `iter.Seq` iterators.
Angles use the `Angle` type (radians inside, `Deg()`/`Rad()` accessors) so
degree-vs-radian bugs cannot compile. See the
[architecture document](docs/public-api-architecture.md) for the full rationale,
decision records, and measurements.

## License & attribution

Apache-2.0. This project reimplements algorithms from Uber's H3; see
[NOTICE](./NOTICE) for attribution.
