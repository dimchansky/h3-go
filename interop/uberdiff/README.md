# uberdiff — differential tests against the official cgo binding

This directory is a **separate Go module** whose tests cross-check the
pure-Go `h3` library at the repository root against
[uber/h3-go](https://github.com/uber/h3-go), the official Go binding that
wraps the H3 C library through cgo. If the two implementations disagree on
any compared operation, a real user migrating between them would see the
difference — that is exactly what this suite guards.

## Why a separate module

The root library has zero dependencies, and keeping it that way is a hard
project invariant. Pulling `uber/h3-go` into the root `go.mod` — even as a
test dependency — would break that, so the differential tests live in their
own module with a `replace` directive pointing at the parent:

```
require github.com/uber/h3-go/v4 v4.4.1
replace github.com/dimchansky/h3-go => ../..
```

Because `uber/h3-go` is a cgo binding, running this suite needs a **C
toolchain** and (first run) **network access** to fetch the binding. Nothing
here is ever part of a consumer's build.

## What is compared

Deterministic pseudo-random coordinates (fixed PCG seed, plus fixed
near-polar and near-antimeridian points) drive both implementations through:

| Test | Operations |
|---|---|
| `TestLatLngToCellParity` | `LatLngToCell` at res 0/4/9/15 — exact index equality |
| `TestCellRoundTripParity` | cell centroid, boundary vertexes, `IsPentagon`/`Resolution`/`BaseCellNumber`/`IsValid` |
| `TestHierarchyParity` | `Parent`, `Children` — exact |
| `TestGridDiskParity` | `GridDisk` k=1,3 — set equality |
| `TestCompactParity` | `CompactCells` over children sets — exact |
| `TestMetricsParity` | `CellAreaKm2` — relative tolerance |

Index-valued results must match **exactly**. Coordinates allow `1e-10`
degrees absolute and areas `1e-9` relative: even at the same upstream
release (the pinned binding vendors H3 C v4.4.1, which differs from this
library's v4.4.0 parity target only by a version-metadata fix), the C
compiler may contract floating-point multiply-adds differently than the Go
compiler, and area computation amplifies the last-ulp differences to
~1e-10–1e-9 relative near pentagons.

## What it proves — and what it does not

This suite demonstrates **drop-in agreement with the ecosystem binding** on
common operations over randomized global input. It is *not* the correctness
anchor: exact, function-by-function equivalence with H3 C v4.4.0 (including
error codes and edge cases) is established by the cgo parity suite in the
root module (`make test-c2go`, 227 test files against pristine upstream
sources). Treat uberdiff as an independent second witness with broad but
shallower coverage.

## Running it

```sh
make test-uberdiff        # from the repository root
# equivalent to: cd interop/uberdiff && go test ./...
```

In CI it runs in the **nightly** workflow (schedule, manual dispatch, and
`v*` release tags) — not on every push/PR, since it needs cgo and an
external dependency; see [docs/ci-policy.md](../../docs/ci-policy.md).

## Updating for a new uber/h3-go release

1. Bump the version in [go.mod](go.mod) (`go get github.com/uber/h3-go/v4@<ver>`
   in this directory, then `go mod tidy`).
2. Run `make test-uberdiff`. If a numeric comparison starts failing, check
   which upstream H3 C release the new binding vendors (its `H3_VERSION`
   file) against this library's parity target; behavioral skew between
   releases needs a tolerance rationale here, not a silent bump.
   Keep the version in step with [interop/uberbench](../uberbench/README.md),
   which benchmarks against the same binding.
3. There are no fixtures or generated files here — inputs are generated at
   run time from a fixed seed.
