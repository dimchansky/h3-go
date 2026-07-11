# h3-go — pure Go port of Uber H3

A **pure-Go** implementation of [Uber's H3](https://github.com/uber/h3) hexagonal
hierarchical geospatial indexing system, behaviorally equivalent to **H3 C v4.4.0**.
No cgo, no external dependencies, and the production library is **safe Go only**
(no `unsafe` — enforced in CI).

All **78/78** public functions of H3 C 4.4.0 are covered: ported
function-by-function, validated against the original C objects by a 227-file cgo
parity suite, and exposed through an idiomatic, strongly typed Go API.

```go
import h3 "github.com/dimchansky/h3-go"

cell, _ := h3.LatLngToCell(h3.LatLngDegs(37.7759, -122.4180), 9)
fmt.Println(cell)                  // 8928308280fffff

disk, _ := cell.GridDisk(1)        // the cell and its 6 neighbors
center, _ := cell.LatLng()         // cell centroid
area, _ := cell.AreaKm2()          // exact spherical area

// Zero-allocation form: reuse a buffer across queries.
buf := make([]h3.Cell, 0, 64)
disk, _ = cell.AppendGridDisk(buf[:0], 1)

// Iterators stream without materializing (0 allocs).
for child := range cell.ChildrenSeq(12) { _ = child }
```

## API shape

- **Typed indexes**: `Cell`, `DirectedEdge`, `Vertex` (distinct `uint64` types)
  with `String`/`ParseCell`/`MarshalText` (hex, JSON-ready).
- **Type-safe angles**: `LatLng` fields are `Angle` values — construct with
  `h3.LatLngDegs(...)` or `h3.Deg`/`h3.Rad`; degree/radian mix-ups don't compile.
- **Errors**: sentinel values (`h3.ErrPentagon`, `h3.ErrCellInvalid`, ...) matched
  with `errors.Is`, mirroring the C error codes.
- **Allocation-aware**: every collection API has an allocating convenience form
  and a zero-allocation `Append*` form; `iter.Seq` iterators for streaming.

See the [package documentation](https://pkg.go.dev/github.com/dimchansky/h3-go)
and [docs/public-api-architecture.md](docs/public-api-architecture.md) (design,
decision records, measurements).

## Layout

- **Root package `h3`** — public API files (`indexing.go`, `traversal.go`, …)
  over a one-file-per-C-function ported implementation
  (`<cfile>_<function>.go`, each with a `// Ported from H3 C: <file>::<name>`
  attribution). Public operations carry `H3 C API: <name>` doc lines.
- **Parity harness** — `*_cgo.go` + `h3lib_*_c2go.c` behind
  `//go:build cgo && c2go` compile the *original* upstream C sources
  (downloaded to `testref/`, never vendored) and compare Go vs C in-process.
- `tools/apiinventory` — generates [docs/c-api-inventory.csv](docs/c-api-inventory.csv)
  and enforces C-API completeness (`make check-api`).
- [docs/DEVIATIONS.md](docs/DEVIATIONS.md) — the intentional differences from C.
- [docs/FUTURE_WORK.md](docs/FUTURE_WORK.md) — the backlog: deferred features
  (GeoJSON, grouped disk distances, polyfill workspaces), profiling-gated
  ideas, rejected designs, and the pre-v1.0.0 checklist.

## Development

```sh
make test              # pure-Go tests (CGO_ENABLED=0)
go test -race ./...    # race detector
make -C testref h3-source && make test-c2go   # cgo parity tests vs original C
make lint              # gofmt + go vet + golangci-lint + smrcptr
make check-unsafe      # DR-007 gate: no unsafe reachable from any normal build
make check-api         # every C public function ported & publicly represented
make bench             # benchmarks with allocation stats
```

Requires Go ≥ 1.24. The parity suite additionally needs a C toolchain and
network access (to fetch upstream sources).

## License & attribution

Apache-2.0. This project reimplements algorithms from Uber's H3; see
[NOTICE](./NOTICE) for attribution.
