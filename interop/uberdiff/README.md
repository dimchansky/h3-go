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
require github.com/uber/h3-go/v4 v4.5.0
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
degrees absolute and areas `1e-8` relative. Both sides run the same H3 C
v4.5.0 algorithm (the pinned binding vendors H3 C v4.5.0, matching this
library's parity target), but `CellAreaKm2` is compared end to end: each
library derives the cell boundary through its own trig — pure Go `math`
here, platform libm in the cgo binding — then sums per-edge Cagnoli terms,
so the tiny boundary-vertex differences are amplified in the area, and the
more so as cell areas shrink with resolution. Measured across the full
deterministic input set on **linux/amd64 and darwin/arm64**, the res-8
cells this test compares reach at most **~2.1e-9 relative**; the `1e-8`
tolerance is a ~4.7x evidence-based margin over that maximum. It is **not**
a claim of bit-exact area equality — exact area equality is not portable
across libm/compiler implementations. (An earlier note quoted ~2e-12,
which was an unrepresentative pentagon-only, macOS-only subset, not the
res-8 comparison this test performs.)

## What it proves — and what it does not

This suite demonstrates **drop-in agreement with the ecosystem binding** on
common operations over randomized global input. It is *not* the correctness
anchor: function-by-function equivalence with H3 C v4.5.0 (including error
codes and edge cases) is established by the cgo parity suite in the root
module (`make test-c2go`, 213 test files against pristine upstream sources).
That suite compares against C **compiled from the same sources on the same
platform**, so most results match bit-for-bit — but **cell areas remain
tolerance-based even there** (an absolute km² tolerance in
`latLng__cellAreaKm2_parity_test.go`, a ~1e-14 relative one in
`area_geoLoopAreaRads2_parity_test.go`), because floating-point area is not
bit-exact across compilers/libms. Treat uberdiff as an independent second
witness with broad but shallower coverage; its looser area tolerance
reflects comparison against a *separately built* binding, end to end.

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
